package StandarManhourService

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
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input standarManhourService) DeleteStandarManhour(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteStandarManhour"
		inputStruct in.StandarManhourRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteStandarManhour)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteStandarManhour, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input standarManhourService) doDeleteStandarManhour(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName         = "doUpdateStandarManhour"
		inputStruct      = inputStructInterface.(in.StandarManhourRequest)
		inputModel       = input.convertDTOToModelDelete(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
		db               = serverconfig.ServerAttribute.DBConnection
		standarManhourDB repository.StandarManhourModel
	)

	//-- Check On DB
	standarManhourDB, err = input.StandarManhourDAO.GetStandarManhourForUpdate(db, inputModel)
	if err.Error != nil {
		return
	}

	if standarManhourDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.StandarManhour)
		return
	}

	//-- Check Limit Own
	err = input.CheckUserLimitedByOwnAccess(contextModel, standarManhourDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	//-- Check Lock
	if standarManhourDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.StandarManhour)
		return
	}

	//-- Update Standar Manhour
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.StandarManhourDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = input.StandarManhourDAO.DeleteStandarManhour(tx, inputModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourService) convertDTOToModelDelete(inputStruct in.StandarManhourRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.StandarManhourModel {
	return repository.StandarManhourModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input standarManhourService) validateDeleteStandarManhour(inputStruct *in.StandarManhourRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
