package ResendOTPService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/common"
	"time"
)

func (input resendOTPService) GenerateResendOTPService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "GenerateResendOTPService"
		inputStruct in.UserVerificationRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateResendOTP)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doGenerateResendOTPService, func(interface{}, applicationModel.ContextModel) {
		//--- function additional
	})

	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GENERATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input resendOTPService) doGenerateResendOTPService(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName         = "GenerateResendOTPService.go"
		funcName         = "doGenerateResendOTPService"
		inputStruct      = inputStructInterface.(in.UserVerificationRequest)
		db               = serverconfig.ServerAttribute.DBConnection
		subject          = "Resend Activation"
		purpose          = " (resend) "
		linkChannel      string
		isReqToAuth      bool
		inputModel       repository.UserRegistrationDetailModel
		validationResult repository.UserRegistrationDetailMapping
		userVerifyModel  repository.UserVerificationModel
		userRegisAdmin   repository.UserRegistrationAdminModel
	)

	//--- Create model and get data mapping on DB
	inputModel = input.convertDTOToModel(inputStruct)
	err, validationResult = dao.UserRegistrationDetailDAO.GetMappingUserValidation(db, inputModel)
	if err.Error != nil {
		return
	}

	//--- Validation mapping
	isReqToAuth, linkChannel, err = input.validationUser(inputModel, validationResult, timeNow, contextModel)
	if err.Error != nil {
		return
	}

	//--- Get Company Name and Branch Name
	if validationResult.ClientMapping.ID.Int64 > 0 {
		userRegisAdmin, err = dao.UserRegistrationAdminDAO.GetUserRegistrationByClientMappingID(db, repository.UserRegistrationAdminModel{
			UniqueID1:       sql.NullString{String: inputModel.UniqueID1.String},
			UniqueID2:       sql.NullString{String: inputModel.UniqueID2.String},
			ClientMappingID: sql.NullInt64{Int64: validationResult.ClientMapping.ID.Int64},
			ClientTypeID:    sql.NullInt64{Int64: inputModel.ClientTypeID.Int64},
		})

		if err.Error != nil {
			return
		}

		if userRegisAdmin.ID.Int64 < 1 {
			err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedRequestUserNotFound)
			return
		}
	}

	validationResult.PKCEClientMapping.CompanyName.String = userRegisAdmin.CompanyName.String
	validationResult.PKCEClientMapping.BranchName.String = userRegisAdmin.BranchName.String

	//--- Resend otp to auth
	if isReqToAuth {
		service.LogMessage(fmt.Sprintf(`Request To Auth Sending...`), 200)
		err = input.ResendOTPToAuth(inputModel, validationResult, *contextModel, linkChannel, purpose)
		if err.Error != nil {
			return
		}

		//--- Return if auth success
		err = errorModel.GenerateNonErrorModel()
		return
	}

	//--- Update user verification
	userVerifyModel = input.createModelForUpdateUserVerification(timeNow, inputModel, validationResult)
	err = dao.UserVerificationDAO.UpdateUserVerification(tx, userVerifyModel)
	if err.Error != nil {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.FailedVerificationProcess)
		return
	}

	//--- Send to Email
	service.LogMessage(fmt.Sprintf(`Send Email From Nextrac`), 200)
	err = input.SendMessageToEmail(validationResult, subject, purpose, userVerifyModel, linkChannel, input.linkQueryEmail)
	if err.Error != nil {
		return
	}

	//--- Send to Gro Chat
	if inputModel.ClientTypeID.Int64 == constanta.ResourceNexstarID {
		err = input.sendMessageToGrochat(validationResult, userVerifyModel, contextModel, linkChannel)
		if err.Error != nil {
			return
		}
	}

	//--- For testing only, delete if unused
	//err = input.sendMessageToGrochat(validationResult, userVerifyModel, contextModel, linkChannel)
	//if err.Error != nil {
	//	return
	//}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input resendOTPService) ResendOTPToAuth(model repository.UserRegistrationDetailModel, dataOnDB repository.UserRegistrationDetailMapping,
	contextModel applicationModel.ContextModel, linkChannel string, purpose string) (err errorModel.ErrorModel) {
	var (
		fileName = "GenerateResendOTPService.go"
		funcName = "ResendOTPToAuth"
		email    = model.Email.String
		auth     authentication_request.ResendUserVerificationRequest
	)

	//--- Check Email
	if util.IsStringEmpty(email) {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Email)
		return
	}

	//--- Create Query Link
	linkChannel, err = input.linkQueryAuth(linkChannel, dataOnDB)
	if err.Error != nil {
		return
	}

	//--- Build Body Request
	auth = authentication_request.ResendUserVerificationRequest{
		UserID:           dataOnDB.UserRegistrationDetail.AuthUserID.Int64,
		Email:            email,
		EmailLinkMessage: linkChannel,
		EmailMessage: common.UserAuthService.GetEmailMessage(authentication_request.ResendUserVerificationMessageParam{
			Purpose:      purpose,
			FirstName:    dataOnDB.User.FirstName.String,
			ClientTypeID: dataOnDB.PKCEClientMapping.ClientTypeID.Int64,
			UniqueID1:    dataOnDB.UserRegistrationDetail.UniqueID1.String,
			CompanyName:  dataOnDB.PKCEClientMapping.CompanyName.String,
			UniqueID2:    dataOnDB.UserRegistrationDetail.UniqueID2.String,
			BranchName:   dataOnDB.PKCEClientMapping.BranchName.String,
			SalesmanID:   dataOnDB.UserRegistrationDetail.SalesmanID.String,
			UserID:       dataOnDB.UserRegistrationDetail.UserID.String,
			Password:     dataOnDB.UserRegistrationDetail.Password.String,
			Email:        dataOnDB.UserRegistrationDetail.Email.String,
			ClientID:     dataOnDB.PKCEClientMapping.ClientID.String,
			AuthUserID:   dataOnDB.UserRegistrationDetail.AuthUserID.Int64,
		}),
	}

	err = common.UserAuthService.HitResendOtpVerificationEmailToAuth(auth, &contextModel, contextModel.AuthAccessTokenModel.ClientID)
	service.LogMessage(fmt.Sprintf(`Success Auth Send Email`), 200)
	return
}

func (input resendOTPService) convertDTOToModel(inputStruct in.UserVerificationRequest) repository.UserRegistrationDetailModel {
	return repository.UserRegistrationDetailModel{
		ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
		UserID:       sql.NullString{String: inputStruct.UserID},
		UniqueID1:    sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:    sql.NullString{String: inputStruct.UniqueID2},
		Email:        sql.NullString{String: inputStruct.Email},
		NoTelp:       sql.NullString{String: fmt.Sprintf(`%s-%s`, inputStruct.CountryCode, inputStruct.Phone)},
		Status:       sql.NullString{String: constanta.StatusActive},
		AuthUserID:   sql.NullInt64{Int64: inputStruct.AuthUserID},
	}
}

func (input resendOTPService) validateResendOTP(inputStruct *in.UserVerificationRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateResendOTP()
}

func (input resendOTPService) createModelForUpdateUserVerification(timeNow time.Time, inputModel repository.UserRegistrationDetailModel, dbData repository.UserRegistrationDetailMapping) (userVerifyModel repository.UserVerificationModel) {
	var (
		expiredTime = timeNow.Add(24 * time.Hour)
		expiredLong = expiredTime.Unix()
		code        = service.GenerateRandomString(6)
	)

	userVerifyModel = repository.UserVerificationModel{
		ID:             sql.NullInt64{Int64: dbData.UserVerification.ID.Int64},
		EmailCode:      sql.NullString{String: code},
		EmailExpires:   sql.NullInt64{Int64: expiredLong},
		FailedOTPEmail: sql.NullInt64{Int64: 0},
		UpdatedAt:      sql.NullTime{Time: timeNow},
	}

	if int(inputModel.ClientTypeID.Int64) == constanta.ResourceNexstarID {
		userVerifyModel.PhoneCode.String = code
		userVerifyModel.PhoneExpires.Int64 = expiredLong
		userVerifyModel.FailedOTPPhone.Int64 = 0
	}

	if inputModel.NoTelp.String != "" {
		userVerifyModel.PhoneCode.String = code
		userVerifyModel.PhoneExpires.Int64 = expiredLong
		userVerifyModel.FailedOTPPhone.Int64 = 0
	}

	return
}
