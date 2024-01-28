package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type CheckLicenseNamedUserRequest struct {
	ClientId     string `json:"client_id"`
	ClientTypeID int64  `json:"client_type_id"`
	UniqueId1    string `json:"unique_id_1"`
	UniqueId2    string `json:"unique_id_2"`
}

func (input CheckLicenseNamedUserRequest) ValidateCheckUserLicense() errorModel.ErrorModel {
	fileName := "CheckLicenseNamedUserDTO.go"
	funcName := "ValidateCheckUserLicense"

	if util.IsStringEmpty(input.ClientId) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingClientID)
	}

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NewClientType)
	}

	if util.IsStringEmpty(input.UniqueId1) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UniqueID1)
	}

	return errorModel.GenerateNonErrorModel()
}
