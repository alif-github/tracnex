package ClientMappingEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ClientMappingService"
	"nexsoft.co.id/nextrac2/service/ClientMappingService/InternalClientMappingService"
)

type clientMappingEndpoint struct {
	endpoint.AbstractEndpoint
}

var ClientMappingEndpoint = clientMappingEndpoint{}.New()

func (input clientMappingEndpoint) getMenuCodeCustomerClientMapping() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterConsumerClientMappingRedesign, constanta.MenuUserMasterConsumerClientMapping)
}

func (input clientMappingEndpoint) New() (output clientMappingEndpoint) {
	output.FileName = "InternalClientMappingEndpoint.go"
	return
}

func (input clientMappingEndpoint) InsertNewBranchToClientMappingEndPoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "InsertNewBranchToClientMappingEndPoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuUserClientMappingBackEndConstanta+common.InsertDataPermissionMustHave, responseWriter, request, ClientMappingService.ClientMappingService.InsertNewBranchToClientMapping)
}

func (input clientMappingEndpoint) DetailMappingEndPoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "DetailMappingEndPoint"
	input.ServeInternalValidationEndpoint(funcName, false, true, responseWriter, request, InternalClientMappingService.GetDetailClientMappingService.GetDetailClientMappings)
}

func (input clientMappingEndpoint) DetailMappingEndPointByClientID(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "DetailMappingEndPointByClientID"
	input.ServeInternalValidationEndpoint(funcName, false, true, responseWriter, request, InternalClientMappingService.GetClientMappingByClientIDService.GetClientMappingsByClientID)
}
