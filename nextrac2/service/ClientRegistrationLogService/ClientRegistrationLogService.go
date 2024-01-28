package ClientRegistrationLogService

import "nexsoft.co.id/nextrac2/service"

type clientRegistrationLogService struct {
	service.AbstractService
	service.GetListData
}

var ClientRegistrationLogService = clientRegistrationLogService{}.New()

func (input clientRegistrationLogService) New() (output clientRegistrationLogService) {
	output.FileName = "ClientRegistrationLogService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id"}
	output.ValidSearchBy = []string{"success_status_auth", "success_status_nexcloud", "resource"}
	return
}