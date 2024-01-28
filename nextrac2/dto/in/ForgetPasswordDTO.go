package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type ResetPasswordRequest struct {
	Email        string `json:"email"`
	EmailLink    string `json:"email_link"`
	EmailMessage string `json:"email_message"`
}

type ChangePasswordByEmailRequest struct {
	AbstractDTO
	Email             string `json:"email"`
	ForgetCode        string `json:"forget_code"`
	NewPassword       string `json:"new_password"`
	VerifyNewPassword string `json:"verify_new_password"`
}

func (input ResetPasswordRequest) ToString() string {
	return util.StructToJSON(input)
}

func (input ResetPasswordRequest) ValidateStructForgetPassword() (err errorModel.ErrorModel) {
	var (
		fileName = "ForgetPasswordDTO.go"
		funcName = "ValidateStructForgetPassword"
	)

	if util.IsStringEmpty(input.Email) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Email)
		return
	}

	return
}

func (input ChangePasswordByEmailRequest) ToString() string {
	return util.StructToJSON(input)
}

func (input ChangePasswordByEmailRequest) ValidateStructChangePassword() (err errorModel.ErrorModel) {
	var (
		fileName = "ForgetPasswordDTO.go"
		funcName = "ValidateStructChangePassword"
	)

	if util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Email)
	}

	if util.IsStringEmpty(input.ForgetCode) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ForgetCode)
	}

	if util.IsStringEmpty(input.NewPassword) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NewPassword)
	}

	if err = input.ValidateMinMaxString(input.NewPassword, constanta.NewPassword, 1, 20); err.Error != nil {
		return
	}

	if util.IsStringEmpty(input.VerifyNewPassword) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.VerifyPassword)
	}

	if err = input.ValidateMinMaxString(input.VerifyNewPassword, constanta.VerifyPassword, 1, 20); err.Error != nil {
		return
	}

	return
}
