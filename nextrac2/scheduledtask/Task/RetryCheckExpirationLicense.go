package Task

import (
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/service/TaskSchedulerService"
)

type retryExpirationLicense struct {
	AbstractScheduledTask
}

var RetryExpirationLicense = retryExpirationLicense{}

func (input retryExpirationLicense) New() (output retryExpirationLicense) {
	output.RunType = "scheduler.retry_check_expired_product_license"
	return
}

func (input retryExpirationLicense) Start()  {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.StartTask(input.RunType, input.retryExpirationLicenseTask)
	}
}

func (input retryExpirationLicense) StartMain()  {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.retryExpirationLicenseTask()
	}
}

func (input retryExpirationLicense) retryExpirationLicenseTask()  {
	TaskSchedulerService.TaskSchedulerService.SchedulerRetryValidateExpirationLicense()
}