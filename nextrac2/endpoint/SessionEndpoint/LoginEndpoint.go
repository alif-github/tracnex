package SessionEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/session/Login"
)

type loginEndpoint struct {
	endpoint.AbstractEndpoint
}

var LoginEndpoint = loginEndpoint{}.New()

func (input loginEndpoint) New() (output loginEndpoint) {
	output.FileName = "LoginEndpoint.go"
	return
}

func (input loginEndpoint) AuthorizeEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "AuthorizeEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.AuthorizeService.StartService)
}

func (input loginEndpoint) VerifyEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "VerifyEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.VerifyService.StartService)
}

func (input loginEndpoint) TokenEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "TokenEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.TokenService.UserTokenService)
}

func (input loginEndpoint) TokenAdminEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "TokenEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.TokenService.NexsoftTokenService)
}

func (input loginEndpoint) TokenClientEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "TokenClientEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.TokenService.ClientTokenService)
}

func (input loginEndpoint) LoginNexmileEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "LoginNexmileEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.NexmileLoginService.LoginNexmileService)
}

func (input loginEndpoint) SysUserLoginEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "SysUserLoginEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.NexmileLoginService.LoginNexmileService)
}

func (input loginEndpoint) LoginGroChatEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "LoginGroChatEndpoint"
	input.ServeWhiteListEndpoint(funcName, false, responseWriter, request, Login.GroChatLoginService.Login)
}
