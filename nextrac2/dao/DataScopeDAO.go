package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type dataScopeDAO struct {
	AbstractDAO
}

var DataScopeDAO = dataScopeDAO{}.New()

func (input dataScopeDAO) New() (output dataScopeDAO) {
	output.FileName = "DataScopeDAO.go"
	output.TableName = "data_scope"
	return
}

func (input dataScopeDAO) CheckIsScopeValid(db *sql.DB, dataGroupScope []string) (result int, err errorModel.ErrorModel) {
	funcName := "CheckIsScopeValid"
	listInterface := ArrayStringToArrayInterface(dataGroupScope)
	inQuery := ListDataToInQuery(listInterface)

	query :=
		"SELECT " +
			"	COUNT(id) " +
			"FROM " +
			"	data_scope " +
			"WHERE " +
			"	scope IN(" + inQuery + ") AND " +
			"	deleted = FALSE "

	results := db.QueryRow(query, listInterface...)

	errorS := results.Scan(&result)

	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataScopeDAO) ViewDataScope(db *sql.DB, userParam repository.DataScopeModel) (dataScope repository.DataScopeModel, err errorModel.ErrorModel) {
	funcName := "ViewDataScope"

	query := "SELECT " +
		"id, scope, description, created_by, updated_at " +
		"FROM " +
		input.TableName +
		" WHERE " +
		"id = $1 AND deleted = FALSE "
	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 != 0 {
		query += "AND created_by = $2"
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := db.QueryRow(query, param...).Scan(
		&dataScope.ID,
		&dataScope.Scope,
		&dataScope.Description,
		&dataScope.CreatedBy,
		&dataScope.UpdatedAt)

	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataScopeDAO) ListDataScopeFromDB(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query :=
		"SELECT " +
			"id, scope, description, created_by, updated_at " +
			"FROM " +
			"data_scope "

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.DataScopeModel
			errorS := rows.Scan(&temp.ID, &temp.Scope, &temp.Description, &temp.CreatedBy, &temp.UpdatedAt)
			return temp, errorS
		}, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input dataScopeDAO) ListDataScopeForDataGroup(db *sql.DB, userParam in.GetListDataDTO) (result []string, err errorModel.ErrorModel) {
	funcName := "ListDataScopeForDataGroup"
	var resultStr sql.NullString

	query := "SELECT json_agg(scope) FROM (SELECT scope FROM data_scope LIMIT $1 OFFSET $2) scope"

	errorS := db.QueryRow(query, userParam.Limit, CountOffset(userParam.Page, userParam.Limit)).Scan(&resultStr)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(resultStr.String), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataScopeDAO) GetDataScopeForDataGroup(db *sql.DB, userParam int64) (result repository.MapOfDataScopeForDataGroupModel, error errorModel.ErrorModel) {
	var jsonResult sql.NullString
	funcName := "GetDataScopeForDataGroup"
	query :=
		"	SELECT	" +
			"	scope as data_scope " +
			"FROM	" +
			"	data_group	" +
			"WHERE	" +
			"	id = $1 AND deleted = FALSE"
	results := db.QueryRow(query, userParam)
	err := results.Scan(&jsonResult)
	if err != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_ = json.Unmarshal([]byte(jsonResult.String), &result.DataGroupScope)

	return
}

func (input dataScopeDAO) UpdateDataScope(tx *sql.Tx, dataScope repository.DataScopeModel, timeNow time.Time) (err errorModel.ErrorModel) {
	funcName := "UpdateDataScope"

	query :=
		"UPDATE " + input.TableName +
			" SET " +
			"scope = $1, " +
			"description = $2, " +
			"updated_by = $3, " +
			"updated_at = $4, " +
			"updated_client = $5 " +
			"WHERE id = $6"
	param := []interface{}{
		dataScope.Scope.String,
		dataScope.Description.String,
		dataScope.UpdatedBy.Int64,
		timeNow,
		dataScope.UpdatedClient.String,
		dataScope.ID.Int64,
	}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
	}
	return
}

func (input dataScopeDAO) InsertDataScope(db *sql.Tx, userParam repository.DataScopeModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertDataScope"
	query := "INSERT INTO data_scope(scope, description, created_by, created_client, created_at, updated_by, updated_client, updated_at) VALUES " +
		"($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	param := []interface{}{userParam.Scope.String, userParam.Description.String, userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time}

	results := db.QueryRow(query, param...)

	errs := results.Scan(&id)
	if errs != nil && errs.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)

		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataScopeDAO) DeleteDataScope(tx *sql.Tx, dataScope repository.DataScopeModel, timeNow time.Time) (err errorModel.ErrorModel) {
	funcName := "DeleteDataScope"

	query :=
		"UPDATE " + input.TableName +
			" SET " +
			"updated_by = $1, " +
			"updated_at = $2, " +
			"updated_client = $3, " +
			"deleted = TRUE " +
			"WHERE id = $4"
	param := []interface{}{
		dataScope.UpdatedBy.Int64,
		timeNow,
		dataScope.UpdatedClient.String,
		dataScope.ID.Int64}

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

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataScopeDAO) GetDataScopeForUpdate(db *sql.Tx, userParam repository.DataScopeModel) (result repository.DataScopeModel, err errorModel.ErrorModel) {
	funcName := "GetDataScopeForUpdate"

	query := "SELECT " +
		"id, created_by, updated_at " +
		"FROM " + input.TableName +
		" WHERE " +
		"id = $1 AND deleted = FALSE "

	param := []interface{}{userParam.ID.Int64}
	if userParam.CreatedBy.Int64 > 0 {
		query += "AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += "FOR UPDATE"

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataScopeDAO) GetCountDataScope(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input dataScopeDAO) GetDataScopeByScope(db *sql.Tx, userParam repository.DataScopeModel) (result repository.DataScopeModel, err errorModel.ErrorModel) {
	funcName := "GetDataScopeByScope"

	query := fmt.Sprintf(`
		SELECT 
			id, created_by, updated_at 
		FROM %s
		WHERE
			scope = $1 AND deleted = FALSE
	`, input.TableName)

	query += "FOR UPDATE"

	param := []interface{}{userParam.Scope.String}
	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
