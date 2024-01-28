package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type postalCodeDAO struct {
	AbstractDAO
}

var PostalCodeDAO = postalCodeDAO{}.New()

func (input postalCodeDAO) New() (output postalCodeDAO) {
	output.FileName = "PostalCodeDAO.go"
	output.TableName = "postal_code"
	return
}

func (input postalCodeDAO) GetPostalCodeIDByMdbID(db *sql.DB, userParam repository.PostalCodeModel) (result int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetPostalCodeIDByMdbID"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT pc.id 
		FROM %s pc 
		LEFT JOIN %s uv ON pc.urban_village_id = uv.id 
		WHERE 
		pc.mdb_postal_code_id = $1 AND pc.deleted = FALSE AND uv.mdb_urban_village_id = $2 `,
		input.TableName, UrbanVillageDAO.TableName)

	params := []interface{}{userParam.MDBPostalCodeID.Int64, userParam.UrbanVillageID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input postalCodeDAO) IsExistPostalCodeForInsert(db *sql.DB, userParam repository.PostalCodeModel) (result bool, err errorModel.ErrorModel) {
	var (
		funcName = "IsExistPostalCodeForInsert"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
		CASE WHEN COUNT(pc.id) > 0 THEN TRUE ELSE FALSE END 
		FROM %s pc 
		LEFT JOIN %s uv 
		ON pc.urban_village_id = uv.id 
		WHERE pc.id = $1 AND pc.deleted = FALSE AND uv.id = $2 `,
		input.TableName, UrbanVillageDAO.TableName)

	param := []interface{}{userParam.ID.Int64, userParam.UrbanVillageID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeDAO) ViewPostalCode(db *sql.DB, userParam repository.PostalCodeModel) (result repository.PostalCodeModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewPostalCode"
		query    string
	)

	query = fmt.Sprintf(
		`SELECT 
		pc.id, pc.urban_village_id, uv.name AS urban_village_name, 
		pc.code, pc.status, pc.created_by, 
		pc.created_at, pc.updated_by, pc.updated_at 
		FROM %s pc 
		LEFT JOIN %s uv ON pc.urban_village_id = uv.id 
		WHERE pc.id = $1 AND pc.deleted = false `,
		input.TableName, UrbanVillageDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.UrbanVillageID,
		&result.UrbanVillageName, &result.Code,
		&result.Status, &result.CreatedBy,
		&result.CreatedAt, &result.UpdatedBy, &result.UpdatedAt,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input postalCodeDAO) GetCountPostalCode(db *sql.DB, searchByParam []in.SearchByParam, inputParam repository.PostalCodeModel) (result int, err errorModel.ErrorModel) {
	var (
		additionalWhere = ""
		params          []interface{}
		additionalQuery string
	)

	additionalQuery = fmt.Sprintf(`
		pc LEFT JOIN %s uv ON pc.urban_village_id = uv.id `,
		UrbanVillageDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "pc." + param.SearchKey
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, params,
		input.TableName+" "+additionalQuery, searchByParam,
		additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "pc.id"},
			Deleted:   FieldStatus{FieldName: "pc.deleted"},
			CreatedBy: FieldStatus{FieldName: "pc.created_by", Value: inputParam.CreatedBy.Int64},
		})
}

func (input postalCodeDAO) GetListPostalCode(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, inputStruct repository.PostalCodeModel) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere = ""
		params          []interface{}
		query           string
	)

	query = fmt.Sprintf(`
		SELECT 
		pc.id as id, pc.urban_village_id, pc.mdb_postal_code_id, 
		pc.code as code, pc.status, pc.created_by, 
		pc.updated_at as updated_at 
		FROM %s pc 
		LEFT JOIN %s uv ON pc.urban_village_id = uv.id `,
		input.TableName, UrbanVillageDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "pc." + param.SearchKey
	}

	//--- Search By Urban Village
	searchByParam = append(searchByParam, in.SearchByParam{
		SearchKey:      "pc.status",
		SearchValue:    constanta.StatusActive,
		SearchOperator: "eq",
		SearchType:     constanta.Filter,
		DataType:       "char",
	})

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.PostalCodeModel
			dbError := rows.Scan(
				&temp.ID, &temp.UrbanVillageID,
				&temp.MDBPostalCodeID, &temp.Code,
				&temp.Status, &temp.CreatedBy,
				&temp.UpdatedAt,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "pc.id"},
			Deleted:   FieldStatus{FieldName: "pc.deleted"},
			CreatedBy: FieldStatus{FieldName: "pc.created_by", Value: inputStruct.CreatedBy.Int64},
		})
}

func (input postalCodeDAO) GetPostalCodeWithUrbanVillageID(db *sql.DB, userParam repository.PostalCodeModel) (result repository.PostalCodeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetPostalCodeWithUrbanVillageID"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
		pc.id, pc.urban_village_id, pc.updated_at, pc.mdb_postal_code_id 
		FROM %s pc 
		LEFT JOIN %s uv ON pc.urban_village_id = uv.id 
		WHERE pc.id = $1 AND pc.deleted = FALSE AND uv.id = $2 `,
		input.TableName, UrbanVillageDAO.TableName)

	param := []interface{}{userParam.ID.Int64, userParam.UrbanVillageID.Int64}
	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.UrbanVillageID, &result.UpdatedAt, &result.MDBPostalCodeID,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeDAO) GetDateLastSyncPostalCode(db *sql.DB) (result repository.PostalCodeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetDateLastSyncPostalCode"
		query    string
		param    []interface{}
		dbError  error
	)

	query = fmt.Sprintf(`
		SELECT 
		CASE WHEN MAX(last_sync) IS NULL 
		THEN '0001-01-01 00:00:00.000000'::timestamp 
		ELSE MAX(last_sync) END max_last_sync 
		FROM %s `, input.TableName)
	dbError = db.QueryRow(query, param...).Scan(&result.LastSync)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeDAO) GetUpdatedDBPostalCode(db *sql.DB, userParam []repository.PostalCodeModel) (result []repository.PostalCodeModel, err errorModel.ErrorModel) {
	var (
		param      []interface{}
		tempResult []interface{}
		tempQuery  string
		query      string
	)

	tempQuery, _ = ListRangeToInQueryWithStartIndex(len(userParam), 1)
	query = fmt.Sprintf(`
		SELECT id, updated_at, last_sync, mdb_postal_code_id 
		FROM %s WHERE mdb_postal_code_id IN ( %s ) `,
		input.TableName, tempQuery)

	for _, itemID := range userParam {
		param = append(param, itemID.MDBPostalCodeID.Int64)
	}

	tempResult, err = GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var (
			temp     repository.PostalCodeModel
			dbErrorS error
		)
		dbErrorS = rows.Scan(&temp.ID, &temp.UpdatedAt, &temp.LastSync, &temp.MDBPostalCodeID)
		return temp, dbErrorS
	}, param)
	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.PostalCodeModel))
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeDAO) UpdateDataPostalCode(db *sql.Tx, userParam repository.PostalCodeModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateDataPostalCode"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		urban_village_id = $1, mdb_postal_code_id = $2, code = $3, 
		status = $4, last_sync = $5, updated_by = $6, 
		updated_client = $7, updated_at = $8 
		WHERE id = $9 `,
		input.TableName)

	param := []interface{}{
		userParam.UrbanVillageID.Int64, userParam.MDBPostalCodeID.Int64, userParam.Code.String,
		userParam.Status.String, userParam.LastSync.Time, userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String, userParam.UpdatedAt.Time, userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeDAO) InsertBulkPostalCode(db *sql.Tx, userParam []repository.PostalCodeModel) (output []int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertBulkPostalCode"
		index    = 1
		paramLen = 12
		params   []interface{}
		result   []interface{}
		query    string
	)

	query = fmt.Sprintf(`INSERT INTO %s 
		(id, urban_village_id, mdb_postal_code_id, 
		code, status, last_sync, 
		created_by, created_client, created_at, 
		updated_by, updated_client, updated_at) 
		VALUES`,
		input.TableName)

	query += CreateDollarParamInMultipleRowsDAO(len(userParam), paramLen, index, "id")
	for i := 0; i < len(userParam); i++ {
		params = append(params,
			userParam[i].MDBPostalCodeID.Int64, userParam[i].UrbanVillageID.Int64, userParam[i].MDBPostalCodeID.Int64,
			userParam[i].Code.String, userParam[i].Status.String, userParam[i].LastSync.Time,
			userParam[i].CreatedBy.Int64, userParam[i].CreatedClient.String, userParam[i].CreatedAt.Time,
			userParam[i].UpdatedBy.Int64, userParam[i].UpdatedClient.String, userParam[i].UpdatedAt.Time)
	}

	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	result, err = RowsCatchResult(rows, func(rws *sql.Rows) (idTemp interface{}, err errorModel.ErrorModel) {
		var (
			funcName = "resultRowsInput"
			errorS   error
			id       int64
		)

		errorS = rows.Scan(&id)
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}

		idTemp = id
		return
	})
	if err.Error != nil {
		return
	}

	for _, itemResult := range result {
		output = append(output, itemResult.(int64))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeDAO) UpdateLastSyncPostalCode(db *sql.Tx, lastSync time.Time) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateLastSyncPostalCode"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET last_sync = $1 `, input.TableName)
	param := []interface{}{lastSync}
	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeDAO) InsertPostalCode(db *sql.Tx, userParam repository.PostalCodeModel) (output int64, err errorModel.ErrorModel) {
	var (
		params   []interface{}
		funcName = "InsertPostalCode"
	)
	query := fmt.Sprintf(`INSERT INTO %s 
		(
		id, mdb_postal_code_id, code, 
		status, created_by, created_at, 
		created_client, updated_by, updated_at, 
		updated_client, last_sync, urban_village_id
		) VALUES ( 
		$1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, (SELECT id FROM %s WHERE mdb_urban_village_id = $12 LIMIT 1) 
		)`,
		input.TableName, UrbanVillageDAO.TableName)

	params = append(params,
		userParam.MDBPostalCodeID.Int64, userParam.MDBPostalCodeID.Int64, userParam.Code.String,
		userParam.Status.String, userParam.CreatedBy.Int64, userParam.CreatedAt.Time,
		userParam.CreatedClient.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.LastSync.Time, userParam.UrbanVillageID.Int64,
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

func (input postalCodeDAO) GetPostalCodeByID(db *sql.DB, postalModel repository.PostalCodeModel) (result repository.PostalCodeModel, err errorModel.ErrorModel) {
	funcName := "GetPostalCodeByID"
	query := `
		SELECT 
			d.id 
		FROM ` + input.TableName + ` d
		WHERE 
			d.deleted = FALSE AND d.id = $1   
	`

	params := []interface{}{postalModel.ID.Int64}

	errorS := db.QueryRow(query, params...).Scan(&result.ID)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}