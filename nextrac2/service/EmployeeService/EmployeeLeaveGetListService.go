package EmployeeService

import (
	"database/sql"
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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input employeeService) InitiateGetListEmployeeLeave(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name"}
		validOrderBy  = []string{"id", "created_at"}
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListEmployeeLeaveValidOperator)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	employeeLeaveFilter, errModel := input.getEmployeeLeaveFilter(request, employee)
	if errModel.Error != nil {
		return
	}

	result, errModel := dao.EmployeeLeaveDAO.InitiateGetListEmployeeLeave(serverconfig.ServerAttribute.DBConnection, searchByParam, employeeLeaveFilter)
	if errModel.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: validSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeLeaveValidOperator,
		CountData:     result,
		ValidSearchParam :[]out.SearchByParam{
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

func (input employeeService) GetListEmployeeLeave(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "el.type"}
		validOrderBy  = []string{"id", "created_at", "e.id_card"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeLeaveValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	employeeLeaveFilter, errModel := input.getEmployeeLeaveFilter(request, employee)
	if errModel.Error != nil {
		return
	}

	results, errModel := dao.EmployeeLeaveDAO.GetListEmployeeLeave(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, employeeLeaveFilter)
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.toEmployeeLeaveResponse(results)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) getEmployeeLeaveFilter(request *http.Request, employee repository.EmployeeModel) (result repository.EmployeeLeaveModel, errModel errorModel.ErrorModel) {
	var (
		leaveDate    string
		memberIdList in.MemberList
		funcName	 = "getEmployeeLeaveFilter"
		idCard       = request.URL.Query().Get("id_card")
		fullName     = request.URL.Query().Get("full_name")
		department   = request.URL.Query().Get("department")
		status		 = request.URL.Query().Get("status")
		onLeave      = request.URL.Query().Get("on_leave")
		leaveType    = request.URL.Query().Get("leave_type")
		key          = request.URL.Query().Get("key")
		keyword      = request.URL.Query().Get("keyword")
		strStartDate = request.URL.Query().Get("start_date")
		strEndDate   = request.URL.Query().Get("end_date")
	)

	/*
		Status
	*/
	if status != "" && !input.isStatusValid(status) {
		errModel = errorModel.GenerateFormatFieldError(input.FileName, funcName, "status")
		return
	}

	/*
		On Leave
	*/
	isOnLeave := onLeave == "true"
	if isOnLeave {
		/*
			Filter by member
		*/
		_ = json.Unmarshal([]byte(employee.Member.String), &memberIdList)

		/*
			Filter by leave date
		*/
		now, _ := time.Parse(constanta.DefaultTimeFormat, time.Now().Format(constanta.DefaultTimeFormat))
		leaveDate = now.Format("2006-01-02")
	}

	result = repository.EmployeeLeaveModel{
		MemberList: memberIdList.MemberID,
		IDCard: 	sql.NullString{String: idCard},
		Name:       sql.NullString{String: fullName},
		Department: sql.NullString{String: department},
		Status:		sql.NullString{String: status},
		OnLeave:    sql.NullBool{Bool: isOnLeave},
		LeaveDate:  sql.NullString{String: leaveDate},
		SearchBy:   sql.NullString{String: key},
		Keyword:    sql.NullString{String: keyword},
		IsYearly:   sql.NullBool{Bool: true},
		Type: 		sql.NullString{String: leaveType},
	}

	/*
		Date Range
	*/
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

	result.StrStartDate.String = strStartDate
	result.StrEndDate.String = strEndDate
	return
}

func (input employeeService) isStatusValid(status string) bool {
	return status == constanta.PendingRequestStatus ||
		status == constanta.PendingCancellationRequestStatus ||
		status == constanta.ApprovedRequestStatus ||
		status == constanta.RejectedRequestStatus ||
		status == constanta.CancelledRequestStatus
}

func (input employeeService) toEmployeeLeaveResponse(data []interface{}) (result []out.EmployeeLeave) {
	for _, item := range data {
		var date []string

		model, _ := item.(repository.EmployeeLeaveModel)

		_ = json.Unmarshal([]byte(model.StrDateList.String), &date)

		res := out.EmployeeLeave{
			ID:            model.ID.Int64,
			IDCard:        model.IDCard.String,
			FirstName:     model.Firstname.String,
			LastName:      model.Lastname.String,
			FullName:      model.Firstname.String + " " + model.Lastname.String,
			Department:    model.Department.String,
			AllowanceName: model.AllowanceName.String,
			Date:          date,
			Value:         model.Value.Int64,
			Type :         model.Type.String,
			Status: 	   model.Status.String,
			StartDate:     model.StartDate.Time,
			EndDate:       model.EndDate.Time,
			LeaveTime:     model.LeaveTime.Time,
		}

		result = append(result, res)
	}

	return
}
