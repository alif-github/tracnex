package EnumEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/EnumService"
)

type enumEndpoint struct {
	endpoint.AbstractEndpoint
}

var EnumEndpoint = enumEndpoint{}.New()

func (input enumEndpoint) New() (output enumEndpoint) {
	output.FileName = "EmployeeEndpoint.go"
	return
}

func (input enumEndpoint) EnumEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "EnumEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, EnumService.EnumService.ViewDetailEnum)
		break
	}
}
