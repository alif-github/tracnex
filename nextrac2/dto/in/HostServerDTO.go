package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type HostServerRequest struct {
	AbstractDTO
	ID           int64  `json:"id"`
	HostName     string `json:"host_name"`
	HostUrl      string `json:"host_url"`
	UpdatedAt     time.Time
	UpdatedAtStr string `json:"updated_at"`
}

func (input *HostServerRequest) ValidateInsertHostServer() (err errorModel.ErrorModel) {
	fileName := "HostServerDTO.go"
	funcName := "ValidateInsertHostServer"

	isMatch, errField := util.IsNexsoftProfileNameStandardValid(input.HostName)
	if !isMatch {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Name, "")
		if err.Error != nil {
			return
		}
	}

	err = input.ValidateMinMaxString(input.HostName, constanta.Name, 5, 64)
	if err.Error != nil {
		return
	}

	err = input.ValidateMinMaxString(input.HostUrl, constanta.HostUrl, 0, 256)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *HostServerRequest) ValidateUpdateHostServer() (err errorModel.ErrorModel) {
	fileName := "HostServerDTO.go"
	funcName := "ValidateDeleteHostServer"
	if input.ID < 0 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	err = input.ValidateMinMaxString(input.HostUrl, constanta.HostUrl, 0, 256)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()

}

