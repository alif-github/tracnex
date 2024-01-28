package repository

import "database/sql"

type ProductLicenseModel struct {
	ID                     sql.NullInt64
	LicenseConfigId        sql.NullInt64
	ProductKey             sql.NullString
	ProductEncrypt         sql.NullString
	ProductSignature       sql.NullString
	ClientId               sql.NullString
	ClientSecret           sql.NullString
	HWID                   sql.NullString
	ActivationDate         sql.NullTime
	LicenseStatus          sql.NullInt32
	TerminationDescription sql.NullString
	DateModified           sql.NullTime
	ModifiedBy             sql.NullString
	CreatedBy              sql.NullInt64
	CreatedAt              sql.NullTime
	CreatedClient          sql.NullString
	UpdatedBy              sql.NullInt64
	UpdatedAt              sql.NullTime
	UpdatedClient          sql.NullString
	Deleted                sql.NullBool
	ProductValidFrom       sql.NullTime
	ProductValidThru       sql.NullTime
	SignatureKey           sql.NullString
	ProductId              sql.NullString
}

type ProductLicenseModelForView struct {
	ID                 sql.NullInt64
	LicenseConfigId    sql.NullInt64
	CustomerName       sql.NullString
	UniqueId1          sql.NullString
	UniqueId2          sql.NullString
	InstallationId     sql.NullInt64
	ProductName        sql.NullString
	LicenseVariantName sql.NullString
	LicenseTypeName    sql.NullString
	ProductValidFrom   sql.NullTime
	ProductValidThru   sql.NullTime
	LicenseStatus      sql.NullInt32
}

type DetailProductLicense struct {
	ID                     sql.NullInt64
	ProductKey             sql.NullString
	ActivationDate         sql.NullTime
	LicenseStatus          sql.NullInt32
	TerminationDescription sql.NullString
	LicenseConfigId        sql.NullInt64
	InstallationId         sql.NullInt64
	ParentCustomerId       sql.NullInt64
	ParentCustomer         sql.NullString
	CustomerId             sql.NullInt64
	SiteId                 sql.NullInt64
	Customer               sql.NullString
	ClientId               sql.NullString
	Product                sql.NullString
	Client                 sql.NullString
	LicenseVariant         sql.NullString
	LicenseType            sql.NullString
	DeploymentMethod       sql.NullString
	NumberOfUser           sql.NullInt64
	ConcurentUser          sql.NullString
	UniqueId1              sql.NullString
	UniqueId2              sql.NullString
	LicenseValidFrom       sql.NullTime
	LicenseValidThru       sql.NullTime
	CreatedAt              sql.NullTime
	UpdatedAt              sql.NullTime
	Module1                sql.NullString
	Module2                sql.NullString
	Module3                sql.NullString
	Module4                sql.NullString
	Module5                sql.NullString
	Module6                sql.NullString
	Module7                sql.NullString
	Module8                sql.NullString
	Module9                sql.NullString
	Module10               sql.NullString
	Components             []ProductComponentModel
	AliasName              sql.NullString
}

type LicenseSalesJournal struct {
	ID               sql.NullInt64
	ClientID         sql.NullString
	LicenseStatusID  sql.NullInt64
	LicenseStatus    sql.NullString
	UniqueID1        sql.NullString
	UniqueID2        sql.NullString
	ProductName      sql.NullString
	ClientType       sql.NullString
	AllowActivation  sql.NullString
	NoOfUser         sql.NullInt64
	ProductValidFrom sql.NullTime
	ProductValidThru sql.NullTime
	IsUserConcurrent sql.NullString
	TotalLicense     sql.NullInt64
	TotalActivated   sql.NullInt64
}

type ProductLicenseResponseForScheduler struct {
	ContentDataOutDetail []ContentDataOutDetail `json:"content_data_out_detail"`
	TotalData            int64                  `json:"total_data"`
}
