package EmployeeService

import (
	"encoding/json"
	"fmt"
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

func (input employeeService) GetListEmployeeNotification(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		db = serverconfig.ServerAttribute.DBConnection
		funcName = "GetListEmployeeNotification"
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validOrderBy  = []string{"created_at"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, []string{}, validOrderBy, applicationModel.GetListEmployeeNotificationValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(db, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	if !employee.ID.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee")
		return
	}

	memberIdList := in.MemberList{}
	_ = json.Unmarshal([]byte(employee.Member.String), &memberIdList)

	auditList, errModel := dao.AuditSystemDAO.GetListEmployeeNotification(db, inputStruct, searchByParam, repository.EmployeeNotification{
		EmployeeId:     employee.ID.Int64,
		MemberIdList:   memberIdList.MemberID,
	})
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.toEmployeeNotificationResponse(auditList)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) toEmployeeNotificationResponse(auditList []interface{}) (results []out.EmployeeNotification) {
	for _, data := range auditList {
		var employeeNotification in.EmployeeNotification

		item, _ := data.(repository.AuditSystemModel)
		_ = json.Unmarshal([]byte(item.Description.String), &employeeNotification)

		messageTitle, messageBody := input.generateNotificationMessage(item.Employee.Firstname, item.Employee.Lastname, employeeNotification)

		result := out.EmployeeNotification{
			ID:                          item.ID.Int64,
			IsRequestingForApproval:     employeeNotification.IsRequestingForApproval,
			IsRequestingForCancellation: employeeNotification.IsRequestingForCancellation,
			IsCancellation:              employeeNotification.IsCancellation,
			IsVerified: 				 employeeNotification.IsVerified,
			EmployeeId:                  employeeNotification.EmployeeId,
			Name:                        fmt.Sprintf("%s %s", item.Employee.Firstname, item.Employee.Lastname),
			RequestType:                 employeeNotification.RequestType,
			Status:                      employeeNotification.Status,
			Date:                        input.formatDateList(employeeNotification.Date),
			MessageTitle:                messageTitle,
			MessageBody:                 messageBody,
			IsRead: 					 employeeNotification.IsRead,
			CreatedAt:                   item.CreatedAt.Time,
		}

		results = append(results, result)
	}

	return
}

func (input employeeService) formatDateList(tl []time.Time) []time.Time {
	for i, t := range tl {
		tl[i], _ = time.Parse(constanta.DefaultTimeFormat, t.Format(constanta.DefaultTimeFormat))
	}

	return tl
}