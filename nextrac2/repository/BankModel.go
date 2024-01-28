package repository

import "database/sql"

type BankModel struct {
	ID            sql.NullInt64
	UUIDKey       sql.NullString
	Name          sql.NullString
	Status        sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	Deleted       sql.NullBool
}

type BankElasticModel struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedBy int64  `json:"created_by"`
}
