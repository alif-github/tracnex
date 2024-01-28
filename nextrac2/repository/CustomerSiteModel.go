package repository

import "database/sql"

type CustomerSiteModel struct {
	ID               sql.NullInt64
	ParentCustomerID sql.NullInt64
	CustomerID       sql.NullInt64
	CreatedBy        sql.NullInt64
	CreatedClient    sql.NullString
	CreatedAt        sql.NullTime
	UpdatedBy        sql.NullInt64
	UpdatedClient    sql.NullString
	UpdatedAt        sql.NullTime
	Deleted          sql.NullBool
}