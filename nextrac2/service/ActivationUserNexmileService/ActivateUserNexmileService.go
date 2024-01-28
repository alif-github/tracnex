package ActivationUserNexmileService

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
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type tempStructCreatePKCEModel struct {
	inputStruct  in.ActivationUserNexmileRequest
	clientOnAuth authentication_response.AdditionalInformationContent
	userLicense  repository.UserLicenseModel
	userRegis    repository.UserRegistrationDetailModel
}

func (input activationUserNexmileService) ActivateUserNexmile(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "ActivateLicense"
	var inputStruct in.ActivationUserNexmileRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateActivateUserNexmile)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.activateUserNexmile, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_ACTIVATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input activationUserNexmileService) activateUserNexmile(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "activateUserNexmile"
	inputStruct := inputStructInterface.(in.ActivationUserNexmileRequest)

	var userRegistrationDetail repository.UserRegistrationDetailModel
	var userLicense repository.UserLicenseModel
	var clientOnAuth authentication_response.CheckClientOrUserResponse

	//validate client
	if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ParentClientID {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	//check license
	if userRegistrationDetail, err = dao.UserRegistrationDetailDAO.GetUserRegistrationDetailForActivation(tx, repository.UserRegistrationDetailModel{
		ID: sql.NullInt64{Int64: inputStruct.UserRegistrationDetailID},
	}); err.Error != nil {
		return
	}

	if userRegistrationDetail.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserRegistrationDetailIDName)
		return
	}

	if clientOnAuth, err = input.checkDataToAuthServer(inputStruct, contextModel); err.Error != nil {
		return
	}

	// Get User License
	if userLicense, err = dao.UserLicenseDAO.ViewDetailUserLicense(serverconfig.ServerAttribute.DBConnection, repository.UserLicenseModel{
		ID: userRegistrationDetail.UserLicenseID,
	}); err.Error != nil {
		return
	}

	if userLicense.ID.Int64 < 1 {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationUserNexmileErrorUserLicense, nil)
		return
	}

	if input.validateActiveDate(userLicense, timeNow) || (userLicense.TotalLicense.Int64-userLicense.TotalActivated.Int64) < 1 || userLicense.LicenseStatus.Int64 != constanta.ProductLicenseStatusActive {
		// change user license
		if userLicense, err = input.getAvailableUserLicense(repository.UserLicenseModel{
			ClientTypeId: userLicense.ClientTypeId,
			ClientID:     sql.NullString{String: inputStruct.ParentClientID},
			UniqueId1:    userRegistrationDetail.UniqueID1,
			UniqueId2:    userRegistrationDetail.UniqueID2,
		}); err.Error != nil {
			return
		}
	}

	return input.doActivateUserNexmile(tx, tempStructCreatePKCEModel{
		inputStruct:  inputStruct,
		clientOnAuth: clientOnAuth.Nexsoft.Payload.Data.Content.AdditionalInformation,
		userLicense:  userLicense,
		userRegis:    userRegistrationDetail,
	}, contextModel, timeNow)
}

func (input activationUserNexmileService) doActivateUserNexmile(tx *sql.Tx, param tempStructCreatePKCEModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var pkceClientMapping repository.PKCEClientMappingModel
	var userModel repository.UserModel
	var insertedID int64
	var isClientDependant bool

	// get PKCE adn User Model
	pkceClientMapping = input.createPKCEClientMappingModel(param, contextModel, timeNow)
	userModel = input.createUserModel(param, contextModel, timeNow)

	// Update User Registration Detail
	if err = dao.UserRegistrationDetailDAO.UpdateActivationUserRegistrationDetail(tx, repository.UserRegistrationDetailModel{
		ID:            param.userRegis.ID,
		UserLicenseID: param.userLicense.ID,
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		Status:        sql.NullString{String: constanta.StatusActive},
	}); err.Error != nil {
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserRegistrationDetailDAO.TableName, param.userRegis.ID.Int64, 0)...)

	// Insert PKCE Client Mapping
	if pkceClientMapping.IsClientDependant.String == constanta.FlagStatusTrue {
		isClientDependant = true
	}

	if insertedID, err = dao.PKCEClientMappingDAO.InsertPKCEClientMapping(tx, &pkceClientMapping, isClientDependant); err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.PKCEClientMappingDAO.TableName, insertedID, 0)...)

	// Insert User
	if insertedID, err = dao.UserDAO.InsertUser(tx, userModel); err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.UserDAO.TableName, insertedID, 0)...)

	insertedID, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: param.inputStruct.ClientID},
		RoleID:        sql.NullInt64{Int64: constanta.RoleUserNexMile},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	// Update User License
	if err = dao.UserLicenseDAO.UpdateTotalActivatedUserLicense(tx, repository.UserLicenseModel{
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		ID:            sql.NullInt64{Int64: param.userLicense.ID.Int64},
	}); err.Error != nil {
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserLicenseDAO.TableName, param.userLicense.ID.Int64, 0)...)

	return
}

func (input activationUserNexmileService) createUserModel(param tempStructCreatePKCEModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (result repository.UserModel) {
	return repository.UserModel{
		ClientID:       sql.NullString{String: param.clientOnAuth.ClientID},
		AuthUserID:     sql.NullInt64{Int64: param.clientOnAuth.UserID},
		IPWhitelist:    sql.NullString{String: param.clientOnAuth.IPWhitelist},
		Locale:         sql.NullString{String: constanta.DefaultApplicationsLanguage},
		SignatureKey:   sql.NullString{String: param.clientOnAuth.SignatureKey},
		AdditionalInfo: sql.NullString{String: param.clientOnAuth.ClientInformation},
		IsSystemAdmin:  sql.NullBool{Bool: false},
		Status:         sql.NullString{String: constanta.StatusActive},
		FirstName:      sql.NullString{String: param.inputStruct.FirstName},
		LastName:       sql.NullString{String: param.inputStruct.LastName},
		Username:       sql.NullString{String: param.clientOnAuth.Username},
		Email:          sql.NullString{String: param.inputStruct.Email},
		Phone:          sql.NullString{String: param.inputStruct.Phone},
		CreatedBy:      sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:      sql.NullTime{Time: timeNow},
		CreatedClient:  sql.NullString{String: constanta.SystemClient},
		UpdatedBy:      sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:      sql.NullTime{Time: timeNow},
		UpdatedClient:  sql.NullString{String: constanta.SystemClient},
		AliasName:      sql.NullString{String: param.inputStruct.AliasName},
	}
}

func (input activationUserNexmileService) createPKCEClientMappingModel(param tempStructCreatePKCEModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (result repository.PKCEClientMappingModel) {
	result = repository.PKCEClientMappingModel{
		ParentClientID:    sql.NullString{String: param.inputStruct.ParentClientID},
		ClientID:          param.userRegis.ClientID,
		ClientTypeID:      param.userLicense.ClientTypeId,
		AuthUserID:        sql.NullInt64{Int64: param.clientOnAuth.UserID},
		Username:          sql.NullString{String: param.clientOnAuth.Username},
		CustomerID:        param.userLicense.CustomerId,
		SiteID:            param.userLicense.SiteId,
		CompanyID:         param.userRegis.UniqueID1,
		BranchID:          param.userRegis.UniqueID2,
		ClientAlias:       sql.NullString{String: param.inputStruct.AliasName},
		CreatedBy:         sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:     sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:         sql.NullTime{Time: timeNow},
		UpdatedBy:         sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:     sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:         sql.NullTime{Time: timeNow},
		IsClientDependant: sql.NullString{String: constanta.StatusNonActive},
	}

	if param.userLicense.ClientTypeId.Int64 == constanta.ResourceNexmileID {
		result.IsClientDependant.String = constanta.FlagStatusTrue
	}
	return
}

func (input activationUserNexmileService) checkDataToAuthServer(inputStruct in.ActivationUserNexmileRequest, contextModel *applicationModel.ContextModel) (checkClientUserResp authentication_response.CheckClientOrUserResponse, err errorModel.ErrorModel) {
	funcName := "checkDataToAuthServer"

	checkClientUserResp, err = service.CheckClientOrUserInAuth(authentication_request.CheckClientOrUser{ClientID: inputStruct.ClientID}, contextModel)
	if err.Error != nil {
		return
	}

	if !checkClientUserResp.Nexsoft.Payload.Data.Content.IsExist {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientMappingClientID)
		return
	}

	if inputStruct.AuthUserID != checkClientUserResp.Nexsoft.Payload.Data.Content.AdditionalInformation.UserID {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientMappingClientID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input activationUserNexmileService) validateActiveDate(model repository.UserLicenseModel, timeNow time.Time) bool {
	layoutFormat := "2006-01-02"
	comparedTime, _ := time.Parse(layoutFormat, timeNow.Format(layoutFormat))

	if model.ProductValidFrom.Time.Unix() <= comparedTime.Unix() && model.ProductValidThru.Time.Unix() >= comparedTime.Unix() {
		return true
	}

	return false
}

func (input activationUserNexmileService) getAvailableUserLicense(inputStruct repository.UserLicenseModel) (result repository.UserLicenseModel, err errorModel.ErrorModel) {
	var countUserLicense int
	funcName := "getAvailableUserLicense"

	if countUserLicense, err = dao.UserLicenseDAO.GetCountUserLicenseForActivationNamedUser(serverconfig.ServerAttribute.DBConnection, inputStruct); err.Error != nil {
		return
	}

	if countUserLicense < 1 {
		//error user license tidak ditemukan
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationUserNexmileErrorUserLicense, nil)
		return
	}

	if result, err = dao.UserLicenseDAO.GetAvailableUserLicenseForActivateNamedUser(serverconfig.ServerAttribute.DBConnection, inputStruct); err.Error != nil {
		return
	}

	return
}

func (input activationUserNexmileService) validateActivateUserNexmile(inputStruct *in.ActivationUserNexmileRequest) errorModel.ErrorModel {
	return inputStruct.ValidateActivationUserNexmile()
}
