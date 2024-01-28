package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type AddResourceClientAuthenticationResponse struct {
	Nexsoft AddResourceClientBodyResponse `json:"nexsoft"`
}

type AddResourceClientBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload AddResourceClientPayload  `json:"payload"`
}

type AddResourceClientPayload struct {
	model.PayloadResponse
	Data AddResourceClientData `json:"data"`
}

type AddResourceClientData struct {
	Content AddResourceClientContent `json:"content"`
}

type AddResourceClientContent struct {
	ClientID     	string           `json:"client_id"`
	ClientSecret    string           `json:"client_secret"`
	SignatureKey 	string           `json:"signature_key"`
}
