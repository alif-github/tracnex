package EmployeeService

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
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"strings"
	"time"
)

func (input employeeService) InitiateGetListEmployeeReimbursement(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "id_card"}
		validOrderBy  = []string{"id", "created_at"}
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListEmployeeReimbursementValidOperator)
	if errModel.Error != nil {
		return
	}

	model, errModel := input.getEmployeeReimbursementFilters(request)
	if errModel.Error != nil {
		return
	}

	result, errModel := dao.EmployeeReimbursementDAO.InitiateGetListEmployeeReimbursement(serverconfig.ServerAttribute.DBConnection, searchByParam, model)
	if errModel.Error != nil {
		return
	}

	results, _ := dao.EmployeeHistoryLeaveDAO.GetYearForFilter(serverconfig.ServerAttribute.DBConnection)
	if len(results) == 0 {
		now := time.Now()
		year , _, _ := now.Date()
		results = append(results, strconv.Itoa(year))
	}

	output.Other = results

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: validSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeReimbursementValidOperator,
		CountData:     result,
		ValidSearchParam: []out.SearchByParam{
			{
				Key:  "id_card",
				Value: "NIK",
			},
			{
				Key : "name",
				Value : "Nama",
			},
			{
				Key : "department",
				Value : "Departemen",
			},
		},
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) InitiateGetListEmployeeReimbursementReport(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "id_card"}
		validOrderBy  = []string{"e.id", "e.created_at"}
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListEmployeeReimbursementValidOperator)
	if errModel.Error != nil {
		return
	}

	model, errModel := input.getEmployeeReimbursementFilters(request)
	if errModel.Error != nil {
		return
	}

	model.IsFilter.Bool = false
	result, errModel := dao.EmployeeReimbursementDAO.InitiateGetListEmployeeReimbursementReport(serverconfig.ServerAttribute.DBConnection, searchByParam, model)
	if errModel.Error != nil {
		return
	}

	results, _ := dao.EmployeeHistoryLeaveDAO.GetYearForFilter(serverconfig.ServerAttribute.DBConnection)
	if len(results) == 0 {
		now := time.Now()
		year , _, _ := now.Date()
		results = append(results, strconv.Itoa(year))
	}

	output.Other = results

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: validSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeReimbursementValidOperator,
		CountData:     result,
		ValidSearchParam: []out.SearchByParam{
			{
				Key:  "id_card",
				Value: "NIK",
			},
			{
				Key : "name",
				Value : "Nama",
			},
			{
				Key : "department",
				Value : "Departemen",
			},
		},
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) GetListEmployeeReimbursement(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "id_card"}
		validOrderBy  = []string{"er.id", "er.created_at"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeReimbursementValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	model, errModel := input.getEmployeeReimbursementFilters(request)
	if errModel.Error != nil {
		return
	}

	results, errModel := dao.EmployeeReimbursementDAO.GetListEmployeeReimbursement(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, model)
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.toEmployeeReimbursementResponse(results)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) GetListEmployeeReimbursementReport(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "id_card"}
		validOrderBy  = []string{"e.id", "e.created_at"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeReimbursementValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	model, errModel := input.getEmployeeReimbursementFilters(request)
	if errModel.Error != nil {
		return
	}

	model.IsFilter.Bool = false
	now := time.Now()
	year , _, _ := now.Date()

	if model.Year.String == "" && model.Month.String == ""{
		model.Year.String = strconv.Itoa(year)
	}

	results, errModel := dao.EmployeeReimbursementDAO.GetListEmployeeReimbursementReport(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, model)
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.responseReimbursementReport(results)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) getEmployeeReimbursementFilters(request *http.Request) (result repository.EmployeeReimbursement, errModel errorModel.ErrorModel) {
	funcName := "getEmployeeReimbursementFilters"

	q := request.URL.Query()

	result.SearchBy.String = q.Get("key")
	result.Keyword.String = q.Get("keyword")
	result.FullName.String = q.Get("full_name")
	result.IDCard.String = q.Get("id_card")
	result.Department.String = q.Get("department")
	result.VerifiedStatus.String = q.Get("verified_status")
	result.Status.String = q.Get("status")
	result.Year.String = q.Get("year")
	result.Month.String = q.Get("month")
	result.ReportType.String = q.Get("report_type")

	now := time.Now()
	year , _, _ := now.Date()

	if result.Year.String == "" && result.Month.String != ""{
		result.Year.String = strconv.Itoa(year)
	}

	result.IsFilter.Bool = true

	if result.VerifiedStatus.String != "" {
		if result.VerifiedStatus.String != constanta.VerifiedReimbursementVerification && result.VerifiedStatus.String != constanta.UnverifiedReimbursementVerification && result.VerifiedStatus.String != constanta.PendingReimbursementVerification {
			errModel = errorModel.GenerateFormatFieldError(input.FileName, funcName, "verified_status")
			return
		}
	}

	if result.Status.String != "" && !input.isStatusValid(result.Status.String) {
		errModel = errorModel.GenerateFormatFieldError(input.FileName, funcName, "status")
		return
	}

	strStartDate := q.Get("start_date")
	strEndDate := q.Get("end_date")

	if strStartDate == "" || strEndDate == "" {
		return
	}

	startDate, err := time.Parse("2006-01-02", strStartDate)
	if err != nil {
		errModel = errorModel.GenerateFormatFieldError(input.FileName, funcName, "start_date")
		return
	}

	endDate, err := time.Parse("2006-01-02", strEndDate)
	if err != nil {
		errModel = errorModel.GenerateFormatFieldError(input.FileName, funcName, "end_date")
		return
	}

	if endDate.Before(startDate) {
		errModel = errorModel.GenerateDateCannotBeLessThanError(input.FileName, funcName, "end_date", "start_date")
		return
	}

	result.StartDate.String = strStartDate
	result.EndDate.String = strEndDate

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) toEmployeeReimbursementResponse(data []interface{}) (result []out.EmployeeReimbursement) {
	for _, item := range data {
		model, _ := item.(repository.EmployeeReimbursement)
		separatedPaths := strings.Split(model.Path.String, "/")
		fileName := separatedPaths[len(separatedPaths) - 1]

		res := out.EmployeeReimbursement{
			ID:                  model.ID.Int64,
			EmployeeId:          model.EmployeeId.Int64,
			IDCard:              model.IDCard.String,
			FirstName:           model.Firstname.String,
			LastName:            model.Lastname.String,
			FullName:            model.Firstname.String + " " + model.Lastname.String,
			Department:          model.Department.String,
			CurrentMedicalValue: model.CurrentMedicalValue.Float64,
			ReceiptNo:           model.ReceiptNo.String,
			Value:               model.Value.Float64,
			Status:              model.Status.String,
			VerifiedStatus:      model.VerifiedStatus.String,
			ApprovedValue:       model.ApprovedValue.Float64,
			Note:                model.Note.String,
			Filename: 			 fileName,
			Attachment:          model.Host.String + model.Path.String,
			CreatedAt:           model.CreatedAt.Time,
			UpdatedAt:           model.UpdatedAt.Time,
		}

		result = append(result, res)
	}

	return
}

func (input employeeService) responseReimbursementReport(data []interface{}) (result []out.EmployeeReimbursementReportResponse) {
	for _, item := range data {
		model, _ := item.(repository.EmployeeReimbursement)
		res := out.EmployeeReimbursementReportResponse{
			ID:                  model.ID.Int64,
			FirstName:           model.Firstname.String,
			LastName:            model.Lastname.String,
			FullName:            model.Firstname.String + " " + model.Lastname.String,
			CurrentMedicalValue: model.CurrentMedicalValue.Float64,
			LastMedicalValue:    model.LastMedicalValue.Float64,
			TotalValue:          model.MonthlyReport.Total.Float64,
			January:             model.MonthlyReportArr[0],
			February:            model.MonthlyReportArr[1],
			March:               model.MonthlyReportArr[2],
			April:               model.MonthlyReportArr[3],
			May:                 model.MonthlyReportArr[4],
			June:                model.MonthlyReportArr[5],
			July:                model.MonthlyReportArr[6],
			August:              model.MonthlyReportArr[7],
			September:           model.MonthlyReportArr[8],
			October:             model.MonthlyReportArr[9],
			November:            model.MonthlyReportArr[10],
			December:            model.MonthlyReportArr[11],
		}

		result = append(result, res)
	}

	return
}
