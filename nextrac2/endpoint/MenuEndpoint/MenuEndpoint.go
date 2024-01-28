package MenuEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MenuService"
)

type menuEndpoint struct {
	endpoint.AbstractEndpoint
}

var MenuEndpoint menuEndpoint

func (input menuEndpoint) ServiceMenuParentSysUserWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ServiceMenuParentSysUserWithoutParam"
	switch request.Method {
	case "GET" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.GetListParentMenu)
	}
}

func (input menuEndpoint) ServiceMenuParentSysAdminWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ServiceMenuParentSysAdminWithoutParam"
	switch request.Method {
	case "GET" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.GetListParentMenuSysadmin)
	}
}

func (input menuEndpoint) ServiceMenuParentSysAdminWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ServiceMenuParentSysAdminWithParam"
	switch request.Method {
	case "PUT" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.UpdateMenuParentSysAdmin)
	}
}

func (input menuEndpoint) ServiceMenuParentSysUserWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ServiceMenuParentSysUserWithParam"
	switch request.Method {
	case "PUT" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.UpdateMenuParentSysUser)
	}
}

func (input menuEndpoint) ServiceMenuServiceWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ServiceMenuServiceWithParam"
	switch request.Method {
	case "PUT" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.UpdateMenuService)
	case "GET" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.GetListServiceMenu)
	}
}

func (input menuEndpoint) ServiceMenuItemWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ServiceMenuItemWithParam"
	switch request.Method {
	case "PUT" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.UpdateMenuItem)
	case "GET" :
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, MenuService.MenuService.GetListMenuItem)
	}
}