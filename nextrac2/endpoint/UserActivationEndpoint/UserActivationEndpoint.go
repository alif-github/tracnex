package UserActivationEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/UserActivationService"
)

type userActivationEndpoint struct {
	endpoint.AbstractEndpoint
}

var UserActivationEndpoint userActivationEndpoint

func (input userActivationEndpoint) UserActivationEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "UserActivationEndpoint"
	input.FileName = "UserActivationEndpoint.go"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, UserActivationService.UserActivationService.EmailActivation)
}