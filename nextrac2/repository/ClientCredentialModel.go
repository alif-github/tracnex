package repository

import "database/sql"

type ClientCredentialModel struct {
	ID            sql.NullInt64
	ClientID      sql.NullString
	ClientSecret  sql.NullString
	SignatureKey  sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	ClientTypeID  sql.NullInt64
}
