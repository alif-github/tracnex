package ProvinceService

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
	"strconv"
	"strings"
)

func (input provinceService) GetListLocalAdminProvince(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListProvinceValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListLocalProvince(inputStruct, searchByParam, contextModel, true)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_PROVINCE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceService) GetListLocalProvince(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListProvinceValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListLocalProvince(inputStruct, searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_PROVINCE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceService) doGetListLocalProvince(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin bool) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult     []interface{}
		scope        map[string]interface{}
		useCreatedBy int64
	)

	if isAdmin {
		scope = make(map[string]interface{})
		scope[constanta.ProvinceDataScope] = []interface{}{"all"}
		useCreatedBy = 0
	} else {
		scope, err = input.validateDataScope(contextModel)
		if err.Error != nil {
			return
		}
		useCreatedBy = contextModel.LimitedByCreatedBy
	}

	dbResult, err = dao.ProvinceDAO.GetListProvinceLocal(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, scope, input.MappingScopeDB, useCreatedBy)
	if err.Error != nil {
		return
	}

	output, err = input.convertToListLocalProvince(dbResult)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceService) convertToListLocalProvince(dbResult []interface{}) (result []out.ProvinceLocalResponse, err errorModel.ErrorModel) {
	fileName := "GetListLocalProvinceService.go"
	funcName := "convertToListLocalProvince"

	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.ListLocalProvinceModel)

		var idDistrictColl []string
		var idDistrictCollInt []int

		repo.DistrictID.String = strings.ReplaceAll(repo.DistrictID.String, "{", "")
		repo.DistrictID.String = strings.ReplaceAll(repo.DistrictID.String, "}", "")

		if repo.DistrictID.String != "NULL" {
			idDistrictColl = strings.Split(repo.DistrictID.String, ",")
			for _, valueStringIDDistrict := range idDistrictColl {
				itemIDDistrictInt, errorS := strconv.Atoi(valueStringIDDistrict)
				if errorS != nil {
					return nil, errorModel.GenerateUnknownError(fileName, funcName, errorS)
				}
				idDistrictCollInt = append(idDistrictCollInt, itemIDDistrictInt)
			}
		}
		result = append(result, out.ProvinceLocalResponse{
			ID:         repo.ID.Int64,
			CountryID:  repo.CountryID.Int64,
			Code:       repo.Code.String,
			Name:       repo.Name.String,
			DistrictID: idDistrictCollInt,
		})
	}

	return result, errorModel.GenerateNonErrorModel()
}
