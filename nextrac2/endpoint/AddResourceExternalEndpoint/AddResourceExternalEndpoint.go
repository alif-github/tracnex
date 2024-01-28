package AddResourceExternalEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/HitAddExternalResource"
)

type addResourceExternal struct {
	endpoint.AbstractEndpoint
}

var AddResourceExternalEndpoint = addResourceExternal{}.New()

func (input addResourceExternal) New() (output addResourceExternal) {
	output.FileName = "AddResourceExternalEndpoint.go"
	return
}

func (input addResourceExternal) AddResourceNexcloudEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "AddResourceNexcloudEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuAddResourceNexcloudConstanta+common.InsertDataPermissionMustHave, responseWriter, request, HitAddExternalResource.HitAddExternalResource.InsertAddResourceNexcloudService)
}

func (input addResourceExternal) AddResourceNexdriveEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "AddResourceNexdriveEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuAddResourceNexdriveConstanta+common.InsertDataPermissionMustHave, responseWriter, request, HitAddExternalResource.HitAddExternalResource.InsertAddResourceNexdriveService)
}

func (input addResourceExternal) ViewLogForAddResourceExternalEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "ViewLogForAddResourceExternalEndpoint"
	input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuAddResourceNexcloudConstanta+common.ViewDataPermissionMustHave, responseWriter, request, HitAddExternalResource.HitAddExternalResource.ViewLogForAddResource)
}