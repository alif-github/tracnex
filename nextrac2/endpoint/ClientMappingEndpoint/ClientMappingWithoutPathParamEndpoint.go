package ClientMappingEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ClientMappingService"
)

func (input clientMappingEndpoint) ClientMappingChangeSocketIDEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "UIClientMappingChangeSocketIDEndPoint"
	switch request.Method {
	case "PUT" :
		input.ServeJWTTokenValidationEndpoint(funcName,false, common.WriteDataAPIMustHave, constanta.MenuUpdateSocketID + common.UpdateDataPermissionMustHave, response, request, ClientMappingService.UpdateSocketIDService.UpdateSocketID)
		break
	}
}