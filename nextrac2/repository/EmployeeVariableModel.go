package repository

import "database/sql"

type EmployeeVariableModel struct {
	ID            sql.NullInt64
	EmployeeID    sql.NullInt64
	RedmineID     sql.NullInt64
	MandaysRate   sql.NullString
	LeadMandays   sql.NullFloat64
	CreatedAt     sql.NullTime
	CreatedBy     sql.NullInt64
	CreatedName   sql.NullString
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedName   sql.NullString
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
	IsUsed        sql.NullBool
}
