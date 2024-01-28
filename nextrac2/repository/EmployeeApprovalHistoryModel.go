package repository

import "database/sql"

type EmployeeApprovalHistory struct {
	ID                  sql.NullInt64
	EmployeeId			sql.NullInt64
	Firstname           sql.NullString
	Lastname            sql.NullString
	IDCard              sql.NullString
	Department			sql.NullString
	RequestType			sql.NullString
	Type				sql.NullString
	LeaveDate			sql.NullString
	Value				sql.NullFloat64
	Date             	sql.NullTime
	Status              sql.NullString
	VerifiedStatus		sql.NullString
	ApprovedValue		sql.NullFloat64
	Host                sql.NullString
	Path                sql.NullString
	TotalLeave			sql.NullInt64
	TotalRemainingLeave sql.NullInt64
	CancellationReason	sql.NullString
	Note				sql.NullString
	ReceiptNo			sql.NullString
	Description         sql.NullString
	CreatedAt           sql.NullTime
	UpdatedAt			sql.NullTime
	CreatedBy			sql.NullInt64
	Deleted             sql.NullBool
}
