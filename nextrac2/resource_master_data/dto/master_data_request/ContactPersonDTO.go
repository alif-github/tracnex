package master_data_request

import (
	"nexsoft.co.id/nextrac2/dto/in"
	"time"
)

type ContactPersonGetListRequest struct {
	in.AbstractDTO
	ID        int64
	ParentID  int64  `json:"parent_id"`
	Connector string `json:"connector"`
	NIK       string `json:"nik"`
}

type ContactPersonWriteRequest struct {
	ID              int64     `json:"id"`
	PersonProfileID int64     `json:"person_profile_id"`
	PersonTitleID   int64     `json:"person_title_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	NIK             string    `json:"nik"`
	Address1        string    `json:"address_1"`
	ParentID        int64     `json:"parent_id"`
	PositionID      int64     `json:"position_id"`
	Connector       string    `json:"connector"`
	Email           string    `json:"email"`
	PhoneCode       string    `json:"phone_code"`
	Phone           string    `json:"phone"`
	JoinDate        time.Time `json:"join_date"`
	ResignDate      time.Time `json:"resign_date"`
	SuperiorID      int64     `json:"superior_id"`
	Status          string    `json:"status"`
	UpdatedAt       time.Time `json:"updated_at"`
}
