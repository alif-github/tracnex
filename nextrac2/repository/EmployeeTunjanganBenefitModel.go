package repository

import "database/sql"

type EmpAllowanceModel struct {
	ID                 sql.NullInt64
	AllowanceName      sql.NullString
	Type               sql.NullString
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

type EmpBenefitModel struct {
	ID                 sql.NullInt64
	BenefitName        sql.NullString
	BenefitType        sql.NullString
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
