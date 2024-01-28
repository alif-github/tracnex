package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type RequestChangePasswordAuthenticationResponse struct {
	Nexsoft RequestChangePasswordBodyResponse 	`json:"nexsoft"`
}

type RequestChangePasswordBodyResponse struct {
	Header  model.HeaderResponse 				`json:"header"`
	Payload RequestChangePasswordPayload  		`json:"payload"`
}

type RequestChangePasswordPayload struct {
	model.PayloadResponse
	Data RequestChangePasswordData 				`json:"data"`
}

type RequestChangePasswordData struct {
	Content RequestChangePasswordContent 		`json:"content"`
}

type RequestChangePasswordContent struct {
	UserID         string                   	`json:"user_id"`
	ClientID       string                   	`json:"client_id"`
	UserName       string                   	`json:"username"`
}