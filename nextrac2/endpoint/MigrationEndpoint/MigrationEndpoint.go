package MigrationEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/MigrationService"
)

type migrationEndpoint struct {
	endpoint.AbstractEndpoint
}

var MigrationEndpoint = migrationEndpoint{}.New()

func (input migrationEndpoint) New() (output migrationEndpoint) {
	output.FileName = "MigrationEndpoint.go"
	return
}

func (input migrationEndpoint) MigrationWithoutParam(response http.ResponseWriter, request *http.Request) {
	funcName := "MigrationWithoutParam"
	switch request.Method {
	case "POST":
		input.ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, "", response, request, MigrationService.MigrationService.ResetMigration)
	}
}
