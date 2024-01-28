package EmployeeService

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
	"nexsoft.co.id/nextrac2/service"
	"strings"
)

func (input employeeService) InitiateGetListApprovalHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
		validSearchBy = []string{"status"}
		validOrderBy  = []string{"created_at"}
		statusParam	  = input.getStatusQueryParam(request)
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListEmployeeApprovalHistoryValidOperator)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	memberIdList := in.MemberList{}
	_ = json.Unmarshal([]byte(employee.Member.String), &memberIdList)

	countData, errModel = dao.EmployeeHistoryDAO.InitiateEmployeeApprovalHistoryByEmployeeIdList(serverconfig.ServerAttribute.DBConnection, searchByParam, memberIdList.MemberID, statusParam)
	if errModel.Error != nil {
		return
	}

	if countData == nil {
		countData = 0
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: validSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeApprovalHistoryValidOperator,
		CountData:     countData.(int),
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) GetListEmployeeApprovalHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"status"}
		validOrderBy  = []string{"created_at"}
		statusParam	  = input.getStatusQueryParam(request)
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeApprovalHistoryValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	memberIdList := in.MemberList{}
	_ = json.Unmarshal([]byte(employee.Member.String), &memberIdList)

	employeeApprovalHistoryList, errModel := dao.EmployeeHistoryDAO.GetListEmployeeApprovalHistoryByEmployeeIdListAndStatus(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, memberIdList.MemberID, statusParam)
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.toEmployeeApprovalHistoryResponse(employeeApprovalHistoryList)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) getStatusQueryParam(request *http.Request) (result []string) {
	statusParam := request.URL.Query().Get("status")

	if statusParam == "" {
		return nil
	}

	separatedStatuses := strings.Split(statusParam, ",")

	for _, status := range separatedStatuses {
		result = append(result, strings.Trim(status, " "))
	}

	return
}

func (input employeeService) toEmployeeApprovalHistoryResponse(data []interface{}) (result []out.EmployeeApprovalHistory) {
	for _, item := range data {
		employeeApprovalHistory, _ := item.(repository.EmployeeApprovalHistory)
		date := input.getApprovalHistoryDate(employeeApprovalHistory)

		if employeeApprovalHistory.RequestType.String == "" {
			employeeApprovalHistory.RequestType.String = constanta.ReimbursementType
		}

		res := out.EmployeeApprovalHistory{
			ID:                  employeeApprovalHistory.ID.Int64,
			Firstname:           employeeApprovalHistory.Firstname.String,
			Lastname:            employeeApprovalHistory.Lastname.String,
			IDCard:              employeeApprovalHistory.IDCard.String,
			Department:          employeeApprovalHistory.Department.String,
			RequestType:         employeeApprovalHistory.RequestType.String,
			Type:                employeeApprovalHistory.Type.String,
			ReceiptNo: 			 employeeApprovalHistory.ReceiptNo.String,
			TotalLeave:          employeeApprovalHistory.TotalLeave.Int64,
			Value: 				 employeeApprovalHistory.Value.Float64,
			Date:                date,
			Status:              employeeApprovalHistory.Status.String,
			VerifiedStatus:      employeeApprovalHistory.VerifiedStatus.String,
			ApprovedValue:       employeeApprovalHistory.ApprovedValue.Float64,
			CancellationReason:  employeeApprovalHistory.CancellationReason.String,
			Note:                employeeApprovalHistory.Note.String,
			Attachment:          employeeApprovalHistory.Host.String + employeeApprovalHistory.Path.String,
			TotalRemainingLeave: employeeApprovalHistory.TotalRemainingLeave.Int64,
			Description:         employeeApprovalHistory.Description.String,
			CreatedAt:           employeeApprovalHistory.CreatedAt.Time,
			UpdatedAt:           employeeApprovalHistory.UpdatedAt.Time,
		}

		result = append(result, res)
	}

	return
}

func (input employeeService) getApprovalHistoryDate(employeeApprovalHistory repository.EmployeeApprovalHistory) (result []string) {
	if !employeeApprovalHistory.Date.Time.IsZero() {
		result = append(result, employeeApprovalHistory.Date.Time.Format("2006-01-02"))
		return
	}

	_ = json.Unmarshal([]byte(employeeApprovalHistory.LeaveDate.String), &result)
	return
}
