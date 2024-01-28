package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type customerCategoryDAO struct {
	AbstractDAO
}

var CustomerCategoryDAO = customerCategoryDAO{}.New()

func (input customerCategoryDAO) New() (output customerCategoryDAO) {
	output.FileName = "CustomerCategoryDAO.go"
	output.TableName = "customer_category"
	return
}

func (input customerCategoryDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "cc." + (*searchByParam)[i].SearchKey
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
		userParam.OrderBy = "cc." + userParam.OrderBy
		break
	}
}

func (input customerCategoryDAO) getCustomerCategoryDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "cc.id"},
		Deleted:   FieldStatus{FieldName: "cc.deleted"},
		CreatedBy: FieldStatus{FieldName: "cc.created_by", Value: createdBy},
	}
}

func (input customerCategoryDAO) GetCustomerCategoryForDelete(db *sql.Tx, userParam repository.CustomerCategoryModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerCategoryModel, err errorModel.ErrorModel) {
	funcName := "GetCustomerCategoryForDelete"
	index := 1
	var tempResult interface{}

	query := fmt.Sprintf(
		`SELECT 
			cc.id, cc.updated_at, cc.created_at, 
			cc.created_by, cc.customer_category_id, 
			CASE WHEN 
				(SELECT COUNT(c.id) FROM customer c WHERE c.customer_category_id = cc.id AND c.deleted = false) > 0 
			THEN TRUE ELSE FALSE END is_used 
		FROM %s cc 
		WHERE cc.id = $1 AND cc.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += len(param)

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, []string{constanta.CustomerCategoryDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerCategoryDefaultMustCheck, index)

	query += " FOR UPDATE "

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerCategoryModel
		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt,
			&temp.CreatedAt, &temp.CreatedBy, &temp.CustomerCategoryID,
			&temp.IsUsed,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerCategoryModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerCategoryDAO) DeleteCustomerCategory(db *sql.Tx, userParam repository.CustomerCategoryModel) (err errorModel.ErrorModel) {
	funcName := "DeleteCustomerCategory"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			deleted = $1,
			customer_category_id = $2,
			updated_by = $3,
			updated_at = $4,
			updated_client = $5
		WHERE
			id = $6 `, input.TableName)

	param := []interface{}{
		true, userParam.CustomerCategoryID.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.ID.Int64,
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

func (input customerCategoryDAO) UpdateCustomerCategory(db *sql.Tx, userParam repository.CustomerCategoryModel) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomerCategory"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			customer_category_name = $1, updated_by = $2,
			updated_at = $3, updated_client = $4
		WHERE
			id = $5 `, input.TableName)

	param := []interface{}{
		userParam.CustomerCategoryName.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.ID.Int64,
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

func (input customerCategoryDAO) GetCustomerCategoryForUpdate(db *sql.Tx, userParam repository.CustomerCategoryModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerCategoryModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetCustomerCategoryForUpdate"
		index      = 1
		tempResult interface{}
		query      string
	)

	query = fmt.Sprintf(
		`SELECT
			cc.id, cc.updated_at, cc.created_at, 
			cc.created_by, 
			CASE WHEN 
				(SELECT COUNT(c.id) FROM customer c WHERE c.customer_category_id = cc.id AND c.deleted = false) > 0 
			THEN TRUE ELSE FALSE END is_used
		FROM %s cc
		WHERE cc.id = $1 AND cc.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += len(param)

	//--- Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.CustomerCategoryDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	//--- Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerCategoryDefaultMustCheck, index)

	query += fmt.Sprintf(` FOR UPDATE `)
	row := db.QueryRow(query, param...)
	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerCategoryModel
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
		result = tempResult.(repository.CustomerCategoryModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerCategoryDAO) ViewCustomerCategory(db *sql.DB, userParam repository.CustomerCategoryModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerCategoryModel, err errorModel.ErrorModel) {
	funcName := "ViewCustomerCategory"
	index := 1
	var tempResult interface{}

	query := fmt.Sprintf(
		`SELECT
			cc.id, cc.customer_category_id, cc.customer_category_name,
			cc.created_by, cc.created_at,
			cc.updated_by, cc.updated_at,
			u.nt_username
		FROM %s cc
		LEFT JOIN "%s" AS u ON cc.updated_by = u.id
		WHERE
			cc.id = $1 AND cc.deleted = FALSE `, input.TableName, UserDAO.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += len(param)

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.CustomerCategoryDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerCategoryDefaultMustCheck, index)

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerCategoryModel
		dbError := rws.Scan(&temp.ID, &temp.CustomerCategoryID, &temp.CustomerCategoryName,
			&temp.CreatedBy, &temp.CreatedAt, &temp.UpdatedBy, &temp.UpdatedAt,
			&temp.UpdatedName)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerCategoryModel)
	}

	return
}

func (input customerCategoryDAO) GetCountCustomerCategory(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (int, errorModel.ErrorModel) {
	var dbParam []interface{}
	var idxStart int = 1
	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "cc." + param.SearchKey
	}

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, idxStart, []string{constanta.CustomerCategoryDataScope})
	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, dbParam, input.TableName+" cc ",
		searchByParam, additionalWhere, input.getCustomerCategoryDefaultMustCheck(createdBy))
}

func (input customerCategoryDAO) GetListCustomerCategory(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var dbParam []interface{}

	query := fmt.Sprintf(
		`SELECT
			cc.id as id, cc.customer_category_id as customer_category_id, 
			cc.customer_category_name as customer_category_name, cc.created_at as created_at, 
			cc.updated_at as updated_at, cc.updated_by as updated_by, u.nt_username
		FROM %s cc
		LEFT JOIN "%s" AS u ON cc.updated_by = u.id `, input.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.CustomerCategoryDataScope})
	dbParam = append(dbParam, param...)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CustomerCategoryModel
			dbError := rows.Scan(
				&temp.ID, &temp.CustomerCategoryID,
				&temp.CustomerCategoryName, &temp.CreatedAt,
				&temp.UpdatedAt, &temp.UpdatedBy, &temp.UpdatedName,
			)
			return temp, dbError
		}, additionalWhere, input.getCustomerCategoryDefaultMustCheck(createdBy))
}

func (input customerCategoryDAO) InsertCustomerCategory(db *sql.Tx, userParam repository.CustomerCategoryModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertCustomerCategory"

	query := fmt.Sprintf(
		`INSERT INTO %s
		(
			customer_category_id, customer_category_name,
			created_by, created_at, created_client,
			updated_by, updated_at, updated_client
		)
		VALUES
		(
			$1, $2, $3,
			$4, $5, $6,
			$7, $8
		)
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.CustomerCategoryID.String, userParam.CustomerCategoryName.String,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time,
		userParam.CreatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input customerCategoryDAO) IsExistCustomerCategory(db *sql.DB, userParam repository.CustomerCategoryModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result bool, err errorModel.ErrorModel) {
	funcName := "IsExistCustomerCategory"
	index := 1

	query := fmt.Sprintf(
		`SELECT 
			CASE WHEN 
				COUNT(cc.id) > 0 
			THEN TRUE ELSE FALSE END 
		FROM %s cc 
		WHERE cc.id = $1 AND cc.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += len(param)

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, []string{constanta.CustomerCategoryDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerCategoryDefaultMustCheck, index)

	dbError := db.QueryRow(query, param...).Scan(
		&result,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
