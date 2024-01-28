package ProductEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ProductService"
)

type productEndpoint struct {
	endpoint.AbstractEndpoint
}

var ProductEndpoint = productEndpoint{}.New()

func (input productEndpoint) New() (output productEndpoint) {
	output.FileName = "ProductEndpoint.go"
	return
}

func (input productEndpoint) getMenuCodeSubProduct() string {
	return endpoint.GetMenuCode(constanta.MenuSubProductRedesign, constanta.MenuSubProduct)
}

func (input productEndpoint) ProductWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ProductWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubProduct()+common.InsertDataPermissionMustHave, responseWriter, request, ProductService.ProductService.InsertProduct)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSubProduct()+common.ViewDataPermissionMustHave, responseWriter, request, ProductService.ProductService.GetListProduct)
	}
}

func (input productEndpoint) ProductWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ProductWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSubProduct()+common.ViewDataPermissionMustHave, responseWriter, request, ProductService.ProductService.ViewProduct)
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubProduct()+common.DeleteDataPermissionMustHave, responseWriter, request, ProductService.ProductService.DeleteProduct)
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubProduct()+common.UpdateDataPermissionMustHave, responseWriter, request, ProductService.ProductService.UpdateProduct)
	}
}

func (input productEndpoint) InitiateProduct(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "InitiateProduct"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSubProduct()+common.ViewDataPermissionMustHave, responseWriter, request, ProductService.ProductService.InitiateGetListProduct)
	}
}
