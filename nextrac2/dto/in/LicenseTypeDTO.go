package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type LicenseTypeRequest struct {
	AbstractDTO
	ID              int64  `json:"id"`
	LicenseTypeName string `json:"license_type_name"`
	LicenseTypeDesc string `json:"license_type_desc"`
	IsMainLicense   bool   `json:"is_main_license"`
	UpdatedAtStr    string `json:"updated_at"`
	UpdatedAt       time.Time
}

func (input *LicenseTypeRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(LicenseTypeDTOFileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *LicenseTypeRequest) ValidateDelete() (err errorModel.ErrorModel) {
	return input.validationForUpdateAndDelete(LicenseTypeDTOFileName, "ValidateDelete")
}

func (input *LicenseTypeRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"
	err = input.validationForUpdateAndDelete(LicenseTypeDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.mandatoryFieldValidation(LicenseTypeDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *LicenseTypeRequest) ValidateInsert() (err errorModel.ErrorModel) {
	funcName := "ValidateInsert"
	return input.mandatoryFieldValidation(LicenseTypeDTOFileName, funcName)
}

func (input *LicenseTypeRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ID)
		return
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *LicenseTypeRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(input.LicenseTypeName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.LicenseTypeName)
	}

	err = input.ValidateMinMaxString(input.LicenseTypeName, constanta.LicenseTypeName, 1, 22)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.LicenseTypeName, input.LicenseTypeName)
	if err.Error != nil {
		return
	}

	err = input.ValidateMinMaxString(input.LicenseTypeDesc, constanta.LicenseTypeDesc, 0, 200)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
