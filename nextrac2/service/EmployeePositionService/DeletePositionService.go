package EmployeePositionService

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

func (input employeePositionService) DeleteEmployeePosition(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteEmployeePosition"
		inputStruct in.EmployeePosition
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteEmployeePosition)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteEmployeePosition, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input employeePositionService) doDeleteEmployeePosition(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var inputStruct = inputStructInterface.(in.EmployeePosition)

	//--- Created Input Model
	inputModel := repository.EmployeePositionModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	//--- Validate Employee Position Check DB
	err = input.validateEmployeePositionOnDB(inputStruct, inputModel, contextModel)
	if err.Error != nil {
		return
	}

	//--- Delete Employee Position
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.EmployeePositionDAO.TableName, inputModel.ID.Int64, 0)...)
	err = dao.EmployeePositionDAO.DeleteEmployeePosition(tx, inputModel)
	return
}

func (input employeePositionService) validateEmployeePositionOnDB(inputStruct in.EmployeePosition, inputModel repository.EmployeePositionModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName     = "DeletePositionService.go"
		funcName     = "validateEmployeePositionOnDB"
		db           = serverconfig.ServerAttribute.DBConnection
		positionOnDB repository.EmployeePositionModel
	)

	positionOnDB, err = dao.EmployeePositionDAO.GetEmployeePositionForUpdate(db, inputModel)
	if err.Error != nil {
		return
	}

	if positionOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Position)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, positionOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if positionOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.Position)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionService) validateDeleteEmployeePosition(inputStruct *in.EmployeePosition) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
