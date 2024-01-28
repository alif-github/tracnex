package RoleEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/RoleService"
)

type roleEndpoint struct {
	endpoint.AbstractEndpoint
}

var RoleEndpoint = roleEndpoint{}.New()

func (input roleEndpoint) New() (output roleEndpoint) {
	output.FileName = "RoleEndpoint.go"
	return
}

func (input roleEndpoint) RoleWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "RoleWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftRoleConstanta+common.InsertDataPermissionMustHave, responseWriter, request, RoleService.RoleService.InsertRole)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoleConstanta+common.ViewDataPermissionMustHave, responseWriter, request, RoleService.RoleService.GetListRole)
	}
}

func (input roleEndpoint) RoleWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "RoleWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftRoleConstanta+common.DeleteDataPermissionMustHave, responseWriter, request, RoleService.RoleService.DeleteRole)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoleConstanta+common.ViewDataPermissionMustHave, responseWriter, request, RoleService.RoleService.ViewRole)
		break
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftRoleConstanta+common.UpdateDataPermissionMustHave, responseWriter, request, RoleService.RoleService.UpdateRole)
	}
}

func (input roleEndpoint) InitiateGetListRole(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetListRole"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoleConstanta+common.ViewDataPermissionMustHave, responseWriter, request, RoleService.RoleService.InitiateGetListRole)
	}
}
