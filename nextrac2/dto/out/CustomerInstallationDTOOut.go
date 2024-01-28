package out

import (
	"time"
)

type CustomerSiteInstallation struct {
	ParentCustomerID   int64          `json:"parent_customer_id"`
	ParentCustomerName string         `json:"parent_customer_name"`
	CustomerSite       []CustomerSite `json:"customer_site"`
}

type CustomerSite struct {
	SiteID               int64                  `json:"site_id"`
	CustomerID           int64                  `json:"customer_id"`
	CustomerSiteName     string                 `json:"customer_site_name"`
	Address              string                 `json:"address"`
	District             string                 `json:"district"`
	Province             string                 `json:"province"`
	Phone                string                 `json:"phone"`
	UpdatedAt            string                 `json:"updated_at"`
	CustomerInstallation []CustomerInstallation `json:"customer_installation"`
}

type CustomerInstallation struct {
	InstallationID     int64  `json:"installation_id"`
	ProductID          int64  `json:"product_id"`
	ProductGroupID     int64  `json:"product_group_id"`
	ProductName        string `json:"product_name"`
	ProductCode        string `json:"product_code"`
	Remark             string `json:"remark"`
	UniqueID1          string `json:"unique_id_1"`
	UniqueID2          string `json:"unique_id_2"`
	InstallationStatus string `json:"installation_status"`
	InstallationDate   string `json:"installation_date"`
	ProductValidFrom   string `json:"product_valid_from"`
	ProductValidThru   string `json:"product_valid_thru"`
	DayRange           int64  `json:"day_range"`
	UpdatedAt          string `json:"updated_at"`
}

type CustomerInstallationDetailList struct {
	InstallationID        int64     `json:"installation_id"`
	ProductID             int64     `json:"product_id"`
	ClientTypeID          int64     `json:"client_type_id"`
	ClientTypeDependantID int64     `json:"client_type_dependant_id"`
	ProductGroupID        int64     `json:"product_group_id"`
	ProductName           string    `json:"product_name"`
	ProductCode           string    `json:"product_code"`
	ProductDescription    string    `json:"product_description"`
	Remark                string    `json:"remark"`
	UniqueID1             string    `json:"unique_id_1"`
	UniqueID2             string    `json:"unique_id_2"`
	InstallationStatus    string    `json:"installation_status"`
	InstallationDate      time.Time `json:"installation_date"`
	ProductValidFrom      time.Time `json:"product_valid_from"`
	ProductValidThru      time.Time `json:"product_valid_thru"`
	DayRange              int64     `json:"day_range"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type CustomerInstallationDetailListForConfig struct {
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
	IsUserConcurrent   bool                   `json:"is_user_concurrent"`
	UniqueID1          string                 `json:"unique_id_1"`
	UniqueID2          string                 `json:"unique_id_2"`
	ProductValidFrom   time.Time              `json:"product_valid_from"`
	ProductValidThru   time.Time              `json:"product_valid_thru"`
	MaxOfflineDays     int64                  `json:"max_offline_days"`
	ClientTypeName     string                 `json:"client_type_name"`
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
	ProductComponent   []ListProductComponent `json:"product_component"`
}

type ListProductComponent struct {
	ComponentName  string `json:"component_name"`
	ComponentValue string `json:"component_value"`
}
