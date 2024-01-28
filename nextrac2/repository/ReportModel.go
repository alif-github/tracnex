package repository

import (
	"database/sql"
	"time"
)

type ReportModel struct {
	NIK                 sql.NullInt64
	Name                sql.NullString
	DepartmentID        sql.NullInt64
	Department          sql.NullString
	BacklogManday       sql.NullFloat64
	ActualManday        sql.NullFloat64
	ActualHistoryManday sql.NullString
	Rate                sql.NullFloat64
	RateStr             sql.NullString
	Tracker             sql.NullString
	RedmineNumber       sql.NullString
}

type ReportHistoryModel struct {
	ID             sql.NullInt64
	Data           sql.NullString
	SuccessTicket  sql.NullString
	DepartmentID   sql.NullInt64
	DepartmentName sql.NullString
	CreatedBy      sql.NullInt64
	CreatedClient  sql.NullString
	CreatedAt      sql.NullTime
	CreatedName    sql.NullString
}

type RedmineModel struct {
	RedmineTicket sql.NullInt64
	Project       sql.NullString
	Signed        sql.NullString
	SignedID      sql.NullInt64
	Subject       sql.NullString
	Sprint        sql.NullString
	Status        sql.NullString
	StatusID      sql.NullInt64
	CreatedAt     sql.NullTime
	Manhour       sql.NullFloat64
	Payment       sql.NullString
	Tracker       sql.NullString
	UpdatedOn     sql.NullTime
}

type RedmineInfraModel struct {
	RedmineTicket sql.NullInt64
	Project       sql.NullString
	Signed        sql.NullString
	SignedID      sql.NullInt64
	Subject       sql.NullString
	Status        sql.NullString
	Value         sql.NullString
	CreatedAt     sql.NullTime
	Notes         sql.NullString
	Manhour       sql.NullFloat64
	Tracker       sql.NullString
	UpdatedOn     sql.NullTime
	Category      sql.NullString
}

type HistoryTimeReportRedmineModel struct {
	User          sql.NullInt64
	RedmineTicket sql.NullInt64
	Tracker       sql.NullString
	TimeHistory   sql.NullString
}

type TicketRedmineModel struct {
	Ticket  int64                 `json:"ticket"`
	History []HistoryRedmineModel `json:"history"`
}

type HistoryRedmineModel struct {
	Subject      string `json:"subject"`
	CreatedOnStr string `json:"created_on"`
	CreatedOn    time.Time
}
