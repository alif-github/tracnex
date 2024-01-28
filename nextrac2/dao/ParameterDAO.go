package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type parameterDAO struct {
	AbstractDAO
}

var ParameterDAO = parameterDAO{}.New()

func (input parameterDAO) New() (output parameterDAO) {
	output.FileName = "Parameter.go"
	output.TableName = "parameter"
	return
}

func (input parameterDAO) GetParameterByNameAndCode(db *sql.DB, userParam repository.ParameterModel) (result repository.ParameterModel, err errorModel.ErrorModel) {
	funcName := "GetParameterByNameAndCode"
	query :=
		"SELECT " +
			"	id, permission, name, value " +
			"FROM " +
			"	parameter " +
			"WHERE " +
			"	permission = $1 and name = $2 "

	param := []interface{}{userParam.Permission.String, userParam.Name.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.Permission, &result.Name, &result.Value)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) GetUserParameter(db *sql.DB, userParam repository.UserParameterModel) (result repository.UserParameterModel, err errorModel.ErrorModel) {
	funcName := "GetUserParameter"
	query :=
		"SELECT " +
			"	id, user_id, parameter_value " +
			"FROM " +
			"	user_parameter " +
			"WHERE " +
			"	user_id = $1 "
	param := []interface{}{userParam.UserID.Int64}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.UserID, &result.ParameterValue)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) InsertParameter(db *sql.Tx, userParam repository.ParameterModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertParameter"
	query := "INSERT INTO parameter(permission, name, value, description, created_at, created_by, created_client, updated_at, updated_by, updated_client) VALUES " +
		"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	param := []interface{}{userParam.Permission.String, userParam.Name.String, userParam.Value.String, userParam.Description.String, userParam.CreatedAt.Time, userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.UpdatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String}

	results := db.QueryRow(query, param...)

	errs := results.Scan(&id)
	if errs != nil && errs.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) GetListParameter(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query :=
		"SELECT " +
			"   id, permission, name, value, description, created_by, updated_at " +
			" FROM " +
			"   parameter "

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ParameterModel
			errors := rows.Scan(
				&temp.ID, &temp.Permission, &temp.Name, &temp.Value, &temp.Description, &temp.CreatedBy, &temp.UpdatedAt)
			return temp, errors
		}, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input parameterDAO) GetCountParameter(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input parameterDAO) ViewParameter(db *sql.DB, userParam repository.ParameterModel) (result repository.ParameterModel, errors errorModel.ErrorModel) {
	funcName := "ViewParameter"

	query :=
		"SELECT " +
			"   id, permission, name, value, description, created_by, updated_at " +
			" FROM " +
			"   parameter " +
			" WHERE " +
			"   id = $1 AND deleted = FALSE "
	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 != 0 {
		query += "created_by = $2"
		param = append(param, userParam.CreatedBy.Int64)
	}

	results := db.QueryRow(query, param...)

	err := results.Scan(&result.ID, &result.Permission, &result.Name, &result.Value, &result.Description, &result.CreatedBy, &result.UpdatedAt)
	if err != nil && err.Error() != "sql: no rows in result set" {
		errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errors = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) UpdateParameter(db *sql.Tx, userParam repository.ParameterModel, timeNow time.Time) (errors errorModel.ErrorModel) {
	funcName := "UpdateParameter"

	query :=
		"UPDATE parameter SET permission = $1, name = $2, value = $3, description = $4, updated_at = $5, updated_by = $6, updated_client = $7 " +
			"WHERE id = $8 "
	param := []interface{}{userParam.Permission.String, userParam.Name.String, userParam.Value.String, userParam.Description.String, timeNow, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.ID.Int64}

	stmt, err := db.Prepare(query)
	if err != nil {
		errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	_, err = stmt.Exec(param...)
	if err != nil {
		errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errors = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) GetParameterForUpdate(db *sql.Tx, userParam repository.ParameterModel) (result repository.ParameterModel, errs errorModel.ErrorModel) {
	funcName := "GetParameterForUpdate"
	query :=
		"SELECT " +
			"   id, created_by, updated_at " +
			" FROM " +
			"   parameter " +
			" WHERE " +
			"   id = $1 AND deleted = FALSE "

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, param...)

	err := results.Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if err != nil && err.Error() != "sql: no rows in result set" {
		errs = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errs = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) DeleteParameter(db *sql.Tx, userParam repository.ParameterModel, timeNow time.Time) (errors errorModel.ErrorModel) {
	funcName := "DeleteParameter"

	query := "UPDATE parameter SET deleted = TRUE, updated_at = $1, updated_by = $2, updated_client = $3 WHERE id = $4 "
	param := []interface{}{timeNow, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.ID.Int64}

	stmt, err := db.Prepare(query)
	if err != nil {
		errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	_, err = stmt.Exec(param...)
	if err != nil {
		errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errors = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) GetParameterForEmployee(db *sql.DB) (results []repository.ParameterModel, err errorModel.ErrorModel) {
	funcName := "GetParameterForEmployee"

	query := "SELECT id, name, value, updated_at FROM parameter WHERE deleted=FALSE and name IN('cutOffAnualLeave','anualLeaveAfterProbation', 'expiredMedicalClaim') "

	rows, errorRows := db.Query(query)

	if errorRows != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
		return
	}

	if rows != nil {
		defer func() {
			errorRows = rows.Close()
			if errorRows != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
				return
			}
		}()

		for rows.Next() {
			var fc repository.ParameterModel
			errorRows = rows.Scan(&fc.ID, &fc.Name, &fc.Value, &fc.UpdatedAt)

			if errorRows != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
				return
			}
			results = append(results, fc)
		}

	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
		return
	}
	return
}

func (input parameterDAO) GetDetailUpdateParameter(db *sql.Tx, id int64) (parameter repository.ParameterModel, err errorModel.ErrorModel) {
	funcName := "GetDetailUpdateParameter"
	query := "SELECT id, updated_at " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND id = $1 "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&parameter.ID, &parameter.UpdatedAt)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterDAO) UpdateParameterEmployee(db *sql.Tx, parameter repository.ParameterModel) errorModel.ErrorModel {
	funcName := "UpdateParameterEmployee"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	value = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4 " +
		"WHERE " +
		"	id = $5 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		parameter.Value.String,
		parameter.UpdatedClient.String,
		parameter.UpdatedAt.Time,
		parameter.UpdatedBy.Int64,
		parameter.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}