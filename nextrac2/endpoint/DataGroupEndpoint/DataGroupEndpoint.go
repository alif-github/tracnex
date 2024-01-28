package DataGroupEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/DataGroupService"
	"nexsoft.co.id/nextrac2/service/DataScopeService"
)

type dataGroupEndpoint struct {
	endpoint.AbstractEndpoint
}

var DataGroupEndpoint dataGroupEndpoint

func (input dataGroupEndpoint) DataGroupEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DataGroupEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftDataGroupConstanta+common.InsertDataPermissionMustHave, response, request, DataGroupService.DataGroupService.InsertDataGroup)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftDataGroupConstanta+common.ViewDataPermissionMustHave, response, request, DataGroupService.DataGroupService.GetListDataGroup)
		break
	}
}

func (input dataGroupEndpoint) DataGroupEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DataGroupEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftDataGroupConstanta+common.ViewDataPermissionMustHave, response, request, DataGroupService.DataGroupService.ViewDataGroup)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftDataGroupConstanta+common.UpdateDataPermissionMustHave, response, request, DataGroupService.DataGroupService.UpdateDataGroup)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftDataGroupConstanta+common.DeleteDataPermissionMustHave, response, request, DataGroupService.DataGroupService.DeleteDataGroup)
		break
	}
}

func (input dataGroupEndpoint) InitiateDataGroupEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateDataGroupEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftDataGroupConstanta+common.ViewDataPermissionMustHave, response, request, DataGroupService.DataGroupService.InitiateGetListDataGroup)
		break
	}
}

func (input dataGroupEndpoint) InsertHelperDataScopeEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InsertHelperDataScopeEndpoint"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftDataGroupConstanta+common.ViewDataPermissionMustHave, response, request, DataScopeService.DataScopeService.InsertDataScope)
		break
	}
}
