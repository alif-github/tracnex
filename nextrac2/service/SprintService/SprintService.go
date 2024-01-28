package SprintService

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/service"
)

type sprintService struct {
	service.AbstractService
	service.GetListData
}

var SprintService = sprintService{}.New()

func (input sprintService) New() (output sprintService) {
	output.FileName = "SprintService.go"
	output.ServiceName = constanta.Sprint
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"order"}
	return
}
