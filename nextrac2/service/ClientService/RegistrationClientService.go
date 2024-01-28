package ClientService

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
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/SocketIDService"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input clientService) RegistrationClient(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName       = "RegistrationClient"
		registerClient interface{}
		inputStruct    in.ClientRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateRegistration)
	if err.Error != nil {
		return
	}

	registerClient, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doRegistrationClient, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output = registerClient.(out.Payload)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) doRegistrationClient(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var (
		fileName              = "RegistrationClientService.go"
		funcName              = "doRegistrationClient"
		inputStruct           = inputStructInterface.(in.ClientRequest)
		detail                string
		resourceIDList        []out.ResourceList
		custInstallationOnDB  []repository.CustomerInstallationForConfig
		registerClientContent authentication_response.RegisterClientContent
	)

	output = out.Payload{}

	//---------- Check data customer installation
	custInstallationOnDB, err = input.checkDataInCustomerInstallation(inputStruct)
	if err.Error != nil {
		return
	}

	//---------- Check in client mapping, is client registered before?
	newInputStructInterface, clientCredential, err := input.DoCheckClientMappingForRegister(tx, inputStructInterface)
	if err.Error != nil {
		return
	}

	//---------- All data has been registered before
	if !util2.IsStringEmpty(clientCredential.ClientID.String) && len(newInputStructInterface.(in.ClientRequest).CompanyData) == 0 {
		var tempDataAudit []repository.AuditSystemModel
		output, dataAudit, err = input.dataRegisteredProcess(tx, newInputStructInterface, clientCredential, contextModel, timeNow)
		if err.Error != nil {
			return
		}

		tempDataAudit, err = input.CheckDataMustUpdateInCustomerInstallation(tx, custInstallationOnDB, clientCredential.ClientID.String, contextModel, timeNow)
		if err.Error != nil {
			return
		}

		isServiceUpdate = true
		dataAudit = append(dataAudit, tempDataAudit...)
		err = errorModel.GenerateNonErrorModel()
		return
		//---------- Any data must reg to new branch
	} else if !util2.IsStringEmpty(clientCredential.ClientID.String) && len(newInputStructInterface.(in.ClientRequest).CompanyData) > 0 {
		var tempDataAudit []repository.AuditSystemModel
		output, dataAudit, err = input.doRegisterNewBranch(tx, newInputStructInterface, clientCredential, contextModel, timeNow)
		if err.Error != nil {
			return
		}

		tempDataAudit, err = input.CheckDataMustUpdateInCustomerInstallation(tx, custInstallationOnDB, clientCredential.ClientID.String, contextModel, timeNow)
		if err.Error != nil {
			return
		}

		isServiceUpdate = true
		dataAudit = append(dataAudit, tempDataAudit...)
		err = errorModel.GenerateNonErrorModel()
		return
	}

	//---------- Check feedback data is empty?
	if len(newInputStructInterface.(in.ClientRequest).CompanyData) == 0 {
		detail = util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientBundle, "DETAIL_ERROR_DATA_HAS_BEEN_REGISTERED", contextModel.AuthAccessTokenModel.Locale, nil)
		err = errorModel.GenerateInvalidRegistrationClient(fileName, funcName, []string{detail})
		return
	}

	//---------- Check is costumer listing exist?
	newInputStructInterface, err = input.doCheckCustomerExistPrepareByErrorForRegClient(tx, newInputStructInterface, contextModel, timeNow, true)
	if err.Error != nil {
		return
	}

	//---------- Hit create new client to authentication server
	registerClientContent, err = input.addClientToAuthenticationServer(newInputStructInterface.(in.ClientRequest), contextModel)
	if err.Error != nil {
		return
	}

	//---------- Add Resource nextrac to auth
	_, err = input.addResourceClient(registerClientContent, contextModel)
	if err.Error != nil {
		return
	} else {
		resourceIDList = append(resourceIDList, out.ResourceList{
			ResourceID: config.ApplicationConfiguration.GetServerResourceID(),
			Status:     "OK",
		})
	}

	//---------- Insert new registered client to user
	idUser, idClientRoleScope, err := input.doInsertToUser(tx, newInputStructInterface.(in.ClientRequest), registerClientContent, timeNow)
	if err.Error != nil {
		return
	}

	//---------- Insert new registered client to client credential
	idClientCredential, err := input.doInsertToClientCredential(tx, registerClientContent, timeNow)
	if err.Error != nil {
		return
	}

	//---------- Insert new registered client to client mapping
	idClientMapping, err := input.doInsertToClientMapping(tx, newInputStructInterface, registerClientContent, timeNow)
	if err.Error != nil {
		return
	}

	//---------- Audit for all insert table process
	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUser},
	}, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientRoleScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idClientRoleScope},
	}, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientCredentialDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idClientCredential},
	})

	for _, idClientMappingElm := range idClientMapping {
		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.ClientMappingDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: idClientMappingElm},
		})
	}

	tempDataAudit, err := input.CheckDataMustUpdateInCustomerInstallation(tx, custInstallationOnDB, registerClientContent.ClientID, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, tempDataAudit...)

	//---------- Hit add resource client to external (nexcloud and nexdrive)
	successResource, failedResource, errS := input.addResourceExternal(newInputStructInterface, registerClientContent, contextModel)
	if errS.Error != nil {
		output, err = input.clientErrorHandle(newInputStructInterface.(in.ClientRequest), contextModel, resourceIDList,
			successResource, failedResource, registerClientContent, timeNow, errS)
		return
	} else {
		resourceIDList = append(resourceIDList, out.ResourceList{
			ResourceID: constanta.NexCloudResourceID,
			Status:     "OK",
		})
	}

	output, err = input.clientSuccessHandle(registerClientContent, contextModel, newInputStructInterface.(in.ClientRequest), resourceIDList, "SUCCESS_REGIST_MESSAGE", timeNow)
	err = errorModel.GenerateNonErrorModel()
	//todo check----------------------------------------
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Message = "Sukses Read Body"
	logModel.Status = 200
	util2.LogInfo(logModel.ToLoggerObject())
	//todo----------------------------------------------
	return
}

func (input clientService) doInsertToClientMapping(tx *sql.Tx, inputStructInterface interface{}, registerClientContent authentication_response.RegisterClientContent, timeNow time.Time) (id []int64, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.ClientRequest)
	var modelClientMapping []repository.ClientMappingModel

	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			modelClientMapping = append(modelClientMapping, repository.ClientMappingModel{
				ClientTypeID:  sql.NullInt64{Int64: inputStruct.ClientTypeID},
				ClientID:      sql.NullString{String: registerClientContent.ClientID},
				CompanyID:     sql.NullString{String: companyDataElm.CompanyID},
				BranchID:      sql.NullString{String: branchDataElm.BranchID},
				SocketID:      sql.NullString{String: inputStruct.SocketID},
				ClientAlias:   sql.NullString{String: branchDataElm.ClientAlias},
				CreatedAt:     sql.NullTime{Time: timeNow},
				UpdatedAt:     sql.NullTime{Time: timeNow},
				CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
				UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
				UpdatedClient: sql.NullString{String: constanta.SystemClient},
				CreatedClient: sql.NullString{String: constanta.SystemClient},
			})
		}
	}

	id, err = dao.ClientMappingDAO.InsertMultipleClientMapping(tx, modelClientMapping)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) doInsertToUser(tx *sql.Tx, inputStruct in.ClientRequest, registerClientContent authentication_response.RegisterClientContent, timeNow time.Time) (idUser int64, idClientRoleScope int64, err errorModel.ErrorModel) {
	var roleID int64
	idUser, err = dao.UserDAO.InsertUser(tx, repository.UserModel{
		ClientID:      sql.NullString{String: registerClientContent.ClientID},
		AuthUserID:    sql.NullInt64{Int64: constanta.AuthUserNonPKCE},
		Locale:        sql.NullString{String: constanta.IndonesianLanguage},
		Status:        sql.NullString{String: constanta.StatusActive},
		FirstName:     sql.NullString{String: inputStruct.ClientName},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		return
	}

	if inputStruct.ClientTypeID == constanta.ResourceND6ID {
		roleID = constanta.RoleUserND6
	}

	idClientRoleScope, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: registerClientContent.ClientID},
		RoleID:        sql.NullInt64{Int64: roleID},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) doInsertToClientCredential(tx *sql.Tx, registerClientContent authentication_response.RegisterClientContent, timeNow time.Time) (idClientCredential int64, err errorModel.ErrorModel) {
	idClientCredential, err = dao.ClientCredentialDAO.InsertClientCredential(tx, &repository.ClientCredentialModel{
		ClientID:      sql.NullString{String: registerClientContent.ClientID},
		ClientSecret:  sql.NullString{String: registerClientContent.ClientSecret},
		SignatureKey:  sql.NullString{String: registerClientContent.SignatureKey},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) addResourceInfo(inputStruct authentication_response.RegisterClientContent, resourceDataList []out.ResourceList) (result authentication_response.RegisterClientContent) {
	inputStruct.ResourceList = resourceDataList
	result = inputStruct
	return
}

func (input clientService) doCheckCustomerExistPrepareByError(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time, isCheckExp bool) (output interface{}, err errorModel.ErrorModel) {
	fileName := "RegistrationClientService.go"
	funcName := "doCheckCustomerExistPrepareByError"

	detail := util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientBundle, "DETAIL_INVALID_BRANCH_COMPANY_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
	errors := errorModel.GenerateInvalidRegistrationClient(fileName, funcName, []string{detail})

	output, err = input.DoCheckCustomerExist(tx, inputStructInterface, errors, contextModel, timeNow, isCheckExp)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) validateRegistration(inputStruct *in.ClientRequest) errorModel.ErrorModel {
	return inputStruct.ValidateRegisClient()
}

func (input clientService) doUpdateSocketID(tx *sql.Tx, clientCredential repository.ClientCredentialModel, inputStruct in.ClientRequest, contextModel *applicationModel.ContextModel,
	timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {

	userOnDB, err := dao.UserDAO.GetAuthUserByClientID(serverconfig.ServerAttribute.DBConnection, repository.UserModel{ClientID: sql.NullString{String: clientCredential.ClientID.String}})
	if err.Error != nil {
		return
	}

	contextModel.AuthAccessTokenModel.ClientID = clientCredential.ClientID.String
	contextModel.LoggerModel.ClientID = clientCredential.ClientID.String
	contextModel.AuthAccessTokenModel.ResourceUserID = userOnDB.AuthUserID.Int64
	clientMappingBody := in.ClientMappingForUIRequest{
		ClientTypeId: constanta.ResourceND6ID,
		ClientId:     clientCredential.ClientID.String,
		SocketID:     inputStruct.SocketID,
	}

	_, dataAudit, err = SocketIDService.DoUpdateSocketID(tx, clientMappingBody, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) dataRegisteredProcess(tx *sql.Tx, newInputStructInterface interface{}, clientCredential repository.ClientCredentialModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		logDB              repository.ClientRegistrationLogModel
		resourceIDList     []out.ResourceList
		successMessage     string
		code               string
		statusResourceTrac = "FAIL"
		statusResourceApi  = "FAIL"
	)

	//---------- Get data user
	userOnDB, err := dao.UserDAO.GetAuthUserByClientID(serverconfig.ServerAttribute.DBConnection, repository.UserModel{ClientID: sql.NullString{String: clientCredential.ClientID.String}})
	if err.Error != nil {
		return
	}

	contextModel.AuthAccessTokenModel.ClientID = clientCredential.ClientID.String
	contextModel.AuthAccessTokenModel.ResourceUserID = userOnDB.ID.Int64

	//---------- Get log
	logDB, err = dao.ClientRegistrationLogDAO.GetDataStatusResource(serverconfig.ServerAttribute.DBConnection, repository.ClientRegistrationLogModel{ClientID: sql.NullString{String: clientCredential.ClientID.String}})
	if err.Error != nil {
		return
	}

	if logDB.SuccessStatusAuth.Bool {
		statusResourceTrac = "OK"
	}
	if logDB.SuccessStatusNexcloud.Bool {
		statusResourceApi = "OK"
	}

	//---------- Set status
	resourceIDList = append(resourceIDList, out.ResourceList{
		ResourceID: config.ApplicationConfiguration.GetServerResourceID(),
		Status:     statusResourceTrac,
	}, out.ResourceList{
		ResourceID: constanta.NexCloudResourceID,
		Status:     statusResourceApi,
	})

	//---------- Set output
	successMessage, code, output = input.customResponseAndAddResourceInfo(authentication_response.RegisterClientContent{
		ClientID:     clientCredential.ClientID.String,
		ClientSecret: clientCredential.ClientSecret.String,
		SignatureKey: clientCredential.SignatureKey.String,
	}, contextModel, resourceIDList, "SUCCESS_REGIST_MESSAGE")

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ClientRegistrationLogDAO.TableName, logDB.ID.Int64, 0)...)

	err = dao.ClientRegistrationLogDAO.UpdateLogAfterReRegistration(tx, repository.ClientRegistrationLogModel{
		ID:               sql.NullInt64{Int64: logDB.ID.Int64},
		AttributeRequest: sql.NullString{String: newInputStructInterface.(in.ClientRequest).BodyRequest.Body},
		MessageAuth:      sql.NullString{String: successMessage},
		Code:             sql.NullString{String: code},
		RequestTimeStamp: sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		RequestCount:     sql.NullInt64{Int64: logDB.RequestCount.Int64 + 1},
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) doRegisterNewBranch(tx *sql.Tx, newInputStructInterface interface{}, clientCredential repository.ClientCredentialModel,
	contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {

	var (
		fileName                 = "RegistrationClientService.go"
		funcName                 = "doRegisterNewBranch"
		userOnDB                 repository.UserModel
		resultClientMappingModel repository.ClientMappingModel
		modelClientMapping       []repository.ClientMappingModel
		idClientMapping          []int64
	)

	//---------- Get data user by client ID in DB user
	userOnDB, err = dao.UserDAO.GetAuthUserByClientID(serverconfig.ServerAttribute.DBConnection, repository.UserModel{ClientID: sql.NullString{String: clientCredential.ClientID.String}})
	if err.Error != nil {
		return
	}

	//---------- Get data header
	resultClientMappingModel, err = dao.ClientMappingDAO.GetDataHeaderClientMapping(tx, repository.ClientMappingModel{ClientID: sql.NullString{String: clientCredential.ClientID.String}})
	if err.Error != nil {
		return
	}

	if util2.IsStringEmpty(resultClientMappingModel.ClientID.String) {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientID)
		return
	}

	//---------- Mapping data
	for _, companyDataElm := range newInputStructInterface.(in.ClientRequest).CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			modelClientMapping = append(modelClientMapping, repository.ClientMappingModel{
				ClientID:      sql.NullString{String: clientCredential.ClientID.String},
				ClientTypeID:  sql.NullInt64{Int64: newInputStructInterface.(in.ClientRequest).ClientTypeID},
				CompanyID:     sql.NullString{String: companyDataElm.CompanyID},
				BranchID:      sql.NullString{String: branchDataElm.BranchID},
				ClientAlias:   sql.NullString{String: branchDataElm.BranchName},
				SocketID:      sql.NullString{String: resultClientMappingModel.SocketID.String},
				CreatedBy:     sql.NullInt64{Int64: userOnDB.ID.Int64},
				CreatedClient: sql.NullString{String: clientCredential.ClientID.String},
				CreatedAt:     sql.NullTime{Time: timeNow},
				UpdatedBy:     sql.NullInt64{Int64: userOnDB.ID.Int64},
				UpdatedClient: sql.NullString{String: clientCredential.ClientID.String},
				UpdatedAt:     sql.NullTime{Time: timeNow},
			})
		}
	}

	idClientMapping, err = dao.ClientMappingDAO.InsertMultipleClientMapping(tx, modelClientMapping)
	if err.Error != nil {
		return
	}

	output, dataAudit, err = input.dataRegisteredProcess(tx, newInputStructInterface, clientCredential, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	for _, idClientMappingElm := range idClientMapping {
		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.ClientMappingDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: idClientMappingElm},
			Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) doCheckCustomerExistPrepareByErrorForRegClient(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time, isCheckExp bool) (output interface{}, err errorModel.ErrorModel) {
	fileName := "RegistrationClientService.go"
	funcName := "doCheckCustomerExistPrepareByErrorForRegClient"

	detail := util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientBundle, "DETAIL_INVALID_BRANCH_COMPANY_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
	errors := errorModel.GenerateInvalidRegistrationClient(fileName, funcName, []string{detail})

	output, err = input.DoCheckCustomerExistForRegistrationClientID(tx, inputStructInterface, errors, timeNow, isCheckExp)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
