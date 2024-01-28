package authentication_response

import (
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
)

type RegisterClientAuthenticationResponse struct {
	Nexsoft RegisterClientBodyResponse `json:"nexsoft"`
}

type RegisterClientBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload RegisterClientPayload  `json:"payload"`
}

type RegisterClientPayload struct {
	model.PayloadResponse
	Data RegisterClientData `json:"data"`
}

type RegisterClientData struct {
	Content RegisterClientContent `json:"content"`
}

type RegisterClientContent struct {
	ClientID     	string                   `json:"client_id"`
	ClientSecret    string                   `json:"client_secret"`
	SignatureKey 	string                   `json:"signature_key"`
	ResourceList	[]out.ResourceList		`json:"resource_list"`
}
