package repository

import "database/sql"

type Allowance struct {
	ID            sql.NullInt64
	AllowanceName sql.NullString
	AllowanceType sql.NullString
	Description   sql.NullString
	Active        sql.NullBool
	Value		  sql.NullString
	EmployeeLevelId sql.NullInt64
	EmployeeGradeId sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	UpdatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}
