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

type urbanVillageDAO struct {
	AbstractDAO
}

var UrbanVillageDAO = urbanVillageDAO{}.New()

func (input urbanVillageDAO) New() (output urbanVillageDAO) {
	output.FileName = "UrbanVillageDAO.go"
	output.TableName = "urban_village"
	return
}

func (input urbanVillageDAO) GetUrbanVillageIDByMdbID(db *sql.DB, userParam repository.UrbanVillageModel) (result int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetUrbanVillageIDByMdbID"
		query    string
	)

	query = fmt.Sprintf(
		`SELECT uv.id FROM %s uv 
		LEFT JOIN %s sd ON uv.sub_district_id = sd.id 
		WHERE 
		uv.mdb_urban_village_id = $1 AND uv.deleted = FALSE AND sd.mdb_sub_district_id = $2 `,
		input.TableName, SubDistrictDAO.TableName)

	params := []interface{}{userParam.MDBUrbanVillageID.Int64, userParam.SubDistrictID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageDAO) ViewUrbanVillage(db *sql.DB, userParam repository.UrbanVillageModel) (result repository.UrbanVillageModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewUrbanVillage"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT uv.id, uv.sub_district_id, sd.name AS sub_district_name, 
		uv.code, uv.name, uv.status, 
		uv.created_by, uv.created_at, uv.updated_by, 
		uv.updated_at 
		FROM %s uv 
		LEFT JOIN %s sd ON uv.sub_district_id = sd.id 
		WHERE 
		uv.id = $1 AND uv.deleted = false `,
		input.TableName, SubDistrictDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.SubDistrictID,
		&result.SubDistrictName, &result.Code,
		&result.Name, &result.Status, &result.CreatedBy,
		&result.CreatedAt, &result.UpdatedBy, &result.UpdatedAt,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input urbanVillageDAO) GetCountUrbanVillage(db *sql.DB, searchByParam []in.SearchByParam, inputParam repository.UrbanVillageModel) (result int, err errorModel.ErrorModel) {
	var (
		additionalWhere = ""
		params          []interface{}
		additionalQuery string
	)

	additionalQuery = fmt.Sprintf(` uv LEFT JOIN %s sd ON uv.sub_district_id = sd.id `, SubDistrictDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "uv." + param.SearchKey
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, params,
		input.TableName+" "+additionalQuery, searchByParam,
		additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "uv.id"},
			Deleted:   FieldStatus{FieldName: "uv.deleted"},
			CreatedBy: FieldStatus{FieldName: "uv.created_by", Value: inputParam.CreatedBy.Int64},
		})
}

func (input urbanVillageDAO) GetListUrbanVillage(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, inputStruct repository.UrbanVillageModel) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere = ""
		params          []interface{}
		query           string
	)

	query = fmt.Sprintf(`
		SELECT uv.id as id, uv.sub_district_id, uv.mdb_urban_village_id, 
		uv.code as code, uv.name as name, uv.status, 
		uv.created_by, uv.updated_at as updated_at 
		FROM %s uv 
		LEFT JOIN %s sd ON uv.sub_district_id = sd.id `,
		input.TableName, SubDistrictDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "uv." + param.SearchKey
	}

	//--- Search By Urban Village
	searchByParam = append(searchByParam, in.SearchByParam{
		SearchKey:      "uv.status",
		SearchValue:    constanta.StatusActive,
		SearchOperator: "eq",
		SearchType:     constanta.Filter,
		DataType:       "char",
	})

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.UrbanVillageModel
			dbError := rows.Scan(
				&temp.ID, &temp.SubDistrictID,
				&temp.MDBUrbanVillageID, &temp.Code,
				&temp.Name, &temp.Status, &temp.CreatedBy,
				&temp.UpdatedAt,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "uv.id"},
			Deleted:   FieldStatus{FieldName: "uv.deleted"},
			CreatedBy: FieldStatus{FieldName: "uv.created_by", Value: inputStruct.CreatedBy.Int64},
		})
}

func (input urbanVillageDAO) GetUrbanVillageByID(db *sql.DB, id int64, isMustNotCheckDeleted bool) (result repository.UrbanVillageModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUrbanVillageByID"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
		uv.id, uv.mdb_urban_village_id, uv.code, 
		uv.name, uv.created_by, uv.updated_at 
		FROM %s uv 
		WHERE uv.mdb_urban_village_id = $1 `,
		input.TableName)

	if !isMustNotCheckDeleted {
		query += fmt.Sprintf(` AND uv.deleted = FALSE `)
	}

	params := []interface{}{id}
	errorS := db.QueryRow(query, params...).Scan(
		&result.ID, &result.MDBUrbanVillageID, &result.Code,
		&result.Name, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageDAO) IsExistUrbanVillageForInsert(db *sql.DB, userParam repository.UrbanVillageModel) (result bool, err errorModel.ErrorModel) {
	var (
		funcName = "IsExistPostalCodeForInsert"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT CASE WHEN COUNT(uv.id) > 0 
		THEN TRUE ELSE FALSE END 
		FROM %s uv 
		LEFT JOIN %s sd ON uv.sub_district_id = sd.id 
		WHERE uv.id = $1 AND uv.deleted = FALSE AND sd.id = $2 `,
		input.TableName, SubDistrictDAO.TableName)

	param := []interface{}{userParam.ID.Int64, userParam.SubDistrictID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageDAO) GetUrbanVillageWithSubDistrictID(db *sql.DB, userParam repository.UrbanVillageModel) (result repository.UrbanVillageModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUrbanVillageWithSubDistrictID"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT uv.id, uv.sub_district_id, uv.updated_at, 
		uv.mdb_urban_village_id 
		FROM %s uv 
		LEFT JOIN %s sd ON uv.sub_district_id = sd.id 
		WHERE uv.id = $1 AND uv.deleted = FALSE AND sd.id = $2 `,
		input.TableName, SubDistrictDAO.TableName)

	param := []interface{}{userParam.ID.Int64, userParam.SubDistrictID.Int64}
	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.SubDistrictID, &result.UpdatedAt, &result.MDBUrbanVillageID,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageDAO) GetDateLastSyncUrbanVillage(db *sql.DB) (result repository.UrbanVillageModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetDateLastSyncUrbanVillage"
		query    string
		param    []interface{}
		dbError  error
	)

	query = fmt.Sprintf(`
		SELECT 
		CASE WHEN MAX(last_sync) IS NULL 
		THEN '0001-01-01 00:00:00.000000'::timestamp ELSE 
		MAX(last_sync) END max_last_sync
		FROM %s `, input.TableName)
	dbError = db.QueryRow(query, param...).Scan(&result.LastSync)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageDAO) GetUpdatedDBUrbanVillage(db *sql.DB, userParam []repository.UrbanVillageModel) (result []repository.UrbanVillageModel, err errorModel.ErrorModel) {
	var (
		param      []interface{}
		tempResult []interface{}
		tempQuery  string
		query      string
	)

	tempQuery, _ = ListRangeToInQueryWithStartIndex(len(userParam), 1)
	query = fmt.Sprintf(`
		SELECT id, updated_at, last_sync, mdb_urban_village_id 
		FROM %s WHERE mdb_urban_village_id IN ( %s ) `,
		input.TableName, tempQuery)

	for _, itemID := range userParam {
		param = append(param, itemID.MDBUrbanVillageID.Int64)
	}

	tempResult, err = GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var (
			temp     repository.UrbanVillageModel
			dbErrorS error
		)
		dbErrorS = rows.Scan(&temp.ID, &temp.UpdatedAt, &temp.LastSync, &temp.MDBUrbanVillageID)
		return temp, dbErrorS
	}, param)
	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.UrbanVillageModel))
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageDAO) UpdateDataUrbanVillage(db *sql.Tx, userParam repository.UrbanVillageModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateDataUrbanVillage"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		sub_district_id = $1, mdb_urban_village_id = $2, code = $3, 
		status = $4, last_sync = $5, updated_by = $6, 
		updated_client = $7, updated_at = $8, name = $9  
		WHERE id = $10 `,
		input.TableName)

	param := []interface{}{
		userParam.SubDistrictID.Int64, userParam.MDBUrbanVillageID.Int64, userParam.Code.String,
		userParam.Status.String, userParam.LastSync.Time, userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String, userParam.UpdatedAt.Time, userParam.Name.String,
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

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageDAO) InsertBulkUrbanVillage(db *sql.Tx, userParam []repository.UrbanVillageModel) (output []int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertBulkUrbanVillage"
		index    = 1
		paramLen = 13
		params   []interface{}
		result   []interface{}
		query    string
	)

	query = fmt.Sprintf(`INSERT INTO %s 
		(id, sub_district_id, mdb_urban_village_id, 
		code, status, last_sync, 
		created_by, created_client, created_at, 
		updated_by, updated_client, updated_at, 
		name) 
		VALUES`,
		input.TableName)

	query += CreateDollarParamInMultipleRowsDAO(len(userParam), paramLen, index, "id")
	for i := 0; i < len(userParam); i++ {
		params = append(params,
			userParam[i].MDBUrbanVillageID.Int64, userParam[i].SubDistrictID.Int64, userParam[i].MDBUrbanVillageID.Int64,
			userParam[i].Code.String, userParam[i].Status.String, userParam[i].LastSync.Time,
			userParam[i].CreatedBy.Int64, userParam[i].CreatedClient.String, userParam[i].CreatedAt.Time,
			userParam[i].UpdatedBy.Int64, userParam[i].UpdatedClient.String, userParam[i].UpdatedAt.Time,
			userParam[i].Name.String)
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

func (input urbanVillageDAO) UpdateLastSyncUrbanVillage(db *sql.Tx, lastSync time.Time) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateLastSyncUrbanVillage"
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

func (input urbanVillageDAO) InsertUrbanVillage(db *sql.Tx, userParam repository.UrbanVillageModel) (output int64, err errorModel.ErrorModel) {
	var (
		params   []interface{}
		funcName = "InsertUrbanVillage"
	)

	query := fmt.Sprintf(`INSERT INTO %s 
		(
		id, mdb_urban_village_id, code, 
		name, status, created_by, 
		created_at, created_client, updated_by, 
		updated_at, updated_client, last_sync, 
		sub_district_id
		) VALUES ( 
		$1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12, 
		(SELECT id FROM %s WHERE mdb_sub_district_id = $13 LIMIT 1) 
		)`, input.TableName, SubDistrictDAO.TableName)

	params = append(params,
		userParam.MDBUrbanVillageID.Int64, userParam.MDBUrbanVillageID.Int64, userParam.Code.String,
		userParam.Name.String, userParam.Status.String, userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time, userParam.CreatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.LastSync.Time,
		userParam.SubDistrictID.Int64,
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

func (input urbanVillageDAO) GetUrbanVillageByIDForGetList(db *sql.DB, id int64, isMustNotCheckDeleted bool) (result repository.UrbanVillageModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUrbanVillageByID"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
		uv.id, uv.mdb_urban_village_id, uv.code, 
		uv.name, uv.created_by, uv.updated_at 
		FROM %s uv 
		WHERE uv.id = $1 `,
		input.TableName)

	if !isMustNotCheckDeleted {
		query += fmt.Sprintf(` AND uv.deleted = FALSE `)
	}

	params := []interface{}{id}
	errorS := db.QueryRow(query, params...).Scan(
		&result.ID, &result.MDBUrbanVillageID, &result.Code,
		&result.Name, &result.CreatedBy, &result.UpdatedAt)
	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
