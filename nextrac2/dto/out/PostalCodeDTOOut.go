package out

import "time"

type PostalCodeResponse struct {
	ID             int64     `json:"id"`
	UrbanVillageID int64     `json:"urban_village_id"`
	Code           string    `json:"code"`
	Status         string    `json:"status"`
	CreatedBy      int64     `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PostalCodeDetailResponse struct {
	ID               int64     `json:"id"`
	UrbanVillageID   int64     `json:"urban_village_id"`
	UrbanVillageName string    `json:"urban_village_name"`
	Code             string    `json:"code"`
	Status           string    `json:"status"`
	CreatedBy        int64     `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedBy        int64     `json:"updated_by"`
	UpdatedAt        time.Time `json:"updated_at"`
}