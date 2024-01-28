package SessionEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/session/GetSessionService"
)

type getSessionEndpoint struct {
	endpoint.AbstractEndpoint
}

var GetSessionEndpoint = getSessionEndpoint{}.New()

func (input getSessionEndpoint) New() (output getSessionEndpoint) {
	output.FileName = "GetSessionEndpoint.go"
	return
}

func (input getSessionEndpoint) GetSessionEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "GetSessionEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, GetSessionService.GetSessionService.StartService)
}

func (input getSessionEndpoint) GetAdminSessionEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "GetAdminSessionEndpoint"
	input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, GetSessionService.GetSessionService.StartService)
}

func (input getSessionEndpoint) GetCurrentDatetimeEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "GetCurrentDatetimeEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, GetSessionService.GetSessionService.GetCurrentDateTimeService)
}

func (input getSessionEndpoint) GetDashboardViewEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "GetCurrentDatetimeEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, GetSessionService.GetSessionService.GetDashboardView)
}

func (input getSessionEndpoint) GetApplicationVersionEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "GetApplicationVersionEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, GetSessionService.GetSessionService.GetSystemVersion)
}
