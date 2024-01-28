package Task

import (
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/service/TaskSchedulerService"
)

type checkExpirationProductLicense struct {
	AbstractScheduledTask
}

var CheckExpirationProductLicense = checkExpirationProductLicense{}.New()

func (input checkExpirationProductLicense) New() (output checkExpirationProductLicense) {
	output.RunType = "scheduler.check_expired_product_license"
	return
}

func (input checkExpirationProductLicense) Start() {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.StartTask(input.RunType, input.checkExpiredProductLicenseTask)
	}
}

func (input checkExpirationProductLicense) StartMain() {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.checkExpiredProductLicenseTask()
	}
}

func (input checkExpirationProductLicense) checkExpiredProductLicenseTask() {
	//--- Check Expired Product License
	TaskSchedulerService.TaskSchedulerService.SchedulerCheckExpiredProductLicense()
}
