package out

type CheckUserLicenseNamedUserResponse struct {
	ID             int64  `json:"id"`
	ProductKey     string `json:"product_key"`
	TotalLicense   int64  `json:"total_license"`
	TotalActivated int64  `json:"total_activated"`
	QuotaLicense   int64  `json:"quota_license"`
}
