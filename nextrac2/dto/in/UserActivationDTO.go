package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type UserActivationRequest struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	EmailCode string `json:"email_code"`
	Username  string `json:"username"`
}

func (input *UserActivationRequest) ValidateActivationUser() errorModel.ErrorModel {
	fileName := "UserActivationDTO.go"
	funcName := "ValidateActivationUser"
	var validationResult bool

	if input.UserID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
	}

	if util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Email)
	}

	validationResult = util.IsEmailAddress(input.Email)
	if !validationResult {
		return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
	}

	if util.IsStringEmpty(input.EmailCode) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.EmailCode)
	}

	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	return errorModel.GenerateNonErrorModel()
}

type UserActivationPhoneRequest struct {
	UserID  int64  `json:"user_id"`
	Phone   string `json:"phone"`
	OTPCode string `json:"otp_code"`
}

func (input *UserActivationPhoneRequest) ValidateActivationUser() errorModel.ErrorModel {
	fileName := "UserActivationDTO.go"
	funcName := "ValidateActivationUser"

	if input.UserID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
	}

	if util.IsStringEmpty(input.Phone) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Phone)
	}

	if !IsPhoneNumberWithCountryCodeMDB(input.Phone) {
		return errorModel.GenerateFieldFormatWithRuleError(userVerificationDTOFileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
	}

	return errorModel.GenerateNonErrorModel()
}
