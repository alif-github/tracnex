package repository

import "database/sql"

type CompanyModel struct {
	ID                      sql.NullInt64
	UUIDKey                 sql.NullString
	Npwp                    sql.NullString
	CompanyTitle            sql.NullString
	CompanyName             sql.NullString
	Address                 sql.NullString
	Address2                sql.NullString
	Hamlet                  sql.NullString
	Neighbourhood           sql.NullString
	ProvinceID              sql.NullInt64
	ProvinceName            sql.NullString
	DistrictID              sql.NullInt64
	DistrictName            sql.NullString
	SubDistrictID           sql.NullInt64
	SubDistrictName         sql.NullString
	UrbanVillageID          sql.NullInt64
	UrbanVillageName        sql.NullString
	PostalCodeID            sql.NullInt64
	PostalCode              sql.NullString
	Longitude               sql.NullString
	Latitude                sql.NullString
	Phone                   sql.NullString
	AlternativePhone        sql.NullString
	Fax                     sql.NullString
	CompanyEmail            sql.NullString
	AlternativeCompanyEmail sql.NullString
	CustomerSource          sql.NullString
	TaxName                 sql.NullString
	TaxAddress              sql.NullString
	Telephone               sql.NullString
	AlternateTelephone      sql.NullString
	CreatedBy               sql.NullInt64
	CreatedAt               sql.NullTime
	CreatedClient           sql.NullString
	CreatedName             sql.NullString
	UpdatedBy               sql.NullInt64
	UpdatedAt               sql.NullTime
	UpdatedClient           sql.NullString
	UpdatedName             sql.NullString
	Deleted                 sql.NullBool
}
