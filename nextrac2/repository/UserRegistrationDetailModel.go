package repository

import "database/sql"

type UserRegistrationDetailModel struct {
	ID               sql.NullInt64
	UserRegDetailID  sql.NullInt64
	UserLicenseID    sql.NullInt64
	ParentCustomerID sql.NullInt64
	CustomerID       sql.NullInt64
	SiteID           sql.NullInt64
	InstallationID   sql.NullInt64
	ParentClientID   sql.NullString
	ClientID         sql.NullString
	UniqueID1        sql.NullString
	UniqueID2        sql.NullString
	AuthUserID       sql.NullInt64
	Firstname        sql.NullString
	Lastname         sql.NullString
	Username         sql.NullString
	UserID           sql.NullString
	Password         sql.NullString
	ClientAliases    sql.NullString
	SalesmanID       sql.NullString
	AndroidID        sql.NullString
	RegDate          sql.NullTime
	Status           sql.NullString
	Email            sql.NullString
	NoTelp           sql.NullString
	NexmileOTP       sql.NullString
	SalesmanCategory sql.NullString
	ProductValidFrom sql.NullTime
	ProductValidThru sql.NullTime
	CreatedClient    sql.NullString
	CreatedBy        sql.NullInt64
	CreatedAt        sql.NullTime
	UpdatedClient    sql.NullString
	UpdatedBy        sql.NullInt64
	UpdatedAt        sql.NullTime
	SalesmanId       sql.NullString
	ClientTypeID     sql.NullInt64
	CustomerName     sql.NullString
	MaxOfflineDays   sql.NullInt64
	LicenseStatus    sql.NullInt64
	LicenseConfigID  sql.NullInt64
}

type UserRegistrationDetailMapping struct {
	UserRegistrationDetail UserRegistrationDetailModel
	UserVerification       UserVerificationModel
	PKCEClientMapping      PKCEClientMappingModel
	ClientMapping          ClientMappingModel
	User                   UserModel
}
