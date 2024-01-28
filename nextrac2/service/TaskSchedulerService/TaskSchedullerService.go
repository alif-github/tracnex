package TaskSchedulerService

import "nexsoft.co.id/nextrac2/service"

type taskSchedulerService struct {
	service.AbstractService
}

var TaskSchedulerService = taskSchedulerService{}.New()

func (input taskSchedulerService) New() (output taskSchedulerService) {
	output.FileName = "TaskSchedulerService.go"
	return
}


