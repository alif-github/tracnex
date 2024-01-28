package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type CheckClientByViewResponse struct {
	Nexsoft CheckClientByViewBodyResponse `json:"nexsoft"`
}

type CheckClientByViewBodyResponse struct {
	Header  model.HeaderResponse     `json:"header"`
	Payload CheckClientByViewPayload `json:"payload"`
}

type CheckClientByViewPayload struct {
	model.PayloadResponse
	Data CheckClientByViewContent `json:"data"`
}

type CheckClientByViewContent struct {
	Content CheckClientByView `json:"content"`
}

type CheckClientByView struct {
	AliasName            string              `json:"alias_name"`
	ClientID             string              `json:"client_id"`
	ResourceIDS          string              `json:"resource_ids"`
	ClientSecret         string              `json:"client_secret"`
	SignatureKey         string              `json:"signature_key"`
	Scope                string              `json:"scope"`
	GrantTypes           string              `json:"grant_types"`
	RedirectUri          string              `json:"redirect_uri"`
	IpWhitelist          string              `json:"ip_whitelist"`
	Authorities          string              `json:"authorities"`
	AccessTokenValidity  int64               `json:"access_token_validity"`
	RefreshTokenValidity int64               `json:"refresh_token_validity"`
	MultipleLogin        bool                `json:"multiple_login"`
	Locale               string              `json:"locale"`
	MaxAuthFail          int64               `json:"max_auth_fail"`
	ClientInformation    []ClientInformation `json:"client_information"`
	UpdatedAt            string              `json:"updated_at"`
	CreatedBy            int64               `json:"created_by"`
}

type ClientInformation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
