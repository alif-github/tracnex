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

type productGroupDAO struct {
	AbstractDAO
}

var ProductGroupDAO = productGroupDAO{}.New()

func (input productGroupDAO) New() (output productGroupDAO) {
	output.FileName = "ProductGroupDAO.go"
	output.TableName = "product_group"
	return
}

func (input productGroupDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "pg." + (*searchByParam)[i].SearchKey
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
		userParam.OrderBy = "pg." + userParam.OrderBy
		break
	}
}

func (input productGroupDAO) getProductGroupDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "pg.id"},
		Deleted:   FieldStatus{FieldName: "pg.deleted"},
		CreatedBy: FieldStatus{FieldName: "pg.created_by", Value: createdBy},
	}
}

func (input productGroupDAO) GetProductGroupForDelete(db *sql.Tx, userParam repository.ProductGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ProductGroupModel, err errorModel.ErrorModel) {
	funcName := "GetProductGroupForDelete"
	index := 1
	var tempResult interface{}

	query := fmt.Sprintf(
		`SELECT
			pg.id, pg.updated_at,
			pg.created_at, pg.created_by, pg.product_group_name,
			CASE WHEN 
				(SELECT COUNT(p.id) FROM product p WHERE p.product_group_id = pg.id AND p.deleted = false) > 0
			THEN TRUE ELSE FALSE END is_used
		FROM %s pg
		WHERE
			pg.id = $1 AND pg.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += len(param)

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.ProductGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getProductGroupDefaultMustCheck, index)

	query += " FOR UPDATE "

	row := db.QueryRow(query, param...)

	tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ProductGroupModel
		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt, &temp.CreatedAt,
			&temp.CreatedBy, &temp.ProductGroupName, &temp.IsUsed,
		)
		return temp, dbError
	}, input.FileName, funcName)

	if tempResult != nil {
		result = tempResult.(repository.ProductGroupModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productGroupDAO) ViewProductGroup(db *sql.DB, userParam repository.ProductGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ProductGroupModel, err errorModel.ErrorModel) {
	funcName := "ViewProductGroup"
	index := 1
	var tempResult interface{}

	query := fmt.Sprintf(
		`SELECT
			pg.id, pg.product_group_name,
			pg.created_by, pg.created_at,
			pg.updated_by, pg.updated_at,
			u.nt_username
		FROM %s pg
		LEFT JOIN "%s" AS u ON pg.updated_by = u.id
		WHERE
			pg.id = $1 AND pg.deleted = FALSE `, input.TableName, UserDAO.TableName)
	param := []interface{}{userParam.ID.Int64}
	index += len(param)

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, []string{constanta.ProductGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getProductGroupDefaultMustCheck, index)

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ProductGroupModel

		dbError := rws.Scan(&temp.ID, &temp.ProductGroupName,
			&temp.CreatedBy, &temp.CreatedAt,
			&temp.UpdatedBy, &temp.UpdatedAt,
			&temp.UpdatedName)

		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.ProductGroupModel)
	}

	return
}

func (input productGroupDAO) DeleteProductGroup(db *sql.Tx, userParam repository.ProductGroupModel) (err errorModel.ErrorModel) {
	funcName := "DeleteProductGroup"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			deleted = $1, product_group_name = $2, updated_by = $3,
			updated_at = $4, updated_client = $5
		WHERE
			id = $6 `, input.TableName)

	param := []interface{}{
		true,
		userParam.ProductGroupName.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
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

func (input productGroupDAO) UpdateProductGroup(db *sql.Tx, userParam repository.ProductGroupModel) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomerGroup"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			product_group_name = $1,
			updated_by = $2,
			updated_at = $3,
			updated_client = $4
		WHERE
			id = $5 `, input.TableName)
	param := []interface{}{
		userParam.ProductGroupName.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
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

func (input productGroupDAO) GetProductGroupForUpdate(db *sql.Tx, userParam repository.ProductGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ProductGroupModel, err errorModel.ErrorModel) {
	funcName := "GetProductGroupForUpdate"
	index := 1
	var tempResult interface{}

	query := fmt.Sprintf(
		`SELECT
			pg.id, pg.updated_at,
			pg.created_at, pg.created_by
		FROM %s pg
		WHERE
			pg.id = $1 AND pg.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += len(param)

	// Add Data scope
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, []string{constanta.ProductGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
		index += len(scopeParam)
	}

	// Check own access
	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getProductGroupDefaultMustCheck, index)

	query += " FOR UPDATE "

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ProductGroupModel

		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt,
			&temp.CreatedAt, &temp.CreatedBy,
		)

		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.ProductGroupModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productGroupDAO) GetCountProductGroup(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (int, errorModel.ErrorModel) {
	var dbParam []interface{}

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.ProductGroupDataScope})
	dbParam = append(dbParam, param...)
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, dbParam, input.TableName+" pg ",
		searchByParam, additionalWhere, input.getProductGroupDefaultMustCheck(createdBy))
}

func (input productGroupDAO) GetListProductGroup(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var dbParam []interface{}
	query := fmt.Sprintf(
		`SELECT
			pg.id, pg.product_group_name,
			pg.created_at, pg.updated_at, pg.updated_by,
			u.nt_username
		FROM %s pg
		LEFT JOIN "%s" AS u ON pg.updated_by = u.id `,
		input.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.ProductGroupDataScope})
	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ProductGroupModel
			dbError := rows.Scan(
				&temp.ID, &temp.ProductGroupName,
				&temp.CreatedAt, &temp.UpdatedAt, &temp.UpdatedBy,
				&temp.UpdatedName,
			)
			return temp, dbError
		}, additionalWhere, input.getProductGroupDefaultMustCheck(createdBy))
}

func (input productGroupDAO) InsertProductGroup(db *sql.Tx, userParam repository.ProductGroupModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertCustomerGroup"

	query := fmt.Sprintf(
		`INSERT INTO %s
		(
			product_group_name,
			created_by, created_at, created_client,
			updated_by, updated_at, updated_client
		)
		VALUES ( $1, $2, $3, $4, $5, $6, $7 )
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.ProductGroupName.String, userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time, userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
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

func (input productGroupDAO) CheckIsProductGroupExist(db *sql.DB, userParam repository.ProductGroupModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ProductGroupModel, err errorModel.ErrorModel) {
	var (
		funcName = "CheckIsProductGroupExist"
		query    string
	)

	query = fmt.Sprintf(`SELECT pg.id FROM %s pg WHERE pg.id = $1 AND pg.deleted = FALSE `, input.TableName)
	param := []interface{}{userParam.ID.Int64}

	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.ProductGroupDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
	}

	dbError := db.QueryRow(query, param...).Scan(&result.ID)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
