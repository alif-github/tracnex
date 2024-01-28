package LicenseTypeEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/LicenseTypeService"
)

type licenseTypeEnpoint struct {
	endpoint.AbstractEndpoint
}

var LicenseTypeEndpoint = licenseTypeEnpoint{}.New()

func (input licenseTypeEnpoint) New() (output licenseTypeEnpoint) {
	output.FileName = "LicenseTypeEndpoint.go"
	return
}

func (input licenseTypeEnpoint) getMenuCodeLicenseType() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterLisensiTipeLisensiRedesign, constanta.MenuUserMasterLisensiTipeLisensi)
}

func (input licenseTypeEnpoint) LicenseTypeEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LicenseTypeEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseType()+common.ViewDataPermissionMustHave, response, request, LicenseTypeService.LicenseTypeService.ViewLicenseType)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseType()+common.UpdateDataPermissionMustHave, response, request, LicenseTypeService.LicenseTypeService.UpdateLicenseType)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseType()+common.DeleteDataPermissionMustHave, response, request, LicenseTypeService.LicenseTypeService.DeleteLicenseType)
		break
	}
}

func (input licenseTypeEnpoint) LicenseTypeEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LicenseTypeEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseType()+common.ViewDataPermissionMustHave, response, request, LicenseTypeService.LicenseTypeService.GetListLicenseType)
		break
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseType()+common.InsertDataPermissionMustHave, response, request, LicenseTypeService.LicenseTypeService.InsertLicenseType)
		break
	}
}

func (input licenseTypeEnpoint) InitiateLicenseTypeEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateLicenseTypeEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseType()+common.ViewDataPermissionMustHave, response, request, LicenseTypeService.LicenseTypeService.InitiateLicenseType)
		break
	}
}
