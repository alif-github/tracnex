package authentication_request

import (
	"nexsoft.co.id/nextrac2/resource_common_service/model"
)

type UserAuthenticationDTO struct {
	ClientID              string                        `json:"client_id"`
	Username              string                        `json:"username"`
	Password              string                        `json:"password"`
	FirstName             string                        `json:"first_name"`
	LastName              string                        `json:"last_name"`
	Email                 string                        `json:"email"`
	CountryCode           string                        `json:"country_code"`
	Phone                 string                        `json:"phone"`
	Device                string                        `json:"device"`
	Locale                string                        `json:"locale"`
	EmailMessage          string                        `json:"email_message"`
	IPWhitelist           string                        `json:"ip_whitelist"`
	EmailLinkMessage      string                        `json:"email_link_message"`
	PhoneMessage          string                        `json:"phone_message"`
	ResourceID            string                        `json:"resource_id"`
	AdditionalInformation []model.AdditionalInformation `json:"additional_information"`
	UpdatedAt             string                        `json:"updated_at"`
}

type GetListUserDTO struct {
	ListID []int64 `json:"list_id"`
}
