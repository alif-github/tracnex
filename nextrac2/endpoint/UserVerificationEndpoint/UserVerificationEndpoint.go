package UserVerificationEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/UserVerificationService"
)

type userVerificationEnpoint struct {
	endpoint.AbstractEndpoint
}

var UserVerificationEndpoint = userVerificationEnpoint{}.New()

func (input userVerificationEnpoint) New() (output userVerificationEnpoint) {
	output.FileName = "UserVerificationEndpoint.go"
	return
}

func (input userVerificationEnpoint) UserVerificationWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserVerificationWithoutParam"
	switch request.Method {
	case "PUT":
		input.ServeWhiteListEndpoint(funcName, false, response, request, UserVerificationService.UserVerificationService.VerifyingUserService)
	}
}
