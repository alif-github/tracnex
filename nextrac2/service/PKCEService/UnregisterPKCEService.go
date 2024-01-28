package PKCEService

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
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input pkceService) UnregisterUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UnregisterUser"

	inputStruct, err := input.readBodyAndValidateRegisUnregisPKCE(request, contextModel, false, input.validateUnregister)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUnregisterUser, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse {
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_PKCE_UNREGIS_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doUnregisterUser(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.PKCERequest)
	var userModelForDelete repository.UserModel

	//------ Check pkce client mapping and get it all data
	userOnDB, err := input.checkPKCEClientMapping(tx, inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	//------ Prepare delete user nexmile
	dataAudit, userModelForDelete, err = input.prepareDeleteUserNexmile(tx, inputStruct, userOnDB, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//------ Delete user nexmile
	err = input.deleteUserNexmile(tx, userModelForDelete, timeNow)
	if err.Error != nil {
		return
	}

	//------ Internal delete user by clientID in auth
	err = resource_common_service.InternalDeleteClientByClientID(userOnDB.ClientID.String, contextModel)
	if err.Error != nil {
		return
	}

	//------ Get list token by client ID
	var listToken []string
	listToken, err = dao.ClientTokenDAO.GetListTokenByClientID(tx, userOnDB.ClientID.String)
	if err.Error != nil {
		return
	}

	//------ Delete token in redis
	go service.DeleteTokenFromRedis(listToken)

	//------ Delete list of token in DB client token
	err = dao.ClientTokenDAO.DeleteListTokenByClientID(tx, userOnDB.ClientID.String)
	if err.Error != nil {
		return
	}

	//------ Hit unregister nexcloud
	go input.hitUnregisterNexcloudToNexcloud(userOnDB, contextModel)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) checkPKCEClientMapping(tx *sql.Tx, inputStruct in.PKCERequest, contextModel *applicationModel.ContextModel) (result repository.PKCEClientMappingModel, err errorModel.ErrorModel) {
	fileName := "UnregisterPKCEService.go"
	funcName := "checkPKCEClientMapping"

	//------ Check valid company and branch ID
	var clientMappingModel []repository.ClientMappingModel

	if inputStruct.CompanyID != "" && inputStruct.BranchID != "" {
		clientMappingModel = append(clientMappingModel, repository.ClientMappingModel{
			CompanyID: 	sql.NullString{String: inputStruct.CompanyID},
			BranchID: 	sql.NullString{String: inputStruct.BranchID},
			ClientID: 	sql.NullString{String: inputStruct.ParentClientID},
		})
		clientMappingModel, err = dao.ClientMappingDAO.CheckClientMapping(tx, clientMappingModel, true)
		if err.Error != nil {
			return
		}

		if len(clientMappingModel) < 1 {
			detail := util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.PKCEUserBundle, "DETAIL_ERROR_INVALID_ID_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
			err = errorModel.GenerateInvalidUnregistrationPKCE(fileName, funcName, []string{detail})
			return
		}
	}

	var modelPkceClientMapping repository.PKCEClientMappingModel

	modelPkceClientMapping = repository.PKCEClientMappingModel {
		ClientTypeID: 	sql.NullInt64{Int64: inputStruct.ClientTypeID},
		Username: 		sql.NullString{String: inputStruct.Username},
		ParentClientID: sql.NullString{String: inputStruct.ParentClientID},
		CompanyID: 		sql.NullString{String: inputStruct.CompanyID},
		BranchID: 		sql.NullString{String: inputStruct.BranchID},
	}

	//------ Unregister with own permission
	modelPkceClientMapping.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	result, err = dao.PKCEClientMappingDAO.CheckPKCEClientMapping(tx, modelPkceClientMapping)
	if err.Error != nil {
		return
	}

	if result.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Username + " : " + inputStruct.Username)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) prepareDeleteUserNexmile(tx *sql.Tx, inputStruct in.PKCERequest, pkceClientMappingModel repository.PKCEClientMappingModel,
	contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, modelUserForDelete repository.UserModel, err errorModel.ErrorModel) {

	fileName := "UnregisterPKCEService.go"
	funcName := "deleteUserPKCEClientMapping"

	var getModel repository.UserModel

	//------ Get id user in table user by client ID
	getModel, err = input.getUserId(pkceClientMappingModel)
	if err.Error != nil {
		return
	}

	if getModel.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Username + " : " + inputStruct.Username)
		return
	}

	modelUserForDelete = repository.UserModel{
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		ID: 			sql.NullInt64{Int64: getModel.ID.Int64},
	}

	if inputStruct.UpdatedAt != getModel.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.User + ": " + inputStruct.Username)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.UserDAO.TableName, getModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) deleteUserNexmile(tx *sql.Tx, modelUserForDelete repository.UserModel, timeNow time.Time) (err errorModel.ErrorModel) {

	//Soft delete user
	err = dao.UserDAO.DeleteUser(tx, modelUserForDelete, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) getUserId(pkceClientMappingModel repository.PKCEClientMappingModel) (result repository.UserModel, err errorModel.ErrorModel) {

	modelUserForGet := repository.UserModel{ClientID: sql.NullString{String: pkceClientMappingModel.ClientID.String}}

	//Get id user in table user by client ID
	result, err = dao.UserDAO.GetIdAndFirstNameUser(serverconfig.ServerAttribute.DBConnection, modelUserForGet)

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) validateUnregister(inputStruct *in.PKCERequest) errorModel.ErrorModel {
	return inputStruct.ValidateUnregisPKCE()
}