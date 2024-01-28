package in

import "nexsoft.co.id/nextrac2/model/errorModel"

type GroChatLoginDTOIn struct {
	RequestId         string `json:"request_id"`
	AuthorizationCode string `json:"authorization_code"`
	CodeVerifier      string `json:"code_verifier"`
}

func (input *GroChatLoginDTOIn) ValidateGroChatLoginDTO() errorModel.ErrorModel {
	var (
		fileName = "GroChatLoginDTO.go"
		funcName = "ValidateGroChatLoginDTO"
	)

	if input.RequestId == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "request_id")
	}

	if input.AuthorizationCode == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "authorization_code")
	}

	if input.CodeVerifier == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "code_verifier")
	}

	return errorModel.GenerateNonErrorModel()
}