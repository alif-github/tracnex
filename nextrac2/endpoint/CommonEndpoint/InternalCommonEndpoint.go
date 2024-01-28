package CommonEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/common"
)

type internalCommonEndpoint struct {
	endpoint.AbstractEndpoint
}

var InternalCommonEndpoint = internalCommonEndpoint{}.New()

func (input internalCommonEndpoint) New() (output internalCommonEndpoint) {
	output.FileName = "InternalCommonEndpoint.go"
	return
}

func (input internalCommonEndpoint) NotifyAddClientEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "StartInternalEndpoint"
	input.ServeInternalValidationEndpoint(funcName, false, true, responseWriter, request, common.NotifyAddClientService.StartService)
}

func (input internalCommonEndpoint) DeleteTokenEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "DeleteTokenEndpoint"
	input.ServeInternalValidationEndpoint(funcName, false, true, responseWriter, request, common.DeleteTokenService.StartService)
}

func (input internalCommonEndpoint) InternalRegisterUserEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "InternalRegisterUserEndpoint"
	input.ServeInternalValidationEndpoint(funcName, false, true, responseWriter, request, common.InsertUserFromInternalService.StartService)
}
