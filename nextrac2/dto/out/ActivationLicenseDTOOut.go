package out

type LicenseResponse struct {
	ProductKey       string `json:"product_key"`
	ProductEncrypt   string `json:"product_encrypt"`
	ProductSignature string `json:"product_signature"`
	ClientTypeID     int64  `json:"client_type_id"`
	UniqueID1        string `json:"unique_id_1"`
	UniqueID2        string `json:"unique_id_2"`
}

type LicenseResponseWithSalesman struct {
	LicenseResponse
	SalesmanList []SalesmanLicenseListOut `json:"list_salesman_status"`
}

type SalesmanLicenseListOut struct {
	ID         int64  `json:"salesman_id"`
	AuthUserID int64  `json:"salesman_auth_id"`
	UserID     string `json:"user_id"`
	Status     string `json:"salesman_status"`
}

type ActivationLicenseErrorDetail struct {
	UniqueID1       string `json:"unique_id_1"`
	UniqueID2       string `json:"unique_id_2"`
	LicenseConfigID int64  `json:"license_config_id"`
	Message         string `json:"message"`
}
