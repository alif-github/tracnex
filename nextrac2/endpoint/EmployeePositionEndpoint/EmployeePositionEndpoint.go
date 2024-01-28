package EmployeePositionEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeePositionService"
)

type employeePositionEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeePositionEndpoint = employeePositionEndpoint{}.New()

func (input employeePositionEndpoint) New() (output employeePositionEndpoint) {
	output.FileName = "EmployeePositionEndpoint.go"
	return
}

func (input employeePositionEndpoint) EmployeePositionEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeePositionEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeePositionService.EmployeePositionService.InsertEmployeePosition)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeePositionService.EmployeePositionService.GetListPosition)
	}
}

func (input employeePositionEndpoint) EmployeePositionEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeePositionEndpointWithParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeePositionService.EmployeePositionService.UpdateEmployeePosition)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeePositionService.EmployeePositionService.ViewDetailEmployeePosition)
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeePositionService.EmployeePositionService.DeleteEmployeePosition)
	}
}

func (input employeePositionEndpoint) InitiateEmployeePositionEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeePositionEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeePositionService.EmployeePositionService.InitiatePosition)
	}
}
