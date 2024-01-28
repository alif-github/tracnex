package repository

import "database/sql"

type EmployeeLevelModel struct {
	ID                 sql.NullInt64
	Level              sql.NullString
	Description        sql.NullString
	CreatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	CreatedName        sql.NullString
	CreatedClient      sql.NullString
	UpdatedBy          sql.NullInt64
	UpdatedAt          sql.NullTime
	UpdatedName        sql.NullString
	UpdatedClient      sql.NullString
	Deleted            sql.NullBool
}

type EmployeeFacilitiesActiveModel struct {
	ID                 sql.NullInt64
	Value              sql.NullString
	Active             sql.NullBool
	BenefitID          sql.NullInt64
	Benefit            sql.NullString
	AllowanceID        sql.NullInt64
	Allowance          sql.NullString
	AllowanceType      sql.NullString
	LevelID            sql.NullInt64
	Level              sql.NullString
	GradeID            sql.NullInt64
	Grade              sql.NullString
	CreatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	CreatedName        sql.NullString
	CreatedClient      sql.NullString
	UpdatedBy          sql.NullInt64
	UpdatedAt          sql.NullTime
	UpdatedName        sql.NullString
	UpdatedClient      sql.NullString
	Deleted            sql.NullBool
}