package EmployeeContractEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EmployeeContractService"
)

type employeeContractEndpoint struct {
	endpoint.AbstractEndpoint
}

var EmployeeContractEndpoint = employeeContractEndpoint{}.New()

func (input employeeContractEndpoint) New() (output employeeContractEndpoint) {
	output.FileName = "EmployeeContractEndpoint.go"
	return
}

func (input employeeContractEndpoint) getMenuCodeEmployeeContract() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterTimesheetEmployeeRedesign, constanta.MenuUserMasterTimesheetEmployee)
}

func (input employeeContractEndpoint) EmployeeContractEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeContractEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployeeContract()+common.InsertDataPermissionMustHave, response, request, EmployeeContractService.EmployeeContractService.InsertEmployeeContract)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployeeContract()+common.ViewDataPermissionMustHave, response, request, EmployeeContractService.EmployeeContractService.GetListEmployeeContract)
		break
	}
}

func (input employeeContractEndpoint) EmployeeContractEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EmployeeContractEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployeeContract()+common.ViewDataPermissionMustHave, response, request, EmployeeContractService.EmployeeContractService.ViewEmployeeContract)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployeeContract()+common.UpdateDataPermissionMustHave, response, request, EmployeeContractService.EmployeeContractService.UpdateEmployeeContract)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployeeContract()+common.DeleteDataPermissionMustHave, response, request, EmployeeContractService.EmployeeContractService.DeleteEmployeeContract)
		break
	}
}

func (input employeeContractEndpoint) InitiateEmployeeContractEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateEmployeeContractEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployeeContract()+common.ViewDataPermissionMustHave, response, request, EmployeeContractService.EmployeeContractService.InitiateEmployeeContract)
		break
	}
}
