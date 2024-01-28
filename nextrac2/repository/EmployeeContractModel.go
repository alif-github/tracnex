package repository

import "database/sql"

type EmployeeContractModel struct {
	ID            sql.NullInt64
	ContractNo    sql.NullString
	Information   sql.NullString
	EmployeeID    sql.NullInt64
	FromDate      sql.NullTime
	ThruDate      sql.NullTime
	Deleted       sql.NullBool
	CreatedBy     sql.NullInt64
	CreatedName   sql.NullString
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedName   sql.NullString
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
}
