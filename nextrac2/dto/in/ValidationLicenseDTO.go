package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/cryptoModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

type ValidationLicenseRequest struct {
	AbstractDTO
	ClientID      string                 `json:"client_id"`
	ClientTypeID  int64                  `json:"client_type_id"`
	SignatureKey  string                 `json:"signature_key"`
	HwID          string                 `json:"hwid"`
	ProductDetail []ProductEncryptDetail `json:"product_detail"`
}

type ProductEncryptDetail struct {
	ProductKey     string `json:"product_key"`
	ProductEncrypt string `json:"product_encrypt"`
}

type ProductEncryptDetailCSReport struct {
	ProductKey       string `json:"product_key"`
	ProductEncrypt   string `json:"product_encrypt"`
	ProductSignature string `json:"product_signature"`
}

type ValidationLicenseJSONFile struct {
	cryptoModel.JSONFileActivationLicenseModel
	LicenseConfigID     int64                 `json:"license_config_id"`
	ProductKey          string                `json:"product_key"`
	ProductEncrypt      string                `json:"product_encrypt"`
	HWID                string                `json:"hwid"`
	ClientTypeID        int64                 `json:"client_type_id"`
	SalesmanLicenseList []SalesmanLicenseList `json:"salesman_license_list"`
}

type SalesmanLicenseList struct {
	ID         int64  `json:"id"`
	AuthUserID int64  `json:"auth_user_id"`
	UserID     string `json:"user_id"`
	Status     string `json:"status"`
}

func (input ValidationLicenseRequest) ValidateRequestValidationLicense() (err errorModel.ErrorModel) {
	if !util.IsStringEmpty(strings.TrimSpace(input.HwID)) {
		err = util2.ValidateMinMaxString(input.HwID, constanta.HWID, 1, 500)
		if err.Error != nil {
			return
		}
	}

	return input.mandatoryFieldValidation(ValidationLicenseDTOFileName, "ValidateRequestValidationLicense")
}

func (input ValidationLicenseRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(input.ClientID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientID)
		return
	}

	if input.ClientTypeID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientTypeID)
		return
	}

	if util.IsStringEmpty(input.SignatureKey) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.SignatureKey)
		return
	}

	if len(input.ProductDetail) < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductDetail)
		return
	}

	for _, detail := range input.ProductDetail {
		if util.IsStringEmpty(detail.ProductKey) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductKey)
			return
		}

		if util.IsStringEmpty(detail.ProductEncrypt) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductEncrypt)
			return
		}
	}

	return
}
