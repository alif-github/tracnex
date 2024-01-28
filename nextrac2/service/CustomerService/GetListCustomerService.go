package CustomerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input customerService) InitiateCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListMasterCustomerValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateCustomer(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListMasterCustomerValidOperator,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerService) doInitiateCustomer(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		scope     map[string]interface{}
		createdBy int64
	)

	output = 0
	createdBy = contextModel.LimitedByCreatedBy
	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.CustomerDAO.GetCountCustomer(serverconfig.ServerAttribute.DBConnection, searchByParam, createdBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerService) GetListCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListMasterCustomerValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListCustomer(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerService) doGetListCustomer(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		scope    map[string]interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.CustomerDAO.GetListCustomer(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input customerService) convertModelToResponseGetList(dbResult []interface{}) (result []out.CustomerListResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.CustomerModel)
		result = append(result, out.CustomerListResponse{
			ID:           item.ID.Int64,
			Npwp:         item.Npwp.String,
			CustomerName: item.CustomerName.String,
			Address:      item.Address.String,
			ProvinceID:   item.ProvinceID.Int64,
			ProvinceName: item.ProvinceName.String,
			DistrictID:   item.DistrictID.Int64,
			DistrictName: item.DistrictName.String,
			Phone:        item.Phone.String,
			Status:       item.Status.String,
			CreatedBy:    item.CreatedBy.Int64,
			CreatedAt:    item.CreatedAt.Time,
			UpdatedBy:    item.UpdatedBy.Int64,
			UpdatedAt:    item.UpdatedAt.Time,
		})
	}

	return result
}
