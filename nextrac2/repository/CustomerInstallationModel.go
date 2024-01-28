package repository

import "database/sql"

type CustomerInstallationModel struct {
	ID                       sql.NullInt64
	ParentCustomerID         sql.NullInt64
	ParentCustomerName       sql.NullString
	CustomerInstallationData []CustomerInstallationData
	CreatedBy                sql.NullInt64
	CreatedClient            sql.NullString
	CreatedAt                sql.NullTime
	UpdatedBy                sql.NullInt64
	UpdatedClient            sql.NullString
	UpdatedAt                sql.NullTime
	ClientMappingID          sql.NullInt64
}

type CustomerInstallationData struct {
	SiteID           sql.NullInt64
	CustomerID       sql.NullInt64
	CustomerSiteName sql.NullString
	Action           sql.NullInt32
	UpdatedAt        sql.NullTime
	Installation     []CustomerInstallationDetail
}

type CustomerInstallationDetail struct {
	InstallationID     sql.NullInt64
	ProductGroupID     sql.NullInt64
	ProductID          sql.NullInt64
	ProductCode        sql.NullString
	ProductName        sql.NullString
	Remark             sql.NullString
	UniqueID1          sql.NullString
	UniqueID2          sql.NullString
	BranchName         sql.NullString
	InstallationStatus sql.NullString
	InstallationDate   sql.NullTime
	ProductValidFrom   sql.NullTime
	ProductValidThru   sql.NullTime
	DayRange           sql.NullInt64
	Action             sql.NullInt32
	UpdatedAt          sql.NullTime
	ProductDescription sql.NullString
	IsUsed             sql.NullBool
	NoOfUser           sql.NullInt64
	ClientTypeID       sql.NullInt64
	ParentClientTypeID sql.NullInt64
	IsLicense          sql.NullBool
	ClientMappingID    sql.NullInt64
}

type CustomerInstallationForConfig struct {
	ID               sql.NullInt64
	ParentCustomerID sql.NullInt64
	CustomerID       sql.NullInt64
	SiteID           sql.NullInt64
	ProductID        sql.NullInt64
	BranchName       sql.NullString
	UniqueID1        sql.NullString
	UniqueID2        sql.NullString
	ClientID         sql.NullString
	ClientTypeID     sql.NullInt64
	LicenseVariantID sql.NullInt64
	LicenseTypeID    sql.NullInt64
	DeploymentMethod sql.NullString
	MaxOfflineDays   sql.NullInt64
	ModuleID1        sql.NullInt64
	ModuleID2        sql.NullInt64
	ModuleID3        sql.NullInt64
	ModuleID4        sql.NullInt64
	ModuleID5        sql.NullInt64
	ModuleID6        sql.NullInt64
	ModuleID7        sql.NullInt64
	ModuleID8        sql.NullInt64
	ModuleID9        sql.NullInt64
	ModuleID10       sql.NullInt64
	IsUserConcurrent sql.NullBool
}

type CustomerInstallationDetailConfig struct {
	InstallationID     sql.NullInt64
	ParentCustomerID   sql.NullInt64
	ParentCustomer     sql.NullString
	CustomerID         sql.NullInt64
	Customer           sql.NullString
	SiteID             sql.NullInt64
	ProductID          sql.NullInt64
	ProductName        sql.NullString
	ClientID           sql.NullString
	LicenseVariantID   sql.NullInt64
	LicenseVariantName sql.NullString
	LicenseTypeID      sql.NullInt64
	LicenseTypeName    sql.NullString
	DeploymentMethod   sql.NullString
	NoOfUser           sql.NullInt64
	IsUserConcurrent   sql.NullBool
	UniqueID1          sql.NullString
	UniqueID2          sql.NullString
	ProductValidFrom   sql.NullTime
	ProductValidThru   sql.NullTime
	MaxOfflineDays     sql.NullInt64
	ClientTypeID       sql.NullInt64
	ClientType         sql.NullString
	ModuleID1          sql.NullInt64
	ModuleIDName1      sql.NullString
	ModuleID2          sql.NullInt64
	ModuleIDName2      sql.NullString
	ModuleID3          sql.NullInt64
	ModuleIDName3      sql.NullString
	ModuleID4          sql.NullInt64
	ModuleIDName4      sql.NullString
	ModuleID5          sql.NullInt64
	ModuleIDName5      sql.NullString
	ModuleID6          sql.NullInt64
	ModuleIDName6      sql.NullString
	ModuleID7          sql.NullInt64
	ModuleIDName7      sql.NullString
	ModuleID8          sql.NullInt64
	ModuleIDName8      sql.NullString
	ModuleID9          sql.NullInt64
	ModuleIDName9      sql.NullString
	ModuleID10         sql.NullInt64
	ModuleIDName10     sql.NullString
	Component          []ProductComponentModel
	ClientMappingID    sql.NullInt64
}

type CustomerInstallationTracking struct {
	KeyID     sql.NullInt64
	UniqueID1 sql.NullString
	UniqueID2 sql.NullString
}
