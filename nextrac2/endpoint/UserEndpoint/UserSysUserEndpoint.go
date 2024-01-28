package UserEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/UserService/CRUDUserService"
)

type userSysUserEndpoint struct {
	endpoint.AbstractEndpoint
}

var UserSysUserEndpoint = userSysUserEndpoint{}.New()

func (input userSysUserEndpoint) New() (output userSysUserEndpoint) {
	output.FileName = "UserSysUserEndpoint.go"
	return
}

func (input userSysUserEndpoint) UserEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuUserAdminSubPengguna+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.ViewUser)
	}
}

func (input userSysUserEndpoint) ProfileSettingSysUser(response http.ResponseWriter, request *http.Request) {
	funcName := "ProfileSettingSysUser"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuUserAdminSubPengguna+common.UpdateDataPermissionMustHave, response, request, CRUDUserService.UserService.UpdateUserProfile)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuProfileSettingUser+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.ViewProfileSettingUser)
		break
	}
}

func (input userSysUserEndpoint) ChangePasswordSysUser(response http.ResponseWriter, request *http.Request) {
	funcName := "ChangePasswordSysUser"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuUserAdminSubPengguna+common.ChangePasswordPermissionMustHave, response, request, CRUDUserService.UserService.ChangePasswordUser)
		break
	}
}

func (input userSysUserEndpoint) ResendVerificationCodeSysUser(response http.ResponseWriter, request *http.Request) {
	funcName := "ResendVerificationCodeSysUser"
	switch request.Method {
	case http.MethodPost:
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, CRUDUserService.UserService.ResendVerificationCode)
		break
	}
}
