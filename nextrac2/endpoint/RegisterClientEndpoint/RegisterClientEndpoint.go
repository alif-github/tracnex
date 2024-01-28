package RegisterClientEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/ClientRegistrationNonOnPremiseService"
	"nexsoft.co.id/nextrac2/service/ClientService"
)

type registerClientEndpoint struct {
	endpoint.AbstractEndpoint
}

var RegisterClientEndpoint registerClientEndpoint

func (input registerClientEndpoint) RegisterClientEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "RegisterClientEndpoint"
	input.FileName = "RegisterClientEndpoint.go"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, ClientService.ClientService.RegistrationClient)
}

func (input registerClientEndpoint) InitiateClientEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "InitiateClientEndpoint"
	input.FileName = "RegisterClientEndpoint.go"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, ClientService.ClientService.InitiateRegistrationClient)
}

func (input registerClientEndpoint) RegisterClientNonOnPremiseEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "RegisterClientNonOnPremiseEndpoint"
	input.FileName = "RegisterClientEndpoint.go"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, ClientRegistrationNonOnPremiseService.ClientRegistrationNonOnPremiseService.InsertClientRegistNonOnPremise)
}