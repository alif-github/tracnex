package repository

import "database/sql"

type MenuModel struct {
	ID              sql.NullInt64
	ParentMenuID    sql.NullInt64
	ServiceMenuID   sql.NullInt64
	Name            sql.NullString
	EnName          sql.NullString
	Sequence        sql.NullInt64
	IconName        sql.NullString
	Background      sql.NullString
	AvailableAction sql.NullString
	MenuCode        sql.NullString
	Status          sql.NullString
	Url             sql.NullString
	TableName       sql.NullString
	CreatedBy       sql.NullInt64
	UpdatedBy       sql.NullInt64
	UpdatedClient   sql.NullString
	UpdatedAt       sql.NullTime
}
