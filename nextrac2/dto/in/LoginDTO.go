package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type VerifyDTOIn struct {
	Authorize     string
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	CheckPassword bool   `json:"check_password"`
}

func (input VerifyDTOIn) ValidateLoginDTO() errorModel.ErrorModel {
	var (
		fileName = "LoginDTO.go"
		funcName = "ValidateLoginDTO"
	)

	if !util.IsStringEmpty(input.Username) && !util.IsStringEmpty(input.Email) {
		return errorModel.GenerateUnauthorizedClientError(fileName, funcName)
	}

	if util.IsStringEmpty(input.Username) && util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	if input.Password == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Password)
	}

	if input.Authorize == "" {
		return errorModel.GenerateUnauthorizedClientError(fileName, funcName)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input VerifyDTOIn) ValidateNexmileLoginDTO() errorModel.ErrorModel {
	var (
		fileName = "LoginDTO.go"
		funcName = "ValidateNexmileLoginDTO"
	)

	if !util.IsStringEmpty(input.Username) && !util.IsStringEmpty(input.Email) {
		return errorModel.GenerateUnauthorizedClientError(fileName, funcName)
	}

	if util.IsStringEmpty(input.Username) && util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "User ID")
	}

	if input.Password == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Password)
	}

	return errorModel.GenerateNonErrorModel()
}
