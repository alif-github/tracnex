package EmployeeLevelEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeLevelService"
)

type employeeLevelEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeLevelEndpoint = employeeLevelEndpoint{}.New()

func (input employeeLevelEndpoint) New() (output employeeLevelEndpoint) {
	output.FileName = "EmployeeLevelEndpoint.go"
	return
}

func (input employeeLevelEndpoint) EmployeeLevelEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeLevelEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeLevelService.EmployeeLevelGetListService.GetEmployeeLevel)
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeLevelService.EmployeeLevelInsertService.InsertEmployeeLevel)

	}
}

func (input employeeLevelEndpoint) EmployeeLevelMatrixEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeLevelMatrixEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeLevelService.EmployeeLevelGetListService.GetEmployeeLevelMatrix)
	}
}

func (input employeeLevelEndpoint) EmployeeLevelEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeLevelEndpointWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeLevelService.EmployeeLevelDeleteService.DeleteEmplooyeeLevel)
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeLevelService.EmployeeLevelUpdateService.UpdateEmployeeLevel)

	}
}

func (input employeeLevelEndpoint) InitiateEmployeeLevelEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeLevelEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeLevelService.EmployeeLevelGetListService.InitiateEmployeeLevel)
	}
}

func (input employeeLevelEndpoint) InitiateEmployeeLevelMatrixEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeLevelMatrixEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeLevelService.EmployeeLevelGetListService.InitiateEmployeeLevelMatrix)
	}
}
