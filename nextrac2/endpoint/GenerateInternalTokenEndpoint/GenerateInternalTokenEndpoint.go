package GenerateInternalTokenEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/GenerateInternalToken"
)

type generateInternalTokenEndpoint struct {
	endpoint.AbstractEndpoint
}

var GenerateInTokenEndpoint = generateInternalTokenEndpoint{}.New()

func (input generateInternalTokenEndpoint) New() (output generateInternalTokenEndpoint) {
	output.FileName = "GenerateInternalTokenEndpoint.go"
	return
}

func (input generateInternalTokenEndpoint) GenerateInternalToken(responseWriter http.ResponseWriter, request * http.Request) {
	funcName := "GenerateInternalToken"
	switch request.Method {
	case "POST":
		input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, GenerateInternalToken.GenerateInTokenService.StartGenerateInToken)
	}
}