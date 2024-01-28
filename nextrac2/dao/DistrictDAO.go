package dao

import (
	"database/sql"
	"fmt"
	"time"

	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type districtDAO struct {
	AbstractDAO
}

var DistrictDAO = districtDAO{}.New()

func (input districtDAO) New() (output districtDAO) {
	output.FileName = "DistrictDAO.go"
	output.TableName = "district"
	return
}

func (input districtDAO) GetDistrictIDByMdbID(db *sql.DB, userParam repository.DistrictModel) (result int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetDistrictIDByMdbID"
		query    string
	)

	query = fmt.Sprintf(`SELECT d.id
		FROM %s d 
		LEFT JOIN %s p ON d.province_id = p.id
		WHERE 
		d.mdb_district_id = $1 AND d.deleted = FALSE AND p.mdb_province_id = $2 `,
		input.TableName, ProvinceDAO.TableName)

	params := []interface{}{userParam.MdbDistrictID.Int64, userParam.ProvinceID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) GetOnlyDistrictIDByMdbID(db *sql.DB, userParam repository.DistrictModel, isMustNotCheckDeleted bool) (result int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetOnlyDistrictIDByMdbID"
		query    string
	)

	query = fmt.Sprintf(`SELECT id FROM %s WHERE mdb_district_id = $1 `, input.TableName)
	if !isMustNotCheckDeleted {
		query += fmt.Sprintf(` AND deleted = FALSE `)
	}

	params := []interface{}{userParam.MdbDistrictID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) IsExistDistrictWithProvinceID(db *sql.DB, districtModel repository.ListLocalDistrictModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result bool, err errorModel.ErrorModel) {
	var (
		funcName = "IsExistDistrictWithProvinceID"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
		CASE WHEN COUNT(d.id) > 0 THEN TRUE ELSE FALSE END 
		FROM %s d 
		LEFT JOIN %s p 
		ON d.province_id = p.id
		WHERE d.deleted = FALSE AND d.id = $1 AND p.id = $2 `,
		input.TableName, ProvinceDAO.TableName)

	params := []interface{}{districtModel.ID.Int64, districtModel.ProvinceID.Int64}

	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 3, []string{constanta.ProvinceDataScope, constanta.DistrictDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		params = append(params, scopeParam...)
	}

	errorS := db.QueryRow(query, params...).Scan(&result)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) GetDistrictByID(db *sql.DB, districtModel repository.DistrictModel) (result repository.DistrictModel, err errorModel.ErrorModel) {
	funcName := "GetDistrictByID"
	query := `
		SELECT 
			d.id 
		FROM ` + input.TableName + ` d
		WHERE 
			d.deleted = FALSE AND d.id = $1   
	`

	params := []interface{}{districtModel.ID.Int64}

	errorS := db.QueryRow(query, params...).Scan(&result.ID)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) CheckIsDistrictExist(db *sql.DB, districtModel repository.DistrictModel) (resultAmount int64, err errorModel.ErrorModel) {
	funcName := "CheckIsDistrictExist"

	query := `
		SELECT 
			COUNT(id)  
		FROM ` + input.TableName + ` 
		WHERE 
			id = $1 AND deleted = FALSE  
	`

	params := []interface{}{districtModel.ID.Int64}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&resultAmount)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) GetDistrictLastSync(db *sql.DB) (result time.Time, err errorModel.ErrorModel) {
	funcName := "GetDistrictLastSync"
	query := fmt.Sprintf(`SELECT 
		CASE WHEN MAX(last_sync) IS NULL
		THEN '0001-01-01 00:00:00.000000'::timestamp ELSE
		MAX(last_sync) END FROM %s`, input.TableName)

	results := db.QueryRow(query)

	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input districtDAO) UpdateDataDistrict(db *sql.Tx, userParam repository.DistrictModel) (err errorModel.ErrorModel) {
	funcName := "UpdateDataDistrict"

	query := fmt.Sprintf(
		`UPDATE %s 
	SET 
		province_id = (SELECT id FROM %s WHERE mdb_province_id = $1 LIMIT 1), code = $2, name = $3, 
		status = $4, updated_at = $5, updated_client = $6,
		updated_by = $7, mdb_district_id = $8, last_sync = $9
	WHERE id = $10 `, input.TableName, ProvinceDAO.TableName)

	param := []interface{}{
		userParam.ProvinceID.Int64, userParam.Code.String, userParam.Name.String,
		userParam.Status.String, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.UpdatedBy.Int64, userParam.MdbDistrictID.Int64, userParam.LastSync.Time,
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

func (input districtDAO) GetUpdatedMDBDistrict(db *sql.DB, userParam []repository.DistrictModel) (result []repository.DistrictModel, err errorModel.ErrorModel) {
	var param []interface{}

	tempQuery, _ := ListRangeToInQueryWithStartIndex(len(userParam), 1)

	query := fmt.Sprintf(`SELECT 
		id, updated_at, last_sync, mdb_district_id 
	FROM %s WHERE mdb_district_id IN ( %s ) `,
		input.TableName, tempQuery)

	for _, model := range userParam {
		param = append(param, model.MdbDistrictID.Int64)
	}

	tempResult, err := GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.DistrictModel
		dbErrorS := rows.Scan(
			&temp.ID, &temp.UpdatedAt, &temp.LastSync, &temp.MdbDistrictID)
		return temp, dbErrorS
	}, param)

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.DistrictModel))
		}
	}

	return
}

func (input districtDAO) GetListDistrictLocal(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	var (
		dbParam []interface{}
		query   string
	)

	query = fmt.Sprintf(`SELECT 
		d.id, d.province_id, d.code, 
		d.name 
		FROM %s d `,
		input.TableName)

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1, []string{constanta.DistrictDataScope})
	dbParam = append(dbParam, param...)

	//--- Search By District
	searchBy = append(searchBy, in.SearchByParam{
		SearchKey:      "d.status",
		SearchValue:    constanta.StatusActive,
		SearchOperator: "eq",
		SearchType:     constanta.Filter,
		DataType:       "char",
	})

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ListLocalDistrictModel

			errors := rows.Scan(
				&temp.ID, &temp.ProvinceID, &temp.Code,
				&temp.Name)

			return temp, errors

		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "d.id"},
			Deleted:   FieldStatus{FieldName: "d.deleted"},
			CreatedBy: FieldStatus{FieldName: "d.created_by", Value: createdBy},
		})
}

func (input districtDAO) GetCountProvinceOnDistrict(db *sql.DB, arrDistrictID []int) (result int, err errorModel.ErrorModel) {
	funcName := "GetCountProvinceOnDistrict"
	query := `
		SELECT 
			COUNT(p.id)
		FROM ` + input.TableName + ` d
		LEFT JOIN ` + ProvinceDAO.TableName + ` p ON d.province_id = p.id
		WHERE d.deleted = FALSE AND p.deleted = FALSE
	`

	query += " AND d.id IN ( "
	tempQuery, _ := ListRangeToInQueryWithStartIndex(len(arrDistrictID), 1)
	query += tempQuery + " ) "

	params := []interface{}{arrDistrictID}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) GetDistrictWithProvinceID(db *sql.DB, districtModel repository.ListLocalDistrictModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.DistrictModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetDistrictWithProvinceID"
		query    string
		params   []interface{}
		errorS   error
	)

	fmt.Println(fmt.Sprintf("%s : Validate District Local", funcName))

	query = fmt.Sprintf(`SELECT d.id, d.province_id, d.updated_at, 
		d.mdb_district_id 
		FROM %s d 
		LEFT JOIN %s p ON d.province_id = p.id
		WHERE 
		d.deleted = FALSE AND d.id = $1 AND p.id = $2 `,
		input.TableName, ProvinceDAO.TableName)

	params = []interface{}{districtModel.ID.Int64, districtModel.ProvinceID.Int64}
	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 3, []string{constanta.ProvinceDataScope, constanta.DistrictDataScope})
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		params = append(params, scopeParam...)
	}

	fmt.Println(fmt.Sprintf("%s : Data Params = %v", funcName, params))

	errorS = db.QueryRow(query, params...).Scan(
		&result.ID, &result.ProvinceID, &result.UpdatedAt,
		&result.MdbDistrictID,
	)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) InsertBulkDistrict(db *sql.Tx, userParam []repository.DistrictModel) (output []int64, err errorModel.ErrorModel) {
	var (
		index    = 1
		paramLen = 12
		params   []interface{}
		funcName = "InsertBulkDistrict"
		result   []interface{}
	)

	query := fmt.Sprintf(`INSERT INTO %s 
	(
		id, mdb_district_id, code, name, status, created_by, 
		created_at, created_client, updated_by, updated_at, updated_client, 
		last_sync, province_id
	) VALUES `, input.TableName)

	for i, model := range userParam {
		query += " ( "
		tempQuery, _ := ListRangeToInQueryWithStartIndex(paramLen, index)

		fkKeyQuery := fmt.Sprintf("(SELECT id FROM %s WHERE mdb_province_id = %d LIMIT 1) ",
			ProvinceDAO.TableName, model.ProvinceID.Int64)

		query += tempQuery
		query = fmt.Sprintf(", %s ) ", fkKeyQuery)

		if i < len(userParam)-1 {
			query += ", "
		}

		params = append(params,
			model.MdbDistrictID.Int64, model.MdbDistrictID.Int64, model.Code.String,
			model.Name.String, model.Status.String, model.CreatedBy.Int64,
			model.CreatedAt.Time, model.CreatedClient.String, model.UpdatedBy.Int64,
			model.UpdatedAt.Time, model.UpdatedClient.String, model.LastSync.Time,
		)
	}

	query += " RETURNING id "

	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	result, err = RowsCatchResult(rows, func(rws *sql.Rows) (idTemp interface{}, err errorModel.ErrorModel) {
		var (
			errorS error
			id     int64
		)

		errorS = rows.Scan(&id)
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}

		idTemp = id
		return
	})

	for _, itemResult := range result {
		output = append(output, itemResult.(int64))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input districtDAO) InsertDistrict(db *sql.Tx, userParam repository.DistrictModel) (output int64, err errorModel.ErrorModel) {
	var (
		params   []interface{}
		query    string
		funcName = "InsertDistrict"
	)

	query = fmt.Sprintf(`
		INSERT INTO %s 
		(
			id, mdb_district_id, code, name, status, created_by, 
			created_at, created_client, updated_by, updated_at, updated_client, 
			last_sync, province_id
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, 
		(SELECT id FROM %s WHERE mdb_province_id = $13 LIMIT 1)) `,
		input.TableName, ProvinceDAO.TableName)

	params = append(params,
		userParam.MdbDistrictID.Int64, userParam.MdbDistrictID.Int64, userParam.Code.String,
		userParam.Name.String, userParam.Status.String, userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time, userParam.CreatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.LastSync.Time,
		userParam.ProvinceID.Int64,
	)

	query += " RETURNING id "
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&output)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
