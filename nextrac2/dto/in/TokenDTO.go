package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type TokenDTOIn struct {
	AuthorizationCode string `json:"authorization_code"`
	CodeVerifier      string `json:"code_verifier"`
}

type TokenClientDTOIn struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	WantedScope  string `json:"wanted_scope"`
}

func (input TokenDTOIn) ValidateToken() errorModel.ErrorModel {
	fileName := "TokenDTO.go"
	funcName := "ValidateToken"
	if input.AuthorizationCode == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.AuthorizationCode)
	}

	if input.CodeVerifier == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CodeVerifier)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input TokenClientDTOIn) ValidateTokenClient() errorModel.ErrorModel {
	fileName := "TokenDTO.go"
	funcName := "ValidateTokenClient"
	if util.IsStringEmpty(input.ClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientID)
	}

	if util.IsStringEmpty(input.ClientSecret) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientSecret)
	}

	if util.IsStringEmpty(input.GrantType) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.GrantTypes)
	}

	if util.IsStringEmpty(input.WantedScope) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Wanted Scope")
	}

	return errorModel.GenerateNonErrorModel()
}