package BacklogEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/BacklogService"
	"nexsoft.co.id/nextrac2/service/DepartmentService"
	"nexsoft.co.id/nextrac2/service/EmployeeService"
	"nexsoft.co.id/nextrac2/service/ProjectService"
	"nexsoft.co.id/nextrac2/service/SprintService"
	"nexsoft.co.id/nextrac2/service/TrackerService"
)

type backlogEndpoint struct {
	endpoint.AbstractEndpoint
}

var BacklogEndpoint = backlogEndpoint{}.New()

func (input backlogEndpoint) New() (output backlogEndpoint) {
	output.FileName = "BacklogEndpoint.go"
	return
}

func (input backlogEndpoint) getMenuCodeBacklog() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterTimesheetBacklogRedesign, constanta.MenuUserMasterTimesheetBacklog)
}

func (input backlogEndpoint) BacklogEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "BacklogEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, BacklogService.BacklogService.GetListParentBacklog)
		break
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.InsertData, response, request, BacklogService.BacklogService.InsertDetailBacklog)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeBacklog()+common.DeleteDataPermissionMustHave, response, request, BacklogService.BacklogService.DeleteBacklogBySprint)
		break
	}
}

func (input backlogEndpoint) InitiateBacklogEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateBacklogEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, BacklogService.BacklogService.InitiateGetListParentBacklog)
		break
	}
}

func (input backlogEndpoint) ImportFileBacklogEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "ImportFileBacklogEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, BacklogService.BacklogService.UnmarshalFileBacklog)
		break
	}
}

func (input backlogEndpoint) BacklogDetailEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "BacklogDetailEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, BacklogService.BacklogService.GetListDetailBacklog)
		break
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeBacklog()+common.InsertDataPermissionMustHave, response, request, BacklogService.BacklogService.InsertDetailBacklog)
		break
	}
}

func (input backlogEndpoint) UpdateStatusBacklogEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "UpdateStatusBacklogEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, BacklogService.BacklogService.MultipleUpdateStatusBacklog)
		break
	}
}

func (input backlogEndpoint) InitiateBacklogDetailEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateBacklogDetailEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, BacklogService.BacklogService.InitiateGetListDetailParentBacklog)
		break
	}
}

func (input backlogEndpoint) StatusBacklogDetailEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "StatusBacklogDetailEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, BacklogService.BacklogService.GetListStatusBacklog)
		break
	}
}

func (input backlogEndpoint) BacklogDetailEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "BacklogDetailEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, BacklogService.BacklogService.ViewDetailBacklog)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeBacklog()+common.DeleteDataPermissionMustHave, response, request, BacklogService.BacklogService.DeleteDetailBacklog)
		break
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeBacklog()+common.UpdateDataPermissionMustHave, response, request, BacklogService.BacklogService.UpdateDetailBacklog)
		break
	}
}

func (input backlogEndpoint) DropDownListSprintEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DropDownListSprintEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, SprintService.SprintService.GetListSprint)
		break
	}
}

func (input backlogEndpoint) DropDownSearchProjectEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DropDownSearchProjectEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, ProjectService.ProjectService.GetListProjectByRedmineAPI)
		break
	}
}

func (input backlogEndpoint) DropDownListTrackerEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DropDownListTrackerEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, TrackerService.TrackerService.GetListTracker)
		break
	}
}

func (input backlogEndpoint) DropDownListEmployeeEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DropDownListEmployeeEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, EmployeeService.EmployeeService.GetListEmployeeDDL)
		break
	}
}

func (input backlogEndpoint) DropDownListDepartmentEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "DropDownListDepartmentEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeBacklog()+common.ViewDataPermissionMustHave, response, request, DepartmentService.DepartmentService.GetListDepartment)
		break
	}
}
