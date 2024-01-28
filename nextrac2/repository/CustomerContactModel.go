package repository

import "database/sql"

type CustomerContactModel struct {
	ID                 sql.NullInt64
	CustomerID         sql.NullInt64
	MdbPersonProfileID sql.NullInt64
	Nik                sql.NullString
	MdbPersonTitle     sql.NullInt64
	PersonTitle        sql.NullString
	FirstName          sql.NullString
	LastName           sql.NullString
	Sex                sql.NullString
	Address            sql.NullString
	Address2           sql.NullString
	Address3           sql.NullString
	Hamlet             sql.NullString
	Neighbourhood      sql.NullString
	ProvinceID         sql.NullInt64
	ProvinceName       sql.NullString
	DistrictID         sql.NullInt64
	DistrictName       sql.NullString
	Phone              sql.NullString
	Email              sql.NullString
	MdbPositionID      sql.NullInt64
	PositionName       sql.NullString
	Status             sql.NullString
	CreatedBy          sql.NullInt64
	CreatedAt          sql.NullTime
	CreatedClient      sql.NullString
	CreatedName        sql.NullString
	UpdatedBy          sql.NullInt64
	UpdatedAt          sql.NullTime
	UpdatedClient      sql.NullString
	UpdatedName        sql.NullString
	Deleted            sql.NullBool
}
