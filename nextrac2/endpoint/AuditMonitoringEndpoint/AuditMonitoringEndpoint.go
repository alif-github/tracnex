package AuditMonitoringEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/AuditMonitoringService"
	"strings"
)

type auditMonitoringEndpoint struct {
	endpoint.AbstractEndpoint
}

var AuditMonitoringEndpoint = auditMonitoringEndpoint{}.New()

func (input auditMonitoringEndpoint) New() (output auditMonitoringEndpoint) {
	output.FileName = "AddResourceExternalEndpoint.go"
	return
}

func (input auditMonitoringEndpoint) AuditMonitoringEndpointWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "AuditMonitoringEndpointWithParam"
	splitPath := strings.Split(request.URL.Path, "/")
	isAdmin := false
	for _, s := range splitPath {
		if s == "admin" {
			isAdmin = !isAdmin
			break
		}
	}

	if isAdmin {
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuNexsoftAuditMonitoring+common.ViewDataPermissionMustHave, response, request, AuditMonitoringService.AuditMonitoringService.GetListAuditMonitoring)
	} else {
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.ReadDataAPIMustHave, constanta.MenuAuditMonitoringUser+common.ViewDataPermissionMustHave, response, request, AuditMonitoringService.AuditMonitoringService.GetListAuditMonitoring)
	}
}
