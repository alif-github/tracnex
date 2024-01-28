package ParameterEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ParameterService"
)

type parameterEndpoint struct {
	endpoint.AbstractEndpoint
}

var ParameterEndpoint = parameterEndpoint{}.New()

func (input parameterEndpoint) New() (output parameterEndpoint) {
	output.FileName = "ParameterEndpoint.go"
	return
}

func (input parameterEndpoint) ParameterEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ParameterEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, ParameterService.EmployeeParameterViewService.ViewParameterEmployee)
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ParameterService.EmployeeParameterUpdateService.UpdateEmployeeParameter)

	}
}

