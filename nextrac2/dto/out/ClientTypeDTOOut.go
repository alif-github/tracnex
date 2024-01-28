package out

import "time"

type InitiateClientTypeResponse struct {
	ClientTypeList []ClientTypeResponse `json:"client_type"`
}

type ClientTypeResponse struct {
	ClientTypeID int64  `json:"client_type_id"`
	ClientType   string `json:"client_type"`
	Description  string `json:"description"`
}

type ListClientTypeResponse struct {
	ID                 int64     `json:"id"`
	ClientType         string    `json:"client_type"`
	ParentClientTypeID int64     `json:"parent_client_type_id"`
	Description        string    `json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedName        string    `json:"updated_name"`
}
