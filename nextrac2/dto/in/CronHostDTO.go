package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type CronHostRequest struct {
	AbstractDTO
	ID           int64  `json:"id"`
	HostID       int64  `json:"host_id"`
	HostName     string `json:"host_name"`
	UpdatedAt    time.Time
	UpdatedAtStr string   `json:"updated_at"`
	ListCron     ListCron `json:"list_cron"`
}

type ListCron []struct {
	CronID    int64 `json:"cron_id"`
	Active    bool  `json:"active"`
	RunStatus bool  `json:"run_status"`
}

func (input *CronHostRequest) ValidateEdit() (err errorModel.ErrorModel) {
	fileName := "CronHostDTO.go"
	funcName := "ValidateEdit"

	if input.HostID < 0 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.HostID)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	if len(input.ListCron) < 0 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.CronHost)
	}
	return errorModel.GenerateNonErrorModel()
}
