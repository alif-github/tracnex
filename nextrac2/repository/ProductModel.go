package repository

import "database/sql"

type ProductModel struct {
	ID                    sql.NullInt64
	ProductID             sql.NullString
	ProductName           sql.NullString
	ProductDescription    sql.NullString
	ProductGroupID        sql.NullInt64
	ProductGroupName      sql.NullString
	ParentClientTypeID    sql.NullInt64
	ClientTypeID          sql.NullInt64
	ClientTypeName        sql.NullString
	IsLicense             sql.NullBool
	LicenseVariantID      sql.NullInt64
	LicenseVariantName    sql.NullString
	LicenseTypeID         sql.NullInt64
	LicenseTypeName       sql.NullString
	DeploymentMethod      sql.NullString
	NoOfUser              sql.NullInt64
	IsUserConcurrent      sql.NullBool
	Module1               sql.NullInt64
	ModuleName1           sql.NullString
	Module2               sql.NullInt64
	ModuleName2           sql.NullString
	Module3               sql.NullInt64
	ModuleName3           sql.NullString
	Module4               sql.NullInt64
	ModuleName4           sql.NullString
	Module5               sql.NullInt64
	ModuleName5           sql.NullString
	Module6               sql.NullInt64
	ModuleName6           sql.NullString
	Module7               sql.NullInt64
	ModuleName7           sql.NullString
	Module8               sql.NullInt64
	ModuleName8           sql.NullString
	Module9               sql.NullInt64
	ModuleName9           sql.NullString
	Module10              sql.NullInt64
	ModuleName10          sql.NullString
	CreatedBy             sql.NullInt64
	CreatedClient         sql.NullString
	CreatedAt             sql.NullTime
	UpdatedBy             sql.NullInt64
	UpdatedClient         sql.NullString
	UpdatedAt             sql.NullTime
	UpdatedName           sql.NullString
	ProductComponentModel []ProductComponentModel
	MaxOfflineDays        sql.NullInt64
}

type ProductComponentModel struct {
	ID             sql.NullInt64
	ComponentID    sql.NullInt64
	ComponentName  sql.NullString
	ComponentValue sql.NullString
	UpdatedBy      sql.NullInt64
	UpdatedClient  sql.NullString
	UpdatedAt      sql.NullTime
	Deleted        sql.NullBool
}

type GetListProductComponent struct {
	ID             sql.NullInt64
	ProductID      sql.NullInt64
	ComponentID    sql.NullInt64
	ComponentName  sql.NullString
	ComponentValue sql.NullString
}

type GetForUpdateProduct struct {
	ID                 sql.NullInt64
	ProductID          sql.NullString
	ProductName        sql.NullString
	UpdatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	ProductComponentID sql.NullInt64
	IsUsed             sql.NullBool
}
