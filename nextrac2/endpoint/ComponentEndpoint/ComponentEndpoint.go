package ComponentEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ComponentService"
)

type componentEndpoint struct {
	endpoint.AbstractEndpoint
}

var ComponentEndpoint = componentEndpoint{}

func (input componentEndpoint) New() (output componentEndpoint) {
	output.FileName = "ComponentEndpoint.go"
	return
}

func (input componentEndpoint) getMenuCodeComponent() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterSetupComponentRedesign, constanta.MenuUserMasterSetupComponent)
}

func (input componentEndpoint) ComponentEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ComponentEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeComponent()+common.ViewDataPermissionMustHave, response, request, ComponentService.ComponentService.GetListComponent)
		break
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeComponent()+common.InsertDataPermissionMustHave, response, request, ComponentService.ComponentService.InsertComponent)
		break
	}
}

func (input componentEndpoint) ComponentEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ComponentEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeComponent()+common.ViewDataPermissionMustHave, response, request, ComponentService.ComponentService.ViewComponent)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeComponent()+common.UpdateDataPermissionMustHave, response, request, ComponentService.ComponentService.UpdateComponent)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeComponent()+common.DeleteDataPermissionMustHave, response, request, ComponentService.ComponentService.DeleteComponent)
		break
	}
}

func (input componentEndpoint) InitiateComponentEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateComponentEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeComponent()+common.ViewDataPermissionMustHave, response, request, ComponentService.ComponentService.InitiateComponent)
		break
	}
}
