package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type RegisterUserAuthenticationResponse struct {
	Nexsoft RegisterUserBodyResponse `json:"nexsoft"`
}

type RegisterUserBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload RegisterUserPayload  `json:"payload"`
}

type RegisterUserPayload struct {
	model.PayloadResponse
	Data RegisterUserData `json:"data"`
}

type RegisterUserData struct {
	Content RegisterUserContent `json:"content"`
}

type RegisterUserContent struct {
	UserID       int64                    `json:"user_id"`
	ClientID     string                   `json:"client_id"`
	SignatureKey string                   `json:"signature_key"`
	NotifyStatus StatusNotifyEmailOrPhone `json:"notify_status"`
}

type StatusNotifyEmailOrPhone struct {
	EmailStatus EmailStatus `json:"email_status"`
	PhoneStatus PhoneStatus `json:"phone_status"`
}

type EmailStatus struct {
	EmailNotify        bool   `json:"email_notify"`
	EmailNotifyStatus  bool   `json:"email_notify_status"`
	EmailNotifyMessage string `json:"email_notify_message"`
}

type PhoneStatus struct {
	PhoneNotify        bool   `json:"phone_notify"`
	PhoneNotifyStatus  bool   `json:"phone_notify_status"`
	PhoneNotifyMessage string `json:"phone_notify_message"`
}
