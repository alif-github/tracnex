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

type customerGroupDAO struct {
	AbstractDAO
}

var CustomerGroupDAO = customerGroupDAO{}.New()

func (input customerGroupDAO) New() (output customerGroupDAO) {
	output.FileName = "CustomerGroupDAO.go"
	output.TableName = "customer_group"
	return
}

func (input customerGroupDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "cg." + (*searchByParam)[i].SearchKey
	}

	switch userParam.OrderBy {
	case "updated_name", "updated_name ASC", "updated_name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "u.nt_username " + strSplit[1]
		} else {
			userParam.OrderBy = "u.nt_username"
		}
		break
	default:
		userParam.OrderBy = "cg." + userParam.OrderBy
		break
	}
}

func (input customerGroupDAO) getCustomerGroupDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "cg.id"},
		Deleted:   FieldStatus{FieldName: "cg.deleted"},
		CreatedBy: FieldStatus{FieldName: "cg.created_by", Value: createdBy},
	}
}

func (input customerGroupDAO) checkOwnPermission(createdBy int64, query *string, param *[]interface{}, index int) int {
	if createdBy > 0 {
		defaultField := input.getCustomerGroupDefaultMustCheck(createdBy)
		queryOwnPermission := " AND " + defaultField.CreatedBy.FieldName + " = $" + strconv.Itoa(index) + " "
		*query += queryOwnPermission
		(*param) = append((*param), createdBy)
		index++
	}
	return index
}

func (input customerGroupDAO) GetCustomerGroupForDelete(db *sql.Tx, userParam repository.CustomerGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerGroupModel, err errorModel.ErrorModel) {
	funcName := "GetCustomerGroupForDelete"
	index := 1
	var tempResult interface{}

	query := `SELECT 
			cg.id, cg.updated_at, cg.created_at, cg.created_by, cg.customer_group_id, 
		    CASE WHEN (SELECT COUNT(c.id) FROM customer c WHERE c.customer_group_id = cg.id AND c.deleted = false) > 0 
				THEN TRUE ELSE FALSE END is_used  
		FROM	
			` + input.TableName + ` cg
		WHERE
			cg.id = $1 AND cg.deleted = FALSE`

	param := []interface{}{userParam.ID.Int64}
	index += 1

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, []string{constanta.CustomerGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerGroupDefaultMustCheck, index)

	query += " FOR UPDATE "

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerGroupModel
		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt,
			&temp.CreatedAt, &temp.CreatedBy, &temp.CustomerGroupID,
			&temp.IsUsed,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerGroupModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerGroupDAO) ViewCustomerGroup(db *sql.DB, userParam repository.CustomerGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerGroupModel, err errorModel.ErrorModel) {
	funcName := "ViewCustomerGroup"
	index := 1
	var tempResult interface{}

	query := `SELECT 
			cg.id, cg.customer_group_id, cg.customer_group_name, 
		    cg.created_by, cg.created_at, 
		    cg.updated_by, cg.updated_at, u.nt_username 
		FROM ` + input.TableName + ` cg 
		LEFT JOIN "` + UserDAO.TableName + `" AS u ON cg.updated_by = u.id 
		WHERE
			cg.id = $1 AND cg.deleted = FALSE`
	param := []interface{}{userParam.ID.Int64}
	index += 1

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, []string{constanta.CustomerGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerGroupDefaultMustCheck, index)

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerGroupModel
		dbError := rws.Scan(&temp.ID, &temp.CustomerGroupID, &temp.CustomerGroupName,
			&temp.CreatedBy, &temp.CreatedAt,
			&temp.UpdatedBy, &temp.UpdatedAt,
			&temp.UpdatedName)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerGroupModel)
	}

	return
}

func (input customerGroupDAO) DeleteCustomerGroup(db *sql.Tx, userParam repository.CustomerGroupModel) (err errorModel.ErrorModel) {
	funcName := "DeleteCustomerGroup"

	query := `UPDATE ` + input.TableName + `
		SET 
			deleted = $1, 
			customer_group_id = $2, 
			updated_by = $3, 
			updated_at = $4, 
			updated_client = $5 
		WHERE 
			id = $6`
	param := []interface{}{
		true,
		userParam.CustomerGroupID.String,
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

	defer stmt.Close()

	return
}

func (input customerGroupDAO) UpdateCustomerGroup(db *sql.Tx, userParam repository.CustomerGroupModel) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomerGroup"

	query := `UPDATE ` + input.TableName + `
		SET 
			customer_group_name = $1, 
		    updated_by = $2, 
		    updated_at = $3,
		    updated_client = $4 
		WHERE
			id = $5`
	param := []interface{}{
		userParam.CustomerGroupName.String,
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

	defer stmt.Close()

	return
}

func (input customerGroupDAO) GetCustomerGroupForUpdate(db *sql.Tx, userParam repository.CustomerGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerGroupModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetCustomerGroupForUpdate"
		index      = 1
		tempResult interface{}
	)

	query := fmt.Sprintf(`SELECT 
		cg.id, cg.updated_at, cg.created_at, 
		cg.created_by, 
		CASE WHEN (SELECT COUNT(c.id) FROM customer c WHERE c.customer_group_id = cg.id AND c.deleted = false) > 0
			THEN TRUE ELSE FALSE END is_used
		FROM %s cg 
		WHERE 
		cg.id = $1 AND cg.deleted = FALSE `,
		input.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += 1

	//--- Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.CustomerGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	//--- Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerGroupDefaultMustCheck, index)

	query += fmt.Sprintf(` FOR UPDATE `)
	row := db.QueryRow(query, param...)
	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerGroupModel
		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt,
			&temp.CreatedAt, &temp.CreatedBy,
			&temp.IsUsed,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerGroupModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerGroupDAO) GetCountCustomerGroup(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (int, errorModel.ErrorModel) {
	var dbParam []interface{}

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.CustomerGroupDataScope})
	dbParam = append(dbParam, param...)
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, dbParam, input.TableName+" cg ",
		searchByParam, additionalWhere, input.getCustomerGroupDefaultMustCheck(createdBy))
}

func (input customerGroupDAO) GetListCustomerGroup(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var dbParam []interface{}

	query := `SELECT 
			cg.id, cg.customer_group_id, cg.customer_group_name, 
			cg.created_at, cg.updated_at, cg.updated_by, u.nt_username 
		FROM ` + input.TableName + ` cg
		LEFT JOIN "` + UserDAO.TableName + `" AS u ON cg.updated_by = u.id `

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.CustomerGroupDataScope})
	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CustomerGroupModel
			dbError := rows.Scan(
				&temp.ID, &temp.CustomerGroupID,
				&temp.CustomerGroupName, &temp.CreatedAt,
				&temp.UpdatedAt, &temp.UpdatedBy, &temp.UpdatedName,
			)
			return temp, dbError
		}, additionalWhere, input.getCustomerGroupDefaultMustCheck(createdBy))
}

func (input customerGroupDAO) InsertCustomerGroup(db *sql.Tx, userParam repository.CustomerGroupModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertCustomerGroup"

	query := `INSERT INTO ` + input.TableName + ` 
		(
		 customer_group_id, customer_group_name, 
		 created_by, created_at, created_client, 
		 updated_by, updated_at, updated_client
		)
		VALUES 
		(
		 $1, $2, $3, $4, $5, $6, $7, $8 
		)
		RETURNING id`

	params := []interface{}{
		userParam.CustomerGroupID.String,
		userParam.CustomerGroupName.String,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
	}
	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input customerGroupDAO) IsExistCustomerGroupForInsert(db *sql.DB, userParam repository.CustomerGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result bool, err errorModel.ErrorModel) {
	funcName := "IsExistCustomerGroupForInsert"
	index := 1

	query := `SELECT 
			CASE WHEN COUNT(cg.id) > 0 THEN TRUE ELSE FALSE END 
		FROM ` + input.TableName + ` cg
		WHERE
			cg.id = $1 AND cg.deleted = FALSE`

	param := []interface{}{userParam.ID.Int64}
	index += 1

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, []string{constanta.CustomerGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerGroupDefaultMustCheck, index)

	dbError := db.QueryRow(query, param...).Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
