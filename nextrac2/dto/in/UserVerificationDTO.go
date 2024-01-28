package in

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

var userVerificationDTOFileName = "UserVerificationDTO.go"

type UserVerificationRequest struct {
	ClientTypeID int64  `json:"client_type_id"`
	ClientID     string `json:"client_id"`
	UserID       string `json:"user_id"`
	Password     string `json:"password"`
	AndroidID    string `json:"android_id"`
	UniqueID1    string `json:"unique_id_1"`
	UniqueID2    string `json:"unique_id_2"`
	OTP          string `json:"otp"`
	CountryCode  string `json:"country_code"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	AuthUserID   int64  `json:"auth_user_id"`
}

func (input *UserVerificationRequest) ValidateVerifyingUser() (err errorModel.ErrorModel) {
	var (
		fileName = "UserVerificationDTO.go"
		funcName = "ValidateVerifyingUser"
	)

	if util.IsStringEmpty(input.OTP) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.OTP)
		return
	}

	if util.IsStringEmpty(input.Password) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Password)
		return
	}

	if !util.IsStringEmpty(input.Phone) {
		if !IsPhoneNumberWithCountryCodeMDB(input.Phone) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
			return
		}
	}

	if !util.IsStringEmpty(input.AndroidID) {
		err = util2.ValidateMinMaxString(input.AndroidID, constanta.AndroidID, 1, 100)
		if err.Error != nil {
			return
		}
	}

	return input.mandatoryValidation(userVerificationDTOFileName, funcName)
}

func (input *UserVerificationRequest) ValidateResendOTP() (err errorModel.ErrorModel) {
	var (
		fileName = userVerificationDTOFileName
		funcName = "ValidateResendOTP"
	)

	if input.AuthUserID < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.AuthUserID)
		return
	}

	if (!util.IsStringEmpty(input.Phone) && !util.IsStringEmpty(input.Email)) || (util.IsStringEmpty(input.Phone) && util.IsStringEmpty(input.Email)) {
		fmt.Println("Input Phone : ", input.Phone)
		fmt.Println("Input Email : ", input.Email)
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.InvalidRequestEmailPhone)
		return
	}

	if !util.IsStringEmpty(input.Phone) {
		if util.IsStringEmpty(input.CountryCode) {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.CountryCode)
			return
		}

		if !IsPhoneNumberWithCountryCodeMDB(fmt.Sprintf(`%s-%s`, input.CountryCode, input.Phone)) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
			return
		}
	} else {
		if !util.IsStringEmpty(input.CountryCode) {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Phone)
			return
		}
	}

	if !util.IsStringEmpty(input.Email) {
		if !util.IsEmailAddress(input.Email) {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
			return
		}
	}

	return input.mandatoryValidation(fileName, funcName)
}

func (input *UserVerificationRequest) mandatoryValidation(fileName, funcName string) (err errorModel.ErrorModel) {
	if input.ClientTypeID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientTypeID)
		return
	}

	if util.IsStringEmpty(input.UserID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UserID)
		return
	}

	if util.IsStringEmpty(input.UniqueID1) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UniqueID1)
		return
	}

	err = util2.ValidateMinMaxString(input.UniqueID1, constanta.UniqueID1, 1, 20)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(input.UniqueID2) {
		err = util2.ValidateMinMaxString(input.UniqueID2, constanta.UniqueID2, 1, 20)
		if err.Error != nil {
			return
		}
	}

	if !util.IsStringEmpty(input.Phone) && !util.IsStringEmpty(input.Email) {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.InvalidOTPRequestCombination)
		return
	}

	if util.IsStringEmpty(input.Phone) && util.IsStringEmpty(input.Email) {
		err = errorModel.GenerateRequestError(fileName, funcName, constanta.InvalidVerifyingRequetCombination)
		return
	}

	if !util.IsStringEmpty(input.Email) {
		if !util.IsEmailAddress(input.Email) {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
			return
		}

		err = util2.ValidateMinMaxString(input.Email, constanta.Email, 1, 100)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
