package repository

import "database/sql"

type CustomerCategoryModel struct {
	ID                   sql.NullInt64
	CustomerCategoryID   sql.NullString
	CustomerCategoryName sql.NullString
	CreatedBy            sql.NullInt64
	CreatedAt            sql.NullTime
	CreatedClient        sql.NullString
	UpdatedBy            sql.NullInt64
	UpdatedAt            sql.NullTime
	UpdatedClient        sql.NullString
	UpdatedName          sql.NullString
	Deleted              sql.NullBool
	IsUsed               sql.NullBool
}
