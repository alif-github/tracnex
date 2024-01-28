package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"regexp"
	"time"
)

type EmployeeGradeRequest struct {
	AbstractDTO
	ID                    int64     `json:"id"`
	Grade                 string    `json:"grade"`
	Description           string    `json:"description"`
	UpdatedAtStr          string    `json:"updated_at"`
	UpdatedAt             time.Time
}

func (input *EmployeeGradeRequest) ValidateEmployeeGrade(isUpdate bool) (err errorModel.ErrorModel) {
	funcName := "ValidateEmployeeGrade"
	fileName := "EmployeeLevelDTO.go"

	if input.Grade == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "grade")
	}

	err = input.ValidateMinMaxString(input.Grade, "grade", 1, 20)
	if err.Error != nil {
		return
	}

	validInput := regexp.MustCompile("^[a-zA-Z0-9 ]*$")
	isValid := validInput.MatchString(input.Grade)
	if !isValid {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "tidak boleh mengandung spesial karakter", "grade", "")
		return err
	}

	if input.Description != ""{
		if len(input.Description) >= 257{
			return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "deskripsi", 256)
		}
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