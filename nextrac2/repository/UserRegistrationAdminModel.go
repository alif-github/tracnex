package repository

import "database/sql"

type UserRegistrationAdminModel struct {
	ID                 sql.NullInt64
	CustomerName       sql.NullString
	ParentCustomerName sql.NullString
	UniqueID1          sql.NullString
	UniqueID2          sql.NullString
	CompanyName        sql.NullString
	BranchName         sql.NullString
	UserAdmin          sql.NullString
	PasswordAdmin      sql.NullString
	ParentCustomerId   sql.NullInt64
	CustomerId         sql.NullInt64
	SiteId             sql.NullInt64
	ClientID           sql.NullString
	ClientMappingID    sql.NullInt64
	ClientTypeID       sql.NullInt64
	CreatedBy          sql.NullInt64
	CreatedClient      sql.NullString
	CreatedAt          sql.NullTime
	UpdatedBy          sql.NullInt64
	UpdatedClient      sql.NullString
	UpdatedAt          sql.NullTime
}
