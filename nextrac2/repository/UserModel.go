package repository

import "database/sql"

type UserModel struct {
	ID              sql.NullInt64
	UUIDKey         sql.NullString
	ClientID        sql.NullString
	AuthUserID      sql.NullInt64
	IPWhitelist     sql.NullString
	Locale          sql.NullString
	SignatureKey    sql.NullString
	AdditionalInfo  sql.NullString
	LastToken       sql.NullTime
	IsSystemAdmin   sql.NullBool
	Status          sql.NullString
	FirstName       sql.NullString
	LastName        sql.NullString
	Username        sql.NullString
	Email           sql.NullString
	Phone           sql.NullString
	CreatedBy       sql.NullInt64
	CreatedAt       sql.NullTime
	CreatedClient   sql.NullString
	UpdatedBy       sql.NullInt64
	UpdatedAt       sql.NullTime
	UpdatedClient   sql.NullString
	Deleted         sql.NullBool
	AliasName       sql.NullString
	IdCard          sql.NullString
	EmployeeId      sql.NullInt64
	Position        sql.NullString
	Department      sql.NullString
	IsHaveMember    sql.NullBool
	Currency        sql.NullString
	PlatformDevice  sql.NullString
	EmployeeLevelId sql.NullInt64
	EmployeeGradeId sql.NullInt64
}

type ListUserModel struct {
	ID             sql.NullInt64
	AuthUserID     sql.NullInt64
	ClientID       sql.NullString
	Status         sql.NullString
	Locale         sql.NullString
	Email          sql.NullString
	Phone          sql.NullString
	FirstName      sql.NullString
	Username       sql.NullString
	IPWhiteList    sql.NullString
	LastName       sql.NullString
	CreatedName    sql.NullString
	CreatedBy      sql.NullInt64
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
	RoleID         sql.NullString
	GroupID        sql.NullString
	PlatformDevice sql.NullString
}

type ViewDetailUserModel struct {
	ID               sql.NullInt64
	ClientID         sql.NullString
	Username         sql.NullString
	Email            sql.NullString
	Phone            sql.NullString
	FirstName        sql.NullString
	LastName         sql.NullString
	Role             sql.NullString
	GroupID          sql.NullString
	IsAdmin          sql.NullBool
	Status           sql.NullString
	CreatedBy        sql.NullInt64
	CreatedFirstName sql.NullString
	CreatedLastName  sql.NullString
	CreatedAt        sql.NullTime
	UpdatedBy        sql.NullInt64
	UpdatedFirstName sql.NullString
	UpdatedLastName  sql.NullString
	UpdatedAt        sql.NullTime
	PlatformDevice   sql.NullString
	Currency         sql.NullString
	IsVerifyPhone    sql.NullBool
	IsVerifyEmail    sql.NullBool
}

type RoleMappingPersonProfileModel struct {
	PersonProfileID sql.NullInt64
	AuthUserID      sql.NullInt64
	RoleName        sql.NullString
	Permissions     sql.NullString
	GroupName       sql.NullString
	Scope           sql.NullString
	IPWhitelist     sql.NullString
	SignatureKey    sql.NullString
	Locale          sql.NullString
	IsAdmin         sql.NullBool
}
