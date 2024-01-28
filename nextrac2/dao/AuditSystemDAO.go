package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
	"time"
)

type auditSystemDAO struct {
	AbstractDAO
}

var AuditSystemDAO = auditSystemDAO{}.New()

func (input auditSystemDAO) New() (output auditSystemDAO) {
	output.FileName = "AuditSystemDAO.go"
	output.TableName = "audit_system"
	return
}

func (input auditSystemDAO) GetDataForAuditByIDTx(db *sql.Tx, action int32, contextModel applicationModel.ContextModel, time time.Time, tableName string, id int64, createdBy int64) (result []repository.AuditSystemModel, error errorModel.ErrorModel) {
	auditField := make(map[string]repository.AuditSystemFieldParam)
	auditField["id"] = repository.AuditSystemFieldParam{
		IsEqual:    true,
		ParamValue: id,
	}
	if action != constanta.ActionAuditDeleteConstanta {
		auditField["deleted"] = repository.AuditSystemFieldParam{
			IsEqual:    true,
			ParamValue: false,
		}
	}
	if createdBy > 0 {
		auditField["created_by"] = repository.AuditSystemFieldParam{
			IsEqual:    true,
			ParamValue: createdBy,
		}
	}

	return input.GetDataForAuditTx(db, action, contextModel, time, tableName, auditField)
}

func (input auditSystemDAO) GetDataForAuditTx(db *sql.Tx, action int32, contextModel applicationModel.ContextModel, time time.Time, tableName string, auditField map[string]repository.AuditSystemFieldParam) (result []repository.AuditSystemModel, error errorModel.ErrorModel) {
	if !config.ApplicationConfiguration.GetAudit().IsActive {
		return
	}

	funcName := "GetDataForAuditTx"
	queryAdd, param := AuditSystemFieldParamToQuery(1, auditField)

	query := "SELECT a.id, uuid_key, row_to_json(a) FROM (SELECT * FROM \"" + tableName + "\" WHERE " + queryAdd + " FOR UPDATE)a "

	rows, err := db.Query(query, param...)
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
			var temp repository.AuditSystemModel
			err = rows.Scan(&temp.PrimaryKey, &temp.UUIDKey, &temp.Data)
			if err != nil {
				error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
				return
			}
			temp.CreatedAt.Time = time
			temp.Action.Int32 = action
			temp.TableName.String = tableName
			temp.CreatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
			temp.CreatedClient.String = contextModel.AuthAccessTokenModel.ClientID
			result = append(result, temp)
		}
	} else {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}

func (input auditSystemDAO) InsertAuditSystem(db *sql.DB, userParam repository.AuditSystemModel) errorModel.ErrorModel {
	var (
		funcName = "InsertAuditSystem"
		params   []interface{}
	)

	if !config.ApplicationConfiguration.GetAudit().IsActive {
		return errorModel.GenerateInactiveAuditSystem(input.FileName, funcName)
	}

	query := "INSERT INTO audit_system(table_name, uuid_key, primary_key, data, action, created_by, created_client, created_at, description) VALUES " +
		"($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	params = []interface{}{
		userParam.TableName.String, userParam.UUIDKey.String, userParam.PrimaryKey.Int64,
		userParam.Data.String, userParam.Action.Int32, userParam.CreatedBy.Int64,
		userParam.CreatedClient.String, userParam.CreatedAt.Time,
	}

	if userParam.Description.String != "" {
		params = append(params, userParam.Description.String)
	} else {
		params = append(params, nil)
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(params...)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input auditSystemDAO) GetDataForAuditByID(db *sql.DB, action int32, contextModel applicationModel.ContextModel, time time.Time, tableName string, id int64, createdBy int64) (result []repository.AuditSystemModel, error errorModel.ErrorModel) {
	auditField := make(map[string]repository.AuditSystemFieldParam)
	auditField["id"] = repository.AuditSystemFieldParam{
		IsEqual:    true,
		ParamValue: id,
	}
	if action != constanta.ActionAuditDeleteConstanta {
		auditField["deleted"] = repository.AuditSystemFieldParam{
			IsEqual:    true,
			ParamValue: false,
		}
	}
	if createdBy > 0 {
		auditField["created_by"] = repository.AuditSystemFieldParam{
			IsEqual:    true,
			ParamValue: createdBy,
		}
	}

	return input.GetDataForAudit(db, action, contextModel, time, tableName, auditField)
}

func (input auditSystemDAO) GetDataForAudit(db *sql.DB, action int32, contextModel applicationModel.ContextModel, time time.Time, tableName string, auditField map[string]repository.AuditSystemFieldParam) (result []repository.AuditSystemModel, error errorModel.ErrorModel) {
	if !config.ApplicationConfiguration.GetAudit().IsActive {
		return
	}

	funcName := "GetDataForAudit"
	queryAdd, param := AuditSystemFieldParamToQuery(1, auditField)

	query := "SELECT a.id, uuid_key, row_to_json(a) FROM (SELECT * FROM \"" + tableName + "\" WHERE " + queryAdd + ")a "

	rows, err := db.Query(query, param...)
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
			var temp repository.AuditSystemModel
			err = rows.Scan(&temp.PrimaryKey, &temp.UUIDKey, &temp.Data)
			if err != nil {
				error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
				return
			}
			temp.CreatedAt.Time = time
			temp.Action.Int32 = action
			temp.TableName.String = tableName
			temp.CreatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
			temp.CreatedClient.String = contextModel.AuthAccessTokenModel.ClientID
			result = append(result, temp)
		}
	} else {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	error = errorModel.GenerateNonErrorModel()
	return
}

func (input auditSystemDAO) GetListAuditData(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64,
	scope map[string]interface{}, mappingScope map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	index := 1
	query := fmt.Sprintf(
		`SELECT
		a.id, a.table_name, a.primary_key, a.action,
		a.created_by, a.created_client, a.created_at,
		CONCAT(u.first_name, ' ', u.last_name) AS created_name
	FROM %s a 
	LEFT JOIN "%s" AS u ON a.created_by = u.id `, input.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, searchByParam)
	tempQuery, tempParam := input.addScopeQueryForLog(scope, mappingScope, index, "a.table_name", "a.primary_key")

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, tempParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.AuditSystemModel
			errorS := rows.Scan(
				&temp.ID, &temp.TableName, &temp.PrimaryKey,
				&temp.Action, &temp.CreatedBy, &temp.CreatedClient,
				&temp.CreatedAt, &temp.CreatedName)
			return temp, errorS
		}, tempQuery, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "a.id"},
			Deleted:   FieldStatus{FieldName: "a.deleted"},
			Status:    FieldStatus{FieldName: "a.status", IsCheck: false},
			CreatedBy: FieldStatus{FieldName: "a.created_by", Value: createdBy},
		})
}

func (input auditSystemDAO) CountAuditData(db *sql.DB, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (int, errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
}

func (input auditSystemDAO) ViewAuditData(db *sql.DB, userParam repository.AuditSystemModel) (result repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "ViewAuditData"
	query :=
		"SELECT " +
			"	id, table_name, uuid_key, primary_key, " +
			"	data, action, created_by, created_client, created_at " +
			"FROM " +
			"	audit_system " +
			"WHERE " +
			"	id = $1"

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 != 0 {
		query += "AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.TableName, &result.UUIDKey, &result.PrimaryKey, &result.Data, &result.Action, &result.CreatedBy, &result.CreatedClient, &result.CreatedAt)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input auditSystemDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam []in.SearchByParam) {
	for i := 0; i < len(searchByParam); i++ {
		searchByParam[i].SearchKey = "a." + searchByParam[i].SearchKey
	}

	switch userParam.OrderBy {
	case "created_name", "created_name ASC", "created_name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "u.first_name " + strSplit[1]
		} else {
			userParam.OrderBy = "u.first_name"
		}
		break
	default:
		userParam.OrderBy = "a." + userParam.OrderBy
		break
	}
}

func (input auditSystemDAO) addScopeQueryForLog(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int, fieldTableName, fieldPrimaryName string) (query string, param []interface{}) {
	var listQueryWithScope = make(map[string]string)
	for key, value := range scopeLimit {
		listData := value.([]interface{})
		var tempQuery string
		for i := 0; i < len(listData); i++ {
			data := listData[i].(string)
			if data != constanta.ScopeAccessAll {
				if i == 0 {
					tempQuery += fieldPrimaryName + " IN ( "
				}
				id, _ := strconv.Atoi(data)
				if id != 0 {
					param = append(param, id)
					tempQuery += " $" + strconv.Itoa(idxStart) + ", "
					idxStart++
				}
			}
		}
		if tempQuery != "" {
			tempQuery = " ( " + fieldTableName + " = $" + strconv.Itoa(idxStart) + " AND " + tempQuery[0:len(tempQuery)-2]
			tempQuery += " ) ) "
			param = append(param, scopeDB[key].View)
			idxStart++
			listQueryWithScope[scopeDB[key].View] = tempQuery
		}
	}

	i := 0
	var tableScopeQuery []string
	for key, value := range listQueryWithScope {
		query += value
		if i < len(listQueryWithScope)-1 {
			query += " OR "
		}
		tableScopeQuery = append(tableScopeQuery, key)
		i++
	}

	if len(tableScopeQuery) > 0 {
		query += " OR " + fieldTableName + " IN ( SELECT DISTINCT(table_name) FROM " + input.TableName + " WHERE table_name NOT IN ( "
		for i := 0; i < len(tableScopeQuery); i++ {
			query += "$" + strconv.Itoa(idxStart)
			idxStart++
			param = append(param, tableScopeQuery[i])
			if i < len(tableScopeQuery)-1 {
				query += ", "
			}
		}
		query += " )) "
	}

	if query != "" {
		query = " AND ( " + query + " ) "
	}
	return
}

func (input auditSystemDAO) GetListEmployeeNotification(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, notification repository.EmployeeNotification) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT
				aud.id, aud.description, e.first_name, 
				e.last_name, aud.created_at
			FROM audit_system AS aud
			LEFT JOIN employee AS e
				ON (aud.description::json->>'employee_id')::bigint = e.id`

	addQuery, params := input.getEmployeeNotificationAddQuery(notification)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.AuditSystemModel

			err := rows.Scan(
				&model.ID, &model.Description, &model.Employee.Firstname,
				&model.Employee.Lastname, &model.CreatedAt)

			return model, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "aud.deleted",
			},
		})
}

func (input auditSystemDAO) getEmployeeNotificationAddQuery(model repository.EmployeeNotification) (result string, params []interface{}) {
	params = append(params, model.EmployeeId)

	memberQuery := ""
	if model.MemberIdList != nil {
		memberQuery = fmt.Sprintf(` OR (
						(
							(aud.description::json->>'is_requesting_for_approval')::bool = true or 	
							(aud.description::json->>'is_requesting_for_cancellation')::bool = true 
						) AND 
						e.id IN (%s)
					)`, strings.Join(model.MemberIdList, ","))
	}

	result = fmt.Sprintf(` 
		AND (aud.description::json->>'is_mobile_notification')::bool = true
		AND (
		(
			(
				(aud.description::json->>'is_requesting_for_approval')::bool = false AND
				(aud.description::json->>'is_requesting_for_cancellation')::bool = false
			) AND 
			e.id = $1  
		)
		%s
	)`, memberQuery)

	if model.FilterByIsRead {
		result += " AND (aud.description::json->>'is_read')::bool = $2 "
		params = append(params, model.IsRead)
	}

	return
}

func (input auditSystemDAO) UpdateDescriptionByIdTx(tx *sql.Tx, model repository.AuditSystemModel) errorModel.ErrorModel {
	funcName := "UpdateDescriptionByIdTx"

	query := `UPDATE ` + input.TableName + ` 
			SET 
				description = $1
			WHERE 
				id = $2 AND 
				deleted = FALSE`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	if _, err = stmt.Exec(model.Description.String, model.ID.Int64); err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}