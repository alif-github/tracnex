package CustomerCategoryService

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

func (input customerCategoryService) GetListScopeCustomerCategory(listScope map[string][]string) (result []out.DetailScope, err errorModel.ErrorModel) {
	//CustomerCategoryOnDB, err :=
	var searchParam []in.SearchByParam
	limit := 10
	lenScope := len(listScope[constanta.CustomerCategoryDataScope])
	arrScope := listScope[constanta.CustomerCategoryDataScope]
	page := math.Ceil(float64(lenScope) / float64(limit))

	inputStruct := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    1,
			Limit:   limit,
			OrderBy: "customer_category_name",
		},
	}

	for i := 0; i < int(page); i++ {
		var dbResult []interface{}
		var newListScope = make(map[string]interface{})
		var appenedArr []interface{}

		newListScope[constanta.CustomerCategoryDataScope] = make(map[string]interface{})
		newListScope[constanta.CustomerCategoryDataScope] = []interface{}{}

		if i == (int(page) - 1) {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):])
			newListScope[constanta.CustomerCategoryDataScope] = appenedArr
		} else {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):((i + 1) * limit)])
			newListScope[constanta.CustomerCategoryDataScope] = appenedArr
		}

		dbResult, err = dao.CustomerCategoryDAO.GetListCustomerCategory(serverconfig.ServerAttribute.DBConnection, inputStruct, searchParam, 0, newListScope, input.MappingScopeDB)
		result = append(result, input.convertModelToResponseScope(dbResult)...)
	}

	return
}

func (input customerCategoryService) convertModelToResponseScope(inputStruct []interface{}) (result []out.DetailScope) {
	for _, itemStruct := range inputStruct {
		item := itemStruct.(repository.CustomerCategoryModel)
		result = append(result, out.DetailScope{
			Label: constanta.CustomerCategoryDataScope + ":" + strconv.Itoa(int(item.ID.Int64)),
			Value: out.ScopeValue{
				ID:   item.ID.Int64,
				Name: item.CustomerCategoryID.String,
			},
		})
	}
	return result
}
