package RegisterUserEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/PKCEService"
)

type registerUserPKCEEndpoint struct {
	endpoint.AbstractEndpoint
}

var RegisterUserPKCEEndpoint = registerUserPKCEEndpoint{}.New()

func (input registerUserPKCEEndpoint) New() (output registerUserPKCEEndpoint) {
	output.FileName = "RegisterUserEndpoint.go"
	return
}

func (input registerUserPKCEEndpoint) RegistrationUserPKCEEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "RegistrationUserPKCEEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuRegistUnregistPKCEUserConstanta+common.InsertDataPermissionMustHave, responseWriter, request, PKCEService.PkceService.RegistrationUserPKCE)
}

func (input registerUserPKCEEndpoint) UnregisterUserPKCEEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "UnregisterUserPKCEEndpoint"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuRegistUnregistPKCEUserConstanta+common.DeleteDataPermissionMustHave, responseWriter, request, PKCEService.PkceService.UnregisterUser)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuRegistUnregistPKCEUserConstanta+common.ViewDataPermissionMustHave, responseWriter, request, PKCEService.PkceService.GetListUserCustomForUnregister)
	}
}

func (input registerUserPKCEEndpoint) ViewForUnregisterPKCEEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ViewForUnregisterPKCEEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuRegistUnregistPKCEUserConstanta+common.ViewDataPermissionMustHave, responseWriter, request, PKCEService.PkceService.ViewUserPKCEForUnregister)
}