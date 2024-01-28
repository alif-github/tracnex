package out

import "time"

type ListProduct struct {
	ID                    int64     `json:"id"`
	ProductName           string    `json:"product_name"`
	ProductDescription    string    `json:"product_description"`
	ProductGroupName      string    `json:"product_group_name"`
	ClientTypeName        string    `json:"client_type_name"`
	LicenseVariantName    string    `json:"license_variant_name"`
	LicenseTypeName       string    `json:"license_type_name"`
	ProductID             string    `json:"product_id"`
	ClientTypeID          int64     `json:"client_type_id"`
	ClientTypeDependantID int64     `json:"client_type_dependant_id"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ViewProduct struct {
	ID                 int64                  `json:"id"`
	ProductID          string                 `json:"product_id"`
	ProductName        string                 `json:"product_name"`
	ProductDescription string                 `json:"product_description"`
	ProductGroupID     int64                  `json:"product_group_id"`
	ProductGroupName   string                 `json:"product_group_name"`
	ClientTypeID       int64                  `json:"client_type_id"`
	ClientTypeName     string                 `json:"client_type_name"`
	IsLicense          bool                   `json:"is_license"`
	LicenseVariantID   int64                  `json:"license_variant_id"`
	LicenseVariantName string                 `json:"license_variant_name"`
	LicenseTypeID      int64                  `json:"license_type_id"`
	LicenseTypeName    string                 `json:"license_type_name"`
	DeploymentMethod   string                 `json:"deployment_method"`
	NoOfUser           int64                  `json:"no_of_user"`
	IsConcurrentUser   bool                   `json:"is_concurrent_user"`
	MaxOfflineDays     int64                  `json:"max_offline_days"`
	ModuleId1          int64                  `json:"module_id_1"`
	ModuleName1        string                 `json:"module_name_1"`
	ModuleId2          int64                  `json:"module_id_2"`
	ModuleName2        string                 `json:"module_name_2"`
	ModuleId3          int64                  `json:"module_id_3"`
	ModuleName3        string                 `json:"module_name_3"`
	ModuleId4          int64                  `json:"module_id_4"`
	ModuleName4        string                 `json:"module_name_4"`
	ModuleId5          int64                  `json:"module_id_5"`
	ModuleName5        string                 `json:"module_name_5"`
	ModuleId6          int64                  `json:"module_id_6"`
	ModuleName6        string                 `json:"module_name_6"`
	ModuleId7          int64                  `json:"module_id_7"`
	ModuleName7        string                 `json:"module_name_7"`
	ModuleId8          int64                  `json:"module_id_8"`
	ModuleName8        string                 `json:"module_name_8"`
	ModuleId9          int64                  `json:"module_id_9"`
	ModuleName9        string                 `json:"module_name_9"`
	ModuleId10         int64                  `json:"module_id_10"`
	ModuleName10       string                 `json:"module_name_10"`
	Component          []ViewProductComponent `json:"component"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	UpdatedName        string                 `json:"updated_name"`
}

type ViewProductComponent struct {
	ID             int64  `json:"id"`
	ComponentID    int64  `json:"component_id"`
	ComponentName  string `json:"component_name"`
	ComponentValue string `json:"component_value"`
}
