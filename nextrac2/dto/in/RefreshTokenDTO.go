package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type RefreshTokenDTOIn struct {
	AbstractDTO
	Authorization string
	RefreshToken  string `json:"refresh_token"`
}

func (input RefreshTokenDTOIn) ValidateRefreshToken() errorModel.ErrorModel {
	fileName := "RefreshTokenDTO.go"
	funcName := "ValidateRefreshToken"
	if input.Authorization == "" {
		return errorModel.GenerateUnauthorizedClientError(fileName, funcName)
	}

	if input.RefreshToken == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.RefreshToken)
	}


	err := input.ValidateMinMaxString(input.RefreshToken, constanta.RefreshToken, 1, 256)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}
