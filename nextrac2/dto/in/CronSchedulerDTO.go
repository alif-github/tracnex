package in

import (
	"github.com/robfig/cron"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type CronSchedulerRequest struct {
	AbstractDTO
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	RunType      string `json:"run_type"`
	Cron         string `json:"cron"`
	Status       bool   `json:"status"`
	UpdatedAt    time.Time
	UpdatedAtStr string `json:"updated_at"`
}

func (input *CronSchedulerRequest) ValidateViewCronScheduler() (err errorModel.ErrorModel) {
	fileName := "CronSchedulerDTO.go"
	funcName := "ValidateViewCronScheduler"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *CronSchedulerRequest) ValidateUpdateCronScheduler() (err errorModel.ErrorModel) {
	fileName := "CronSchedulerDTO.go"
	funcName := "ValidateUpdateCronScheduler"
	if input.ID < 0 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	c := cron.New()
	_, errs := c.AddFunc(input.Cron, func() {})
	if errs != nil {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Cron)
		return
	}

	return errorModel.GenerateNonErrorModel()

}
