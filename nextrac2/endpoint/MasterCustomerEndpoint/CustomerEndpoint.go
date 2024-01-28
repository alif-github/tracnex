package MasterCustomerEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/CustomerService"
)

type masterCustomerEndpoint struct {
	endpoint.AbstractEndpoint
}

var MasterCustomerEndpoint = masterCustomerEndpoint{}.New()

func (input masterCustomerEndpoint) New() (output masterCustomerEndpoint) {
	output.FileName = "MasterCustomerEndpoint.go"
	return
}

func (input masterCustomerEndpoint) getMenuCodeSubCustomer() string {
	return endpoint.GetMenuCode(constanta.MenuCustomerSubCustomerRedesign, constanta.MenuCustomerSubCustomer)
}

func (input masterCustomerEndpoint) CustomerEndpointWithoutPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerEndpointWithoutPathParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubCustomer()+common.InsertDataPermissionMustHave, response, request, CustomerService.CustomerService.InsertCustomer)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSubCustomer()+common.ViewDataPermissionMustHave, response, request, CustomerService.CustomerService.GetListCustomer)
		break
	}
}

func (input masterCustomerEndpoint) CustomerEndpointWithPathParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CustomerEndpointWithPathParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSubCustomer()+common.ViewDataPermissionMustHave, response, request, CustomerService.CustomerService.ViewCustomer)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubCustomer()+common.UpdateDataPermissionMustHave, response, request, CustomerService.CustomerService.UpdateCustomer)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubCustomer()+common.DeleteDataPermissionMustHave, response, request, CustomerService.CustomerService.DeleteCustomer)
		break
	}
}

func (input masterCustomerEndpoint) InitiateCustomerEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateCustomerEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubCustomer()+common.ViewDataPermissionMustHave, response, request, CustomerService.CustomerService.InitiateCustomer)
		break
	}
}

func (input masterCustomerEndpoint) GetListCustomerNonParentEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "GetListCustomerNonParentEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubCustomer()+common.ViewDataPermissionMustHave, response, request, CustomerService.CustomerService.GetListCustomerNonParent)
		break
	}
}

func (input masterCustomerEndpoint) InitiateCustomerNonParentEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateCustomerNonParentEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubCustomer()+common.ViewDataPermissionMustHave, response, request, CustomerService.CustomerService.InitiateCustomerNonParent)
		break
	}
}

func (input masterCustomerEndpoint) GetListCustomerParentEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "GetListCustomerParentEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSubCustomer()+common.ViewDataPermissionMustHave, response, request, CustomerService.CustomerService.GetListCustomerParent)
		break
	}
}

func (input masterCustomerEndpoint) InitiateCustomerParentEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateCustomerParentEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSubCustomer()+common.ViewDataPermissionMustHave, response, request, CustomerService.CustomerService.InitiateCustomerParent)
		break
	}
}

func (input masterCustomerEndpoint) InternalGetListCustomerEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InternalGetListCustomerEndpoint"
	switch request.Method {
	case "POST":
		input.ServeInternalValidationEndpoint(funcName, false, true, response, request, CustomerService.CustomerService.InternalGetListCustomer)
		break
	}
}

func (input masterCustomerEndpoint) InternalGetListDistributorEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InternalGetListDistributorEndpoint"
	switch request.Method {
	case "POST":
		input.ServeInternalValidationEndpoint(funcName, false, true, response, request, CustomerService.CustomerService.InternalGetListDistributor)
		break
	}
}

func (input masterCustomerEndpoint) InternalCountCustomerEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InternalGetListCustomerEndpoint"
	switch request.Method {
	case "POST":
		input.ServeInternalValidationEndpoint(funcName, false, true, response, request, CustomerService.CustomerService.InternalCountCustomer)
		break
	}
}
