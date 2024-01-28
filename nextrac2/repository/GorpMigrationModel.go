package repository

import "database/sql"

type MigrationModel struct {
	ID sql.NullString
}