package out

import "time"

type ProvinceResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int64     `json:"updated_by"`
}

type ProvinceLocalResponse struct {
	ID         int64  `json:"id"`
	CountryID  int64  `json:"country_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	DistrictID []int  `json:"district_id"`
}
