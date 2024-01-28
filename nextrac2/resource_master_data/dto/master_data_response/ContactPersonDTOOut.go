package master_data_response

import "time"

type ContactPersonResponse struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	NIK             string    `json:"nik"`
	PersonProfileID int64     `json:"person_profile_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PositionID      int64     `json:"position_id"`
	Position        string    `json:"position"`
	Connector       string    `json:"connector"`
	ParentID        int64     `json:"parent_id"`
	CreatedBy       int64     `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ViewContactPersonResponse struct {
	ID                int64     `json:"id"`
	NIK               string    `json:"nik"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	PersonProfileID   int64     `json:"person_profile_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	PositionID        int64     `json:"position_id"`
	Position          string    `json:"position"`
	ParentID          int64     `json:"parent_id"`
	Connector         string    `json:"connector"`
	JoinDate          time.Time `json:"join_date"`
	ResignDate        time.Time `json:"resign_date"`
	SuperiorID        int64     `json:"superior_id"`
	SuperiorFirstName string    `json:"superior_first_name"`
	SuperiorLastName  string    `json:"superior_last_name"`
	Status            string    `json:"status"`
	CreatedBy         int64     `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
}
