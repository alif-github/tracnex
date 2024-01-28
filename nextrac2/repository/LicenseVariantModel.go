package repository

import "database/sql"

type LicenseVariantModel struct {
	ID                 sql.NullInt64
	LicenseVariantName sql.NullString
	CreatedBy          sql.NullInt64
	CreatedClient      sql.NullString
	CreatedAt          sql.NullTime
	UpdatedBy          sql.NullInt64
	UpdatedClient      sql.NullString
	UpdatedAt          sql.NullTime
	UpdatedName        sql.NullString
	IsUsed             sql.NullBool
}

type LicenseVariantListModel struct {
	ID                 sql.NullInt64
	LicenseVariantName sql.NullString
	CreatedAt          sql.NullTime
	UpdatedBy          sql.NullInt64
	UpdatedName        sql.NullString
	UpdatedAt          sql.NullTime
}
