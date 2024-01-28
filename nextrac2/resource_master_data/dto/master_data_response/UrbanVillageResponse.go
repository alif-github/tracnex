package master_data_response

import "time"

type UrbanVillageResponse struct {
	ID              int64     `json:"id"`
	SubDistrictID   int64     `json:"sub_district_id"`
	SubDistrictName string    `json:"sub_district_name"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	CreatedBy       int64     `json:"created_by"`
	UpdatedAtStr    string    `json:"updated_at"`
	UpdatedAt       time.Time `json:"-"`
}
