package EmployeeGradeEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeGradeService"
)

type employeeGradeEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeGradeEndpoint = employeeGradeEndpoint{}.New()

func (input employeeGradeEndpoint) New() (output employeeGradeEndpoint) {
	output.FileName = "EmployeeGradeEndpoint.go"
	return
}

func (input employeeGradeEndpoint) EmployeeGradeEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeGradeEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeGradeService.EmployeeGradeGetListService.GetEmployeeGrade)
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeGradeService.EmployeeGradeInsertService.InsertEmployeeGrade)

	}
}

func (input employeeGradeEndpoint) EmployeeGradeMatrixEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeGradeMatrixEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeGradeService.EmployeeGradeGetListService.GetEmployeeGradeMatrix)
	}
}

func (input employeeGradeEndpoint) EmployeeGradeEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeGradeEndpointWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeGradeService.EmployeeGradeDeleteService.DeleteEmployeeGrade)
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeGradeService.EmployeeGradeUpdateService.UpdateEmployeeGrade)

	}
}

func (input employeeGradeEndpoint) InitiateEmployeeGradeEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeGradeEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeGradeService.EmployeeGradeGetListService.InitiateEmployeeGrade)
	}
}

func (input employeeGradeEndpoint) InitiateEmployeeGradeMatrixEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeGradeMatrixEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeGradeService.EmployeeGradeGetListService.InitiateEmployeeGradeMatrix)
	}
}
