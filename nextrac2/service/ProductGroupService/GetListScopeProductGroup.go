package ProductGroupService

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

func (input productGroupService) GetListScopeProductGroup(listScope map[string][]string) (result []out.DetailScope, err errorModel.ErrorModel) {
	var searchParam []in.SearchByParam
	limit := 10
	lenScope := len(listScope[constanta.ProductGroupDataScope])
	arrScope := listScope[constanta.ProductGroupDataScope]
	page := math.Ceil(float64(lenScope) / float64(limit))

	inputStruct := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    1,
			Limit:   limit,
			OrderBy: "product_group_name",
		},
	}

	for i := 0; i < int(page); i++ {
		var dbResult []interface{}
		var newListScope = make(map[string]interface{})
		var appenedArr []interface{}

		newListScope[constanta.ProductGroupDataScope] = make(map[string]interface{})
		newListScope[constanta.ProductGroupDataScope] = []interface{}{}

		if i == (int(page) - 1) {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):])
			newListScope[constanta.ProductGroupDataScope] = appenedArr
		} else {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):((i + 1) * limit)])
			newListScope[constanta.ProductGroupDataScope] = appenedArr
		}

		dbResult, err = dao.ProductGroupDAO.GetListProductGroup(serverconfig.ServerAttribute.DBConnection, inputStruct, searchParam, 0, newListScope, input.MappingScopeDB)
		result = append(result, input.convertModelToResponseScope(dbResult)...)
	}

	return
}

func (input productGroupService) convertModelToResponseScope(inputStruct []interface{}) (result []out.DetailScope) {
	for _, itemStruct := range inputStruct {
		item := itemStruct.(repository.ProductGroupModel)
		result = append(result, out.DetailScope{
			Label: constanta.ProductGroupDataScope + ":" + strconv.Itoa(int(item.ID.Int64)),
			Value: out.ScopeValue{
				ID:   item.ID.Int64,
				Name: item.ProductGroupName.String,
			},
		})
	}
	return result
}
