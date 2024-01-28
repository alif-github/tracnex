package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type DeleteUserAuthenticationResponse struct {
	Nexsoft AuthorizeBodyResponse `json:"nexsoft"`
}

type DeleteUserBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload AuthorizePayload     `json:"payload"`
}

type DeleteUserPayload struct {
	model.PayloadResponse
}
