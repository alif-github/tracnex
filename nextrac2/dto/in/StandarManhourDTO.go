package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type StandarManhourRequest struct {
	AbstractDTO
	ID           int64   `json:"id"`
	Case         string  `json:"case_name"`
	DepartmentID int64   `json:"department_id"`
	Manhour      float64 `json:"manhour"`
	UpdatedAtStr string  `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *StandarManhourRequest) ValidateInsert() (err errorModel.ErrorModel) {
	return input.mandatoryValidation("StandarManhourDTO.go", "ValidateInsert")
}

func (input *StandarManhourRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	var (
		fileName = "StandarManhourDTO.go"
		funcName = "ValidateUpdate"
	)

	err = input.mandatoryValidation(fileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(fileName, funcName)
	return
}

func (input *StandarManhourRequest) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		fileName = "StandarManhourDTO.go"
		funcName = "ValidateDelete"
	)

	err = input.validationForUpdateAndDelete(fileName, funcName)
	return
}

func (input *StandarManhourRequest) ValidateView() (err errorModel.ErrorModel) {
	var (
		fileName = "StandarManhourDTO.go"
		funcName = "ValidateDelete"
	)

	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *StandarManhourRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *StandarManhourRequest) mandatoryValidation(funcName string, fileName string) (err errorModel.ErrorModel) {
	//-- Validate Case
	if util.IsStringEmpty(input.Case) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Case)
	}

	err = input.ValidateMinMaxString(input.Case, constanta.Case, 1, 256)
	if err.Error != nil {
		return
	}

	//err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Case, input.Case)
	//if err.Error != nil {
	//	return
	//}

	//-- Validate Department
	if input.DepartmentID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.DepartmentId)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
