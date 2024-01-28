package UserRegistrationDetailEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/RegistrationNamedUserService"
	"nexsoft.co.id/nextrac2/service/ResendOTPService"
	"nexsoft.co.id/nextrac2/service/UserRegistrationService"
)

type userRegistrationDetailEndpoint struct {
	endpoint.AbstractEndpoint
}

var UserRegistrationDetailEndpoint = userRegistrationDetailEndpoint{}.New()

func (input userRegistrationDetailEndpoint) New() (output userRegistrationDetailEndpoint) {
	output.FileName = "UserRegistrationDetailEndpoint.go"
	return
}

func (input userRegistrationDetailEndpoint) UserRegistrationDetailNamedUserWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "UserRegistrationDetailNamedUserWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuRegisterNamedOtherUser+common.InsertDataPermissionMustHave, responseWriter, request, RegistrationNamedUserService.RegistrationNamedUserService.InsertNamedUser)
	}
}

func (input userRegistrationDetailEndpoint) UserRegistrationDetailNamedUserClientMappingWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "UserRegistrationDetailNamedUserClientMappingWithoutParam"
	switch request.Method {
	case "POST":
		//input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuRegisterNexmileUser+common.InsertDataPermissionMustHave, responseWriter, request, RegistrationNamedUserService.RegistrationNamedUserService.RegisterOrRenewLicenseNamedUser)
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", responseWriter, request, RegistrationNamedUserService.RegistrationNamedUserService.RegisterOrRenewLicenseNamedUser)
	}
}

func (input userRegistrationDetailEndpoint) UserRegistrationDetailNamedUserCheckWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "UserRegistrationDetailNamedUserCheckWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", responseWriter, request, RegistrationNamedUserService.RegistrationNamedUserService.CheckEmailAndPhoneBeforeRegisterNamedUser)
	}
}

func (input userRegistrationDetailEndpoint) UserRegistrationDetailNamedUserWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "UserRegistrationDetailNamedUserWithParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuUnregisterUserLicense+common.UpdateDataPermissionMustHave, responseWriter, request, RegistrationNamedUserService.RegistrationNamedUserService.UnregisterNamedUser)
	}
}

func (input userRegistrationDetailEndpoint) CheckLicenseNamedUser(response http.ResponseWriter, request *http.Request) {
	funcName := "CheckLicenseNamedUser"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuCheckLicenseQuota+common.ViewDataPermissionMustHave, response, request, UserRegistrationService.UserRegistrationService.CheckLicenseNamedUser)
	}
}

//func (input userRegistrationDetailEndpoint) RegisterAndActivateUserNexmileNexstar(response http.ResponseWriter, request *http.Request) {
//	funcName := "RegisterAndActivateUserNexmileNexstar"
//	switch request.Method {
//	case "POST":
//		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, RegistrationNamedUserService.RegistrationNamedUserService.RegisterOrRenewLicenseNamedUser)
//	}
//}

func (input userRegistrationDetailEndpoint) ResendActivationNexmileNexstar(response http.ResponseWriter, request *http.Request) {
	funcName := "ResendActivationNexmileNexstar"
	switch request.Method {
	case "PUT":
		input.ServeWhiteListEndpoint(funcName, false, response, request, ResendOTPService.ResendOTPService.GenerateResendOTPService)
	}
}
