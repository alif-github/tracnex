package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type OTPGenerator struct {
	GrochatAccount string           `json:"grochat_account"`
	DataDetail     DetailDataBundle `json:"data"`
}

type DetailDataBundle struct {
	ClientID     string `json:"client_id"`
	HWID         string `json:"hwid"`
	TimestampStr string `json:"timestamp"`
	Timestamp    time.Time
}

func (input *OTPGenerator) ValidateOTPGenerator() (err errorModel.ErrorModel) {
	var (
		fileName = "OtpGeneratorDTO.go"
		funcName = "ValidateOTPGenerator"
	)

	if util.IsStringEmpty(input.GrochatAccount) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "Grochat")
		return
	}

	if util.IsStringEmpty(input.DataDetail.ClientID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientID)
		return
	}

	if util.IsStringEmpty(input.DataDetail.HWID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.HWID)
		return
	}

	if util.IsStringEmpty(input.DataDetail.TimestampStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Timestamp")
	}

	input.DataDetail.Timestamp, err = TimeStrToTime(input.DataDetail.TimestampStr, "Timestamp")
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}
