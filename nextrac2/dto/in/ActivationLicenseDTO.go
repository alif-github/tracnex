package in

import (
	"github.com/Azure/go-autorest/autorest/date"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/cryptoModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

type ActivationLicenseRequest struct {
	AbstractDTO
	ClientID     string           `json:"client_id"`
	ClientTypeID int64            `json:"client_type_id"`
	SignatureKey string           `json:"signature_key"`
	Hwid         string           `json:"hwid"`
	DetailClient []UniqueIDClient `json:"detail_client"`
}

type ActivateLicenseDataRequest struct {
	LicenseConfigID  int64                 `json:"license_config_id"`
	ProductKey       string                `json:"product_key"`
	ProductEncrypt   string                `json:"product_encrypt"`
	ProductSignature string                `json:"product_signature"`
	IsUserConcurrent string                `json:"is_user_concurrent"`
	ClientID         string                `json:"client_id"`
	ClientSecret     string                `json:"client_secret"`
	Hwid             string                `json:"hwid"`
	UniqueID1        string                `json:"unique_id_1"`
	UniqueID2        string                `json:"unique_id_2"`
	ParentCustomerID int64                 `json:"parent_customer_id"`
	CustomerID       int64                 `json:"customer_id"`
	SiteID           int64                 `json:"site_id"`
	InstallationID   int64                 `json:"installation_id"`
	NumberOfUser     int64                 `json:"number_of_user"`
	ProductValidFrom date.Date             `json:"product_valid_from"`
	ProductValidThru date.Date             `json:"product_valid_thru"`
	ClientTypeID     int64                 `json:"client_type_id"`
	SalesmanList     []SalesmanLicenseList `json:"salesman_list"`
}

type JSONActivationLicenseRequest struct {
	LicenseConfigID int64 `json:"license_config_id"`
	cryptoModel.EncryptLicenseRequestModel
}

func (input *ActivationLicenseRequest) ValidateActivateLicense() (err errorModel.ErrorModel) {

	if err = input.mandatoryValidation(ActivationLicenseDTOFileName, "ValidateActivateLicense"); err.Error != nil {
		return
	}

	return input.optionalValidation(ActivationLicenseDTOFileName, "ValidateActivateLicense")
}

func (input *ActivationLicenseRequest) optionalValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if !util.IsStringEmpty(strings.TrimSpace(input.Hwid)) {
		err = util2.ValidateMinMaxString(input.Hwid, constanta.HWID, 1, 500)
		if err.Error != nil {
			return
		}
	}

	return
}

func (input *ActivationLicenseRequest) mandatoryValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(input.ClientID) {
		err = errorModel.GenerateEmptyFieldError(ActivationLicenseDTOFileName, funcName, constanta.ClientMappingClientID)
		return
	}

	if util.IsStringEmpty(input.SignatureKey) {
		err = errorModel.GenerateEmptyFieldError(ActivationLicenseDTOFileName, funcName, constanta.SignatureKey)
		return
	}

	if len(input.DetailClient) < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ApplicationDetail)
		return
	}

	if input.ClientTypeID < 1 {
		err = errorModel.GenerateEmptyFieldError(ActivationLicenseDTOFileName, funcName, constanta.ClientTypeID)
	}

	return
}
