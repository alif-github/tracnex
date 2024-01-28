package PositionEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MasterDataService/PositionService"
)

type positionEndpoint struct {
	endpoint.AbstractEndpoint
}

var PositionEndpoint = positionEndpoint{}.New()

func (input positionEndpoint) New() (output positionEndpoint) {
	output.FileName = "PositionEndpoint.go"
	return
}

func (input positionEndpoint) PositionEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "PositionEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PositionService.PositionService.GetListPosition)
		break
	}
}

func (input positionEndpoint) PositionEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "PositionEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PositionService.PositionService.ViewPosition)
		break
	}
}
