package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type cronSchedulerDAO struct {
	AbstractDAO
}

var CronSchedulerDAO = cronSchedulerDAO{}.New()

func (input cronSchedulerDAO) New() (output cronSchedulerDAO) {
	output.FileName = "CronSchedulerDAO.go"
	output.TableName = "cron_scheduler"
	return
}

func (input cronSchedulerDAO) GetListCronScheduler(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (result []interface{}, class errorModel.ErrorModel) {
	query :=
		" SELECT " +
			"	id, name, run_type, created_by, updated_at, status " +
			"FROM " +
			"	cron_scheduler "

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CRONSchedulerModel
			errorS := rows.Scan(&temp.ID, &temp.Name, &temp.RunType, &temp.CreatedBy, &temp.UpdatedAt, &temp.Status)
			return temp, errorS
		}, "AND status = TRUE ", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
}

func (input cronSchedulerDAO) CountCronScheduler(db *sql.DB, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (int, errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
}

func (input cronSchedulerDAO) ViewCronScheduler(db *sql.DB, userParam repository.CRONSchedulerModel) (result repository.CRONSchedulerModel, err errorModel.ErrorModel) {
	funcName := "ViewCronScheduler"
	query :=
		"SELECT " +
			"	id, name, run_type, cron, created_by, updated_at, status " +
			"FROM " +
			"	cron_scheduler " +
			"WHERE " +
			"	id = $1"

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 != 0 {
		query += "AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.Name, &result.RunType, &result.CRON, &result.CreatedBy, &result.UpdatedAt, &result.Status)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input cronSchedulerDAO) GetCronSchedulerForUpdate(db *sql.Tx, userParam repository.CRONSchedulerModel) (result repository.CRONSchedulerModel, err errorModel.ErrorModel) {
	funcName := "GetCronSchedulerForUpdate"
	query :=
		"SELECT " +
			"	id, created_by, updated_at " +
			"FROM " +
			" cron_scheduler " +
			"WHERE " +
			"	id = $1 AND deleted = FALSE"

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += "AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input cronSchedulerDAO) UpdatedCronScheduler(db *sql.Tx, userParam repository.CRONSchedulerModel, timeNow time.Time) (err errorModel.ErrorModel){
	funcName := "UpdatedCronScheduler"
	query :=
		"UPDATE cron_scheduler " +
			"	SET cron = $1, status = $2, updated_at = $3, " +
			"	updated_client = $4, updated_by = $5 " +
			"WHERE " +
			"	id = $6 AND deleted = FALSE "

	param := []interface{}{
		userParam.CRON.String, userParam.Status.Bool, timeNow,
		userParam.UpdatedClient.String,	userParam.UpdatedBy.Int64, userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $7 "
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

func (input cronSchedulerDAO) GetCronScheduler(db *sql.Tx, userParam repository.CRONSchedulerModel) (result repository.CRONSchedulerModel, err errorModel.ErrorModel) {
	funcName := "GetCronSchedulerForUpdate"
	query :=
		"SELECT " +
			"	id, name, run_type, created_by, updated_at " +
			"FROM " +
			" cron_scheduler " +
			"WHERE " +
			"	id = $1 AND deleted = FALSE"

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += "AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.Name, &result.RunType, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input cronSchedulerDAO) UpdatedCronSchedulerByRunType(db *sql.Tx, userParam repository.CRONSchedulerModel) (err errorModel.ErrorModel){
	funcName := "UpdatedCronSchedulerByRunType"

	query :=
		"UPDATE cron_scheduler " +
			"	SET cron = $1 WHERE " +
			"	run_type = $2 AND deleted = FALSE "

	param := []interface{}{
		userParam.CRON.String, userParam.RunType.String,}

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