package out

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type AddResourceNexcloudResponse struct {
	Nexsoft AddResourceNexcloudBodyResponse `json:"nexsoft"`
}

type AddResourceNexcloudBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload AddResourceNexcloudPayload  `json:"payload"`
}

type AddResourceNexcloudPayload struct {
	model.PayloadResponse
}

type ViewDetailLogForAddResourceDTOOut struct {
	ClientID	string	`json:"client_id"`
	AuthUserID 	int64	`json:"auth_user_id"`
	ClientType	string	`json:"client_type"`
	Status		string	`json:"status"`
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	Resource	string	`json:"resource"`
	UpdatedAt	string	`json:"updated_at"`
}
