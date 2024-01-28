package out

import "time"

type ListLicenseConfigModel struct {
	LicenseConfigID    int64     `json:"license_config_id"`
	CustomerName       string    `json:"customer_name"`
	UniqueID1          string    `json:"unique_id_1"`
	UniqueID2          string    `json:"unique_id_2"`
	InstallationID     int64     `json:"installation_id"`
	ProductName        string    `json:"product_name"`
	ClientTypeID       int64     `json:"client_type_id"`
	LicenseVariantName string    `json:"license_variant_name"`
	LicenseTypeName    string    `json:"license_type_name"`
	ProductValidFrom   string    `json:"product_valid_from"`
	ProductValidThru   string    `json:"product_valid_thru"`
	AllowActivation    string    `json:"allow_activation"`
	PaymentStatus      string    `json:"payment_status"`
	IsExtendChecklist  bool      `json:"is_extend_checklist"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ViewDetailLicenseConfig struct {
	ID                 int64                  `json:"id"`
	InstallationID     int64                  `json:"installation_id"`
	ParentCustomerID   int64                  `json:"parent_customer_id"`
	ParentCustomer     string                 `json:"parent_customer"`
	CustomerID         int64                  `json:"customer_id"`
	Customer           string                 `json:"customer"`
	SiteID             int64                  `json:"site_id"`
	ProductID          int64                  `json:"product_id"`
	ProductName        string                 `json:"product_name"`
	ClientID           string                 `json:"client_id"`
	LicenseVariantName string                 `json:"license_variant_name"`
	LicenseTypeName    string                 `json:"license_type_name"`
	DeploymentMethod   string                 `json:"deployment_method"`
	NoOfUser           int64                  `json:"no_of_user"`
	IsUserConcurrent   string                 `json:"is_user_concurrent"`
	UniqueID1          string                 `json:"unique_id_1"`
	UniqueID2          string                 `json:"unique_id_2"`
	ProductValidFrom   string                 `json:"product_valid_from"`
	ProductValidThru   string                 `json:"product_valid_thru"`
	MaxOfflineDays     int64                  `json:"max_offline_days"`
	ClientTypeName     string                 `json:"client_type_name"`
	AllowActivation    string                 `json:"allow_activation"`
	ModuleID1          string                 `json:"module_id_1"`
	ModuleID2          string                 `json:"module_id_2"`
	ModuleID3          string                 `json:"module_id_3"`
	ModuleID4          string                 `json:"module_id_4"`
	ModuleID5          string                 `json:"module_id_5"`
	ModuleID6          string                 `json:"module_id_6"`
	ModuleID7          string                 `json:"module_id_7"`
	ModuleID8          string                 `json:"module_id_8"`
	ModuleID9          string                 `json:"module_id_9"`
	ModuleID10         string                 `json:"module_id_10"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	UpdatedName        string                 `json:"updated_name"`
	ProductComponent   []ListProductComponent `json:"product_component"`
}

type ListLicenseConfigIDs struct {
	TotalLicenseConfigID int64   `json:"total_license_config_id"`
	LicenseConfigID      []int64 `json:"license_config_id"`
}
