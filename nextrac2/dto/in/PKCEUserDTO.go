package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
	"time"
)

type PKCERequest struct {
	AbstractDTO
	BodyRequest
	ID                  int64  `json:"id"`
	ParentClientID      string `json:"parent_client_id"`
	PKCEClientMappingID int64  `json:"pkce_client_mapping_id"`
	ClientTypeID        int64  `json:"client_type_id"`
	ClientAlias         string `json:"client_alias"`
	CompanyID           string `json:"company_id"`
	BranchID            string `json:"branch_id"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	Phone               string `json:"phone"`
	UpdatedAtStr        string `json:"updated_at"`
	UpdatedAt           time.Time
}

type PKCEReRequest struct {
	ID                  int64  `json:"id"`
	ClientID            string `json:"client_id"`
	AuthUserID          int64  `json:"auth_user_id"`
	PKCEClientMappingID int64  `json:"pkce_client_mapping_id"`
	PKCERequest         PKCERequest
}

func (input *PKCERequest) ValidateUnregisPKCE() errorModel.ErrorModel {
	fileName := "PKCEUserDTOIn.go"
	funcName := "ValidateUnregisPKCE"
	var err errorModel.ErrorModel

	err = input.validateRegisUnregisPKCE()
	if err.Error != nil {
		return err
	}

	//---------- Username is mandatory
	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	input.Username = strings.Trim(strings.ToLower(input.Username), " ")

	err = input.ValidateMinMaxString(input.Username, constanta.Username, 1, 20)
	if err.Error != nil {
		return err
	}

	//---------- Updated at is mandatory
	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *PKCERequest) ValidateRegistrationPKCE() errorModel.ErrorModel {
	var err errorModel.ErrorModel

	err = input.validateRegisUnregisPKCE()
	if err.Error != nil {
		return err
	}

	return input.validateMandatoryPKCE()
}

func (input *PKCERequest) validateRegisUnregisPKCE() errorModel.ErrorModel {
	fileName := "PKCEUserDTO.go"
	funcName := "validateRegisUnregisPKCE"
	var err errorModel.ErrorModel
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

	return errorModel.GenerateNonErrorModel()
}

func (input *PKCERequest) validateMandatoryPKCE() errorModel.ErrorModel {
	fileName := "PKCEUserDTOIn.go"
	funcName := "validateMandatoryPKCE"
	var validationResult bool
	var errField string
	var additionalInfo string
	var err errorModel.ErrorModel

	//---------- Username is mandatory
	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	validationResult, errField, additionalInfo = util.IsNexsoftUsernameStandardValid(input.Username)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Username, additionalInfo)
	}

	input.Username = strings.Trim(strings.ToLower(input.Username), " ")

	err = input.ValidateMinMaxString(input.Username, constanta.Username, 1, 20)
	if err.Error != nil {
		return err
	}

	//---------- Password is mandatory
	if util.IsStringEmpty(input.Password) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Password)
	}

	validationResult, errField, additionalInfo = util.IsNexsoftPasswordStandardValid(input.Password)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Password, additionalInfo)
	}

	//---------- Firstname is mandatory
	if util.IsStringEmpty(input.FirstName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.FirstName)
	}

	err = input.ValidateMinMaxString(input.FirstName, constanta.FirstName, 1, 50)
	if err.Error != nil {
		return err
	}

	//---------- Email is mandatory
	if util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Email)
	}

	if !util.IsEmailAddress(input.Email) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
	}

	input.Email = strings.ToLower(input.Email)

	//---------- Phone is mandatory
	if util.IsStringEmpty(input.Phone) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Phone)
	}

	number, isValid, errField := input.IsPhoneNumberValid(input.Phone)
	if !isValid {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Phone, "")
	}

	input.Phone = strconv.Itoa(number)

	return errorModel.GenerateNonErrorModel()
}

func (input *PKCERequest) ValidateViewUserCustomForUnregister() errorModel.ErrorModel {
	fileName := "PKCEUserDTO.go"
	funcName := "ValidateViewUserCustomForUnregister"

	if util.IsStringEmpty(input.ParentClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ParentClientID)
	}

	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *PKCERequest) checkParentClientID(fileName string, funcName string) errorModel.ErrorModel {
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

func (input *PKCERequest) checkCompanyID(fileName string, funcName string) errorModel.ErrorModel {
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

func (input *PKCERequest) checkBranchID(fileName string, funcName string) errorModel.ErrorModel {
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

func (input *PKCERequest) ValidateViewForUnregisterPKCE() errorModel.ErrorModel {
	fileName := "PKCEUserDTO.go"
	funcName := "ValidateViewForUnregisterPKCE"

	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *PKCERequest) ValidateUpdateClientName() (err errorModel.ErrorModel) {
	var (
		fileName = "PKCEUserDTO.go"
		funcName = "ValidateUpdateClientName"
	)

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

	if !util.IsStringEmpty(input.ClientAlias) {
		err = input.ValidateMinMaxString(input.ClientAlias, constanta.ClientAlias, 1, 50)
		if err.Error != nil {
			return
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ClientAlias, input.ClientAlias)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
