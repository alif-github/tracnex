package TrackerService

import (
	"nexsoft.co.id/nextrac2/service"
)

type trackerService struct {
	service.AbstractService
	service.GetListData
}

var TrackerService = trackerService{}.New()

func (input trackerService) New() (output trackerService) {
	output.FileName = "TrackerService.go"
	output.ServiceName = "TRACKER"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"order"}
	return
}
