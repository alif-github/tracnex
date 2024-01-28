package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type permissionDAO struct {
	AbstractDAO
}

var PermissionDAO = permissionDAO{}.New()

func (input permissionDAO) New() (output permissionDAO) {
	output.FileName = "PermissionDAO.go"
	output.TableName = "permission"
	return
}

func (input permissionDAO) CheckIsPermissionValid(db *sql.DB, dataPermission []string) (result int, err errorModel.ErrorModel) {
	funcName := "CheckIsPermissionValid"
	listInterface := ArrayStringToArrayInterface(dataPermission)
	inQuery := ListDataToInQuery(listInterface)

	query :=
		"SELECT " +
			"COUNT(id) " +
			"FROM " +
			"permission " +
			"WHERE " +
			"permission IN(" + inQuery + ")"

	results := db.QueryRow(query, listInterface...)

	errorS := results.Scan(&result)

	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}