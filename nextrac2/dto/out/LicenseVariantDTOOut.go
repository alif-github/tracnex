package out

import "time"

type LicenseVariantListResponse struct {
	ID                 int64     `json:"id"`
	LicenseVariantName string    `json:"license_variant_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedBy          int64     `json:"updated_by"`
	UpdatedName        string    `json:"updated_name"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type LicenseVariantViewResponse struct {
	ID                 int64     `json:"id"`
	LicenseVariantName string    `json:"license_variant_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedName        string    `json:"updated_name"`
}
