package RegistrationNamedUserService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input registrationNamedUserService) InsertNamedUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName     = "InsertNamedUser"
		inputStruct  in.RegisterNamedUserRequest
		registerUser interface{}
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	registerUser, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doInsertNamedUser, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Additional Function
	})

	if err.Error != nil {
		return
	}

	//--- Output Content
	output.Data.Content = registerUser

	//--- Output Status
	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_USER_REGISTRATION_DETAIL", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) checkDataToAuthServer(registerNamedModel *repository.UserRegistrationDetailModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName            = "InsertNamedUserService.go"
		funcName            = "checkDataToAuthServer"
		checkClientUserResp authentication_response.CheckClientOrUserResponse
	)

	if checkClientUserResp, err = service.CheckClientOrUserInAuth(authentication_request.CheckClientOrUser{ClientID: registerNamedModel.ClientID.String}, contextModel); err.Error != nil {
		return
	}

	if !checkClientUserResp.Nexsoft.Payload.Data.Content.IsExist {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientMappingClientID)
		return
	}

	if registerNamedModel.AuthUserID.Int64 != checkClientUserResp.Nexsoft.Payload.Data.Content.AdditionalInformation.UserID {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientMappingClientID)
		return
	}

	//--- Get Username From Auth
	registerNamedModel.Username.String = checkClientUserResp.Nexsoft.Payload.Data.Content.AdditionalInformation.Username

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) doInsertNamedUser(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var (
		userLicenseOnDb          repository.UserLicenseModel
		idUserRegistrationDetail int64
		inputStruct              = inputStructInterface.(in.RegisterNamedUserRequest)
		pkceClientMappingModel   repository.PKCEClientMappingModel
	)

	//--- Service Update
	isServiceUpdate = true

	//--- Create Model Regist
	registerNamedModel := input.createModelRegistrationNamedUserLicense(inputStruct, contextModel, timeNow)

	//--- Check Data In Auth Server
	err = input.checkDataToAuthServer(&registerNamedModel, contextModel)
	if err.Error != nil {
		return
	}

	//--- Check License Quota
	userLicenseOnDb, err = input.checkLicenseQuota(&registerNamedModel)
	if err.Error != nil {
		return
	}

	//--- Add New PKCE Client Mapping
	pkceClientMappingModel = input.createModelForInsertPKCEClientMapping(registerNamedModel, contextModel, timeNow)
	err = input.addNewPKCEClientMapping(tx, registerNamedModel, &dataAudit, pkceClientMappingModel)
	if err.Error != nil {
		return
	}

	//--- Insert To User Registration Detail
	idUserRegistrationDetail, err = input.addNewToUserRegistrationDetail(tx, registerNamedModel, &dataAudit, false)
	if err.Error != nil {
		return
	}

	//--- Add New User
	err = input.addNewUser(tx, registerNamedModel, timeNow, &dataAudit)
	if err.Error != nil {
		return
	}

	//--- Update Total Activated User License
	err = input.updateTotalActivatedUserLicense(tx, contextModel, timeNow, userLicenseOnDb, registerNamedModel, &dataAudit)
	if err.Error != nil {
		return
	}

	output = out.RegisterNamedUserResponse{UserRegistrationID: idUserRegistrationDetail}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) updateTotalActivatedUserLicense(tx *sql.Tx, contextModel *applicationModel.ContextModel, timeNow time.Time, userLicenseOnDb repository.UserLicenseModel,
	registerNamedModel repository.UserRegistrationDetailModel, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {

	*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserLicenseDAO.TableName, userLicenseOnDb.ID.Int64, 0)...)

	//----------------------- Update Total Activated User License
	if err = dao.UserLicenseDAO.UpdateTotalActivatedUserLicense(tx, repository.UserLicenseModel{
		UpdatedBy:     sql.NullInt64{Int64: registerNamedModel.UpdatedBy.Int64},
		UpdatedClient: sql.NullString{String: registerNamedModel.UpdatedClient.String},
		UpdatedAt:     sql.NullTime{Time: registerNamedModel.UpdatedAt.Time},
		ID:            sql.NullInt64{Int64: registerNamedModel.UserLicenseID.Int64},
	}); err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) addNewToUserRegistrationDetail(tx *sql.Tx, registerNamedModel repository.UserRegistrationDetailModel, dataAudit *[]repository.AuditSystemModel, isFromClientMapping bool) (idUserRegistrationDetail int64, err errorModel.ErrorModel) {

	if idUserRegistrationDetail, err = dao.UserRegistrationDetailDAO.InsertUserRegistrationDetail(tx, registerNamedModel, isFromClientMapping); err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserRegistrationDetailDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUserRegistrationDetail},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) addNewToUserRegistrationDetailWithReturnModel(tx *sql.Tx, registerNamedModel repository.UserRegistrationDetailModel, dataAudit *[]repository.AuditSystemModel, isFromClientMapping bool) (idUserRegistrationDetail int64, err errorModel.ErrorModel) {

	if idUserRegistrationDetail, err = dao.UserRegistrationDetailDAO.InsertUserRegistrationDetail(tx, registerNamedModel, isFromClientMapping); err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserRegistrationDetailDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUserRegistrationDetail},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) validateInsert(inputStruct *in.RegisterNamedUserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateRegisterNamedUserRequest()
}

func (input registrationNamedUserService) createModelForInsertPKCEClientMapping(registerNamedModel repository.UserRegistrationDetailModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (pkceClientMappingModel repository.PKCEClientMappingModel) {
	return repository.PKCEClientMappingModel{
		Username:       sql.NullString{String: registerNamedModel.Username.String},
		ClientID:       sql.NullString{String: registerNamedModel.ClientID.String},
		CompanyID:      sql.NullString{String: registerNamedModel.UniqueID1.String},
		BranchID:       sql.NullString{String: registerNamedModel.UniqueID2.String},
		ClientTypeID:   sql.NullInt64{Int64: registerNamedModel.ClientTypeID.Int64},
		ParentClientID: sql.NullString{String: registerNamedModel.ParentClientID.String},
		AuthUserID:     sql.NullInt64{Int64: registerNamedModel.AuthUserID.Int64},
		InstallationID: sql.NullInt64{Int64: registerNamedModel.InstallationID.Int64},
		CustomerID:     sql.NullInt64{Int64: registerNamedModel.CustomerID.Int64},
		SiteID:         sql.NullInt64{Int64: registerNamedModel.SiteID.Int64},
		ClientAlias:    sql.NullString{String: registerNamedModel.ClientAliases.String},
		CreatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:      sql.NullTime{Time: timeNow},
		UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:      sql.NullTime{Time: timeNow},
	}
}
