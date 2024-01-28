package UrbanVillageEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/UrbanVillageService"
)

type urbanVillageEndpoint struct {
	endpoint.AbstractEndpoint
}

var UrbanVillageEndpoint = urbanVillageEndpoint{}.New()

func (input urbanVillageEndpoint) New() (output urbanVillageEndpoint) {
	output.FileName = "UrbanVillageEndpoint.go"
	return
}

func (input urbanVillageEndpoint) UrbanVillageEndpointWhitoutParam(response http.ResponseWriter, request *http.Request)  {
	funcName := "UrbanVillageEndpointWhitoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, UrbanVillageService.UrbanVillageService.GetListUrbanVillage)
		break
	}
}

func (input urbanVillageEndpoint) UrbanVillageEndpointWhitPathParam(response http.ResponseWriter, request *http.Request)  {
	funcName := "UrbanVillageEndpointWhitPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, UrbanVillageService.UrbanVillageService.ViewUrbanVillage)
		break
	}
}

func (input urbanVillageEndpoint) InitiateUrbanVillageEndpoint(response http.ResponseWriter, request *http.Request)  {
	funcName := "InitiateUrbanVillageEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, UrbanVillageService.UrbanVillageService.InitiateUrbanVillage)
		break
	}
}
