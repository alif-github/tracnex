package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type ChangePassword struct {
	AbstractDTO
	ID                      int64      	`json:"id"`
	ParentClientID			string		`json:"parent_client_id"`
	ClientTypeID            int64      	`json:"client_type_id"`
	CompanyID               string     	`json:"company_id"`
	BranchID                string     	`json:"branch_id"`
	Username                string     	`json:"username"`
	CurrentPassword         string     	`json:"current_password"`
	NewPassword             string     	`json:"new_password"`
	ConfirmationPassword    string     	`json:"confirmation_password"`
	UpdatedAtStr            string     	`json:"updated_at"`
	UpdatedAt               time.Time
}

func (input *ChangePassword) ValidateReqChangePassword() errorModel.ErrorModel {
	err := input.validateMandatoryField()
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ChangePassword) validateMandatoryField() errorModel.ErrorModel {
	fileName := "RequestChangePasswordDTO.go"
	funcName := "validateMandatoryField"
	var validationResult bool
	var errField string
	var err errorModel.ErrorModel
	var additionalInfo string
	isNexmileValidation := input.ClientTypeID == constanta.ResourceNexmileID

	if isNexmileValidation {
		//---------- Parent client ID is mandatory
		err = input.checkParentClientID(fileName, funcName)
		if err.Error != nil {
			return err
		}
		//---------- Company ID is mandatory
		err = input.checkCompanyID(fileName, funcName)
		if err.Error != nil {
			return err
		}
		//---------- Branch ID is mandatory
		err = input.checkBranchID(fileName, funcName)
		if err.Error != nil {
			return err
		}
	} else {
		//---------- Optional Parent Client ID
		if input.ParentClientID != "" {
			err = input.checkParentClientID(fileName, funcName)
			if err.Error != nil {
				return err
			}
		}
		//---------- Optional Company ID
		if input.CompanyID != "" {
			err = input.checkCompanyID(fileName, funcName)
			if err.Error != nil {
				return err
			}
		}
		//---------- Optional Branch ID
		if input.BranchID != "" {
			err = input.checkBranchID(fileName, funcName)
			if err.Error != nil {
				return err
			}
		}
	}

	validationResult = util.IsStringEmpty(input.Username)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	err = input.ValidateMinMaxString(input.Username, constanta.Username, 1, 20)
	if err.Error != nil {
		return err
	}

	validationResult = util.IsStringEmpty(input.CurrentPassword)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CurrentPassword)
	}

	validationResult = util.IsStringEmpty(input.NewPassword)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.NewPassword)
	}

	validationResult = util.IsStringEmpty(input.ConfirmationPassword)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ConfirmationPassword)
	}

	if input.NewPassword != input.ConfirmationPassword {
		return errorModel.GenerateInvalidDifferentCompareData(fileName, funcName, constanta.NewPassword, constanta.ConfirmationPassword)
	}

	validationResult, errField, additionalInfo = util.IsNexsoftPasswordStandardValid(input.NewPassword)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Password, additionalInfo)
	}

	if input.NewPassword == input.CurrentPassword {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "Kata sandi baru dan Kata sandi lama harus berbeda", constanta.Password, "")
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ChangePassword) checkParentClientID(fileName string, funcName string) errorModel.ErrorModel {
	var validationResult bool
	var errField string

	if util.IsStringEmpty(input.ParentClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ParentClientID)
	}

	validationResult, errField, _ = util2.IsClientIDValid(input.ParentClientID)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.ParentClientID, "")
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ChangePassword) checkCompanyID(fileName string, funcName string) errorModel.ErrorModel {
	var err errorModel.ErrorModel

	if util.IsStringEmpty(input.CompanyID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyID)
	}

	err = input.ValidateMinMaxString(input.CompanyID, constanta.CompanyID, 1, 20)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ChangePassword) checkBranchID(fileName string, funcName string) errorModel.ErrorModel {
	var err errorModel.ErrorModel

	if util.IsStringEmpty(input.BranchID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.BranchID)
	}

	err = input.ValidateMinMaxString(input.BranchID, constanta.BranchID, 1, 20)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}