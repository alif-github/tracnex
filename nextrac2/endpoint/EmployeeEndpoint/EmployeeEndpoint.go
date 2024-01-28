package EmployeeEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeService"
)

type employeeEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeEndpoint = employeeEndpoint{}.New()

func (input employeeEndpoint) New() (output employeeEndpoint) {
	output.FileName = "EmployeeEndpoint.go"
	return
}

func (input employeeEndpoint) getMenuCodeEmployee() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterTimesheetEmployeeRedesign, constanta.MenuUserMasterTimesheetEmployee)
}

//--- Start Of Employee

func (input employeeEndpoint) EmployeeEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.InsertDataPermissionMustHave, response, request, EmployeeService.EmployeeService.InsertEmployee)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployee)
		break
	}
}

func (input employeeEndpoint) EmployeeEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeEndpointWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.DeleteDataPermissionMustHave, response, request, EmployeeService.EmployeeService.DeleteEmployee)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.UpdateDataPermissionMustHave, response, request, EmployeeService.EmployeeService.UpdateEmployee)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.ViewEmployee)
		break
	}
}

func (input employeeEndpoint) EmployeeAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeAdminEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftDistrict+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployeeByAdmin)
	}
}

func (input employeeEndpoint) InitiateEmployeeEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.InitiateEmployee)
		break
	}
}

//--- End Of Employee

//--- Start of Employee Timesheet

func (input employeeEndpoint) EmployeeTimeSheetEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeTimeSheetEndpointWithParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.UpdateDataPermissionMustHave, response, request, EmployeeService.EmployeeService.UpdateEmployeeTimeSheet)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.ViewEmployeeTimeSheet)
		break
	}
}

func (input employeeEndpoint) EmployeeTimeSheetEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeTimesheetEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployeeTimeSheet)
		break
	}
}

func (input employeeEndpoint) InitiateEmployeeTimeSheetEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeTimeSheetEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.InitiateEmployeeTimeSheet)
		break
	}
}

func (input employeeEndpoint) EmployeeTimeSheetCheckEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetEmployeeTimeSheetRedmineByNIK)
		break
	}
}

//--- End Of Employee Timesheet

//func (input employeeEndpoint) EmployeeAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
//	funcName := "EmployeeAdminEndpointWithoutParam"
//	switch request.Method {
//	case "GET":
//		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftDistrict+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployeeTimeSheetByAdmin)
//	}
//}

//--- Start of Employee Leave

func (input employeeEndpoint) EmployeeLeaveEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeLeaveEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeLeave)
		break
	case "POST":
		//input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.InsertDataPermissionMustHave, response, request, EmployeeService.EmployeeService.InsertEmployeeLeave)
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InsertEmployeeLeave)
		break
	}
}

func (input employeeEndpoint) EmployeeLeaveEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeLeaveEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeLeaveDetailService.DetailEmployeeLeave)
		break
	}
}

func (input employeeEndpoint) InitiateGetListEmployeeLeave(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetListEmployeeLeave"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetListEmployeeLeave)
}

func (input employeeEndpoint) GetRemainingLeave(response http.ResponseWriter, request *http.Request) {
	funcName := "GetRemainingLeave"
	//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.InsertDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetRemainingLeave)
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetRemainingLeave)
}

func (input employeeEndpoint) DownloadEmployeeLeaveReport(response http.ResponseWriter, request *http.Request) {
	funcName := "DownloadEmployeeLeaveReport"
	//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.InsertDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetRemainingLeave)
	input.ServeJWTTokenValidationEndpointWithFileResult(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.DownloadEmployeeLeaveReport)
}

func (input employeeEndpoint) GetEmployeeLeaveTypes(response http.ResponseWriter, request *http.Request) {
	funcName := "GetEmployeeLeaveTypes"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetEmployeeLeaveTypes)
}

func (input employeeEndpoint) InitiateGetEmployeeLeaveTypes(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetEmployeeLeaveTypes"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetEmployeeLeaveTypes)
}

//--- End of Employee Leave

//--- Start of Employee Reimbursement

func (input employeeEndpoint) EmployeeReimbursementEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeReimbursementEndpointWithoutParam"
	switch request.Method {
	case "GET" :
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeReimbursement)
		break
	case "POST":
		//input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.InsertDataPermissionMustHave, response, request, EmployeeService.EmployeeService.InsertEmployeeReimbursement)
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InsertEmployeeReimbursement)
		break
	}
}

func (input employeeEndpoint) InitiateGetListEmployeeReimbursement(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetListEmployeeReimbursement"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetListEmployeeReimbursement)
}

func (input employeeEndpoint) GetMedicalRemainingBalance(response http.ResponseWriter, request *http.Request) {
	funcName := "GetMedicalRemainingBalance"
	//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.InsertDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetMedicalRemainingBalance)
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetMedicalRemainingBalance)
}

func (input employeeEndpoint) GetEmployeeReimbursementTypes(response http.ResponseWriter, request *http.Request) {
	funcName := "GetEmployeeReimbursementTypes"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetEmployeeReimbursementTypes)
}

func (input employeeEndpoint) InitiateGetEmployeeReimbursementTypes(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetEmployeeReimbursementTypes"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetEmployeeReimbursementTypes)
}

//--- End of Employee Reimbursement

//--- Start of Employee History

func (input employeeEndpoint) EmployeeRequestHistoryEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeRequestHistoryEndpointWithoutParam"
	switch request.Method {
	case "GET":
		//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployeeRequestHistory)
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeRequestHistory)
		break
	}
}

func (input employeeEndpoint) InitiateEmployeeRequestHistoryEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeRequestHistoryEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetListRequestHistory)
}

func (input employeeEndpoint) CancelEmployeeRequestEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "CancelEmployeeRequestEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.CancelEmployeeRequest)
}

func (input employeeEndpoint) EmployeeApprovalHistoryEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeApprovalHistoryEndpointWithoutParam"
	switch request.Method {
	case "GET":
		//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployeeApprovalHistory)
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeApprovalHistory)
		break
	}
}

func (input employeeEndpoint) InitiateEmployeeApprovalHistoryEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeApprovalHistoryEndpoint"
	//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.InitiateGetListApprovalHistory)
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetListApprovalHistory)
}

func (input employeeEndpoint) UpdateStatusEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UpdateStatusEndpointWithoutParam"
	//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.UpdateApprovalStatus)
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.UpdateApprovalStatus)
}

func (input employeeEndpoint) GetListEmployeeLeaveYearlyEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "GetListEmployeeLeaveYearlyEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeLeaveYearly)
}

func (input employeeEndpoint) IntiateEmployeeLeaveYearlyEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "IntiateEmployeeLeaveYearlyEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetListEmployeeLeaveYearly)
}

func (input employeeEndpoint) VerifyReimbursementEnpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "VerifyReimbursementEnpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeReimbursementVerifyService.VerifyReimbursement)
}

func (input employeeEndpoint) DownloadAnnualLeaveReport(response http.ResponseWriter, request *http.Request) {
	funcName := "DownloadEmployeeLeaveReport"
	input.ServeJWTTokenValidationEndpointWithFileResult(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.DownloadAnnualReport)
}

//--- Start of Employee History

func (input employeeEndpoint) EmployeeNotificationEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeNotificationEndpointWithoutParam"

	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeNotification)
		break
	}
}

func (input employeeEndpoint) ReadNotificationEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ReadNotificationEndpoint"

	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.ReadEmployeeNotification)
}

func (input employeeEndpoint) DownloadReimbursementReport(response http.ResponseWriter, request *http.Request) {
	funcName := "DownloadReimbursementReport"
	input.ServeJWTTokenValidationEndpointWithFileResult(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.DownloadReimbursementReport)
}

func (input employeeEndpoint) GetListEmployeeReimbursementReport(response http.ResponseWriter, request *http.Request) {
	funcName := "GetListEmployeeReimbursementReport"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeReimbursementReport)
}

func (input employeeEndpoint) InitiateGetListEmployeeReimbursementReport(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetListEmployeeReimbursementReport"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetListEmployeeReimbursementReport)
}