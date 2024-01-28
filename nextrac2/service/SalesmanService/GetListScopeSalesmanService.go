package SalesmanService

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

func (input salesmanService) GetListScopeSalesman(listScope map[string][]string) (result []out.DetailScope, err errorModel.ErrorModel) {
	var (
		searchParam []in.SearchByParam
		limit       = 10
		lenScope    = len(listScope[constanta.SalesmanDataScope])
		arrScope    = listScope[constanta.SalesmanDataScope]
		page        = math.Ceil(float64(lenScope) / float64(limit))
	)

	inputStruct := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    1,
			Limit:   limit,
			OrderBy: "full_name",
		},
	}

	for i := 0; i < int(page); i++ {
		var (
			dbResult     []interface{}
			newListScope = make(map[string]interface{})
			appendArr    []interface{}
		)

		newListScope[constanta.SalesmanDataScope] = make(map[string]interface{})
		newListScope[constanta.SalesmanDataScope] = []interface{}{}

		if i == (int(page) - 1) {
			appendArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):])
			newListScope[constanta.SalesmanDataScope] = appendArr
		} else {
			appendArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):((i + 1) * limit)])
			newListScope[constanta.SalesmanDataScope] = appendArr
		}

		dbResult, err = dao.SalesmanDAO.GetListSalesman(serverconfig.ServerAttribute.DBConnection, inputStruct, searchParam, 0, newListScope, input.MappingScopeDB, false)
		result = append(result, input.convertModelToResponseScope(dbResult)...)
	}

	return
}

func (input salesmanService) convertModelToResponseScope(inputStruct []interface{}) (result []out.DetailScope) {
	for _, itemStruct := range inputStruct {
		item := itemStruct.(repository.ListSalesmanModel)
		result = append(result, out.DetailScope{
			Label: constanta.SalesmanDataScope + ":" + strconv.Itoa(int(item.ID.Int64)),
			Value: out.ScopeValue{
				ID:   item.ID.Int64,
				Name: item.FirstName.String,
			},
		})
	}
	return result
}
