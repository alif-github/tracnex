package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type UserLicenseRequest struct {
	AbstractDTO
	ID int64 `json:"id"`

}

type ViewUserLicenseRequest struct {
	UserLicenseId int64 `json:"user_license_id"`
}

func (input UserLicenseRequest) ValidateViewDetailUserLicense() (err errorModel.ErrorModel) {

	err = input.ValidateZeroIDUserLicense()
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input UserLicenseRequest) ValidateZeroIDUserLicense() errorModel.ErrorModel {
	fileName := "UserLicenseDTO.go"
	funcName := "ValidateZeroID"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input ViewUserLicenseRequest) ValidateViewTransferKeyUserLicense() (err errorModel.ErrorModel) {
	fileName := "UserLicenseDTO.go"
	funcName := "ValidateViewTransferKeyUserLicense"

	if input.UserLicenseId < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.UserLicenseID)
	}

	return errorModel.GenerateNonErrorModel()
}