package EmployeeFacilitiesActiveEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeFacilitiesActiveService"
)

type employeeMatrixEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeMatrixEndpoint = employeeMatrixEndpoint{}.New()

func (input employeeMatrixEndpoint) New() (output employeeMatrixEndpoint) {
	output.FileName = "EmployeeAllowanceEndpoint.go"
	return
}

func (input employeeMatrixEndpoint) EmployeeMatrixEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeMatrixEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeFacilitiesActiveService.EmployeeMatrixGetListService.GetEmployeeMatrix)
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeFacilitiesActiveService.EmployeeMatrixInsertService.InsertEmployeeMatrix)
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeFacilitiesActiveService.EmployeeMatrixDeleteService.DeleteEmployeeMatrix)
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeFacilitiesActiveService.EmployeeMatrixUpdateService.UpdateEmployeeMatrix)
	}
}

func (input employeeMatrixEndpoint) EmployeeMatrixEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeMatrixEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeFacilitiesActiveService.EmployeeMatrixDetailService.DetailEmployeeMatrix)
	case "PUT":
		//input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeAllowanceService.EmployeeAllowanceUpdateService.UpdateEmployeeALlowance)

	}
}

func (input employeeMatrixEndpoint) InitiateEmployeeMatrixEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeMatrixEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeFacilitiesActiveService.EmployeeMatrixGetListService.InitiateEmployeeMatrix)
	}
}
