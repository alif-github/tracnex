package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type LicenseConfigRequest struct {
	ID                  int64  `json:"id"`
	InstallationID      int64  `json:"installation_id"`
	NoOfUser            int64  `json:"no_of_user"`
	IsUserConcurrent    string `json:"is_user_concurrent"`
	ProductValidFromStr string `json:"product_valid_from"`
	ProductValidThruStr string `json:"product_valid_thru"`
	UpdatedAtStr        string `json:"updated_at"`
	ProductValidFrom    time.Time
	ProductValidThru    time.Time
	UpdatedAt           time.Time
}

type LicenseConfigMultipleRequest struct {
	ProductValidThruStr string  `json:"product_valid_thru"`
	LicenseConfigID     []int64 `json:"license_config_id"`
	ProductValidThru    time.Time
}

func (input *LicenseConfigRequest) ValidateInsertLicenseConfig() errorModel.ErrorModel {
	var (
		fileName = input.fileNameFuncNameLicenseConfig()
		funcName = "ValidateInsertLicenseConfig"
		err      errorModel.ErrorModel
	)

	if input.InstallationID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.InstallationID)
	}

	if util.IsStringEmpty(input.ProductValidFromStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductValidFrom)
	}

	input.ProductValidFrom, err = TimeStrToTimeWithTimeFormat(input.ProductValidFromStr, constanta.ProductValidFrom, constanta.DefaultInstallationTimeFormat)
	if err.Error != nil {
		return err
	}

	timeNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.FixedZone(config.ApplicationConfiguration.GetLocalTimezone().Zone, config.ApplicationConfiguration.GetLocalTimezone().Offset))
	if input.ProductValidFrom.Before(timeNow) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.MaxValueProductValidFrom, constanta.ProductValidFrom, "")
	}

	if util.IsStringEmpty(input.ProductValidThruStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductValidThru)
	}

	input.ProductValidThru, err = TimeStrToTimeWithTimeFormat(input.ProductValidThruStr, constanta.ProductValidThru, constanta.DefaultInstallationTimeFormat)
	if err.Error != nil {
		return err
	}

	timeLimit := time.Date(input.ProductValidThru.Year()+3, time.January, 1, 0, 0, 0, 0, time.FixedZone(config.ApplicationConfiguration.GetLocalTimezone().Zone, config.ApplicationConfiguration.GetLocalTimezone().Offset))
	timeLimit = timeLimit.Add(time.Hour * -24)
	if input.ProductValidThru.After(timeLimit) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.MaxValueProductValidThru, constanta.ProductValidThru, "")
	}

	//--- Request high authority 15/05/2023
	//if input.NoOfUser < 1 {
	//	err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NumberOfUser)
	//	return err
	//}

	isValidFromBeforeThru := input.ProductValidFrom.Before(input.ProductValidThru)
	if !isValidFromBeforeThru {
		err = errorModel.GenerateDateValidateFromThru(fileName, funcName, "E-6-TRAC-SRV-016")
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *LicenseConfigMultipleRequest) ValidateInsertMultipleLicenseConfig() errorModel.ErrorModel {
	var (
		fileName            = "LicenseConfigDTO.go"
		funcName            = "ValidateInsertMultipleLicenseConfig"
		err                 errorModel.ErrorModel
		isValidThruAfterNow bool
	)

	if util.IsStringEmpty(input.ProductValidThruStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductValidThru)
	}

	input.ProductValidThru, err = TimeStrToTimeWithTimeFormat(input.ProductValidThruStr, constanta.ProductValidThru, constanta.DefaultInstallationTimeFormat)
	if err.Error != nil {
		return err
	}

	isValidThruAfterNow = input.ProductValidThru.After(time.Now())
	if !isValidThruAfterNow {
		err = errorModel.GenerateDateValidateFromThru(fileName, funcName, "E-4-ETR-TRAC-SRV-005")
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *LicenseConfigRequest) ValidationDeleteLicenseConfig() (err errorModel.ErrorModel) {
	var (
		fileName = input.fileNameFuncNameLicenseConfig()
		funcName = "ValidationDeleteLicenseConfig"
	)

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

func (input *LicenseConfigRequest) ValidationUpdateLicenseConfig() (err errorModel.ErrorModel) {
	var (
		fileName = input.fileNameFuncNameLicenseConfig()
		funcName = "ValidationUpdateLicenseConfig"
	)

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

func (input *LicenseConfigRequest) ValidationViewLicenseConfig() (err errorModel.ErrorModel) {
	var (
		fileName = input.fileNameFuncNameLicenseConfig()
		funcName = "ValidationViewLicenseConfig"
	)

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input LicenseConfigRequest) fileNameFuncNameLicenseConfig() (fileName string) {
	return "LicenseConfigDTO.go"
}
