package repository

import "database/sql"

type DataGroupModel struct {
	ID            sql.NullInt64
	GroupID       sql.NullString
	Description   sql.NullString
	Scope         sql.NullString
	Level         sql.NullInt32
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	CreatedName   sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	UpdatedName   sql.NullString
	Deleted       sql.NullBool
	IsUsed        sql.NullBool
}
