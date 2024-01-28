package UserVerificationService

import (
	"database/sql"
	"github.com/Azure/go-autorest/autorest/date"
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
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/common"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input userVerificationService) VerifyingUserService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "VerifyingUserService"
		inputStruct in.UserVerificationRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateVerifyingUser)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doVerifyingUser, func(interface{}, applicationModel.ContextModel) {
		//func Additional
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VERIFICATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input userVerificationService) doVerifyingUser(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.UserVerificationRequest)
	var (
		funcName              = "doVerifyingUser"
		UserOnDB              repository.UserModel
		userRegisOnDB         repository.UserRegistrationDetailModel
		clientMappingOnDB     repository.ClientMappingModel
		pkceClientMappingOnDB repository.PKCEClientMappingModel
		tempAuditData         []repository.AuditSystemModel
		userOnAuth            authentication_response.CheckClientOrUserContent
	)

	// Get User Registration Detail
	userRegisOnDB, err = dao.UserRegistrationDetailDAO.GetUserActiveRegistrationForVerifying(serverconfig.ServerAttribute.DBConnection, repository.UserRegistrationDetailModel{
		UserID:    sql.NullString{String: inputStruct.UserID},
		UniqueID1: sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2: sql.NullString{String: inputStruct.UniqueID2},
	})
	if err.Error != nil {
		return
	}

	if userRegisOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserRegistrationDetailName)
		return
	}

	// Client Validation
	inputStruct.ClientID = userRegisOnDB.ClientID.String
	pkceClientMappingOnDB, clientMappingOnDB, err = input.clientValidation(inputStruct)
	if err.Error != nil {
		return
	}

	if userRegisOnDB.ClientID.String != pkceClientMappingOnDB.ClientID.String {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	// User Validation
	UserOnDB, userOnAuth, err = input.checkUserValidation(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	// Check Verifying User Channel
	tempAuditData, err = input.verifyingUser(tx, inputStruct, userRegisOnDB, UserOnDB, userOnAuth, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, tempAuditData...)

	// Activating User
	if err = dao.UserDAO.UpdateUserStatus(tx, repository.UserModel{
		AuthUserID:    UserOnDB.AuthUserID,
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		Status:        sql.NullString{String: constanta.StatusActive},
	}); err.Error != nil {
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserDAO.TableName, UserOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	// Update User Registration Detail
	if err = dao.UserRegistrationDetailDAO.UpdateAndroidID(tx, repository.UserRegistrationDetailModel{
		ID:            userRegisOnDB.ID,
		AndroidID:     sql.NullString{String: inputStruct.AndroidID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
	}); err.Error != nil {
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserRegistrationDetailDAO.TableName, userRegisOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	result, err = input.getUserVerificationResponse(inputStruct, userRegisOnDB, clientMappingOnDB)
	return
}

func (input userVerificationService) getUserVerificationResponse(inputStruct in.UserVerificationRequest, userRegisOnDB repository.UserRegistrationDetailModel, clientMappingOnDB repository.ClientMappingModel) (output out.UserVerificationResponse, err errorModel.ErrorModel) {
	var (
		parameterValue     []repository.ParameterValueModel
		userRegisAdminOnDB repository.UserRegistrationAdminModel
	)
	// Get Nexmile Parameter
	parameterValue, err = dao.NexmileParameterDAO.GetFieldNexmileParameter(serverconfig.ServerAttribute.DBConnection, repository.NexmileParameterModel{
		ClientID:     clientMappingOnDB.ClientID,
		UniqueID1:    sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:    sql.NullString{String: inputStruct.UniqueID2},
		ClientTypeID: clientMappingOnDB.ClientTypeID,
	})

	if err.Error != nil {
		return
	}

	// Get User Registration Admin
	userRegisAdminOnDB, err = dao.UserRegistrationAdminDAO.GetFieldForNexmileParameter(serverconfig.ServerAttribute.DBConnection, repository.UserRegistrationAdminModel{
		UniqueID1:       sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:       sql.NullString{String: inputStruct.UniqueID2},
		ClientTypeID:    sql.NullInt64{Int64: inputStruct.ClientTypeID},
		ClientMappingID: clientMappingOnDB.ID,
	})
	if err.Error != nil {
		return
	}

	if userRegisAdminOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, "getUserVerificationResponse", constanta.UserRegistrationAdmin)
		return
	}

	// Convert to Response Struct
	output = input.convertToResponseUserVerification(parameterValue, userRegisOnDB, userRegisAdminOnDB)
	return
}

func (input userVerificationService) convertToResponseUserVerification(parameterValue []repository.ParameterValueModel, userRegistOnDB repository.UserRegistrationDetailModel, userRegisAdminOnDB repository.UserRegistrationAdminModel) (output out.UserVerificationResponse) {
	var parameterValueResp []out.ParameterValueResponse
	output = out.UserVerificationResponse{
		CompanyID:        userRegistOnDB.UniqueID1.String,
		BranchID:         userRegistOnDB.UniqueID2.String,
		CompanyName:      userRegisAdminOnDB.CompanyName.String,
		ProductValidFrom: date.Date{Time: userRegistOnDB.ProductValidFrom.Time},
		ProductValidThru: date.Date{Time: userRegistOnDB.ProductValidThru.Time},
		LicenseStatus:    userRegistOnDB.LicenseStatus.Int64,
		AdminPassword:    userRegisAdminOnDB.PasswordAdmin.String,
		AdminUsername:    userRegisAdminOnDB.UserAdmin.String,
		MaxOfflineDays:   userRegistOnDB.MaxOfflineDays.Int64,
	}

	for _, item := range parameterValue {
		parameterValueResp = append(parameterValueResp, out.ParameterValueResponse{
			ParameterID:    item.ParameterID.String,
			ParameterValue: item.ParameterValue.String,
		})
	}

	output.ParameterValue = parameterValueResp
	return
}

func (input userVerificationService) verifyingUser(tx *sql.Tx, inputStruct in.UserVerificationRequest, userRegistOnDB repository.UserRegistrationDetailModel,
	UserOnDB repository.UserModel, userOnAuth authentication_response.CheckClientOrUserContent, contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		userVerificationOnDB repository.UserVerificationModel
	)

	userVerificationModel := repository.UserVerificationModel{
		UserRegistrationDetailID: userRegistOnDB.ID,
		Email:                    sql.NullString{String: inputStruct.Email},
		EmailCode:                sql.NullString{String: inputStruct.OTP},
		Phone:                    sql.NullString{String: inputStruct.Phone},
		PhoneCode:                sql.NullString{String: inputStruct.OTP},
	}

	if userOnAuth.AdditionalInformation.UserStatus == 1 {
		// Verify Nextrac
		// Get User Verification
		userVerificationOnDB, err = dao.UserVerificationDAO.GetUserVerificationForVerifying(serverconfig.ServerAttribute.DBConnection, userVerificationModel)
		if err.Error != nil {
			return
		}

		// Check Expired
		err = input.verifyingOTPFromNextrac(userVerificationOnDB, inputStruct, contextModel, timeNow)
		if err.Error != nil {
			return
		}

		// Delete OTP on User Verification
		err = dao.UserVerificationDAO.HardDeleteUserVerification(tx, userVerificationOnDB)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.UserVerificationDAO.TableName, userVerificationOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)
	} else {
		// Verify Auth
		err = input.verifyingByAuth(inputStruct, UserOnDB, contextModel)
		if err.Error != nil {
			return
		}
	}

	return
}

func (input userVerificationService) verifyingOTPFromNextrac(userVerificationOnDB repository.UserVerificationModel, inputStruct in.UserVerificationRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var (
		expiredTime time.Time
		funcName    = "verifyingOTPFromNextrac"
	)

	if util.IsStringEmpty(inputStruct.Email) {
		// Validate OTP
		if userVerificationOnDB.FailedOTPPhone.Int64 < 3 {
			if userVerificationOnDB.PhoneCode.String != inputStruct.OTP {
				err = input.updateOTPFailed(repository.UserVerificationModel{
					ID:           userVerificationOnDB.ID,
					PhoneExpires: sql.NullInt64{Int64: timeNow.Add(24 * time.Hour).Unix()},
				}, contextModel, timeNow)
				if err.Error != nil {
					return
				}
				err = errorModel.GenerateRequestError(input.FileName, funcName, constanta.WrongOTPCode)
				return
			}

			expiredTime = time.Unix(userVerificationOnDB.PhoneExpires.Int64, 0)
			if timeNow.After(expiredTime) {
				err = errorModel.GenerateRequestError(input.FileName, funcName, constanta.ExpiredActivationCode)
				return
			}
		} else {
			err = errorModel.GenerateRequestError(input.FileName, funcName, constanta.ExpiredActivationCode)
			return
		}

	} else {
		// Validate OTP
		if userVerificationOnDB.FailedOTPEmail.Int64 < 3 {
			if userVerificationOnDB.EmailCode.String != inputStruct.OTP {
				err = input.updateOTPFailed(repository.UserVerificationModel{
					ID:           userVerificationOnDB.ID,
					EmailExpires: sql.NullInt64{Int64: timeNow.Add(24 * time.Hour).Unix()},
				}, contextModel, timeNow)
				if err.Error != nil {
					return
				}
				err = errorModel.GenerateRequestError(input.FileName, funcName, constanta.WrongOTPCode)
				return
			}

			expiredTime = time.Unix(userVerificationOnDB.EmailExpires.Int64, 0)
			if timeNow.After(expiredTime) {
				err = errorModel.GenerateRequestError(input.FileName, funcName, constanta.ExpiredActivationCode)
				return
			}
		} else {
			err = errorModel.GenerateRequestError(input.FileName, funcName, constanta.ExpiredActivationCode)
			return
		}
	}
	return
}

func (input userVerificationService) updateOTPFailed(userVerificationOnDB repository.UserVerificationModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	//update user verification failed otp
	userVerificationOnDB.UpdatedAt.Time = timeNow
	userVerificationOnDB.UpdatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
	userVerificationOnDB.UpdatedClient.String = contextModel.AuthAccessTokenModel.ClientID
	err = dao.UserVerificationDAO.UpdateFailedOTP(serverconfig.ServerAttribute.DBConnection, userVerificationOnDB)
	return
}

func (input userVerificationService) verifyingByAuth(inputStruct in.UserVerificationRequest, UserOnDB repository.UserModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName = "VerifyingUserService.go"
		funcName = "verifyingByAuth"
		message  string
	)

	//--- Verifying Auth
	if util.IsStringEmpty(inputStruct.Email) {
		userActivationReq := in.UserActivationPhoneRequest{
			UserID:  UserOnDB.AuthUserID.Int64,
			Phone:   inputStruct.Phone,
			OTPCode: inputStruct.OTP,
		}

		err = common.UserAuthService.HitActivationPhoneToAuth(userActivationReq, contextModel, "")
		if err.Error != nil {
			service.LogMessage("error activation phone", 200)
			service.LogMessage(err.Error.Error(), 200)
			service.LogMessage(err.CausedBy.Error(), 200)

			if err.Error.Error() == constanta.ExpiredActivationCodeError || err.CausedBy.Error() == constanta.ExpiredActivationCodeMessage {
				message = util2.GenerateConstantaI18n(constanta.ExpiredActivationCode, contextModel.AuthAccessTokenModel.Locale, nil)
				err = errorModel.GenerateErrorCustomActivationCode(fileName, funcName, message)
			} else if err.Error.Error() == constanta.FormatPhoneActivationCodeError || err.CausedBy.Error() == constanta.FormatPhoneActivationCodeMessage {
				message = util2.GenerateConstantaI18n(constanta.WrongFormatPhone, contextModel.AuthAccessTokenModel.Locale, nil)
				err = errorModel.GenerateErrorCustomActivationCode(fileName, funcName, message)
			} else if err.Error.Error() == constanta.EmailPhoneUnknownDataError || err.CausedBy.Error() == constanta.PhoneUnknownMessage {
				message = util2.GenerateConstantaI18n(constanta.WrongPhoneData, contextModel.AuthAccessTokenModel.Locale, nil)
				err = errorModel.GenerateErrorCustomActivationCode(fileName, funcName, message)
			} else if err.Error.Error() == constanta.OTPCodeActivationError || err.CausedBy.Error() == constanta.OTPCodeActivationMessage {
				message = util2.GenerateConstantaI18n(constanta.WrongOTPCode, contextModel.AuthAccessTokenModel.Locale, nil)
				err = errorModel.GenerateErrorCustomActivationCode(fileName, funcName, message)
			}
			return
		}
	} else {
		userActivationReq := in.UserActivationRequest{
			UserID:    UserOnDB.AuthUserID.Int64,
			Email:     inputStruct.Email,
			EmailCode: inputStruct.OTP,
		}

		err = common.UserAuthService.HitActivationEmailToAuth(userActivationReq, contextModel, "")
		if err.Error != nil {
			service.LogMessage("error activation email", 200)
			service.LogMessage(err.Error.Error(), 200)
			service.LogMessage(err.CausedBy.Error(), 200)

			if err.Error.Error() == constanta.ActivationCodeError || err.CausedBy.Error() == constanta.ActivationCodeMessage {
				message = util2.GenerateConstantaI18n(constanta.WrongActivationCode, contextModel.AuthAccessTokenModel.Locale, nil)
				err = errorModel.GenerateErrorCustomActivationCode(fileName, funcName, message)
			} else if err.Error.Error() == constanta.EmailPhoneUnknownDataError || err.CausedBy.Error() == constanta.EmailUnknownMessage {
				message = util2.GenerateConstantaI18n(constanta.WrongEmailData, contextModel.AuthAccessTokenModel.Locale, nil)
				err = errorModel.GenerateErrorCustomActivationCode(fileName, funcName, message)
			} else if err.Error.Error() == constanta.ExpiredActivationCodeError || err.CausedBy.Error() == constanta.ExpiredActivationCodeMessage {
				message = util2.GenerateConstantaI18n(constanta.ExpiredActivationCode, contextModel.AuthAccessTokenModel.Locale, nil)
				err = errorModel.GenerateErrorCustomActivationCode(fileName, funcName, message)
			}
			return
		}
	}

	return
}

func (input userVerificationService) checkUserValidation(inputStruct in.UserVerificationRequest, contextModel *applicationModel.ContextModel) (userOnDB repository.UserModel,
	userOnAuth authentication_response.CheckClientOrUserContent, err errorModel.ErrorModel) {
	var (
		funcName = "checkUserValidation"
	)

	// Check user on DB
	userOnDB, err = dao.UserDAO.GetAuthUserByClientIDForVerifying(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
		ClientID: sql.NullString{String: inputStruct.ClientID},
	})
	if err.Error != nil {
		return
	}

	if userOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.User)
		return
	}

	// Check user on Auth
	userOnAuthResp, err := service.CheckClientOrUserInAuth(authentication_request.CheckClientOrUser{ClientID: inputStruct.ClientID}, contextModel)
	if err.Error != nil {
		return
	}

	userOnAuth = userOnAuthResp.Nexsoft.Payload.Data.Content

	if !userOnAuth.IsExist {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserAuth)
		return
	}

	return
}

func (input userVerificationService) clientValidation(inputStruct in.UserVerificationRequest) (PKCEClientMappingOnDB repository.PKCEClientMappingModel, clientMappingOnDB repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "clientValidation"

	PKCEModel := repository.PKCEClientMappingModel{
		ClientID:     sql.NullString{String: inputStruct.ClientID},
		CompanyID:    sql.NullString{String: inputStruct.UniqueID1},
		BranchID:     sql.NullString{String: inputStruct.UniqueID2},
		ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
	}

	// Validate Client Type
	clientTypeOnDB, err := dao.ClientTypeDAO.ValidateClientTypeByID(serverconfig.ServerAttribute.DBConnection, repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
	})
	if err.Error != nil {
		return
	}

	if clientTypeOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientType)
		return
	}

	// Validate PKCE Client Mapping
	PKCEClientMappingOnDB, clientMappingOnDB, err = dao.PKCEClientMappingDAO.GetPKCEClientWithClientIDAndUniqueID(serverconfig.ServerAttribute.DBConnection, PKCEModel)
	if err.Error != nil {
		return
	}

	if PKCEClientMappingOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateRequestError(input.FileName, funcName, "ERROR_REQUEST_CLIENT_CREDENTIAL")
		return
	}

	return
}

func (input userVerificationService) validateVerifyingUser(inputStruct *in.UserVerificationRequest) errorModel.ErrorModel {
	return inputStruct.ValidateVerifyingUser()
}
