package UserLicenseEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/UserLicenseService"
)

type userLicenseEndpoint struct {
	endpoint.AbstractEndpoint
}

var UserLicenseEndpoint = userLicenseEndpoint{}.New()

func (input userLicenseEndpoint) New() (output userLicenseEndpoint) {
	output.FileName = "UserLicenseEndpoint.go"
	return
}

func (input userLicenseEndpoint) getMenuCodeUserLicense() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterLisensiUserLisensiRedesign, constanta.MenuUserMasterLisensiUserLisensi)
}

func (input userLicenseEndpoint) InitiateUserLicenseEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateUserLicenseEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserLicense()+common.ViewDataPermissionMustHave, response, request, UserLicenseService.UserLicenseService.InitiateGetListUserLicense)
	}

}

func (input userLicenseEndpoint) UserLicenseWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserLicenseWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserLicense()+common.ViewDataPermissionMustHave, response, request, UserLicenseService.UserLicenseService.GetListUserLicense)
	}
}

func (input userLicenseEndpoint) ViewDetailUserLicenseEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ViewDetailUserLicenseEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserLicense()+common.ViewDataPermissionMustHave, response, request, UserLicenseService.UserLicenseService.ViewDetailUserLicense)
	}
}

func (input userLicenseEndpoint) InitiateViewDetailUserLicenseEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateViewDetailUserLicenseEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserLicense()+common.ViewDataPermissionMustHave, response, request, UserLicenseService.UserLicenseService.InitiateViewUserLicense)
	}
}

func (input userLicenseEndpoint) TransferUserLicenseEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "TransferUserLicenseEndpointWithPathParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeUserLicense()+common.UpdateDataPermissionMustHave, response, request, UserLicenseService.UserLicenseService.TransferUserLicense)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeUserLicense()+common.ViewDataPermissionMustHave, response, request, UserLicenseService.UserLicenseService.ViewTransferUserLicense)
		break
	}
}
