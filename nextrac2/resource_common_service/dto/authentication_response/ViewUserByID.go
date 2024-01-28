package authentication_response

import (
	"nexsoft.co.id/nextrac2/resource_common_service/model"
)

type ViewUserByIDAuthenticationResponse struct {
	Nexsoft ViewUserByIDBodyResponse `json:"nexsoft"`
}

type ViewUserByIDBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload ViewUserByIDPayload  `json:"payload"`
}

type ViewUserByIDPayload struct {
	model.PayloadResponse
	Data ViewUserByIDData `json:"data"`
}

type ViewUserByIDData struct {
	Content ViewUserByIDContent `json:"content"`
}

type ViewUserByIDContent struct {
	ID              int64                         `json:"id"`
	ClientID        string                        `json:"client_id"`
	Username        string                        `json:"username"`
	FirstName       string                        `json:"first_name"`
	LastName        string                        `json:"last_name"`
	Email           string                        `json:"email"`
	Phone           string                        `json:"phone"`
	Locale          string                        `json:"locale"`
	Status          int32                         `json:"status"`
	RoleID          string                        `json:"role_id"`
	GroupID         string                        `json:"group_id"`
	ResourceIDS     string                        `json:"resource_ids"`
	UserInformation []model.AdditionalInformation `json:"user_information"`
	UpdatedAt       string                        `json:"updated_at"`
	CreatedBy       int64                         `json:"created_by"`
}
