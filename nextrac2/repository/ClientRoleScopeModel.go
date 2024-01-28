package repository

import "database/sql"

type ClientRoleScopeModel struct {
	ID            sql.NullInt64
	UUIDKey       sql.NullString
	ClientID      sql.NullString
	RoleID        sql.NullInt64
	GroupID       sql.NullInt64
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	Deleted       sql.NullBool
}
