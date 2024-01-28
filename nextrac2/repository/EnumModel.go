package repository

import "database/sql"

type EnumModel struct {
	Oid           sql.NullInt64
	EnumTypeId    sql.NullInt64
	EnumSortOrder sql.NullInt64
	EnumLabel     sql.NullString
}