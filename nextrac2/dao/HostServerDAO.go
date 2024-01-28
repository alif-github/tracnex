package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"time"
)

type hostServerDAO struct {
	AbstractDAO
}

var HostServerDAO = hostServerDAO{}.New()

func (input hostServerDAO) New() (output hostServerDAO) {
	output.FileName = "HostServerDAO.go"
	output.TableName = "host_server"
	return
}

func (input hostServerDAO) getHostServerDefaultMustCheck(createdBy int64, isChecked bool) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "hs.id"},
		Deleted:   FieldStatus{FieldName: "hs.deleted", IsCheck: isChecked},
		CreatedBy: FieldStatus{FieldName: "hs.created_by", Value: createdBy},
	}
}

func (input hostServerDAO) checkOwnPermission(createdBy int64, query *string, param *[]interface{}, index int) int {
	if createdBy > 0 {
		queryOwnPermission := " AND created_by = $" + strconv.Itoa(index) + " "
		*query += queryOwnPermission
		(*param) = append((*param), createdBy)
		index += 1
	}
	return index
}

func (input hostServerDAO) GetHostByHostName(db *sql.Tx, userParam repository.HostServerModel) (result repository.HostServerModel, error errorModel.ErrorModel) {
	funcName := "GetHostByHostName"
	query := fmt.Sprintf(
		`SELECT
			id
		FROM
			%s
		WHERE
			host_name = $1 `, input.TableName)

	results := db.QueryRow(query, userParam.HostName.String)

	err := results.Scan(&result.ID)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerDAO) InsertHostnameAndIP(db *sql.Tx, userParam repository.HostServerModel) (result int64, error errorModel.ErrorModel) {
	funcName := "InsertHostnameAndIP"
	query := fmt.Sprintf(
		`INSERT INTO %s
			(host_name, host_url, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4) RETURNING id `, input.TableName)

	results := db.QueryRow(query, userParam.HostName.String, userParam.HostURL.String, userParam.CreatedAt.Time, userParam.UpdatedAt.Time)

	err := results.Scan(&result)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerDAO) AutoUpdateHostnameAndIP(db *sql.Tx, userParam repository.HostServerModel) errorModel.ErrorModel {
	funcName := "AutoUpdateHostnameAndIP"
	query := fmt.Sprintf(
		`UPDATE %s set
				host_url = $1, updated_at = $2, deleted = FALSE
			WHERE
				id = $3 `, input.TableName)


	stmt, err := db.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(userParam.HostURL.String, userParam.UpdatedAt.Time, userParam.ID.Int64)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input hostServerDAO) GetListHostServer(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (result []interface{}, class errorModel.ErrorModel) {
	query := fmt.Sprintf(
		`SELECT
			id, host_name, created_by, updated_at
		FROM %s `, input.TableName )


	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.HostServerModel
			errorS := rows.Scan(&temp.ID, &temp.HostName, &temp.CreatedBy, &temp.UpdatedAt)
			return temp, errorS
		}, "", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
}

func (input hostServerDAO) CountHostServer(db *sql.DB, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (int, errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
}

func (input hostServerDAO) GetHostServerForUpdate(db *sql.Tx, userParam repository.HostServerModel) (result repository.HostServerModel, err errorModel.ErrorModel) {
	funcName := "GetHostServerForUpdate"
	query := fmt.Sprintf(
		`SELECT
			id, created_by, updated_at
		FROM
			%s
		WHERE
			id = $1 AND deleted = FALSE `, input.TableName)


	param := []interface{}{userParam.ID.Int64}

	_ = input.checkOwnPermission(userParam.CreatedBy.Int64, &query, &param, 2)

	query += " FOR UPDATE"

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerDAO) UpdateHostnameAndIP(db *sql.Tx, userParam repository.HostServerModel, timeNow time.Time) (err errorModel.ErrorModel) {
	funcName := "UpdateHostnameAndIP"
	query := fmt.Sprintf(
		`UPDATE %s
		SET 
			host_name = $1, host_url = $2, updated_at = $3,
			updated_client = $4, updated_by = $5
		WHERE
			id = $6 AND deleted = FALSE `, input.TableName)


	param := []interface{}{
		userParam.HostName.String, userParam.HostURL.String, timeNow, userParam.UpdatedClient.String,
		userParam.UpdatedBy.Int64, userParam.ID.Int64}

	_ = input.checkOwnPermission(userParam.CreatedBy.Int64, &query, &param, 7)

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

func (input hostServerDAO) DeleteHostServer(db *sql.Tx, userParam repository.HostServerModel, timeNow time.Time) (err errorModel.ErrorModel) {
	funcName := "DeleteHostServer"
	query := fmt.Sprintf(
		`UPDATE %s
		SET 
			deleted = TRUE, updated_at = $1, updated_client = $2,
			updated_by = $3
		WHERE
			id = $4 AND deleted = FALSE `, input.TableName)

	param := []interface{}{
		timeNow, userParam.UpdatedClient.String, userParam.UpdatedBy.Int64,
		userParam.ID.Int64}

	_ = input.checkOwnPermission(userParam.CreatedBy.Int64, &query, &param, 5)

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

func (input hostServerDAO) CheckIsHostServerUsed(db *sql.DB, userParam repository.HostServerModel) (result bool, err errorModel.ErrorModel) {
	funcName := "CheckIsHostServerUsed"
	query := fmt.Sprintf(
		`SELECT
			CASE WHEN
				count_cron_host count_server_run= 0
			THEN FALSE ELSE TRUE END
		FROM
			(SELECT
				CASE WHEN
				  (SELECT 1 FROM cron_host WHERE host_id = $1 AND DELETED = false LIMIT 1)
			is NULL THEN 0 ELSE 1 END AS count_cron_host) a,
			(SELECT
				CASE WHEN
				  (SELECT 1 FROM server_run WHERE host_id = $1 AND DELETED = false LIMIT 1)
			is NULL THEN 0 ELSE 1 END AS count_server_run) b `)

	param := []interface{}{userParam.ID.Int64}

	errorS := db.QueryRow(query, param...).Scan(&result)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerDAO) GetHostServer(db *sql.Tx, userParam repository.HostServerModel) (result repository.HostServerModel, err errorModel.ErrorModel) {
	funcName := "GetHostServerForUpdate"
	query := fmt.Sprintf(
		`SELECT
			id, created_by, updated_at
		FROM
			%s
		WHERE
			id = $1 AND deleted = FALSE `, input.TableName)


	param := []interface{}{userParam.ID.Int64}

	_ = input.checkOwnPermission(userParam.CreatedBy.Int64, &query, &param, 2)

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerDAO) GetHostServerDb(db *sql.DB, userParam repository.HostServerModel) (result repository.HostServerModel, err errorModel.ErrorModel) {
	funcName := "GetHostServerDb"
	query := fmt.Sprintf(
		`SELECT
			id, created_by, updated_at,
			host_name, host_url
		FROM
			%s
		WHERE
			id = $1 AND deleted = FALSE `, input.TableName)


	param := []interface{}{userParam.ID.Int64}
	_ = input.checkOwnPermission(userParam.CreatedBy.Int64, &query, &param, 2)

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt, &result.HostName, &result.HostURL)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerDAO) ViewDetailHostServer(db *sql.DB, userParam in.GetListDataDTO, parameter repository.HostServerModel, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(
		`SELECT
			cs.id, cs.cron, cs.name, sr.status
		FROM cron_host ch 
		INNER JOIN %s hs ON hs.id = ch.host_id AND hs.deleted = FALSE
		LEFT JOIN cron_scheduler cs ON cs.id = ch.cron_id AND cs.deleted = FALSE
		LEFT JOIN server_run sr ON sr.run_type = cs.run_type AND hs.id = sr.host_id AND sr.deleted = FALSE`,
		input.TableName)


	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{parameter.ID.Int64}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ListScheduler
			errorS := rows.Scan(&temp.ID, &temp.Cron, &temp.Name, &temp.Status)
			return temp, errorS
		}, " AND hs.id = $1 AND ch.status = TRUE ", input.getHostServerDefaultMustCheck(createdBy, false))
}

func (input hostServerDAO) CountRunningCron(db *sql.DB, userParam repository.HostServerModel, searchBy []in.SearchByParam) (result int, err errorModel.ErrorModel) {
	table :=
		"	cron_host ch INNER JOIN host_server " +
			"	hs ON hs.id = ch.host_id AND hs.deleted = FALSE " +
			"	LEFT JOIN cron_scheduler cs ON cs.id = ch.cron_id AND cs.deleted = FALSE " +
			"	LEFT JOIN server_run sr ON sr.run_type = cs.run_type AND hs.id = sr.host_id AND sr.deleted = FALSE "

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{userParam.ID.Int64}, table, searchBy, " AND hs.id = $1 AND ch.status = TRUE ",input.getHostServerDefaultMustCheck(0, false))
}

func (input hostServerDAO) GetHostUrlDbByHostname(db *sql.DB, userParam repository.HostServerModel) (result repository.HostServerModel, err errorModel.ErrorModel) {
	funcName := "GetHostServerDb"
	query := fmt.Sprintf(
		`SELECT
			id, host_url
		FROM
			%s
		WHERE
			host_name = $1 AND deleted = FALSE `,
			input.TableName)


	param := []interface{}{userParam.HostName.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.HostURL)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}