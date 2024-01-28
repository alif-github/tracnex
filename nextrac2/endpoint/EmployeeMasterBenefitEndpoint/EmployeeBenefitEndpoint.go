package EmployeeMasterBenefitEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeMasterBenefitService"
)

type employeeBenefitEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeBenefitEndpoint = employeeBenefitEndpoint{}.New()

func (input employeeBenefitEndpoint) New() (output employeeBenefitEndpoint) {
	output.FileName = "EmployeeBenefitEndpoint.go"
	return
}

func (input employeeBenefitEndpoint) EmployeeBenefitEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeAllowanceEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeMasterBenefitService.EmployeeBenefitGetListService.GetEmployeeBenefit)
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeMasterBenefitService.EmployeeBenefitInsertService.InsertEmployeeBenefit)

	}
}

func (input employeeBenefitEndpoint) EmployeeBenefitEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeAllowanceEndpointWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeMasterBenefitService.EmployeeBenefitDeleteService.DeleteEmployeeBenefit)
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EmployeeMasterBenefitService.EmployeeBenefitUpdateService.UpdateEmployeeBenefit)

	}
}

func (input employeeBenefitEndpoint) InitiateEmployeeBenefitEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeBenefitEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, EmployeeMasterBenefitService.EmployeeBenefitGetListService.InitiateEmployeeBenefit)
	}
}