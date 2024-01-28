package master_data_response

import "time"

type SubDistrictResponse struct {
	ID           int64     `json:"id"`
	DistrictID   int64     `json:"district_id"`
	DistrictName string    `json:"district_name"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CreatedBy    int64     `json:"created_by"`
	UpdatedAtStr string    `json:"updated_at"`
	UpdatedAt    time.Time `json:"-"`
}
