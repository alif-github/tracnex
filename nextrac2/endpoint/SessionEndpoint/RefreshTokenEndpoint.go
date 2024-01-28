package SessionEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/session/RefreshToken"
)

type refreshTokenEndpoint struct {
	endpoint.AbstractEndpoint
}

var RefreshTokenEndpoint = refreshTokenEndpoint{}.New()

func (input refreshTokenEndpoint) New() (output refreshTokenEndpoint) {
	output.FileName = "RefreshTokenEndpoint.go"
	return
}

func (input refreshTokenEndpoint) RefreshTokenEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "TokenEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, RefreshToken.RefreshTokenService.StartService)
}
