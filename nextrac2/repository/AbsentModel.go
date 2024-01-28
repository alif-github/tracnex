package repository

import "database/sql"

type AbsentModel struct {
	ID              sql.NullInt64
	EmployeeID      sql.NullInt64
	EmployeeName    sql.NullString
	IDCard          sql.NullString
	AbsentID        sql.NullInt64
	NormalDays      sql.NullInt64
	ActualDays      sql.NullInt64
	Absent          sql.NullInt64
	Overdue         sql.NullInt64
	LeaveEarly      sql.NullInt64
	Overtime        sql.NullInt64
	NumberOfLeave   sql.NullInt64
	LeavingDuties   sql.NullInt64
	NumbersIn       sql.NullInt64
	NumbersOut      sql.NullInt64
	Scan            sql.NullInt64
	SickLeave       sql.NullInt64
	PaidLeave       sql.NullInt64
	PermissionLeave sql.NullInt64
	WorkHours       sql.NullInt64
	PercentAbsent   sql.NullFloat64
	PeriodStart     sql.NullTime
	PeriodEnd       sql.NullTime
	CreatedBy       sql.NullInt64
	CreatedClient   sql.NullString
	CreatedAt       sql.NullTime
	UpdatedBy       sql.NullInt64
	UpdatedClient   sql.NullString
	UpdatedAt       sql.NullTime
}
