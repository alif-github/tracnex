package CompanyTitleEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MasterDataService/CompanyTitleService"
)

type companyTitleEndpoint struct {
	endpoint.AbstractEndpoint
}

var CompanyTitleEndpoint = companyTitleEndpoint{}.New()

func (input companyTitleEndpoint) New() (output companyTitleEndpoint) {
	output.FileName = "CompanyTitleEndpoint.go"
	return
}

func (input companyTitleEndpoint) CompanyTitleEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CompanyTitleEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, CompanyTitleService.CompanyTitleService.GetListCompanyTitle)
		break
	}
}

func (input companyTitleEndpoint) CompanyTitleEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CompanyTitleEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, CompanyTitleService.CompanyTitleService.ViewCompanyTitle)
		break
	}
}
