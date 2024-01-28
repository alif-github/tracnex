package repository

import "database/sql"

type NexmileParameterModel struct {
	ID            sql.NullInt64
	ClientID      sql.NullString
	ClientTypeID  sql.NullInt64
	UniqueID1     sql.NullString
	UniqueID2     sql.NullString
	SalesmanID    sql.NullInt64
	ParameterData []ParameterValueModel
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	AuthUserID    sql.NullInt64
	UserID        sql.NullString
	Password      sql.NullString
	AndroidID     sql.NullString
}

type NexmileParameterModelMap struct {
	ID            sql.NullInt64
	CLientID      sql.NullString
	ClientTypeID  sql.NullInt64
	UniqueID1     sql.NullString
	UniqueID2     sql.NullString
	SalesmanID    sql.NullInt64
	ParameterData map[string]ParameterValueModel
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	AuthUserID    sql.NullInt64
	UserID        sql.NullString
	Password      sql.NullString
	AndroidID     sql.NullString
}

type ParameterValueModel struct {
	ID             sql.NullInt64
	ParameterID    sql.NullString
	ParameterValue sql.NullString
	UpdatedBy      sql.NullInt64
	UpdatedClient  sql.NullString
	UpdatedAt      sql.NullTime
}
