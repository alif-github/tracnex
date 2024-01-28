package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type serverRunDAO struct {
	AbstractDAO
}

var ServerRunDAO = serverRunDAO{}.New()

func (input serverRunDAO) New() (output serverRunDAO) {
	output.FileName = "ServerRunDAO.go"
	output.TableName = "server_run"
	return
}

func (input serverRunDAO) CheckIsServerRunExist(db *sql.DB, userParam repository.ServerRunModel) (result repository.ServerRunModel, err errorModel.ErrorModel) {
	funcName := "CheckIsServerRunExist"
	query :=
		"SELECT " +
			"	id, created_by, updated_at " +
			"FROM " +
			"	server_run " +
			"WHERE " +
			"	host_id = $1 " +
			"AND " +
			"	run_type = $2 AND deleted = FALSE"

	param := []interface{}{userParam.HostID.Int64, userParam.RunType.String}
	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $3 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input serverRunDAO) InsertServerRun(db *sql.Tx, userParam repository.ServerRunModel) (result int64, err errorModel.ErrorModel){
	funcName := "InsertServerRun"
	query :=
		" INSERT INTO server_run " +
			"	(name, run_type, host_id, " +
			"	created_by, created_client, updated_by, " +
			"	updated_client, created_at, updated_at, " +
			"	status) "+
			"VALUES " +
			"	($1, $2, $3," +
			"	$4, $5, $6, " +
			"	$7, $8, $9, " +
			"	$10) " +
			" RETURNING id "

	param := []interface{}{
		userParam.Name.String, userParam.RunType.String, userParam.HostID.Int64, userParam.CreatedBy.Int64,
		userParam.CreatedClient.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.CreatedAt.Time, userParam.UpdatedAt.Time, userParam.Status.Bool}

	errorS := db.QueryRow(query, param...).Scan(&result)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input serverRunDAO) UpdateServerRun(db *sql.Tx, userParam repository.ServerRunModel, timeNow time.Time) (err errorModel.ErrorModel){
	funcName := "UpdateServerRun"
	query :=
		"UPDATE server_run " +
			"	SET status = $1, updated_at = $2, updated_client = $3, " +
			"	updated_by = $4 " +
			"WHERE " +
			"	id = $5 AND deleted = FALSE "

	param := []interface{}{
		userParam.Status.Bool, timeNow, userParam.UpdatedClient.String,
		userParam.UpdatedBy.Int64, userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $6 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	return errorModel.GenerateNonErrorModel()

}