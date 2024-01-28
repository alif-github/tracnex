package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type cronHostDAO struct {
	AbstractDAO
}

var CronHostDAO = cronHostDAO{}.New()

func (input cronHostDAO) New() (output cronHostDAO) {
	output.FileName = "CronHostDAO.go"
	output.TableName = "cron_host"
	return
}

func (input cronHostDAO) CheckIsCronHostExist(db *sql.DB, userParam repository.CRONHostModel) (result repository.CRONHostModel, err errorModel.ErrorModel){
	funcName := "checkIsCronHostExist"
	query :=
		"SELECT " +
			"	id, created_by, updated_at " +
			"FROM " +
			"	cron_host " +
			"WHERE " +
			"	host_id = $1 " +
			"AND " +
			"	cron_id = $2 AND deleted = FALSE"

	param := []interface{}{userParam.HostID.Int64, userParam.CronID.Int64}
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

func (input cronHostDAO) InsertCronHost(db *sql.Tx, userParam repository.CRONHostModel) (result int64, err errorModel.ErrorModel){
	funcName := "InsertCronHost"
	query :=
		" INSERT INTO cron_host " +
			"	(cron_id, host_id, created_by, " +
			"	created_client, updated_by, updated_client, " +
			"	created_at, updated_at, status) " +
			"VALUES " +
			"	($1, $2, $3," +
			"	$4, $5, $6, " +
			"	$7, $8, $9) " +
			" RETURNING id "

	param := []interface{}{
		userParam.CronID.Int64, userParam.HostID.Int64, userParam.CreatedBy.Int64,
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

func (input cronHostDAO) UpdateCronHost(db *sql.Tx, userParam repository.CRONHostModel, timeNow time.Time) (err errorModel.ErrorModel) {
	funcName := "UpdateCronHost"
	query :=
		"UPDATE cron_host " +
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