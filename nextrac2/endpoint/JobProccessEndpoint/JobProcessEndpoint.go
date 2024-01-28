package JobProccessEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/JobProcessService"
	"nexsoft.co.id/nextrac2/service/TaskSchedulerService"
)

type jobProcessEndpoint struct {
	endpoint.AbstractEndpoint
}

var JobProcessEndpoint = jobProcessEndpoint{}.New()

func (input jobProcessEndpoint) New() (output jobProcessEndpoint) {
	output.FileName = "JobProcessEndpoint.go"
	return
}

func (input jobProcessEndpoint) JobProcessEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "JobProcessEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftJobProcessConstanta+common.ViewDataPermissionMustHave, response, request, JobProcessService.JobProcessService.GetListJobProcess)
	}
}

func (input jobProcessEndpoint) JobProcessEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "JobProcessEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftJobProcessConstanta+common.ViewDataPermissionMustHave, response, request, JobProcessService.JobProcessService.ViewJobProcess)
	}
}

func (input jobProcessEndpoint) InitiateGetListJobProcess(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetListJobProcess"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftJobProcessConstanta+common.ViewDataPermissionMustHave, response, request, JobProcessService.JobProcessService.InitiateGetListJobProcess)
	}
}

func (input jobProcessEndpoint) SynchronizeRegionalDataEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "SynchronizeRegionalDataEndpoint"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, TaskSchedulerService.TaskSchedulerService.APISynchronizeRegionalData)
	}
}
