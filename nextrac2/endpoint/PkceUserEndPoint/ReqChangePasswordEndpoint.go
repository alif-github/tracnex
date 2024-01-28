package PkceUserEndPoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/PkceUserService"
)

type changePasswordEndPoint struct {
	endpoint.AbstractEndpoint
}

var ChangePasswordEndPoint changePasswordEndPoint

func (input changePasswordEndPoint) ChangePasswordEndPoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ChangePasswordEndPoint"
	input.FileName = "ChangePasswordEndPoint.go"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuChangePasswordPKCENexmile+common.ChangePasswordPermissionMustHave, responseWriter, request, PkceUserService.ReqChangePasswordService.RequestChangePassword)
}