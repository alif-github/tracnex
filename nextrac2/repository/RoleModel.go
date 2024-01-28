package repository

import "database/sql"

type RoleModel struct {
	ID            sql.NullInt64
	RoleID        sql.NullString
	Description   sql.NullString
	Permission    sql.NullString
	Status        sql.NullString
	Level         sql.NullInt32
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	CreatedName   sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
	IsUsed        sql.NullBool
}
