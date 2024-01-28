package repository

import (
	"database/sql"
)

type ClientRegistrationLogModel struct {
	ID                    sql.NullInt64
	ClientID              sql.NullString
	ClientTypeID          sql.NullInt64
	AttributeRequest      sql.NullString
	SuccessStatusAuth     sql.NullBool
	SuccessStatusNexcloud sql.NullBool
	SuccessStatusNexdrive sql.NullBool
	SuccessStatus		  sql.NullBool //for success status without dependent resource
	Resource              sql.NullString
	MessageAuth           sql.NullString
	MessageNexcloud       sql.NullString
	MessageNexdrive       sql.NullString
	Message				  sql.NullString //for message without dependent resource
	Details               sql.NullString
	Code                  sql.NullString
	RequestTimeStamp      sql.NullTime
	CreatedBy             sql.NullInt64
	CreatedAt             sql.NullTime
	CreatedClient         sql.NullString
	UpdatedBy             sql.NullInt64
	UpdatedAt             sql.NullTime
	UpdatedClient         sql.NullString
	RequestCount		  sql.NullInt64
}

type ParamClientRegistrationLogModel struct {
	Code				sql.NullString
	Detail				sql.NullString
	ClientTypeID		sql.NullInt64
}

type ViewClientRegistrationLogModel struct {
	ClientID			sql.NullString
	AuthUserID			sql.NullInt64
	ClientType			sql.NullString
	Status				sql.NullString
	FirstName			sql.NullString
	LastName			sql.NullString
	Resource			sql.NullString
	UpdatedAt			sql.NullString
}

type ListClientRegistrationLogModel struct {
	ID						sql.NullInt64
	ClientID				sql.NullString
	ClientTypeID			sql.NullInt64
	SuccessStatusAuth		sql.NullBool
	SuccessStatusNexcloud	sql.NullBool
	Resource				sql.NullString
}