package CustomerInstallationEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/CustomerInstallationService"
)

type customerInstallationEndpoint struct {
	endpoint.AbstractEndpoint
}

var CustomerInstallationEndpoint = customerInstallationEndpoint{}.New()

func (input customerInstallationEndpoint) New() (output customerInstallationEndpoint) {
	output.FileName = "CustomerInstallationEndpoint.go"
	return
}

func (input customerInstallationEndpoint) getMenuCodeCustomerInstallation() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterKonsumenInstallationRedesign, constanta.MenuUserMasterKonsumenInstallation)
}

func (input customerInstallationEndpoint) CustomerInstallationEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerInstallationEndpointWithParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustomerInstallation()+common.UpdateDataPermissionMustHave, response, request, CustomerInstallationService.CustomerInstallationService.UpdateCustomerInstallation)
		break
	}
}

func (input customerInstallationEndpoint) CustomerSiteInSiteInstallationEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerSiteInSiteInstallationEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustomerInstallation()+common.ViewDataPermissionMustHave, response, request, CustomerInstallationService.CustomerInstallationService.ViewCustomerSiteInInstallationService)
		break
	}
}

func (input customerInstallationEndpoint) CustomerInstallationInSiteInstallationEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerInstallationInSiteInstallationEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustomerInstallation()+common.ViewDataPermissionMustHave, response, request, CustomerInstallationService.CustomerInstallationService.ViewCustomerInstallationInInstallationService)
		break
	}
}

func (input customerInstallationEndpoint) DetailInstallationEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DetailInstallationEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustomerInstallation()+common.ViewDataPermissionMustHave, response, request, CustomerInstallationService.CustomerInstallationService.ViewCustomerSiteInstallationByInstallationIDService)
		break
	}
}
