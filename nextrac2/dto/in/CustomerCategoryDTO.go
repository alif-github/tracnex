package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type CustomerCategoryRequest struct {
	AbstractDTO
	ID                   int64  `json:"id"`
	CustomerCategoryID   string `json:"customer_category_id"`
	CustomerCategoryName string `json:"customer_category_name"`
	UpdatedAtStr         string `json:"updated_at"`
	UpdatedAt            time.Time
}

func (input *CustomerCategoryRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(CustomerGroupDTOFileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *CustomerCategoryRequest) ValidateDelete() (err errorModel.ErrorModel) {
	funcName := "ValidateDelete"
	return input.validationForUpdateAndDelete(CustomerCategoryDTOFileName, funcName)
}

func (input *CustomerCategoryRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	if util.IsStringEmpty(input.CustomerCategoryName) {
		err = errorModel.GenerateEmptyFieldError(CustomerCategoryDTOFileName, funcName, constanta.CustomerCategoryName)
		return
	}
	err = input.ValidateMinMaxString(input.CustomerCategoryName, constanta.CustomerCategoryName, 1, 50)
	if err.Error != nil {
		return
	}

	if err = util2.ValidateSpecialCharacter(CustomerCategoryDTOFileName, funcName, constanta.CustomerCategoryName, input.CustomerCategoryName); err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(CustomerCategoryDTOFileName, funcName)
	if err.Error != nil {
		return
	}
	return
}

func (input *CustomerCategoryRequest) ValidateInsert() (err errorModel.ErrorModel) {
	funcName := "ValidateInsertCustomerCategory"
	return input.mandatoryFieldValidation(CustomerCategoryDTOFileName, funcName)
}

func (input *CustomerCategoryRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *CustomerCategoryRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	//var isValid bool
	//var errField string

	//---------- Check is string empty for Customer Category ID
	if util.IsStringEmpty(input.CustomerCategoryID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CustomerCategoryID)
		return
	}
	err = input.ValidateMinMaxString(input.CustomerCategoryID, constanta.CustomerCategoryID, 1, 22)
	if err.Error != nil {
		return
	}

	//isValid, errField = util.IsNexsoftDirectoryNameStandardValid(input.CustomerCategoryID)
	//if !isValid {
	//	return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.CustomerCategoryID, "")
	//}

	if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.CustomerCategoryID, input.CustomerCategoryID); err.Error != nil {
		return
	}

	//---------- Check is string empty for Customer Category Name
	if util.IsStringEmpty(input.CustomerCategoryName) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CustomerCategoryName)
		return
	}
	err = input.ValidateMinMaxString(input.CustomerCategoryName, constanta.CustomerCategoryName, 1, 50)
	if err.Error != nil {
		return
	}

	if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.CustomerCategoryName, input.CustomerCategoryName); err.Error != nil {
		return
	}
	return
}
