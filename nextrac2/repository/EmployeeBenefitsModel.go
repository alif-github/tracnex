package repository

import "database/sql"

type EmployeeBenefitsModel struct {
	ID                 sql.NullInt64
	EmployeeID         sql.NullInt64
	EmployeeLevelID    sql.NullInt64
	EmployeeGradeID    sql.NullInt64
	Salary             sql.NullFloat64
	BPJSNo             sql.NullString
	BPJSTkNo           sql.NullString
	CurrentAnnualLeave sql.NullInt64
	LastAnnualLeave    sql.NullInt64
	CurrentMedicalValue sql.NullFloat64
	LastMedicalValue    sql.NullFloat64
	VehicleLimit       sql.NullFloat64
	CreatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	CreatedName        sql.NullString
	CreatedClient      sql.NullString
	UpdatedBy          sql.NullInt64
	UpdatedAt          sql.NullTime
	UpdatedName        sql.NullString
	UpdatedClient      sql.NullString
	Deleted            sql.NullBool
	Year               sql.NullString
	JoinDate           sql.NullTime
	CutOffLeaveValue   sql.NullInt64
	NoteCutOff         sql.NullString
}
