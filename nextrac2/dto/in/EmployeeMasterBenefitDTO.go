package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type EmpBenefitRequest struct {
	AbstractDTO
	ID                  int64     `json:"id"`
	BenefitName         string    `json:"benefit_name"`
	BenefitType         string    `json:"benefit_type"`
	Value               string    `json:"value"`
	Active              bool      `json:"active"`
	UpdatedAtStr        string    `json:"updated_at"`
	UpdatedAt           time.Time
}

func (input *EmpBenefitRequest) ValidateEmployeeMasterBenefit(isUpdate bool) (err errorModel.ErrorModel) {
	funcName := "ValidateEmployeeMasterBenefit"
	fileName := "EmployeeMasterBenefitDTO.go"

	if input.BenefitName == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "benefit_name")
	}

	err = input.ValidateMinMaxString(input.BenefitName, "name", 3, 50)
	if err.Error != nil {
		return
	}

	if input.BenefitType != ""{
		if len(input.BenefitType) > 256{
			return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "type", 256)
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