package repository

import "database/sql"

type AuditSystemModel struct {
	ID            sql.NullInt64
	UUIDKey       sql.NullString
	TableName     sql.NullString
	PrimaryKey    sql.NullInt64
	Data          sql.NullString
	Description   sql.NullString
	Action        sql.NullInt32
	Employee	  AuditEmployee
	CreatedName   sql.NullString
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
}

type AuditEmployee struct {
	Firstname string
	Lastname  string
}

type GetAuditDataSystemModel struct {
	ID   int64
	Data string
}

type AuditSystemFieldParam struct {
	IsEqual    bool
	ParamValue interface{}
}
