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

func (input standarManhourService) UpdateStandarManhour(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateStandarManhour"
		inputStruct in.StandarManhourRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateStandarManhour)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateStandarManhour, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input standarManhourService) doUpdateStandarManhour(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName         = "doUpdateStandarManhour"
		inputStruct      = inputStructInterface.(in.StandarManhourRequest)
		inputModel       = input.convertDTOToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
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
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.StandarManhourDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = input.StandarManhourDAO.UpdateStandarManhour(tx, inputModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourService) convertDTOToModelUpdate(inputStruct in.StandarManhourRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.StandarManhourModel {
	return repository.StandarManhourModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Case:          sql.NullString{String: inputStruct.Case},
		DepartmentID:  sql.NullInt64{Int64: inputStruct.DepartmentID},
		Manhour:       sql.NullFloat64{Float64: inputStruct.Manhour},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input standarManhourService) validateUpdateStandarManhour(inputStruct *in.StandarManhourRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
