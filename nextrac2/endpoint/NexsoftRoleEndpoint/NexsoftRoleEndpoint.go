package NexsoftRoleEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/NexsoftRoleService"
)

type nexsoftRoleEndpoint struct {
	endpoint.AbstractEndpoint
}

var NexsoftRoleEndpoint = nexsoftRoleEndpoint{}.New()

func (input nexsoftRoleEndpoint) New() (output nexsoftRoleEndpoint) {
	output.FileName = "NexsoftRoleEndpoint.go"
	return
}

func (input nexsoftRoleEndpoint) NexsoftRoleWithoutPathParamEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "NexsoftRoleEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftSysadminRoleConstanta + common.InsertDataPermissionMustHave, response, request, NexsoftRoleService.InsertNexsoftRoleService.InsertNexsoftRole)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName,false, common.ReadDataAPIMustHave, constanta.MenuNexsoftSysadminRoleConstanta + common.ViewDataPermissionMustHave, response, request, NexsoftRoleService.GetListNexsoftRoleService.GetListNexsoftRole)
		break
	}
}

func (input nexsoftRoleEndpoint) NexsoftRoleWithPathParamEndpoint(response http.ResponseWriter, request *http.Request)  {
	funcName := "NexsoftRoleWithPathParamEndpoint"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName,false, common.WriteDataAPIMustHave, constanta.MenuNexsoftSysadminRoleConstanta + common.UpdateDataPermissionMustHave, response, request, NexsoftRoleService.UpdateNexsoftRoleService.UpdateNexsoftRole)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftSysadminRoleConstanta + common.ViewDataPermissionMustHave, response, request, NexsoftRoleService.ViewNexsoftRoleService.ViewNexsoftRole)
		break
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftSysadminRoleConstanta + common.DeleteDataPermissionMustHave, response, request, NexsoftRoleService.DeleteNexsoftRoleService.DeleteNexsoftRole)
		break
	}
}

func (input nexsoftRoleEndpoint) InitiateNexsoftRoleEnpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateNexsoftRoleEnpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftSysadminRoleConstanta + common.ViewDataPermissionMustHave, response, request, NexsoftRoleService.GetListNexsoftRoleService.InitiateNexsoftRole)
		break
	}
	
}