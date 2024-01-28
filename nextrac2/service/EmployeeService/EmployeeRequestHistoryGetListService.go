package EmployeeService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

func (input employeeService) InitiateGetListRequestHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		funcName = "InitiateGetListRequestHistory"
		searchByParam []in.SearchByParam
		countData     interface{}
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeRequestHistoryValidOperator)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	if !employee.ID.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee")
		return
	}

	countData, errModel = dao.EmployeeHistoryDAO.InitiateEmployeeRequestHistory(serverconfig.ServerAttribute.DBConnection, searchByParam, employee.ID.Int64)
	if errModel.Error != nil {
		return
	}

	if countData == nil {
		countData = 0
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  []string{"created_at"},
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeRequestHistoryValidOperator,
		CountData:     countData.(int),
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) GetListEmployeeRequestHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		funcName = "GetListEmployeeRequestHistory"
		inputStruct in.GetListDataDTO
		searchByParam []in.SearchByParam
		validOrderBy = []string{"created_at"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, input.ValidSearchBy, validOrderBy, applicationModel.GetListEmployeeRequestHistoryValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	userId, _ := service.CheckIsOnlyHaveOwnPermission(*contextModel)

	employee, errModel := dao.EmployeeDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	if !employee.ID.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee")
		return
	}

	employeeRequestHistoryList, errModel := dao.EmployeeHistoryDAO.GetListEmployeeRequestHistory(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, userId, employee.ID.Int64)
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.toEmployeeRequestHistoryResponse(employeeRequestHistoryList)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)

	return
}

func (input employeeService) toEmployeeRequestHistoryResponse(data []interface{}) (result []out.EmployeeRequestHistoryResponse) {
	for _, item := range data {
		employeeRequestHistory, _ := item.(repository.EmployeeRequestHistory)
		date := input.getRequestHistoryDate(employeeRequestHistory)

		result = append(result, out.EmployeeRequestHistoryResponse{
			ID:                 employeeRequestHistory.ID.Int64,
			ReceiptNo: 			employeeRequestHistory.ReceiptNo.String,
			Description:        employeeRequestHistory.Description.String,
			Date: 				date,
			RequestType:        employeeRequestHistory.RequestType.String,
			Type: 				employeeRequestHistory.Type.String,
			TotalLeave:         employeeRequestHistory.TotalLeave.Int64,
			Value:              employeeRequestHistory.Value.Float64,
			ApprovedValue:      employeeRequestHistory.ApprovedValue.Float64,
			Status:             employeeRequestHistory.Status.String,
			VerifiedStatus:     employeeRequestHistory.VerifiedStatus.String,
			CancellationReason: employeeRequestHistory.CancellationReason.String,
			Note:               employeeRequestHistory.Note.String,
			Attachment:         employeeRequestHistory.Host.String + employeeRequestHistory.Path.String,
			CreatedAt:          employeeRequestHistory.CreatedAt.Time,
			UpdatedAt:          employeeRequestHistory.UpdatedAt.Time,
		})
	}

	return
}

func (input employeeService) getRequestHistoryDate(employeeRequestHistory repository.EmployeeRequestHistory) (result []string) {
	if !employeeRequestHistory.Date.Time.IsZero() {
		result = append(result, employeeRequestHistory.Date.Time.Format("2006-01-02"))
		return
	}

	_ = json.Unmarshal([]byte(employeeRequestHistory.LeaveDate.String), &result)
	return
}