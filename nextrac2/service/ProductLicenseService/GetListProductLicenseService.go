package ProductLicenseService

import (
	"database/sql"
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

func (input productLicenseService) InitiateGetListProductLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var countData interface{}
	var searchByParam []in.SearchByParam

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListProductLicenseValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListProductLicense(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidOperator: applicationModel.GetListProductLicenseValidOperator,
		CountData:     countData.(int),
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_PRODUCT_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) doInitiateListProductLicense(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		scope map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.ProductLicenseDAO.GetCountProductLicense(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy,
		scope, input.MappingScopeDB, input.ListScope)
	if err.Error != nil {
		return
	}
	return
}

func (input productLicenseService) GetListProductLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListProductLicenseValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListProductLicense(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input productLicenseService) doGetListProductLicense(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		scope    map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.ProductLicenseDAO.GetListProductLicense(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, repository.ProductLicenseModel{
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}, scope, input.MappingScopeDB, input.ListScope)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponse(dbResult)

	return
}

func (input productLicenseService) convertModelToResponse(dbResult []interface{}) (result []out.ProductLicenseResponse) {
	for _, resultItem := range dbResult {
		item := resultItem.(repository.ProductLicenseModelForView)
		result = append(result, out.ProductLicenseResponse{
			ID:             item.ID.Int64,
			LicenseConfig:  item.LicenseConfigId.Int64,
			CustomerName:   item.CustomerName.String,
			UniqueId1:      item.UniqueId1.String,
			UniqueId2:      item.UniqueId2.String,
			InstallationId: item.InstallationId.Int64,
			ProductName:    item.ProductName.String,
			LicenseVariant: item.LicenseVariantName.String,
			LicenseType:    item.LicenseTypeName.String,
			ValidFrom:      item.ProductValidFrom.Time,
			ValidThru:      item.ProductValidThru.Time,
			Status:         item.LicenseStatus.Int32,
		})
	}

	return result
}
