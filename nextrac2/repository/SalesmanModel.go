package repository

import "database/sql"

type SalesmanModel struct {
	ID            sql.NullInt64
	PersonTitleID sql.NullInt64
	PersonTitle   sql.NullString
	Sex           sql.NullString
	Nik           sql.NullString
	FirstName     sql.NullString
	LastName      sql.NullString
	Address       sql.NullString
	Hamlet        sql.NullString
	Neighbourhood sql.NullString
	ProvinceID    sql.NullInt64
	DistrictID    sql.NullInt64
	Phone         sql.NullString
	Email         sql.NullString
	Status        sql.NullString
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	IsUsed        sql.NullBool
}

type ViewSalesmanModel struct {
	ID            sql.NullInt64
	PersonTitleID sql.NullInt64
	PersonTitle   sql.NullString
	Sex           sql.NullString
	Nik           sql.NullString
	FirstName     sql.NullString
	LastName      sql.NullString
	Address       sql.NullString
	Hamlet        sql.NullString
	Neighbourhood sql.NullString
	MdbProvinceID sql.NullInt64
	Province      sql.NullString
	MdbDistrictID sql.NullInt64
	District      sql.NullString
	Phone         sql.NullString
	Email         sql.NullString
	Status        sql.NullString
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	UpdatedName   sql.NullString
}

type ListSalesmanModel struct {
	ID        sql.NullInt64
	FirstName sql.NullString
	LastName  sql.NullString
	Address   sql.NullString
	Province  sql.NullString
	District  sql.NullString
	Phone     sql.NullString
	Email     sql.NullString
	Status    sql.NullString
	UpdatedAt sql.NullTime
}
