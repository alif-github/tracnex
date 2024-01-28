package repository

import "database/sql"

type ParameterModel struct {
	ID            sql.NullInt64
	Permission    sql.NullString
	Name          sql.NullString
	Value         sql.NullString
	Description   sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}

type UserParameterModel struct {
	ID             sql.NullInt64
	UserID         sql.NullInt64
	ParameterValue sql.NullString
	CreatedBy      sql.NullInt64
	CreatedAt      sql.NullTime
	CreatedClient  sql.NullString
	UpdatedBy      sql.NullInt64
	UpdatedAt      sql.NullTime
	UpdatedClient  sql.NullString
	Deleted        sql.NullBool
}
