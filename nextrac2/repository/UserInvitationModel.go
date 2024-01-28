package repository

import "database/sql"

type UserInvitation struct {
	Id             sql.NullInt64
	InvitationCode sql.NullString
	Email          sql.NullString
	ClientId	   sql.NullString
	RoleId         sql.NullInt64
	DataGroupId    sql.NullInt64
	ExpiresOn      sql.NullTime
	CreatedClient  sql.NullString
	UpdatedClient  sql.NullString
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
	CreatedBy      sql.NullInt64
	UpdatedBy      sql.NullInt64
	Deleted        sql.NullBool
}
