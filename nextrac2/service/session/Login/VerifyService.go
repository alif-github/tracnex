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

type verifyService struct {
	service.AbstractService
}

var VerifyService = verifyService{}.New()

func (input verifyService) New() (output verifyService) {
	output.FileName = "VerifyService.go"
	return
}

func (input verifyService) StartService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.VerifyDTOIn

	inputStruct, err = input.readBodyAndValidate(request, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.HitVerifyAuthenticationServer(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("VERIFY_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input verifyService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.VerifyDTOIn, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	inputStruct.Authorize = request.Header.Get(constanta.TokenHeaderNameConstanta)
	err = inputStruct.ValidateLoginDTO()
	return
}

func (input verifyService) HitVerifyAuthenticationServer(inputStruct in.VerifyDTOIn, contextModel *applicationModel.ContextModel) (result out.VerifyDTOOut, err errorModel.ErrorModel) {
	var (
		funcName             = "HitVerifyAuthenticationServer"
		authenticationServer = config.ApplicationConfiguration.GetAuthenticationServer()
		verifyPath           = authenticationServer.Host + authenticationServer.PathRedirect.Verify
		requestBody          authentication_request.VerifyRequestDTO
	)

	requestBody = authentication_request.VerifyRequestDTO{
		Password: inputStruct.Password,
	}

	if !util.IsStringEmpty(inputStruct.Email) {
		requestBody.Email = inputStruct.Email
	} else {
		requestBody.Username = inputStruct.Username
	}

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{inputStruct.Authorize}

	statusCode, header, bodyResult, errorS := common.HitAPI(verifyPath, headerRequest, util.StructToJSON(requestBody), "POST", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		var authorizeResult authentication_response.AuthorizeAuthenticationResponse
		_ = json.Unmarshal([]byte(bodyResult), &authorizeResult)
		result.AuthorizationCode = header[constanta.CodeResponseValidateConstanta][0]
		err = errorModel.GenerateNonErrorModel()
	} else {
		var errorResult authentication_response.AuthenticationErrorResponse
		_ = json.Unmarshal([]byte(bodyResult), &errorResult)
		causedBy := errors.New(errorResult.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, errorResult.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}
