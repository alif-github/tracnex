package EmployeeService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

func (input employeeService) InitiateGetListEmployeeLeaveYearly(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name"}
		validOrderBy  = []string{"e.id", "e.created_at"}
		key           = request.URL.Query().Get("key")
		keyword       = request.URL.Query().Get("keyword")
		year          = request.URL.Query().Get("year")
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListEmployeeLeaveValidOperator)
	if errModel.Error != nil {
		return
	}

	result, errModel := dao.EmployeeLeaveDAO.InitiateEmployeeLeaveYearly(serverconfig.ServerAttribute.DBConnection, searchByParam, repository.EmployeeLeaveModel{
		SearchBy: sql.NullString{String: key},
		Keyword:  sql.NullString{String: keyword},
		IsYearly: sql.NullBool{Bool: true},
		Year:     sql.NullString{String: year},
	})
	if errModel.Error != nil {
		return
	}

	results, _ := dao.EmployeeHistoryLeaveDAO.GetYearForFilter(serverconfig.ServerAttribute.DBConnection)
	if len(results) == 0 {
		now := time.Now()
		year, _, _ := now.Date()
		results = append(results, strconv.Itoa(year))
	}

	output.Other = results
	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: validSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeLeaveValidOperator,
		CountData:     result,
		ValidSearchParam: []out.SearchByParam{
			{
				Key:   "e.id_card",
				Value: "NIK",
			},
			{
				Key:   "employee_name",
				Value: "Nama",
			},
			{
				Key:   "d.name",
				Value: "Departemen",
			},
		},
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) GetListEmployeeLeaveYearly(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "el.type"}
		validOrderBy  = []string{"e.id", "e.created_at"}
		key           = request.URL.Query().Get("key")
		keyword       = request.URL.Query().Get("keyword")
		year          = request.URL.Query().Get("year")
		content       interface{}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeLeaveValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	content, errModel = input.DoGetListEmployeeLeaveYearly(inputStruct, searchByParam, key, keyword, year)
	if errModel.Error != nil {
		return
	}

	output.Data.Content = content
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) DoGetListEmployeeLeaveYearly(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, key, keyword, year string) (content interface{}, err errorModel.ErrorModel) {
	var (
		results    []interface{}
		modelLeave repository.EmployeeLeaveModel
		db         = serverconfig.ServerAttribute.DBConnection
	)

	//--- For Report
	if inputStruct.Page == 0 {
		inputStruct.Page = -99
		inputStruct.Limit = -99
	}

	//--- Model Employee Leave
	modelLeave = repository.EmployeeLeaveModel{
		SearchBy: sql.NullString{String: key},
		Keyword:  sql.NullString{String: keyword},
		IsYearly: sql.NullBool{Bool: true},
		Year:     sql.NullString{String: year},
	}

	//--- Get Leave Yearly
	results, err = dao.EmployeeLeaveDAO.GetEmployeeLeaveYearly(db, inputStruct, searchByParam, modelLeave)
	if err.Error != nil {
		return
	}

	content = input.toEmployeeLeaveYearlyResponse(results)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) toEmployeeLeaveYearlyResponse(data []interface{}) (result []out.EmployeeLeaveYearly) {
	for _, item := range data {
		model, _ := item.(repository.EmployeeLeaveModel)

		if model.CurrentAnnualLeave.Int64 < 0 {
			model.OwingLeave.Int64 = -1 * model.CurrentAnnualLeave.Int64
			model.CurrentAnnualLeave.Int64 = 0
		}

		res := out.EmployeeLeaveYearly{
			ID:                     model.ID.Int64,
			IDCard:                 model.IDCard.String,
			FirstName:              model.Firstname.String,
			LastName:               model.Lastname.String,
			FullName:               model.Firstname.String + " " + model.Lastname.String,
			Department:             model.Department.String,
			Level:                  model.Level.String,
			Grade:                  model.Grade.String,
			CurrentLeaveThisPeriod: model.CurrentAnnualLeave.Int64,
			LastLeaveBeforePeriod:  model.LastAnnualLeave.Int64,
			OwingLeave:             model.OwingLeave.Int64,
		}

		result = append(result, res)
	}

	return
}
