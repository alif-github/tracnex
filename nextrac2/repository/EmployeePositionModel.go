package repository

import "database/sql"

type EmployeePositionModel struct {
	ID            sql.NullInt64
	Name          sql.NullString
	Description   sql.NullString
	CompanyID     sql.NullInt64
	CompanyName   sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	CreatedName   sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	UpdatedName   sql.NullString
}
