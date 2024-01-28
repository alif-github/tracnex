package ClientRegistrationNonOnPremiseService

import (
	"database/sql"
	"math"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"sync"
	"time"
)

func (input clientRegistrationNonOnPremiseService) InsertClientRegistNonOnPremise(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName       = "InsertClientRegistrationNonOnPremise"
		inputStruct    in.ClientRegistrationNonOnPremiseRequest
		additionalInfo interface{}
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	additionalInfo, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doInsertClientRegistrationNonOnPremise, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- func additional
	})

	if err.Error != nil {
		if additionalInfo != nil {
			output.Other = additionalInfo
			return
		}
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_CLIENT_REGISTRATION_NON_ON_PREMISE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationNonOnPremiseService) doInsertClientRegistrationNonOnPremise(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var (
		fileName                   = "InsertClientRegistrationNonOnPremiseService.go"
		funcName                   = "doInsertClientRegistrationNonOnPremise"
		inputStruct                = inputStructInterface.(in.ClientRegistrationNonOnPremiseRequest)
		clientRegistModel          repository.ClientRegistNonOnPremiseModel
		resultClientAuth           repository.ClientRegistNonOnPremiseModel
		resultError, resultSuccess []repository.DetailUniqueID
		userOnDB                   repository.UserModel
		db                         = serverconfig.ServerAttribute.DBConnection
		idClientType               int64
		isParent                   bool
	)

	isServiceUpdate = true

	//--- Create Model
	clientRegistModel = input.createModelClientRegistrationNonOnPremise(inputStruct, contextModel, timeNow)

	//--- Check Client Type
	idClientType, isParent, err = dao.ClientTypeDAO.CheckClientTypeIsParentAndExist(db, &repository.ClientTypeModel{ID: sql.NullInt64{Int64: clientRegistModel.ClientTypeID.Int64}})
	if err.Error != nil {
		return
	}

	if idClientType < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientTypeID)
		return
	}

	if !isParent {
		err = errorModel.GenerateErrorClientTypeNotParent(fileName, funcName)
		return
	}

	//--- Check Client Credential On DB
	err = input.CheckClientCredentialOnDB(fileName, &clientRegistModel)
	if err.Error != nil {
		return
	}

	//--- Hit Authentication Server For Get Client Credential
	resultClientAuth, err = input.getClientForValidateToAuthenticationServer(clientRegistModel, contextModel)
	if err.Error != nil {
		return
	}

	//--- Comparing Client Credential
	err = input.compareClientCredential(&clientRegistModel, resultClientAuth)
	if err.Error != nil {
		return
	}

	//--- Check to Customer Installation By Thread
	input.checkCustomerInstallationByThread(clientRegistModel, &resultError, &resultSuccess, contextModel, timeNow)

	//--- If Result Error Exist Then Return
	if len(resultError) > 0 {
		var clientResponseOut []out.ClientRegisterNonOnPremiseResponse
		err = errorModel.GenerateMultipleErrorAcquired(fileName, funcName)

		for _, itemResultError := range resultError {
			clientResponseOut = append(clientResponseOut, out.ClientRegisterNonOnPremiseResponse{
				UniqueID1:    itemResultError.UniqueID1.String,
				UniqueID2:    itemResultError.UniqueID2.String,
				IsError:      itemResultError.IsError.Bool,
				ErrorMessage: itemResultError.ErrorMessage.String,
			})
		}

		output = clientResponseOut
		return
	}

	clientRegistModel.DetailUnique = resultSuccess
	for _, itemData := range clientRegistModel.DetailUnique {
		err = input.insertToClientMappingUpdateCustomerInstallation(tx, itemData, clientRegistModel, &dataAudit, contextModel, timeNow)
		if err.Error != nil {
			return
		}
	}

	//--- Check User By Client ID
	userModel := repository.UserModel{ClientID: sql.NullString{String: clientRegistModel.ClientID.String}}
	userOnDB, err = dao.UserDAO.GetUserByClientID(db, userModel)
	if err.Error != nil {
		return
	}

	//--- Add User
	if userOnDB.ID.Int64 < 1 {
		err = input.addNewUser(tx, clientRegistModel, timeNow, &dataAudit)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationNonOnPremiseService) checkCustomerInstallationByThread(clientData repository.ClientRegistNonOnPremiseModel, resultError *[]repository.DetailUniqueID, resultSuccess *[]repository.DetailUniqueID,
	contextModel *applicationModel.ContextModel, timeNow time.Time) {

	var (
		wg        sync.WaitGroup
		totalPage = int(math.Ceil(float64(len(clientData.DetailUnique)) / float64(constanta.TotalDataProductLicensePerChannel)))
		resultSc  = make(chan []repository.DetailUniqueID, len(clientData.DetailUnique))
		resultEr  = make(chan []repository.DetailUniqueID, len(clientData.DetailUnique))
	)

	for i := 1; i <= totalPage; i++ {
		var (
			newData       []repository.DetailUniqueID
			offset, until int
		)

		offset = dao.CountOffset(i, constanta.TotalDataProductLicensePerChannel)
		until = offset + constanta.TotalDataProductLicensePerChannel

		if i == totalPage {
			newData = append(newData, clientData.DetailUnique[offset:]...)
		} else {
			newData = append(newData, clientData.DetailUnique[offset:until]...)
		}

		wg.Add(1)
		go input.threadProc(resultSc, resultEr, newData, contextModel, &wg, timeNow)
	}

	for j := 1; j <= totalPage; j++ {
		tempResult := <-resultSc
		*resultSuccess = append(*resultSuccess, tempResult...)
	}

	for j := 1; j <= totalPage; j++ {
		tempResult := <-resultEr
		*resultError = append(*resultError, tempResult...)
	}

	wg.Wait()
	close(resultSc)
	close(resultEr)
}

func (input clientRegistrationNonOnPremiseService) threadProc(resultSc chan []repository.DetailUniqueID, resultEr chan []repository.DetailUniqueID, jobData []repository.DetailUniqueID,
	contextModel *applicationModel.ContextModel, wg *sync.WaitGroup, timeNow time.Time) {

	input.getCustomerInstall(resultSc, resultEr, jobData, contextModel, wg, timeNow)
}

func (input clientRegistrationNonOnPremiseService) getCustomerInstall(resultSc chan []repository.DetailUniqueID, resultEr chan []repository.DetailUniqueID, jobData []repository.DetailUniqueID,
	contextModel *applicationModel.ContextModel, wg *sync.WaitGroup, timeNow time.Time) {

	var (
		fileName  = "InsertClientRegistrationNonOnPremiseService.go"
		funcName  = "getCustomerInstall"
		resultScc []repository.DetailUniqueID
		resultErr []repository.DetailUniqueID
	)

	defer wg.Done()

	for _, itemJobData := range jobData {
		var (
			err                  errorModel.ErrorModel
			resultTemp           repository.CustomerInstallationDetailConfig
			idInstallationStr    string
			idInstallationIntCol []int
			db                   = serverconfig.ServerAttribute.DBConnection
		)

		inputDataModel := repository.CustomerInstallationDetail{
			UniqueID1: sql.NullString{String: itemJobData.UniqueID1.String},
			UniqueID2: sql.NullString{String: itemJobData.UniqueID2.String},
		}

		_, resultTemp, idInstallationStr, err = dao.CustomerInstallationDAO.GetCustomerInstallation(db, inputDataModel)
		if err.Error != nil {
			input.writeErrorToModel(&resultErr, contextModel, err, inputDataModel)
		}

		idInstallationIntCol, err = service.RefactorArrayAggInt(idInstallationStr)
		if err.Error != nil {
			input.writeErrorToModel(&resultErr, contextModel, err, inputDataModel)
		}

		if len(idInstallationIntCol) < 1 {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Installation)
			input.writeErrorToModel(&resultErr, contextModel, err, inputDataModel)
		}

		//if isExistClientMapping {
		//	err = errorModel.GenerateDataUsedError(fileName, funcName, fmt.Sprintf(`%s - %s`, constanta.UniqueID1, constanta.UniqueID2))
		//	input.writeErrorToModel(&resultErr, contextModel, err, inputDataModel)
		//}

		input.writeSuccessToModel(&resultScc, resultTemp, timeNow, idInstallationIntCol)
		err = errorModel.GenerateNonErrorModel()
	}

	resultSc <- resultScc
	resultEr <- resultErr
}

func (input clientRegistrationNonOnPremiseService) writeErrorToModel(resultError *[]repository.DetailUniqueID, contextModel *applicationModel.ContextModel, err errorModel.ErrorModel, inputDataModel repository.CustomerInstallationDetail) {
	*resultError = append(*resultError, repository.DetailUniqueID{
		UniqueID1:    sql.NullString{String: inputDataModel.UniqueID1.String},
		UniqueID2:    sql.NullString{String: inputDataModel.UniqueID2.String},
		IsError:      sql.NullBool{Bool: true},
		ErrorMessage: sql.NullString{String: service.GetErrorMessage(err, *contextModel)},
	})
}

func (input clientRegistrationNonOnPremiseService) writeSuccessToModel(resultSuccess *[]repository.DetailUniqueID, resultTemp repository.CustomerInstallationDetailConfig, timeNow time.Time, idInstallationCol []int) {
	var installationIDCol []repository.InstallationIDColInt

	for _, valueIDInstallation := range idInstallationCol {
		installationIDCol = append(installationIDCol, repository.InstallationIDColInt{InstallationID: sql.NullInt64{Int64: int64(valueIDInstallation)}})
	}

	*resultSuccess = append(*resultSuccess, repository.DetailUniqueID{
		InstallationID:    sql.NullInt64{Int64: resultTemp.InstallationID.Int64},
		ParentCustomerID:  sql.NullInt64{Int64: resultTemp.ParentCustomerID.Int64},
		CustomerID:        sql.NullInt64{Int64: resultTemp.CustomerID.Int64},
		SiteID:            sql.NullInt64{Int64: resultTemp.SiteID.Int64},
		UniqueID1:         sql.NullString{String: resultTemp.UniqueID1.String},
		UniqueID2:         sql.NullString{String: resultTemp.UniqueID2.String},
		CreatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:     sql.NullString{String: constanta.SystemClient},
		CreatedAt:         sql.NullTime{Time: timeNow},
		UpdatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:     sql.NullString{String: constanta.SystemClient},
		UpdatedAt:         sql.NullTime{Time: timeNow},
		InstallationIDCol: installationIDCol,
	})
}

func (input clientRegistrationNonOnPremiseService) CheckClientCredentialOnDB(fileName string, clientRegistModel *repository.ClientRegistNonOnPremiseModel) (err errorModel.ErrorModel) {
	var (
		funcName                   = "checkClientCredentialOnDB"
		resultClientCredentialOnDB repository.ClientCredentialModel
	)

	clientModel := repository.ClientCredentialModel{ClientID: sql.NullString{String: clientRegistModel.ClientID.String}}
	resultClientCredentialOnDB, err = dao.ClientCredentialDAO.GetClientCredentialByClientID(serverconfig.ServerAttribute.DBConnection, clientModel)
	if err.Error != nil {
		return
	}

	if resultClientCredentialOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientID)
		return
	}

	clientRegistModel.ClientID.String = resultClientCredentialOnDB.ClientID.String
	clientRegistModel.ClientSecret.String = resultClientCredentialOnDB.ClientSecret.String
	clientRegistModel.SignatureKey.String = resultClientCredentialOnDB.SignatureKey.String

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationNonOnPremiseService) validateInsert(inputStruct *in.ClientRegistrationNonOnPremiseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertClientRegistNonOnPremise()
}

func (input clientRegistrationNonOnPremiseService) addNewUser(tx *sql.Tx, clientRegistModel repository.ClientRegistNonOnPremiseModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	var (
		idUser            int64
		roleID            int64
		idClientRoleScope int64
	)

	//------------ Insert New User
	idUser, err = dao.UserDAO.InsertUser(tx, repository.UserModel{
		ClientID:      sql.NullString{String: clientRegistModel.ClientID.String},
		AuthUserID:    sql.NullInt64{Int64: constanta.AuthUserNonPKCE},
		Locale:        sql.NullString{String: constanta.IndonesianLanguage},
		Status:        sql.NullString{String: "A"},
		AliasName:     sql.NullString{String: clientRegistModel.AliasName.String},
		FirstName:     sql.NullString{String: clientRegistModel.FirstName.String},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUser},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	//------------ Insert New Client Credential
	roleID = constanta.RoleUserND6
	idClientRoleScope, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: clientRegistModel.ClientID.String},
		RoleID:        sql.NullInt64{Int64: roleID},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientRoleScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idClientRoleScope},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationNonOnPremiseService) insertToClientMappingUpdateCustomerInstallation(tx *sql.Tx, itemData repository.DetailUniqueID, clientRegistModel repository.ClientRegistNonOnPremiseModel,
	dataAudit *[]repository.AuditSystemModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {

	var (
		clientMappingModelTemp repository.ClientMappingModel
		idClientMappingTemp    int64
		idFound                int64
		db                     = serverconfig.ServerAttribute.DBConnection
	)

	clientMappingModelTemp = repository.ClientMappingModel{
		ClientID:         sql.NullString{String: clientRegistModel.ClientID.String},
		CompanyID:        sql.NullString{String: itemData.UniqueID1.String},
		BranchID:         sql.NullString{String: itemData.UniqueID2.String},
		ClientTypeID:     sql.NullInt64{Int64: clientRegistModel.ClientTypeID.Int64},
		ClientAlias:      sql.NullString{String: clientRegistModel.AliasName.String},
		ParentCustomerID: sql.NullInt64{Int64: itemData.ParentCustomerID.Int64},
		CustomerID:       sql.NullInt64{Int64: itemData.CustomerID.Int64},
		SiteID:           sql.NullInt64{Int64: itemData.SiteID.Int64},
		CreatedBy:        sql.NullInt64{Int64: itemData.CreatedBy.Int64},
		CreatedClient:    sql.NullString{String: itemData.CreatedClient.String},
		CreatedAt:        sql.NullTime{Time: itemData.CreatedAt.Time},
		UpdatedBy:        sql.NullInt64{Int64: itemData.UpdatedBy.Int64},
		UpdatedClient:    sql.NullString{String: itemData.UpdatedClient.String},
		UpdatedAt:        sql.NullTime{Time: itemData.UpdatedAt.Time},
	}

	idFound, err = dao.ClientMappingDAO.CheckClientMappingByUniqueID12(db, clientMappingModelTemp)
	if err.Error != nil {
		return
	}

	if idFound < 1 {
		idClientMappingTemp, err = dao.ClientMappingDAO.InsertClientMapping(tx, &clientMappingModelTemp)
		if err.Error != nil {
			return
		}

		clientMappingModelTemp.ID.Int64 = idClientMappingTemp
		*dataAudit = append(*dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.ClientMappingDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: idClientMappingTemp},
			Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
		})
	} else {
		clientMappingModelTemp.ID.Int64 = idFound
	}

	for _, valueIDInstallation := range itemData.InstallationIDCol {
		*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerInstallationDAO.TableName, valueIDInstallation.InstallationID.Int64, 0)...)
	}

	ciModel := repository.CustomerInstallationModel{
		UpdatedBy:     sql.NullInt64{Int64: clientMappingModelTemp.UpdatedBy.Int64},
		UpdatedClient: sql.NullString{String: clientMappingModelTemp.UpdatedClient.String},
		UpdatedAt:     sql.NullTime{Time: clientMappingModelTemp.UpdatedAt.Time},
	}

	err = dao.CustomerInstallationDAO.UpdateCustomerInstallationByMultipleInstallationID(tx, ciModel, clientMappingModelTemp.ID.Int64, itemData)
	return
}
