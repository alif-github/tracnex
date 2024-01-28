package repository

import (
	"database/sql"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"time"
)

type JobProcessModel struct {
	ID             sql.NullInt64
	UUIDKey        sql.NullString
	ParentJobID    sql.NullString
	Level          sql.NullInt32
	JobID          sql.NullString
	Group          sql.NullString
	Type           sql.NullString
	Name           sql.NullString
	Counter        sql.NullInt32
	Total          sql.NullInt32
	Status         sql.NullString
	MessageAlert   sql.NullString
	AlertId        sql.NullString
	AlertContent   sql.NullString
	Parameter      sql.NullString
	ContentDataOut sql.NullString
	UrlIn          sql.NullString
	FilenameIn     sql.NullString
	CreatedBy      sql.NullInt64
	CreatedAt      sql.NullTime
	CreatedClient  sql.NullString
	UpdatedAt      sql.NullTime
}

type ListJobProcessModel struct {
	Level     sql.NullInt32
	JobID     sql.NullString
	Group     sql.NullString
	Type      sql.NullString
	Name      sql.NullString
	Counter   sql.NullInt32
	Total     sql.NullInt32
	Status    sql.NullString
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

type ViewJobProcessModel struct {
	ID              sql.NullInt64
	ParentJobID     sql.NullString
	Level           sql.NullInt32
	JobID           sql.NullString
	Group           sql.NullString
	Type            sql.NullString
	Name            sql.NullString
	Counter         sql.NullInt32
	Total           sql.NullInt32
	Status          sql.NullString
	UpdatedAt       sql.NullTime
	CreatedAt       sql.NullTime
	UrlIn           sql.NullString
	FileNameIn      sql.NullString
	ContentDataOut  sql.NullString
	ChildJobProcess ChildJobProcess
}

type ContentDataOutDetail struct {
	ID      int64  `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ChildJobProcess []struct {
	Level          int     `json:"level"`
	JobID          string  `json:"job_id"`
	Group          string  `json:"group"`
	Type           string  `json:"type"`
	Name           string  `json:"name"`
	UrlIn          string  `json:"url_in"`
	FileNameIn     string  `json:"file_name_in"`
	ContentDataOut string  `json:"content_data_out"`
	Counter        int     `json:"counter"`
	Total          int     `json:"total"`
	Status         string  `json:"status"`
	UpdatedAt      string  `json:"updated_at"`
	CreatedAt      string  `json:"created_at"`
	Duration       float64 `json:"duration"`
}

type GetListFileUploadJobProcess struct {
	JobID       sql.NullString
	Status      sql.NullString
	Progress    sql.NullFloat64
	FileUrl     sql.NullString
	CreatedName sql.NullString
	CreatedAt   sql.NullTime
}

func GenerateAddResourceNexcloudUserJobProcessModel(level int32, timeNow time.Time, createdBy int64) JobProcessModel {
	return generateJobProcessTask(level, constanta.JobProcessUserGroup, constanta.JobProcessAddResourceType, constanta.JobProcessAddResourceNexcloudName, timeNow, createdBy)
}

func GenerateUpdateLogAfterAddResourceUserJobProcessModel(level int32, uuid string, timeNow time.Time, createdBy int64) JobProcessModel {
	return generateChildJobProcessTask(level, uuid, constanta.JobProcessUserGroup, constanta.JobProcessAddResourceType, constanta.JobProcessUpdateLogAfterAddResourceNexcloudName, timeNow, createdBy)
}

func generateJobProcessTask(level int32, group string, jobType string, name string, timeNow time.Time, createdBy int64) JobProcessModel {
	return JobProcessModel{
		Level:     sql.NullInt32{Int32: level},
		JobID:     sql.NullString{String: util.GetUUID()},
		Group:     sql.NullString{String: group},
		Type:      sql.NullString{String: jobType},
		Name:      sql.NullString{String: name},
		Status:    sql.NullString{String: constanta.JobProcessOnProgressStatus},
		CreatedBy: sql.NullInt64{Int64: createdBy},
		CreatedAt: sql.NullTime{Time: timeNow},
		UpdatedAt: sql.NullTime{Time: timeNow},
	}
}

func generateChildJobProcessTask(level int32, uuid string, group string, jobType string, name string, timeNow time.Time, createdBy int64) JobProcessModel {
	return JobProcessModel{
		Level:       sql.NullInt32{Int32: level},
		ParentJobID: sql.NullString{String: uuid},
		JobID:       sql.NullString{String: util.GetUUID()},
		Group:       sql.NullString{String: group},
		Type:        sql.NullString{String: jobType},
		Name:        sql.NullString{String: name},
		Status:      sql.NullString{String: constanta.JobProcessOnProgressStatus},
		CreatedBy:   sql.NullInt64{Int64: createdBy},
		CreatedAt:   sql.NullTime{Time: timeNow},
		UpdatedAt:   sql.NullTime{Time: timeNow},
	}
}
