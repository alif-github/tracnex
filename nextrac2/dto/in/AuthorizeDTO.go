package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type AuthorizeDTOIn struct {
	AbstractDTO
	CodeChallenger string `json:"code_challenger"`
}

func (input AuthorizeDTOIn) ValidateAuthorize() errorModel.ErrorModel {
	if input.CodeChallenger == "" {
		return errorModel.GenerateEmptyFieldError("AuthorizeDTO", "ValidateAuthorize", constanta.CodeChallenger)
	}

	return errorModel.GenerateNonErrorModel()
}
