package master_data_response

import "time"

type DistrictResponse struct {
	ID             int64     `json:"id"`
	ProvinceID     int64     `json:"province_id"`
	ProvinceName   string    `json:"province_name"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	KemendagriCode string    `json:"kemendagri_code"`
	Status         string    `json:"status"`
	CreatedBy      int64     `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}
