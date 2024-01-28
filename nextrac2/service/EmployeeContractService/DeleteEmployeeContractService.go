package EmployeeContractService

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

func (input employeeContractService) DeleteEmployeeContract(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteEmployeeContract"
		inputStruct in.EmployeeContractRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteEmployeeContract)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteEmployeeContract, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input employeeContractService) doDeleteEmployeeContract(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var inputStruct = inputStructInterface.(in.EmployeeContractRequest)
	inputModel := repository.EmployeeContractModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	//--- Validate Employee Contract On DB
	err = input.validateEmployeeContractOnDB(inputStruct, inputModel, contextModel)
	if err.Error != nil {
		return
	}

	//--- Store Audit
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.EmployeeContractDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	//--- Delete Employee Contract
	err = dao.EmployeeContractDAO.DeleteEmployeeContract(tx, inputModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractService) validateEmployeeContractOnDB(inputStruct in.EmployeeContractRequest, inputModel repository.EmployeeContractModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName     = "DeleteEmployeeContractService.go"
		funcName     = "validateEmployeeContractOnDB"
		db           = serverconfig.ServerAttribute.DBConnection
		contractOnDB repository.EmployeeContractModel
	)

	//--- Get Employee Contract On DB
	contractOnDB, err = dao.EmployeeContractDAO.GetEmployeeContractForUpdate(db, inputModel)
	if err.Error != nil {
		return
	}

	if contractOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, contractOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if contractOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(fileName, funcName, "Contract")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractService) validateDeleteEmployeeContract(inputStruct *in.EmployeeContractRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
