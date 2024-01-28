package out

import "time"

type SubDistrictResponse struct {
	ID           int64     `json:"id"`
	DistrictID   int64     `json:"district_id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CreatedBy    int64     `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SubDistrictDetailResponse struct {
	ID           int64     `json:"id"`
	DistrictID   int64     `json:"district_id"`
	DistrictName string    `json:"district_name"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CreatedBy    int64     `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedBy    int64     `json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}
