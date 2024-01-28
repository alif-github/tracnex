package in

import "nexsoft.co.id/nextrac2/model/errorModel"

type UserInvitationRequest struct {
	Id			int64  `json:"-"`
	Email       string `json:"email"`
	RoleId      int64  `json:"role_id"`
	DataGroupId int64  `json:"data_group_id"`
}

func (input UserInvitationRequest) ValidateInsert() errorModel.ErrorModel {
	fileName := "UserInvitationDTO.go"
	funcName := "ValidateInsert"

	if input.Email == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "email")
	}

	if input.RoleId <= 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "role_id")
	}

	if input.DataGroupId <= 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "data_group_id")
	}

	return errorModel.GenerateNonErrorModel()
}