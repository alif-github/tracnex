package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type AuthorizeAuthenticationResponse struct {
	Nexsoft AuthorizeBodyResponse `json:"nexsoft"`
}

type AuthorizeBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload AuthorizePayload     `json:"payload"`
}

type AuthorizePayload struct {
	model.PayloadResponse
}
