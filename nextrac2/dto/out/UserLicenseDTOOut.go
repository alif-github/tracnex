package out

import "time"

type UserLicenseResponse struct {
	ID              int64  `json:"license_id"`
	LicenseConfigId int64  `json:"license_config_id"`
	CustomerName    string `json:"customer_name"`
	UniqueId1       string `json:"unique_id_1"`
	UniqueId2       string `json:"unique_id_2"`
	InstallationId  int64  `json:"installation_id"`
	TotalLicense    int64  `json:"total_license"`
	TotalActivated  int64  `json:"active"`
}

type UserLicenseDetailResponse struct {
	ID                      int64                            `json:"license_id"`
	LicenseConfigId         int64                            `json:"license_config_id"`
	InstallationId          int64                            `json:"installation_id"`
	ParentCustomerId        int64                            `json:"parent_customer_id"`
	ParentCustomer          string                           `json:"parent_customer"`
	CustomerId              int64                            `json:"customer_id"`
	SiteId                  int64                            `json:"site_id"`
	CustomerName            string                           `json:"customer"`
	UniqueId1               string                           `json:"company_id"`
	UniqueId2               string                           `json:"branch_id"`
	ProductName             string                           `json:"product"`
	TotalActivated          int64                            `json:"active_user"`
	TotalLicense            int64                            `json:"of"`
	LicenseValidFrom        time.Time                        `json:"license_valid_from"`
	LicenseValidThru        time.Time                        `json:"license_valid_thru"`
	UpdatedAt               time.Time                        `json:"updated_at"`
	UserRegistrationDetails []UserRegistrationDetailResponse `json:"user_registration_details"`
}

type UserLicenseTransferKeyResponse struct {
	ID               int64     `json:"license_id"`
	LicenseConfigId  int64     `json:"license_config_id"`
	InstallationId   int64     `json:"installation_id"`
	ParentCustomerId int64     `json:"parent_customer_id"`
	ParentCustomer   string    `json:"parent_customer"`
	CustomerId       int64     `json:"customer_id"`
	SiteId           int64     `json:"site_id"`
	CustomerName     string    `json:"customer"`
	UniqueId1        string    `json:"company_id"`
	UniqueId2        string    `json:"branch_id"`
	ProductName      string    `json:"product"`
	TotalActivated   int64     `json:"active_user"`
	TotalLicense     int64     `json:"of"`
	LicenseValidFrom time.Time `json:"license_valid_from"`
	LicenseValidThru time.Time `json:"license_valid_thru"`
	UpdatedAt        time.Time `json:"updated_at"`
	ClientTypeId     int64     `json:"client_type_id"`
}

type UserRegistrationDetailResponse struct {
	UserRegistrationDetailID int64     `json:"user_registration_detail_id"`
	UserId                   string    `json:"user_id"`
	SalesmanId               string    `json:"salesman_id"`
	Email                    string    `json:"email"`
	NoTelp                   string    `json:"no_telp"`
	SalesmanCategory         string    `json:"salesman_category"`
	RegDate                  time.Time `json:"registration_date"`
	AndroidId                string    `json:"android_id"`
	Status                   string    `json:"status"`
}
