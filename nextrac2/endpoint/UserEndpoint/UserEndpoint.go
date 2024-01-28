package UserEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/UserService/ForgetPassword"
)

type userEndpoint struct {
	endpoint.AbstractEndpoint
}

var UserEndpoint = userEndpoint{}.New()

func (input userEndpoint) New() (output userEndpoint) {
	output.FileName = "UserEndpoint.go"
	return
}

func (input userEndpoint) UserEndpointWithoutParamResetPassword(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointWithoutParamResetPassword"
	switch request.Method {
	case http.MethodPost:
		input.ServeWhiteListEndpoint(funcName, false, response, request, ForgetPassword.ForgetPasswordService.ResetPassword)
		break
	}
}

func (input userEndpoint) UserEndpointWithoutParamChangePassword(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointWithoutParamChangePassword"
	switch request.Method {
	case http.MethodPost:
		input.ServeWhiteListEndpoint(funcName, false, response, request, ForgetPassword.ForgetPasswordService.ChangePassword)
		break
	}
}