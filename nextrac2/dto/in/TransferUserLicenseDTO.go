package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type TransferUserLicenseRequest struct {
	AbstractDTO
	ID              int64     `json:"id"`
	InstallationID  int64     `json:"installation_id"`
	NoOfUser        int64     `json:"no_of_user"`
	UpdatedAtStr    string    `json:"updated_at"`
	UpdatedAt       time.Time
}

func (input *TransferUserLicenseRequest) ValidateTransferredUser() (err errorModel.ErrorModel) {
	if err = input.validateForUpdate("TransferUserLicenseDTO.go", "ValidateTransferredUser"); err.Error != nil {
		return
	}

	return input.mandatoryValidation("TransferUserLicenseDTO.go", "ValidateTransferredUser")
}

func (input *TransferUserLicenseRequest) validateForUpdate(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return
}

func (input *TransferUserLicenseRequest) mandatoryValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.InstallationID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.InstallationID)
		return
	}

	if input.NoOfUser < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.NumberOfUser)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}