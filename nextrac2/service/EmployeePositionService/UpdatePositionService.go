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

func (input employeePositionService) UpdateEmployeePosition(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateEmployeePosition"
		inputStruct in.EmployeePosition
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateUpdateEmployeePosition)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateEmployeePosition, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input employeePositionService) doUpdateEmployeePosition(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName     = "doUpdateEmployeePosition"
		inputStruct  = inputStructInterface.(in.EmployeePosition)
		inputModel   = input.convertDTOToModelUpdate(inputStruct, contextModel, timeNow)
		db           = serverconfig.ServerAttribute.DBConnection
		positionOnDB repository.EmployeePositionModel
	)

	positionOnDB, err = dao.EmployeePositionDAO.GetEmployeePositionForUpdate(db, repository.EmployeePositionModel{
		ID: inputModel.ID,
	})
	if err.Error != nil {
		return
	}

	if positionOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Position)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, positionOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if positionOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Position)
		return
	}

	//--- Check Company
	err = input.checkCompany(tx, inputModel.CompanyID.Int64)
	if err.Error != nil {
		return
	}

	//--- Update Position
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.EmployeePositionDAO.TableName, inputModel.ID.Int64, 0)...)
	err = dao.EmployeePositionDAO.UpdateEmployeePosition(tx, inputModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionService) convertDTOToModelUpdate(inputStruct in.EmployeePosition, contextModel *applicationModel.ContextModel, timeNow time.Time) repository.EmployeePositionModel {
	return repository.EmployeePositionModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Name:          sql.NullString{String: inputStruct.Name},
		Description:   sql.NullString{String: inputStruct.Description},
		CompanyID:     sql.NullInt64{Int64: inputStruct.CompanyID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}
}

func (input employeePositionService) ValidateUpdateEmployeePosition(inputStruct *in.EmployeePosition) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
