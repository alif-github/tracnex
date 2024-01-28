package NexmileParameter

import (
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/NexmileParameterService"
)

type nexmileParameterEndpoint struct {
	endpoint.AbstractEndpoint
}

var NexmileParametersEndpoint = nexmileParameterEndpoint{}.New()

func (input nexmileParameterEndpoint) New() (output nexmileParameterEndpoint) {
	output.FileName = "NexmileParameterEndpoint.go"
	return
}

func (input nexmileParameterEndpoint) GetNexmileParameterWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "GetNexmileParameterWithoutParam"
	switch request.Method {
	case "POST":
		fmt.Println("GetNexmileParameterWithoutParam, hit ServeJWTTokenValidationEndpoint")
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexmileParameter+common.ViewDataPermissionMustHave, response, request, NexmileParameterService.NexmileParameterService.ViewNexmileParameter)
		//input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexmileParameter + common.ViewDataPermissionMustHave, response, request, NexmileParameterService.NexmileParameterService.ViewNexmileParameter)
		break
	}
}

func (input nexmileParameterEndpoint) AddNexmileParameterWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "AddNexmileParameterWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexmileParameter+common.InsertDataPermissionMustHave, response, request, NexmileParameterService.NexmileParameterService.InsertNexmileParameter)
		break
	}
}
