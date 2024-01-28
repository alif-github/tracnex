package repository

import "database/sql"

type ClientRegistNonOnPremiseModel struct {
	ClientID      sql.NullString
	SignatureKey  sql.NullString
	ClientSecret  sql.NullString
	AliasName     sql.NullString
	FirstName     sql.NullString
	ClientTypeID  sql.NullInt64
	DetailUnique  []DetailUniqueID
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
}

type DetailUniqueID struct {
	InstallationID    sql.NullInt64
	ParentCustomerID  sql.NullInt64
	CustomerID        sql.NullInt64
	SiteID            sql.NullInt64
	UniqueID1         sql.NullString
	UniqueID2         sql.NullString
	IsError           sql.NullBool
	ErrorMessage      sql.NullString
	CreatedBy         sql.NullInt64
	CreatedClient     sql.NullString
	CreatedAt         sql.NullTime
	UpdatedBy         sql.NullInt64
	UpdatedClient     sql.NullString
	UpdatedAt         sql.NullTime
	InstallationIDCol []InstallationIDColInt
}

type InstallationIDColInt struct {
	InstallationID sql.NullInt64
}
