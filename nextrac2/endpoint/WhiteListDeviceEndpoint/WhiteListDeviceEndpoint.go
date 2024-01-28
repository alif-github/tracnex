package WhiteListDeviceEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/WhitelistDeviceService"
)

type whiteListDevice struct {
	endpoint.AbstractEndpoint
}

var WhiteListDeviceEndpoint = whiteListDevice{}.New()

func (input whiteListDevice) New() (output whiteListDevice) {
	output.FileName = "WhiteListDevice.go"
	return
}

func (input whiteListDevice) getMenuCodeWhiteListDevice() string {
	return endpoint.GetMenuCode(constanta.MenuWhiteListDeviceRedesign, constanta.MenuWhiteListDevice)
}

func (input whiteListDevice) WhiteListDeviceEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "WhiteListDeviceEndpointWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeWhiteListDevice()+common.InsertDataPermissionMustHave, response, request, WhitelistDeviceService.WhitelistDeviceService.InsertWhiteListDevice)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeWhiteListDevice()+common.ViewDataPermissionMustHave, response, request, WhitelistDeviceService.WhitelistDeviceService.GetListWhiteListDevice)
		break
	}
}

func (input whiteListDevice) WhiteListDeviceEndpointWithParam(response http.ResponseWriter, request *http.Request) {
	funcName := "WhiteListDeviceEndpointWithoutParam"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeWhiteListDevice()+common.UpdateDataPermissionMustHave, response, request, WhitelistDeviceService.WhitelistDeviceService.UpdateWhiteListDevice)
		break
	case "DELETE":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, input.getMenuCodeWhiteListDevice()+common.DeleteDataPermissionMustHave, response, request, WhitelistDeviceService.WhitelistDeviceService.DeleteWhiteListDevice)
		break
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeWhiteListDevice()+common.ViewDataPermissionMustHave, response, request, WhitelistDeviceService.WhitelistDeviceService.ViewWhiteListDevice)
		break
	}
}

func (input whiteListDevice) InitiateGetListWhiteListDeviceEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateGetListWhiteListDeviceEndpoint"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, input.getMenuCodeWhiteListDevice()+common.ViewDataPermissionMustHave, response, request, WhitelistDeviceService.WhitelistDeviceService.InitiateWhiteListDevice)
		break
	}
}