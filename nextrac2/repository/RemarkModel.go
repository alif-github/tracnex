package repository

import "database/sql"

type RemarkModel struct {
	ID            sql.NullInt64
	Name          sql.NullString
	Value         sql.NullString
	Level         sql.NullInt64
	ParentID      sql.NullInt64
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
}
