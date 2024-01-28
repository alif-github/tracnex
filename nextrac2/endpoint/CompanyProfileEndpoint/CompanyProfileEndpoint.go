package CompanyProfileEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MasterDataService/CompanyProfileService"
)

type companyProfileEndpoint struct {
	endpoint.AbstractEndpoint
}

var CompanyProfileEndpoint = companyProfileEndpoint{}.New()

func (input companyProfileEndpoint) New() (output companyProfileEndpoint) {
	output.FileName = "CompanyProfileEndpoint.go"
	return
}

func (input companyProfileEndpoint) CompanyProfileEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CompanyProfileEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, CompanyProfileService.CompanyProfileService.GetListCompanyProfile)
		break
	}
}

func (input companyProfileEndpoint) CompanyProfileEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CompanyProfileEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, CompanyProfileService.CompanyProfileService.ViewCompanyProfile)
		break
	}
}
