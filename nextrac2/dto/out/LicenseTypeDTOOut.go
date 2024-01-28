package out

import "time"

type LicenseTypeResponse struct {
	ID              int64     `json:"id"`
	LicenseTypeName string    `json:"license_type_name"`
	LicenseTypeDesc string    `json:"license_type_desc"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       int64     `json:"updated_by"`
	UpdatedName     string    `json:"updated_name"`
}
