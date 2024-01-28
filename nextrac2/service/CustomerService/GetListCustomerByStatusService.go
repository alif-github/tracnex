package CustomerService

import (
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input customerService) doInitiateCustomerByStatusParent(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isParent bool) (output interface{}, err errorModel.ErrorModel) {
	var createdBy int64
	var scope map[string]interface{}
	output = 0
	createdBy = contextModel.LimitedByCreatedBy

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.CustomerDAO.GetCountCustomerByStatusParent(serverconfig.ServerAttribute.DBConnection, searchByParam, createdBy, scope, input.MappingScopeDB, isParent)
	return
}

func (input customerService) doGetListCustomerByStatusParent(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isParent bool) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}
	var scope map[string]interface{}

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.CustomerDAO.GetListCustomerByStatusParent(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB, isParent)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input customerService) convertModelToResponseGetListByStatus(dbResult []interface{}) (result []out.CustomerListByStatusResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.CustomerModel)
		result = append(result, out.CustomerListByStatusResponse{
			ID:                  item.ID.Int64,
			MDBCompanyProfileID: item.MDBCompanyProfileID.Int64,
			Npwp:                item.Npwp.String,
			CustomerName:        item.CustomerName.String,
			Address:             item.Address.String,
			ProvinceID:          item.ProvinceID.Int64,
			ProvinceName:        item.ProvinceName.String,
			DistrictID:          item.DistrictID.Int64,
			DistrictName:        item.DistrictName.String,
			Phone:               item.Phone.String,
			Status:              item.Status.String,
			CreatedBy:           item.CreatedBy.Int64,
			CreatedAt:           item.CreatedAt.Time,
			UpdatedBy:           item.UpdatedBy.Int64,
			UpdatedAt:           item.UpdatedAt.Time,
		})
	}

	return result
}
