package repository

import "database/sql"

type EmployeeRequestHistory struct {
	ID                 sql.NullInt64
	EmployeeId		   sql.NullInt64
	ReceiptNo		   sql.NullString
	Description        sql.NullString
	LeaveDate          sql.NullString
	Date               sql.NullTime
	RequestType        sql.NullString
	Type			   sql.NullString
	Status             sql.NullString
	VerifiedStatus     sql.NullString
	CancellationReason sql.NullString
	TotalLeave         sql.NullInt64
	Value              sql.NullFloat64
	ApprovedValue      sql.NullFloat64
	Note               sql.NullString
	Host               sql.NullString
	Path               sql.NullString
	CreatedBy          sql.NullInt64
	CreatedAt          sql.NullTime
	UpdatedAt          sql.NullTime
	Deleted			   sql.NullBool
}