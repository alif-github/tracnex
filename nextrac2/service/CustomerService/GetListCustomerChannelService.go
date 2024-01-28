package CustomerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

func (input customerService) GetListCustomerNonParent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListMasterCustomerValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListCustomerByStatusParent(inputStruct, searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input customerService) InitiateCustomerNonParent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListMasterCustomerValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateCustomerByStatusParent(searchByParam, contextModel, false)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListMasterCustomerValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input customerService) GetListCustomerParent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListMasterCustomerValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListCustomerByStatusParent(inputStruct, searchByParam, contextModel, true)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input customerService) InitiateCustomerParent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListMasterCustomerValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateCustomerByStatusParent(searchByParam, contextModel, true)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListMasterCustomerValidOperator,
		CountData:     countData.(int),
	}
	return
}