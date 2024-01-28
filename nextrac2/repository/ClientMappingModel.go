package repository

import "database/sql"

type ClientMappingModel struct {
	ID               sql.NullInt64
	ClientID         sql.NullString
	ClientTypeID     sql.NullInt64
	InstallationID   sql.NullInt64
	CustomerID       sql.NullInt64
	SiteID           sql.NullInt64
	CompanyID        sql.NullString
	BranchID         sql.NullString
	ClientAlias      sql.NullString
	SocketID         sql.NullString
	CreatedBy        sql.NullInt64
	CreatedClient    sql.NullString
	CreatedAt        sql.NullTime
	UpdatedBy        sql.NullInt64
	UpdatedClient    sql.NullString
	UpdatedAt        sql.NullTime
	AuthUserId       sql.NullString
	UserName         sql.NullString
	ParentCustomerID sql.NullInt64
}

type ClientMappingForDetailModel struct {
	ClientTypeID sql.NullInt64
	CompanyData  []CompanyDataModel
	ClientData   []ClientDataModel
}

type CompanyDataModel struct {
	CompanyID  sql.NullString
	BranchData []BranchDataModel
}

type BranchDataModel struct {
	BranchID sql.NullString
}

type CLientMappingDetailForViewModel struct {
	ID             sql.NullInt64
	ClientId       sql.NullString
	ClientTypeId   sql.NullInt64
	CompanyId      sql.NullString
	BranchId       sql.NullString
	Aliases        sql.NullString
	AuthUserId     sql.NullInt64
	Username       sql.NullString
	SocketPassword sql.NullString
	SocketID       sql.NullString
	CreatedAt      sql.NullTime
	CreatedBy      sql.NullInt64
	UpdatedAt      sql.NullTime
	UpdatedBy      sql.NullInt64
}

type ClientMappingForRemoveModel struct {
	CompanyID   sql.NullString
	BranchID    sql.NullString
	BranchName  sql.NullString
	ClientAlias sql.NullString
}

type ClientMappingForViewModel struct {
	ID                    sql.NullInt64
	ClientID              sql.NullString
	SocketID              sql.NullString
	ClientType            sql.NullString
	CompanyID             sql.NullString
	BranchID              sql.NullString
	Aliases               sql.NullString
	SuccessStatusAuth     sql.NullBool
	SuccessStatusNexcloud sql.NullBool
	SuccessStatusNexdrive sql.NullBool
	CreatedAt             sql.NullTime
	CreatedBy             sql.NullInt64
	UpdatedAt             sql.NullTime
	UpdatedBy             sql.NullInt64
}

type ClientDataModel struct {
	ClientID sql.NullString
}
