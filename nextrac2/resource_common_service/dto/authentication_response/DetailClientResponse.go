package authentication_response

import (
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"time"
)

type DetailClientResponse struct {
	ClientID             string                         `json:"client_id"`
	ClientSecret         string                         `json:"client_secret"`
	ResourceIDS          string                         `json:"resource_ids"`
	Scope                string                         `json:"scope"`
	GrantTypes           string                         `json:"grant_types"`
	RedirectUri          string                         `json:"redirect_uri"`
	IPWhitelist          string                         `json:"ip_whitelist"`
	AccessTokenValidity  int64                          `json:"access_token_validity"`
	RefreshTokenValidity int64                          `json:"refresh_token_validity"`
	MultipleLogin        bool                           `json:"multiple_login"`
	MaxAuthFail          int                            `json:"max_auth_fail"`
	Locale               string                         `json:"locale"`
	Group                string                         `json:"group"`
	Role                 string                         `json:"role"`
	ClientInformation    []model2.AdditionalInformation `json:"client_information"`
	UpdatedAtString      string                         `json:"updated_at"`
	UpdatedAt            time.Time						`json:"-"`
	Username             string                         `json:"username"`
	Email                string                         `json:"email"`
	Phone                string                         `json:"phone"`
	AliasName            string                         `json:"alias_name"`
	CreatedClient        string                         `json:"created_client"`
	UpdatedClient        string                         `json:"updated_client"`
}
