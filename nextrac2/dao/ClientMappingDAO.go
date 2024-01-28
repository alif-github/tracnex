package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type clientMappingDAO struct {
	AbstractDAO
}

var ClientMappingDAO = clientMappingDAO{}.New()

func (input clientMappingDAO) New() (output clientMappingDAO) {
	output.FileName = "ClientMappingDAO.go"
	output.TableName = "client_mapping"
	return
}

func (input clientMappingDAO) InsertClientMapping(tx *sql.Tx, userParam *repository.ClientMappingModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertClientMapping"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`INSERT INTO %s (client_id, company_id, branch_id, 
		client_type_id, created_by, created_client, 
		created_at, updated_by, updated_client, 
		updated_at, client_alias, parent_customer_id, 
		customer_id, site_id) 
		VALUES ($1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12, 
		$13, $14) returning id `,
		input.TableName)

	param = append(param, userParam.ClientID.String, userParam.CompanyID.String, userParam.BranchID.String,
		userParam.ClientTypeID.Int64, userParam.CreatedBy.Int64, userParam.CreatedClient.String,
		userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time)

	if userParam.ClientAlias.String != "" {
		param = append(param, userParam.ClientAlias.String)
	} else {
		param = append(param, nil)
	}

	if userParam.ParentCustomerID.Int64 > 0 {
		param = append(param, userParam.ParentCustomerID.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.CustomerID.Int64 > 0 {
		param = append(param, userParam.CustomerID.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.SiteID.Int64 > 0 {
		param = append(param, userParam.SiteID.Int64)
	} else {
		param = append(param, nil)
	}

	errorS := tx.QueryRow(query, param...).Scan(&id)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) InsertMultipleClientMapping(tx *sql.Tx, userParam []repository.ClientMappingModel) (id []int64, err errorModel.ErrorModel) {
	funcName := "InsertMultipleClientMapping"
	parameterClientMapping := 16
	jVar := 1

	query := "INSERT INTO " + input.TableName + " " +
		"(client_id, company_id, branch_id, " +
		"client_type_id, created_by, created_client, " +
		"created_at, updated_by, updated_client, " +
		"updated_at, socket_id, client_alias, " +
		"installation_id, customer_id, site_id, " +
		"parent_customer_id" +
		") " +
		"VALUES "

	for i := 1; i <= len(userParam); i++ {
		query += "("

		for j := jVar; j <= parameterClientMapping; j++ {
			query += " $" + strconv.Itoa(j) + ""
			if parameterClientMapping-j != 0 {
				query += ","
			} else {
				query += ")"
			}
		}

		if len(userParam)-i != 0 {
			query += ","
		} else {
			query += " returning id"
		}

		jVar += 16
		parameterClientMapping += 16
	}
	var param []interface{}
	for i := 0; i < len(userParam); i++ {
		param = append(param,
			userParam[i].ClientID.String, userParam[i].CompanyID.String, userParam[i].BranchID.String,
			userParam[i].ClientTypeID.Int64)

		if userParam[i].CreatedBy.Int64 != 0 {
			param = append(param, userParam[i].CreatedBy.Int64)
		} else {
			param = append(param, nil)
		}

		if userParam[i].CreatedClient.String != "" {
			param = append(param, userParam[i].CreatedClient.String)
		} else {
			param = append(param, nil)
		}

		if !userParam[i].CreatedAt.Time.IsZero() {
			param = append(param, userParam[i].CreatedAt.Time)
		} else {
			param = append(param, nil)
		}

		if userParam[i].UpdatedBy.Int64 != 0 {
			param = append(param, userParam[i].UpdatedBy.Int64)
		} else {
			param = append(param, nil)
		}

		if userParam[i].UpdatedClient.String != "" {
			param = append(param, userParam[i].UpdatedClient.String)
		} else {
			param = append(param, nil)
		}

		if !userParam[i].UpdatedAt.Time.IsZero() {
			param = append(param, userParam[i].UpdatedAt.Time)
		} else {
			param = append(param, nil)
		}

		if userParam[i].SocketID.String != "" {
			param = append(param, userParam[i].SocketID.String)
		} else {
			param = append(param, nil)
		}

		if userParam[i].ClientAlias.String != "" {
			param = append(param, userParam[i].ClientAlias.String)
		} else {
			param = append(param, nil)
		}

		if userParam[i].InstallationID.Int64 > 0 {
			param = append(param, userParam[i].InstallationID.Int64)
		} else {
			param = append(param, nil)
		}

		if userParam[i].CustomerID.Int64 > 0 {
			param = append(param, userParam[i].CustomerID.Int64)
		} else {
			param = append(param, nil)
		}

		if userParam[i].SiteID.Int64 > 0 {
			param = append(param, userParam[i].SiteID.Int64)
		} else {
			param = append(param, nil)
		}

		if userParam[i].ParentCustomerID.Int64 > 0 {
			param = append(param, userParam[i].ParentCustomerID.Int64)
		} else {
			param = append(param, nil)
		}
	}

	rows, errorS := tx.Query(query, param...)
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
			var idTemp int64
			errorS = rows.Scan(&idTemp)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			id = append(id, idTemp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) CheckClientMapping(tx *sql.Tx, userParam []repository.ClientMappingModel, isCheckByClientID bool) (result []repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "CheckClientMapping"
	number := 1
	lengthUserParam := len(userParam)

	query := "select company_id, branch_id, client_id, client_alias, id from " + input.TableName + " where id IN ("

	for i := 1; i <= lengthUserParam; i++ {
		query += "(select id from " + input.TableName + " where company_id = $" + strconv.Itoa(number) + " AND branch_id = $" + strconv.Itoa(number+1)

		if userParam[i-1].ClientID.String != "" && isCheckByClientID {
			query += " AND client_id = $" + strconv.Itoa(number+2)
		}

		query += " order by id limit 1)"

		if lengthUserParam-i != 0 {
			query += ","
		} else {
			query += ")"
		}

		if userParam[i-1].ClientID.String != "" && isCheckByClientID {
			number += 3
		} else {
			number += 2
		}
	}
	var param []interface{}
	for i := 0; i < lengthUserParam; i++ {
		param = append(param, userParam[i].CompanyID.String, userParam[i].BranchID.String)

		if userParam[i].ClientID.String != "" {
			param = append(param, userParam[i].ClientID.String)
		}
	}

	rows, errorS := tx.Query(query, param...)
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
			var temp repository.ClientMappingModel
			errorS = rows.Scan(&temp.CompanyID, &temp.BranchID, &temp.ClientID, &temp.ClientAlias, &temp.ID)
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

func (input clientMappingDAO) GetParentClientMappingForClientValidation(db *sql.DB, userParam repository.ClientMappingModel) (output repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetParentClientMappingForClientValidation"

	query := fmt.Sprintf(`
			SELECT DISTINCT (ci.site_id), cm.id, cm.client_id, cm.parent_customer_id, cm.customer_id  
			FROM %s cm 
				JOIN %s ci
					ON ci.client_mapping_id = cm.id
			WHERE cm.client_id = $1 AND cm.company_id = $2 `,
		input.TableName, CustomerInstallationDAO.TableName)

	param := []interface{}{userParam.ClientID.String, userParam.CompanyID.String}

	if !util.IsStringEmpty(userParam.BranchID.String) {
		query += " AND cm.branch_id = $3"
		param = append(param, userParam.BranchID.String)
	}

	rows := db.QueryRow(query, param...)

	var tempResult interface{}
	if tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ClientMappingModel
		dbError := rws.Scan(
			&temp.SiteID, &temp.ID, &temp.ClientID, &temp.ParentCustomerID, &temp.CustomerID,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		output = tempResult.(repository.ClientMappingModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) GetClientMappingForRegistrationNexMileNexStar(db *sql.DB, userParam repository.ClientMappingModel) (result repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetClientMappingForRegistrationNexMileNexStar"

	query := fmt.Sprintf(`SELECT id, parent_customer_id, site_id FROM %s WHERE client_id = $1 and deleted = FALSE`, input.TableName)

	params := []interface{}{
		userParam.ClientID.String,
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.ParentCustomerID, &result.SiteID)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return

}

func (input clientMappingDAO) GetClientMapping(db *sql.DB, userParam in.DetailClientMappingDTO) (output []interface{}, err errorModel.ErrorModel) {
	funcName := "GetClientMapping"

	fmt.Println("table name : ", input.TableName)
	query := "SELECT " +
		"	cm.client_id, " +
		"	cm.client_type_id, " +
		"	cm.auth_user_id, " +
		"	cm.username, " +
		"	cm.company_id, " +
		"	cm.branch_id, " +
		"	cm.client_alias " +
		"FROM " + input.TableName + " cm " +
		//"LEFT JOIN user us ON cm.id = us.id " +
		"WHERE client_type_id = $1 " +
		"   AND company_id = $2 " +
		"   AND branch_id = $3 "

	//var result []repository.ClientMappingModel
	resultQuery, error := db.Query(query, "1", "NS0013060000371", "1473326058063")
	//resultQuery, error := db.Query(query, userParam.ClientTypeID, userParam.CompanyID, userParam.BranchID)

	if error != nil {
		fmt.Println("error sql ", error)
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, error)
		return
	}

	if resultQuery != nil {
		defer func() {
			error = resultQuery.Close()
			if error != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, error)
				return
			}
		}()
		for resultQuery.Next() {
			var data repository.ClientMappingModel
			error = resultQuery.Scan(
				&data.ClientID,
				&data.ClientTypeID,
				&data.AuthUserId,
				&data.UserName,
				&data.CompanyID,
				&data.BranchID,
				&data.ClientAlias,
			)
			if error != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, error)
				return
			}
			output = append(output, data)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, error)
		return
	}

	fmt.Println("result Query : ", resultQuery)
	fmt.Println("result final : ", output)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) GetDetailClientMapping(db *sql.DB, clients repository.ClientMappingForDetailModel, userParam in.GetListDataDTO, searchBy []in.SearchByParam, clientId string) (results []interface{}, err errorModel.ErrorModel) {
	//funcName := "GetDetailClientMapping"

	query := "SELECT " +
		"	client_mapping.id, " +
		"	client_mapping.client_id, " +
		"	client_mapping.client_type_id, " +
		"	client_mapping.company_id, " +
		"	client_mapping.branch_id, " +
		"	client_mapping.client_alias, " +
		"	user_client.auth_user_id, " +
		"	user_client.nt_username, " +
		"	client_mapping.socket_id, " +
		"	client_mapping.updated_at, " +
		"	client_mapping.updated_by, " +
		"	client_mapping.created_by, " +
		"	client_mapping.created_at " +
		" FROM " + input.TableName + " "

	getListData := getListJoinDataDAO{
		Query: query,
		Table: input.TableName,
	}

	getListData.LeftJoin(ClientTypeDAO.TableName, "client_mapping.client_type_id", "client_type.id")
	getListData.InnerJoinAlias("\""+UserDAO.TableName+"\"", "user_client", "client_mapping.client_id", "user_client.client_id")
	getListData.InnerJoinAlias(ClientRegistrationLogDAO.TableName, "crl", "client_mapping.client_id", "crl.client_id")

	getListData.SetWhere("client_mapping.client_type_id", strconv.Itoa(int(clients.ClientTypeID.Int64)))
	getListData.SetWhereAdditional("( crl.success_status_auth = TRUE AND crl.success_status_nexcloud = TRUE )")

	var whereParams []string
	whereParams = append(whereParams, " client_mapping.company_id IN ( ", " client_mapping.branch_id IN ( ")

	for i, companyData := range clients.CompanyData {
		whereParams[0] += " '" + companyData.CompanyID.String + "' "

		for j, branchData := range companyData.BranchData {
			whereParams[1] += " '" + branchData.BranchID.String + "' "
			if j < len(companyData.BranchData)-1 {
				whereParams[1] += ", "
			}
		}

		if i < len(clients.CompanyData)-1 {
			whereParams[0] += ", "
			whereParams[1] += ", "
		}
	}

	whereParams[0] += " ) "
	whereParams[1] += " ) "

	additionalQuery := strings.Join(whereParams, " AND ")
	getListData.SetWhereAdditional(additionalQuery)
	if clientId != "" {
		getListData.SetWhere("client_mapping.client_id", clientId)
	}

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var result repository.CLientMappingDetailForViewModel

		dbError := rows.Scan(
			&result.ID,
			&result.ClientId,
			&result.ClientTypeId,
			&result.CompanyId,
			&result.BranchId,
			&result.Aliases,
			&result.AuthUserId,
			&result.Username,
			&result.SocketID,
			&result.UpdatedAt,
			&result.UpdatedBy,
			&result.CreatedBy,
			&result.CreatedAt)

		return result, dbError
	}

	return getListData.GetListJoinDataWithoutDeleted(db, userParam, searchBy, 0, mappingFunc)
}

func (input clientMappingDAO) GetClientMappingForUpdateByType(db *sql.Tx, clientMapping repository.ClientMappingModel, clientType string, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ClientMappingModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetClientMappingForUpdateByType"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
		client_mapping.id, client_mapping.created_by, client_mapping.updated_at 
		FROM %s 
		INNER JOIN %s ON client_mapping.client_type_id = client_type.id 
		WHERE 
		client_mapping.id = $1 AND client_type.client_type = $2 AND client_mapping.deleted = FALSE `,
		input.TableName, ClientTypeDAO.TableName)

	params := []interface{}{clientMapping.ID.Int64, clientType}
	if clientMapping.ClientID.String != "" {
		query += " AND client_mapping.client_id = $3 "
		params = append(params, clientMapping.ClientID.String)
	}

	additionalWhere := input.PrepareScopeInClientMapping(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	query += " FOR UPDATE"
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.CreatedBy, &result.UpdatedAt)
	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) UpdateClientName(db *sql.Tx, dataToBeUpdated repository.ClientMappingModel) (err errorModel.ErrorModel) {
	funcName := "UpdateClientName"

	query := "UPDATE " + input.TableName + " " +
		"SET " +
		"	client_alias = $1, " +
		"	updated_by = $2, " +
		"	updated_at = $3, " +
		"	updated_client = $4 " +
		"WHERE " +
		"	id = $5 AND " +
		"	deleted = FALSE "
	params := []interface{}{
		dataToBeUpdated.ClientAlias.String,
		dataToBeUpdated.UpdatedBy.Int64,
		dataToBeUpdated.UpdatedAt.Time,
		dataToBeUpdated.UpdatedClient.String,
		dataToBeUpdated.ID.Int64,
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

func (input clientMappingDAO) UpdateSocketID(db *sql.Tx, dataToBeUpdated repository.ClientMappingModel) (err errorModel.ErrorModel) {
	funcName := "UpdateClientName"

	query := "UPDATE " + input.TableName + " " +
		"SET " +
		"	socket_id = $1, " +
		"	updated_by = $2, " +
		"	updated_at = $3, " +
		"	updated_client = $4 " +
		"WHERE " +
		"	client_type_id = $5 AND " +
		"	client_id = $6 AND " +
		"	deleted = FALSE "

	var param []interface{}

	if dataToBeUpdated.SocketID.String != "" {
		param = append(param, dataToBeUpdated.SocketID.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, dataToBeUpdated.UpdatedBy.Int64, dataToBeUpdated.UpdatedAt.Time, dataToBeUpdated.UpdatedClient.String,
		dataToBeUpdated.ClientTypeID.Int64, dataToBeUpdated.ClientID.String)

	stmnt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmnt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input clientMappingDAO) GetClientMappingsForChangeSocketID(db *sql.Tx, clientMapping repository.ClientMappingModel) (results []repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetClientMappingsForChangeSocketID"
	query := "SELECT " +
		"	id, created_by, updated_at " +
		" FROM " + input.TableName + " " +
		" WHERE " +
		"	client_type_id = $1 AND " +
		"	client_id = $2 AND " +
		"deleted = FALSE "
	params := []interface{}{
		clientMapping.ClientTypeID.Int64,
		clientMapping.ClientID.String}
	//if clientMapping.CreatedBy.Int64 > 0 {
	//	query += " AND created_by = $4 "
	//	params = append(params, clientMapping.CreatedBy.Int64)
	//}

	query += " FOR UPDATE"
	rows, errorDB := db.Query(query, params...)
	if errorDB != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	if rows != nil {
		defer func() {
			errorDB = rows.Close()
			if errorDB != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
				return
			}
		}()
		for rows.Next() {
			var result repository.ClientMappingModel
			errorDB = rows.Scan(
				&result.ID,
				&result.CreatedBy,
				&result.UpdatedAt)
			if errorDB != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
				return
			}
			results = append(results, result)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) GetDataHeaderClientMapping(db *sql.Tx, clientMapping repository.ClientMappingModel) (result repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetDataHeaderClientMapping"

	query := "SELECT " +
		" client_id, socket_id, id " +
		" FROM " + input.TableName + " " +
		" WHERE " +
		" client_id = $1 AND deleted = FALSE " +
		" order by id asc limit 1"

	params := []interface{}{clientMapping.ClientID.String}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&result.ClientID, &result.SocketID, &result.ID)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) GetListClientMapping(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (results []interface{}, err errorModel.ErrorModel) {
	var (
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT client_mapping.id, client_mapping.client_id, client_mapping.socket_id, 
		client_type.client_type, client_mapping.company_id, client_mapping.branch_id, 
		client_mapping.client_alias, client_mapping.updated_at, client_mapping.updated_by, 
		client_mapping.created_by, client_mapping.created_at, crl.success_status_auth, 
		crl.success_status_nexcloud, crl.success_status_nexdrive 
		FROM %s `, input.TableName)

	for index := range searchBy {
		switch searchBy[index].SearchKey {
		case "client_alias":
			searchBy[index].SearchKey = "client_mapping." + searchBy[index].SearchKey
		case "client_id":
			searchBy[index].SearchKey = "client_mapping." + searchBy[index].SearchKey
		case "success_status_nexcloud":
			searchBy[index].SearchKey = "crl." + searchBy[index].SearchKey
		}
	}

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "client_mapping.created_by",
			SearchOperator: "eq",
			SearchValue:    strconv.Itoa(int(createdBy)),
			DataType:       "number",
			SearchType:     constanta.Filter,
		})
	}

	switch order := userParam.OrderBy; order {
	case "client_id", "client_id ASC", "client_id DESC":
		userParam.OrderBy = "client_mapping." + userParam.OrderBy
	case "socket_id", "socket_id ASC", "socket_id DESC":
		userParam.OrderBy = "client_mapping." + userParam.OrderBy
	case "client_type", "client_type ASC", "client_type DESC":
		userParam.OrderBy = "client_type." + userParam.OrderBy
	case "company_id", "company_id ASC", "company_id DESC":
		userParam.OrderBy = "client_mapping." + userParam.OrderBy
	case "branch_id", "branch_id ASC", "branch_id DESC":
		userParam.OrderBy = "client_mapping." + userParam.OrderBy
	case "client_alias", "client_alias ASC", "client_alias DESC":
		userParam.OrderBy = "client_mapping." + userParam.OrderBy
	default:
		userParam.OrderBy = "client_mapping." + userParam.OrderBy
	}

	additionalWhere = input.PrepareScopeInClientMapping(scopeLimit, scopeDB, 1)

	getListData := getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: additionalWhere}
	getListData.LeftJoin(ClientTypeDAO.TableName, "client_mapping.client_type_id", "client_type.id")
	getListData.InnerJoin(ClientCredentialDAO.TableName, "client_mapping.client_id", "client_credential.client_id")
	getListData.LeftJoinAlias(ClientRegistrationLogDAO.TableName, "crl", "client_mapping.client_id", "crl.client_id")
	getListData.InnerJoinAlias("\""+UserDAO.TableName+"\"", "user_client", "client_credential.client_id", "user_client.client_id")

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var result repository.ClientMappingForViewModel

		dbError := rows.Scan(
			&result.ID,
			&result.ClientID,
			&result.SocketID,
			&result.ClientType,
			&result.CompanyID,
			&result.BranchID,
			&result.Aliases,
			&result.UpdatedAt,
			&result.UpdatedBy,
			&result.CreatedBy,
			&result.CreatedAt,
			&result.SuccessStatusAuth,
			&result.SuccessStatusNexcloud,
			&result.SuccessStatusNexdrive,
		)

		return result, dbError
	}

	return getListData.GetListJoinDataWithoutDeleted(db, userParam, searchBy, 0, mappingFunc)
}

func (input clientMappingDAO) GetCountListClientMapping(db *sql.DB, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (results int, err errorModel.ErrorModel) {
	var (
		query           string
		getListData     getListJoinDataDAO
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT COUNT(client_mapping.id) FROM %s `, input.TableName)

	for index := range searchBy {
		switch searchBy[index].SearchKey {
		case "client_alias":
			searchBy[index].SearchKey = fmt.Sprintf(`client_mapping.%s`, searchBy[index].SearchKey)
		case "client_id":
			searchBy[index].SearchKey = fmt.Sprintf(`client_mapping.%s`, searchBy[index].SearchKey)
		case "success_status_nexcloud":
			searchBy[index].SearchKey = fmt.Sprintf(`crl.%s`, searchBy[index].SearchKey)
		}
	}

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "client_mapping.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}

	additionalWhere = input.PrepareScopeInClientMapping(scopeLimit, scopeDB, 1)

	getListData = getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: additionalWhere}
	getListData.LeftJoin(ClientTypeDAO.TableName, "client_mapping.client_type_id", "client_type.id")
	getListData.InnerJoin(ClientCredentialDAO.TableName, "client_mapping.client_id", "client_credential.client_id")
	getListData.LeftJoinAlias(ClientRegistrationLogDAO.TableName, "crl", "client_mapping.client_id", "crl.client_id")
	getListData.InnerJoinAlias("\""+UserDAO.TableName+"\"", "user_client", "client_credential.client_id", "user_client.client_id")

	return getListData.GetCountJoinDataWithoutDeleted(db, searchBy, 0)
}

func (input clientMappingDAO) ViewClientMapping(db *sql.DB, userParam repository.ClientMappingModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ClientMappingForViewModel, err errorModel.ErrorModel) {
	funcName := "ViewClientMapping"

	query := "SELECT " +
		"	client_mapping.id, " +
		"	client_mapping.client_id, " +
		"	client_mapping.socket_id, " +
		"	client_type.client_type, " +
		"	client_mapping.company_id, " +
		"	client_mapping.branch_id, " +
		"	client_mapping.client_alias, " +
		"	client_mapping.updated_at, " +
		"	client_mapping.updated_by, " +
		"	client_mapping.created_by, " +
		"	client_mapping.created_at, " +
		"	crl.success_status_auth, " +
		"	crl.success_status_nexcloud, " +
		"	crl.success_status_nexdrive " +
		"FROM " + input.TableName + " " +
		"LEFT JOIN " + ClientTypeDAO.TableName + " " +
		"	ON client_mapping.client_type_id = client_type.id " +
		"INNER JOIN " + ClientCredentialDAO.TableName + " " +
		"	ON client_mapping.client_id = client_credential.client_id " +
		"INNER JOIN \"" + UserDAO.TableName + "\" AS user_client " +
		"	ON client_credential.client_id = user_client.client_id " +
		"LEFT JOIN " + ClientRegistrationLogDAO.TableName + " AS crl " +
		"	ON client_mapping.client_id = crl.client_id " +
		"WHERE " +
		" client_mapping.id = $1 " +
		" AND " + ClientCredentialDAO.TableName + ".deleted = FALSE "

	params := []interface{}{
		userParam.ID.Int64,
	}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND client_mapping.created_by = $2 "
		params = append(params, userParam.CreatedBy.Int64)
	}

	additionalWhere := input.PrepareScopeInClientMapping(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	results := db.QueryRow(query, params...)

	errorDB := results.Scan(
		&result.ID,
		&result.ClientID,
		&result.SocketID,
		&result.ClientType,
		&result.CompanyID,
		&result.BranchID,
		&result.Aliases,
		&result.UpdatedAt,
		&result.UpdatedBy,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.SuccessStatusAuth,
		&result.SuccessStatusNexcloud,
		&result.SuccessStatusNexdrive,
	)
	if errorDB != nil && errorDB.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) GetDetailClientMappingByClientID(db *sql.DB, clients repository.ClientMappingForDetailModel, userParam in.GetListDataDTO, searchBy []in.SearchByParam, clientId string) (results []interface{}, err errorModel.ErrorModel) {
	//funcName := "GetDetailClientMappingByClientID"

	query := "SELECT " +
		"	client_mapping.id, " +
		"	client_mapping.client_id, " +
		"	client_mapping.client_type_id, " +
		"	client_mapping.company_id, " +
		"	client_mapping.branch_id, " +
		"	client_mapping.client_alias, " +
		"	user_client.auth_user_id, " +
		"	user_client.nt_username, " +
		"	client_mapping.socket_id, " +
		"	client_mapping.updated_at, " +
		"	client_mapping.updated_by, " +
		"	client_mapping.created_by, " +
		"	client_mapping.created_at " +
		" FROM " + input.TableName + " "

	getListData := getListJoinDataDAO{
		Query: query,
		Table: input.TableName,
	}

	getListData.LeftJoin(ClientTypeDAO.TableName, "client_mapping.client_type_id", "client_type.id")
	getListData.InnerJoinAlias("\""+UserDAO.TableName+"\"", "user_client", "client_mapping.client_id", "user_client.client_id")
	getListData.InnerJoinAlias(ClientRegistrationLogDAO.TableName, "crl", "client_mapping.client_id", "crl.client_id")

	getListData.SetWhere("client_mapping.client_type_id", strconv.Itoa(int(clients.ClientTypeID.Int64)))
	getListData.SetWhereAdditional("( crl.success_status_auth = TRUE AND crl.success_status_nexcloud = TRUE )")

	var strClientIDs []string

	for _, clientData := range clients.ClientData {
		strClientIDs = append(strClientIDs, " '"+clientData.ClientID.String+"' ")
	}

	getListData.SetWhereAdditional(" client_mapping.client_id IN ( " + strings.Join(strClientIDs, ", ") + ") ")

	if clientId != "" {
		getListData.SetWhere("client_mapping.client_id", clientId)
	}

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var result repository.CLientMappingDetailForViewModel

		dbError := rows.Scan(
			&result.ID,
			&result.ClientId,
			&result.ClientTypeId,
			&result.CompanyId,
			&result.BranchId,
			&result.Aliases,
			&result.AuthUserId,
			&result.Username,
			&result.SocketID,
			&result.UpdatedAt,
			&result.UpdatedBy,
			&result.CreatedBy,
			&result.CreatedAt)

		return result, dbError
	}

	return getListData.GetListJoinDataWithoutDeleted(db, userParam, searchBy, 0, mappingFunc)
}

func (input clientMappingDAO) GetClientIDWithID(db *sql.DB, clientMapping repository.ClientMappingModel) (result repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "GetClientIDWithID"

	query := "SELECT " +
		" cm.id, cm.client_id " +
		" FROM " + input.TableName + " AS cm " +
		" WHERE " +
		" cm.id = $1 AND " +
		" cm.deleted = FALSE "

	params := []interface{}{clientMapping.ID.Int64}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&result.ID, &result.ClientID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) CheckClientMappingWithClientID(db *sql.DB, clientMapping repository.ClientMappingModel) (result repository.ClientMappingModel, err errorModel.ErrorModel) {
	funcName := "CheckClientMappingWithClientID"

	query := "SELECT id, installation_id, customer_id, " +
		" site_id, company_id, branch_id, " +
		" parent_customer_id, client_type_id, client_id " +
		" FROM " + input.TableName + " " +
		" WHERE " +
		" client_id = $1 AND deleted = FALSE AND company_id = $2 "

	params := []interface{}{clientMapping.ClientID.String, clientMapping.CompanyID.String}

	if clientMapping.BranchID.String != "" {
		query += " AND branch_id = $3 "
		params = append(params, clientMapping.BranchID.String)
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.ID, &result.InstallationID, &result.CustomerID,
		&result.SiteID, &result.CompanyID, &result.BranchID,
		&result.ParentCustomerID, &result.ClientTypeID, &result.ClientID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) PrepareScopeInClientMapping(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int) (additionalWhere []string) {
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

func (input clientMappingDAO) CheckClientMappingByUniqueID12(db *sql.DB, clientMapping repository.ClientMappingModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "CheckClientMappingByUniqueID12"
		query    string
		params   []interface{}
		results  *sql.Row
		dbError  error
	)

	query = fmt.Sprintf(`SELECT id FROM %s WHERE company_id = $1 AND branch_id = $2 AND deleted = false `, input.TableName)
	params = []interface{}{clientMapping.CompanyID.String, clientMapping.BranchID.String}
	results = db.QueryRow(query, params...)
	dbError = results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingDAO) GetClientMappingByUniqueID12(db *sql.Tx, clientMapping repository.ClientMappingModel) (output repository.ClientMappingModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetClientMappingByUniqueID12"
		query    string
		params   []interface{}
		dbError  error
	)

	query = fmt.Sprintf(`
		SELECT id, company_id, branch_id 
		FROM %s 
		WHERE 
		deleted = false AND company_id = $1 AND branch_id = $2 AND 
		client_id = $3 `, input.TableName)

	params = append(params, clientMapping.CompanyID.String, clientMapping.BranchID.String, clientMapping.ClientID.String)
	rows := db.QueryRow(query, params...)
	dbError = rows.Scan(&output.ID, &output.CompanyID, &output.BranchID)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
