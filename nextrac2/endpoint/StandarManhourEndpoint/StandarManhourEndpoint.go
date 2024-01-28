package StandarManhourEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/StandarManhourService"
)

type standarManhourEndpoint struct {
	endpoint.AbstractEndpoint
}

var StandarManhourEndpoint = standarManhourEndpoint{}.New()

func (input standarManhourEndpoint) New() (output standarManhourEndpoint) {
	output.FileName = "StandarManhourEndpoint.go"
	return
}

func (input standarManhourEndpoint) getMenuCodeEmployee() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterTimesheetEmployeeRedesign, constanta.MenuUserMasterTimesheetEmployee)
}

func (input standarManhourEndpoint) StandarManhourWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "StandarManhourWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.InsertDataPermissionMustHave, response, request, StandarManhourService.StandarManhourService.InsertStandarManhour)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, StandarManhourService.StandarManhourService.GetListStandarManhour)
	}
}

func (input standarManhourEndpoint) StandarManhourWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "StandarManhourWithParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.UpdateDataPermissionMustHave, response, request, StandarManhourService.StandarManhourService.UpdateStandarManhour)
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeEmployee()+common.DeleteDataPermissionMustHave, response, request, StandarManhourService.StandarManhourService.DeleteStandarManhour)
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, StandarManhourService.StandarManhourService.ViewStandarManhour)
	}
}

func (input standarManhourEndpoint) InitiateStandarManhour(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateStandarManhour"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeEmployee()+common.ViewDataPermissionMustHave, response, request, StandarManhourService.StandarManhourService.InitiateStandarManhour)
	}
}
