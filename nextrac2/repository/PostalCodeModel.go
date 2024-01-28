package repository

import "database/sql"

type PostalCodeModel struct {
	ID               sql.NullInt64
	UrbanVillageID   sql.NullInt64
	UrbanVillageName sql.NullString
	MDBPostalCodeID  sql.NullInt64
	Code             sql.NullString
	Status           sql.NullString
	CreatedAt        sql.NullTime
	CreatedBy        sql.NullInt64
	CreatedClient    sql.NullString
	UpdatedAt        sql.NullTime
	UpdatedBy        sql.NullInt64
	UpdatedClient    sql.NullString
	Deleted          sql.NullBool
	LastSync         sql.NullTime
}
