package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type UnregisterNamedUserRequest struct {
	ID int64 `json:"id"`
}

func (input *UnregisterNamedUserRequest) ValidateUnregisterNamedUser() (err errorModel.ErrorModel) {
	fileName := "UserRegistrationDetailDTO.go"
	funcName := "ValidateUnregisterNamedUser"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.UserRegistrationDetailID)
	}

	return errorModel.GenerateNonErrorModel()
}
