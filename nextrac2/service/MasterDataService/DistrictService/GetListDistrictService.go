package DistrictService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
)

func (input districtService) GetListDistrict(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListDistrictValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListDistrict(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse {
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_DISTRICT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtService) doGetListDistrict(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var districtStruct in.DistrictRequest
	var pageLimit in.AbstractDTO

	pageLimit = in.AbstractDTO{
		Page:  inputStruct.Page,
		Limit: inputStruct.Limit,
	}

	districtStruct = in.DistrictRequest{AbstractDTO: pageLimit}

	var intParse int
	for _, searchByParamValue := range searchByParam {
		switch searchByParamValue.SearchKey {
		case "province_id":
			intParse, _ = strconv.Atoi(searchByParamValue.SearchValue)
			districtStruct.ProvinceID = int64(intParse)
		case "id":
			intParse, _ = strconv.Atoi(searchByParamValue.SearchValue)
			districtStruct.ID = int64(intParse)
		case "code":
			districtStruct.Code = searchByParamValue.SearchValue
		case "name":
			districtStruct.Name = searchByParamValue.SearchValue
		}
	}

	output, err = master_data_dao.GetListDistrictFromMasterData(districtStruct, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}