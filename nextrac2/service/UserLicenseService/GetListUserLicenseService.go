package UserLicenseService

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
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input userLicenseService) InitiateGetListUserLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var countData interface{}
	var searchByParam []in.SearchByParam

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListUserLicenseValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateUserLicense(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListUserLicenseValidOperator,
		CountData:     countData.(int),
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_USER_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseService) doInitiateUserLicense(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var scope map[string]interface{}

	scope, err = input.validateDataScopeUserLicense(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.UserLicenseDAO.GetCountUserLicense(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}
	return
}

func (input userLicenseService) GetListUserLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListUserLicenseValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListUserLicense(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_USER_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseService) doGetListUserLicense(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		scope    map[string]interface{}
	)

	scope, err = input.validateDataScopeUserLicense(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.UserLicenseDAO.GetListUserLicense(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, 0, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponse(dbResult)
	return
}

func (input userLicenseService) validateDataScopeUserLicense(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "validateDataScopeUserLicense"

	output = service.ValidateScope(contextModel, []string{
		constanta.ProvinceDataScope,
		constanta.DistrictDataScope,
		constanta.CustomerGroupDataScope,
		constanta.CustomerCategoryDataScope,
		constanta.SalesmanDataScope,
		constanta.ProductGroupDataScope,
		constanta.ClientTypeDataScope,
	})

	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseService) convertModelToResponse(dbResult []interface{}) (result []out.UserLicenseResponse) {
	for _, resultItem := range dbResult {
		item := resultItem.(repository.UserLicenseModel)
		result = append(result, out.UserLicenseResponse{
			ID:              item.ID.Int64,
			LicenseConfigId: item.LicenseConfigId.Int64,
			CustomerName:    item.CustomerName.String,
			UniqueId1:       item.UniqueId1.String,
			UniqueId2:       item.UniqueId2.String,
			InstallationId:  item.InstallationId.Int64,
			TotalLicense:    item.TotalLicense.Int64,
			TotalActivated:  item.TotalActivated.Int64,
		})
	}

	return result
}
