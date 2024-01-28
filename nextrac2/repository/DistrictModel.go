package repository

import "database/sql"

type DistrictModel struct {
	ID            sql.NullInt64
	ProvinceID    sql.NullInt64
	MdbDistrictID sql.NullInt64
	Code          sql.NullString
	Name          sql.NullString
	Status        sql.NullString
	UpdatedAt     sql.NullTime
	LastSync      sql.NullTime
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	CreatedClient sql.NullString
	UpdatedClient sql.NullString
}

type ListLocalDistrictModel struct {
	ID         sql.NullInt64
	ProvinceID sql.NullInt64
	Code       sql.NullString
	Name       sql.NullString
}
