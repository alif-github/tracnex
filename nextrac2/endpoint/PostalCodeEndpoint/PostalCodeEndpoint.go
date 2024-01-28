package PostalCodeEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/PostalCodeService"
)

type postalCodeEndpoint struct {
	endpoint.AbstractEndpoint
}

var PostalCodeEndpoint = postalCodeEndpoint{}.New()

func (input postalCodeEndpoint) New() (output postalCodeEndpoint) {
	output.FileName = "PostalCodeEndpoint.go"
	return
}

func (input postalCodeEndpoint) PostalCodeEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "PostalCodeEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PostalCodeService.PostalCodeService.GetListPostalCode)
		break
	}
}

func (input postalCodeEndpoint) PostalCodeEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "PostalCodeEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PostalCodeService.PostalCodeService.ViewPostalCode)
		break
	}
}

func (input postalCodeEndpoint) InitiatePostalCodeEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiatePostalCodeEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PostalCodeService.PostalCodeService.InitiatePostalCode)
		break
	}
}
