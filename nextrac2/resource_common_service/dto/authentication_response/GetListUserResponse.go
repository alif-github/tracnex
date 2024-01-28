package authentication_response

import (
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"time"
)

type GetListUserAuthenticationResponse struct {
	Nexsoft GetListUserBodyResponse `json:"nexsoft"`
}

type GetListUserBodyResponse struct {
	Header  model.HeaderResponse `json:"header"`
	Payload GetListUserPayload   `json:"payload"`
}

type GetListUserPayload struct {
	model.PayloadResponse
	Data GetListUserData `json:"data"`
}

type GetListUserData struct {
	Content []GetListUserContent `json:"content"`
}

type GetListUserContent struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
}
