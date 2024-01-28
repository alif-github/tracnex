package ClientRegistrationLogEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ClientRegistrationLogService"
)

type clientRegistrationLogEndpoint struct {
	endpoint.AbstractEndpoint
}

var ClientRegistrationLogEndpoint = clientRegistrationLogEndpoint{}.New()

func (input clientRegistrationLogEndpoint) New() (output clientRegistrationLogEndpoint) {
	output.FileName = "ClientRegistrationLogEndpoint.go"
	return
}

func (input clientRegistrationLogEndpoint) ClientRegistrationLogWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ClientRegistrationLogWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftCustomerConstanta+common.ViewDataPermissionMustHave, response, request, ClientRegistrationLogService.ClientRegistrationLogService.GetListLog)
	}
}