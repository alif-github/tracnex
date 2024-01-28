package repository

import "database/sql"

type CRONSchedulerModel struct {
	ID            sql.NullInt64
	Name          sql.NullString
	RunType       sql.NullString
	CRON          sql.NullString
	Status        sql.NullBool
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}

type HostServerModel struct {
	ID            sql.NullInt64
	HostName      sql.NullString
	HostURL       sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}

type CRONHostModel struct {
	ID            sql.NullInt64
	HostID        sql.NullInt64
	CronID        sql.NullInt64
	HostName      sql.NullString
	HostURL       sql.NullString
	Status        sql.NullBool
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}

type RefreshScheduler struct {
	HostName sql.NullString
	HostURL  sql.NullString
	CRON     sql.NullString
}

type ServerRunModel struct {
	ID            sql.NullInt64
	Name          sql.NullString
	HostID        sql.NullInt64
	RunType       sql.NullString
	HostName      sql.NullString
	Status        sql.NullBool
	Parameter     sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}

type ListScheduler struct {
	ID     sql.NullInt64
	Name   sql.NullString
	Cron   sql.NullString
	Status sql.NullString
}
