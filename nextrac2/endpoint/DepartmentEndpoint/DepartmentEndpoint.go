package DepartmentEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/DepartmentService"
)

type departmentEndpoint struct {
	endpoint.AbstractEndpoint
}

var DepartmentEndpoint = departmentEndpoint{}.New()

func (input departmentEndpoint) New() (output departmentEndpoint) {
	output.FileName = "DepartmentEndpoint.go"
	return
}

func (input departmentEndpoint) DepartmentEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DepartmentEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, DepartmentService.DepartmentService.GetListDepartment)
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, DepartmentService.DepartmentService.InsertDepartment)
	}
}

func (input departmentEndpoint) InitiateDepartmentEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DepartmentEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, DepartmentService.DepartmentService.InitiateDepartment)
	}
}

func (input departmentEndpoint) DepartmentEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "DepartmentEndpointWithParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, "", response, request, DepartmentService.DepartmentService.ViewDepartment)
	//case "DELETE":
	//	input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, DepartmentService.DepartmentService.DeleteDepartment)
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, DepartmentService.DepartmentService.UpdateDepartment)
	}
}
