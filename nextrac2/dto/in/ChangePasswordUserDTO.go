package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

type ChangePasswordUserRequestDTO struct {
	ID					int64		`json:"id"`
	Username			string		`json:"username"`
	OldPassword			string		`json:"old_password"`
	NewPassword			string		`json:"new_password"`
	VerifyNewPassword	string		`json:"verify_new_password"`
}

func (input *ChangePasswordUserRequestDTO) ValidateChangePassword() errorModel.ErrorModel {
	fileName := "ChangePasswordUserDTO.go"
	funcName := "ValidateChangePassword"
	var validationResult bool

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	validationResult = util.IsStringEmpty(input.Username)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	validationResult = util.IsStringEmpty(input.OldPassword)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.OldPassword)
	}

	validationResult = util.IsStringEmpty(input.NewPassword)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.NewPassword)
	}

	validationResult = util.IsStringEmpty(input.VerifyNewPassword)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.VerifyPassword)
	}

	compareResult := strings.Compare(input.NewPassword, input.VerifyNewPassword)
	if compareResult != 0 {
		detail := util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.UserBundle, "FAILED_COMPARING_MESSAGE_CHANGE_PASSWORD", constanta.DefaultApplicationsLanguage, nil)
		return errorModel.GenerateFailedChangePassword(fileName, funcName, []string{detail})
	}

	return errorModel.GenerateNonErrorModel()
}
