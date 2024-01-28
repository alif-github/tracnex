package out

import (
	"time"
)

type ProductLicenseResponse struct {
	ID             int64     `json:"id"`
	LicenseConfig  int64     `json:"license_config"`
	CustomerName   string    `json:"customer_name"`
	UniqueId1      string    `json:"unique_id_1"`
	UniqueId2      string    `json:"unique_id_2"`
	InstallationId int64     `json:"installation_id"`
	ProductName    string    `json:"product_name"`
	LicenseVariant string    `json:"license_variant"`
	LicenseType    string    `json:"license_type"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidThru      time.Time `json:"valid_thru"`
	Status         int32     `json:"status"`
}

type ProductLicenseDetailResponse struct {
	LicenseID              int64                                `json:"license_id"`
	ProductKey             string                               `json:"product_key"`
	ActivationDate         time.Time                            `json:"activation_date"`
	LicenseStatus          int32                                `json:"license_status"`
	TerminationDescription string                               `json:"termination_description"`
	LicenseConfigId        int64                                `json:"license_config_id"`
	InstallationId         int64                                `json:"installation_id"`
	ParentCustomerId       int64                                `json:"parent_customer_id"`
	ParentCustomer         string                               `json:"parent_customer"`
	CustomerId             int64                                `json:"customer_id"`
	SiteId                 int64                                `json:"site_id"`
	Customer               string                               `json:"customer"`
	ClientId               string                               `json:"client_id"`
	Product                string                               `json:"product"`
	Client                 string                               `json:"client"`
	LicenseVariant         string                               `json:"license_variant"`
	LicenseType            string                               `json:"license_type"`
	DeploymentMethod       string                               `json:"deployment_method"`
	NumberOfUser           int64                                `json:"number_of_user"`
	ConcurentUser          string                               `json:"concurent_user"`
	UniqueId1              string                               `json:"unique_id_1"`
	UniqueId2              string                               `json:"unique_id_2"`
	LicenseValidFrom       time.Time                            `json:"license_valid_from"`
	LicenseValidThru       time.Time                            `json:"license_valid_thru"`
	Created                time.Time                            `json:"created"`
	Modified               time.Time                            `json:"modified"`
	ModifiedBy             string                               `json:"modified_by"`
	Module1                string                               `json:"module_1"`
	Module2                string                               `json:"module_2"`
	Module3                string                               `json:"module_3"`
	Module4                string                               `json:"module_4"`
	Module5                string                               `json:"module_5"`
	Module6                string                               `json:"module_6"`
	Module7                string                               `json:"module_7"`
	Module8                string                               `json:"module_8"`
	Module9                string                               `json:"module_9"`
	Module10               string                               `json:"module_10"`
	Components             []ListProductComponentProductLicense `json:"components"`
}

type ListProductComponentProductLicense struct {
	ComponentID    int64  `json:"component_id"`
	ComponentName  string `json:"component_name"`
	ComponentValue string `json:"component_value"`
}

type ListLicenseJournal struct {
	ID                int64  `json:"id"`
	ClientID          string `json:"client_id"`
	LicenseStatus     int64  `json:"license_status"`
	StatusDescription string `json:"status_description"`
	UniqueID1         string `json:"unique_id_1"`
	UniqueID2         string `json:"unique_id_2"`
	ProductName       string `json:"product_name"`
	ClientType        string `json:"client_type"`
	AllowActivation   string `json:"allow_activation"`
	NoOfUser          int64  `json:"no_of_user"`
	ProductValidFrom  string `json:"product_valid_from"`
	ProductValidThru  string `json:"product_valid_thru"`
	IsUserConcurrent  string `json:"is_user_concurrent"`
	TotalLicense      int64  `json:"total_license"`
	TotalActivated    int64  `json:"total_activated"`
}
