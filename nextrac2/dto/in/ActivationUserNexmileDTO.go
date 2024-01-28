package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strconv"
	"strings"
)

type ActivationUserNexmileRequest struct {
	AbstractDTO
	UserRegistrationDetailID int64  `json:"user_registration_detail_id"`
	ParentClientID           string `json:"parent_client_id"`
	FirstName                string `json:"first_name"`
	LastName                 string `json:"last_name"`
	AndroidID                string `json:"android_id"`
	ClientID                 string `json:"client_id"`
	UserID                   string `json:"user_id"`
	AuthUserID               int64  `json:"auth_user_id"`
	AliasName                string `json:"alias_name"`
	Email                    string `json:"email"`
	Phone                    string `json:"phone"`
}

func (input *ActivationUserNexmileRequest) ValidateActivationUserNexmile() (err errorModel.ErrorModel) {
	funcName := "ValidateActivationUserNexmile"

	if err = input.mandatoryValidation(ActivationUserNexmileDTOFileName, funcName); err.Error != nil {
		return
	}

	return input.validateOptionalField(ActivationUserNexmileDTOFileName, funcName)
}

func (input *ActivationUserNexmileRequest) mandatoryValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.UserRegistrationDetailID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UserRegistrationDetailIDName)
		return
	}

	if util.IsStringEmpty(input.ParentClientID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingParentClientID)
		return
	}

	if util.IsStringEmpty(input.ClientID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingClientID)
		return
	}

	if util.IsStringEmpty(input.UserID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UserID)
		return
	}

	if input.AuthUserID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.AuthUserID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ActivationUserNexmileRequest) validateOptionalField(fileName string, funcName string) (err errorModel.ErrorModel) {
	var splitPhone []string
	var phone int
	var validationResult bool

	if !util.IsStringEmpty(input.Email) {
		if err = input.ValidateMinMaxString(input.Email, constanta.Email, 1, 100); err.Error != nil {
			return
		}

		if !util.IsEmailAddress(input.Email) {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
			return
		}
	}

	if !util.IsStringEmpty(input.Phone) {
		splitPhone = strings.Split(input.Phone, "-")
		if !util.IsCountryCode(splitPhone[0]) {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
			return
		}

		if phone, validationResult = util.IsPhoneNumber(splitPhone[1]); !validationResult {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
			return
		}

		input.Phone = splitPhone[0] + "-" + strconv.Itoa(phone)
	}

	return
}
