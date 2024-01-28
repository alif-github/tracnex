package RefreshToken

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/session"
	util2 "nexsoft.co.id/nextrac2/util"
)

type refreshTokenService struct {
	service.AbstractService
}

var RefreshTokenService = refreshTokenService{}.New()

func (input refreshTokenService) New() (output refreshTokenService) {
	output.FileName = "RefreshTokenService.go"
	return
}

func (input refreshTokenService) StartService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.RefreshTokenDTOIn

	inputStruct, err = input.readBodyAndValidate(request, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content, header, err = input.hitAuthenticationServerRefreshToken(inputStruct, contextModel)

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("REFRESH_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input refreshTokenService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.RefreshTokenDTOIn, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	inputStruct.Authorization = request.Header.Get(constanta.TokenHeaderNameConstanta)

	err = inputStruct.ValidateRefreshToken()
	return
}

func (input refreshTokenService) hitAuthenticationServerRefreshToken(inputStruct in.RefreshTokenDTOIn, contextModel *applicationModel.ContextModel) (result out.TokenDTOOut, headerResult map[string]string, err errorModel.ErrorModel) {
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	refreshTokenPath := authenticationServer.Host + authenticationServer.PathRedirect.Token
	newToken, errorS := resource_common_service.RefreshToken(serverconfig.ServerAttribute.RedisClient, inputStruct.Authorization, config.ApplicationConfiguration.GetJWTToken().JWT, refreshTokenPath, constanta.ExpiredTokenOnRedisConstanta, inputStruct.RefreshToken, contextModel)
	if errorS.Error != nil {
		err.Error = errorS.Error
		err.Code = errorS.Code
		err.CausedBy = errorS.CausedBy
		return
	}

	headerResult = make(map[string]string)
	headerResult[constanta.TokenHeaderNameConstanta] = newToken
	result.RefreshToken = inputStruct.RefreshToken
	err = errorModel.GenerateNonErrorModel()
	return
}
