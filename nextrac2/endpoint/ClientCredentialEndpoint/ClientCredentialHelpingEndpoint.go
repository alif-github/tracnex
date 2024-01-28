package ClientCredentialEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ClientCredentialService"
)

type clientCredentialEndpoint struct {
	endpoint.AbstractEndpoint
}

var ClientCredentialEndpoint = clientCredentialEndpoint{}.New()

func (input clientCredentialEndpoint) New() (output clientCredentialEndpoint) {
	output.FileName = "ClientCredentialHelpingEndpoint.go"
	return
}

func (input clientCredentialEndpoint) ClientCredentialHelping(response http.ResponseWriter, request *http.Request) {
	funcName := "ClientCredentialHelping"
	input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ClientCredentialService.ClientCredentialService.InsertClientCredential)
}
