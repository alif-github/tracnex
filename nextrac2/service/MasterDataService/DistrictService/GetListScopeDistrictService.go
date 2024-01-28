package DistrictService

import (
	"math"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

func (input districtService) GetListScopeDistrict(listScope map[string][]string) (result []out.DetailScope, err errorModel.ErrorModel) {
	var searchParam []in.SearchByParam
	limit := 10
	lenScope := len(listScope[constanta.DistrictDataScope])
	arrScope := listScope[constanta.DistrictDataScope]
	page := math.Ceil(float64(lenScope) / float64(limit))

	inputStruct := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    1,
			Limit:   limit,
			OrderBy: "d.name",
		},
	}

	for i := 0; i < int(page); i++ {
		var dbResult []interface{}
		var newListScope = make(map[string]interface{})
		var appenedArr []interface{}

		newListScope[constanta.DistrictDataScope] = make(map[string]interface{})
		newListScope[constanta.DistrictDataScope] = []interface{}{}

		if i == (int(page) - 1) {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):])
			newListScope[constanta.DistrictDataScope] = appenedArr
		} else {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):((i + 1) * limit)])
			newListScope[constanta.DistrictDataScope] = appenedArr
		}

		dbResult, err = dao.DistrictDAO.GetListDistrictLocal(serverconfig.ServerAttribute.DBConnection, inputStruct, searchParam, newListScope, input.MappingScopeDB, 0)
		result = append(result, input.convertModelToResponseScope(dbResult)...)
	}

	return
}

func (input districtService) convertModelToResponseScope(inputStruct []interface{}) (result []out.DetailScope) {
	for _, itemStruct := range inputStruct {
		item := itemStruct.(repository.ListLocalDistrictModel)
		result = append(result, out.DetailScope{
			Label: constanta.DistrictDataScope + ":" + strconv.Itoa(int(item.ID.Int64)),
			Value: out.ScopeValue{
				ID:       item.ID.Int64,
				ParentID: item.ProvinceID.Int64,
				Name:     item.Name.String,
			},
		})
	}
	return result
}
