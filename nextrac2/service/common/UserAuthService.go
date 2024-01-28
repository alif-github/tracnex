package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
)

type userAuthService struct {
	service.AbstractService
}

var UserAuthService = userAuthService{}.New()

func (input userAuthService) New() (output userAuthService) {
	output.FileName = "UserAuthService.go"
	return
}

func (input userAuthService) HitActivationEmailToAuth(inputStruct in.UserActivationRequest, contextModel *applicationModel.ContextModel, clientID string) (err errorModel.ErrorModel) {
	var (
		internalToken     = resource_common_service.GenerateInternalToken("auth", 0, clientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
		authConfig        = config.ApplicationConfiguration.GetAuthenticationServer()
		authActivationUrl = authConfig.Host + authConfig.PathRedirect.InternalUser.Activation.Email
	)

	return input.doHitToAuthGeneralFunction(internalToken, authActivationUrl, http.MethodPost, util.StructToJSON(inputStruct), contextModel)
}

func (input userAuthService) HitActivationPhoneToAuth(inputStruct in.UserActivationPhoneRequest, contextModel *applicationModel.ContextModel, clientID string) (err errorModel.ErrorModel) {
	var (
		internalToken     = resource_common_service.GenerateInternalToken("auth", 0, clientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
		authConfig        = config.ApplicationConfiguration.GetAuthenticationServer()
		authActivationUrl = authConfig.Host + authConfig.PathRedirect.InternalUser.Activation.Phone
	)

	return input.doHitToAuthGeneralFunction(internalToken, authActivationUrl, http.MethodPost, util.StructToJSON(inputStruct), contextModel)
}

func (input userAuthService) HitResendOtpVerificationEmailToAuth(inputStruct authentication_request.ResendUserVerificationRequest, contextModel *applicationModel.ContextModel, clientID string) (err errorModel.ErrorModel) {
	var (
		funcName          = "HitResendOtpVerificationEmailToAuth"
		serverResource    = config.ApplicationConfiguration.GetServerResourceID()
		authConfig        = config.ApplicationConfiguration.GetAuthenticationServer()
		internalToken     = resource_common_service.GenerateInternalToken("auth", 0, clientID, serverResource, constanta.IndonesianLanguage)
		authActivationUrl = authConfig.Host + authConfig.PathRedirect.InternalUser.ResendActivation.Email
		payloadMessage    authentication_response.AuthenticationErrorResponse
	)

	statusCode, bodyResult, errorS := common.HitResendAuthenticationServer(internalToken, authActivationUrl, inputStruct, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &payloadMessage)
	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(payloadMessage.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, payloadMessage.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func (input userAuthService) GetEmailMessage(inputStruct authentication_request.ResendUserVerificationMessageParam) string {
	var (
		commonBundle = serverconfig.ServerAttribute.CommonServiceBundle
		param        = make(map[string]interface{})
		locale       = constanta.DefaultApplicationsLanguage
	)

	param[constanta.PurposeTableParam] = inputStruct.Purpose
	param[constanta.NameTableParam] = inputStruct.FirstName

	switch inputStruct.ClientTypeID {
	case constanta.ResourceNexmileID:
		param[constanta.ClientTypeTableParam] = constanta.Nexmile
	case constanta.ResourceNexstarID:
		param[constanta.ClientTypeTableParam] = constanta.Nexstar
	case constanta.ResourceNextradeID:
		param[constanta.ClientTypeTableParam] = constanta.Nextrade
	default:
		param[constanta.ClientTypeTableParam] = " - "
	}

	param[constanta.CompanyIDTableParam] = inputStruct.UniqueID1
	param[constanta.CompanyNameTableParam] = inputStruct.CompanyName
	param[constanta.BranchIDTableParam] = inputStruct.UniqueID2
	param[constanta.BranchNameTableParam] = inputStruct.BranchName
	param[constanta.SalesmanIDTableParam] = inputStruct.SalesmanID
	param[constanta.UserTableParam] = inputStruct.UserID
	param[constanta.PasswordTableParam] = inputStruct.Password
	param[constanta.OTPTableParam] = fmt.Sprintf(`{{.%s}}`, constanta.OTPTableParam)
	param[constanta.LinkTableParam] = fmt.Sprintf(`{{.%s}}`, constanta.LinkTableParam)
	param[constanta.EmailTableParam] = inputStruct.Email
	param[constanta.ClientIDTableParam] = inputStruct.ClientID
	param[constanta.RegistrationIDTableParam] = fmt.Sprintf(`{{.%s}}`, constanta.RegistrationIDTableParam)

	return util2.GenerateI18NServiceMessage(commonBundle, "OTP_MESSAGE_AUTH3", locale, param)
}

func (input userAuthService) HitResendEmailVerificationToEmail(inputStruct authentication_request.ResendEmailVerificationRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		internalToken            = resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, "", config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
		authConfig               = config.ApplicationConfiguration.GetAuthenticationServer()
		authResendEmailVerifyUrl = authConfig.Host + authConfig.PathRedirect.InternalUser.ResendActivation.Email
	)

	return input.doHitToAuthGeneralFunction(internalToken, authResendEmailVerifyUrl, http.MethodPut, util.StructToJSON(inputStruct), contextModel)
}

func (input userAuthService) doHitToAuthGeneralFunction(internalToken, authActivationUrl, method, inputStruct string, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName       = "doHitToAuthGeneralFunction"
		header         = make(map[string][]string)
		payloadMessage authentication_response.AuthenticationErrorResponse
	)

	header[common.AuthorizationHeaderConstanta] = []string{internalToken}
	statusCode, _, bodyResult, errorS := common.HitAPI(authActivationUrl, header, inputStruct, method, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &payloadMessage)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(payloadMessage.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, payloadMessage.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

//func (input userAuthService) HitUpdateEmailUser(inputStruct authentication_request.ResendEmailVerificationRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
//	var (
//		payloadMessage          authentication_response.AuthenticationErrorResponse
//		funcName                = "HitUpdateEmailUser"
//		internalToken           = resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, "", config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
//		authConfig              = config.ApplicationConfiguration.GetAuthenticationServer()
//		authResendEmailVerifUrl = authConfig.Host + authConfig.PathRedirect.InternalUser.ResendActivation.Email
//	)
//
//	header := make(map[string][]string)
//	header[common.AuthorizationHeaderConstanta] = []string{internalToken}
//
//	statusCode, _, bodyResult, errorS := common.HitAPI(authResendEmailVerifUrl, header, util.StructToJSON(inputStruct), http.MethodPut, *contextModel)
//
//	if errorS != nil {
//		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
//		return
//	}
//
//	_ = json.Unmarshal([]byte(bodyResult), &payloadMessage)
//
//	if statusCode == 200 {
//		err = errorModel.GenerateNonErrorModel()
//	} else {
//		causedBy := errors.New(payloadMessage.Nexsoft.Payload.Status.Message)
//		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, payloadMessage.Nexsoft.Payload.Status.Code, causedBy)
//		return
//	}
//
//	return
//}
