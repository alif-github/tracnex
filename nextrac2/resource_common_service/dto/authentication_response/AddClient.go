package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type AddClientAuthenticationResponse struct {
	Nexsoft AddClientBodyResponse `json:"nexsoft"`
}

type AddClientBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload AddClientPayload     `json:"payload"`
}

type AddClientPayload struct {
	model.PayloadResponse
	Data AddClientData `json:"data"`
}

type AddClientData struct {
	Content AddClientContent `json:"content"`
}

type AddClientContent struct {
	ClientID          string `json:"client_id"`
	UserID            int64  `json:"user_id"`
	Username          string `json:"username"`
	SignatureKey      string `json:"signature_key"`
	GrantTypes        string `json:"grant_types"`
	ResourceID        string `json:"resource_id"`
	IPWhitelist       string `json:"ip_whitelist"`
	UserStatus        int    `json:"user_status"`
	Scope             string `json:"scope"`
	Locale            string `json:"locale"`
	ClientInformation string `json:"client_information"`
	UserInformation   string `json:"user_information"`
}
