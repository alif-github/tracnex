package repository

import "database/sql"

type Benefit struct {
	ID            sql.NullInt64
	BenefitName   sql.NullString
	BenefitType   sql.NullString
	Description   sql.NullString
	Active        sql.NullBool
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
