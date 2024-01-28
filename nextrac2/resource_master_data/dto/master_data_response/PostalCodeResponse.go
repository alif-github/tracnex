package master_data_response

import "time"

type PostalCodeResponse struct {
	ID               int64     `json:"id"`
	UrbanVillageID   int64     `json:"urban_village_id"`
	UrbanVillageName string    `json:"urban_village_name"`
	Code             string    `json:"code"`
	Status           string    `json:"status"`
	CreatedBy        int64     `json:"created_by"`
	UpdatedAtStr     string    `json:"updated_at"`
	UpdatedAt        time.Time `json:"-"`
}
