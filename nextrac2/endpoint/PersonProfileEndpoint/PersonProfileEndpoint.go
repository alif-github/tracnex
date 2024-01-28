package PersonProfileEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MasterDataService/PersonProfileService"
)

type personProfileEndpoint struct {
	endpoint.AbstractEndpoint
}

var PersonProfileEndpoint = personProfileEndpoint{}.New()

func (input personProfileEndpoint) New() (output personProfileEndpoint) {
	output.FileName = "PersonProfileEndpoint.go"
	return
}

func (input personProfileEndpoint) PersonProfileEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "PersonProfileEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PersonProfileService.PersonProfileServie.GetListPersonProfile)
	}
}

func (input personProfileEndpoint) PersonProfileEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "PersonProfileEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PersonProfileService.PersonProfileServie.ViewPersonProfile)
	}
}
