package ImportFileCustomerListEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ProcessFileListCustomerService"
)

type importFileCustomerListEndpoint struct {
	endpoint.AbstractEndpoint
}

var ImportFileCustomerListEndpoint = importFileCustomerListEndpoint{}.New()

func (input importFileCustomerListEndpoint) New() (output importFileCustomerListEndpoint){
	output.FileName = "ImportFileCustomerListEndpoint.go"
	return
}

func (input importFileCustomerListEndpoint) ImportFileCustomerListEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UserEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftCustomerConstanta + common.InsertDataPermissionMustHave, response, request, ProcessFileListCustomerService.ProcessFileListCustomerService.ImportFileCustomerList)
	}
}

func (input importFileCustomerListEndpoint) InitiateImportEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateImportEndpoint"
	input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftImportConstanta + common.ViewDataPermissionMustHave, response, request, ProcessFileListCustomerService.ImportService.InitiateImportData)
}

func (input importFileCustomerListEndpoint) ImportAndValidateDataEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ImportAndValidateDataEndpoint"
	input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftImportConstanta + common.InsertDataPermissionMustHave, response, request, ProcessFileListCustomerService.ImportService.ImportAndValidateData)
}

func (input importFileCustomerListEndpoint) ConfirmImportEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ConfirmImportEndpoint"
	input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuNexsoftImportConstanta + common.InsertDataPermissionMustHave, response, request, ProcessFileListCustomerService.ImportService.ConfirmImport)
}