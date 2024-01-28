package out

type ValidationLicenseErrorDetail struct {
	ProductKey      string `json:"product_key"`
	ProductEncrypt  string `json:"product_encrypt"`
	LicenseConfigID int64  `json:"license_config_id"`
	Message         string `json:"message"`
}
