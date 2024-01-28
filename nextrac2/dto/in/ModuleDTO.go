package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type ModuleRequest struct {
	AbstractDTO
	ID           int64  `json:"id"`
	ModuleName   string `json:"module_name"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *ModuleRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(ModuleDTOFileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ModuleRequest) ValidateDelete() (err errorModel.ErrorModel) {
	funcName := "ValidateDelete"
	return input.validationForUpdateAndDelete(ModuleDTOFileName, funcName)
}

func (input *ModuleRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	err = input.mandatoryValidation(ModuleDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(ModuleDTOFileName, funcName)
	return
}

func (input *ModuleRequest) ValidateInsert() (err errorModel.ErrorModel) {
	funcName := "ValidateInsert"
	return input.mandatoryValidation(ModuleDTOFileName, funcName)
}

func (input *ModuleRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
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

func (input *ModuleRequest) mandatoryValidation(funcName string, fileName string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(input.ModuleName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ModuleName)
	}

	err = input.ValidateMinMaxString(input.ModuleName, constanta.ModuleName, 1, 22)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ModuleName, input.ModuleName)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
