package ActivationLicenseEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ActivationLicenseService"
)

type activationLicenseEndpoint struct {
	endpoint.AbstractEndpoint
}

var ActivationLicenseEndpoint = activationLicenseEndpoint{}.New()

func (input activationLicenseEndpoint) New() (output activationLicenseEndpoint) {
	output.FileName = "ActivationLicenseEndpoint.go"
	return
}

func (input activationLicenseEndpoint) ActivationLicenseEndpointWithoutPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ActivationLicenseEndpointWithoutPathParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuUserClientMappingBackEndConstanta+common.InsertDataPermissionMustHave, response, request, ActivationLicenseService.ActivationLicenseService.ActivateLicense)
		break
	}
}
