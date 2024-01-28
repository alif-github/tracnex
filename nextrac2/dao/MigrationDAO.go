package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type migrationDAO struct {
	AbstractDAO
}

var MigrationDAO = migrationDAO{}.New()

func (input migrationDAO) New() (output migrationDAO) {
	output.FileName = "MigrationDAO.go"
	output.TableName = "gorp_migrations"
	return
}

func (input migrationDAO) ResetMigration(db *sql.Tx, userParam repository.MigrationModel) (err errorModel.ErrorModel) {
	var (
		funcName = "ResetMigration"
	)

	query := fmt.Sprintf(`DELETE FROM %s `, input.TableName)
	params := []interface{}{}

	if !util.IsStringEmpty(userParam.ID.String) {
		query += "WHERE id = $1"
		params = append(params, userParam.ID.String)
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	return
}
