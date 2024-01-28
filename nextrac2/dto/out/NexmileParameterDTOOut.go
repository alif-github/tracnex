package out

import "time"

type ViewNexmileParameterResponse struct {
	UniqueId1         string                   `json:"unique_id_1"`
	CompanyName       string                   `json:"company_name"`
	ProductValidaFrom time.Time                `json:"product_valida_from"`
	ProductValidThru  time.Time                `json:"product_valid_thru"`
	PasswordAdmin     string                   `json:"password_admin"`
	UniqueId2         string                   `json:"unique_id_2"`
	UserAdmin         string                   `json:"user_admin"`
	LicenseStatus     int64                    `json:"license_status"`
	MaxOfflineDays    int64                    `json:"max_offline_days"`
	Parameters        []ParameterValueResponse `json:"parameters"`
}

type ParameterValueResponse struct {
	ParameterID    string `json:"parameter_id"`
	ParameterValue string `json:"parameter_value"`
}
