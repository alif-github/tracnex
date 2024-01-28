package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type schedulerDAO struct {
	AbstractDAO
}

var SchedulerDAO = schedulerDAO{}.New()

func (input schedulerDAO) New() (output schedulerDAO) {
	output.FileName = "SchedulerDAO.go"
	return
}

func (input schedulerDAO) GetListSchedulerByHostName(db *sql.DB, hostName string) (result []repository.CRONSchedulerModel, error errorModel.ErrorModel) {
	funcName := "GetListSchedulerByHostName"
	query :=
		"SELECT " +
			"	cs.run_type, cs.cron " +
			"FROM " +
			"	host_server hs LEFT JOIN cron_host ch ON ch.host_id = hs.id LEFT JOIN " +
			"	cron_scheduler cs ON cs.id = ch.cron_id " +
			"WHERE " +
			"	hs.host_name = $1 AND hs.deleted = FALSE AND " +
			"	ch.deleted = FALSE AND ch.status = TRUE AND " +
			"	cs.deleted = FALSE AND cs.status = TRUE"

	rows, err := db.Query(query, hostName)
	if err != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}
	if rows != nil {
		defer func() {
			err = rows.Close()
			if err != nil {
				error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
			}
		}()
		for rows.Next() {
			var temp repository.CRONSchedulerModel
			err = rows.Scan(&temp.RunType, &temp.CRON)
			if err != nil {
				error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
				return
			}
			result = append(result, temp)
		}
	} else {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}

func (input schedulerDAO) GetServerRunByHostNameAndRunType(db *sql.DB, hostName string, runType string) (result repository.ServerRunModel, error errorModel.ErrorModel) {
	funcName := "GetServerRunByHostNameAndRunType"
	query :=
		"SELECT " +
			"	sr.id " +
			"FROM " +
			"	server_run sr LEFT JOIN host_server hs ON sr.host_id = hs.id " +
			"WHERE " +
			"	sr.run_type = $1 AND hs.host_name = $2 AND " +
			"	sr.status = TRUE and sr.deleted = FALSE "

	_rows := db.QueryRow(query, runType, hostName)
	err := _rows.Scan(&result.ID)
	if err != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}

func (input schedulerDAO) GetDataForRestartScheduler(db *sql.DB, runType string) (result []repository.RefreshScheduler, error errorModel.ErrorModel) {
	funcName := "GetDataForRestartScheduler"
	query :=
		"SELECT " +
			"	host_name, host_url, cs.cron " +
			"FROM " +
			"	cron_host ch LEFT JOIN cron_scheduler cs ON ch.cron_id = cs.id LEFT JOIN " +
			"	host_server hs ON ch.host_id = hs.id " +
			"WHERE " +
			"	cs.run_type = $1 "

	rows, err := db.Query(query, runType)
	if err != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}
	if rows != nil {
		defer func() {
			err = rows.Close()
			if err != nil {
				error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
			}
		}()
		for rows.Next() {
			var temp repository.RefreshScheduler
			err = rows.Scan(&temp.HostName, &temp.HostURL, &temp.CRON)
			if err != nil {
				error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
				return
			}
			result = append(result, temp)
		}
	} else {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}

func (input schedulerDAO) GetDataForRestartOwnScheduler(db *sql.DB, runType string, hostName string) (result repository.RefreshScheduler, error errorModel.ErrorModel) {
	funcName := "GetDataForRestartOwnScheduler"
	query :=
		"SELECT " +
			"	cs.cron " +
			"FROM " +
			"	cron_host ch LEFT JOIN cron_scheduler cs ON ch.cron_id = cs.id LEFT JOIN " +
			"	host_server hs ON ch.host_id = hs.id " +
			"WHERE " +
			"	cs.run_type = $1 AND host_name = $2 AND ch.status = TRUE AND ch.deleted = FALSE "

	_rows := db.QueryRow(query, runType, hostName)
	err := _rows.Scan(&result.CRON)
	if err != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}
