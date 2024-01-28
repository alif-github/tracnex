package PersonTitleEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MasterDataService/PersonTitleService"
)

type personTitleEndpoint struct {
	endpoint.AbstractEndpoint
}

var PersonTitleEndpoint = personTitleEndpoint{}.New()

func (input personTitleEndpoint) New() (output personTitleEndpoint) {
	output.FileName = "PersonTitleEndpoint.go"
	return
}

func (input personTitleEndpoint) PersonTitleEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "PersonTitleEndpointWithoutParam"
	switch request.Method {
	case "GET":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftRoot+common.ViewDataPermissionMustHave, response, request, PersonTitleService.PersonTitleService.GetListPersonTitle)
	}
}
