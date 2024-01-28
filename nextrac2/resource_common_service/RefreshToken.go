package resource_common_service

import (
	"database/sql"
	"encoding/json"
	"github.com/bukalapak/go-redis"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/resource_common_service/token"
	"nexsoft.co.id/nextrac2/serverconfig"
	"time"
)

func RefreshToken(client *redis.Client, oldToken string, tokenKey string, urlRefreshToken string, expiredIn int64, refreshToken string, contextModel *applicationModel.ContextModel) (newToken string, err model.ResourceCommonErrorModel) {
	redisResult, errors := client.Get(oldToken).Result()
	if errors != nil && errors.Error() != "redis: nil" {
		err = common.GenerateUnauthorizedTokenError()
		return
	} else {
		statusCode, header, resultBody, errorS := refreshTokenOnAuthenticationServer(oldToken, refreshToken, urlRefreshToken, *contextModel)
		if errorS != nil {
			err = common.GenerateUnknownError(errorS)
			util.LogError(contextModel.LoggerModel.ToLoggerObject())
			return
		} else {
			if statusCode != 200 {
				var authServerError authentication_response.AuthenticationErrorResponse
				_ = json.Unmarshal([]byte(resultBody), &authServerError)
				err = common.GenerateAuthenticationServerErrorWithMessage(statusCode, authServerError.Nexsoft.Payload.Status.Code, authServerError.Nexsoft.Payload.Status.Message)
			} else {
				newToken = header["Authorization"][0]

				payload, _ := ConvertJWTToPayload(newToken)
				tx, errs := serverconfig.ServerAttribute.DBConnection.Begin()

				defer func() {
					if errs != nil || err.Error != nil {
						_ = tx.Rollback()
						return
					} else {
						_ = tx.Commit()
					}
				}()

				errors := dao.UserDAO.UpdateLastTokenUser(tx, repository.UserModel{ClientID: sql.NullString{String: payload.ClientID}})
				if errors.Error != nil {
					err.Error = errors.CausedBy
					return
				}

				go refreshTokenOnRedis(client, redisResult, tokenKey, newToken, oldToken, expiredIn)
			}
		}
	}
	return
}

func refreshTokenOnRedis(client *redis.Client, redisResult string, tokenKey string, newToken string, oldToken string, expiredIn int64) {
	var accessTokenModel model.AuthAccessTokenModel
	_ = json.Unmarshal([]byte(redisResult), &accessTokenModel)

	payload, _ := token.ValidateJWT(newToken, tokenKey)
	accessTokenModel.Locale = payload.Locale

	var expiration time.Duration
	if expiredIn == 0 {
		expiration = time.Until(time.Unix(payload.ExpiresAt, 0))
	} else {
		expiration = time.Duration(expiredIn)
	}

	client.Set(newToken, util.StructToJSON(accessTokenModel), -1)
	client.Expire(newToken, expiration)
	client.Del(oldToken)
	return
}

func refreshTokenOnAuthenticationServer(oldToken string, refreshToken string, urlRefreshToken string, contextModel applicationModel.ContextModel) (statusCode int, headerResult map[string][]string, bodyResult string, err error) {
	header := make(map[string][]string)
	header["Authorization"] = []string{oldToken}

	refreshTokenBody := model.RefreshTokenBody{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	return common.HitAPI(urlRefreshToken, header, util.StructToJSON(refreshTokenBody), "POST", contextModel)
}
