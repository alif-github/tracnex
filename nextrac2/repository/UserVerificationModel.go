package repository

import "database/sql"

type UserVerificationModel struct {
	ID                       sql.NullInt64
	UserRegistrationDetailID sql.NullInt64
	Email                    sql.NullString
	EmailCode                sql.NullString
	EmailExpires             sql.NullInt64
	Phone                    sql.NullString
	PhoneCode                sql.NullString
	PhoneExpires             sql.NullInt64
	CreatedBy                sql.NullInt64
	CreatedAt                sql.NullTime
	CreatedClient            sql.NullString
	UpdatedBy                sql.NullInt64
	UpdatedAt                sql.NullTime
	UpdatedClient            sql.NullString
	FailedOTPEmail           sql.NullInt64
	FailedOTPPhone           sql.NullInt64
}
