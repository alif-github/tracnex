package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type ClientCredential struct {
	AbstractDTO
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	SignatureKey string `json:"signature_key"`
}

func (input ClientCredential) MandatoryValidationClientCredential() (err errorModel.ErrorModel) {
	fileName := "ClientCredentialDTO.go"
	funcName := "MandatoryValidationClientCredential"
	if util.IsStringEmpty(input.ClientID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingClientID)
		return
	}

	err = input.ValidateMinMaxString(input.ClientID, constanta.ClientMappingClientID, 1, 256)
	if err.Error != nil {
		return
	}

	if util.IsStringEmpty(input.ClientSecret) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientSecret)
		return
	}

	err = input.ValidateMinMaxString(input.ClientSecret, constanta.ClientSecret, 1, 256)
	if err.Error != nil {
		return
	}

	if util.IsStringEmpty(input.SignatureKey) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.SignatureKey)
		return
	}

	err = input.ValidateMinMaxString(input.SignatureKey, constanta.SignatureKey, 1, 256)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
