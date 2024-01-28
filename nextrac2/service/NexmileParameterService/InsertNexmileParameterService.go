package NexmileParameterService

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
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input nexmileParameterService) InsertNexmileParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertNexmileParameter"
		inputStruct in.NexmileParameterRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doInsertNexmileParameter, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Additional Function
	})

	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterService) doInsertNexmileParameter(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName          = "InsertNexmileParameterService.go"
		funcName          = "doInsertNexmileParameter"
		inputStruct       = inputStructInterface.(in.NexmileParameterRequest)
		db                = serverconfig.ServerAttribute.DBConnection
		nexParamID        []int64
		parameterOnDB     []repository.ParameterValueModel
		clientMappingOnDB repository.ClientMappingModel
	)

	// ---------- Create Model For Nexmile Parameter
	nexmileParameterModel, nexmileParameterModelMap := input.createModelNexmileParameter(inputStruct, contextModel, timeNow)

	//// ---------- Check PKCE Client Mapping Get Parent And Client Type ID
	//if resultPKCEClientMapping, err = dao.PKCEClientMappingDAO.GetFieldForViewNexmileParameter(db, repository.PKCEClientMappingModel{
	//	ClientID: nexmileParameterModel.ClientID,
	//}); err.Error != nil {
	//	return
	//}
	//
	//if resultPKCEClientMapping.ClientMappingID.Int64 < 1 {
	//	err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientMappingClientID)
	//	return
	//}

	// ---------- Check Client Mapping With Client ID
	clientMappingOnDB, err = dao.ClientMappingDAO.CheckClientMappingWithClientID(db, repository.ClientMappingModel{
		ClientID:  nexmileParameterModel.ClientID,
		CompanyID: nexmileParameterModel.UniqueID1,
		BranchID:  nexmileParameterModel.UniqueID2,
	})
	if err.Error != nil {
		return
	}

	// ---------- Check User Login (Parent) Same as Parent Client ID in Client Mapping
	if contextModel.AuthAccessTokenModel.ClientID != clientMappingOnDB.ClientID.String {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	// ---------- Check Client Type
	if clientMappingOnDB.ClientTypeID.Int64 != nexmileParameterModel.ClientTypeID.Int64 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientMappingClientID)
		return
	}

	//---------- Check Nexmile Parameter Exist or Not
	parameterOnDB, err = dao.NexmileParameterDAO.CheckNexmileParameterExistOrNot(db, nexmileParameterModel)
	if err.Error != nil {
		return
	}

	//---------- Sort to Insert and Update
	for idx, itemParameterOnDB := range parameterOnDB {
		v, ok := nexmileParameterModelMap.ParameterData[itemParameterOnDB.ParameterID.String]
		if !ok {
			continue
		}

		delete(nexmileParameterModelMap.ParameterData, itemParameterOnDB.ParameterID.String)
		parameterOnDB[idx].ParameterValue.String = v.ParameterValue.String
	}

	// ---------- Insert Multi Nexmile Parameter
	if len(nexmileParameterModelMap.ParameterData) > 0 {
		if nexParamID, err = dao.NexmileParameterDAO.InsertMultiNexmileParameter(tx, nexmileParameterModelMap); err.Error != nil {
			err = input.checkDuplicateError(err)
			return
		}
	}

	// ---------- Noted to Data Audit
	for _, itemNexParamID := range nexParamID {
		dataAudit = append(dataAudit, service.GetAuditData(
			tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.NexmileParameterDAO.TableName, itemNexParamID, 0)...)
	}

	// ---------- Update Nexmile Parameter
	for _, itemUpdateNexmileParameter := range parameterOnDB {
		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.NexmileParameterDAO.TableName, itemUpdateNexmileParameter.ID.Int64, 0)...)
		itemUpdateNexmileParameter.UpdatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
		itemUpdateNexmileParameter.UpdatedClient.String = contextModel.AuthAccessTokenModel.ClientID
		itemUpdateNexmileParameter.UpdatedAt.Time = timeNow
		err = dao.NexmileParameterDAO.UpdateNexmileParameter(tx, itemUpdateNexmileParameter)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterService) validateInsert(inputStruct *in.NexmileParameterRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertNexmileParameter()
}
