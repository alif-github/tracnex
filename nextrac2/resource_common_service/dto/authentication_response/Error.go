package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type AuthenticationErrorResponse struct {
	Nexsoft RegisterUserBodyResponse `json:"nexsoft"`
}

type ErrorBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload ErrorPayload         `json:"payload"`
}

type ErrorPayload struct {
	model.PayloadResponse
}
