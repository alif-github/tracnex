package ProvinceService

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

func (input provinceService) GetListProvince(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListProvinceValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListProvince(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse {
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_PROVINCE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceService) doGetListProvince(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var provinceStruct in.ProvinceRequest
	var pageLimit in.AbstractDTO

	pageLimit = in.AbstractDTO{
		Page:  inputStruct.Page,
		Limit: inputStruct.Limit,
	}

	provinceStruct = in.ProvinceRequest{AbstractDTO: pageLimit}

	var intParse int
	for _, searchByParamValue := range searchByParam {
		switch searchByParamValue.SearchKey {
		case "country_id":
			intParse, _ = strconv.Atoi(searchByParamValue.SearchValue)
			provinceStruct.CountryID = int64(intParse)
		case "id":
			intParse, _ = strconv.Atoi(searchByParamValue.SearchValue)
			provinceStruct.ID = int64(intParse)
		case "code":
			provinceStruct.Code = searchByParamValue.SearchValue
		case "name":
			provinceStruct.Name = searchByParamValue.SearchValue
		}
	}

	output, err = master_data_dao.GetListProvinceFromMasterData(provinceStruct, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
