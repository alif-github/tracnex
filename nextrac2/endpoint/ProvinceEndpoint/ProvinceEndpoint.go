package ProvinceEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MasterDataService/ProvinceService"
)

type provinceEndpoint struct {
	endpoint.AbstractEndpoint
}

var ProvinceEndpoint = provinceEndpoint{}.New()

func (input provinceEndpoint) New() (output provinceEndpoint) {
	output.FileName = "ProvinceEndpoint.go"
	return
}

func (input provinceEndpoint) ProvinceEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ProvinceEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, ProvinceService.ProvinceService.GetListLocalProvince)
	}
}

func (input provinceEndpoint) ProvinceAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ProvinceAdminEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftProvince+common.ViewDataPermissionMustHave, response, request, ProvinceService.ProvinceService.GetListLocalAdminProvince)
	}
}

func (input provinceEndpoint) ResetLastSyncProvinceEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ResetLastSyncProvinceEndpoint"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ProvinceService.ProvinceService.ResetProvinceService)
	}
}
