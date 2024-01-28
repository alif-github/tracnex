package authentication_response

import (
	"nexsoft.co.id/nextrac2/resource_common_service/model"
)

type UserAuthenticationResponse struct {
	Nexsoft UserBodyResponse `json:"nexsoft"`
}

type UserBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload UserPayload          `json:"payload"`
}

type UserPayload struct {
	model.PayloadResponse
	Data UserData `json:"data"`
}

type UserData struct {
	Content UserContent `json:"content"`
}

type UserContent struct {
	AliasName         string                        `json:"alias_name"`
	ID                int64                         `json:"id"`
	UserID            int64                         `json:"user_id"`
	Username          string                        `json:"username"`
	FirstName         string                        `json:"first_name"`
	LastName          string                        `json:"last_name"`
	Email             string                        `json:"email"`
	Phone             string                        `json:"phone"`
	ClientID          string                        `json:"client_id"`
	MaxAuthFail       int                           `json:"max_auth_fail"`
	Locale            string                        `json:"locale"`
	ResourceID        string                        `json:"resource_id"`
	UserStatus        int                           `json:"user_status"`
	SignatureKey      string                        `json:"signature_key"`
	GrantTypes        string                        `json:"grant_types"`
	Scope             string                        `json:"scope"`
	IPWhitelist       string                        `json:"ip_whitelist"`
	RedirectURI       string                        `json:"redirect_uri"`
	ClientInformation string                        `json:"client_information"`
	UserInformation   []model.AdditionalInformation `json:"user_information"`
	UpdatedAt         string                        `json:"updated_at"`
}
