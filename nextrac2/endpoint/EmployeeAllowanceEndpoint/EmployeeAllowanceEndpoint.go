package EmployeeAllowanceEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeAllowanceService"
)

type employeeAllowanceEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeAllowanceEndpoint = employeeAllowanceEndpoint{}.New()

func (input employeeAllowanceEndpoint) New() (output employeeAllowanceEndpoint) {
	output.FileName = "EmployeeAllowanceEndpoint.go"
	return
}

func (input employeeAllowanceEndpoint) EmployeeAllowanceEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeAllowanceEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeAllowanceService.EmployeeAllowanceGetListService.GetEmployeeAllowance)
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeAllowanceService.EmployeeAllowanceInsertService.InsertEmployeeAllowance)

	}
}

func (input employeeAllowanceEndpoint) EmployeeAllowanceEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeAllowanceEndpointWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeAllowanceService.EmployeeAllowanceDeleteService.DeleteEmployeeAllowance)
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeAllowanceService.EmployeeAllowanceUpdateService.UpdateEmployeeALlowance)

	}
}

func (input employeeAllowanceEndpoint) InitiateEmployeeAllowanceEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeAllowanceEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeAllowanceService.EmployeeAllowanceGetListService.InitiateEmployeeAllowance)
	}
}