package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type ProvinceRequest struct {
	AbstractDTO
	CountryID   int64  `json:"country_id"`
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	LastSyncStr string `json:"last_sync"`
	LastSync    time.Time
}

func (input *ProvinceRequest) ValidateReset() errorModel.ErrorModel {
	var errorS error
	if !util.IsStringEmpty(input.LastSyncStr) {
		input.LastSync, errorS = time.Parse(constanta.DefaultTimeFormat, input.LastSyncStr)
		if errorS != nil {
			return errorModel.GenerateFormatFieldError("AbstractDTO.go", "TimeStrToTime", "Last Sync")
		}
	}
	return errorModel.GenerateNonErrorModel()
}
