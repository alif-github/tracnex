package DistrictService

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
	"nexsoft.co.id/nextrac2/util"
)

func (input districtService) GetListAdminLocalDistrict(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListDistrictValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListLocalDistrict(inputStruct, searchByParam, contextModel, true)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_DISTRICT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtService) GetListLocalDistrict(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListDistrictValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListLocalDistrict(inputStruct, searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_DISTRICT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtService) doGetListLocalDistrict(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin bool) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}
	var scope map[string]interface{}
	var useCreated int64

	if isAdmin {
		scope = make(map[string]interface{})
		scope[constanta.DistrictDataScope] = []interface{}{"all"}
		useCreated = 0
	} else {
		scope, err = input.validateDataScope(contextModel)
		if err.Error != nil {
			return
		}
		useCreated = contextModel.LimitedByCreatedBy
	}

	dbResult, err = dao.DistrictDAO.GetListDistrictLocal(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, scope, input.MappingScopeDB, useCreated)
	if err.Error != nil {
		return
	}

	output = input.convertToListLocalDistrict(dbResult)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtService) convertToListLocalDistrict(dbResult []interface{}) (result []out.DistrictLocalResponse) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.ListLocalDistrictModel)
		result = append(result, out.DistrictLocalResponse{
			ID:         repo.ID.Int64,
			ProvinceID: repo.ProvinceID.Int64,
			Code:       repo.Code.String,
			Name:       repo.Name.String,
		})
	}

	return result
}
