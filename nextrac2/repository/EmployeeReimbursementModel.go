package repository

import "database/sql"

type EmployeeReimbursement struct {
	ID                  sql.NullInt64
	IDCard              sql.NullString
	Firstname           sql.NullString
	Lastname            sql.NullString
	FullName			sql.NullString
	Department          sql.NullString
	Name                sql.NullString
	ReceiptNo           sql.NullString
	BenefitId           sql.NullInt64
	Description         sql.NullString
	Date                sql.NullTime
	Value               sql.NullFloat64
	CurrentMedicalValue sql.NullFloat64
	LastMedicalValue    sql.NullFloat64
	Status              sql.NullString
	VerifiedStatus      sql.NullString
	ApprovedValue       sql.NullFloat64
	FileUploadId        sql.NullInt64
	Host				sql.NullString
	Path 				sql.NullString
	Note				sql.NullString
	EmployeeId          sql.NullInt64
	StartDate			sql.NullString
	EndDate				sql.NullString
	CancellationReason	sql.NullString
	CreatedBy           sql.NullInt64
	CreatedAt           sql.NullTime
	CreatedClient       sql.NullString
	UpdatedAt           sql.NullTime
	UpdatedBy           sql.NullInt64
	UpdatedClient       sql.NullString
	Deleted             sql.NullBool
	SearchBy			sql.NullString
	Keyword				sql.NullString
	DateJoin            sql.NullTime
	DateOut             sql.NullTime
	MonthlyReport       EmployeeReimbursementMonthlyReport
	Year                sql.NullString
	Month               sql.NullString
	IsFilter            sql.NullBool
	MonthlyReportArr    [12]float64
	ReportType          sql.NullString
	Email				sql.NullString
	ClientId			sql.NullString
}

type EmployeeReimbursementMonthlyReport struct {
     Total              sql.NullFloat64
     January            sql.NullFloat64
     February           sql.NullFloat64
     March              sql.NullFloat64
     April              sql.NullFloat64
     May                sql.NullFloat64
     June               sql.NullFloat64
     July               sql.NullFloat64
     August             sql.NullFloat64
     September          sql.NullFloat64
     October            sql.NullFloat64
     November           sql.NullFloat64
     December           sql.NullFloat64
}