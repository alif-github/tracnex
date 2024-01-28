package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type TokenAuthenticationResponse struct {
	Nexsoft TokenBodyResponse `json:"nexsoft"`
}

type TokenBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload TokenPayload         `json:"payload"`
}

type TokenPayload struct {
	model.PayloadResponse
	Data TokenData `json:"data"`
}

type TokenData struct {
	Content TokenContent `json:"content"`
}

type TokenContent struct {
	RefreshToken string `json:"refresh_token"`
	State        string `json:"state"`
}
