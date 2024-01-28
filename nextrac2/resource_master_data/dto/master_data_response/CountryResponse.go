package master_data_response

import "time"

type CountryResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	PhoneCode    string    `json:"phone_code"`
	Status       string    `json:"status"`
	CreatedBy    int64     `json:"created_by"`
	CreatedName  string    `json:"created_name"`
	UpdatedAtStr string    `json:"updated_at"`
	UpdatedAt    time.Time `json:"-"`
}
