package SessionEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	Login "nexsoft.co.id/nextrac2/service/session/Logout"
)

type logoutEndpoint struct {
	endpoint.AbstractEndpoint
}

var LogoutEndpoint = logoutEndpoint{}.New()

func (input logoutEndpoint) New() (output logoutEndpoint) {
	output.FileName = "LogoutEndpoint.go"
	return
}

func (input logoutEndpoint) LogoutEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "LogoutEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.LogoutService.StartService)
}
