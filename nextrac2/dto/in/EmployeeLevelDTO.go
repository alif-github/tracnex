package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"regexp"
	"time"
)

type EmployeeLevelRequest struct {
	AbstractDTO
	ID                    int64     `json:"id"`
	Level                 string    `json:"level"`
	Description           string    `json:"description"`
	UpdatedAtStr          string    `json:"updated_at"`
	UpdatedAt             time.Time
}

func (input *EmployeeLevelRequest) ValidateEmployeeLevel(isUpdate bool) (err errorModel.ErrorModel) {
	funcName := "ValidateEmployeeLevel"
	fileName := "EmployeeLevelDTO.go"

	if input.Level == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "level")
	}

	err = input.ValidateMinMaxString(input.Level, "level", 1, 20)
	if err.Error != nil {
		return
	}

	if input.Description != ""{
		if len(input.Description) >= 257{
			return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "deskripsi", 256)
		}
	}

	validInput := regexp.MustCompile("^[a-zA-Z0-9 ]*$")
	isValid := validInput.MatchString(input.Level)
	if !isValid {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "tidak boleh mengandung spesial karakter", "level", "")
	    return err
	}

	if isUpdate {
		if util.IsStringEmpty(input.UpdatedAtStr) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
		}

		input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
		if err.Error != nil {
			return
		}
	}

	return errorModel.GenerateNonErrorModel()
}