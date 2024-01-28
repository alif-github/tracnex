package repository

import (
	"database/sql"
)

type UserLicenseModel struct {
	ID                 sql.NullInt64
	ProductLicenseID   sql.NullInt64
	ParentCustomerId   sql.NullInt64
	ParentCustomerName sql.NullString
	CustomerId         sql.NullInt64
	CustomerName       sql.NullString
	SiteId             sql.NullInt64
	InstallationId     sql.NullInt64
	ClientID           sql.NullString
	UniqueId1          sql.NullString
	UniqueId2          sql.NullString
	ProductValidFrom   sql.NullTime
	ProductValidThru   sql.NullTime
	TotalLicense       sql.NullInt64
	TotalActivated     sql.NullInt64
	LicenseConfigId    sql.NullInt64
	ProductKey         sql.NullString
	ProductSignature   sql.NullString
	ProductEncrypt     sql.NullString
	ProductName        sql.NullString
	QuotaLicense       sql.NullInt64
	CreatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	CreatedClient      sql.NullString
	UpdatedAt          sql.NullTime
	UpdatedBy          sql.NullInt64
	UpdatedClient      sql.NullString
	LicenseStatus      sql.NullInt64
	ClientTypeId       sql.NullInt64
}
