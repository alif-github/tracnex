package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type jobProcessDAO struct {
	AbstractDAO
}

var JobProcessDAO = jobProcessDAO{}.New()

func (input jobProcessDAO) New() (output jobProcessDAO) {
	output.FileName = "JobProcessDAO.go"
	output.TableName = "job_process"
	return
}

func (input jobProcessDAO) InsertJobProcess(tx *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	var (
		funcName = "InsertJobProcess"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`INSERT INTO %s 
		(parent_job_id, level, job_id, 
		"group", parameter, type, 
		name, counter, total, 
		created_by, created_at, created_client, 
		updated_at) 
		VALUES 
		($1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12, 
		$13) `,
		input.TableName)

	if userParam.ParentJobID.String != "" {
		param = append(param, userParam.ParentJobID.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Level.Int32 > 0 {
		param = append(param, userParam.Level.Int32)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.JobID.String, userParam.Group.String)

	if userParam.Parameter.String != "" {
		param = append(param, userParam.Parameter.String)
	} else {
		param = append(param, nil)
	}

	param = append(param,
		userParam.Type.String, userParam.Name.String, userParam.Counter.Int32,
		userParam.Total.Int32, userParam.CreatedBy.Int64, userParam.CreatedAt.Time,
		userParam.CreatedClient.String, userParam.UpdatedAt.Time)

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateJobProcessCounter(db *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessCounter"

	if userParam.Status.String == constanta.JobProcessOnProgressStatus || userParam.Status.String == constanta.JobProcessOnProgressErrorStatus {
		if userParam.Counter.Int32 == userParam.Total.Int32 {
			if userParam.Status.String == constanta.JobProcessOnProgressStatus {
				userParam.Status.String = constanta.JobProcessDoneStatus
			} else {
				userParam.Status.String = constanta.JobProcessErrorStatus
			}
		}
	}

	query := "UPDATE " + input.TableName + " " +
		"SET " +
		"counter = $1, status = $2, updated_at = $3, " +
		"content_data_out = $4, url_in = $5, filename_in = $6 " +
		"WHERE " +
		"job_id = $7 "

	param := []interface{}{
		userParam.Counter.Int32, userParam.Status.String, userParam.UpdatedAt.Time,
		userParam.ContentDataOut.String, userParam.UrlIn.String, userParam.FilenameIn.String,
		userParam.JobID.String}

	stmt, errorS := db.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateErrorJobProcess(tx *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessCounter"

	query := "UPDATE " +
		"" + input.TableName + " " +
		"SET " +
		"status = $1, " +
		"updated_at = $2, " +
		"content_data_out = $3 " +
		"WHERE " +
		"job_id = $4 "

	param := []interface{}{
		userParam.Status.String, userParam.UpdatedAt.Time, userParam.ContentDataOut.String,
		userParam.JobID.String}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateErrorParentJobProcess(db *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateErrorParentJobProcess"
	var parent repository.JobProcessModel
	tx, errs := db.Begin()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	parent.JobID = userParam.ParentJobID
	parent, err = input.GetJobProcessForUpdate(tx, parent)
	if err.Error != nil {
		return
	}

	parent.Status.String = constanta.JobProcessOnProgressErrorStatus
	parent.Counter.Int32 += 1
	parent.UpdatedAt.Time = userParam.UpdatedAt.Time

	err = input.UpdateJobProcessCounterTx(tx, parent)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateJobProcessCounterTx(tx *sql.Tx, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessCounterTx"

	if userParam.Status.String == constanta.JobProcessOnProgressStatus || userParam.Status.String == constanta.JobProcessOnProgressErrorStatus {
		if userParam.Counter.Int32 == userParam.Total.Int32 {
			if userParam.Status.String == constanta.JobProcessOnProgressStatus {
				userParam.Status.String = constanta.JobProcessDoneStatus
			} else {
				userParam.Status.String = constanta.JobProcessErrorStatus
			}
		}
	}

	query := "UPDATE " + input.TableName + " SET counter = $1, status = $2, updated_at = $3 WHERE job_id = $4 "

	param := []interface{}{userParam.Counter.Int32, userParam.Status.String, userParam.UpdatedAt.Time, userParam.JobID.String}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) UpdateParentJobProcessCounter(db *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateParentJobProcessCounter"
		parent   repository.JobProcessModel
	)

	tx, errs := db.Begin()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	parent.JobID = userParam.ParentJobID
	parent, err = input.GetJobProcessForUpdate(tx, parent)
	if err.Error != nil {
		return
	}

	parent.Counter.Int32 += 1
	parent.UpdatedAt.Time = userParam.UpdatedAt.Time

	err = input.UpdateJobProcessCounterTx(tx, parent)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) GetJobProcessForUpdate(db *sql.Tx, useParam repository.JobProcessModel) (result repository.JobProcessModel, err errorModel.ErrorModel) {
	funcName := "GetJobProcessForUpdate"
	query :=
		" SELECT " +
			"	job_id, counter, total, status " +
			" FROM " +
			input.TableName +
			" WHERE " +
			"	job_id = $1 FOR UPDATE "

	param := []interface{}{useParam.JobID.String}

	errorS := db.QueryRow(query, param...).Scan(&result.JobID, &result.Counter, &result.Total, &result.Status)

	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input jobProcessDAO) UpdateJobProcessUpdateAt(tx *sql.Tx, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessUpdateAt"

	query := "UPDATE " + input.TableName + " SET updated_at = $1 WHERE job_id = $2 "

	param := []interface{}{userParam.UpdatedAt.Time, userParam.JobID.String}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) GetListJobProcess(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query :=
		"SELECT " +
			"	level, job_id, \"group\", type, name, " +
			"	counter, total, status, created_at, updated_at " +
			"FROM " +
			"	job_process "

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ListJobProcessModel
			errorS := rows.Scan(&temp.Level, &temp.JobID, &temp.Group, &temp.Type, &temp.Name, &temp.Counter, &temp.Total, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt)
			return temp, errorS
		}, " AND (parent_job_id is null OR parent_job_id = '') ", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
}

func (input jobProcessDAO) GetCountJobProcess(db *sql.DB, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (int, errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, " AND (parent_job_id is null OR parent_job_id = '') ", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
}

func (input jobProcessDAO) ViewJobProcess(db *sql.DB, userParam repository.JobProcessModel) (result repository.ViewJobProcessModel, err errorModel.ErrorModel) {
	funcName := "ViewJobProcess"
	var child sql.NullString
	query :=
		"SELECT " +
			"	job.parent_job_id, job.level, job.job_id, job.\"group\", " +
			"	job.type, job.name, job.counter, job.total, job.status, " +
			"	job.created_at, job.updated_at, job.url_in, " +
			"	job.filename_in, job.content_data_out, job.id, " +
			"	(SELECT " +
			"		json_agg(a) " +
			"	FROM " +
			"		(SELECT " +
			"			level, job_id, \"group\", type, name, " +
			"			counter, total, status, created_at, updated_at, " +
			"			url_in, filename_in, content_data_out " +
			"		FROM " +
			"			job_process " +
			"		WHERE " +
			"			parent_job_id = job.job_id " +
			"		ORDER BY id) a) " +
			"FROM " +
			"	job_process job " +
			"WHERE " +
			"	job.job_id = $1"

	param := []interface{}{userParam.JobID.String}

	if userParam.CreatedBy.Int64 != 0 {
		query += " AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := db.QueryRow(query, param...).Scan(
		&result.ParentJobID, &result.Level, &result.JobID, &result.Group,
		&result.Type, &result.Name, &result.Counter, &result.Total, &result.Status,
		&result.CreatedAt, &result.UpdatedAt, &result.UrlIn,
		&result.FileNameIn, &result.ContentDataOut, &result.ID, &child)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(child.String), &result.ChildJobProcess)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input jobProcessDAO) UpdateErrorJobProcessWithCounter(db *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateJobProcessCounter"

	if userParam.Status.String == constanta.JobProcessOnProgressStatus || userParam.Status.String == constanta.JobProcessOnProgressErrorStatus {
		if userParam.Counter.Int32 == userParam.Total.Int32 {
			if userParam.Status.String == constanta.JobProcessOnProgressStatus {
				userParam.Status.String = constanta.JobProcessDoneStatus
			} else {
				userParam.Status.String = constanta.JobProcessErrorStatus
			}
		} else {
			userParam.Status.String = constanta.JobProcessErrorStatus
		}
	}

	query := "UPDATE " + input.TableName + " SET counter = $1, status = $2, updated_at = $3 WHERE job_id = $4 "

	param := []interface{}{userParam.Counter.Int32, userParam.Status.String, userParam.UpdatedAt.Time, userParam.JobID.String}

	stmt, errorS := db.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) GetJobProcessError(db *sql.DB, userParam repository.JobProcessModel) (result repository.JobProcessModel, err errorModel.ErrorModel) {
	funcName := "GetJobProcessError"
	query := fmt.Sprintf(`SELECT 
		id, parent_job_id, level, 
		job_id, "group", type, 
		name, counter, total, 
		status, content_data_out
	FROM %s 
	WHERE 
		"group" = $1 AND type = $2 AND
		name = $3 AND status = $4 AND
		counter <> total AND deleted = FALSE AND 
		created_at > now()::date AND created_at < now()::date + INTERVAL '1 DAY'`, input.TableName)

	param := []interface{}{
		userParam.Group.String,
		userParam.Type.String,
		userParam.Name.String,
		userParam.Status.String,
	}

	row := db.QueryRow(query, param...)

	tempData, err := RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.JobProcessModel
		dbErrorS := rws.Scan(
			&temp.ID, &temp.ParentJobID, &temp.Level,
			&temp.JobID, &temp.Group, &temp.Type,
			&temp.Name, &temp.Counter, &temp.Total,
			&temp.Status, &temp.ContentDataOut)
		return temp, dbErrorS
	}, input.FileName, funcName)

	if err.Error != nil {
		return
	}

	result = tempData.(repository.JobProcessModel)

	return
}

func (input jobProcessDAO) UpdateFullJobProcess(tx *sql.DB, userParam repository.JobProcessModel) (err errorModel.ErrorModel) {
	funcName := "UpdateFullJobProcess"

	query := "UPDATE " +
		"" + input.TableName + " " +
		"SET " +
		"status = $1, updated_at = $2, content_data_out = $3, " +
		"counter = $4, total = $5, message_alert = $6 " +
		"WHERE " +
		"job_id = $7 "

	param := []interface{}{userParam.Status.String, userParam.UpdatedAt.Time}

	if userParam.ContentDataOut.String != "null" {
		param = append(param, userParam.ContentDataOut.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.Counter.Int32, userParam.Total.Int32)

	if userParam.MessageAlert.String != "" {
		param = append(param, userParam.MessageAlert.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.JobID.String)

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input jobProcessDAO) GetJobProcessRunning(db *sql.DB, userParam repository.JobProcessModel) (result repository.JobProcessModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetJobProcessRunning"
		query    string
		param    []interface{}
		errorS   error
	)

	query = fmt.Sprintf(`
		SELECT id, job_id, name, 
		status, created_at 
		FROM %s 
		WHERE
		"group" = $1 AND type = $2 AND name = $3 AND  
		deleted = false 
		ORDER BY created_at DESC LIMIT 1 `,
		input.TableName)

	param = []interface{}{userParam.Group.String, userParam.Type.String, userParam.Name.String}
	errorS = db.QueryRow(query, param...).Scan(
		&result.ID, &result.JobID, &result.Name,
		&result.Status, &result.CreatedAt)
	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
