package CustomerCategoryEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/CustomerCategoryService"
)

type customerCategoryEndpoint struct {
	endpoint.AbstractEndpoint
}

var CustomerCategoryEndpoint = customerCategoryEndpoint{}

func (input customerCategoryEndpoint) getMenuCodeCustCategory() string {
	return endpoint.GetMenuCode(constanta.MenuCustomerCategoryRedesign, constanta.MenuCustomerCategory)
}

func (input customerCategoryEndpoint) CustomerCategoryEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerCategoryEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustCategory()+common.InsertDataPermissionMustHave, response, request, CustomerCategoryService.CustomerCategoryService.InsertCustomerCategory)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustCategory()+common.ViewDataPermissionMustHave, response, request, CustomerCategoryService.CustomerCategoryService.GetListCustomerCategoryByUser)
		break
	}
}

func (input customerCategoryEndpoint) CustomerCategoryEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerCategoryEndpointWithPathParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustCategory()+common.UpdateDataPermissionMustHave, response, request, CustomerCategoryService.CustomerCategoryService.UpdateCustomerCategory)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustCategory()+common.DeleteDataPermissionMustHave, response, request, CustomerCategoryService.CustomerCategoryService.DeleteCustomerCategory)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustCategory()+common.ViewDataPermissionMustHave, response, request, CustomerCategoryService.CustomerCategoryService.ViewCustomerCategory)
		break
	}
}

func (input customerCategoryEndpoint) InitiateCustomerCategoryEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerCategoryEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustCategory()+common.ViewDataPermissionMustHave, response, request, CustomerCategoryService.CustomerCategoryService.InitiateCustomerCategory)
		break
	}
}

func (input customerCategoryEndpoint) CustomerCategoryAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerCategoryAdminEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftCustomerCategory+common.ViewDataPermissionMustHave, response, request, CustomerCategoryService.CustomerCategoryService.GetListCustomerCategoryByAdmin)
		break
	}
}
