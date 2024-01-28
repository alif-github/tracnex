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
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/session"
	util2 "nexsoft.co.id/nextrac2/util"
)

type authorizeService struct {
	service.AbstractService
}

var AuthorizeService = authorizeService{}.New()

func (input authorizeService) New() (output authorizeService) {
	output.FileName = "AuthorizeService.go"
	return
}

func (input authorizeService) StartService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.AuthorizeDTOIn

	inputStruct, err = input.readBodyAndValidate(request, contextModel)
	if err.Error != nil {
		return
	}

	header, err = input.HitAuthorizeAuthenticationServer(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("AUTHORIZE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input authorizeService) readBodyAndValidate(request *http.Request, loggerModel *applicationModel.ContextModel) (inputStruct in.AuthorizeDTOIn, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, loggerModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	err = inputStruct.ValidateAuthorize()
	return
}

func (input authorizeService) HitAuthorizeAuthenticationServer(inputStruct in.AuthorizeDTOIn, loggerModel *applicationModel.ContextModel) (headerResult map[string]string, err errorModel.ErrorModel) {
	var (
		funcName             = "HitAuthorizeAuthenticationServer"
		authenticationServer = config.ApplicationConfiguration.GetAuthenticationServer()
		authorizePath        = authenticationServer.Host + authenticationServer.PathRedirect.Authorize
		requestBody          authentication_request.AuthorizeRequestDTO
	)

	requestBody = authentication_request.AuthorizeRequestDTO{
		CodeChallenger: inputStruct.CodeChallenger,
		ResponseType:   "code_pkce",
	}

	statusCode, header, bodyResult, errorS := common.HitAPI(authorizePath, nil, util.StructToJSON(requestBody), "POST", *loggerModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		var authorizeResult authentication_response.AuthorizeAuthenticationResponse
		_ = json.Unmarshal([]byte(bodyResult), &authorizeResult)
		headerResult = make(map[string]string)
		headerResult[constanta.TokenHeaderNameConstanta] = header[constanta.TokenHeaderNameConstanta][0]
		err = errorModel.GenerateNonErrorModel()
	} else {
		var errorResult authentication_response.AuthenticationErrorResponse
		_ = json.Unmarshal([]byte(bodyResult), &errorResult)
		causedBy := errors.New(errorResult.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, errorResult.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}
