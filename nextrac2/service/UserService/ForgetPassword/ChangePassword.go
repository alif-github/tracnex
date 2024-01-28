package ForgetPassword

import (
	"bytes"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input forgetPasswordService) ChangePassword(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.ChangePasswordByEmailRequest
	)

	inputStruct, err = input.readBodyForChangePassword(request, input.validateChangePassword)
	if err.Error != nil {
		return
	}

	header, err = input.doChangePassword(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_CHANGE_PASSWORD", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input forgetPasswordService) validateChangePassword(inputStruct *in.ChangePasswordByEmailRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateStructChangePassword()
}

func (input forgetPasswordService) doChangePassword(inputStruct in.ChangePasswordByEmailRequest, contextModel *applicationModel.ContextModel) (header map[string]string, err errorModel.ErrorModel) {
	return input.hitChangePasswordAuth(inputStruct, contextModel)
}

func (input forgetPasswordService) hitChangePasswordAuth(inputStruct in.ChangePasswordByEmailRequest, contextModel *applicationModel.ContextModel) (header map[string]string, err errorModel.ErrorModel) {
	var (
		fileName               = input.FileName
		funcName               = "hitChangePasswordAuth"
		client                 = &http.Client{}
		requestBody            = bytes.NewBuffer([]byte(inputStruct.ToString()))
		hostAuth               = config.ApplicationConfiguration.GetAuthenticationServer().Host
		pathChangePasswordAuth = config.ApplicationConfiguration.GetAuthenticationServer().PathRedirect.InternalUser.Forget.ChangePassword.Email
	)

	// Hit Authentication Server
	stmt, errs := http.NewRequest(http.MethodPost, hostAuth+pathChangePasswordAuth, requestBody)
	if errs != nil {
		err = errorModel.GenerateUnauthorizedClientError(input.FileName, funcName)
		return
	}

	// Set Internal Token For Hit Authentication Server
	var (
		resourceDestination = "auth"
		clientID            = contextModel.AuthAccessTokenModel.ClientID
		issuer              = config.ApplicationConfiguration.GetServerResourceID()
		locale              = constanta.IndonesianLanguage
		userID              int64
		internalToken       string
	)

	internalToken = resource_common_service.GenerateInternalToken(resourceDestination, userID, clientID, issuer, locale)

	stmt.Header.Add(constanta.TokenHeaderNameConstanta, internalToken)
	//stmt.Header.Add(constanta.DefaultTokenKeyConstanta, "authServer2020")

	response, errs := client.Do(stmt)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}
	stmt.Body.Close()

	// Check Response Message
	if response.StatusCode != 200 {
		var APIError out.APIResponse
		json.NewDecoder(response.Body).Decode(&APIError)
		err = errorModel.GenerateErrorModelWithoutCaused(response.StatusCode,APIError.Nexsoft.Payload.Status.Message, fileName, funcName)
		return
	}

	header = make(map[string]string)
	header[constanta.TokenHeaderNameConstanta] = response.Header.Get(constanta.TokenHeaderNameConstanta)
	header[constanta.CodeResponseValidateConstanta] = response.Header.Get(constanta.CodeResponseValidateConstanta)
	return
}
