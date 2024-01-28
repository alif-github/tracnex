package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strings"
)

type UpdateLicenseHWIDRequest struct {
	Hwid     string                 `json:"hwid"`
	Licenses []ProductEncryptDetail `json:"license"`
}

type UpdateLicenseHWIDByPassRequest struct {
	ClientTypeID         int64                          `json:"client_type_id"`
	ClientID             string                         `json:"client_id"`
	ClientIDInternal     string                         `json:"client_id_internal"`
	ClientSecret         string                         `json:"client_secret"`
	ClientSecretInternal string                         `json:"client_secret_internal"`
	SignatureKey         string                         `json:"signature_key"`
	SignatureKeyInternal string                         `json:"signature_key_internal"`
	HWID                 string                         `json:"hwid"`
	HWIDInternal         string                         `json:"hwid_internal"`
	License              []ProductEncryptDetailCSReport `json:"license"`
}

func (input *UpdateLicenseHWIDRequest) ValidateUpdateHWID() (err errorModel.ErrorModel) {
	var (
		fileName = UpdateLicenseHIWDDTOFileName
		funcName = "ValidateUpdateHWID"
	)

	input.Hwid = strings.TrimSpace(input.Hwid)
	if util.IsStringEmpty(input.Hwid) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.HWID)
	}

	if len(input.Licenses) == 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductLicense)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *UpdateLicenseHWIDByPassRequest) ValidateUpdateHWIDByPass() (err errorModel.ErrorModel) {
	var (
		fileName = "UpdateLicenseHIWDDTOFileName"
		funcName = "ValidateUpdateHWIDByPass"
	)

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ClientTypeID)
	}

	if util.IsStringEmpty(input.ClientID) || util.IsStringEmpty(input.ClientIDInternal) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientID)
	}

	if util.IsStringEmpty(input.SignatureKey) || util.IsStringEmpty(input.SignatureKeyInternal) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.SignatureKey)
	}

	if util.IsStringEmpty(input.ClientSecret) || util.IsStringEmpty(input.ClientSecretInternal) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientSecret)
	}

	if util.IsStringEmpty(input.HWIDInternal) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.HardwareIDInternal)
	}

	if len(input.License) < 1 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductLicense)
	}

	for _, itemProductEncrypt := range input.License {
		if util.IsStringEmpty(itemProductEncrypt.ProductKey) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductKey)
		}

		if util.IsStringEmpty(itemProductEncrypt.ProductEncrypt) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductEncrypt)
		}

		if util.IsStringEmpty(itemProductEncrypt.ProductSignature) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, "Product Signature")
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
