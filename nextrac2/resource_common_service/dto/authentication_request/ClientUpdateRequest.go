package authentication_request

import model2 "nexsoft.co.id/nextrac2/resource_common_service/model"

type ClientUpdateRequest struct {
	Scope                string                         `json:"scope"`
	RedirectUri          string                         `json:"redirect_uri"`
	IPWhitelist          string                         `json:"ip_whitelist"`
	AccessTokenValidity  int64                          `json:"access_token_validity"`
	RefreshTokenValidity int64                          `json:"refresh_token_validity"`
	MultipleLogin        bool                           `json:"multiple_login"`
	MaxAuthFail          int                            `json:"max_auth_fail"`
	Locale               string                         `json:"locale"`
	ClientInformation    []model2.AdditionalInformation `json:"client_information"`
	UpdatedAtString      string                         `json:"updated_at"`
}
