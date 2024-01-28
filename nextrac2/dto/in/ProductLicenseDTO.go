package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type ProductLicenseRequest struct {
	AbstractDTO
	ID                     int64  `json:"id"`
	LicenseConfigID        int64  `json:"license_config_id"`
	ProductKey             string `json:"product_key"`
	ProductEncrypt         string `json:"product_encrypt"`
	ProductSignature       string `json:"product_signature"`
	ClientID               string `json:"client_id"`
	ClientSecret           string `json:"client_secret"`
	Hwid                   string `json:"hwid"`
	ActivationDateStr      string `json:"activation_date"`
	ActivationDate         time.Time
	LicenseStatus          int32  `json:"license_status"`
	TerminationDescription string `json:"termination_description"`
	UpdateAtStr            string `json:"update_at"`
	UpdatedAt              time.Time
}

func (input ProductLicenseRequest) ValidateViewProductLicense() (err errorModel.ErrorModel) {
	fileName := "ProductLicenseDTO.go"
	funcName := "ValidateViewProductLicense"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ProductLicenseID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ProductLicenseRequest) ValidateUpdateProductLicense() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateProductLicense"

	err = input.validateLicenseStatus(ProductLicenseDTOFilename, funcName)
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(ProductGroupDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ProductLicenseRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdateAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdateAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ProductLicenseRequest) validateLicenseStatus(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.LicenseStatus <= 0 || input.LicenseStatus > 3 {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.LicenseStatus)
		return
	}

	return
}
