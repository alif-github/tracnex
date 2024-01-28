package repository

import "database/sql"

type ClientTypeModel struct {
	ID                 sql.NullInt64
	ClientType         sql.NullString
	Description        sql.NullString
	ParentClientTypeID sql.NullInt64
	CreatedBy          sql.NullInt64
	CreatedClient      sql.NullString
	CreatedAt          sql.NullTime
	UpdatedBy          sql.NullInt64
	UpdatedName        sql.NullString
	UpdatedClient      sql.NullString
	UpdatedAt          sql.NullTime
}
