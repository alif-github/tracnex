package DashboardEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/AbsentService"
	"nexsoft.co.id/nextrac2/service/DashboardService"
	"nexsoft.co.id/nextrac2/service/TodaysLeaveService"
)

type dashboardEndpoint struct {
	endpoint.AbstractEndpoint
}

var DashboardEndpoint = dashboardEndpoint{}.New()

func (input dashboardEndpoint) New() (output dashboardEndpoint) {
	output.FileName = "DashboardEndpoint.go"
	return
}

func (input dashboardEndpoint) ReimbursementDashboardPanelWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ReimbursementDashboardPanelWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, DashboardService.ReimbursementService.ViewDashboardCountReimbursement)
	}
}

func (input dashboardEndpoint) LeaveDashboardPanelWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LeaveDashboardPanelWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, DashboardService.LeaveDashboardService.ViewDashboardCountLeave)
	}
}

func (input dashboardEndpoint) AbsentDashboardPanelWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "AbsentDashboardPanelWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, DashboardService.AbsentDashboardService.ViewAbsentAverage)
	}
}

func (input dashboardEndpoint) AbsentDashboardWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "AbsentDashboardWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, AbsentService.AbsentService.UploadFileAbsent)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, AbsentService.AbsentService.GetListAbsent)
	}
}

func (input dashboardEndpoint) InitiateAbsentDashboardEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateAbsentDashboardEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, AbsentService.AbsentService.InitiateAbsent)
	}
}

func (input dashboardEndpoint) GetListAbsentPeriodDashboardEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "GetListAbsentPeriodDashboardEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, AbsentService.AbsentService.GetListAbsentPeriod)
	}
}

func (input dashboardEndpoint) LeaveDashboardWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LeaveDashboardWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, TodaysLeaveService.TodaysLeave.GetListTodaysLeave)
	}
}

func (input dashboardEndpoint) InitiateTodayLeaveEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateTodayLeaveEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, TodaysLeaveService.TodaysLeave.InitiateTodaysLeave)
	}
}
