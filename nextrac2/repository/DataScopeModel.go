package repository

import "database/sql"

type DataScopeModel struct {
	ID            sql.NullInt64
	UUIDKey       sql.NullString
	Scope         sql.NullString
	Description   sql.NullString
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	CreatedAt     sql.NullTime
	UpdatedBy     sql.NullInt64
	UpdatedClient sql.NullString
	UpdatedAt     sql.NullTime
	Deleted       sql.NullBool
}

type DataScopeForDataGroupModel struct {
	Scope string `json:"scope"`
}

type MapOfDataScopeForDataGroupModel struct {
	DataScope      map[string][]string `json:"data_scope"`
	DataGroupScope map[string][]string `json:"data_group_scope"`
}


