package Login

import (
	"encoding/json"
	"errors"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/session"
	util2 "nexsoft.co.id/nextrac2/util"
	"sync"
)

type logoutService struct {
	service.AbstractService
}

var LogoutService = logoutService{}.New()

func (input logoutService) New() (output logoutService) {
	output.FileName = "LogoutService.go"
	return
}

func (input logoutService) StartService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	token := request.Header.Get(constanta.TokenHeaderNameConstanta)

	err = input.deleteTokenFromRedis(token)
	if err.Error != nil {
		return
	}

	err = input.deleteTokenFromRedisSession(token)
	if err.Error != nil {
		return
	}

	err = HitLogoutAuthenticationServer(token, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("LOGOUT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input logoutService) deleteTokenFromRedis(token string) errorModel.ErrorModel {
	funcName := "deleteTokenFromRedis"
	if token == "" {
		return errorModel.GenerateUnauthorizedClientError(input.FileName, funcName)
	}
	serverconfig.ServerAttribute.RedisClient.Del(token)
	return errorModel.GenerateNonErrorModel()
}

func (input logoutService) deleteTokenFromRedisSession(token string) errorModel.ErrorModel {
	funcName := "deleteTokenFromRedisSession"
	if token == "" {
		return errorModel.GenerateUnauthorizedClientError(input.FileName, funcName)
	}

	serverconfig.ServerAttribute.RedisClientSession.Del(constanta.SessionUser + token)
	serverconfig.ServerAttribute.RedisClientSession.Del(constanta.SessionAdmin + token)

	return errorModel.GenerateNonErrorModel()
}

func HitLogoutAuthenticationServer(token string, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	fileName := "TokenService.go"
	funcName := "HitLogoutAuthenticationServer"
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	tokenPath := authenticationServer.Host + authenticationServer.PathRedirect.Logout

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{token}

	statusCode, _, bodyResult, errorS := common.HitAPI(tokenPath, headerRequest, "", "GET", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		var errorResult authentication_response.AuthenticationErrorResponse
		_ = json.Unmarshal([]byte(bodyResult), &errorResult)
		causedBy := errors.New(errorResult.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, errorResult.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}

func LogoutAuthServerAutomatic(listToken []string, contextModel applicationModel.ContextModel) {
	var wg sync.WaitGroup

	for i := 0; i < len(listToken); i++ {
		wg.Add(1)

		a := i

		go func() {
			defer wg.Done()
			HitLogoutAuthenticationServer(listToken[a], &contextModel)
		}()
	}

	wg.Wait()
}
