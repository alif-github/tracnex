package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

type clientRegistrationLogDAO struct {
	AbstractDAO
}

var ClientRegistrationLogDAO = clientRegistrationLogDAO{}.New()

func (input clientRegistrationLogDAO) New() (output clientRegistrationLogDAO) {
	output.FileName = "ClientRegistrationLogDAO.go"
	output.TableName = "client_registration_log"
	return
}

func (input clientRegistrationLogDAO) InsertClientRegistrationLog(tx *sql.Tx, userParam *repository.ClientRegistrationLogModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertClientRegistrationLog"

	query := "INSERT INTO "+ input.TableName +" " +
		"(client_id, client_type_id, success_status_auth, " + // 1 , 2 , 3
		"message_auth, code, request_timestamp, " + // 4 , 5 , 6
		"created_by, created_at, created_client, " + // 7 , 8 , 9
		"updated_by, updated_at, updated_client, " + // 10 , 11 , 12
		"resource, details, success_status_nexcloud, " + // 13 , 14 , 15
		"success_status_nexdrive, message_nexcloud, message_nexdrive, " + // 16 , 17 , 18
		"attribute_request, request_count) " +  //19, 20
		"VALUES " +
		"($1, $2, $3, " +
		"$4, $5, $6, " +
		"$7, $8, $9, " +
		"$10, $11, $12, " +
		"$13, $14, $15, " +
		"$16, $17, $18, " +
		"$19, $20) returning id"

	var param []interface{}

	param = append(param, userParam.ClientID.String, userParam.ClientTypeID.Int64, userParam.SuccessStatusAuth.Bool,
		userParam.MessageAuth.String, userParam.Code.String, userParam.RequestTimeStamp.Time,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.Resource.String)

	if userParam.Details.String != "" {
		param = append(param, userParam.Details.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.SuccessStatusNexcloud.Bool, userParam.SuccessStatusNexdrive.Bool)

	if userParam.MessageNexcloud.String != "" {
		param = append(param, userParam.MessageNexcloud.String)
	} else {
		param = append(param, nil)
	}

	if userParam.MessageNexdrive.String != "" {
		param = append(param, userParam.MessageNexdrive.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.AttributeRequest.String, userParam.RequestCount.Int64)

	errorS := tx.QueryRow(query, param...).Scan(&id)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogDAO) GetClientRegistrationLogForScheduler(db *sql.DB, userParam repository.ParamClientRegistrationLogModel) (result []repository.ClientRegistrationLogModel, err errorModel.ErrorModel){
	funcName := "GetClientRegistrationLogForScheduler"

	query := "select id, client_id, resource " +
		"from "+ input.TableName +" " +
		"where " +
		"(success_status_auth = true AND " +
		"success_status_nexcloud = false) AND " +
		"client_type_id = $1"

	param := []interface{}{userParam.ClientTypeID.Int64}

	rows, errorS := db.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if rows != nil {
		defer func() {
			errorS = rows.Close()
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
		}()
		for rows.Next() {
			var temp repository.ClientRegistrationLogModel
			errorS = rows.Scan(&temp.ID, &temp.ClientID, &temp.Resource)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			result = append(result, temp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogDAO) GetDataClientRegistrationLogForUpdate(db *sql.DB,
	userParam repository.ClientRegistrationLogModel) (result repository.ClientRegistrationLogModel, err errorModel.ErrorModel) {

	fileName := "ClientRegistrationLogDAO.go"
	funcName := "GetDataClientRegistrationLogForUpdate"

	query := "SELECT " +
		"id, resource, client_type_id, " +
		"updated_at " +
		"FROM "+ input.TableName +" " +
		"WHERE " +
		"client_id = $1 AND deleted = false"

	param := []interface{}{userParam.ClientID.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID.Int64, &result.Resource.String, &result.ClientTypeID.Int64,
		&result.UpdatedAt.Time)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogDAO) UpdateRegistrationLogForAddResource(tx *sql.Tx, userParam repository.ClientRegistrationLogModel, doForResource string) (err errorModel.ErrorModel) {
	funcName := "UpdateRegistrationLogForAddResource"
	var daoKeyValueTemp in.PreparedDaoUpdateClientRegisterLog

	switch doForResource {
	case constanta.NexCloudResourceID :
		daoKeyValueTemp = in.PreparedDaoUpdateClientRegisterLog {
			StatusKey: 	"success_status_nexcloud",
			MessageKey:	"message_nexcloud",
		}
		break
	case constanta.NexdriveResourceID :
		daoKeyValueTemp = in.PreparedDaoUpdateClientRegisterLog {
			StatusKey: 	"success_status_nexdrive",
			MessageKey:	"message_nexdrive",
		}
	}

	daoKeyValueTemp.StatusValue = userParam.SuccessStatus.Bool
	daoKeyValueTemp.MessageValue = userParam.Message.String

	query := "UPDATE "+ input.TableName +" " +
		"SET " +
		""+ daoKeyValueTemp.StatusKey +" = $1, "+ daoKeyValueTemp.MessageKey +" = $2, details = $3, " +
		"code = $4, request_timestamp = $5, updated_by = $6, " +
		"updated_client = $7, updated_at = $8, resource = $9, " +
		"attribute_request = $10 " +
		"WHERE " +
		"client_id = $11 AND " +
		"deleted = false"

	var param []interface{}

	param = append(param, daoKeyValueTemp.StatusValue, daoKeyValueTemp.MessageValue)

	if userParam.Details.String != "" {
		param = append(param, userParam.Details.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Code.String != "" {
		param = append(param, userParam.Code.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.RequestTimeStamp.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.Resource.String, userParam.AttributeRequest.String,
		userParam.ClientID.String)

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

func (input clientRegistrationLogDAO) UpdateRegistrationLogForAddResourceNexcloudScheduler(tx *sql.Tx, userParam repository.ClientRegistrationLogModel) (err errorModel.ErrorModel) {
	funcName := "UpdateRegistrationLogForAddResourceNexcloudScheduler"

	query := "UPDATE "+ input.TableName +" " +
		"SET " +
		"success_status_nexcloud = $1, message_nexcloud = $2, details = $3, " +
		"code = $4, request_timestamp = $5, updated_by = $6, " +
		"updated_client = $7, updated_at = $8, resource = $9, " +
		"attribute_request = $10 " +
		"WHERE " +
		"id = $11 AND " +
		"deleted = false"

	var param []interface{}

	param = append(param, userParam.SuccessStatusNexcloud.Bool, userParam.MessageNexcloud.String)

	if userParam.Details.String != "" {
		param = append(param, userParam.Details.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Code.String != "" {
		param = append(param, userParam.Code.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.RequestTimeStamp.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.Resource.String, userParam.AttributeRequest.String,
		userParam.ID.Int64)

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

func (input pkceClientMappingDAO) ViewDetailLogRegistrationJoinUser(db *sql.DB, logModel repository.ViewClientRegistrationLogModel) (result repository.ViewClientRegistrationLogModel, err errorModel.ErrorModel) {
	funcName := "ViewDetailClientMappingPKCEJoinUser"

	query := "SELECT " +
		"logs.client_id, userTable.auth_user_id, ct.client_type, " + 			//1, 2, 3
		"userTable.status, userTable.first_name, userTable.last_name, " + 		//4, 5, 6
		"logs.resource, logs.updated_at " + 									//7, 8
		"FROM " +
		"client_registration_log logs " +
		"INNER JOIN " +
		"client_type ct ON logs.client_type_id = ct.id " +
		"INNER JOIN " +
		"\"user\" userTable ON logs.client_id = userTable.client_id " +
		"WHERE " +
		"logs.client_id = $1 AND " +
		"userTable.deleted = FALSE "

	params := []interface{}{logModel.ClientID.String}

	query += " FOR UPDATE"
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.ClientID, &result.AuthUserID, &result.ClientType,
		&result.Status, &result.FirstName, &result.LastName,
		&result.Resource, &result.UpdatedAt)
	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogDAO) GetListLog(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {

	query :=
		" SELECT " +
			"	id, client_id, client_type_id, " +
			"	success_status_auth, success_status_nexcloud, resource " +
			" FROM "+ input.TableName +" "

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "created_by",
			SearchOperator: "eq",
			SearchValue:    strconv.Itoa(int(createdBy)),
		})
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ListClientRegistrationLogModel
			errorS := rows.Scan(
				&temp.ID, &temp.ClientID, &temp.ClientTypeID,
				&temp.SuccessStatusAuth, &temp.SuccessStatusNexcloud, &temp.Resource)
			return temp, errorS
		}, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input clientRegistrationLogDAO) GetDataStatusResource(db *sql.DB, userParam repository.ClientRegistrationLogModel) (result repository.ClientRegistrationLogModel, err errorModel.ErrorModel) {

	fileName := "ClientRegistrationLogDAO.go"
	funcName := "GetDataStatusResource"

	query := "SELECT " +
		"id, success_status_auth, success_status_nexcloud, request_count " +
		"FROM "+ input.TableName +" " +
		"WHERE " +
		"client_id = $1 AND deleted = false"

	param := []interface{}{userParam.ClientID.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.SuccessStatusAuth, &result.SuccessStatusNexcloud, &result.RequestCount)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogDAO) UpdateLogAfterReRegistration(tx *sql.Tx, userParam repository.ClientRegistrationLogModel) (err errorModel.ErrorModel) {
	funcName := "UpdateLogAfterReRegistration"

	query := "UPDATE "+ input.TableName +" " +
		"SET " +
		"attribute_request = $1, message_auth = $2, code = $3, " +
		"request_timestamp = $4, updated_by = $5, updated_client = $6, " +
		"updated_at = $7, request_count = $8 " +
		"WHERE " +
		"id = $9 AND " +
		"deleted = false"

	var param []interface{}

	param = append(param, userParam.AttributeRequest.String, userParam.MessageAuth.String, userParam.Code.String,
		userParam.RequestTimeStamp.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.RequestCount.Int64, userParam.ID.Int64)

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