package ForgetPassword

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input forgetPasswordService) ResetPassword(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.ResetPasswordRequest
	)

	inputStruct, err = input.readBodyForResetPassword(request, input.validateResetPassword)
	if err.Error != nil {
		return
	}

	header, err = input.doResetPassword(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_RESET_PASSWORD", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input forgetPasswordService) doResetPassword(inputStruct in.ResetPasswordRequest, contextModel *applicationModel.ContextModel) (header map[string]string, err errorModel.ErrorModel) {
	return input.hitForgetPasswordAuth(inputStruct, contextModel)
}

func (input forgetPasswordService) validateResetPassword(inputStruct *in.ResetPasswordRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateStructForgetPassword()
}

func (input forgetPasswordService) hitForgetPasswordAuth(inputStruct in.ResetPasswordRequest, contextModel *applicationModel.ContextModel) (header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName                 = "hitForgetPasswordAuth"
		client                   = &http.Client{}
		hostAuth                 = config.ApplicationConfiguration.GetAuthenticationServer().Host
		hostNextrac              = config.ApplicationConfiguration.GetNextracFrontend().Host
		pathResetPasswordNextrac = config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.ResetPasswordPath
		pathForgetPasswordAuth   = config.ApplicationConfiguration.GetAuthenticationServer().PathRedirect.InternalUser.Forget.Email
	)

	inputStruct.EmailLink = hostNextrac + pathResetPasswordNextrac
	inputStruct.EmailMessage = constanta.EmailResetPassword

	// Hit Authentication Server
	requestBody := bytes.NewBuffer([]byte(inputStruct.ToString()))
	stmt, errs := http.NewRequest(http.MethodPost, hostAuth+pathForgetPasswordAuth, requestBody)
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

	response, errs := client.Do(stmt)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	stmt.Body.Close()

	// Check Response Auth
	if response.StatusCode != 200 {
		bodyResultByte, _ := ioutil.ReadAll(response.Body)
		bodyResult := string(bodyResultByte)
		err = common.ReadAuthServerError(funcName, response.StatusCode, bodyResult, contextModel)

		return
	}

	var APIResponse out.APIResponse
	json.NewDecoder(response.Body).Decode(&APIResponse)
	contextModel.LoggerModel.Message = APIResponse.Nexsoft.Payload.Status.Message

	header = make(map[string]string)
	header[constanta.TokenHeaderNameConstanta] = response.Header.Get(constanta.TokenHeaderNameConstanta)
	header[constanta.X_NEXCODE] = response.Header.Get(constanta.X_NEXCODE)
	return
}
