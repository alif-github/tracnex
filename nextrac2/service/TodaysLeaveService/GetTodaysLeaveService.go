package TodaysLeaveService

import (
	"encoding/json"
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
	"strings"
	"time"
)

func (input todaysLeaveService) InitiateTodaysLeave(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
		db            = serverconfig.ServerAttribute.DBConnection
		timeNow       = time.Now()
	)

	//--- Read And Validate
	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListTodaysLeaveValidOperator)
	if err.Error != nil {
		return
	}

	//--- Validate Search By Param Type
	err = input.validateSearchByParamType(searchByParam)
	if err.Error != nil {
		return
	}

	//--- Count Data
	countData, err = dao.EmployeeLeaveDAO.GetCountTodayLeave(db, searchByParam, timeNow)
	if err.Error != nil {
		return
	}

	//--- Create Enum For Type
	enumType := make(map[string][]string)
	enumType["type"] = []string{"leave", "sick", "permit"}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		EnumData:      enumType,
		ValidOperator: applicationModel.GetListTodaysLeaveValidOperator,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input todaysLeaveService) GetListTodaysLeave(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListTodaysLeaveValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	//--- Main Get List
	output.Data.Content, err = input.doGetListTodaysLeave(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input todaysLeaveService) doGetListTodaysLeave(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult   []interface{}
		sumResult  []interface{}
		db         = serverconfig.ServerAttribute.DBConnection
		timeNow    = time.Now()
		outputTemp = make(map[string]interface{})
	)

	//--- Validate Search By Param Type
	err = input.validateSearchByParamType(searchByParam)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.EmployeeLeaveDAO.GetListTodayLeave(db, inputStruct, searchByParam, timeNow)
	if err.Error != nil {
		return
	}

	sumResult, err = dao.EmployeeLeaveDAO.GetSummaryLeaveToday(db, timeNow)
	if err.Error != nil {
		return
	}

	outputTemp["leave"] = 0
	outputTemp["sick"] = 0
	outputTemp["permit"] = 0

	if len(sumResult) > 0 {
		for _, itemType := range sumResult {
			var (
				typeTemp = itemType.(repository.EmployeeLeaveModel)
				key      = strings.ToLower(typeTemp.Type.String)
			)

			outputTemp[key] = typeTemp.CountType.Int64
		}
	}

	outputTemp["detail"] = input.convertToListDTOOut(dbResult, contextModel)
	output = outputTemp
	return
}

func (input todaysLeaveService) convertToListDTOOut(dbResult []interface{}, contextModel *applicationModel.ContextModel) (result []out.ListTodaysLeave) {
	for _, dbResultItem := range dbResult {
		var (
			dateTemp []string
			repo     = dbResultItem.(repository.EmployeeLeaveModel)
		)

		switch repo.Type.String {
		case "LEAVE":
			repo.Type.String = util.GenerateConstantaI18n(constanta.LeaveStatus, contextModel.AuthAccessTokenModel.Locale, nil)
		case "SICK":
			repo.Type.String = util.GenerateConstantaI18n(constanta.SickStatus, contextModel.AuthAccessTokenModel.Locale, nil)
		case "PERMIT":
			repo.Type.String = util.GenerateConstantaI18n(constanta.PermitStatus, contextModel.AuthAccessTokenModel.Locale, nil)
		default:
		}

		_ = json.Unmarshal([]byte(repo.Date.String), &dateTemp)
		result = append(result, out.ListTodaysLeave{
			IDCard:     repo.IDCard.String,
			Name:       repo.Name.String,
			Department: repo.Department.String,
			Date:       dateTemp,
			Type:       repo.Type.String,
		})
	}

	return
}

func (input todaysLeaveService) validateSearchByParamType(searchByParam []in.SearchByParam) (err errorModel.ErrorModel) {
	var (
		fileName = "GetTodaysLeaveService.go"
		funcName = "validateSearchByParamType"
	)

	for _, item := range searchByParam {
		if item.SearchKey == "type" {
			if item.SearchValue != "leave" && item.SearchValue != "sick" && item.SearchValue != "permit" {
				err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Type)
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
