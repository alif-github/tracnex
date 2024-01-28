package repository

import "database/sql"

type StandarManhourModel struct {
	ID            sql.NullInt64
	Case          sql.NullString
	DepartmentID  sql.NullInt64
	Department    sql.NullString
	Manhour       sql.NullFloat64
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	UpdatedName   sql.NullString
	Deleted       sql.NullBool
}
