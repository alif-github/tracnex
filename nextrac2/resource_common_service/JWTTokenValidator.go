package resource_common_service

import (
	"database/sql"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/bukalapak/go-redis"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/token"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

func ValidateJWTToken(client *redis.Client, jwtTokenStr string, expiredIn int64, checkTokenEndpoint string, resourceID string, scope string, _ string,
	roleMapping func(clientID string, token string, payload token.PayloadJWTToken) (model.AuthAccessTokenModel, model.AuthenticationRoleModel, model.AuthenticationDataModel, error),
	contextModel *applicationModel.ContextModel) (result model.AuthAccessTokenModel, err model.ResourceCommonErrorModel) {
	if jwtTokenStr == "" {
		err = common.GenerateUnauthorizedTokenError()
		return
	} else {
		var payload token.PayloadJWTToken
		//payload, err = MandatoryValidateJWTToken(jwtTokenStr, scope, resourceID, key)
		payload, err = ValidateTokenWithoutCheckSignature(jwtTokenStr, scope, resourceID)
		if err.Error != nil {
			return
		}
		return checkTokenInRedis(client, jwtTokenStr, checkTokenEndpoint, expiredIn, resourceID, scope, payload, roleMapping, contextModel)
	}
}

func checkTokenInRedis(client *redis.Client, token string, checkTokenEndpoint string, expiredIn int64, resourceID string, scope string, jwtPayload token.PayloadJWTToken,
	roleMapping func(clientID string, token string, payload token.PayloadJWTToken) (model.AuthAccessTokenModel, model.AuthenticationRoleModel, model.AuthenticationDataModel, error),
	contextModel *applicationModel.ContextModel) (result model.AuthAccessTokenModel, err model.ResourceCommonErrorModel) {
	var (
		fileName   = "JWTTokenValidator.go"
		funcName   = "checkTokenInRedis"
		redisModel model.RedisAuthAccessTokenModel
	)

	redisResult, errors := client.Get(token).Result()
	if errors != nil && errors.Error() != "redis: nil" {
		return
	}

	if redisResult == "" {
		var (
			body       string
			statusCode int
		)

		usedScope := jwtPayload.Scope
		result.Locale = jwtPayload.Locale

		if scope != "" {
			usedScope = scope
		}

		statusCode, body, errors = common.HitCheckTokenURL(checkTokenEndpoint, token, resourceID, usedScope, contextModel)
		if errors != nil {
			return
		}

		if statusCode != 200 {
			err = common.GenerateUnauthorizedTokenError()
			return
		} else {
			var (
				bodyObject          authentication_response.CheckTokenAuthenticationResponse
				authenticationModel model.AuthenticationModel
			)

			_ = json.Unmarshal([]byte(body), &bodyObject)
			_ = json.Unmarshal([]byte(bodyObject.Nexsoft.Payload.Data.Content.Authentication), &authenticationModel)

			result, authenticationModel.Role, authenticationModel.Data, errors = roleMapping(jwtPayload.ClientID, token, jwtPayload)
			if errors != nil {
				err = common.GenerateUnknownError(errors)
				return
			}

			if result.ResourceUserID == 0 {
				err = common.GenerateUnauthorizedTokenError()
				return
			}

			//--- Check Job Process Running
			if authenticationModel.Role.Role != constanta.RoleIDUserND6 && authenticationModel.Role.Role != constanta.RoleIDUserNexmile && !result.IsAdmin {
				jobProcess, errJobProcess := dao.JobProcessDAO.GetJobProcessRunning(serverconfig.ServerAttribute.DBConnection, repository.JobProcessModel{
					Group: sql.NullString{String: constanta.JobProcessSynchronizeGroup},
					Type:  sql.NullString{String: constanta.JobProcessMasterDataType},
					Name:  sql.NullString{String: constanta.JobProcessSynchronizeRegional},
				})

				if errJobProcess.Error != nil {
					err = common.GenerateUnknownError(errors2.New(errJobProcess.CausedBy.Error()))
					return
				}

				if jobProcess.ID.Int64 > 0 {
					timeNow := time.Now()
					timeCreatedAddHalfHour := jobProcess.CreatedAt.Time.Add(30 * time.Minute)
					if (jobProcess.Status.String != constanta.JobProcessDoneStatus && jobProcess.Status.String != constanta.JobProcessErrorStatus) && timeNow.Before(timeCreatedAddHalfHour) {
						errJobProcess = errorModel.GenerateErrorSpawnSynchronize(fileName, funcName)
						errJobProcessName := util2.GenerateI18NErrorMessage(errJobProcess, constanta.DefaultApplicationsLanguage)
						err = common.GenerateAuthenticationServerError(errJobProcess.Code, errJobProcessName)
						return
					}
				}
			}

			authenticationModel.Oauth.IsAdmin = result.IsAdmin
			authenticationServerUserID, _ := strconv.Atoi(jwtPayload.Subject)
			result.Authentication = util.StructToJSON(authenticationModel)
			result.ClientID = jwtPayload.ClientID
			result.AuthenticationServerUserID = int64(authenticationServerUserID)
			result.IPWhiteList = bodyObject.Nexsoft.Payload.Data.Content.IPWhitelist
			result.Scope = jwtPayload.Scope

			var expiration time.Duration
			fmt.Println("Expiration Exist : ", expiredIn)
			fmt.Println("Expiration Exist JWT : ", jwtPayload.ExpiresAt)
			if expiredIn == 0 {
				expiration = time.Until(time.Unix(jwtPayload.ExpiresAt, 0))
				//expiration = time.Duration(2 * time.Minute)
				fmt.Println("Time Expiration New : ", expiration)
			} else {
				expiration = time.Duration(expiredIn)
				fmt.Println("Time Expiration Exist : ", expiration)
			}

			redisModel = result.ConvertToRedisModel()
			client.Set(token, util.StructToJSON(redisModel), -1)
			client.Expire(token, expiration)
			fmt.Println("Redis Client Expiration : ", expiration)

			if !result.IsAdmin {
				serverconfig.ServerAttribute.RedisClientSession.Set(constanta.SessionUser+token, util.StructToJSON(redisModel), -1)
				serverconfig.ServerAttribute.RedisClientSession.Expire(constanta.SessionUser+token, expiration)
			} else {
				serverconfig.ServerAttribute.RedisClientSession.Set(constanta.SessionAdmin+token, util.StructToJSON(redisModel), -1)
				serverconfig.ServerAttribute.RedisClientSession.Expire(constanta.SessionAdmin+token, expiration)
			}

			return
		}
	} else {
		_ = json.Unmarshal([]byte(redisResult), &redisModel)
		authenticationServerUserID, _ := strconv.Atoi(jwtPayload.Subject)

		result.ClientID = jwtPayload.ClientID
		result.Locale = jwtPayload.Locale
		result.AuthenticationServerUserID = int64(authenticationServerUserID)
		result.RedisAuthAccessTokenModel = redisModel
		result.Scope = jwtPayload.Scope
		return
	}
}
