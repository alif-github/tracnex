package ProductLicenseEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ProductLicenseService"
)

type productLicenseEndpoint struct {
	endpoint.AbstractEndpoint
}

var ProductLicenseEndpoint = productLicenseEndpoint{}.New()

func (input productLicenseEndpoint) New() (output productLicenseEndpoint) {
	output.FileName = "ProductLicenseEndpoint.go"
	return
}

func (input productLicenseEndpoint) getMenuCodeProductLicense() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterProdukLisensiRedesign, constanta.MenuUserMasterProdukLisensi)
}

func (input productLicenseEndpoint) ProductLicenseWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ProductLicenseWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeProductLicense()+common.ViewDataPermissionMustHave, response, request, ProductLicenseService.ProductLicenseService.GetListProductLicense)
	}
}

func (input productLicenseEndpoint) InitiateProductLicenseEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateProductLicenseEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeProductLicense()+common.ViewDataPermissionMustHave, response, request, ProductLicenseService.ProductLicenseService.InitiateGetListProductLicense)
		break
	}
}

func (input productLicenseEndpoint) DetailProductLicenseWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DetailProductLicenseWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeProductLicense()+common.ViewDataPermissionMustHave, response, request, ProductLicenseService.ProductLicenseService.ViewProductLicense)
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeProductLicense()+common.UpdateDataPermissionMustHave, response, request, ProductLicenseService.ProductLicenseService.UpdateProductLicense)
	}
}

func (input productLicenseEndpoint) DecryptProductLicenseEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DecryptProductLicenseEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeProductLicense()+common.ViewDataPermissionMustHave, response, request, ProductLicenseService.ProductLicenseService.DecryptProductLicense)
		break
	}
}

func (input productLicenseEndpoint) UpdateProductLicenseHWIDEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UpdateProductLicenseHWIDEndpoint"
	switch request.Method {
	case http.MethodPut:
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuActivatedNexmileUser+common.UpdateDataPermissionMustHave, response, request, ProductLicenseService.ProductLicenseService.UpdateHWIDProductLicense)
		break
	}
}

func (input productLicenseEndpoint) UpdateProductLicenseHWIDByPassEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UpdateProductLicenseHWIDByPassEndpoint"
	switch request.Method {
	case http.MethodPut:
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, ProductLicenseService.ProductLicenseService.UpdateHWIDByPassProductLicense)
		break
	}
}
