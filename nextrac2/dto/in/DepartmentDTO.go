package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

var fileNameDepartmentDTO = "DepartmentDTO.go"

type DepartmentRequest struct {
	ID             int64  `json:"id"`
	DepartmentName string `json:"department_name"`
	Description    string `json:"description"`
	UpdatedAtStr   string `json:"updated_at"`
	UpdatedAt      time.Time
}

func (input *DepartmentRequest) ValidateInsert() errorModel.ErrorModel {
	return input.mandatoryValidation()
}

func (input *DepartmentRequest) mandatoryValidation() (err errorModel.ErrorModel) {
	var (
		funcName = "mandatoryValidation"
		fileName = fileNameDepartmentDTO
	)

	if util.IsStringEmpty(input.DepartmentName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.DepartmentName)
	}

	err = util2.ValidateMinMaxString(input.DepartmentName, constanta.DepartmentName, 1, 100)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacterAlphabet(fileName, funcName, constanta.DepartmentName, input.DepartmentName)
	if err.Error != nil {
		return
	}

	return
}

func (input *DepartmentRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileNameDepartmentDTO, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *DepartmentRequest) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateDelete"
	)

	return input.validationForUpdateAndDelete(fileNameDepartmentDTO, funcName)
}

func (input *DepartmentRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *DepartmentRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	err = input.mandatoryValidation()
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(fileNameDepartmentDTO, funcName)
	return
}
