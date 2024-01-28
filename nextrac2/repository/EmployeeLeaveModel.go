package repository

import "database/sql"

type EmployeeLeaveModel struct {
	ID                 sql.NullInt64
	IDCard             sql.NullString
	Firstname          sql.NullString
	Lastname           sql.NullString
	Name               sql.NullString
	Type               sql.NullString
	CountType          sql.NullInt64
	Department         sql.NullString
	AllowanceId        sql.NullInt64
	AllowanceName      sql.NullString
	AllowanceType      sql.NullString
	Description        sql.NullString
	Date               sql.NullString
	Level              sql.NullString
	Grade              sql.NullString
	StrDateList        sql.NullString
	StartDate          sql.NullTime
	EndDate            sql.NullTime
	Value              sql.NullInt64
	Status             sql.NullString
	CancellationReason sql.NullString
	FileUploadId       sql.NullInt64
	EmployeeId         sql.NullInt64
	CreatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	CreatedClient      sql.NullString
	UpdatedAt          sql.NullTime
	UpdatedBy          sql.NullInt64
	UpdatedClient      sql.NullString
	Deleted            sql.NullBool
	MemberList         []string
	SearchBy           sql.NullString
	Keyword            sql.NullString
	LeaveTime          sql.NullTime
	IsYearly           sql.NullBool
	OnLeave            sql.NullBool
	LeaveDate          sql.NullString
	CurrentAnnualLeave sql.NullInt64
	LastAnnualLeave    sql.NullInt64
	OwingLeave         sql.NullInt64
	Year               sql.NullString
	StrStartDate       sql.NullString
	StrEndDate         sql.NullString
	CurrentMedicalValue sql.NullFloat64
	LastMedicalValue   sql.NullFloat64
	ClientID		   sql.NullString
	Email			   sql.NullString
	Host               sql.NullString
	Path               sql.NullString
}

type EmployeeLeaveReportModel struct {
	RowNumber       sql.NullInt64
	IDCard          sql.NullString
	Name            sql.NullString
	Position        sql.NullString
	Department      sql.NullString
	DateJoin        sql.NullString
	DateProbation   sql.NullString
	DetailLeave     sql.NullString
	TotalLeave      sql.NullInt64
	CurrentLeave    sql.NullInt64
	MonthLeave      sql.NullInt64
	YearLeave       sql.NullInt64
	DetailLeaveList []DateDetailReportModel
}

type DateDetailReportModel struct {
	Date        []sql.NullTime
	Type        sql.NullString
	Description sql.NullString
}
