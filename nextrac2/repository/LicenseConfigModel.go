package repository

import "database/sql"

type LicenseConfigModel struct {
	ID                     sql.NullInt64
	UUIDKey                sql.NullString
	LicenseConfigIDs       []LicenseConfigIDsModel
	InstallationID         sql.NullInt64
	ParentCustomerID       sql.NullInt64
	ParentCustomer         sql.NullString
	CustomerID             sql.NullInt64
	Customer               sql.NullString
	SiteID                 sql.NullInt64
	ClientID               sql.NullString
	ProductID              sql.NullInt64
	ProductCode            sql.NullString
	ProductName            sql.NullString
	ClientTypeID           sql.NullInt64
	ClientType             sql.NullString
	LicenseVariantID       sql.NullInt64
	LicenseVariantName     sql.NullString
	LicenseTypeID          sql.NullInt64
	LicenseTypeName        sql.NullString
	DeploymentMethod       sql.NullString
	NoOfUser               sql.NullInt64
	IsUserConcurrent       sql.NullString
	MaxOfflineDays         sql.NullInt64
	UniqueID1              sql.NullString
	UniqueID2              sql.NullString
	ProductValidFrom       sql.NullTime
	ProductValidThru       sql.NullTime
	ModuleID1              sql.NullInt64
	ModuleID2              sql.NullInt64
	ModuleID3              sql.NullInt64
	ModuleID4              sql.NullInt64
	ModuleID5              sql.NullInt64
	ModuleID6              sql.NullInt64
	ModuleID7              sql.NullInt64
	ModuleID8              sql.NullInt64
	ModuleID9              sql.NullInt64
	ModuleID10             sql.NullInt64
	ModuleIDName1          sql.NullString
	ModuleIDName2          sql.NullString
	ModuleIDName3          sql.NullString
	ModuleIDName4          sql.NullString
	ModuleIDName5          sql.NullString
	ModuleIDName6          sql.NullString
	ModuleIDName7          sql.NullString
	ModuleIDName8          sql.NullString
	ModuleIDName9          sql.NullString
	ModuleIDName10         sql.NullString
	AllowActivation        sql.NullString
	PaymentStatus          sql.NullString
	IsExtendChecklist      sql.NullBool
	IsHasPrevLicenseConfig sql.NullBool
	PrevLicenseConfigID    sql.NullInt64
	CreatedBy              sql.NullInt64
	CreatedClient          sql.NullString
	CreatedAt              sql.NullTime
	UpdatedBy              sql.NullInt64
	UpdatedName            sql.NullString
	UpdatedClient          sql.NullString
	UpdatedAt              sql.NullTime
	IsUsed                 sql.NullBool
	CheckSum               sql.NullString
	ComponentSting         sql.NullString
	Component              []ProductComponentModel
	ProductLicenseStatus   sql.NullInt64
	ProductKey             sql.NullString
	SalesmanString         sql.NullString
}

type LicenseConfigComponent struct {
	LicenseConfigID sql.NullInt64
	ProductID       sql.NullInt64
	ComponentID     sql.NullInt64
	ComponentName   sql.NullString
	ComponentValue  sql.NullString
	CreatedBy       sql.NullInt64
	CreatedClient   sql.NullString
	CreatedAt       sql.NullTime
	UpdatedBy       sql.NullInt64
	UpdatedClient   sql.NullString
	UpdatedAt       sql.NullTime
}

type LicenseConfigIDsModel struct {
	ID sql.NullInt64
}
