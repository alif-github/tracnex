package SubDistrictEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/SubDistrictService"
)

type subDistrictEndpoint struct {
	endpoint.AbstractEndpoint
}

var SubDistrictEndpoint = subDistrictEndpoint{}.New()

func (input subDistrictEndpoint) New() (output subDistrictEndpoint) {
	output.FileName = "SubDistrictEndpoint.go"
	return
}

func (input subDistrictEndpoint) SubDistrictEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "SubDistrictEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, SubDistrictService.SubDistrictService.GetListSubDistrict)
		break
	}
}

func (input subDistrictEndpoint) SubDistrictEndpointWhithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "SubDistrictEndpointWhithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, SubDistrictService.SubDistrictService.ViewSubDistrictService)
		break
	}
}

func (input subDistrictEndpoint) InitiateSubDistrictEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateSubDistrictEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, SubDistrictService.SubDistrictService.InitiateSubDistrict)
		break
	}
}
