package ReportAnnualLeaveEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ReportAnnualLeaveService"
)

type reportAnnualLeaveEndpoint struct {
	endpoint.AbstractEndpoint
}

var ReportAnnualLeaveEndpoint = reportAnnualLeaveEndpoint{}.New()

func (input reportAnnualLeaveEndpoint) New() (output reportAnnualLeaveEndpoint) {
	output.FileName = "ReportAnnualLeaveEndpoint.go"
	return
}

func (input reportAnnualLeaveEndpoint) ReportAnnualLeaveEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ReportAnnualLeaveEndpoint"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, ReportAnnualLeaveService.ReportAnnualLeaveService.ReportAnnualLeaveService)
	}
}

func (input reportAnnualLeaveEndpoint) ReportAnnualLeaveJobListEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ReportAnnualLeaveJobListEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, ReportAnnualLeaveService.ReportAnnualLeaveService.GetListJobReportAnnualLeaveService)
	}
}
