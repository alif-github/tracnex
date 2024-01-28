package repository

import "database/sql"

type EmployeeHistoryModel struct {
	ID           sql.NullInt64
	Description1 sql.NullString
	Description2 sql.NullString
	CreatedAt    sql.NullTime
	Editor       sql.NullString
	IDBenefit    sql.NullInt64
}
