package AbsentService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

func (input absentService) GetListAbsentPeriod(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"period"}
		validOrderBy  = []string{"period"}
		validLimit    = service.DefaultLimit
		validOperator = make(map[string]applicationModel.DefaultOperator)
	)

	validOperator["period"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, validOperator, validLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListAbsentPeriod(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentService) doGetListAbsentPeriod(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	dbResult, err = dao.AbsentDAO.GetListPeriodAbsent(db, inputStruct, searchByParam)
	if err.Error != nil {
		return
	}

	output = input.convertToListPeriodDTOOut(dbResult)
	return
}

func (input absentService) convertToListPeriodDTOOut(dbResult []interface{}) (result []string) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(string)
		result = append(result, repo)
	}

	return result
}
