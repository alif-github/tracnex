package CustomerCategoryService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

func (input customerCategoryService) GetListCustomerCategoryByUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListCustomerCategoryValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListCustomerCategory(inputStruct, searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input customerCategoryService) GetListCustomerCategoryByAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListCustomerCategoryValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListCustomerCategory(inputStruct, searchByParam, contextModel, true)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}
