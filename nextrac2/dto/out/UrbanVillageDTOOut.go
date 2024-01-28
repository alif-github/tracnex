package out

import "time"

type UrbanVillageResponse struct {
	ID              int64     `json:"id"`
	SubDistrictID   int64     `json:"sub_district_id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	CreatedBy       int64     `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UrbanVillageDetailResponse struct {
	ID              int64     `json:"id"`
	SubDistrictID   int64     `json:"sub_district_id"`
	SubDistrictName string    `json:"sub_district_name"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	CreatedBy       int64     `json:"created_by"`
	CreatedAt 		time.Time `json:"created_at"`
	UpdatedBy       int64     `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}
