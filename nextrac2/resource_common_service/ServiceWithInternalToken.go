package resource_common_service

import (
	"database/sql"
	"encoding/json"
	"github.com/bukalapak/go-redis"
	"github.com/dgrijalva/jwt-go"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	token2 "nexsoft.co.id/nextrac2/resource_common_service/token"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/token"
	"strconv"
	"time"
)

func ValidateJWTInternal(isCheckClientID bool, client *redis.Client, addResourceUrl string,
	jwtTokenStr string, key string, resourceID string, tokenInternal string,
	roleMapping func(clientID string, userID int64) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error),
	saveClientToDB func(result authentication_response.AddClientAuthenticationResponse, createdBy int64) (int64, error), contextModel *applicationModel.ContextModel) (result model2.AuthAccessTokenModel, err model2.ResourceCommonErrorModel) {
	var payload token.PayloadJWTInternal

	if jwtTokenStr == "" {
		err = common.GenerateUnauthorizedTokenError()
		return
	} else {
		payload, err = token2.ValidateJWTInternal(jwtTokenStr, key)
		result.Locale = payload.Locale
		if err.Error != nil {
			return
		}

		if isCheckClientID {
			result, err = checkInternalTokenInRedis(client, resourceID, payload, addResourceUrl, tokenInternal, roleMapping, saveClientToDB, contextModel)
			if err.Error != nil {
				return
			}

		}

		if !common.CheckIsResourceIDExist(payload.Resource, resourceID) {
			err = common.GenerateForbiddenByResourceID()
			return
		}
	}

	return
}

func checkInternalTokenInRedis(client *redis.Client, resourceID string, payload token.PayloadJWTInternal, addResourceUrl string,
	tokenInternal string, roleMapping func(clientID string, userID int64) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error), saveClientToDB func(result authentication_response.AddClientAuthenticationResponse, createdBy int64) (int64, error),
	contextModel *applicationModel.ContextModel) (result model2.AuthAccessTokenModel, err model2.ResourceCommonErrorModel) {

	var authenticationRoleModel model2.AuthenticationRoleModel
	var authenticationDataModel model2.AuthenticationDataModel
	//var redisModel model2.RedisAuthAccessTokenModel
	redisResult, errors := client.Get(payload.ClientID).Result()
	if errors != nil && errors.Error() != "redis: nil" {
		err.Error = errors
		return
	}

	if redisResult == "" {
		result, authenticationRoleModel, authenticationDataModel, err = checkInternalTokenInDB(payload, roleMapping, resourceID, addResourceUrl, tokenInternal, saveClientToDB, contextModel)
		if err.Error != nil {
			return
		}

		result.Authentication = util.StructToJSON(model2.AuthenticationModel{
			Role: authenticationRoleModel,
			Data: authenticationDataModel,
		})
		result.Locale = payload.Locale
		client.Set(payload.ClientID, util.StructToJSON(result), -1)
	} else {
		_ = json.Unmarshal([]byte(redisResult), &result)
	}

	return
}

func checkInternalTokenInDB(payload token.PayloadJWTInternal, roleMapping func(clientID string, userID int64) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error),
	resourceID string, addResourceUrl string, tokenInternal string, saveClientToDB func(result authentication_response.AddClientAuthenticationResponse, createdBy int64) (int64, error),
	contextModel *applicationModel.ContextModel) (result model2.AuthAccessTokenModel, result2 model2.AuthenticationRoleModel, result3 model2.AuthenticationDataModel, err model2.ResourceCommonErrorModel) {
	var errors error
	var userID int

	userID, errors = strconv.Atoi(payload.Subject)
	if errors != nil {
		err = common.GenerateUnauthorizedTokenError()
		return
	}

	result, result2, result3, errors = roleMapping(payload.ClientID, int64(userID))
	if errors != nil {
		err.Error = errors
		err.Code = 500
		return
	}

	if result.ResourceUserID == 0 {
		if !config.ApplicationConfiguration.GetServerAutoAddClient() {
			err = common.GenerateUnauthorizedTokenError()
			return
		}

		var addClientResourceResult authentication_response.AddClientAuthenticationResponse
		var resourceUserID int64

		statusCode, bodyResult, errors := common.HitAddClientResource(tokenInternal, addResourceUrl, payload.ClientID, resourceID, contextModel)
		if errors != nil {
			err.Error = errors
			err.Code = 500
			return
		}

		if statusCode == 200 {
			_ = json.Unmarshal([]byte(bodyResult), &addClientResourceResult)
			resourceUserID, errors = saveClientToDB(addClientResourceResult, 0)
			if errors != nil {
				err.Error = errors
				err.Code = 500
				return
			}
			result.ClientID = payload.ClientID
			result.ResourceUserID = resourceUserID
		} else {
			err = common.GenerateUnauthorizedTokenError()
		}
	}

	sub, _ := strconv.Atoi(payload.Subject)

	if result.ResourceUserID != int64(sub) {
		userModel, errs := dao.UserDAO.CheckIsAuthUserExist(serverconfig.ServerAttribute.DBConnection, repository.UserModel{AuthUserID: sql.NullInt64{Int64: int64(sub)}})
		if errs.Error != nil {
			return result, result2, result3, common.GenerateUnknownError(errs.CausedBy)
		}
		if userModel.ID.Int64 == 0 {
			return result, result2, result3, common.GenerateUnauthorizedTokenError()
		}
		result.ResourceUserID = userModel.ID.Int64
	}

	result = ReadAuthTokenAndPayload(result, payload)
	return
}

func ReadAuthTokenAndPayload(authModel model2.AuthAccessTokenModel, payload token.PayloadJWTInternal) model2.AuthAccessTokenModel {
	authModel.ClientID = payload.UserClient
	userID, _ := strconv.Atoi(payload.Subject)
	authModel.AuthenticationServerUserID = int64(userID)

	return authModel
}

func GenerateInternalToken(resourceDestination string, userID int64, clientID string, issuer string, locale string) string {
	userClientID := config.ApplicationConfiguration.GetClientCredentialsClientID()
	if clientID != "" {
		userClientID = clientID
	}

	usedUserID := config.ApplicationConfiguration.GetClientCredentialsAuthUserID()
	if userID > 0 {
		usedUserID = userID
	}

	tokenCode := token.PayloadJWTInternal{
		Locale:     locale,
		ClientID:   config.ApplicationConfiguration.GetClientCredentialsClientID(),
		UserClient: userClientID,
		Resource:   resourceDestination,
		Version:    config.ApplicationConfiguration.GetServerVersion(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			Subject:   strconv.Itoa(int(usedUserID)),
		},
	}

	jwtToken, _ := token.JWTToken{}.GenerateToken(tokenCode, config.ApplicationConfiguration.GetJWTToken().Internal)

	return jwtToken
}
