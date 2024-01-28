package LicenseTypeService

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

func (input licenseTypeService) InitiateLicenseType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListLicenseTypeValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateLicenseType(searchByParam, contextModel)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListLicenseTypeValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input licenseTypeService) doInitiateLicenseType(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var createdBy int64
	output = 0
	createdBy = contextModel.LimitedByCreatedBy

	output, err = dao.LicenseTypeDAO.GetCountLicenseType(serverconfig.ServerAttribute.DBConnection, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	return
}

func (input licenseTypeService) GetListLicenseType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListLicenseTypeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}
	
	output.Data.Content, err = input.doGetListLicenseType(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input licenseTypeService) doGetListLicenseType(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}

	dbResult, err = dao.LicenseTypeDAO.GetListLicenseType(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return 
	}
	
	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input licenseTypeService) convertModelToResponseGetList(dbResult []interface{}) (result []out.LicenseTypeResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.LicenseTypeModel)
		result = append(result, out.LicenseTypeResponse{
			ID:              item.ID.Int64,
			LicenseTypeName: item.LicenseTypeName.String,
			LicenseTypeDesc: item.LicenseTypeDesc.String,
			CreatedAt:       item.CreatedAt.Time,
			UpdatedAt:       item.UpdatedAt.Time,
			UpdatedBy:       item.UpdatedBy.Int64,
			UpdatedName:     item.UpdatedName.String,
		})
	}

	return result
}
