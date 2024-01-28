package EmployeeContractService

import (
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

func (input employeeContractService) InitiateEmployeeContract(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
		scope         map[string]interface{}
		db            = serverconfig.ServerAttribute.DBConnection
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeContractValidOperator)
	if err.Error != nil {
		return
	}

	scope, err = input.ValidateMultipleDataScope(contextModel, []string{constanta.EmployeeDataScope})
	if err.Error != nil {
		return
	}

	//--- Count Data
	countData, err = dao.EmployeeContractDAO.GetCountEmployeeContract(db, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeContractValidOperator,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractService) GetListEmployeeContract(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeContractValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListEmployeeContract(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeContractService) doGetListEmployeeContract(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult  []interface{}
		scope     map[string]interface{}
		createdBy int64
		db        = serverconfig.ServerAttribute.DBConnection
	)

	scope, err = input.ValidateMultipleDataScope(contextModel, []string{constanta.EmployeeDataScope})
	if err.Error != nil {
		return
	}

	dbResult, err = dao.EmployeeContractDAO.GetListEmployeeContract(db, inputStruct, searchByParam, createdBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input employeeContractService) convertModelToResponseGetList(dbResult []interface{}) (result interface{}) {
	var resultGetList []out.GetListEmployeeContractResponse
	for _, dbResultItem := range dbResult {
		var item = dbResultItem.(repository.EmployeeContractModel)

		//--- Get List
		resultGetList = append(resultGetList, out.GetListEmployeeContractResponse{
			ID:          item.ID.Int64,
			ContractNo:  item.ContractNo.String,
			Information: item.Information.String,
			EmployeeID:  item.EmployeeID.Int64,
			FromDate:    item.FromDate.Time,
			ThruDate:    item.ThruDate.Time,
			CreatedName: item.CreatedName.String,
			CreatedAt:   item.CreatedAt.Time,
			UpdatedName: item.UpdatedName.String,
			UpdatedAt:   item.UpdatedAt.Time,
		})
	}

	result = resultGetList
	return
}
