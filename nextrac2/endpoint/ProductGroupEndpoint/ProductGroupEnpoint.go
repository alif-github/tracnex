package ProductGroupEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ProductGroupService"
)

type productGroupEndpoint struct {
	endpoint.AbstractEndpoint
}

var ProductGroupEndpoint = productGroupEndpoint{}.New()

func (input productGroupEndpoint) getMenuCodeProductGroup() string {
	return endpoint.GetMenuCode(constanta.MenuProductGroupRedesign, constanta.MenuProductGroup)
}

func (input productGroupEndpoint) New() (output productGroupEndpoint) {
	output.FileName = "ProductEndpoint.go"
	return
}

func (input productGroupEndpoint) ProductGroupEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ProductGroupEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeProductGroup()+common.ViewDataPermissionMustHave, response, request, ProductGroupService.ProductGroupService.GetListProductGroupByUser)
		break
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeProductGroup()+common.InsertDataPermissionMustHave, response, request, ProductGroupService.ProductGroupService.InsertProductGroup)
		break
	}
}

func (input productGroupEndpoint) ProductGroupEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ProductGroupEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeProductGroup()+common.ViewDataPermissionMustHave, response, request, ProductGroupService.ProductGroupService.ViewProductGroup)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeProductGroup()+common.UpdateDataPermissionMustHave, response, request, ProductGroupService.ProductGroupService.UpdateProductGroup)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeProductGroup()+common.DeleteDataPermissionMustHave, response, request, ProductGroupService.ProductGroupService.DeleteProductGroup)
		break
	}
}

func (input productGroupEndpoint) InitiateProductGroupEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateProductGroupEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeProductGroup()+common.ViewDataPermissionMustHave, response, request, ProductGroupService.ProductGroupService.GetInitiateProductGroup)
		break
	}
}

func (input productGroupEndpoint) ProductGroupAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ProductGroupAdminEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftProductGroup+common.ViewDataPermissionMustHave, response, request, ProductGroupService.ProductGroupService.GetListProductGroupByAdmin)
		break
	}
}
