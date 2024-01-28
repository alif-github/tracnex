package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type ValidationNamedUserRequest struct {
	AbstractDTO
	ClientId   string `json:"client_id"`
	AuthUserId int64  `json:"auth_user_id"`
	UserId     string `json:"user_id"`
	UniqueId1  string `json:"unique_id_1"`
	UniqueId2  string `json:"unique_id_2"`
}

func (input ValidationNamedUserRequest) ValidationDTONamedUser() errorModel.ErrorModel {
	return input.mandatoryFieldValidation()
}

func (input ValidationNamedUserRequest) mandatoryFieldValidation() (err errorModel.ErrorModel) {
	var (
		funcName = "mandatoryFieldValidation"
	)
	if util.IsStringEmpty(input.UniqueId1) {
		err = errorModel.GenerateEmptyFieldError(ValidationNamedUserDTOFilename, funcName, constanta.UniqueID1)
		return
	}

	if input.AuthUserId < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(ValidationNamedUserDTOFilename, funcName, constanta.AuthUserID)
		return
	}

	return
}
