package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type CheckTokenAuthenticationResponse struct {
	Nexsoft CheckTokenBodyResponse `json:"nexsoft"`
}

type CheckTokenBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload CheckTokenPayload    `json:"payload"`
}

type CheckTokenPayload struct {
	model.PayloadResponse
	Data CheckTokenData `json:"data"`
}

type CheckTokenData struct {
	Content CheckTokenContent `json:"content"`
}

type CheckTokenContent struct {
	Authentication string `json:"authentication"`
	IPWhitelist    string `json:"ip_whitelist"`
}
