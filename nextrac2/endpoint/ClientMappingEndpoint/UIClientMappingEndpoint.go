package ClientMappingEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ClientMappingService"
)

func (input clientMappingEndpoint) UIClientMappingChangeNameND6Endpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIClientMappingChangeNameND6Endpoint"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustomerClientMapping()+common.UpdateDataPermissionMustHave, response, request, ClientMappingService.UIUpdateClientNameService.UpdateClientName)
		break
	}
}

func (input clientMappingEndpoint) UIClientMappingWithoutPathParamEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIClientMappingWithoutPathParamEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustomerClientMapping()+common.ViewDataPermissionMustHave, response, request, ClientMappingService.UIGetListClientMappingService.GetClientMappings)
		break
	}
}

func (input clientMappingEndpoint) UIClientMappingWithPathParamEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIClientMappingWithPathParamEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustomerClientMapping()+common.ViewDataPermissionMustHave, response, request, ClientMappingService.UIViewClientMappingService.ViewClientMapping)
		break
	}
}

func (input clientMappingEndpoint) UIGetListCLientMappingInitiateEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIGetListCLientMappingInitiateEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustomerClientMapping()+common.ViewDataPermissionMustHave, response, request, ClientMappingService.UIGetListClientMappingService.InitiateGetListClientMappings)
}
