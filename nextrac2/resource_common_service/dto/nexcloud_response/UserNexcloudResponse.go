package nexcloud_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type UserNexcloudResponse struct {
	Nexsoft UserNexcloudBodyResponse `json:"nexsoft"`
}

type UserNexcloudBodyResponse struct {
	Header 		model.HeaderResponse 		`json:"header"`
	Payload 	UserNexcloudPayload			`json:"payload"`
}

type UserNexcloudPayload struct {
	model.PayloadResponse
	Data UserNexcloudData `json:"data"`
}

type UserNexcloudData struct {
	Content UserNexcloudContent `json:"content"`
}

type UserNexcloudContent struct {
	ID			string	`json:"id"`
	ClientID	string	`json:"client_id"`
	AuthUserID 	int64	`json:"auth_user_id"`
	UpdatedAt	string	`json:"updated_at"`
	CreatedBy	int64	`json:"created_by"`
}