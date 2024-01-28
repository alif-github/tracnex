package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type InitiateInsertUserAuthenticationResponse struct {
	Nexsoft InitiateInsertUserBodyResponse `json:"nexsoft"`
}

type InitiateInsertUserBodyResponse struct {
	Header  model.HeaderResponse      `json:"header"`
	Payload InitiateInsertUserPayload `json:"payload"`
}

type InitiateInsertUserPayload struct {
	model.PayloadResponse
	Data InitiateInsertUserData `json:"data"`
}

type InitiateInsertUserData struct {
	Content InitiateInsertUserContent `json:"content"`
}

type InitiateInsertUserContent struct {
	Country []struct {
		CountryName string `json:"country_name"`
		CountryCode string `json:"country_code"`
	} `json:"country"`
	LocaleLanguage []string `json:"locale"`
}
