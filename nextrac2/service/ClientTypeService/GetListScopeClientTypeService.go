package ClientTypeService

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

func (input clientTypeService) GetListScopeClientType(listScope map[string][]string) (result []out.DetailScope, err errorModel.ErrorModel) {
	var searchParam []in.SearchByParam
	limit := 10
	lenScope := len(listScope[constanta.ClientTypeDataScope])
	arrScope := listScope[constanta.ClientTypeDataScope]
	page := math.Ceil(float64(lenScope) / float64(limit))

	inputStruct := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    1,
			Limit:   limit,
			OrderBy: "client_type",
		},
	}

	for i := 0; i < int(page); i++ {
		var dbResult []interface{}
		var newListScope = make(map[string]interface{})
		var appenedArr []interface{}

		newListScope[constanta.ClientTypeDataScope] = make(map[string]interface{})
		newListScope[constanta.ClientTypeDataScope] = []interface{}{}

		if i == (int(page) - 1) {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):])
			newListScope[constanta.ClientTypeDataScope] = appenedArr
		} else {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):((i + 1) * limit)])
			newListScope[constanta.ClientTypeDataScope] = appenedArr
		}

		dbResult, err = dao.ClientTypeDAO.GetListClientType(serverconfig.ServerAttribute.DBConnection, inputStruct, searchParam, newListScope, input.MappingScopeDB, 0)
		result = append(result, input.convertModelToResponseScope(dbResult)...)
	}

	return
}

func (input clientTypeService) convertModelToResponseScope(inputStruct []interface{}) (result []out.DetailScope) {
	for _, itemStruct := range inputStruct {
		item := itemStruct.(repository.ClientTypeModel)
		result = append(result, out.DetailScope{
			Label: constanta.ClientTypeDataScope + ":" + strconv.Itoa(int(item.ID.Int64)),
			Value: out.ScopeValue{
				ID:   item.ID.Int64,
				Name: item.ClientType.String,
			},
		})
	}
	return result
}
