package TodaysLeaveService

import "nexsoft.co.id/nextrac2/service"

type todaysLeaveService struct {
	service.AbstractService
	service.GetListData
}

var TodaysLeave = todaysLeaveService{}.New()

func (input todaysLeaveService) New() (output todaysLeaveService) {
	output.ServiceName = "Todays Leave"
	output.FileName = "TodaysLeaveService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{
		"id_card",
		"name",
		"department",
		"type",
	}
	output.ValidOrderBy = []string{
		"id_card",
		"name",
		"department",
	}
	return
}
