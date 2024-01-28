package ProjectService

import (
	"nexsoft.co.id/nextrac2/service"
)

type projectService struct {
	service.AbstractService
	service.GetListData
}

var ProjectService = projectService{}.New()

func (input projectService) New() (output projectService) {
	output.FileName = "ProjectService.go"
	output.ServiceName = "PROJECT"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"name"}
	output.ValidSearchBy = []string{"name"}
	return
}
