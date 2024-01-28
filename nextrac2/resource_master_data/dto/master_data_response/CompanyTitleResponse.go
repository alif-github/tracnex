package master_data_response

import "time"

type CompanyTitleResponse struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	CreatedBy     int64     `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
}