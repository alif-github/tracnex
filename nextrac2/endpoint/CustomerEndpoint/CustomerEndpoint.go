package CustomerEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/CustomerListService"
)

type customerEndpoint struct {
	endpoint.AbstractEndpoint
}

var CustomerEndpoint = customerEndpoint{}.New()

func (input customerEndpoint) New() (output customerEndpoint) {
	output.FileName = "CustomerEndpoint.go"
	return
}

func (input customerEndpoint) CustomerEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerEndpointWithoutParam"

	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftCustomerConstanta+common.InsertDataPermissionMustHave, response, request, CustomerListService.CustomerListService.InsertCustomer)
		break
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftCustomerConstanta+common.ViewDataPermissionMustHave, response, request, CustomerListService.CustomerListService.GetListCustomerList)
	}
}

func (input customerEndpoint) CustomerEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerEndpointWithParam"

	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftCustomerConstanta+common.ViewDataPermissionMustHave, response, request, CustomerListService.CustomerListService.ViewCustomer)
	}
}

func (input customerEndpoint) InitiateCustomer(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateCustomer"
	input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftCustomerConstanta+common.ViewDataPermissionMustHave, response, request, CustomerListService.CustomerListService.InitiateGetListCustomerList)
}