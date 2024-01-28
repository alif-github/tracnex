package LicenseJournalEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/LicenseJournalService"
)

type licenseJournalEndpoint struct {
	endpoint.AbstractEndpoint
}

var LicenseJournalEndpoint = licenseJournalEndpoint{}.New()

func (input licenseJournalEndpoint) New() (output licenseJournalEndpoint) {
	output.FileName = "LicenseJournalEndpoint.go"
	return
}

func (input licenseJournalEndpoint) InitiateLicenseJournalEksternalEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "InitiateLicenseJournalEksternalEndpoint"
	switch request.Method {
	case "POST":
		input.ServeInternalValidationEndpoint(funcName, false, true, response, request, LicenseJournalService.LicenseJournalService.InitiateGetListLicenseJournal)
	}
}

func (input licenseJournalEndpoint) ListLicenseJournalEksternalEndpoint(response http.ResponseWriter, request *http.Request) {
	funcName := "ListLicenseJournalEksternalEndpoint"
	switch request.Method {
	case "POST":
		input.ServeInternalValidationEndpoint(funcName, false, true, response, request, LicenseJournalService.LicenseJournalService.GetListLicenseJournal)
	}
}
