package PKCEClientMappingEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/PKCEClientMappingService"
)

func (input pkceClientMappingEndpoints) ChangeNamePCKEClientMappingEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ChangeNamePCKEClientMappingEndpoint"
	input.FileName = "UIPKCEClientMappingEndpoint.go"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodePKCEClientMapping()+common.UpdateDataPermissionMustHave, response, request, PKCEClientMappingService.UpdatePKCEClientNameService.UpdateClientName)
}

func (input pkceClientMappingEndpoints) UIPKCEClientMappingWithoutPathParamEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIPKCEClientMappingWithoutPathParamEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodePKCEClientMapping()+common.ViewDataPermissionMustHave, response, request, PKCEClientMappingService.PKCEClientMappingService.GetListPKCEClientMapping)
		break
	}
}

func (input pkceClientMappingEndpoints) UIPKCEClientMappingWithPathParamEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIPKCEClientMappingWithPathParamEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodePKCEClientMapping()+common.ViewDataPermissionMustHave, response, request, PKCEClientMappingService.PKCEClientMappingService.ViewPKCEClientMapping)
		break
	}
}

func (input pkceClientMappingEndpoints) UIGetListPKCEClientMappingInitiateEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIGetListPKCEClientMappingInitiateEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodePKCEClientMapping()+common.ViewDataPermissionMustHave, response, request, PKCEClientMappingService.PKCEClientMappingService.InitiateGetListPKCEClientMapping)
}
