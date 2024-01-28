package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type clientTypeDAO struct {
	AbstractDAO
}

var ClientTypeDAO = clientTypeDAO{}.New()

func (input clientTypeDAO) New() (output clientTypeDAO) {
	output.FileName = "ClientTypeDAO.go"
	output.TableName = "client_type"
	return
}

func (input clientTypeDAO) CheckClientType(db *sql.DB, userParam *repository.ClientTypeModel) (result repository.ClientTypeModel, err errorModel.ErrorModel) {
	funcName := "CheckClientType"
	query := "SELECT " +
		"id " +
		"FROM " + input.TableName + " " +
		"WHERE " +
		"client_type = $1"

	results := db.QueryRow(query, userParam.ClientType.String)
	errorS := results.Scan(&result.ID)

	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) DeleteClientType(db *sql.Tx, userParam repository.ClientTypeModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteClientType"
	)

	query := fmt.Sprintf(`UPDATE %s SET 
									deleted = $1, 
									updated_by = $2, 
									updated_at = $3, 
									updated_client = $4 
								WHERE id = $5 `,
		input.TableName)

	param := []interface{}{
		true,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input clientTypeDAO) ViewClientType(db *sql.DB, userParam repository.ClientTypeModel) (result repository.ClientTypeModel, err errorModel.ErrorModel) {
	funcName := "ViewClientType"

	query := fmt.Sprintf(`
			SELECT
				ct.id, ct.client_type, ct.description, 
				ct.created_at, u.nt_username AS updated_name, ct.updated_at
			FROM %s ct
			LEFT JOIN "%s" u ON ct.updated_by = u.id 
			WHERE ct.id = $1 AND ct.deleted = FALSE `,
		input.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(
		&result.ID, &result.ClientType, &result.Description,
		&result.CreatedAt, &result.UpdatedName, &result.UpdatedAt,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input clientTypeDAO) GetClientTypeByID(db *sql.DB, userParam repository.ClientTypeModel) (result repository.ClientTypeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetClientTypeByID"
	)

	query := fmt.Sprintf(`SELECT id, client_type, created_by, updated_at FROM %s WHERE id = $1 AND deleted = FALSE `, input.TableName)

	results := db.QueryRow(query, userParam.ID.Int64)
	errorS := results.Scan(&result.ID, &result.ClientType, &result.CreatedBy, &result.UpdatedAt)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) GetInitiateInsertClientService(db *sql.DB) (result []out.ClientTypeResponse, err errorModel.ErrorModel) {
	funcName := "GetInitiateInsertClientService"
	query := "SELECT " +
		"client_type_id, client_type " +
		"FROM " + input.TableName + ""

	results, errorS := db.Query(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	errorS = results.Scan(&result)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) ValidationClientType(db *sql.DB, clientType repository.ClientTypeModel) (result repository.ClientTypeModel, err errorModel.ErrorModel) {
	funcName := "ValidationClientType"
	query := "SELECT " +
		"	id, client_type " +
		"FROM " + input.TableName + " " +
		"WHERE " +
		"	id = $1 AND " +
		"	client_type = $2 "

	params := []interface{}{clientType.ID.Int64, clientType.ClientType.String}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&result.ID, &result.ClientType)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) CheckClientTypeByID(db *sql.DB, userParam *repository.ClientTypeModel) (result repository.ClientTypeModel, err errorModel.ErrorModel) {
	var (
		funcName = "CheckClientTypeByID"
		query    string
		results  *sql.Row
		errorS   error
	)

	query = fmt.Sprintf(`SELECT id, client_type FROM %s WHERE id = $1`, input.TableName)

	results = db.QueryRow(query, userParam.ID.Int64)
	errorS = results.Scan(&result.ID, &result.ClientType)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) ValidateClientTypeByID(db *sql.DB, userParam repository.ClientTypeModel) (result repository.ClientTypeModel, err errorModel.ErrorModel) {
	funcName := "ValidateClientTypeByID"
	query := "SELECT " +
		"id, client_type " +
		"FROM " + input.TableName + " " +
		"WHERE " +
		"id = $1"

	results := db.QueryRow(query, userParam.ID.Int64)
	errorS := results.Scan(&result.ID, &result.ClientType)

	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) GetListClientType(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	var dbParam []interface{}
	query :=
		fmt.Sprintf(`
			SELECT 
				ct.id, ct.client_type, ct.description, 
				ct.created_at, u.nt_username AS updated_name, ct.updated_at, 
				ct.parent_client_type_id
			FROM %s ct 
			LEFT JOIN "%s" u ON ct.updated_by = u.id`,
			input.TableName, UserDAO.TableName)

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.ClientTypeDataScope})
	dbParam = append(dbParam, param...)

	input.convertUserParamAndSearchBy(&userParam, searchBy)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ClientTypeModel
			errors := rows.Scan(&temp.ID, &temp.ClientType, &temp.Description,
				&temp.CreatedAt, &temp.UpdatedName, &temp.UpdatedAt,
				&temp.ParentClientTypeID)
			return temp, errors

		}, additionalWhere, input.getDefaultMustCheck(createdBy))
}

func (input clientTypeDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "ct.id"},
		Deleted:   FieldStatus{FieldName: "ct.deleted"},
		CreatedBy: FieldStatus{FieldName: "ct.created_by", Value: createdBy},
	}
}

func (input clientTypeDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam []in.SearchByParam) {
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "client_type_id" {
			searchByParam[i].SearchKey = "ct.id"
		} else {
			searchByParam[i].SearchKey = "ct." + searchByParam[i].SearchKey
		}
	}

	switch userParam.OrderBy {
	case "updated_name", "updated_name ASC", "updated_name DESC":
		userParam.OrderBy = "u.nt_username"
		break
	default:
		userParam.OrderBy = "ct." + userParam.OrderBy
		break
	}
}

func (input clientTypeDAO) InsertClientType(db *sql.Tx, userParam repository.ClientTypeModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertClientType"
	)

	query := fmt.Sprintf(`
		INSERT INTO %s (
			client_type, description, created_by, 
			created_at, created_client, updated_by, 
			updated_at, updated_client
		`, input.TableName)

	if userParam.ParentClientTypeID.Int64 != 0 {
		query += fmt.Sprintf(`, parent_client_type_id`)
	}

	query += fmt.Sprintf(`) 
		VALUES (
			$1, $2, $3, 
			$4, $5, $6, 
			$7, $8`)

	if userParam.ParentClientTypeID.Int64 != 0 {
		query += fmt.Sprintf(`, $9`)
	}

	query += fmt.Sprintf(`) RETURNING id`)

	params := []interface{}{
		userParam.ClientType.String,
		userParam.Description.String,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
	}

	if userParam.ParentClientTypeID.Int64 != 0 {
		params = append(params, userParam.ParentClientTypeID.Int64)
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input clientTypeDAO) UpdateClientType(db *sql.Tx, userParam repository.ClientTypeModel) (err errorModel.ErrorModel) {
	var (
		funcName   = "UpdateClientType"
		startIndex = 6
	)

	query := fmt.Sprintf(`UPDATE %s SET
		client_type = $1, description = $2, updated_at = $3, 
		updated_by = $4, updated_client = $5 `,
		input.TableName)

	params := []interface{}{
		userParam.ClientType.String, userParam.Description.String, userParam.UpdatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
	}

	if userParam.ParentClientTypeID.Int64 > 0 {
		query += fmt.Sprintf(`, parent_client_type_id = $%d `, startIndex)
		params = append(params, userParam.ParentClientTypeID.Int64)
		startIndex++
	}

	query += fmt.Sprintf(`WHERE id = $%d`, startIndex)
	params = append(params, userParam.ID.Int64)

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input clientTypeDAO) CheckClientTypeByIDWithScope(db *sql.DB, userParam *repository.ClientTypeModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ClientTypeModel, err errorModel.ErrorModel) {
	var (
		funcName = "CheckClientTypeByIDWithScope"
		query    string
		results  *sql.Row
		errorS   error
	)

	query = fmt.Sprintf(`SELECT ct.id, ct.client_type FROM %s ct WHERE ct.id = $1 `, input.TableName)
	param := []interface{}{userParam.ID.Int64}

	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.ClientTypeDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
	}

	results = db.QueryRow(query, param...)
	errorS = results.Scan(&result.ID, &result.ClientType)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) CheckClientTypeIsParentAndExist(db *sql.DB, userParam *repository.ClientTypeModel) (id int64, isParent bool, err errorModel.ErrorModel) {
	var (
		funcName = "CheckClientTypeIsParentAndExist"
		query    string
		results  *sql.Row
		errorS   error
	)

	query = fmt.Sprintf(`SELECT id, CASE WHEN parent_client_type_id = 0 or parent_client_type_id is null
		THEN true ELSE false END is_parent 
		FROM %s 
		WHERE id = $1`, input.TableName)

	results = db.QueryRow(query, userParam.ID.Int64)
	errorS = results.Scan(&id, &isParent)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeDAO) GetCountClientType(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName+" ct ", searchByParam, "", input.getDefaultMustCheck(createdBy))
}

func (input clientTypeDAO) IsClientTypeUsed(db *sql.DB, userParam repository.ClientTypeModel) (output bool, err errorModel.ErrorModel) {
	funcName := "IsClientTypeUsed"
	var tempData repository.ClientTypeModel
	query := fmt.Sprintf(
		`SELECT ct.id
	FROM %s ct 
	WHERE 
		ct.id = $1 AND (EXISTS (SELECT 1 from %s cm where ct.id = cm.client_type_id) 
		OR EXISTS (SELECT 1 from %s lc where ct.id = lc.client_type_id)
		OR EXISTS (SELECT 1 from %s np where ct.id = np.client_type_id)
		OR EXISTS (SELECT 1 from %s pcm where ct.id = pcm.client_type_id)
		OR EXISTS (SELECT 1 from %s p where ct.id = p.client_type_id))`,
		input.TableName, ClientMappingDAO.TableName, LicenseConfigDAO.TableName,
		NexmileParameterDAO.TableName, PKCEClientMappingDAO.TableName, ProductDAO.TableName)

	results := db.QueryRow(query, userParam.ID.Int64)
	errorS := results.Scan(&tempData.ID)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if tempData.ID.Int64 > 0 {
		output = true
	}

	return
}
