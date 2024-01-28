package Login

import (
	"encoding/json"
	"errors"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/session"
	"nexsoft.co.id/nextrac2/token"
	util2 "nexsoft.co.id/nextrac2/util"
)

type tokenService struct {
	service.AbstractService
}

var TokenService = tokenService{}.New()

func (input tokenService) New() (output tokenService) {
	output.FileName = "TokenService.go"
	return
}

func (input tokenService) NexsoftTokenService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	return input.startTokenService(request, contextModel, endpoint.RoleMappingUserNexsoft)
}

func (input tokenService) UserTokenService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	return input.startTokenService(request, contextModel, endpoint.RoleMappingUser)
}

func (input tokenService) startTokenService(request *http.Request, contextModel *applicationModel.ContextModel, roleMapping func(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error)) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.TokenDTOIn
	inputStruct, err = input.readBodyAndValidate(request, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content, _, header, err = input.HitTokenAuthenticationServer(inputStruct, contextModel, roleMapping)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("TOKEN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input tokenService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.TokenDTOIn, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	err = inputStruct.ValidateToken()
	return
}

func (input tokenService) HitTokenAuthenticationServer(inputStruct in.TokenDTOIn, contextModel *applicationModel.ContextModel, roleMapping func(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error)) (result out.TokenDTOOut, authAccessTokenModel model2.AuthAccessTokenModel, headerResult map[string]string, err errorModel.ErrorModel) {
	var (
		funcName             = "HitTokenAuthenticationServer"
		authenticationServer = config.ApplicationConfiguration.GetAuthenticationServer()
		tokenPath            = authenticationServer.Host + authenticationServer.PathRedirect.Token
		requestBody          authentication_request.TokenRequestDTO
	)

	requestBody = authentication_request.TokenRequestDTO{
		CodeVerifier:      inputStruct.CodeVerifier,
		AuthorizationCode: inputStruct.AuthorizationCode,
	}

	statusCode, header, bodyResult, errorS := common.HitAPI(tokenPath, nil, util.StructToJSON(requestBody), "POST", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		var (
			tokenResult   authentication_response.TokenAuthenticationResponse
			errs          model2.ResourceCommonErrorModel
			checkTokenUrl = authenticationServer.Host + authenticationServer.PathRedirect.CheckToken
		)

		_ = json.Unmarshal([]byte(bodyResult), &tokenResult)
		headerResult = make(map[string]string)
		headerResult[constanta.TokenHeaderNameConstanta] = header[constanta.TokenHeaderNameConstanta][0]

		authAccessTokenModel, errs = resource_common_service.ValidateJWTToken(serverconfig.ServerAttribute.RedisClient, header[constanta.TokenHeaderNameConstanta][0], 0, checkTokenUrl, config.ApplicationConfiguration.GetServerResourceID(), "read write", config.ApplicationConfiguration.GetJWTToken().JWT, roleMapping, contextModel)
		if errs.Error != nil {
			err = endpoint.ReadError(errs)
			return
		}

		result.RefreshToken = tokenResult.Nexsoft.Payload.Data.Content.RefreshToken
		err = errorModel.GenerateNonErrorModel()
	} else {
		var errorResult authentication_response.AuthenticationErrorResponse
		_ = json.Unmarshal([]byte(bodyResult), &errorResult)
		causedBy := errors.New(errorResult.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, errorResult.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}

func (input tokenService) ClientTokenService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	return input.startClientTokenService(request, contextModel)
}

func (input tokenService) startClientTokenService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.TokenClientDTOIn

	inputStruct, err = input.readBodyAndValidateForClient(request, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content, _, header, err = input.hitTokenClientAuthenticationServer(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("TOKEN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input tokenService) readBodyAndValidateForClient(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.TokenClientDTOIn, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	err = inputStruct.ValidateTokenClient()
	return
}

func (input tokenService) hitTokenClientAuthenticationServer(inputStruct in.TokenClientDTOIn, contextModel *applicationModel.ContextModel) (result out.TokenDTOOut, authAccessTokenModel model2.AuthAccessTokenModel, headerResult map[string]string, err errorModel.ErrorModel) {
	funcName := "hitTokenClientAuthenticationServer"

	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	tokenPath := authenticationServer.Host + authenticationServer.PathRedirect.Token

	statusCode, header, bodyResult, errorS := common.HitAPI(tokenPath, nil, util.StructToJSON(inputStruct), "POST", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		var tokenResult authentication_response.TokenAuthenticationResponse
		_ = json.Unmarshal([]byte(bodyResult), &tokenResult)
		headerResult = make(map[string]string)
		headerResult[constanta.TokenHeaderNameConstanta] = header[constanta.TokenHeaderNameConstanta][0]

		result.RefreshToken = tokenResult.Nexsoft.Payload.Data.Content.RefreshToken
		err = errorModel.GenerateNonErrorModel()
	} else {
		var errorResult authentication_response.AuthenticationErrorResponse
		_ = json.Unmarshal([]byte(bodyResult), &errorResult)
		causedBy := errors.New(errorResult.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, errorResult.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}
