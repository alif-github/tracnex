package ValidationLicenseEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/service/ValidationLicenseService"
)

type validationLicenseEndpoint struct {
	endpoint.AbstractEndpoint
}

var ValidationLicenseEndpoint = validationLicenseEndpoint{}.New()

func (input validationLicenseEndpoint) New() (output validationLicenseEndpoint) {
	output.FileName = "ValidationLicenseEndpoint.go"
	return
}

func (input validationLicenseEndpoint) ValidateLicenseEndpoint(response http.ResponseWriter, request *http.Request)  {
	funcName := "ValidateLicenseEndpoint"
	switch request.Method {
	case "PUT":
		input.ServeJWTTokenValidationEndpoint(funcName, false, common.WriteDataAPIMustHave, constanta.MenuUserClientMappingBackEndConstanta + common.UpdateDataPermissionMustHave, response, request, ValidationLicenseService.ValidationLicenseService.ValidateLicense)
		break
	}
}
