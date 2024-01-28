package UserEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/UserService/CRUDUserService"
	"nexsoft.co.id/nextrac2/service/UserService/UserInvitationService"
)

type userSysAdminEndpoint struct {
	endpoint.AbstractEndpoint
}

var UserSysAdminEndpoint = userSysAdminEndpoint{}.New()

func (input userSysAdminEndpoint) New() (output userSysAdminEndpoint) {
	output.FileName = "UserSysAdminEndpoint.go"
	return
}

func (input userSysAdminEndpoint) UserEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.InsertDataPermissionMustHave, response, request, CRUDUserService.UserService.InsertUserSysAdmin)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.GetListUser)
	}
}

func (input userSysAdminEndpoint) InitiateUserParam(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateUserParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.InitiateUser)
	}
}

func (input userSysAdminEndpoint) UserVerificationEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserVerificationEndpointWithParam"
	switch request.Method {
	case http.MethodPost:
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, CRUDUserService.UserService.ResendVerificationCode)
	}
}

func (input userSysAdminEndpoint) CheckUserAuthEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CheckUserAuthEndpointWithoutParam"
	switch request.Method {
	case http.MethodPost:
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.CheckUserAuth)
		break
	}
}

func (input userSysAdminEndpoint) SysUserLoginEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "SysUserLoginEndpointWithoutParam"
	switch request.Method {
	case http.MethodGet:
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.GetListUserSysUserActive)
		break
	case http.MethodPost:
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.UpdateDataPermissionMustHave, response, request, CRUDUserService.UserService.KillSessionUserActive)
		break
	}
}

func (input userSysAdminEndpoint) CheckUsernameAuthEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CheckUsernameAuthEndpointWithoutParam"
	switch request.Method {
	case http.MethodPost:
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.InsertDataPermissionMustHave, response, request, CRUDUserService.UserService.CheckUsernameAuth)
		break
	}
}

func (input userSysAdminEndpoint) UserEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointWithoutParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.DeleteDataPermissionMustHave, response, request, CRUDUserService.UserService.DeleteUserSysAdmin)
		break
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.UpdateDataPermissionMustHave, response, request, CRUDUserService.UserService.UpdateUserSysAdmin)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.ViewUser)
	}
}

func (input userSysAdminEndpoint) ChangePasswordAdminUser(response http.ResponseWriter, request *http.Request) {
	funcName := "ChangePasswordAdminUser"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftChangePassword+common.UpdateDataPermissionMustHave, response, request, CRUDUserService.UserService.ChangePasswordUser)
		break
	}
}

func (input userSysAdminEndpoint) ProfileAdminUser(response http.ResponseWriter, request *http.Request) {
	funcName := "ProfileAdminUser"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftProfileSetting+common.UpdateDataPermissionMustHave, response, request, CRUDUserService.UserService.UpdateAdminProfile)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftProfileSetting+common.ViewDataPermissionMustHave, response, request, CRUDUserService.UserService.ViewProfileSettingUser)
		break
	}
}

func (input userSysAdminEndpoint) UserEndpointHelpingDeleted(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointHelpingDeleted"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.UpdateDataPermissionMustHave, response, request, CRUDUserService.UserService.UpdateHelpingUserDeleted)
		break
	}
}

func (input userSysAdminEndpoint) UserEndpointHelpingFirstName(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointHelpingFirstName"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta+common.UpdateDataPermissionMustHave, response, request, CRUDUserService.UserService.UpdateHelpingClientND6FirstName)
		break
	}
}

func (input userSysAdminEndpoint) Invitation(response http.ResponseWriter, request *http.Request) {
	funcName := "Invitation"
	input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftUserConstanta + common.InsertDataPermissionMustHave, response, request, UserInvitationService.InvitationService.SendInvitation)
}
