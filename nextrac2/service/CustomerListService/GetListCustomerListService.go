package CustomerListService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input customerListService) InitiateGetListCustomerList(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListCustomerValidOperator)
	if err.Error != nil {return}

	countData, err = input.doInitiateGetListCustomerList(searchByParam, *contextModel)
	if err.Error != nil {return}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy: 	input.ValidOrderBy,
		ValidSearchBy: 	input.ValidSearchBy,
		ValidLimit: 	input.ValidLimit,
		ValidOperator: 	applicationModel.GetListCustomerValidOperator,
		CountData: 		countData.(int),
	}

	output.Status = out.StatusResponse {
		Code: 		util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message:	GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_CUSTOMER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) doInitiateGetListCustomerList(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (count int, err errorModel.ErrorModel) {
	count, err = dao.CustomerListDAO.GetCountCustomer(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {return}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) GetListCustomerList(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListCustomerValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListCustomerList(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_CUSTOMER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) doGetListCustomerList(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, _ *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}

	dbResult, err = dao.CustomerListDAO.GetListCustomerList(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, 0)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input customerListService) convertToListDTOOut(dbResult []interface{}) (result []out.GetListCustomerResponse) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.CustomerListModel)
		result = append(result, out.GetListCustomerResponse{
			ID: 			repo.ID.Int64,
			CompanyID: 		repo.CompanyID.String,
			BranchID: 		repo.BranchID.String,
			CompanyName: 	repo.CompanyName.String,
			Product: 		repo.Product.String,
			UserAmount: 	repo.UserAmount.Int64,
			ExpDate: 		repo.ExpDate.Time,
			UpdatedAt: 		repo.UpdatedAt.Time,
		})
	}

	return result
}