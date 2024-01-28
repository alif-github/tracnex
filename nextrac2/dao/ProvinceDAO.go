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
	"time"
)

type provinceDAO struct {
	AbstractDAO
}

var ProvinceDAO = provinceDAO{}.New()

func (input provinceDAO) New() (output provinceDAO) {
	output.FileName = "ProvinceDAO.go"
	output.TableName = "province"
	return
}

func (input provinceDAO) GetProvinceIDByMdbID(db *sql.DB, userParam repository.ProvinceModel, isMustNotCheckDeleted bool) (result int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetProvinceIDByMdbID"
		query    string
	)

	query += fmt.Sprintf(`SELECT id FROM %s WHERE mdb_province_id = $1 `, input.TableName)
	if !isMustNotCheckDeleted {
		query += fmt.Sprintf(` AND deleted = FALSE `)
	}

	params := []interface{}{userParam.MDBProvinceID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceDAO) IsExistProvince(db *sql.DB, userParam repository.ProvinceModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result bool, err errorModel.ErrorModel) {
	funcName := "IsExistProvince"

	query := fmt.Sprintf(`SELECT 
			CASE WHEN COUNT(p.id) > 0 THEN TRUE ELSE FALSE END 
			FROM %s p 
			WHERE p.id = $1 AND p.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.ProvinceDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
	}

	dbError := db.QueryRow(query, param...).Scan(&result)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceDAO) GetProvinceForInsert(db *sql.DB, userParam repository.ProvinceModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ProvinceModel, err errorModel.ErrorModel) {
	funcName := "GetProvinceForInsert"

	query := fmt.Sprintf(` SELECT 
			p.id, p.country_id, p.updated_at, 
			p.created_by 
			FROM %s p 
			WHERE 
			p.id = $1 AND p.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.ProvinceDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
	}

	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.CountryID,
		&result.UpdatedAt, &result.CreatedBy,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceDAO) GetMaxLastSync(db *sql.DB) (result repository.ProvinceModel, err errorModel.ErrorModel) {
	funcName := "GetMaxLastSync"

	query := " SELECT MAX(last_sync) " +
		" FROM " + input.TableName + " " +
		" WHERE deleted = FALSE "
	dbError := db.QueryRow(query).Scan(&result.LastSync)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceDAO) CheckIsProvinceExist(db *sql.DB, provinceModel repository.ProvinceModel) (resultAmount int64, err errorModel.ErrorModel) {
	funcName := "CheckIsProvinceExist"

	query := "SELECT COUNT(id) " +
		" FROM " + input.TableName + " " +
		" WHERE " +
		" id = $1 AND " +
		" deleted = FALSE "

	params := []interface{}{provinceModel.ID.Int64}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&resultAmount)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceDAO) GetListProvinceLocal(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	var (
		newQuery        string
		additionalWhere string
	)

	newQuery = fmt.Sprintf(`SELECT 
		province.id, province.country_id, province.code, 
		province.name, (array_agg(district.id)) 
		FROM %s `, input.TableName)

	input.setScopeOnDAO(scopeLimit, scopeDB, &additionalWhere)
	additionalWhere += fmt.Sprintf(` GROUP BY 
		province.id, province.country_id, province.code, 
		province.name `)
	input.setSearchByProvince(&searchBy)
	input.setCreatedByProvince(createdBy, &searchBy)

	getListData := getListJoinDataDAO{Table: input.TableName, Query: newQuery, AdditionalWhere: []string{additionalWhere}}
	input.setGetListJoinProvince(&getListData)
	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var temp repository.ListLocalProvinceModel
		errors := rows.Scan(
			&temp.ID, &temp.CountryID, &temp.Code,
			&temp.Name, &temp.DistrictID)

		return temp, errors
	}

	return getListData.GetListJoinDataAndFreeAdditionalWhere(db, userParam, searchBy, 0, mappingFunc)
}

func (input provinceDAO) setSearchByProvince(searchBy *[]in.SearchByParam) {
	temp := *searchBy
	for index := range temp {
		switch temp[index].SearchKey {
		case "country_id":
			temp[index].SearchKey = "province." + temp[index].SearchKey
		case "id":
			temp[index].SearchKey = "province." + temp[index].SearchKey
		case "mdb_province_id":
			temp[index].SearchKey = "province." + temp[index].SearchKey
		case "code":
			temp[index].SearchKey = "province." + temp[index].SearchKey
		case "name":
			temp[index].SearchKey = "province." + temp[index].SearchKey
		}
	}

	*searchBy = append(*searchBy, in.SearchByParam{
		SearchKey:      "province.status",
		SearchValue:    constanta.StatusActive,
		SearchOperator: "eq",
		SearchType:     constanta.Filter,
		DataType:       "char",
	})
}

func (input provinceDAO) setCreatedByProvince(createdBy int64, searchBy *[]in.SearchByParam) {
	if createdBy > 0 {
		*searchBy = append(*searchBy, in.SearchByParam{
			SearchKey:      "province.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}
}

func (input provinceDAO) setGetListJoinProvince(getListData *getListJoinDataDAO) {
	getListData.LeftJoin(DistrictDAO.TableName, "province.id", "district.province_id")
}

func (input provinceDAO) PrepareScopeInProvince(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, additionalWhere *string) {
	var dbParam []interface{}

	_, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.ProvinceDataScope})
	dbParam = append(dbParam, param...)

	if len(dbParam) > 0 {
		*additionalWhere = " " + scopeDB[constanta.ProvinceDataScope].View + " IN ("
	}

	for idx, valueScope := range dbParam {
		idScope := valueScope.(int64)
		if len(dbParam)-(idx+1) == 0 {
			*additionalWhere += strconv.Itoa(int(idScope)) + ")"
		} else {
			*additionalWhere += strconv.Itoa(int(idScope)) + ","
		}
	}
}

func (input provinceDAO) GetListScopeProvinceLocal(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	var dbParam []interface{}
	newQuery := "SELECT " +
		"	province.id, province.country_id, province.code, " +
		"	province.name " +
		"FROM " + input.TableName + " "

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.ProvinceDataScope})
	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, newQuery, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ListLocalProvinceModel
			dbError := rows.Scan(
				&temp.ID, &temp.CountryID, &temp.Code,
				&temp.Name,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "province.id"},
			Deleted:   FieldStatus{FieldName: "province.deleted"},
			CreatedBy: FieldStatus{FieldName: "province.created_by", Value: createdBy},
		})
}

func (input provinceDAO) setScopeOnDAO(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, additionalWhere *string) {
	var additionalWhereTemp string
	PrepareScopeOnDAO(scopeLimit, scopeDB, &additionalWhereTemp, 1, constanta.ProvinceDataScope, true)
	if additionalWhereTemp != "" {
		*additionalWhere = " AND " + additionalWhereTemp
	}
}

func (input provinceDAO) GetProvinceForCustomer(db *sql.DB, userParam repository.ProvinceModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ProvinceModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetProvinceForCustomer"
		query    string
		param    []interface{}
		dbError  error
	)

	query = fmt.Sprintf(`
		SELECT id, country_id, updated_at, 
		created_by, mdb_province_id 
		FROM %s 
		WHERE 
		id = $1 AND deleted = FALSE `,
		input.TableName)

	param = []interface{}{userParam.ID.Int64}
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.ProvinceDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		param = append(param, scopeParam...)
	}

	dbError = db.QueryRow(query, param...).Scan(
		&result.ID, &result.CountryID, &result.UpdatedAt,
		&result.CreatedBy, &result.MDBProvinceID,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceDAO) GetUpdatedMDBProvince(db *sql.DB, userParam []repository.ProvinceModel) (result []repository.ProvinceModel, err errorModel.ErrorModel) {
	var (
		param      []interface{}
		tempResult []interface{}
	)

	tempQuery, _ := ListRangeToInQueryWithStartIndex(len(userParam), 1)
	query := fmt.Sprintf(`SELECT 
		id, updated_at, last_sync, mdb_province_id 
		FROM %s WHERE mdb_province_id IN ( %s ) FOR UPDATE `,
		input.TableName, tempQuery)

	for _, model := range userParam {
		param = append(param, model.MDBProvinceID.Int64)
	}

	tempResult, err = GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.ProvinceModel
		dbErrorS := rows.Scan(
			&temp.ID, &temp.UpdatedAt, &temp.LastSync, &temp.MDBProvinceID)
		return temp, dbErrorS
	}, param)

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.ProvinceModel))
		}
	}

	return
}

func (input provinceDAO) GetProvinceLastSync(db *sql.DB) (result time.Time, err errorModel.ErrorModel) {
	var (
		funcName = "GetProvinceLastSync"
		query    string
	)

	query = fmt.Sprintf(
		`SELECT 
		CASE WHEN MAX(last_sync) IS NULL
		THEN '0001-01-01 00:00:00.000000'::timestamp ELSE
		MAX(last_sync) END FROM %s`,
		input.TableName)

	row := db.QueryRow(query)
	dbError := row.Scan(&result)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input provinceDAO) UpdateDataProvince(db *sql.Tx, userParam repository.ProvinceModel) (err errorModel.ErrorModel) {
	funcName := "UpdateDataProvince"

	query := fmt.Sprintf(
		`UPDATE %s SET 
		country_id = $1, code = $2, name = $3, 
		status = $4, updated_at = $5, updated_client = $6,
		updated_by = $7, mdb_province_id = $8, last_sync = $9
		WHERE id = $10 `, input.TableName)

	param := []interface{}{
		userParam.CountryID.Int64, userParam.Code.String, userParam.Name.String,
		userParam.Status.String, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.UpdatedBy.Int64, userParam.MDBProvinceID.Int64, userParam.LastSync.Time,
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

func (input provinceDAO) InsertBulkProvince(db *sql.Tx, userParam []repository.ProvinceModel) (output []int64, err errorModel.ErrorModel) {
	var (
		tempQuery string
		index     = 1
		paramLen  = 13
		params    []interface{}
		funcName  = "InsertBulkProvince"
	)
	query := fmt.Sprintf(`INSERT INTO %s 
		(
			id, mdb_province_id, country_id, 
			code, name, status, 
			created_by, created_at, created_client, 
			updated_by, updated_at, updated_client, 
			last_sync
		) VALUES `,
		input.TableName)

	tempQuery, index, params = ListValuesToInsertBulk(userParam, paramLen, index, func(inputVal interface{}) []interface{} {
		tempInputValue := inputVal.(repository.ProvinceModel)
		param := []interface{}{
			tempInputValue.MDBProvinceID.Int64, tempInputValue.MDBProvinceID.Int64, tempInputValue.CountryID.Int64,
			tempInputValue.Code.String, tempInputValue.Name.String, tempInputValue.Status.String,
			tempInputValue.CreatedBy.Int64, tempInputValue.CreatedAt.Time, tempInputValue.CreatedClient.String,
			tempInputValue.UpdatedBy.Int64, tempInputValue.UpdatedAt.Time, tempInputValue.UpdatedClient.String,
			tempInputValue.LastSync.Time,
		}
		return param
	})

	query += tempQuery
	query += " RETURNING id "
	rows, errorS := db.Query(query, params...)
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
			output = append(output, idTemp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceDAO) UpdateLastSyncResetProvince(db *sql.Tx, userParam repository.ProvinceModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateLastSyncResetProvince"
		query    = fmt.Sprintf(`UPDATE %s SET last_sync = $1`, input.TableName)
		params   []interface{}
	)

	if userParam.LastSync.Time.IsZero() {
		query = fmt.Sprintf(`UPDATE %s SET last_sync = NULL, status = 'N'`, input.TableName)
	}

	params = append(params, userParam.LastSync.Time)
	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	if userParam.LastSync.Time.IsZero() {
		_, dbError = stmt.Exec()
	} else {
		_, dbError = stmt.Exec(params...)
	}

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
