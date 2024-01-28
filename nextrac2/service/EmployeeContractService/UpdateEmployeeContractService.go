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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input employeeContractService) UpdateEmployeeContract(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.EmployeeContractRequest
		funcName    = "UpdateEmployeeContract"
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateUpdateEmployeeContract)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateEmployeeContract, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractService) doUpdateEmployeeContract(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.EmployeeContractRequest)
		repoModel   repository.EmployeeContractModel
	)

	//--- Convert To DTO Model
	repoModel = input.convertDTOToModel(inputStruct, contextModel, timeNow, true)

	//--- Validate Employee Contract On DB
	err = input.validateEmployeeContractOnDB(inputStruct, repoModel, contextModel)
	if err.Error != nil {
		return
	}

	//--- Store Audit
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.EmployeeContractDAO.TableName, repoModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	//--- Update Employee Contract
	err = dao.EmployeeContractDAO.UpdateEmployeeContract(tx, repoModel)
	if err.Error != nil {
		err = CheckDuplicateError(err)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractService) ValidateUpdateEmployeeContract(inputStruct *in.EmployeeContractRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
