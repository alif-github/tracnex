package AbsentService

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
)

func (input absentService) InitiateAbsent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam           []in.SearchByParam
		countData               interface{}
		enumDataPeriod          []string
		db                      = serverconfig.ServerAttribute.DBConnection
		tempAbsentValidOperator = make(map[string]applicationModel.DefaultOperator)
		tempValidSearchBy       = []string{"id_card", "name"}
	)

	tempAbsentValidOperator["id_card"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	tempAbsentValidOperator["name"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	//--- Read And Validate
	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListAbsentValidOperator)
	if err.Error != nil {
		return
	}

	//--- Check Period
	if err = input.PeriodCheck(&searchByParam); err.Error != nil {
		return
	}

	countData, err = dao.AbsentDAO.GetCountAbsent(db, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	if countData == nil {
		countData = 0
	}

	enumDataPeriod, err = input.createEnumPeriod()
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: tempValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: tempAbsentValidOperator,
		EnumData:      enumDataPeriod,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentService) createEnumPeriod() (result []string, err errorModel.ErrorModel) {
	return dao.AbsentDAO.GetPeriodAbsent(serverconfig.ServerAttribute.DBConnection, true, 100)
}

func (input absentService) addedNewPeriodSearchByParam(periodLast repository.AbsentModel, searchByParam *[]in.SearchByParam) {
	*searchByParam = append(*searchByParam,
		in.SearchByParam{
			SearchKey:      "period_start",
			DataType:       "char",
			SearchOperator: "eq",
			SearchValue:    periodLast.PeriodStart.Time.Format(constanta.DefaultTimeFormat),
			SearchType:     constanta.Filter,
		},
		in.SearchByParam{
			SearchKey:      "period_end",
			DataType:       "char",
			SearchOperator: "eq",
			SearchValue:    periodLast.PeriodEnd.Time.Format(constanta.DefaultTimeFormat),
			SearchType:     constanta.Filter,
		})
}
