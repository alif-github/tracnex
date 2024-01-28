package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type EmployeePosition struct {
	AbstractDTO
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CompanyID    int64  `json:"company_id"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *EmployeePosition) ValidateView() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeePositionDTO.go"
		funcName = "ValidateView"
	)

	if input.ID < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *EmployeePosition) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeePositionDTO.go"
		funcName = "ValidateDelete"
	)

	return input.validationForUpdateAndDelete(fileName, funcName)
}

func (input *EmployeePosition) ValidateInsert() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeePositionDTO.go"
		funcName = "ValidateInsert"
	)

	if !util.IsStringEmpty(input.Description) {
		err = input.ValidateMinMaxString(input.Description, constanta.Description, 1, 256)
		if err.Error != nil {
			return
		}
	}

	return input.mandatoryValidation(fileName, funcName)
}

func (input *EmployeePosition) ValidateUpdate() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeePositionDTO.go"
		funcName = "ValidateUpdate"
	)

	err = input.mandatoryValidation(fileName, funcName)
	if err.Error != nil {
		return
	}

	return input.validationForUpdateAndDelete(fileName, funcName)
}

func (input *EmployeePosition) mandatoryValidation(fileName, funcName string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(input.Name) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Position)
	}

	err = input.ValidateMinMaxString(input.Name, constanta.Position, 1, 100)
	if err.Error != nil {
		return
	}

	if input.CompanyID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.CompanyID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *EmployeePosition) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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
