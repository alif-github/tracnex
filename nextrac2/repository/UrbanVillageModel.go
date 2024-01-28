package repository

import "database/sql"

type UrbanVillageModel struct {
	ID                sql.NullInt64
	SubDistrictID     sql.NullInt64
	SubDistrictName   sql.NullString
	MDBUrbanVillageID sql.NullInt64
	Code              sql.NullString
	Name              sql.NullString
	Status            sql.NullString
	CreatedAt         sql.NullTime
	CreatedBy         sql.NullInt64
	CreatedClient     sql.NullString
	UpdatedAt         sql.NullTime
	UpdatedBy         sql.NullInt64
	UpdatedClient     sql.NullString
	Deleted           sql.NullBool
	LastSync          sql.NullTime
}
