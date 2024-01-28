package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type LicenseVariantRequest struct {
	AbstractDTO
	ID                 int64  `json:"id"`
	LicenseVariantName string `json:"license_variant_name"`
	CreatedBy          int64  `json:"created_by"`
	CreatedAtStr       string `json:"created_at"`
	CreatedClient      string `json:"created_client"`
	UpdatedBy          int64  `json:"updated_by"`
	UpdatedAtStr       string `json:"updated_at"`
	UpdatedClient      string `json:"updated_client"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (input LicenseVariantRequest) ValidationInsertLicenseVariant() (err errorModel.ErrorModel) {
	return input.mandatoryValidation()
}

func (input *LicenseVariantRequest) ValidationUpdateLicenseVariant() (err errorModel.ErrorModel) {
	var (
		fileName = input.fileNameFuncNameLicenseVariant()
		funcName = "ValidationUpdateLicenseVariant"
	)

	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
		return
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
		return
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return input.mandatoryValidation()
}

func (input LicenseVariantRequest) ValidateViewLicenseVariant() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameLicenseVariant()
	funcName := "ValidateViewLicenseVariant"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *LicenseVariantRequest) ValidationDeleteLicenseVariant() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameLicenseVariant()
	funcName := "ValidationDeleteLicenseVariant"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
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

func (input LicenseVariantRequest) mandatoryValidation() (err errorModel.ErrorModel) {
	var (
		fileName = input.fileNameFuncNameLicenseVariant()
		funcName = "mandatoryValidation"
	)

	if util.IsStringEmpty(input.LicenseVariantName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.LicenseVariantName)
	}

	err = input.ValidateMinMaxString(input.LicenseVariantName, constanta.LicenseVariantName, 1, 22)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.LicenseVariantName, input.LicenseVariantName)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input LicenseVariantRequest) fileNameFuncNameLicenseVariant() (fileName string) {
	return "LicenseVariantDTO.go"
}
