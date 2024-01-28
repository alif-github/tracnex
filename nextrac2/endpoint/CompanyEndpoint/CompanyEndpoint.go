package CompanyEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/CompanyService"
)

type companyEndpoint struct {
	endpoint.AbstractEndpoint
}

var CompanyEndpoint = companyEndpoint{}.New()

func (input companyEndpoint) New() (output companyEndpoint) {
	output.FileName = "CompanyEndpoint.go"
	return
}

//func (input companyEndpoint) getMenuCodeEmployeeSetting() string {
//	return endpoint.GetMenuCode(constanta.MenuSettingHRIS, constanta.MenuSettingHRISRedesign)
//}

func (input companyEndpoint) CompanyEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CompanyEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, CompanyService.CompanyGetListService.GetCompanyList)
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, CompanyService.CompanyInsertService.InsertCompany)

	}
}

func (input companyEndpoint) CompanyEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "CompanyEndpointWithParam"
	switch request.Method {
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, CompanyService.CompanyDeleteService.DeleteCompany)
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, CompanyService.CompanyUpdateService.UpdateCompany)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, CompanyService.CompanyDetailService.DetailCompany)
	}
}

func (input companyEndpoint) InitiateCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateCompanyEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, CompanyService.CompanyGetListService.InitiateCompany)
	}
}
