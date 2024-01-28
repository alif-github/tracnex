package repository

import "database/sql"

type PKCEClientMappingModel struct {
	ID                sql.NullInt64
	ParentClientID    sql.NullString
	ClientID          sql.NullString
	ClientTypeID      sql.NullInt64
	AuthUserID        sql.NullInt64
	Username          sql.NullString
	InstallationID    sql.NullInt64
	CustomerID        sql.NullInt64
	SiteID            sql.NullInt64
	CompanyID         sql.NullString
	CompanyName       sql.NullString
	BranchID          sql.NullString
	BranchName        sql.NullString
	ClientAlias       sql.NullString
	IsClientDependant sql.NullString
	CreatedBy         sql.NullInt64
	CreatedClient     sql.NullString
	CreatedAt         sql.NullTime
	UpdatedBy         sql.NullInt64
	UpdatedClient     sql.NullString
	UpdatedAt         sql.NullTime
	UpdatedAtStr      sql.NullString
	ClientMappingID   sql.NullInt64
}

type ViewPKCEClientMappingModel struct {
	ID             sql.NullInt64
	ClientID       sql.NullString
	ParentClientID sql.NullString
	FirstName      sql.NullString
	LastName       sql.NullString
	Username       sql.NullString
	ClientType     sql.NullString
	CompanyID      sql.NullString
	BranchID       sql.NullString
	ClientAlias    sql.NullString
	CreatedAt      sql.NullTime
	CreatedBy      sql.NullInt64
	UpdatedAt      sql.NullTime
	UpdatedBy      sql.NullInt64
}

type CheckPKCEClientMappingModel struct {
	ID                  sql.NullInt64
	PKCEClientMappingID sql.NullInt64
	ClientID            sql.NullString
	AuthUserID          sql.NullInt64
	Username            sql.NullString
	Email               sql.NullString
	Phone               sql.NullString
	Firstname           sql.NullString
	Lastname            sql.NullString
	IsRegisteredBefore  sql.NullBool
}
