package CustomerGroupService

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

func (input customerGroupService) GetListScopeCustomerGroup(listScope map[string][]string) (result []out.DetailScope, err errorModel.ErrorModel) {
	var searchParam []in.SearchByParam
	limit := 10
	lenScope := len(listScope[constanta.CustomerGroupDataScope])
	arrScope := listScope[constanta.CustomerGroupDataScope]
	page := math.Ceil(float64(lenScope) / float64(limit))

	inputStruct := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    1,
			Limit:   limit,
			OrderBy: "customer_group_name",
		},
	}

	for i := 0; i < int(page); i++ {
		var dbResult []interface{}
		var newListScope = make(map[string]interface{})
		var appenedArr []interface{}

		newListScope[constanta.CustomerGroupDataScope] = make(map[string]interface{})
		newListScope[constanta.CustomerGroupDataScope] = []interface{}{}

		if i == (int(page) - 1) {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):])
			newListScope[constanta.CustomerGroupDataScope] = appenedArr
		} else {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):((i + 1) * limit)])
			newListScope[constanta.CustomerGroupDataScope] = appenedArr
		}

		dbResult, err = dao.CustomerGroupDAO.GetListCustomerGroup(serverconfig.ServerAttribute.DBConnection, inputStruct, searchParam, 0, newListScope, input.MappingScopeDB)
		result = append(result, input.convertModelToResponseScope(dbResult)...)
	}

	return
}

func (input customerGroupService) convertModelToResponseScope(inputStruct []interface{}) (result []out.DetailScope) {
	for _, itemStruct := range inputStruct {
		item := itemStruct.(repository.CustomerGroupModel)
		result = append(result, out.DetailScope{
			Label: constanta.CustomerGroupDataScope + ":" + strconv.Itoa(int(item.ID.Int64)),
			Value: out.ScopeValue{
				ID:   item.ID.Int64,
				Name: item.CustomerGroupName.String,
			},
		})
	}
	return result
}
