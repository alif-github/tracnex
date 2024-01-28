package CustomerSIteEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/CustomerSiteService"
)

type customerSiteEndpoint struct {
	endpoint.AbstractEndpoint
}

var CustomerSiteEndpoint = customerSiteEndpoint{}.New()

func (input customerSiteEndpoint) New() (output customerSiteEndpoint) {
	output.FileName = "CustomerSiteEndpoint.go"
	return
}

func (input customerSiteEndpoint) getMenuCodeCustomerInstallation() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterKonsumenInstallationRedesign, constanta.MenuUserMasterKonsumenInstallation)
}

func (input customerSiteEndpoint) CustomerSiteEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerSiteEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustomerInstallation()+common.InsertDataPermissionMustHave, response, request, CustomerSiteService.CustomerSiteService.InsertCustomerSite)
		break
	}
}
