package LicenseConfigEndpoint

import (
	"net/http"

	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/LicenseConfigService"
)

type licenseConfigEndpoint struct {
	endpoint.AbstractEndpoint
}

var LicenseConfigEndpoint = licenseConfigEndpoint{}.New()

func (input licenseConfigEndpoint) New() (output licenseConfigEndpoint) {
	output.FileName = "LicenseConfigEndpoint.go"
	return
}

func (input licenseConfigEndpoint) getMenuCodeLicenseConfig() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterLicensePreLicenseRedesign, constanta.MenuUserMasterLicensePreLicense)
}

func (input licenseConfigEndpoint) LicenseConfigWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LicenseConfigWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.InsertDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.InsertLicenseConfig)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.ViewDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.GetListLicenseConfig)
	}
}

func (input licenseConfigEndpoint) LicenseConfigWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "LicenseConfigWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.DeleteDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.DeleteLicenseConfig)
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.UpdateDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.UpdateLicenseConfig)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.ViewDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.ViewDetailLicenseConfigService)
	}
}

func (input licenseConfigEndpoint) InitiateLicenseConfig(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateLicenseConfig"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.ViewDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.InitiateGetListLicenseConfig)
	}
}

func (input licenseConfigEndpoint) SelectAllLicenseConfig(response http.ResponseWriter, request *http.Request) {
	funcName := "SelectAllLicenseConfig"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.ViewDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.SelectAllLicenseConfig)
	}
}

func (input licenseConfigEndpoint) InsertMultipleLicenseConfig(response http.ResponseWriter, request *http.Request) {
	funcName := "InsertMultipleLicenseConfig"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeLicenseConfig()+common.InsertDataPermissionMustHave, response, request, LicenseConfigService.LicenseConfigService.InsertMultipleLicenseConfig)
	}
}
