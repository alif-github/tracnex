package RegistrationNamedUserService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input registrationNamedUserService) InsertNamedUserClientMapping(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName     = "InsertNamedUserClientMapping"
		inputStruct  in.RegisterNamedUserRequest
		registerUser interface{}
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsertUserLicenseClientMapping)
	if err.Error != nil {
		return
	}

	registerUser, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doInsertNamedUserClientMapping, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	//--- Output Content
	output.Data.Content = registerUser

	//--- Output Status
	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_USER_REGISTRATION_DETAIL_CLIENT_MAPPING", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) doInsertNamedUserClientMapping(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var (
		fileName                 = "InsertNamedUserClientMappingService.go"
		funcName                 = "doInsertNamedUserClientMapping"
		userLicenseOnDb          repository.UserLicenseModel
		inputStruct              = inputStructInterface.(in.RegisterNamedUserRequest)
		idUserRegistrationDetail int64
	)

	//--- Service Update
	isServiceUpdate = true

	//--- Create Model Regist
	registerNamedModel := input.createModelRegistrationNamedClientMappingUserLicense(inputStruct, contextModel, timeNow)

	//--- Check Data In Auth Server
	if err = input.checkDataToAuthServer(&registerNamedModel, contextModel); err.Error != nil {
		return
	}

	//--- Get User License
	if userLicenseOnDb, err = dao.UserLicenseDAO.CheckLicenseNamedUserForUserLicense(serverconfig.ServerAttribute.DBConnection, in.CheckLicenseNamedUserRequest{
		ClientId:     registerNamedModel.ParentClientID.String,
		UniqueId1:    registerNamedModel.UniqueID1.String,
		UniqueId2:    registerNamedModel.UniqueID2.String,
		ClientTypeID: registerNamedModel.ClientTypeID.Int64,
	}); err.Error != nil {
		return
	}

	if userLicenseOnDb.ID.Int64 < 1 {
		err = errorModel.GenerateNotFoundActiveLicense(fileName, funcName)
		return
	}

	registerNamedModel.InstallationID.Int64 = userLicenseOnDb.InstallationId.Int64
	registerNamedModel.SiteID.Int64 = userLicenseOnDb.SiteId.Int64
	registerNamedModel.UserLicenseID.Int64 = userLicenseOnDb.ID.Int64
	registerNamedModel.ProductValidFrom.Time = userLicenseOnDb.ProductValidFrom.Time
	registerNamedModel.ProductValidThru.Time = userLicenseOnDb.ProductValidThru.Time
	registerNamedModel.ParentCustomerID.Int64 = userLicenseOnDb.ParentCustomerId.Int64
	registerNamedModel.CustomerID.Int64 = userLicenseOnDb.CustomerId.Int64

	//--- Insert To User Registration Detail
	idUserRegistrationDetail, err = input.addNewToUserRegistrationDetail(tx, registerNamedModel, &dataAudit, true)
	if err.Error != nil {
		return
	}

	output = out.RegisterNamedUserResponse{UserRegistrationID: idUserRegistrationDetail}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) validateInsertUserLicenseClientMapping(inputStruct *in.RegisterNamedUserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateRegisterNamedUserClientMappingRequest()
}
