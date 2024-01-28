package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
	"time"
)

type subDistrictDAO struct {
	AbstractDAO
}

var SubDistrictDAO = subDistrictDAO{}.New()

func (input subDistrictDAO) New() (output subDistrictDAO) {
	output.FileName = "SubDistrictDAO.go"
	output.TableName = "sub_district"
	return
}

func (input subDistrictDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "sd." + (*searchByParam)[i].SearchKey
	}

	//--- Search By Sub District
	*searchByParam = append(*searchByParam, in.SearchByParam{
		SearchKey:      "sd.status",
		SearchValue:    constanta.StatusActive,
		SearchOperator: "eq",
		SearchType:     constanta.Filter,
		DataType:       "char",
	})

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
		userParam.OrderBy = "sd." + userParam.OrderBy
		break
	}
}

func (input subDistrictDAO) getSubDistrictDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "sd.id"},
		Deleted:   FieldStatus{FieldName: "sd.deleted"},
		CreatedBy: FieldStatus{FieldName: "sd.created_by", Value: createdBy},
	}
}

func (input subDistrictDAO) GetSubDistrictIDByMdbID(db *sql.DB, userParam repository.SubDistrictModel) (result int64, err errorModel.ErrorModel) {
	funcName := "GetSubDistrictIDByMdbID"
	query := fmt.Sprintf(
		`SELECT
			sd.id
		FROM %s sd
		LEFT JOIN %s d ON sd.district_id = d.id
		WHERE
			sd.mdb_sub_district_id = $1 AND sd.deleted = FALSE AND
			d.mdb_district_id = $2 `,
		input.TableName, DistrictDAO.TableName)

	params := []interface{}{userParam.MDBSubDistrictID.Int64, userParam.DistrictID.Int64}

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

func (input subDistrictDAO) IsExistSubDistrictWithDistrictID(db *sql.DB, userParam repository.SubDistrictModel) (result bool, err errorModel.ErrorModel) {
	funcName := "IsExistSubDistrictWithDistrictID"
	query := fmt.Sprintf(
		`SELECT
			CASE WHEN COUNT(sd.id) > 0 THEN TRUE ELSE FALSE END
		FROM %s sd
		LEFT JOIN %s d ON sd.district_id = d.id
		WHERE
			sd.id = $1 AND sd.deleted = FALSE
			AND sd.district_id = $2 `,
		input.TableName, DistrictDAO.TableName)

	param := []interface{}{userParam.ID.Int64, userParam.DistrictID.Int64}

	dbError := db.QueryRow(query, param...).Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input subDistrictDAO) ViewSubDistrict(db *sql.DB, userParam repository.SubDistrictModel) (result repository.SubDistrictModel, err errorModel.ErrorModel) {
	funcName := "ViewSubDistrict"
	var tempResult interface{}

	query := fmt.Sprintf(
		`SELECT
			sd.id, sd.district_id, d.name AS district_name,
			sd.code, sd.name, sd.status, sd.created_by,
			sd.created_at, sd.updated_by, sd.updated_at
		FROM %s sd
		LEFT JOIN %s d ON sd.district_id = d.id
		WHERE
			sd.id = $1 AND sd.deleted = FALSE `,
		input.TableName, DistrictDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	row := db.QueryRow(query, params...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.SubDistrictModel
		dbError := rws.Scan(
			&temp.ID, &temp.DistrictID,
			&temp.DistrictName, &temp.Code,
			&temp.Name, &temp.Status, &temp.CreatedBy,
			&temp.CreatedAt, &temp.UpdatedBy, &temp.UpdatedAt,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.SubDistrictModel)
	}

	return
}

func (input subDistrictDAO) GetCountSubDistrict(db *sql.DB, searchByParam []in.SearchByParam, inputParam repository.SubDistrictModel) (result int, err errorModel.ErrorModel) {
	additionalWhere := ""
	var params []interface{}

	additionalQuery := fmt.Sprintf(
		` sd LEFT JOIN %s d ON sd.district_id = d.id `, DistrictDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "sd." + param.SearchKey
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, params,
		input.TableName+" "+additionalQuery, searchByParam,
		additionalWhere, input.getSubDistrictDefaultMustCheck(inputParam.CreatedBy.Int64))
}

func (input subDistrictDAO) GetListSubDistrict(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, inputStruct repository.SubDistrictModel) (result []interface{}, err errorModel.ErrorModel) {
	additionalWhere := ""
	var params []interface{}

	query := fmt.Sprintf(
		`SELECT
			sd.id, sd.district_id,
			sd.code, sd.name, sd.status,
			sd.created_by, sd.updated_at
		FROM %s sd
		LEFT JOIN %s d ON sd.district_id = d.id `,
		input.TableName, DistrictDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.SubDistrictModel
			dbError := rows.Scan(
				&temp.ID, &temp.DistrictID, &temp.Code,
				&temp.Name, &temp.Status, &temp.CreatedBy,
				&temp.UpdatedAt,
			)
			return temp, dbError
		}, additionalWhere, input.getSubDistrictDefaultMustCheck(inputStruct.CreatedBy.Int64))
}

func (input subDistrictDAO) GetSubDistrictByID(db *sql.DB, id int64, _ int64, isMustNotCheckDeleted bool) (result repository.SubDistrictModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetSubDistrictByID"
		tempResult interface{}
		query      string
	)

	query = fmt.Sprintf(`SELECT
			sd.id, sd.mdb_sub_district_id,
			sd.code, sd.name,
			sd.created_by, sd.updated_at
		FROM %s sd
		WHERE sd.mdb_sub_district_id = $1 `, input.TableName)

	if !isMustNotCheckDeleted {
		query += fmt.Sprintf(` AND sd.deleted = FALSE `)
	}

	params := []interface{}{id}
	row := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.SubDistrictModel
		dbError := rws.Scan(
			&temp.ID, &temp.MDBSubDistrictID, &temp.Code,
			&temp.Name, &temp.CreatedBy, &temp.UpdatedAt,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.SubDistrictModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input subDistrictDAO) GetSubDistrictWithDistrictID(db *sql.DB, userParam repository.SubDistrictModel) (result repository.SubDistrictModel, err errorModel.ErrorModel) {
	funcName := "GetSubDistrictWithDistrictID"
	query := fmt.Sprintf(
		`SELECT
			sd.id, sd.district_id, sd.updated_at, sd.mdb_sub_district_id
		FROM %s sd
		LEFT JOIN %s d ON sd.district_id = d.id
		WHERE
			sd.id = $1 AND sd.deleted = FALSE
			AND sd.district_id = $2 `,
		input.TableName, DistrictDAO.TableName)

	param := []interface{}{userParam.ID.Int64, userParam.DistrictID.Int64}

	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.DistrictID, &result.UpdatedAt, &result.MDBSubDistrictID,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input subDistrictDAO) GetSubDistrictLastSync(db *sql.DB) (result time.Time, err errorModel.ErrorModel) {
	funcName := "GetSubDistrictLastSync"
	query := fmt.Sprintf(
		`SELECT 
		CASE WHEN MAX(last_sync) IS NULL
		THEN '0001-01-01 00:00:00.000000'::timestamp ELSE
		MAX(last_sync) END 
	FROM %s`, input.TableName)

	row := db.QueryRow(query)

	dbError := row.Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input subDistrictDAO) GetUpdatedMDBSubDistrict(db *sql.DB, userParam []repository.SubDistrictModel) (result []repository.SubDistrictModel, err errorModel.ErrorModel) {
	var param []interface{}

	tempQuery, _ := ListRangeToInQueryWithStartIndex(len(userParam), 1)

	query := fmt.Sprintf(`SELECT 
		id, updated_at, last_sync, mdb_sub_district_id 
	FROM %s WHERE mdb_sub_district_id IN ( %s )  FOR UPDATE `,
		input.TableName, tempQuery)

	for _, model := range userParam {
		param = append(param, model.MDBSubDistrictID.Int64)
	}

	tempResult, err := GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.SubDistrictModel
		dbErrorS := rows.Scan(
			&temp.ID, &temp.UpdatedAt, &temp.LastSync, &temp.MDBSubDistrictID)
		return temp, dbErrorS
	}, param)

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.SubDistrictModel))
		}
	}

	return
}

func (input subDistrictDAO) UpdateDataSubDistrict(db *sql.Tx, userParam repository.SubDistrictModel) (err errorModel.ErrorModel) {
	funcName := "UpdateDataSubDistrict"

	query := fmt.Sprintf(
		`UPDATE %s 
	SET 
		district_id = (SELECT id FROM %s WHERE mdb_district_id = $1 LIMIT 1), code = $2, name = $3, 
		status = $4, updated_at = $5, updated_client = $6,
		updated_by = $7, mdb_sub_district_id = $8, last_sync = $9
	WHERE id = $10 `, input.TableName, DistrictDAO.TableName)

	param := []interface{}{
		userParam.DistrictID.Int64, userParam.Code.String, userParam.Name.String,
		userParam.Status.String, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.UpdatedBy.Int64, userParam.MDBSubDistrictID.Int64, userParam.LastSync.Time,
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

func (input subDistrictDAO) InsertBulkSubDistrict(db *sql.Tx, userParam []repository.SubDistrictModel) (output []int64, err errorModel.ErrorModel) {
	var (
		tempQuery string
		index     = 1
		paramLen  = 12
		params    []interface{}
		funcName  = "InsertBulkSubDistrict"
	)
	query := fmt.Sprintf(`INSERT INTO %s 
	(
		id, mdb_sub_district_id, code, name, status, created_by, 
		created_at, created_client, updated_by, updated_at, updated_client, 
		last_sync, district_id
	) VALUES `, input.TableName)

	for i, model := range userParam {
		query += " ( "
		tempQuery, index = ListRangeToInQueryWithStartIndex(paramLen, index)

		fkKeyQuery := fmt.Sprintf("(SELECT id FROM %s WHERE mdb_district_id = %d LIMIT 1) ",
			DistrictDAO.TableName, model.DistrictID.Int64)

		query += tempQuery
		query += fmt.Sprintf(", %s ) ", fkKeyQuery)

		if i < len(userParam)-1 {
			query += ", "
		}

		params = append(params,
			model.MDBSubDistrictID.Int64, model.MDBSubDistrictID.Int64, model.Code.String,
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

func (input subDistrictDAO) InsertSubDistrict(db *sql.Tx, userParam repository.SubDistrictModel) (output int64, err errorModel.ErrorModel) {
	var (
		params   []interface{}
		funcName = "InsertDistrict"
	)
	query := fmt.Sprintf(`INSERT INTO %s 
		(id, mdb_sub_district_id, code, 
		name, status, created_by, 
		created_at, created_client, updated_by, 
		updated_at, updated_client, last_sync, 
		district_id
		) VALUES 
		($1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12, 
		(SELECT id FROM %s WHERE mdb_district_id = $13 LIMIT 1)) `,
		input.TableName, DistrictDAO.TableName)

	params = append(params,
		userParam.MDBSubDistrictID.Int64, userParam.MDBSubDistrictID.Int64, userParam.Code.String,
		userParam.Name.String, userParam.Status.String, userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time, userParam.CreatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.LastSync.Time,
		userParam.DistrictID.Int64,
	)

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input subDistrictDAO) GetSubDistrictByIDForGetList(db *sql.DB, id int64, _ int64, isMustNotCheckDeleted bool) (result repository.SubDistrictModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetSubDistrictByID"
		tempResult interface{}
		query      string
	)

	query = fmt.Sprintf(`SELECT
			sd.id, sd.mdb_sub_district_id,
			sd.code, sd.name,
			sd.created_by, sd.updated_at
		FROM %s sd
		WHERE sd.id = $1 `, input.TableName)

	if !isMustNotCheckDeleted {
		query += fmt.Sprintf(` AND sd.deleted = FALSE `)
	}

	params := []interface{}{id}
	row := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.SubDistrictModel
		dbError := rws.Scan(
			&temp.ID, &temp.MDBSubDistrictID, &temp.Code,
			&temp.Name, &temp.CreatedBy, &temp.UpdatedAt,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.SubDistrictModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input subDistrictDAO) GetSubDistrictByIDForUpdate(db *sql.DB, districtModel repository.SubDistrictModel) (result repository.SubDistrictModel, err errorModel.ErrorModel) {
	funcName := "GetSubDistrictByIDForUpdate"
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