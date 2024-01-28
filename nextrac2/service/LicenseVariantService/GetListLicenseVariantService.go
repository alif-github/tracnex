package LicenseVariantService

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

func (input licenseVariantService) GetListLicenseVariant(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListLicenseVariantValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListLicenseVariant(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_LICENSE_VARIANT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) InitiateGetListLicenseVariant(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var countData interface{}
	var searchByParam []in.SearchByParam

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListLicenseVariantValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListLicenseVariant(searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListLicenseVariantValidOperator,
		CountData:     countData.(int),
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_LICENSE_VARIANT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) doGetListLicenseVariant(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}

	dbResult, err = dao.LicenseVariantDAO.GetListLicenseVariant(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertToListLicenseVariant(dbResult)
	return
}

func (input licenseVariantService) doInitiateListLicenseVariant(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = dao.LicenseVariantDAO.GetCountLicenseVariant(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}
	return
}

func (input licenseVariantService) convertToListLicenseVariant(dbResult []interface{}) (result []out.LicenseVariantListResponse) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.LicenseVariantListModel)
		result = append(result, out.LicenseVariantListResponse{
			ID:                 repo.ID.Int64,
			LicenseVariantName: repo.LicenseVariantName.String,
			CreatedAt:          repo.CreatedAt.Time,
			UpdatedBy:          repo.UpdatedBy.Int64,
			UpdatedName:        repo.UpdatedName.String,
			UpdatedAt:          repo.UpdatedAt.Time,
		})
	}

	return result
}
