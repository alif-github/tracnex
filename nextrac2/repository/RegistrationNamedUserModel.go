package repository

import "database/sql"

type RegistrationNamedUserModel struct {
	UserLicenseID    sql.NullInt64
	ParentCustomerID sql.NullInt64
	CustomerID       sql.NullInt64
	SiteID           sql.NullInt64
	InstallationID   sql.NullInt64
	ClientID         sql.NullString
	UniqueID1        sql.NullString
	UniqueID2        sql.NullString
	AuthUserID       sql.NullInt64
	UserID           sql.NullString
	Password         sql.NullString
	SalesmanID       sql.NullString
	AndroidID        sql.NullString
	RegDate          sql.NullTime
	Status           sql.NullString
	Email            sql.NullString
	NoTelp           sql.NullString
	SalesmanCategory sql.NullString
	ProductValidFrom sql.NullTime
	ProductValidThru sql.NullTime
	CreatedBy        sql.NullInt64
	CreatedClient    sql.NullString
	CreatedAt        sql.NullTime
	UpdatedBy        sql.NullInt64
	UpdatedClient    sql.NullString
	UpdatedAt        sql.NullTime
}
