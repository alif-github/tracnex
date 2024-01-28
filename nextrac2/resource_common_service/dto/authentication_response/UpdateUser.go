package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type UpdateUserAuthenticationResponse struct {
	Nexsoft UpdateUserBodyResponse `json:"nexsoft"`
}

type UpdateUserBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload UpdateUserPayload  `json:"payload"`
}

type UpdateUserPayload struct {
	model.PayloadResponse
	Data UpdateUserData `json:"data"`
}

type UpdateUserData struct {
	Content UpdateUserContent `json:"content"`
}

type UpdateUserContent struct {
	EmailStatus EmailStatus `json:"email_status"`
	PhoneStatus PhoneStatus `json:"phone_status"`
}