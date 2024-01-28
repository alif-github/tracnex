package RegistrationNamedUserService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_request"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/UserService/CRUDUserService"
	common2 "nexsoft.co.id/nextrac2/service/common"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input registrationNamedUserService) RegisterOrRenewLicenseNamedUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = input.FileName
		inputStruct in.RegisterNamedUserRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validationRegisterOrRenewLicenseNamedUser)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doRegisterOrRenewLicenseNamedUser, func(i interface{}, model applicationModel.ContextModel) {
		// code additional information
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_USER_REGISTRATION_DETAIL", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) validationRegisterOrRenewLicenseNamedUser(inputStruct *in.RegisterNamedUserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateRegisterOrRenewLicenseNamedUser()
}

func (input registrationNamedUserService) doRegisterOrRenewLicenseNamedUser(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var funcName = "doRegisterOrRenewLicenseNamedUser"
	var (
		inputStruct      = inputStructInterface.(in.RegisterNamedUserRequest)
		userNexmileModel = input.createModelRegistrationNamedUserLicense(inputStruct, contextModel, timeNow)
	)

	clientMappingModel := repository.ClientMappingModel{
		ClientID:  sql.NullString{String: inputStruct.ParentClientID},
		CompanyID: sql.NullString{String: inputStruct.UniqueID1},
		BranchID:  sql.NullString{String: inputStruct.UniqueID2},
	}

	resultClientMapping, err := dao.ClientMappingDAO.GetParentClientMappingForClientValidation(serverconfig.ServerAttribute.DBConnection, clientMappingModel)
	if err.Error != nil {
		return
	}

	if contextModel.AuthAccessTokenModel.ClientID != resultClientMapping.ClientID.String {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	// Step 1 - Cek License Kuota
	_, err = input.checkLicenseQuota(&userNexmileModel)
	if err.Error != nil {
		return
	}

	lastUserLicenseID := userNexmileModel.UserLicenseID

	clientTypeRequest := input.setModelForGetClientType(inputStruct)
	clientTypeOnDB, err := dao.ClientTypeDAO.GetClientTypeByID(serverconfig.ServerAttribute.DBConnection, clientTypeRequest)
	if err.Error != nil {
		return
	}

	inputStruct.ClientType = clientTypeOnDB.ClientType.String

	// Step 2 - Mapping user baru atau perpanjang berdasarkan data di auth dan data di nextrac
	output, err = input.mappingCreateNamedUserOrRenewNamedUser(tx, inputStruct, resultClientMapping, contextModel, timeNow, &dataAudit, lastUserLicenseID.Int64)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) setModelForGetUserRegistrationDetail(inputStruct in.RegisterNamedUserRequest) repository.UserRegistrationDetailModel {
	return repository.UserRegistrationDetailModel{
		UserID:     sql.NullString{String: inputStruct.UserID},
		UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
		AuthUserID: sql.NullInt64{Int64: inputStruct.AuthUserID},
	}
}

func (input registrationNamedUserService) setModelForGetClientType(inputStruct in.RegisterNamedUserRequest) repository.ClientTypeModel {
	return repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
	}
}

func (input registrationNamedUserService) updateNamedUser(tx *sql.Tx, inputStruct in.RegisterNamedUserRequest, authUser authentication_response.UserContent, contextModel *applicationModel.ContextModel, userRegDetailOnDB repository.UserRegistrationDetailModel, timeNow time.Time, lastUserLicenseID int64) (output out.RegisterOrRenewLicenseUserResponse, err errorModel.ErrorModel) {
	var (
		funcName = "updateNamedUser"
		userOnDB repository.UserModel
	)

	// Step 1 - Update status user menjadi active
	err = dao.UserRegistrationDetailDAO.UpdateRenewUserRegistrationDetail(tx, repository.UserRegistrationDetailModel{
		ID:               userRegDetailOnDB.ID,
		UserLicenseID:    sql.NullInt64{Int64: lastUserLicenseID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UserID:           sql.NullString{String: inputStruct.UserID},
		SalesmanCategory: sql.NullString{String: inputStruct.SalesmanCategory},
		SalesmanID:       sql.NullString{String: inputStruct.SalesmanID},
		Email:            sql.NullString{String: inputStruct.Email},
		NoTelp:           sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.NoTelp},
		Status:           sql.NullString{String: constanta.StatusActive},
	})
	if err.Error != nil {
		fmt.Println("Error UpdateRenewUserRegistrationDetail : ", err.Error.Error())
		return
	}

	// Step 2 - Update Untuk menambahkan total user aktif di user license
	err = dao.UserLicenseDAO.UpdateTotalActivatedUserLicense(tx, repository.UserLicenseModel{
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		ID:            sql.NullInt64{Int64: lastUserLicenseID},
	})
	if err.Error != nil {
		fmt.Println("Error UpdateTotalActivatedUserLicense : ", err.Error.Error())
		return
	}

	// Step 2.2 - Update Status User Kembali Pending
	err = dao.UserDAO.UpdateRenewUserStatus(tx, repository.UserModel{
		AuthUserID:    sql.NullInt64{Int64: userRegDetailOnDB.AuthUserID.Int64},
		FirstName:     sql.NullString{String: inputStruct.Firstname},
		LastName:      sql.NullString{String: inputStruct.Lastname},
		Email:         sql.NullString{String: inputStruct.Email},
		Phone:         sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.NoTelp},
		Status:        sql.NullString{String: constanta.PendingOnApproval},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
	})
	if err.Error != nil {
		fmt.Println("Error UpdateRenewUserStatus : ", err.Error.Error())
		return
	}

	// Step 3 - Menambahkan resource trac
	err = CRUDUserService.HitAuthenticationServerForAddResourceToUserAuth(authentication_response.UserContent{ClientID: inputStruct.ClientID}, contextModel, config.ApplicationConfiguration.GetServerResourceID())
	if err.Error != nil {
		fmt.Println("Error HitAuthenticationServerForAddResourceToUserAuth : ", err.Error.Error())
		return
	}

	// Data send Email
	dataMapping := repository.UserRegistrationDetailMapping{
		User: repository.UserModel{
			FirstName: sql.NullString{String: inputStruct.Firstname},
		},
		UserRegistrationDetail: repository.UserRegistrationDetailModel{
			UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
			UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
			SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
			UserID:     sql.NullString{String: inputStruct.UserID},
			Password:   sql.NullString{String: inputStruct.Password},
			Email:      sql.NullString{String: inputStruct.Email},
			NoTelp:     sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.NoTelp},
			AuthUserID: sql.NullInt64{Int64: userRegDetailOnDB.AuthUserID.Int64},
		},
		PKCEClientMapping: repository.PKCEClientMappingModel{
			CompanyName:  sql.NullString{String: inputStruct.CompanyName},
			BranchName:   sql.NullString{String: inputStruct.BranchName},
			ClientID:     sql.NullString{String: inputStruct.ClientID},
			ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
		},
	}

	var linkEksternal string
	if inputStruct.ClientType == constanta.Nexmile {
		linkEksternal = config.ApplicationConfiguration.GetNexmile().Host + config.ApplicationConfiguration.GetNexmile().PathRedirect.ActivationUser
	} else if inputStruct.ClientType == constanta.Nexstar {
		linkEksternal = config.ApplicationConfiguration.GetNexstar().Host + config.ApplicationConfiguration.GetNexstar().PathRedirect.ActivationUser
	} else if inputStruct.ClientType == constanta.Nextrade {
		linkEksternal = config.ApplicationConfiguration.GetNextrade().Host + config.ApplicationConfiguration.GetNextrade().PathRedirect.ActivationUser
	}

	var isSendEmailFromAuth bool
	if util.IsStringEmpty(authUser.Email) {
		fmt.Println("Masuk send email : ")
		if util.IsStringEmpty(inputStruct.Email) {
			err = errorModel.GenerateEmailEmptyAuthNexstarForResendVerification(input.FileName, funcName)
			return
		}

		userOnDB, err = dao.UserDAO.ViewDetailUserForRegisterNamedUser(serverconfig.ServerAttribute.DBConnection, repository.UserModel{AuthUserID: userRegDetailOnDB.AuthUserID})
		if err.Error != nil {
			fmt.Println("ViewDetailUserForRegisterNamedUser : ", err.Error.Error())
			return
		}

		authUpdateEmailReq := in.UserRequest{
			FirstName:   userOnDB.FirstName.String,
			Username:    userOnDB.Username.String,
			Email:       inputStruct.Email,
			Locale:      contextModel.AuthAccessTokenModel.Locale,
			AuthUserID:  userRegDetailOnDB.AuthUserID.Int64,
			LastName:    userOnDB.LastName.String,
			CountryCode: inputStruct.CountryCode,
		}

		dataMapping = repository.UserRegistrationDetailMapping{
			User: repository.UserModel{
				FirstName: sql.NullString{String: inputStruct.Firstname},
			},
			UserRegistrationDetail: repository.UserRegistrationDetailModel{
				UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
				UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
				SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
				UserID:     sql.NullString{String: inputStruct.UserID},
				Password:   sql.NullString{String: inputStruct.Password},
				Email:      sql.NullString{String: inputStruct.Email},
			},
			PKCEClientMapping: repository.PKCEClientMappingModel{
				CompanyName:  sql.NullString{String: inputStruct.CompanyName},
				BranchName:   sql.NullString{String: inputStruct.BranchName},
				ClientID:     sql.NullString{String: inputStruct.ClientID},
				ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
			},
		}

		_, err = input.updateUserToAuthenticationServer(authUpdateEmailReq, dataMapping, linkEksternal, contextModel)
		if err.Error != nil {
			return
		}

		isSendEmailFromAuth = true
	}

	// Step 4 - Generate OTP Nextrac untuk perpanjang named user
	var otpModel repository.UserVerificationModel
	if inputStruct.ClientType == constanta.Nexmile || inputStruct.ClientType == constanta.Nextrade {
		otpModel = input.setModelForRenewNexMileLicense(userRegDetailOnDB, contextModel, timeNow)
	} else if inputStruct.ClientType == constanta.Nexstar {
		otpModel = input.setModelForRenewNexStarLicense(userRegDetailOnDB, contextModel, timeNow)
	}

	// insert otp ke tabel user_verification
	if isSendEmailFromAuth == false {
		_, err = dao.UserVerificationDAO.InsertOTPUser(tx, otpModel)
		if err.Error != nil {
			fmt.Println("InsertOTPUser : ", err.Error.Error())
			err = input.CheckOTPDuplicateError(err)
			fmt.Println("CheckOTPDuplicateError : ", err.Error.Error())

			return
		}

		err = input.SendMessageToEmail(dataMapping, constanta.SubjectActivationEmail, " ", otpModel, linkEksternal, input.linkQueryEmail)
		if err.Error != nil {
			fmt.Println("SendMessageToEmail : ", err.Error.Error())
			return
		}
	}

	//--- Send Gro Chat Message
	if inputStruct.ClientType == constanta.Nexstar {
		err = input.sendMessageToGrochat(dataMapping, otpModel, contextModel, linkEksternal)
		if err.Error != nil {
			fmt.Println("sendMessageToGrochat : ", err.Error.Error())
			return
		}
	}

	output.UserRegistrationDetailId = userRegDetailOnDB.ID.Int64
	output.ClientId = userRegDetailOnDB.ClientID.String
	output.AuthUserId = userRegDetailOnDB.AuthUserID.Int64
	return
}

func (input registrationNamedUserService) CheckOTPDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "user_registration_detail_unique_userid_uniqueid_1_2") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.OTP)
		}
	}

	return err
}

func (input registrationNamedUserService) sendEmailForActivationLicense(inputStruct in.RegisterNamedUserRequest, otpModel repository.UserVerificationModel, linkEksternal string) (err errorModel.ErrorModel) {
	var (
		funcName    = "sendEmailForActivationLicense"
		templateSrc = config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath + config.ApplicationConfiguration.GetDataDirectory().Template + "/template-activation-nexmile-nexstar.html"
		errorEmail  error
		clientType  string
	)

	switch inputStruct.ClientTypeID {
	case constanta.ResourceNexmileID:
		clientType = constanta.Nexmile
	case constanta.ResourceNexstarID:
		clientType = constanta.Nexstar
	case constanta.ResourceNextradeID:
		clientType = constanta.Nextrade
	default:
		message := fmt.Sprintf(`%s %s %s`, constanta.Nexmile, constanta.Nexstar, constanta.Nextrade)
		err = errorModel.GenerateErrorEksternalClientTypeMustHave(input.FileName, funcName, message)
		return
	}

	reqEmail := util2.NewRequestMail("Activation " + clientType)
	reqEmail.To = []mail.Address{
		{
			Address: inputStruct.Email,
		},
	}

	linkEksternal += fmt.Sprint("?client_id=", inputStruct.ClientID, "&unique_id_1=",
		inputStruct.UniqueID1, "&unique_id_2=", inputStruct.UniqueID2, "&salesman_id=", inputStruct.SalesmanID, "&userid=",
		inputStruct.UserID, "&otp=", otpModel.EmailCode.String, "&email=", inputStruct.Email)

	// Generate email template
	emailTemplateData := util2.TemplateDataActivationUserNexMile{
		Name:       inputStruct.Firstname,
		ClientType: inputStruct.ClientType,
		UniqueID1:  inputStruct.UniqueID1,
		UniqueID2:  inputStruct.UniqueID2,
		SalesmanID: inputStruct.SalesmanID,
		UserID:     inputStruct.UserID,
		Password:   inputStruct.Password,
		OTP:        otpModel.EmailCode.String,
		Link:       linkEksternal,
	}

	// Get Template File
	templateSrc, errorEmail = filepath.Abs(templateSrc)
	if errorEmail != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorEmail)
		service.LogMessageWithErrorModel(err)
		return
	}

	// Generate Email Template
	errorEmail = reqEmail.GenerateEmailTemplate(templateSrc, emailTemplateData, constanta.IndonesianLanguage)
	if errorEmail != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorEmail)
		service.LogMessageWithErrorModel(err)
		return
	}

	// Send Email
	errorEmail = reqEmail.SendEmail()
	if errorEmail != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorEmail)
		service.LogMessageWithErrorModel(err)
		return
	}

	return
}

func (input registrationNamedUserService) setModelForRenewNexMileLicense(repo repository.UserRegistrationDetailModel, contextModel *applicationModel.ContextModel, timeNow time.Time) repository.UserVerificationModel {
	var otp = service.GenerateRandomString(6)
	var expires = timeNow.Add(24 * time.Hour).Unix()

	return repository.UserVerificationModel{
		UserRegistrationDetailID: repo.ID,
		Email:                    repo.Email,
		EmailCode:                sql.NullString{String: otp},
		EmailExpires:             sql.NullInt64{Int64: expires},
		UpdatedBy:                sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:            sql.NullString{String: constanta.SystemClient},
		UpdatedAt:                sql.NullTime{Time: timeNow},
		CreatedBy:                sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:            sql.NullString{String: constanta.SystemClient},
		CreatedAt:                sql.NullTime{Time: timeNow},
	}
}

func (input registrationNamedUserService) setModelForRenewNexStarLicense(repo repository.UserRegistrationDetailModel, contextModel *applicationModel.ContextModel, timeNow time.Time) repository.UserVerificationModel {
	var otp = service.GenerateRandomString(6)
	var expires = timeNow.Add(24 * time.Hour).Unix()

	return repository.UserVerificationModel{
		UserRegistrationDetailID: repo.ID,
		Email:                    repo.Email,
		EmailCode:                sql.NullString{String: otp},
		EmailExpires:             sql.NullInt64{Int64: expires},
		Phone:                    sql.NullString{String: constanta.IndonesianCodeNumber + "-" + repo.NoTelp.String},
		PhoneCode:                sql.NullString{String: otp},
		PhoneExpires:             sql.NullInt64{Int64: expires},
		UpdatedBy:                sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:            sql.NullString{String: constanta.SystemClient},
		UpdatedAt:                sql.NullTime{Time: timeNow},
		CreatedBy:                sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:            sql.NullString{String: constanta.SystemClient},
		CreatedAt:                sql.NullTime{Time: timeNow},
	}
}

func (input registrationNamedUserService) mappingCreateNamedUserOrRenewNamedUser(tx *sql.Tx, inputStruct in.RegisterNamedUserRequest, clientMappingOnDB repository.ClientMappingModel, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel, lastUserLicenseID int64) (output out.RegisterOrRenewLicenseUserResponse, err errorModel.ErrorModel) {
	var (
		funcName                    = "mappingCreateNamedUserOrRenewNamedUser"
		authUser                    authentication_response.UserContent
		userAuthByEmailAndPhone     authentication_response.UserAuthenticationResponse
		userAuthByEmail             authentication_response.UserAuthenticationResponse
		userAuthByPhone             authentication_response.UserAuthenticationResponse
		dataUserAuthByPhone         authentication_response.UserContent
		dataUserAuthByEmail         authentication_response.UserContent
		userRegDetailOnDB           repository.UserRegistrationDetailModel
		userRegDetailOnDbByUserAuth repository.UserRegistrationDetailModel
		isUpdate                    bool
	)

	if inputStruct.ClientType == constanta.Nexmile || inputStruct.ClientType == constanta.Nextrade {
		if util.IsStringEmpty(inputStruct.Email) {
			err = errorModel.GenerateMandatoryField(constanta.Email, input.FileName, funcName)
			return
		}
	}

	if inputStruct.ClientType == constanta.Nexstar {
		if util.IsStringEmpty(inputStruct.NoTelp) {
			err = errorModel.GenerateMandatoryField(constanta.Phone, input.FileName, funcName)
			return
		}
	}

	// Step 1 - Cek Kombinasi Email dan Phone di Authentication Server
	if !util.IsStringEmpty(inputStruct.Email) && !util.IsStringEmpty(inputStruct.NoTelp) {
		userAuthByEmailAndPhone, err = CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Email: inputStruct.Email, Phone: constanta.IndonesianCodeNumber + "-" + inputStruct.NoTelp}, contextModel)
		if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
			return
		}
	}

	if userAuthByEmailAndPhone.Nexsoft.Payload.Data.Content.UserID < 1 {
		if !util.IsStringEmpty(inputStruct.Email) {
			userAuthByEmail, err = CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Email: inputStruct.Email}, contextModel)
			dataUserAuthByEmail = userAuthByEmail.Nexsoft.Payload.Data.Content
			if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
				return
			}
		}

		if !util.IsStringEmpty(inputStruct.NoTelp) {
			userAuthByPhone, err = CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Phone: inputStruct.CountryCode + "-" + inputStruct.NoTelp}, contextModel)
			dataUserAuthByPhone = userAuthByPhone.Nexsoft.Payload.Data.Content
			if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
				return
			}
		}

		if dataUserAuthByEmail.UserID > 0 && dataUserAuthByPhone.UserID > 0 {
			err = errorModel.GenerateBothEmailAndPhoneAlreadyRegisteredAuth(input.FileName, funcName)
			return
		}

		if dataUserAuthByEmail.UserID > 0 {
			userAuthByEmailAndPhone = userAuthByEmail
		} else if dataUserAuthByPhone.UserID > 0 {
			userAuthByEmailAndPhone = userAuthByPhone
		}
	}

	authUser = userAuthByEmailAndPhone.Nexsoft.Payload.Data.Content

	if authUser.UserID > 0 && inputStruct.AuthUserID > 0 {
		if authUser.UserID != inputStruct.AuthUserID {
			err = errorModel.GenerateMismatchAuthUserID(input.FileName, funcName)
			return
		}
	}

	if authUser.UserID > 0 {
		inputStruct.AuthUserID = authUser.UserID
	}

	userRegDetailOnDB, err = dao.UserRegistrationDetailDAO.GetUserForCheckRegistrationNamedUserOrRenewNamedUser(serverconfig.ServerAttribute.DBConnection, input.setModelForGetUserRegistrationDetail(inputStruct))
	if err.Error != nil {
		return
	}

	if inputStruct.AuthUserID != 0 {
		if userRegDetailOnDB.Status.String == constanta.StatusActive {
			err = errorModel.GenerateUserStatusActive(input.FileName, funcName)
			return
		}
	} else {
		if userRegDetailOnDB.Status.String == constanta.StatusActive {
			err = errorModel.GenerateUserIdHasBeenActivated(input.FileName, funcName)
			return
		}
	}

	// Step 3 - Mapping User, kondisi jika ditemukan user di nextrac
	if userRegDetailOnDB.ID.Int64 != 0 {
		inputStruct.ClientID = userRegDetailOnDB.ClientID.String
		userRegDetailOnDB.UniqueID1.String = inputStruct.UniqueID1
		userRegDetailOnDB.UniqueID1.String = inputStruct.UniqueID1
		userRegDetailOnDB.NoTelp.String = strings.ReplaceAll(userRegDetailOnDB.NoTelp.String, "+62-", "")

		if util.IsStringEmpty(userRegDetailOnDB.NoTelp.String) {
			userRegDetailOnDB.NoTelp.String = inputStruct.NoTelp
		}

		// proses update user
		if inputStruct.AuthUserID != 0 {
			output, err = input.updateNamedUser(tx, inputStruct, authUser, contextModel, userRegDetailOnDB, timeNow, lastUserLicenseID)
			if err.Error != nil {
				return
			}

			return
		} else {
			// kondisi jika ditemukan data di nextrac by user_id dan ditemukan user di auth
			if authUser.UserID != 0 {
				userRegDetailOnDbByUserAuth, err = dao.UserRegistrationDetailDAO.GetUserForCheckRegistrationNamedUserOrRenewNamedUser(serverconfig.ServerAttribute.DBConnection, repository.UserRegistrationDetailModel{
					UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
					UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
					AuthUserID: sql.NullInt64{Int64: authUser.UserID},
				})
				if err.Error != nil {
					return
				}

				if userRegDetailOnDbByUserAuth.Status.String == constanta.StatusActive {
					err = errorModel.GenerateUserIdHasBeenActivated(input.FileName, funcName)
					return
				}

				// data updated on db
				inputStruct.ClientID = userRegDetailOnDbByUserAuth.ClientID.String
				userRegDetailOnDbByUserAuth.UniqueID1.String = inputStruct.UniqueID1
				userRegDetailOnDbByUserAuth.UniqueID1.String = inputStruct.UniqueID1
				userRegDetailOnDbByUserAuth.NoTelp.String = strings.ReplaceAll(userRegDetailOnDB.NoTelp.String, "+62-", "")
				if util.IsStringEmpty(userRegDetailOnDbByUserAuth.NoTelp.String) {
					userRegDetailOnDbByUserAuth.NoTelp.String = inputStruct.NoTelp
				}

				// update
				fmt.Println(userRegDetailOnDbByUserAuth)
				if userRegDetailOnDbByUserAuth.ID.Int64 != 0 {
					output, err = input.updateNamedUser(tx, inputStruct, authUser, contextModel, userRegDetailOnDbByUserAuth, timeNow, lastUserLicenseID)
					if err.Error != nil {
						return
					}

					return
				}

				// kondisi jika request adalah user baru, berbeda dengan yg ditemukan di nextrac by user_id
				if inputStruct.ClientType == constanta.Nexmile || inputStruct.ClientType == constanta.Nextrade {
					output, _, _, err = input.createNamedUserForNexmileOrNextradeFromAuth(tx, inputStruct, authUser, clientMappingOnDB, contextModel, timeNow, dataAudit)
					if err.Error != nil {
						return
					}

					return
				}

				if inputStruct.ClientType == constanta.Nexstar {
					output, _, _, err = input.createNamedUserForNexstar(tx, inputStruct, authUser, contextModel, timeNow, dataAudit)
					if err.Error != nil {
						return
					}

					return
				}
			} else {
				output, err = input.createNamedUserForNexmileOrNextrade(tx, inputStruct, clientMappingOnDB, contextModel, timeNow, dataAudit)
				if err.Error != nil {
					return
				}

			}
		}

		err = errorModel.GenerateNonErrorModel()

		return
	} else if userRegDetailOnDB.ID.Int64 == 0 && authUser.UserID != 0 {
		if inputStruct.ClientType == constanta.Nexstar {
			output, userRegDetailOnDbByUserAuth, isUpdate, err = input.createNamedUserForNexstar(tx, inputStruct, authUser, contextModel, timeNow, dataAudit)
			if err.Error != nil {
				return
			}

			if isUpdate {
				inputStruct.ClientID = userRegDetailOnDbByUserAuth.ClientID.String
				userRegDetailOnDbByUserAuth.UniqueID1.String = inputStruct.UniqueID1
				userRegDetailOnDbByUserAuth.UniqueID1.String = inputStruct.UniqueID1
				userRegDetailOnDbByUserAuth.NoTelp.String = strings.ReplaceAll(userRegDetailOnDB.NoTelp.String, "+62-", "")
				if util.IsStringEmpty(userRegDetailOnDbByUserAuth.NoTelp.String) {
					userRegDetailOnDbByUserAuth.NoTelp.String = inputStruct.NoTelp
				}

				output, err = input.updateNamedUser(tx, inputStruct, authUser, contextModel, userRegDetailOnDbByUserAuth, timeNow, lastUserLicenseID)
				if err.Error != nil {
					return
				}
			}

		} else if inputStruct.ClientType == constanta.Nexmile || inputStruct.ClientType == constanta.Nextrade {
			output, isUpdate, userRegDetailOnDbByUserAuth, err = input.createNamedUserForNexmileOrNextradeFromAuth(tx, inputStruct, authUser, clientMappingOnDB, contextModel, timeNow, dataAudit)
			if err.Error != nil {
				return
			}

			if isUpdate {
				inputStruct.ClientID = userRegDetailOnDbByUserAuth.ClientID.String
				userRegDetailOnDbByUserAuth.UniqueID1.String = inputStruct.UniqueID1
				userRegDetailOnDbByUserAuth.UniqueID1.String = inputStruct.UniqueID1
				userRegDetailOnDbByUserAuth.NoTelp.String = strings.ReplaceAll(userRegDetailOnDB.NoTelp.String, "+62-", "")
				if util.IsStringEmpty(userRegDetailOnDbByUserAuth.NoTelp.String) {
					userRegDetailOnDbByUserAuth.NoTelp.String = inputStruct.NoTelp
				}

				output, err = input.updateNamedUser(tx, inputStruct, authUser, contextModel, userRegDetailOnDbByUserAuth, timeNow, lastUserLicenseID)
				if err.Error != nil {
					return
				}
			}
		}

		err = errorModel.GenerateNonErrorModel()
		return

	} else {
		output, err = input.createNamedUserForNexmileOrNextrade(tx, inputStruct, clientMappingOnDB, contextModel, timeNow, dataAudit)
		if err.Error != nil {
			return
		}

		err = errorModel.GenerateNonErrorModel()
		return
	}
}

func (input registrationNamedUserService) createNamedUserForNexmileOrNextrade(tx *sql.Tx, inputStruct in.RegisterNamedUserRequest, clientMappingOnDB repository.ClientMappingModel, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (output out.RegisterOrRenewLicenseUserResponse, err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "createNamedUserForNexmileOrNextrade"
	)

	if inputStruct.ClientType == constanta.Nexstar {
		err = errorModel.GenerateUserNexstarNotFoundInAuthentication(fileName, funcName)
		return
	}

	// create username auth
	var usernameNamedUser string
	if inputStruct.ClientType == constanta.Nexmile {
		usernameNamedUser = fmt.Sprintf("nd6_%d_%s", clientMappingOnDB.ID.Int64, service.GenerateRandomNumberToString(8))
	} else if inputStruct.ClientType == constanta.Nextrade {
		usernameNamedUser = fmt.Sprintf("nd6_%d_%s", clientMappingOnDB.ID.Int64, service.GenerateRandomNumberToString(8))
	}

	inputStruct.Username = usernameNamedUser

	responseCheckUsername, errs := CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Username: usernameNamedUser}, contextModel)
	if errs.Error != nil && errs.CausedBy.Error() != constanta.AuthenticationDataNotFound {
		return
	}

	if !util.IsStringEmpty(responseCheckUsername.Nexsoft.Payload.Data.Content.Username) {
		err = errorModel.GenerateUsernameAlreadyUsed(fileName, funcName)
		return
	}

	// add new auth user
	var userRequestAuth = setRequestForAddUserAuth(inputStruct, usernameNamedUser)
	linkEksternal := config.ApplicationConfiguration.GetNexmile().Host + config.ApplicationConfiguration.GetNexmile().PathRedirect.ActivationUser

	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Status = 200
	logModel.Message = "Link External : " + linkEksternal
	util.LogInfo(logModel.ToLoggerObject())

	dataMapping := repository.UserRegistrationDetailMapping{
		User: repository.UserModel{
			FirstName: sql.NullString{String: inputStruct.Firstname},
		},
		UserRegistrationDetail: repository.UserRegistrationDetailModel{
			UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
			UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
			SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
			UserID:     sql.NullString{String: inputStruct.UserID},
			Password:   sql.NullString{String: inputStruct.Password},
			Email:      sql.NullString{String: inputStruct.Email},
		},
		PKCEClientMapping: repository.PKCEClientMappingModel{
			CompanyName:  sql.NullString{String: inputStruct.CompanyName},
			BranchName:   sql.NullString{String: inputStruct.BranchName},
			ClientID:     sql.NullString{String: inputStruct.ClientID},
			ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
		},
	}
	responseAddAuth, err := input.AddUserToAuthenticationServer(userRequestAuth, dataMapping, contextModel, true, linkEksternal)

	if err.Error != nil {
		return
	}

	// add new nexTrac user
	inputStruct.AuthUserID = responseAddAuth.Nexsoft.Payload.Data.Content.UserID
	inputStruct.ClientID = responseAddAuth.Nexsoft.Payload.Data.Content.ClientID
	idUserRegistrationDetail, arrDataAudit, _, err := input.doInsertUserWithStatusPending(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	*dataAudit = append(*dataAudit, arrDataAudit...)

	// add new user registration
	err = input.CheckUserRegistrationAdmin(tx, contextModel, inputStruct, timeNow, dataAudit)
	if err.Error != nil {
		return
	}

	output.UserRegistrationDetailId = idUserRegistrationDetail
	output.ClientId = inputStruct.ClientID
	output.AuthUserId = responseAddAuth.Nexsoft.Payload.Data.Content.UserID
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) createNamedUserForNexmileOrNextradeFromAuth(tx *sql.Tx, inputStruct in.RegisterNamedUserRequest, userAuth authentication_response.UserContent, clientMappingOnDB repository.ClientMappingModel, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (output out.RegisterOrRenewLicenseUserResponse, isUpdate bool, userRegDetailOnDbByUserAuth repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "createNamedUserForNexmileOrNextradeFromAuth"
	)

	if inputStruct.ClientType == constanta.Nexstar {
		err = errorModel.GenerateUserNexstarNotFoundInAuthentication(fileName, funcName)
		return
	}

	// validation user di nextrac by auth_user_id
	userRegDetailModel := repository.UserRegistrationDetailModel{
		UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
		AuthUserID: sql.NullInt64{Int64: userAuth.UserID},
	}

	userRegDetailOnDbByUserAuth, err = dao.UserRegistrationDetailDAO.GetUserForCheckRegistrationNamedUserOrRenewNamedUser(serverconfig.ServerAttribute.DBConnection, userRegDetailModel)
	if err.Error != nil {
		return
	}

	if userRegDetailOnDbByUserAuth.Status.String == constanta.StatusActive {
		err = errorModel.GenerateUserStatusActive(input.FileName, funcName)
		return
	}

	// if true, user must be update
	if userRegDetailOnDbByUserAuth.ID.Int64 > 0 {
		isUpdate = true
		return
	}

	// resend email verification auth
	var isSendEmailFromAuth bool
	if util.IsStringEmpty(userAuth.Email) {
		if util.IsStringEmpty(inputStruct.Email) {
			err = errorModel.GenerateEmailEmptyAuthNexstarForResendVerification(fileName, funcName)
			return
		}

		authUpdateEmailReq := in.UserRequest{
			FirstName:   userAuth.FirstName,
			Username:    userAuth.Username,
			Email:       inputStruct.Email,
			Locale:      userAuth.Locale,
			AuthUserID:  userAuth.UserID,
			LastName:    userAuth.LastName,
			CountryCode: inputStruct.CountryCode,
		}

		linkEksternal := config.ApplicationConfiguration.GetNexstar().Host + config.ApplicationConfiguration.GetNexstar().PathRedirect.ActivationUser
		dataMapping := repository.UserRegistrationDetailMapping{
			User: repository.UserModel{
				FirstName: sql.NullString{String: inputStruct.Firstname},
			},
			UserRegistrationDetail: repository.UserRegistrationDetailModel{
				UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
				UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
				SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
				UserID:     sql.NullString{String: inputStruct.UserID},
				Password:   sql.NullString{String: inputStruct.Password},
				Email:      sql.NullString{String: inputStruct.Email},
			},
			PKCEClientMapping: repository.PKCEClientMappingModel{
				CompanyName:  sql.NullString{String: inputStruct.CompanyName},
				BranchName:   sql.NullString{String: inputStruct.BranchName},
				ClientID:     sql.NullString{String: inputStruct.ClientID},
				ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
			},
		}

		_, err = input.updateUserToAuthenticationServer(authUpdateEmailReq, dataMapping, linkEksternal, contextModel)
		if err.Error != nil {
			return
		}

		isSendEmailFromAuth = true
	}

	// add new nexTrac user
	inputStruct.Username = userAuth.Username
	inputStruct.AuthUserID = userAuth.UserID
	inputStruct.ClientID = userAuth.ClientID

	idUserRegistrationDetail, arrDataAudit, _, err := input.doInsertUserWithStatusPending(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	// Data send Email
	dataMapping := repository.UserRegistrationDetailMapping{
		User: repository.UserModel{
			FirstName: sql.NullString{String: inputStruct.Firstname},
		},
		UserRegistrationDetail: repository.UserRegistrationDetailModel{
			UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
			UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
			SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
			UserID:     sql.NullString{String: inputStruct.UserID},
			Password:   sql.NullString{String: inputStruct.Password},
			Email:      sql.NullString{String: inputStruct.Email},
			NoTelp:     sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.NoTelp},
			AuthUserID: sql.NullInt64{Int64: userAuth.UserID},
		},
		PKCEClientMapping: repository.PKCEClientMappingModel{
			CompanyName:  sql.NullString{String: inputStruct.CompanyName},
			BranchName:   sql.NullString{String: inputStruct.BranchName},
			ClientID:     sql.NullString{String: inputStruct.ClientID},
			ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
		},
	}

	var linkEksternal string
	if inputStruct.ClientType == constanta.Nexmile {
		linkEksternal = config.ApplicationConfiguration.GetNexmile().Host + config.ApplicationConfiguration.GetNexmile().PathRedirect.ActivationUser
	} else if inputStruct.ClientType == constanta.Nexstar {
		linkEksternal = config.ApplicationConfiguration.GetNexstar().Host + config.ApplicationConfiguration.GetNexstar().PathRedirect.ActivationUser
	} else if inputStruct.ClientType == constanta.Nextrade {
		linkEksternal = config.ApplicationConfiguration.GetNextrade().Host + config.ApplicationConfiguration.GetNextrade().PathRedirect.ActivationUser
	}

	// Step 4 - Generate OTP Nextrac untuk perpanjang named user
	var otpModel repository.UserVerificationModel
	if inputStruct.ClientType == constanta.Nexmile || inputStruct.ClientType == constanta.Nextrade {
		otpModel = input.setModelForRenewNexMileLicense(repository.UserRegistrationDetailModel{
			Email: sql.NullString{String: inputStruct.Email},
			ID:    sql.NullInt64{Int64: idUserRegistrationDetail},
		}, contextModel, timeNow)
	}

	// insert otp ke tabel user_verification
	if isSendEmailFromAuth == false {
		_, err = dao.UserVerificationDAO.InsertOTPUser(tx, otpModel)
		if err.Error != nil {
			err = input.CheckOTPDuplicateError(err)
			return
		}

		err = input.SendMessageToEmail(dataMapping, constanta.SubjectActivationEmail, " ", otpModel, linkEksternal, input.linkQueryEmail)
		if err.Error != nil {
			return
		}
	}

	*dataAudit = append(*dataAudit, arrDataAudit...)

	// add new user registration
	err = input.CheckUserRegistrationAdmin(tx, contextModel, inputStruct, timeNow, dataAudit)
	if err.Error != nil {
		return
	}

	// add resource ke auth
	resourceID := fmt.Sprint(config.ApplicationConfiguration.GetServerResourceID())
	if err = CRUDUserService.HitAuthenticationServerForAddResourceToUserAuth(userAuth, contextModel, resourceID); err.Error != nil {
		return
	}

	output.UserRegistrationDetailId = idUserRegistrationDetail
	output.ClientId = inputStruct.ClientID
	output.AuthUserId = userAuth.UserID

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) createNamedUserForNexstar(tx *sql.Tx, inputStruct in.RegisterNamedUserRequest, userAuth authentication_response.UserContent, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (output out.RegisterOrRenewLicenseUserResponse, userRegDetailOnDbByUserAuth repository.UserRegistrationDetailModel, isUpdate bool, err errorModel.ErrorModel) {
	var (
		fileName      = input.FileName
		funcName      = "createNamedUserForNexstar"
		isResendEmail bool
	)

	// validasi untuk data user di authentication ditemukan dan tipe usernya nexmile atau nextrade
	if inputStruct.ClientType == constanta.Nexmile || inputStruct.ClientType == constanta.Nextrade {
		err = errorModel.GenerateUserNexmileNextradeFoundInAuthentication(fileName, funcName)
		return
	}

	// validation user di nextrac by auth_user_id
	userRegDetailModel := repository.UserRegistrationDetailModel{
		UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
		AuthUserID: sql.NullInt64{Int64: userAuth.UserID},
	}

	userRegDetailOnDbByUserAuth, err = dao.UserRegistrationDetailDAO.GetUserForCheckRegistrationNamedUserOrRenewNamedUser(serverconfig.ServerAttribute.DBConnection, userRegDetailModel)
	if err.Error != nil {
		return
	}

	if userRegDetailOnDbByUserAuth.Status.String == constanta.StatusActive {
		err = errorModel.GenerateUserStatusActive(input.FileName, funcName)
		return
	}

	// if true, user must be update
	if userRegDetailOnDbByUserAuth.ID.Int64 > 0 {
		isUpdate = true
		return
	}

	// validasi untuk data user di authentication ditemukan dan tipe usernya nexstar tapi tidak memiliki resource chat
	if inputStruct.ClientType == constanta.Nexstar {
		if !strings.Contains(userAuth.ResourceID, constanta.GroChatResourceID) {
			err = errorModel.GenerateUserNexstarHasNoGrochatResource(fileName, funcName)
			return
		}
	}

	// resend email verification auth
	if util.IsStringEmpty(userAuth.Email) {
		if util.IsStringEmpty(inputStruct.Email) {
			err = errorModel.GenerateEmailEmptyAuthNexstarForResendVerification(fileName, funcName)
			return
		}

		authUpdateEmailReq := in.UserRequest{
			FirstName:   userAuth.FirstName,
			Username:    userAuth.Username,
			Email:       inputStruct.Email,
			Locale:      userAuth.Locale,
			AuthUserID:  userAuth.UserID,
			LastName:    userAuth.LastName,
			CountryCode: inputStruct.CountryCode,
		}

		linkEksternal := config.ApplicationConfiguration.GetNexstar().Host + config.ApplicationConfiguration.GetNexstar().PathRedirect.ActivationUser
		dataMapping := repository.UserRegistrationDetailMapping{
			User: repository.UserModel{
				FirstName: sql.NullString{String: inputStruct.Firstname},
			},
			UserRegistrationDetail: repository.UserRegistrationDetailModel{
				UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
				UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
				SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
				UserID:     sql.NullString{String: inputStruct.UserID},
				Password:   sql.NullString{String: inputStruct.Password},
				Email:      sql.NullString{String: inputStruct.Email},
			},
			PKCEClientMapping: repository.PKCEClientMappingModel{
				CompanyName:  sql.NullString{String: inputStruct.CompanyName},
				BranchName:   sql.NullString{String: inputStruct.BranchName},
				ClientID:     sql.NullString{String: inputStruct.ClientID},
				ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
			},
		}

		_, err = input.updateUserToAuthenticationServer(authUpdateEmailReq, dataMapping, linkEksternal, contextModel)
		if err.Error != nil {
			return
		}

		isResendEmail = true
	}

	// data user di authentication ditemukan dengan tipe nexstar dan memiliki resource grochat maka dilakukan add resource trac
	resourceID := fmt.Sprint(config.ApplicationConfiguration.GetServerResourceID())
	if err = CRUDUserService.HitAuthenticationServerForAddResourceToUserAuth(userAuth, contextModel, resourceID); err.Error != nil {
		return
	}

	// insert data dengan flow yang lama
	inputStruct.ClientID = userAuth.ClientID
	inputStruct.AuthUserID = userAuth.UserID
	idUserRegDetail, arrDataAudit, _, err := input.doInsertUserWithStatusNonAktif(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	*dataAudit = append(*dataAudit, arrDataAudit...)

	// Penambahan atau Perubahan user_registration_admin
	err = input.CheckUserRegistrationAdmin(tx, contextModel, inputStruct, timeNow, dataAudit)
	if err.Error != nil {
		return
	}

	// Generate otp untuk nexstar
	userRegistrationDetailModel := repository.UserRegistrationDetailModel{
		ID:     sql.NullInt64{Int64: idUserRegDetail},
		Email:  sql.NullString{String: inputStruct.Email},
		NoTelp: sql.NullString{String: inputStruct.NoTelp},
	}

	otpNexstar := input.setModelForRenewNexStarLicense(userRegistrationDetailModel, contextModel, timeNow)
	_, err = dao.UserVerificationDAO.InsertOTPUser(tx, otpNexstar)
	if err.Error != nil {
		return
	}

	if !isResendEmail {
		inputStruct.ClientID = userAuth.ClientID

		linkEksternal := config.ApplicationConfiguration.GetNexstar().PathRedirect.ActivationUser
		dataMapping := repository.UserRegistrationDetailMapping{
			User: repository.UserModel{
				FirstName: sql.NullString{String: inputStruct.Firstname},
			},
			UserRegistrationDetail: repository.UserRegistrationDetailModel{
				UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
				UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
				SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
				UserID:     sql.NullString{String: inputStruct.UserID},
				Password:   sql.NullString{String: inputStruct.Password},
				Email:      sql.NullString{String: inputStruct.Email},
				AuthUserID: sql.NullInt64{Int64: userAuth.UserID},
			},
			PKCEClientMapping: repository.PKCEClientMappingModel{
				CompanyName:  sql.NullString{String: inputStruct.CompanyName},
				BranchName:   sql.NullString{String: inputStruct.BranchName},
				ClientID:     sql.NullString{String: inputStruct.ClientID},
				ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
			},
		}

		err = input.SendMessageToEmail(dataMapping, constanta.SubjectActivationEmail, " ", otpNexstar, linkEksternal, input.linkQueryEmail)
		if err.Error != nil {
			return
		}
	}

	//--- Send Gro Chat Message
	if inputStruct.ClientType == constanta.Nexstar {
		inputStruct.ClientID = userAuth.ClientID

		linkEksternal := config.ApplicationConfiguration.GetNexstar().Host + config.ApplicationConfiguration.GetNexstar().PathRedirect.ActivationUser
		dataMapping := repository.UserRegistrationDetailMapping{
			User: repository.UserModel{
				FirstName: sql.NullString{String: inputStruct.Firstname},
			},
			UserRegistrationDetail: repository.UserRegistrationDetailModel{
				UniqueID1:  sql.NullString{String: inputStruct.UniqueID1},
				UniqueID2:  sql.NullString{String: inputStruct.UniqueID2},
				SalesmanID: sql.NullString{String: inputStruct.SalesmanID},
				UserID:     sql.NullString{String: inputStruct.UserID},
				Password:   sql.NullString{String: inputStruct.Password},
				Email:      sql.NullString{String: inputStruct.Email},
				NoTelp:     sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.NoTelp},
				AuthUserID: sql.NullInt64{Int64: userAuth.UserID},
			},
			PKCEClientMapping: repository.PKCEClientMappingModel{
				CompanyName:  sql.NullString{String: inputStruct.CompanyName},
				BranchName:   sql.NullString{String: inputStruct.BranchName},
				ClientID:     sql.NullString{String: inputStruct.ClientID},
				ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
			},
		}
		err = input.sendMessageToGrochat(dataMapping, otpNexstar, contextModel, linkEksternal)
		if err.Error != nil {
			return
		}
	}

	output.UserRegistrationDetailId = idUserRegDetail
	output.ClientId = inputStruct.ClientID
	output.AuthUserId = userAuth.UserID
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) updateUserToAuthenticationServer(inputStruct in.UserRequest, dataMapping repository.UserRegistrationDetailMapping, linkEksternal string, contextModel *applicationModel.ContextModel) (updateResposne authentication_response.UpdateUserAuthenticationResponse, err errorModel.ErrorModel) {
	emailMessages := getEmailMessage(inputStruct, dataMapping)
	return resource_common_service.InternalUpdateUser(inputStruct, contextModel, emailMessages)
}

func (input registrationNamedUserService) linkQueryAuth(linkStr string, dataOnDB repository.UserRegistrationDetailMapping) (linkQuery string, err errorModel.ErrorModel) {
	var (
		queryUrlEmailLink url.Values
		urlEmailLink      *url.URL
	)

	queryUrlEmailLink, urlEmailLink, err = input.linkQuery(linkStr, dataOnDB)
	if err.Error != nil {
		return
	}

	urlEmailLink.RawQuery = queryUrlEmailLink.Encode()
	linkQuery = urlEmailLink.String()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) linkQueryEmail(linkStr string, dataOnDB repository.UserRegistrationDetailMapping, userVerifyModel repository.UserVerificationModel) (linkQuery string, err errorModel.ErrorModel) {
	var (
		queryUrlEmailLink url.Values
		urlEmailLink      *url.URL
	)

	queryUrlEmailLink, urlEmailLink, err = input.linkQuery(linkStr, dataOnDB)
	if err.Error != nil {
		return
	}

	if util.IsStringEmpty(queryUrlEmailLink.Get(constanta.UserIDQueryParam)) {
		queryUrlEmailLink.Set(constanta.UserIDQueryParam, strconv.Itoa(int(dataOnDB.UserRegistrationDetail.AuthUserID.Int64)))
	}

	if !util.IsStringEmpty(userVerifyModel.EmailCode.String) {
		queryUrlEmailLink.Set(constanta.ActivationCodeQueryParam, userVerifyModel.EmailCode.String)
	}

	queryUrlEmailLink.Set(constanta.EmailQueryParam, dataOnDB.UserRegistrationDetail.Email.String)
	urlEmailLink.RawQuery = queryUrlEmailLink.Encode()
	linkQuery = urlEmailLink.String()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) linkQuery(linkStr string, dataOnDB repository.UserRegistrationDetailMapping) (queryUrlEmailLink url.Values, urlEmailLink *url.URL, err errorModel.ErrorModel) {
	var (
		fileName = "PerubahanRegisterUserNexmileNexstar.go"
		funcName = "linkQuery"
		errS     error
	)

	urlEmailLink, errS = url.Parse(linkStr)
	if errS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errS)
		return
	}

	queryUrlEmailLink = urlEmailLink.Query()
	queryUrlEmailLink.Set(constanta.UniqueID1QueryParam, dataOnDB.UserRegistrationDetail.UniqueID1.String)
	queryUrlEmailLink.Set(constanta.UniqueID2QueryParam, dataOnDB.UserRegistrationDetail.UniqueID2.String)
	queryUrlEmailLink.Set(constanta.SalesmanIDQueryParam, dataOnDB.UserRegistrationDetail.SalesmanID.String)
	queryUrlEmailLink.Set(constanta.UserQueryParam, dataOnDB.UserRegistrationDetail.UserID.String)
	if !util.IsStringEmpty(dataOnDB.PKCEClientMapping.ClientID.String) {
		queryUrlEmailLink.Set(constanta.ClientIDQueryParam, dataOnDB.PKCEClientMapping.ClientID.String)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) CheckUserRegistrationAdmin(tx *sql.Tx, contextModel *applicationModel.ContextModel, inputStruct in.RegisterNamedUserRequest, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	clientMappingOnDB, err := dao.ClientMappingDAO.GetParentClientMappingForClientValidation(serverconfig.ServerAttribute.DBConnection, repository.ClientMappingModel{
		ClientID:  sql.NullString{String: inputStruct.ParentClientID},
		CompanyID: sql.NullString{String: inputStruct.UniqueID1},
		BranchID:  sql.NullString{String: inputStruct.UniqueID2},
	})
	if err.Error != nil {
		return
	}

	userLicenseOnDB, err := dao.UserLicenseDAO.GetCustomerForAccountRegistration(serverconfig.ServerAttribute.DBConnection, repository.UserLicenseModel{
		ClientID:     sql.NullString{String: inputStruct.ParentClientID},
		UniqueId1:    sql.NullString{String: inputStruct.UniqueID1},
		UniqueId2:    sql.NullString{String: inputStruct.UniqueID2},
		ClientTypeId: sql.NullInt64{Int64: inputStruct.ClientTypeID},
	})
	if err.Error != nil {
		return
	}

	userRegistrationOnDB, err := dao.UserRegistrationAdminDAO.GetUserRegistrationForRegistrationNexMileNexStar(serverconfig.ServerAttribute.DBConnection, repository.UserRegistrationAdminModel{
		SiteId:           clientMappingOnDB.SiteID,
		UniqueID1:        sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:        sql.NullString{String: inputStruct.UniqueID2},
		ClientTypeID:     sql.NullInt64{Int64: inputStruct.ClientTypeID},
		ParentCustomerId: userLicenseOnDB.ParentCustomerId,
	})
	if err.Error != nil {
		return
	}

	userRegistrationModel := repository.UserRegistrationAdminModel{
		ID:               sql.NullInt64{Int64: userRegistrationOnDB.ID.Int64},
		ParentCustomerId: sql.NullInt64{Int64: userLicenseOnDB.ParentCustomerId.Int64},
		CustomerId:       sql.NullInt64{Int64: userLicenseOnDB.CustomerId.Int64},
		SiteId:           sql.NullInt64{Int64: clientMappingOnDB.SiteID.Int64},
		UniqueID1:        sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:        sql.NullString{String: inputStruct.UniqueID2},
		CompanyName:      sql.NullString{String: inputStruct.CompanyName},
		BranchName:       sql.NullString{String: inputStruct.BranchName},
		UserAdmin:        sql.NullString{String: inputStruct.UserAdmin},
		PasswordAdmin:    sql.NullString{String: inputStruct.PasswordAdmin},
		ClientTypeID:     sql.NullInt64{Int64: inputStruct.ClientTypeID},
		ClientMappingID:  sql.NullInt64{Int64: clientMappingOnDB.ID.Int64},
		CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	if userRegistrationOnDB.ID.Int64 < 1 {
		idNamedUser, errs := dao.UserRegistrationAdminDAO.InsertUserRegistrationAdmin(tx, userRegistrationModel)
		if errs.Error != nil {
			return
		}

		*dataAudit = append(*dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.UserRegistrationAdminDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: idNamedUser},
			Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
		})

		fmt.Println("Insert User Registration Admin: ", idNamedUser)
	} else {
		namedUser, errs := dao.UserRegistrationAdminDAO.UpdateUserRegistrationAdmin(tx, userRegistrationModel)
		if errs.Error != nil {
			return
		}

		*dataAudit = append(*dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.UserRegistrationAdminDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: namedUser.ID.Int64},
			Action:     sql.NullInt32{Int32: constanta.ActionAuditUpdateConstanta},
		})

		fmt.Println("Update User Registration Admin: ", namedUser)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) doInsertUserWithStatusPending(tx *sql.Tx, inputStruct in.RegisterNamedUserRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (idUserRegistrationDetail int64, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	registerNamedModel := input.createModelRegistrationNamedUserLicense(inputStruct, contextModel, timeNow)

	// Check License Quota
	userLicenseOnDb, err := input.checkLicenseQuota(&registerNamedModel)
	if err.Error != nil {
		return
	}

	// Add New PKCE Client Mapping
	pkceClientMappingModel := input.createModelForInsertPKCEClientMapping(registerNamedModel, contextModel, timeNow)
	err = input.addNewPKCEClientMapping(tx, registerNamedModel, &dataAudit, pkceClientMappingModel)
	if err.Error != nil {
		return
	}

	// Insert To User Registration Detail
	idUserRegistrationDetail, err = input.addNewToUserRegistrationDetail(tx, registerNamedModel, &dataAudit, false)
	if err.Error != nil {
		return
	}

	// Add New User
	err = input.addNewUserStatusPending(tx, registerNamedModel, timeNow, &dataAudit)
	if err.Error != nil {
		return
	}

	// Update Total Activated User License
	err = input.updateTotalActivatedUserLicense(tx, contextModel, timeNow, userLicenseOnDb, registerNamedModel, &dataAudit)
	if err.Error != nil {
		return
	}

	return
}

func (input registrationNamedUserService) doInsertUserWithStatusNonAktif(tx *sql.Tx, inputStruct in.RegisterNamedUserRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (idUserRegistrationDetail int64, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	registerNamedModel := input.createModelRegistrationNamedUserLicense(inputStruct, contextModel, timeNow)

	// Check Data In Auth Server
	err = input.checkDataToAuthServer(&registerNamedModel, contextModel)
	if err.Error != nil {
		return
	}

	// Check License Quota
	userLicenseOnDb, err := input.checkLicenseQuota(&registerNamedModel)
	if err.Error != nil {
		return
	}

	// Add New PKCE Client Mapping
	pkceClientMappingModel := input.createModelForInsertPKCEClientMapping(registerNamedModel, contextModel, timeNow)
	err = input.addNewPKCEClientMapping(tx, registerNamedModel, &dataAudit, pkceClientMappingModel)
	if err.Error != nil {
		return
	}

	// Insert To User Registration Detail
	idUserRegistrationDetail, err = input.addNewToUserRegistrationDetail(tx, registerNamedModel, &dataAudit, false)
	if err.Error != nil {
		return
	}

	// Add New User
	err = input.addNewUserStatusNonAktif(tx, registerNamedModel, timeNow, &dataAudit)
	if err.Error != nil {
		return
	}

	// Update Total Activated User License
	err = input.updateTotalActivatedUserLicense(tx, contextModel, timeNow, userLicenseOnDb, registerNamedModel, &dataAudit)
	if err.Error != nil {
		return
	}

	return
}

func setRequestForAddUserAuth(inputStruct in.RegisterNamedUserRequest, usernameNexmile string) in.UserRequest {
	return in.UserRequest{
		Username:    usernameNexmile,
		Password:    inputStruct.Password,
		FirstName:   inputStruct.Firstname,
		LastName:    inputStruct.Lastname,
		Email:       inputStruct.Email,
		Phone:       inputStruct.NoTelp,
		CountryCode: inputStruct.CountryCode,
		ResourceID:  config.ApplicationConfiguration.GetServerResourceID(),
		Scope:       constanta.ScopeClient,
		Locale:      constanta.DefaultApplicationsLanguage,
	}
}

func (input registrationNamedUserService) AddUserToAuthenticationServer(inputStruct in.UserRequest, dataMapping repository.UserRegistrationDetailMapping, contextModel *applicationModel.ContextModel,
	_ bool, linkEksternal string) (registerUserResponse authentication_response.RegisterUserAuthenticationResponse, err errorModel.ErrorModel) {

	var (
		funcName   = "AddUserToAuthenticationServer"
		authServer = config.ApplicationConfiguration.GetAuthenticationServer()
	)

	registerAuthentication := authentication_request.UserAuthenticationDTO{
		Username:              inputStruct.Username,
		Password:              inputStruct.Password,
		FirstName:             inputStruct.FirstName,
		LastName:              inputStruct.LastName,
		Email:                 inputStruct.Email,
		CountryCode:           inputStruct.CountryCode,
		Phone:                 inputStruct.Phone,
		Device:                inputStruct.Device,
		Locale:                inputStruct.Locale,
		EmailMessage:          getEmailMessage(inputStruct, dataMapping),
		IPWhitelist:           inputStruct.IPWhitelist,
		EmailLinkMessage:      "",
		PhoneMessage:          constanta.PhoneMessageEmptyDefault,
		ResourceID:            config.ApplicationConfiguration.GetServerResourceID(),
		AdditionalInformation: inputStruct.AdditionalInformation,
	}
	//log
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Status = 200
	logModel.Message = "Link External 2 : " + registerAuthentication.EmailLinkMessage
	util.LogInfo(logModel.ToLoggerObject())
	//param := make(map[string]interface{})
	//param[constanta.OTPTableParam] = fmt.Sprintf(`{{.%s}}`, constanta.OTPTableParam)
	//userVerifyModel := repository.UserVerificationModel{EmailCode: sql.NullString{String: fmt.Sprintf("%v", param[constanta.OTPTableParam])}}
	fmt.Println("After : ", registerAuthentication.EmailLinkMessage)

	registerAuthentication.EmailLinkMessage, _ = input.linkQueryAuth(linkEksternal, dataMapping)

	internalToken := resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, "", config.ApplicationConfiguration.GetServerResourceID(), constanta.DefaultApplicationsLanguage)
	registerUserUrl := authServer.Host + authServer.PathRedirect.InternalUser.CrudUser

	statusCode, bodyResult, errorS := common.HitRegisterUserAuthenticationServer(internalToken, registerUserUrl, registerAuthentication, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, "AddUserToAuthenticationServer", errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &registerUserResponse)
	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(registerUserResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, registerUserResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func getEmailMessage(inputStruct in.UserRequest, dataMapping repository.UserRegistrationDetailMapping) string {
	param := make(map[string]interface{})

	switch dataMapping.PKCEClientMapping.ClientTypeID.Int64 {
	case constanta.ResourceNexmileID:
		param[constanta.ClientTypeTableParam] = constanta.Nexmile
	case constanta.ResourceNexstarID:
		param[constanta.ClientTypeTableParam] = constanta.Nexstar
	case constanta.ResourceNextradeID:
		param[constanta.ClientTypeTableParam] = constanta.Nextrade
	default:
		param[constanta.ClientTypeTableParam] = " - "
	}

	param[constanta.PurposeTableParam] = " "
	param[constanta.CompanyIDTableParam] = dataMapping.UserRegistrationDetail.UniqueID1.String
	param[constanta.CompanyNameTableParam] = dataMapping.PKCEClientMapping.CompanyName.String
	param[constanta.BranchIDTableParam] = dataMapping.UserRegistrationDetail.UniqueID2.String
	param[constanta.BranchNameTableParam] = dataMapping.PKCEClientMapping.BranchName.String
	param[constanta.SalesmanIDTableParam] = dataMapping.UserRegistrationDetail.SalesmanID.String
	param[constanta.UserTableParam] = dataMapping.UserRegistrationDetail.UserID.String
	param[constanta.PasswordTableParam] = inputStruct.Password
	param[constanta.OTPTableParam] = fmt.Sprintf(`{{.%s}}`, constanta.OTPTableParam)
	param[constanta.LinkTableParam] = fmt.Sprintf(`{{.%s}}`, constanta.LinkTableParam)
	param[constanta.EmailTableParam] = dataMapping.UserRegistrationDetail.Email.String
	param[constanta.ClientIDTableParam] = dataMapping.PKCEClientMapping.ClientID.String
	param[constanta.RegistrationIDTableParam] = fmt.Sprintf(`{{.%s}}`, constanta.RegistrationIDTableParam)
	param[constanta.NameTableParam] = inputStruct.FirstName

	return util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, "OTP_MESSAGE_AUTH3", inputStruct.Locale, param)
}

func (input registrationNamedUserService) sendMessageToGrochat(dataOnDB repository.UserRegistrationDetailMapping, userVerifyModel repository.UserVerificationModel, contextModel *applicationModel.ContextModel, linkChannel string) (err errorModel.ErrorModel) {
	var (
		otp         = userVerifyModel.PhoneCode.String
		groMsgParam in.RegisterNamedUserRequest
		groRequest  grochat_request.SendMessageGroChatRequest
	)

	dataOnDB.UserRegistrationDetail.NoTelp.String = strings.ReplaceAll(dataOnDB.UserRegistrationDetail.NoTelp.String, "-", "")
	groMsgParam = in.RegisterNamedUserRequest{
		ClientID:     dataOnDB.PKCEClientMapping.ClientID.String,
		Firstname:    dataOnDB.User.FirstName.String,
		UniqueID1:    dataOnDB.UserRegistrationDetail.UniqueID1.String,
		CompanyName:  dataOnDB.PKCEClientMapping.CompanyName.String,
		UniqueID2:    dataOnDB.UserRegistrationDetail.UniqueID2.String,
		BranchName:   dataOnDB.PKCEClientMapping.BranchName.String,
		SalesmanID:   dataOnDB.UserRegistrationDetail.SalesmanID.String,
		UserID:       dataOnDB.UserRegistrationDetail.UserID.String,
		Password:     dataOnDB.UserRegistrationDetail.Password.String,
		ClientTypeID: dataOnDB.PKCEClientMapping.ClientTypeID.Int64,
		NoTelp:       dataOnDB.UserRegistrationDetail.NoTelp.String,
		Email:        dataOnDB.UserRegistrationDetail.Email.String,
		AuthUserID:   dataOnDB.UserRegistrationDetail.AuthUserID.Int64,
	}

	linkChannel, err = input.linkQueryGroChat(linkChannel, dataOnDB, userVerifyModel)
	if err.Error != nil {
		return
	}

	groRequest = grochat_request.SendMessageGroChatRequest{
		Data: grochat_request.DetailData{
			PhoneNumber: groMsgParam.NoTelp,
			MessageContent: grochat_request.MessageContent{
				Message: common2.GroChatService.GetGroChatMessage(groMsgParam, otp, linkChannel, " Registration "),
			},
		},
	}

	//--- Get Default
	groRequest.GetDefault(&groRequest)

	return common2.GroChatService.SendGroChatMessage(groRequest, contextModel)
}

func (input registrationNamedUserService) linkQueryGroChat(linkStr string, dataOnDB repository.UserRegistrationDetailMapping, userVerifyModel repository.UserVerificationModel) (linkQuery string, err errorModel.ErrorModel) {
	var (
		queryUrlEmailLink url.Values
		urlEmailLink      *url.URL
	)

	queryUrlEmailLink, urlEmailLink, err = input.linkQuery(linkStr, dataOnDB)
	if err.Error != nil {
		return
	}

	if util.IsStringEmpty(queryUrlEmailLink.Get(constanta.UserIDQueryParam)) {
		queryUrlEmailLink.Set(constanta.UserIDQueryParam, strconv.Itoa(int(dataOnDB.UserRegistrationDetail.AuthUserID.Int64)))
	}

	queryUrlEmailLink.Set(constanta.OTPQueryParam, userVerifyModel.PhoneCode.String)
	queryUrlEmailLink.Set(constanta.PhoneQueryParam, dataOnDB.UserRegistrationDetail.NoTelp.String)
	urlEmailLink.RawQuery = queryUrlEmailLink.Encode()
	linkQuery = urlEmailLink.String()

	err = errorModel.GenerateNonErrorModel()
	return
}
