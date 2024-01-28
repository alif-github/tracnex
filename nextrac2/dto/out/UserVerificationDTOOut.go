package out

import "github.com/Azure/go-autorest/autorest/date"

type UserVerificationResponse struct {
	CompanyID        string                   `json:"company_id"`
	BranchID         string                   `json:"branch_id"`
	CompanyName      string                   `json:"company_name"`
	ProductValidFrom date.Date                `json:"product_valid_from"`
	ProductValidThru date.Date                `json:"product_valid_thru"`
	LicenseStatus    int64                    `json:"license_status"`
	AdminPassword    string                   `json:"admin_password"`
	AdminUsername    string                   `json:"admin_username"`
	MaxOfflineDays   int64                    `json:"max_offline_days"`
	ParameterValue   []ParameterValueResponse `json:"parameters"`
}
