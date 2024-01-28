package EmployeeService

import (
	"database/sql"
	"github.com/google/uuid"
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

func (input employeeService) DeleteEmployee(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteEmployee"
		inputStruct in.EmployeeRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteEmployee)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteEmployee, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input employeeService) doDeleteEmployee(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct  = inputStructInterface.(in.EmployeeRequest)
		employeeOnDB repository.EmployeeModel
	)

	inputModel := repository.EmployeeModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	//--- Validate Employee On DB
	employeeOnDB, err = input.validateEmployeeOnDB(tx, inputStruct, inputModel, contextModel)
	if err.Error != nil {
		return
	}

	//--- Create Random Uniq NIK
	input.randTokenGenerator(&inputModel, employeeOnDB)

	//--- Store Audit
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.EmployeeDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	//--- Delete Employee
	err = dao.EmployeeDAO.DeleteEmployee(tx, inputModel)
	if err.Error != nil {
		return
	}

	//--- Delete Employee Variable
	if err = dao.EmployeeVariableDAO.DeleteEmployeeVariableByEmployeeID(tx, repository.EmployeeVariableModel{
		EmployeeID:    inputModel.ID,
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}); err.Error != nil {
		return
	}

	//--- Delete Employee Benefits
	if err = dao.EmployeeBenefitsDAO.DeleteEmployeeBenefitsByEmployeeID(tx, repository.EmployeeBenefitsModel{
		EmployeeID:    inputModel.ID,
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}); err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) validateEmployeeOnDB(_ *sql.Tx, inputStruct in.EmployeeRequest, inputModel repository.EmployeeModel, contextModel *applicationModel.ContextModel) (employeeOnDB repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "validateEmployeeOnDB"
		db       = serverconfig.ServerAttribute.DBConnection
	)

	//--- Get Employee ID
	employeeOnDB, err = input.EmployeeDAO.GetEmployeeForUpdate(db, repository.EmployeeModel{ID: inputModel.ID})
	if err.Error != nil {
		return
	}

	if employeeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, employeeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if employeeOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.EmployeeConstanta)
		return
	}

	if employeeOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.EmployeeConstanta)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) validateDeleteEmployee(inputStruct *in.EmployeeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}

func (input employeeService) randTokenGenerator(inputModel *repository.EmployeeModel, employeeOnDB repository.EmployeeModel) {
	encodedStr := service.RandTimeToken(10, uuid.New().String())
	inputModel.IDCard.String = employeeOnDB.IDCard.String + encodedStr
}
