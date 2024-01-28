package EmployeeHistoryEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeService"
)

type employeeHistoryEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeHistoryEndpoint = employeeHistoryEndpoint{}.New()

func (input employeeHistoryEndpoint) New() (output employeeHistoryEndpoint) {
	output.FileName = "EmployeeHistoryEndpoint.go"
	return
}

func (input employeeHistoryEndpoint) EmployeeHistoryWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeHistoryWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.GetListEmployeeHistory)
	}
}

func (input employeeHistoryEndpoint) InitiateEmployeeHistoryEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeHistoryEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeService.EmployeeService.InitiateGetListEmployeeHistory)
	}
}
