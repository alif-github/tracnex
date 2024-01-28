package ReportEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/DepartmentService"
	"nexsoft.co.id/nextrac2/service/EmployeeService"
	"nexsoft.co.id/nextrac2/service/ReportService"
)

type reportEndpoint struct {
	endpoint.AbstractEndpoint
}

var ReportEndpoint = reportEndpoint{}.New()

func (input reportEndpoint) New() (output reportEndpoint) {
	output.FileName = "ReportEndpoint.go"
	return
}

func (input reportEndpoint) getMenuCodeReport() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterTimesheetReportRedesign, constanta.MenuUserMasterTimesheetReport)
}

func (input reportEndpoint) ReportEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ReportEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, ReportService.ReportService.GetListReport)
		break
	}
}

func (input reportEndpoint) InitiateReportEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateReportEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, ReportService.ReportService.InitiateReport)
		break
	}
}

func (input reportEndpoint) DownloadReportEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DownloadReportEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithValidationEndpointWithFileCSVResult(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, ReportService.ReportService.DownloadListReport)
		break
	}
}

func (input reportEndpoint) PaidPaymentReportEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "PaidPaymentReportEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, ReportService.ReportService.UpdatePaymentReport)
		break
	}
}

func (input reportEndpoint) PaidPaymentReportEndpointTry(response http.ResponseWriter, request *http.Request) {
	funcName := "PaidPaymentReportEndpointTry"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, ReportService.ReportService.CheckUpdateSprint)
		break
	}
}

func (input reportEndpoint) HelperSetDefaultRedmineSprintPaid(response http.ResponseWriter, request *http.Request) {
	funcName := "HelperSetDefaultRedmineSprintPaid"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, ReportService.ReportService.SetInitPaidSprint)
		break
	}
}

func (input reportEndpoint) HelperGetListPaymentHistory(response http.ResponseWriter, request *http.Request) {
	funcName := "HelperGetListPaymentHistory"
	switch request.Method {
	case "GET":
		input.ServeWhiteListEndpoint(funcName, false, response, request, ReportService.ReportService.GetListReportHistory)
		break
	}
}

func (input reportEndpoint) HelperViewDetailPaymentHistory(response http.ResponseWriter, request *http.Request) {
	funcName := "HelperViewDetailPaymentHistory"
	switch request.Method {
	case "GET":
		input.ServeWhiteListEndpoint(funcName, false, response, request, ReportService.ReportService.ViewReportHistory)
		break
	}
}

func (input reportEndpoint) DropDownListEmployeeEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DropDownListEmployeeEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployeeDDL)
		break
	}
}

func (input reportEndpoint) DropDownListDepartmentEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DropDownListDepartmentEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeReport()+common.ViewDataPermissionMustHave, response, request, DepartmentService.DepartmentService.GetListDepartment)
		break
	}
}
