package ModuleEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ModuleService"
)

type moduleEndpoint struct {
	endpoint.AbstractEndpoint
}

var ModuleEndpoint = moduleEndpoint{}.New()

func (input moduleEndpoint) New() (output moduleEndpoint) {
	output.FileName = "ModuleEndpoint.go"
	return
}

func (input moduleEndpoint) getMenuCodeModule() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterSetupModuleRedesign, constanta.MenuUserMasterSetupModule)
}

func (input moduleEndpoint) ModuleEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ModuleEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeModule()+common.ViewDataPermissionMustHave, response, request, ModuleService.ModuleService.GetListModule)
		break
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeModule()+common.InsertDataPermissionMustHave, response, request, ModuleService.ModuleService.InsertModule)
		break
	}
}

func (input moduleEndpoint) ModuleEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ModuleEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeModule()+common.ViewDataPermissionMustHave, response, request, ModuleService.ModuleService.ViewModule)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeModule()+common.UpdateDataPermissionMustHave, response, request, ModuleService.ModuleService.UpdateModule)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeModule()+common.DeleteDataPermissionMustHave, response, request, ModuleService.ModuleService.DeleteModule)
		break
	}
}

func (input moduleEndpoint) InitiateModuleEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateModuleEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeModule()+common.ViewDataPermissionMustHave, response, request, ModuleService.ModuleService.InitiateModule)
		break
	}
}
