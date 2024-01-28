package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type pkceClientMappingDAO struct {
	AbstractDAO
}

var PKCEClientMappingDAO = pkceClientMappingDAO{}.New()

func (input pkceClientMappingDAO) New() (output pkceClientMappingDAO) {
	output.FileName = "PKCEClientMappingDAO.go"
	output.TableName = "pkce_client_mapping"
	return
}

func (input pkceClientMappingDAO) InsertPKCEClientMapping(tx *sql.Tx, userParam *repository.PKCEClientMappingModel, isClientDependant bool) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertClientMapping"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`INSERT INTO %s 
		(client_id, company_id, branch_id,
		client_type_id, created_by, created_client,
		created_at, updated_by, updated_client,
		updated_at, parent_client_id, auth_user_id,
		username, installation_id, customer_id,
		site_id, client_alias, is_client_dependant)
		VALUES 
		($1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10, $11, $12, 
		$13, $14, $15, $16, $17, $18) 
		returning id `,
		input.TableName)

	param = append(param, userParam.ClientID.String)

	HandleOptionalParam([]interface{}{
		userParam.CompanyID.String,
		userParam.BranchID.String,
	}, &param)

	param = append(param, userParam.ClientTypeID.Int64, userParam.CreatedBy.Int64, userParam.CreatedClient.String,
		userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time)

	HandleOptionalParam([]interface{}{
		userParam.ParentClientID.String,
	}, &param)

	param = append(param, userParam.AuthUserID.Int64)

	HandleOptionalParam([]interface{}{
		userParam.Username.String,
		userParam.InstallationID.Int64,
		userParam.CustomerID.Int64,
		userParam.SiteID.Int64,
		userParam.ClientAlias.String,
	}, &param)

	if isClientDependant {
		userParam.IsClientDependant.String = "Y"
	} else {
		userParam.IsClientDependant.String = "N"
	}

	param = append(param, userParam.IsClientDependant.String)

	errorS := tx.QueryRow(query, param...).Scan(&id)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) CheckPKCEClientMapping(tx *sql.Tx, userParam repository.PKCEClientMappingModel) (result repository.PKCEClientMappingModel, err errorModel.ErrorModel) {
	funcName := "CheckPKCEClientMapping"
	paramCount := 2

	query := fmt.Sprintf(
		`SELECT 
			id, auth_user_id, client_id
		FROM %s
		WHERE 
			client_type_id = $1 AND
			username = $2 `, input.TableName)

	var param []interface{}
	param = append(param, userParam.ClientTypeID.Int64, userParam.Username.String)

	if userParam.ParentClientID.String != "" {
		paramCount++
		query += " AND parent_client_id = $" + strconv.Itoa(paramCount) + " "
		param = append(param, userParam.ParentClientID.String)
	}

	if userParam.CompanyID.String != "" {
		paramCount++
		query += " AND company_id = $" + strconv.Itoa(paramCount) + " "
		param = append(param, userParam.CompanyID.String)
	}

	if userParam.BranchID.String != "" {
		paramCount++
		query += " AND branch_id = $" + strconv.Itoa(paramCount) + " "
		param = append(param, userParam.BranchID.String)
	}

	if userParam.CreatedBy.Int64 > 0 {
		paramCount++
		query += " AND created_by = $" + strconv.Itoa(paramCount) + " "
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := tx.QueryRow(query, param...).Scan(
		&result.ID, &result.AuthUserID, &result.ClientID)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) GetPKCEClientMappingForUpdateByType(db *sql.Tx, pkceClient repository.PKCEClientMappingModel, clientType string, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.PKCEClientMappingModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetPKCEClientMappingForUpdateByType"
		query    string
	)

	query = fmt.Sprintf(`SELECT pkce_client_mapping.id, pkce_client_mapping.created_by, pkce_client_mapping.updated_at
		FROM %s pkce_client_mapping
		INNER JOIN %s ON pkce_client_mapping.client_type_id = client_type.id
		WHERE
		pkce_client_mapping.id = $1 AND
		client_type.client_type = $2 AND
		pkce_client_mapping.deleted = FALSE `,
		input.TableName, ClientTypeDAO.TableName)

	params := []interface{}{pkceClient.ID.Int64, clientType}
	if pkceClient.ClientID.String != "" {
		query += " AND pkce_client_mapping.client_id = $3 "
		params = append(params, pkceClient.ClientID.String)
	}

	additionalWhere := input.PrepareScopeInPKCEClientMapping(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	query += " FOR UPDATE"
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) UpdatePKCECLientMapping(db *sql.Tx, dataToUpdate repository.PKCEClientMappingModel) (err errorModel.ErrorModel) {
	funcName := "UpdatePKCECLientMapping"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			client_alias = $1, updated_by = $2,
			updated_at = $3, updated_client = $4
		WHERE
			id = $5 AND
			deleted = FALSE `,
		input.TableName)

	params := []interface{}{
		dataToUpdate.ClientAlias.String,
		dataToUpdate.UpdatedBy.Int64,
		dataToUpdate.UpdatedAt.Time,
		dataToUpdate.UpdatedClient.String,
		dataToUpdate.ID.Int64,
	}

	stmnt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmnt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input pkceClientMappingDAO) GetListPKCEClientMappingByJoin(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {

	query := fmt.Sprintf(
		`SELECT 
			userTable.id, pkce.parent_client_id, userTable.client_id, 
			pkce.username, userTable.created_by, userTable.updated_at
		FROM %s AS pkce `, input.TableName)

	for i := 0; i < len(searchBy); i++ {
		if searchBy[i].SearchKey == "username" {
			searchBy[i].SearchKey = "pkce." + searchBy[i].SearchKey
			searchBy[i].SearchType = constanta.Filter
			continue
		} else if searchBy[i].SearchKey == "parent_client_id" {
			searchBy[i].SearchKey = "pkce." + searchBy[i].SearchKey
			searchBy[i].SearchType = constanta.Filter
			continue
		}
	}

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "userTable.created_by",
			SearchType:     constanta.Filter,
			SearchOperator: "eq",
			SearchValue:    strconv.Itoa(int(createdBy)),
			DataType:       "number",
		})
	}

	createdBy = 0

	getListData := getListJoinDataDAO{Table: input.TableName, Query: query}
	getListData.InnerJoinAlias("\""+UserDAO.TableName+"\"", "userTable", "pkce.client_id", "userTable.client_id")

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var result repository.PKCEClientMappingModel

		dbError := rows.Scan(
			&result.ID,
			&result.ParentClientID,
			&result.ClientID,
			&result.Username,
			&result.CreatedBy,
			&result.UpdatedAt)

		return result, dbError
	}

	return getListData.GetListJoinDataWithoutDeleted(db, userParam, searchBy, createdBy, mappingFunc)
}

func (input pkceClientMappingDAO) GetListPKCEClientMapping(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query, additionalWhereStr string
		additionalWhere           []string
	)

	query = fmt.Sprintf(`SELECT pkce_client_mapping.id, client_mapping.client_id, userTable.first_name,
			userTable.last_name, userTable.nt_username, client_type.client_type,
			pkce_client_mapping.company_id, pkce_client_mapping.branch_id, pkce_client_mapping.client_alias,
			pkce_client_mapping.updated_at, pkce_client_mapping.updated_by, pkce_client_mapping.created_by,
			pkce_client_mapping.created_at
		FROM %s 
		INNER JOIN %s ON client_mapping.branch_id = pkce_client_mapping.branch_id 
		INNER JOIN "%s" AS userTable ON userTable.client_id = pkce_client_mapping.client_id 
		INNER JOIN %s ON client_type.id = pkce_client_mapping.client_type_id `,
		input.TableName, ClientMappingDAO.TableName, UserDAO.TableName,
		ClientTypeDAO.TableName)

	for index := range searchBy {
		switch searchBy[index].SearchKey {
		case "client_alias":
			searchBy[index].SearchKey = "pkce_client_mapping." + searchBy[index].SearchKey
		case "parent_client_id":
			searchBy[index].SearchKey = "pkce_client_mapping." + searchBy[index].SearchKey
		}
	}

	switch order := userParam.OrderBy; order {
	case "parent_client_id", "parent_client_id ASC", "parent_client_id DESC":
		userParam.OrderBy = "pkce_client_mapping." + userParam.OrderBy
	case "nt_username", "nt_username ASC", "nt_username DESC":
		userParam.OrderBy = "userTable." + userParam.OrderBy
	case "client_type", "client_type ASC", "client_type DESC":
		userParam.OrderBy = "client_type." + userParam.OrderBy
	case "company_id", "company_id ASC", "company_id DESC":
		userParam.OrderBy = "pkce_client_mapping." + userParam.OrderBy
	case "branch_id", "branch_id ASC", "branch_id DESC":
		userParam.OrderBy = "pkce_client_mapping." + userParam.OrderBy
	case "client_alias", "client_alias ASC", "client_alias DESC":
		userParam.OrderBy = "pkce_client_mapping." + userParam.OrderBy
	default:
		userParam.OrderBy = "client_mapping." + userParam.OrderBy
	}

	additionalWhere = input.PrepareScopeInPKCEClientMapping(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		additionalWhereStr = strWhere
	}

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.ViewPKCEClientMappingModel

		dbError := rows.Scan(&resultTemp.ID, &resultTemp.ParentClientID, &resultTemp.FirstName,
			&resultTemp.LastName, &resultTemp.Username, &resultTemp.ClientType,
			&resultTemp.CompanyID, &resultTemp.BranchID, &resultTemp.ClientAlias,
			&resultTemp.UpdatedAt, &resultTemp.UpdatedBy, &resultTemp.CreatedBy,
			&resultTemp.CreatedAt)

		return resultTemp, dbError
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		mappingFunc, additionalWhereStr, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "pkce_client_mapping.id"},
			Deleted:   FieldStatus{FieldName: "pkce_client_mapping.deleted"},
			CreatedBy: FieldStatus{FieldName: "pkce_client_mapping.created_by", Value: createdBy},
		})
}

func (input pkceClientMappingDAO) GetCountPKCEListClientMapping(db *sql.DB, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (results int, err errorModel.ErrorModel) {
	var (
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT COUNT(pkce_client_mapping.id) FROM %s `, input.TableName)

	for index := range searchBy {
		switch searchBy[index].SearchKey {
		case "client_alias":
			searchBy[index].SearchKey = "pkce_client_mapping." + searchBy[index].SearchKey
		case "parent_client_id":
			searchBy[index].SearchKey = "pkce_client_mapping." + searchBy[index].SearchKey
		}
	}

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "pkce_client_mapping.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}

	additionalWhere = input.PrepareScopeInPKCEClientMapping(scopeLimit, scopeDB, 1)

	getListData := getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: additionalWhere}
	getListData.InnerJoin(ClientMappingDAO.TableName, "client_mapping.branch_id", "pkce_client_mapping.branch_id")
	getListData.InnerJoinAlias("\""+UserDAO.TableName+"\"", "userTable", "userTable.client_id", "pkce_client_mapping.client_id")
	getListData.InnerJoin(ClientTypeDAO.TableName, "client_type.id", "pkce_client_mapping.client_type_id")

	return getListData.GetCountJoinDataWithoutDeleted(db, searchBy, 0)
}

func (input pkceClientMappingDAO) GetViewPKCEClientMapping(db *sql.DB, userParam repository.PKCEClientMappingModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ViewPKCEClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetViewPKCEClientMapping"
	query := fmt.Sprintf(
		`SELECT
			pkce_client_mapping.id,
			pkce_client_mapping.client_id,
			userTable.first_name,
			userTable.last_name,
			userTable.nt_username,
			client_type.client_type,
			pkce_client_mapping.company_id,
			pkce_client_mapping.branch_id,
			pkce_client_mapping.client_alias,
			pkce_client_mapping.updated_at,
			pkce_client_mapping.updated_by,
			pkce_client_mapping.created_by,
			pkce_client_mapping.created_at
		FROM %s
		INNER JOIN %s ON client_mapping.branch_id = pkce_client_mapping.branch_id
		INNER JOIN "%s" AS userTable ON userTable.client_id = pkce_client_mapping.client_id
		INNER JOIN %s ON client_type.id = pkce_client_mapping.client_type_id
		WHERE
			pkce_client_mapping.id = $1
			AND %s.deleted = FALSE
			AND %s.deleted = FALSE `,
		input.TableName, ClientMappingDAO.TableName, UserDAO.TableName,
		ClientTypeDAO.TableName, input.TableName, ClientMappingDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	if userParam.ClientID.String != "" {
		query += "AND (pkce_client_mapping.client_id = $2 OR pkce_client_mapping.parent_client_id = $3) "
		params = append(params, userParam.ClientID.String, userParam.ClientID.String)
	}

	additionalWhere := input.PrepareScopeInPKCEClientMapping(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID,
		&result.ClientID,
		&result.FirstName,
		&result.LastName,
		&result.Username,
		&result.ClientType,
		&result.CompanyID,
		&result.BranchID,
		&result.ClientAlias,
		&result.UpdatedAt,
		&result.UpdatedBy,
		&result.CreatedBy,
		&result.CreatedAt)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) CheckUserPKCEUnregisterBefore(db *sql.DB, userParam repository.CheckPKCEClientMappingModel) (result repository.CheckPKCEClientMappingModel, err errorModel.ErrorModel) {
	funcName := "CheckUserPKCEUnregisterBefore"

	query := fmt.Sprintf(
		`SELECT
			userTable.id, userTable.client_id, 
			userTable.auth_user_id, pcm.id
		FROM "%s" AS userTable
		INNER JOIN %s AS pcm ON userTable.client_id = pcm.client_id
		INNER JOIN %s AS cm ON pcm.parent_client_id = cm.client_id
		WHERE
			userTable.nt_username = $1 AND userTable.deleted = TRUE AND
			pcm.deleted = FALSE AND cm.deleted = FALSE AND
			userTable.email = $2 AND userTable.phone = $3 AND
			userTable.first_name = $4 AND userTable.last_name = $5 `,
		UserDAO.TableName, input.TableName, ClientMappingDAO.TableName)

	params := []interface{}{
		userParam.Username.String, userParam.Email.String, userParam.Phone.String,
		userParam.Firstname.String, userParam.Lastname.String,
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.ClientID, &result.AuthUserID,
		&result.PKCEClientMappingID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) UpdatePKCEClientMappingWithClientID(db *sql.Tx, dataToUpdate repository.PKCEClientMappingModel) (err errorModel.ErrorModel) {

	funcName := "UpdatePKCEClientMappingWithClientID"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			parent_client_id = $1, company_id = $2,
			branch_id = $3, client_alias = $4,
			updated_by = $5, updated_at = $6,
			updated_client = $7, created_by = $8,
			created_at = $9, created_client = $10
		WHERE
			id = $11 AND deleted = FALSE `, input.TableName)

	params := []interface{}{
		dataToUpdate.ParentClientID.String,
		dataToUpdate.CompanyID.String,
		dataToUpdate.BranchID.String,
		dataToUpdate.ClientAlias.String,
		dataToUpdate.UpdatedBy.Int64,
		dataToUpdate.UpdatedAt.Time,
		dataToUpdate.UpdatedClient.String,
		dataToUpdate.CreatedBy.Int64,
		dataToUpdate.CreatedAt.Time,
		dataToUpdate.CreatedClient.String,
		dataToUpdate.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input pkceClientMappingDAO) GetPKCEClientMappingForAddResource(db *sql.DB, pkceClient repository.PKCEClientMappingModel) (result repository.PKCEClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetPKCEClientMappingForAddResource"
	query := fmt.Sprintf(
		`SELECT
			pkce.id, pkce.created_by, pkce.updated_at
		FROM %s AS pkce
		INNER JOIN %s ON pkce.parent_client_id = client_mapping.client_id
		WHERE
			pkce.client_id = $1 AND
			pkce.parent_client_id = $2 AND
			pkce.deleted = FALSE `,
		input.TableName, ClientMappingDAO.TableName)
	params := []interface{}{pkceClient.ClientID.String, pkceClient.ParentClientID.String}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) ViewForUnregisterUserPKCE(db *sql.DB, userParam repository.PKCEClientMappingModel) (result repository.PKCEClientMappingModel, err errorModel.ErrorModel) {
	funcName := "ViewForUnregisterUserPKCE"

	query := fmt.Sprintf(
		`SELECT  
			userTable.id, pkce.parent_client_id, userTable.client_id,  
			pkce.username, userTable.created_by, userTable.updated_at  
		FROM %s AS pkce
		INNER JOIN "%s" AS userTable ON userTable.client_id = pkce.client_id
		WHERE 
			userTable.nt_username = $1 AND userTable.deleted = FALSE
			AND pkce.deleted = FALSE `,
		input.TableName, UserDAO.TableName)

	params := []interface{}{userParam.Username.String}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND userTable.created_by = $2"
		params = append(params, userParam.CreatedBy.Int64)
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID,
		&result.ParentClientID,
		&result.ClientID,
		&result.Username,
		&result.CreatedBy,
		&result.UpdatedAt)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) GetPKCEClientMappingForActivateUser(db *sql.DB, userParam repository.PKCEClientMappingModel) (result repository.PKCEClientMappingModel, err errorModel.ErrorModel) {
	var tempResult interface{}
	funcName := "GetPKCEClientMappingForActivateUser"
	query := fmt.Sprintf(`SELECT id,
		parent_client_id, client_id, client_type_id, 
		auth_user_id, installation_id, company_id, branch_id
	FROM %s
	WHERE 
		client_id = $1 AND auth_user_id = $2 AND
		company_id = $3 AND branch_id = $4 AND
		deleted = FALSE
	`, input.TableName)

	param := []interface{}{
		userParam.ClientID.String,
		userParam.AuthUserID.Int64,
		userParam.CompanyID.String,
		userParam.BranchID.String,
	}
	row := db.QueryRow(query, param...)
	tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.PKCEClientMappingModel
		errorS := rws.Scan(
			&temp.ID, &temp.ParentClientID, &temp.ClientID,
			&temp.ClientTypeID, &temp.AuthUserID, &temp.InstallationID,
			&temp.CompanyID, &temp.BranchID,
		)
		return temp, errorS
	}, input.FileName, funcName)
	if err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.PKCEClientMappingModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) GetFieldForViewNexmileParameter(db *sql.DB, userParam repository.PKCEClientMappingModel) (result repository.PKCEClientMappingModel, clientMapping repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetFieldForViewNexmileParameter"

	query := fmt.Sprintf(`SELECT 
			pkce.parent_client_id, cm.id, pkce.client_type_id, pkce.client_id,
			cm.client_type_id
		FROM %s pkce 
		JOIN %s cm ON pkce.parent_client_id = cm.client_id 
		WHERE 
		pkce.client_id = $1 AND cm.deleted = FALSE AND pkce.deleted = FALSE `,
		PKCEClientMappingDAO.TableName, ClientMappingDAO.TableName)

	params := []interface{}{
		userParam.ClientID.String,
	}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&clientMapping.ClientID,
		&clientMapping.ID,
		&result.ClientTypeID,
		&result.ClientID,
		&clientMapping.ClientTypeID,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingDAO) PrepareScopeInPKCEClientMapping(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int) (additionalWhere []string) {
	for key := range scopeLimit {
		var (
			keyDataScope   string
			additionalTemp string
		)

		switch key {
		case constanta.ClientTypeDataScope:
			keyDataScope = constanta.ClientTypeDataScope
		default:
			//--- No Key Limit Found
		}

		if keyDataScope != "" {
			PrepareScopeOnDAO(scopeLimit, scopeDB, &additionalTemp, idxStart, keyDataScope, true)
			if additionalTemp != "" {
				additionalWhere = append(additionalWhere, additionalTemp)
			}
		}
	}

	return
}

func (input pkceClientMappingDAO) GetPKCEClientWithClientIDAndUniqueID(db *sql.DB, userParam repository.PKCEClientMappingModel) (result repository.PKCEClientMappingModel, clientMapping repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetPKCEClientWithClientIDAndUniqueID"

	query := fmt.Sprintf(
		`SELECT 
		pkce.id, pkce.client_id, pkce.client_type_id, 
		cm.id, cm.client_id, cm.client_type_id
	FROM %s pkce
	LEFT JOIN %s cm ON cm.client_id = pkce.parent_client_id
	WHERE 
		pkce.client_id = $1 AND pkce.company_id = $2 AND pkce.branch_id = $3 
		AND cm.company_id = $2 AND cm.branch_id = $3
		AND pkce.client_type_id = $4 AND pkce.deleted = FALSE AND cm.deleted = FALSE `,
		input.TableName, ClientMappingDAO.TableName)

	param := []interface{}{
		userParam.ClientID.String, userParam.CompanyID.String, userParam.BranchID.String,
		userParam.ClientTypeID.Int64,
	}

	rows := db.QueryRow(query, param...)
	dbError := rows.Scan(
		&result.ID,
		&result.ClientID,
		&result.ClientTypeID,
		&clientMapping.ID,
		&clientMapping.ClientID,
		&clientMapping.ClientTypeID,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	return
}
