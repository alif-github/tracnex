package ResendOTPService

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input resendOTPService) validationUser(inputModel repository.UserRegistrationDetailModel, validationResult repository.UserRegistrationDetailMapping, timeNow time.Time, contextModel *applicationModel.ContextModel) (isReqToAuth bool, linkChannel string, err errorModel.ErrorModel) {
	var (
		fileName = "GenerateResendOTPService.go"
		funcName = "validationUser"
		ctID     = int(inputModel.ClientTypeID.Int64)
	)

	if validationResult.UserRegistrationDetail.ID.Int64 < 1 {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestUserNotFound)
		return
	}

	if inputModel.UniqueID1.String != validationResult.UserRegistrationDetail.UniqueID1.String {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestUserNotFound)
		return
	}

	if !util.IsStringEmpty(inputModel.UniqueID2.String) {
		if inputModel.UniqueID2.String != validationResult.UserRegistrationDetail.UniqueID2.String {
			err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestUserNotFound)
			return
		}
	}

	if validationResult.PKCEClientMapping.ID.Int64 < 1 {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestClientCredential)
		return
	}

	if ctID == constanta.ResourceNexmileID || ctID == constanta.ResourceNextradeID {
		linkChannel, isReqToAuth, err = input.validationNexmileNextrade(ctID, timeNow, validationResult, contextModel)
		if err.Error != nil {
			return
		}
	} else if ctID == constanta.ResourceNexstarID {
		linkChannel, isReqToAuth, err = input.validationNexstar(inputModel, timeNow, validationResult, contextModel)
		if err.Error != nil {
			return
		}
	} else {
		msg := fmt.Sprintf(`%s %s %s`, constanta.Nexmile, constanta.Nexstar, constanta.Nextrade)
		err = errorModel.GenerateErrorEksternalClientTypeMustHave(fileName, funcName, msg)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input resendOTPService) doValidationNexmileNextrade(validationResult repository.UserRegistrationDetailMapping, timeNow time.Time, contextModel *applicationModel.ContextModel) (isRequestCodeToAuth bool, err errorModel.ErrorModel) {
	var (
		fileName       = "GenerateResendOTPService.go"
		funcName       = "doValidationNexmileNextrade"
		userOnAuthResp authentication_response.CheckClientOrUserResponse
	)

	if validationResult.User.ID.Int64 < 1 {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestUserNotFound)
		return
	}

	if (validationResult.User.Status.String != constanta.StatusNonActive) && (validationResult.User.Status.String != constanta.PendingOnApproval) {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.UserActiveStatus, constanta.Status, "")
		return
	}

	//--- Check user on Auth
	userOnAuthResp, err = service.CheckClientOrUserInAuth(authentication_request.CheckClientOrUser{ClientID: validationResult.PKCEClientMapping.ClientID.String}, contextModel)
	if err.Error != nil {
		return
	}

	userOnAuth := userOnAuthResp.Nexsoft.Payload.Data.Content
	if validationResult.User.Status.String == constanta.PendingOnApproval {
		if userOnAuth.IsExist {
			if userOnAuth.AdditionalInformation.UserStatus != 1 {
				isRequestCodeToAuth = true
				return
			}
		} else {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserAuth)
			return
		}
	}

	err = input.validationUserVerification(validationResult, timeNow)
	return
}

func (input resendOTPService) doValidationNexstar(_ repository.UserRegistrationDetailModel, validationResult repository.UserRegistrationDetailMapping, timeNow time.Time) (isRequestCodeToAuth bool, err errorModel.ErrorModel) {
	var (
		fileName = "GenerateResendOTPService.go"
		funcName = "doValidationNexstar"
	)

	if validationResult.User.ID.Int64 < 1 {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestUserNotFound)
		return
	}

	if (validationResult.User.Status.String != constanta.StatusNonActive) && (validationResult.User.Status.String != constanta.PendingOnApproval) {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.UserActiveStatus, constanta.Status, "")
		return
	}

	//---- Fendy 19/05/2023
	//if inputModel.NoTelp.String != validationResult.UserVerification.Phone.String {
	//	err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestPhoneNotFound)
	//	return
	//}

	err = input.validationUserVerification(validationResult, timeNow)
	return
}

func (input resendOTPService) validationUserVerification(validationResult repository.UserRegistrationDetailMapping, timeNow time.Time) (err errorModel.ErrorModel) {
	var (
		fileName = "ValidationGenerateResendOTPService.go"
		funcName = "validationUserVerification"
	)

	if validationResult.UserVerification.ID.Int64 < 1 {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestUserNotFound)
		return
	}

	timeStr := timeNow.Format(constanta.DefaultTimeFormat)
	timeNowNew, _ := time.Parse(constanta.DefaultTimeFormat, timeStr)

	difTime := timeNowNew.Sub(validationResult.UserVerification.UpdatedAt.Time)
	if difTime < (2 * time.Minute) {
		err = errorModel.GenerateErrorEksternalSpawnTimeResendOTP(fileName, funcName)
		return
	}

	return
}

func (input resendOTPService) validationNexmileNextrade(clientType int, timeNow time.Time, validationResult repository.UserRegistrationDetailMapping, contextModel *applicationModel.ContextModel) (linkChannel string, isReqToAuth bool, err errorModel.ErrorModel) {
	var (
		cfgNexmile  = config.ApplicationConfiguration.GetNexmile()
		cfgNextrade = config.ApplicationConfiguration.GetNextrade()
	)

	isReqToAuth, err = input.doValidationNexmileNextrade(validationResult, timeNow, contextModel)
	if err.Error != nil {
		return
	}

	if clientType == constanta.ResourceNextradeID {
		linkChannel = cfgNextrade.Host + cfgNextrade.PathRedirect.ActivationUser
	} else {
		linkChannel = cfgNexmile.Host + cfgNexmile.PathRedirect.ActivationUser
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input resendOTPService) validationNexstar(inputModel repository.UserRegistrationDetailModel, timeNow time.Time, validationResult repository.UserRegistrationDetailMapping, _ *applicationModel.ContextModel) (linkChannel string, isReqToAuth bool, err errorModel.ErrorModel) {
	var cfgNexstar = config.ApplicationConfiguration.GetNexstar()
	isReqToAuth, err = input.doValidationNexstar(inputModel, validationResult, timeNow)
	if err.Error != nil {
		return
	}

	linkChannel = cfgNexstar.Host + cfgNexstar.PathRedirect.ActivationUser
	err = errorModel.GenerateNonErrorModel()
	return
}
