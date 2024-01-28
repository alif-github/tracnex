package repository

import "database/sql"

type CustomerListModel struct {
	ID             sql.NullInt64
	CompanyID      sql.NullString
	BranchName     sql.NullString
	BranchID       sql.NullString
	CompanyName    sql.NullString
	City           sql.NullString
	Implementer    sql.NullString
	Implementation sql.NullTime
	Product        sql.NullString
	Version        sql.NullString
	LicenseType    sql.NullString
	UserAmount     sql.NullInt64
	ExpDate        sql.NullTime
	CreatedBy      sql.NullInt64
	CreatedAt      sql.NullTime
	CreatedClient  sql.NullString
	UpdatedBy      sql.NullInt64
	UpdatedAt      sql.NullTime
	UpdatedClient  sql.NullString
	Deleted        sql.NullBool
}
