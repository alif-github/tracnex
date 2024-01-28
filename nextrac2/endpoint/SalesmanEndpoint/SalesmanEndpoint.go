package SalesmanEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/SalesmanService"
)

type salesmanEndpoint struct {
	endpoint.AbstractEndpoint
}

var SalesmanEndpoint = salesmanEndpoint{}.New()

func (input salesmanEndpoint) New() (output salesmanEndpoint) {
	output.FileName = "SalesmanEndpoint.go"
	return
}

func (input salesmanEndpoint) getMenuCodeSalesman() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterSetupSalesmanRedesign, constanta.MenuUserMasterSetupSalesman)
}

func (input salesmanEndpoint) SalesmanWithoutParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "SalesmanWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSalesman()+common.InsertDataPermissionMustHave, responseWriter, request, SalesmanService.SalesmanService.InsertSalesman)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSalesman()+common.ViewDataPermissionMustHave, responseWriter, request, SalesmanService.SalesmanService.GetListSalesman)
	}
}

func (input salesmanEndpoint) SalesmanWithParam(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "SalesmanWithParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSalesman()+common.UpdateDataPermissionMustHave, responseWriter, request, SalesmanService.SalesmanService.UpdateSalesman)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSalesman()+common.ViewDataPermissionMustHave, responseWriter, request, SalesmanService.SalesmanService.ViewSalesman)
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeSalesman()+common.UpdateDataPermissionMustHave, responseWriter, request, SalesmanService.SalesmanService.DeleteSalesman)
	}
}

func (input salesmanEndpoint) InitiateSalesman(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "InitiateSalesman"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeSalesman()+common.ViewDataPermissionMustHave, responseWriter, request, SalesmanService.SalesmanService.InitiateGetListSalesman)
}

func (input salesmanEndpoint) SalesmanWithoutParamAdmin(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "SalesmanWithoutParamAdmin"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", responseWriter, request, SalesmanService.SalesmanService.GetListAdminSalesman)
	}
}
