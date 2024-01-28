package repository

import "database/sql"

type ProvinceModel struct {
	ID            sql.NullInt64
	CountryID	  sql.NullInt64
	MDBProvinceID sql.NullInt64
	Code          sql.NullString
	Name          sql.NullString
	Status        sql.NullString
	LastSync      sql.NullTime
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}

type ListLocalProvinceModel struct {
	ID         sql.NullInt64
	CountryID  sql.NullInt64
	Code       sql.NullString
	Name       sql.NullString
	DistrictID sql.NullString
}
