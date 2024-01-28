package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type EmployeeAllowanceRequest struct {
	AbstractDTO
	ID                    int64     `json:"id"`
	AllowanceName         string    `json:"allowance_name"`
	AllowanceType         string    `json:"allowance_type"`
	Value                 string    `json:"value"`
	Active                bool      `json:"active"`
	UpdatedAtStr          string    `json:"updated_at"`
	UpdatedAt             time.Time
}

func (input *EmployeeAllowanceRequest) ValidateEmployeeAllowance(isUpdate bool) (err errorModel.ErrorModel) {
	funcName := "ValidateEmployeeAllowance"
	fileName := "EmployeeAllowanceDTO.go"

	if input.AllowanceName == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "allowance_name")
	}

	err = input.ValidateMinMaxString(input.AllowanceName, "name", 3, 100)
	if err.Error != nil {
		return
	}

	if input.AllowanceType != ""{
		if len(input.AllowanceType) > 50{
			return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "type", 50)
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