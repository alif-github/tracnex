package DistrictEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MasterDataService/DistrictService"
)

type districtEndpoint struct {
	endpoint.AbstractEndpoint
}

var DistrictEndpoint = districtEndpoint{}.New()

func (input districtEndpoint) New() (output districtEndpoint) {
	output.FileName = "DistrictEndpoint.go"
	return
}

func (input districtEndpoint) DistrictEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DistrictEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, DistrictService.DistrictService.GetListLocalDistrict)
	}
}

func (input districtEndpoint) DistrictAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DistrictAdminEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftDistrict+common.ViewDataPermissionMustHave, response, request, DistrictService.DistrictService.GetListAdminLocalDistrict)
	}
}
