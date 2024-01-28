package CRUDUserService

import (
	"database/sql"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/serverconfig"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	"nexsoft.co.id/nextrac2/service/common"
	"nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input userService) ResendVerificationCode(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.UserRequest
	)

	inputStruct, err = input.readBodyAndValidateForViewAndResendOTP(request, contextModel, input.validateVerificationCode)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doResendVerificationCode(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_RESEND_VERIFICATION_CODE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateVerificationCode(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewUserAndResendOTP()
}

func (input userService) doResendVerificationCode(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		userOnDB repository.UserModel
	)

	userOnDB, err = dao.UserDAO.GetUserForResendVerificationCode(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})

	err = input.validationUserOnDBForResendVerificationCode(userOnDB)
	if err.Error != nil {
		return
	}

	err = input.ResendUserToAuthenticationServer(userOnDB, contextModel)
	if err.Error != nil {
		return
	}

	return
}

func (input userService) validationUserOnDBForResendVerificationCode(userOnDB repository.UserModel) (err errorModel.ErrorModel) {
	var (
		funcName = "validationUserOnDBForResendVerificationCode"
	)

	if userOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	if userOnDB.Status.String != constanta.StatusPending {
		err = errorModel.GenerateResendOTPStatusActive(input.FileName, funcName)
		return
	}

	if util2.IsStringEmpty(userOnDB.Email.String) {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(input.FileName, funcName, constanta.Email)
		return
	}

	return
}

func (input userService) ResendUserToAuthenticationServer(userOnDB repository.UserModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		dataEmailNotification in.UserRequest
		dataRequest           authentication_request.ResendUserVerificationRequest
	)

	dataEmailNotification = in.UserRequest{
		FirstName: userOnDB.FirstName.String,
		Username:  userOnDB.Username.String,
		Email:     userOnDB.Email.String,
		Phone:     strings.Trim(userOnDB.Phone.String, constanta.IndonesianCodeNumberWithDash),
	}

	dataRequest = authentication_request.ResendUserVerificationRequest{
		UserID:           userOnDB.AuthUserID.Int64,
		Email:            userOnDB.Email.String,
		EmailLinkMessage: config.ApplicationConfiguration.GetNextracFrontend().Host + config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath,
		EmailMessage:     GetEmailMessage(dataEmailNotification, false, false),
	}

	if !userOnDB.IsSystemAdmin.Bool {
		hostSysUser := config.ApplicationConfiguration.GetNextracFrontend().Host
		pathSysUser := strings.Replace(config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath, "/nexsoft-admin", "", 1)
		dataRequest.EmailLinkMessage = hostSysUser + pathSysUser
	}

	err = common.UserAuthService.HitResendOtpVerificationEmailToAuth(dataRequest, contextModel, contextModel.AuthAccessTokenModel.ClientID)
	if err.Error != nil {
		return
	}

	return
}
