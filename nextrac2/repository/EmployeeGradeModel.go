package repository

import "database/sql"

type EmployeeGradeModel struct {
	ID                 sql.NullInt64
	Grade              sql.NullString
	Description        sql.NullString
	CreatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	CreatedName        sql.NullString
	CreatedClient      sql.NullString
	UpdatedBy          sql.NullInt64
	UpdatedAt          sql.NullTime
	UpdatedName        sql.NullString
	UpdatedClient      sql.NullString
	Deleted            sql.NullBool
}
