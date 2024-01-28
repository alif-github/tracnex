package ActivationUserNexmileEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ActivationUserNexmileService"
)

type activationUserNexmileEndpoint struct {
	endpoint.AbstractEndpoint
}

var ActivationUserNexmileEndpoint = activationUserNexmileEndpoint{}.New()

func (input activationUserNexmileEndpoint) New() (output activationUserNexmileEndpoint) {
	output.FileName = "ActivationUserNexmileEndpoint.go"
	return
}

func (input activationUserNexmileEndpoint) ActivateUserNexmileEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ActivationLicenseEndpointWithoutPathParam"

	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuActivatedNexmileUser+common.UpdateDataPermissionMustHave, response, request, ActivationUserNexmileService.ActivationUserNexmileService.ActivateUserNexmile)
		break
	}
}
