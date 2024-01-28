package repository

import "database/sql"

type FileUpload struct {
	ID            sql.NullInt64
	UUIDKey       sql.NullString
	FileName      sql.NullString
	FileSize      sql.NullInt64
	Category      sql.NullString
	Konektor      sql.NullString
	ParentID      sql.NullInt64
	Host          sql.NullString
	Path          sql.NullString
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	Deleted       sql.NullBool
}
