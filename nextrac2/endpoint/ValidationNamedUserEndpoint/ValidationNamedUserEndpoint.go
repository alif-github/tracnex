package ValidationNamedUserEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ValidationNamedUserService"
)

type validationNamedUserEndpoint struct {
	endpoint.AbstractEndpoint
}

var ValidationNamedUserEndpoint = validationNamedUserEndpoint{}.New()

func (input validationNamedUserEndpoint) New() (output validationNamedUserEndpoint) {
	output.FileName = "ValidationNamedUserEndpoint.go"
	return
}

func (input validationNamedUserEndpoint) ValidateNamedUserEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ValidateNamedUserEndpoint"
	switch request.Method {
	case "POST":
		// todo keperluan testing lepas permission
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ValidationNamedUserService.ValidationNamedUserService.ValidateNamedUser)
		//input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuValidationNamedUser+common.InsertDataPermissionMustHave, response, request, ValidationNamedUserService.ValidationNamedUserService.ValidateNamedUser)
	}
}
