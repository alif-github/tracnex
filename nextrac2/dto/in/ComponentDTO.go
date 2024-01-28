package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type ComponentRequest struct {
	AbstractDTO
	ID            int64  `json:"id"`
	ComponentName string `json:"component_name"`
	UpdatedAtStr  string `json:"updated_at"`
	UpdatedAt     time.Time
}

func (input *ComponentRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(ComponentDTOFileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ComponentRequest) ValidateDelete() (err errorModel.ErrorModel) {
	funcName := "ValidateDelete"
	return input.validationForUpdateAndDelete(ComponentDTOFileName, funcName)
}

func (input *ComponentRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	err = input.mandatoryValidation(ComponentDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(ComponentDTOFileName, funcName)
	return
}

func (input *ComponentRequest) ValidateInsert() (err errorModel.ErrorModel) {
	funcName := "ValidateInsert"
	return input.mandatoryValidation(ComponentDTOFileName, funcName)
}

func (input *ComponentRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *ComponentRequest) mandatoryValidation(funcName string, fileName string) (err errorModel.ErrorModel) {
	err = input.ValidateMinMaxString(input.ComponentName, constanta.ComponentName, 1, 22)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ComponentName, input.ComponentName)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
