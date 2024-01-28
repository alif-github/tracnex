package UserRegistrationAdminEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/UserRegistrationAdminService"
)

type userRegistrationAdminEndpoint struct {
	endpoint.AbstractEndpoint
}

var UserRegistrationAdminEndpoint = userRegistrationAdminEndpoint{}.New()

func (input userRegistrationAdminEndpoint) New() (output userRegistrationAdminEndpoint) {
	output.FileName = "UserRegistrationAdminEndpoint.go"
	return
}

func (input userRegistrationAdminEndpoint) getMenuCodeUserRegistrationAdmin() string {
	return endpoint.GetMenuCode(constanta.MenuViewRegisterUserAdminRedesign, constanta.MenuViewRegisterUserAdmin)
}

func (input userRegistrationAdminEndpoint) UserRegistrationAdminWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserRegistrationAdminWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserRegistrationAdmin()+common.ViewDataPermissionMustHave, response, request, UserRegistrationAdminService.UserRegistrationAdminService.ViewUserRegistrationAdmin)
	}
}

func (input userRegistrationAdminEndpoint) UserRegistrationAdminWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserRegistrationAdminWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserRegistrationAdmin()+common.ViewDataPermissionMustHave, response, request, UserRegistrationAdminService.UserRegistrationAdminService.GetListUserRegistrationAdmin)
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuRegisterUserAdmin+common.InsertDataPermissionMustHave, response, request, UserRegistrationAdminService.UserRegistrationAdminService.InsertUserRegistrationAdmin)
	}
}

func (input userRegistrationAdminEndpoint) UserRegistrationAdminInitiate(response http.ResponseWriter, request *http.Request) {
	funcName := "UserRegistrationAdminInitiate"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserRegistrationAdmin()+common.ViewDataPermissionMustHave, response, request, UserRegistrationAdminService.UserRegistrationAdminService.InitiateGetListUserRegistrationAdmin)
}
