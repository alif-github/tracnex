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
)

func (input employeeContractService) ViewEmployeeContract(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.EmployeeContractRequest
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewEmployeeContract)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewEmployeeContract(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input employeeContractService) doViewEmployeeContract(inputStruct in.EmployeeContractRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		fileName = "ViewEmployeeContractService.go"
		funcName = "doViewEmployeeContract"
		db       = serverconfig.ServerAttribute.DBConnection
		dbResult repository.EmployeeContractModel
	)

	dbResult, err = dao.EmployeeContractDAO.ViewEmployeeContract(db, repository.EmployeeContractModel{ID: sql.NullInt64{Int64: inputStruct.ID}})
	if err.Error != nil {
		return
	}

	if dbResult.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, dbResult.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(dbResult)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractService) convertModelToResponseDetail(inputModel repository.EmployeeContractModel) out.GetListEmployeeContractResponse {
	return out.GetListEmployeeContractResponse{
		ID:          inputModel.ID.Int64,
		ContractNo:  inputModel.ContractNo.String,
		Information: inputModel.Information.String,
		EmployeeID:  inputModel.EmployeeID.Int64,
		FromDate:    inputModel.FromDate.Time,
		ThruDate:    inputModel.ThruDate.Time,
		CreatedName: inputModel.CreatedName.String,
		CreatedAt:   inputModel.CreatedAt.Time,
		UpdatedName: inputModel.UpdatedName.String,
		UpdatedAt:   inputModel.UpdatedAt.Time,
	}
}

func (input employeeContractService) validateViewEmployeeContract(inputStruct *in.EmployeeContractRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
