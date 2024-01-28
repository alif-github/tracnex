package repository

import "database/sql"

type ClientTokenModel struct {
	ID            sql.NullInt64
	ClientID      sql.NullString
	AuthUserID    sql.NullInt64
	Token         sql.NullString
	ExpiredAt     sql.NullTime
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	Deleted       sql.NullBool
}
