package CustomerGroupEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/CustomerGroupService"
)

type customerGroupEndpoint struct {
	endpoint.AbstractEndpoint
}

var CustomerGroupEndpoint customerGroupEndpoint

func (input customerGroupEndpoint) getMenuCodeCustGroup() string {
	return endpoint.GetMenuCode(constanta.MenuCustomerGroupRedesign, constanta.MenuCustomerGroup)
}

func (input customerGroupEndpoint) CustomerGroupEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerGroupEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustGroup()+common.InsertDataPermissionMustHave, response, request, CustomerGroupService.CustomerGroupService.InsertCustomerGroup)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustGroup()+common.ViewDataPermissionMustHave, response, request, CustomerGroupService.CustomerGroupService.GetListCustomerGroupByUser)
		break
	}
}

func (input customerGroupEndpoint) InitiateCustomerGroupEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateCustomerGroupEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustGroup()+common.ViewDataPermissionMustHave, response, request, CustomerGroupService.CustomerGroupService.GetInitiateCustomerGroup)
		break
	}
}

func (input customerGroupEndpoint) CustomerGroupEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerGroupEndpointWithPathParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustGroup()+common.UpdateDataPermissionMustHave, response, request, CustomerGroupService.CustomerGroupService.UpdateCustomerGroup)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeCustGroup()+common.DeleteDataPermissionMustHave, response, request, CustomerGroupService.CustomerGroupService.DeleteCustomerGroup)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeCustGroup()+common.ViewDataPermissionMustHave, response, request, CustomerGroupService.CustomerGroupService.ViewCustomerGroup)
		break
	}
}

func (input customerGroupEndpoint) CustomerGroupAdminEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerGroupAdminEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftCustomerGroup+common.ViewDataPermissionMustHave, response, request, CustomerGroupService.CustomerGroupService.GetListCustomerGroupByAdmin)
		break
	}
}
