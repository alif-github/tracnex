package EmployeeService

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

func (input employeeService) GetListScopeEmployee(listScope map[string][]string) (result []out.DetailScope, err errorModel.ErrorModel) {
	var (
		searchParam []in.SearchByParam
		limit       = 10
		lenScope    = len(listScope[constanta.EmployeeDataScope])
		arrScope    = listScope[constanta.EmployeeDataScope]
		page        = math.Ceil(float64(lenScope) / float64(limit))
	)

	inputStruct := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    1,
			Limit:   limit,
			OrderBy: "name",
		},
	}

	for i := 0; i < int(page); i++ {
		var (
			dbResult     []interface{}
			appenedArr   []interface{}
			newListScope = make(map[string]interface{})
			db           = serverconfig.ServerAttribute.DBConnection
		)

		newListScope[constanta.EmployeeDataScope] = make(map[string]interface{})
		newListScope[constanta.EmployeeDataScope] = []interface{}{}

		if i == (int(page) - 1) {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):])
			newListScope[constanta.EmployeeDataScope] = appenedArr
		} else {
			appenedArr = service.GetArrayInterfaceFromStringCollection(arrScope[(i * limit):((i + 1) * limit)])
			newListScope[constanta.EmployeeDataScope] = appenedArr
		}

		dbResult, err = dao.EmployeeDAO.GetListEmployee(db, inputStruct, searchParam, 0, newListScope, input.MappingScopeDB, false)
		result = append(result, input.convertModelToResponseScope(dbResult)...)
	}

	return
}

func (input employeeService) convertModelToResponseScope(inputStruct []interface{}) (result []out.DetailScope) {
	for _, itemStruct := range inputStruct {
		item := itemStruct.(repository.EmployeeModel)
		result = append(result, out.DetailScope{
			Label: constanta.EmployeeDataScope + ":" + strconv.Itoa(int(item.ID.Int64)),
			Value: out.ScopeValue{
				ID:   item.ID.Int64,
				Name: item.Name.String,
			},
		})
	}
	return result
}
