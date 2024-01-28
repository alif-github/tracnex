package EmployeePositionService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

func (input employeePositionService) InsertEmployeePosition(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertEmployeePosition"
		inputStruct in.EmployeePosition
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertEmployeePosition)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertEmployeePosition, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input employeePositionService) doInsertEmployeePosition(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.EmployeePosition)
		inputModel  = input.convertDTOToModel(inputStruct, *contextModel, timeNow)
		resultID    int64
	)

	//--- Check Company
	err = input.checkCompany(tx, inputModel.CompanyID.Int64)
	if err.Error != nil {
		return
	}

	//--- Insert Employee Position
	resultID, err = dao.EmployeePositionDAO.InsertEmployeePosition(tx, inputModel)
	if err.Error != nil {
		return
	}

	//--- Data Audit
	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeePositionDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: resultID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionService) convertDTOToModel(inputStruct in.EmployeePosition, contextModel applicationModel.ContextModel, timeNow time.Time) repository.EmployeePositionModel {
	return repository.EmployeePositionModel{
		Name:          sql.NullString{String: inputStruct.Name},
		Description:   sql.NullString{String: inputStruct.Description},
		CompanyID:     sql.NullInt64{Int64: inputStruct.CompanyID},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}
}

func (input employeePositionService) ValidateInsertEmployeePosition(inputStruct *in.EmployeePosition) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}
