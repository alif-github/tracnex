package ClientTypeEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ClientTypeService"
)

type clientTypeEndpoint struct {
	endpoint.AbstractEndpoint
}

var ClientTypeEndpoint = clientTypeEndpoint{}.New()

func (input clientTypeEndpoint) New() (output clientTypeEndpoint) {
	output.FileName = "ClientTypeEndpoint.go"
	return
}

func (input clientTypeEndpoint) ClientTypeEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ClientTypeEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, ClientTypeService.ClientTypeService.GetListClientType)
		break
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ClientTypeService.ClientTypeService.InsertClientType)
		break
	}
}

func (input clientTypeEndpoint) InitiateClientTypeEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ClientTypeEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, ClientTypeService.ClientTypeService.InitiateClientType)
		break
	}
}

func (input clientTypeEndpoint) ClientTypeEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ClientTypeEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, ClientTypeService.ClientTypeService.ViewClientType)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ClientTypeService.ClientTypeService.UpdateClientType)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ClientTypeService.ClientTypeService.DeleteClientType)
		break
	}
}

func (input clientTypeEndpoint) ClientTypeAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ClientTypeAdminEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftClientType+common.ViewDataPermissionMustHave, response, request, ClientTypeService.ClientTypeService.GetListClientTypeAdmin)
		break
	}
}
