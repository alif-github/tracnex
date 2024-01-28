package LicenseVariantEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/LicenseVariantService"
)

type licenseVariantEndpoint struct {
	endpoint.AbstractEndpoint
}

var LicenseVariantEndpoint = licenseVariantEndpoint{}.New()

func (input licenseVariantEndpoint) New() (output licenseVariantEndpoint) {
	output.FileName = "LicenseVariantEndpoint.go"
	return
}

func (input licenseVariantEndpoint) getMenuCodeLicenseVariant() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterLisensiVariantLisensiRedesign, constanta.MenuUserMasterLisensiVariantLisensi)
}

func (input licenseVariantEndpoint) LicenseVariantWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LicenseVariantWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseVariant()+common.InsertDataPermissionMustHave, response, request, LicenseVariantService.LicenseVariantService.InsertLicenseVariant)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseVariant()+common.ViewDataPermissionMustHave, response, request, LicenseVariantService.LicenseVariantService.GetListLicenseVariant)
	}
}

func (input licenseVariantEndpoint) LicenseVariantWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LicenseVariantWithParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseVariant()+common.UpdateDataPermissionMustHave, response, request, LicenseVariantService.LicenseVariantService.UpdateLicenseVariant)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseVariant()+common.ViewDataPermissionMustHave, response, request, LicenseVariantService.LicenseVariantService.ViewLicenseVariant)
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseVariant()+common.DeleteDataPermissionMustHave, response, request, LicenseVariantService.LicenseVariantService.DeleteLicenseVariant)
	}
}

func (input licenseVariantEndpoint) InitiateLicenseVariant(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateLicenseVariant"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseVariant()+common.ViewDataPermissionMustHave, response, request, LicenseVariantService.LicenseVariantService.InitiateGetListLicenseVariant)
}
