package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type whiteListDeviceDAO struct {
	AbstractDAO
}

var WhiteListDevice = whiteListDeviceDAO{}.New()

func (input whiteListDeviceDAO) New() (output whiteListDeviceDAO) {
	output.FileName = "WhiteListDevice.go"
	output.TableName = "device"
	return
}

func (input whiteListDeviceDAO) InsertWhiteListDevice(db *sql.Tx, userParam repository.WhiteListDeviceModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertWhiteListDevice"
	)

	query := fmt.Sprintf(
		`INSERT INTO %s
		(
			name, description,
			created_by, created_at, created_client,
			updated_by, updated_at, updated_client
		)
		VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 )
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.Device.String, userParam.Description.String,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input whiteListDeviceDAO) GetWhiteListDeviceForUpdateOrDelete(db *sql.Tx, userParam repository.WhiteListDeviceModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.WhiteListDeviceModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetWhiteListDeviceForUpdateOrDelete"
		tempResult interface{}
		//index := 1
	)

	query := fmt.Sprintf(
		`SELECT
			d.id, d.updated_at, d.created_at, d.created_by
		FROM %s d
		WHERE
			d.id = $1 AND d.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	//index += len(param)

	// Add Data scope
	//scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, []string{constanta.ProductGroupDataScope})
	//if scopeAdditionalWhere != "" {
	//	query += " " + scopeAdditionalWhere
	//	param = append(param, scopeParam...)
	//	index += len(scopeParam)
	//}

	// Check own access
	//_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getProductGroupDefaultMustCheck, index)

	query += " FOR UPDATE "

	row := db.QueryRow(query, param...)

	tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.WhiteListDeviceModel
		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt, &temp.CreatedAt, &temp.CreatedBy,
		)
		return temp, dbError
	}, input.FileName, funcName)

	if tempResult != nil {
		result = tempResult.(repository.WhiteListDeviceModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input whiteListDeviceDAO) UpdateWhiteListDevice(db *sql.Tx, userParam repository.WhiteListDeviceModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateWhiteListDevice"
	)

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			name = $1,
			description = $2,
			updated_by = $3,
			updated_at = $4,
			updated_client = $5
		WHERE
			id = $6 `, input.TableName)
	param := []interface{}{
		userParam.Device.String, userParam.Description.String, userParam.UpdatedBy.Int64,
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

func (input whiteListDeviceDAO) DeleteWhiteListDevice(db *sql.Tx, userParam repository.WhiteListDeviceModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteWhiteListDevice"
	)

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			deleted = $1, updated_by = $2, updated_at = $3, updated_client = $4
		WHERE
			id = $5 `, input.TableName)

	param := []interface{}{
		true,
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

	return
}

func (input whiteListDeviceDAO) HardDeleteWhiteListDevice(db *sql.Tx, userParam repository.WhiteListDeviceModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteWhiteListDevice"
	)

	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1 `, input.TableName)

	param := []interface{}{
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

func (input whiteListDeviceDAO) GetCountWhiteListDevice(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	for i, _ := range searchByParam {
		if searchByParam[i].SearchKey == "device" {
			searchByParam[i].SearchKey = "name"
		}
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input whiteListDeviceDAO) GetListWhiteListDevice(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	var query string
	query = fmt.Sprintf(`SELECT wd.id, wd.name, wd.description, 
		wd.created_at, wd.updated_at, u.nt_username
		FROM %s wd
		LEFT JOIN "%s" AS u 
		ON wd.updated_by = u.id `,
		input.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.WhiteListDeviceModel
			dbError := rows.Scan(
				&temp.ID, &temp.Device, &temp.Description, &temp.CreatedAt,
				&temp.UpdatedAt, &temp.UpdatedName,
			)
			return temp, dbError
		}, " ", input.getDefaultMustCheck(createdBy))
}

func (input whiteListDeviceDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "wd.id"},
		Deleted:   FieldStatus{FieldName: "wd.deleted"},
		CreatedBy: FieldStatus{FieldName: "wd.created_by", Value: createdBy},
	}
}

func (input whiteListDeviceDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		if (*searchByParam)[i].SearchKey == "device" {
			(*searchByParam)[i].SearchKey = "wd.name"
		} else {
			(*searchByParam)[i].SearchKey = "wd." + (*searchByParam)[i].SearchKey
		}
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
	case "device", "device ASC", "device DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "wd.name " + strSplit[1]
		} else {
			userParam.OrderBy = "wd.name"
		}
		break
	default:
		userParam.OrderBy = "wd." + userParam.OrderBy
		break
	}
}

func (input whiteListDeviceDAO) ViewWhiteListDevice(db *sql.DB, userParam repository.WhiteListDeviceModel) (result repository.WhiteListDeviceModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewWhiteListDevice"
		query    string
		params   []interface{}
	)

	query = fmt.Sprintf(`SELECT 
		wd.id, wd.name, wd.description, 
		wd.created_at, wd.updated_at, wd.updated_by, 
		u.nt_username, wd.created_by
		FROM %s wd
		LEFT JOIN "%s" AS u ON wd.updated_by = u.id
		WHERE 
		wd.id = $1 AND wd.deleted = FALSE `,
		input.TableName, UserDAO.TableName)

	params = []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.Device, &result.Description,
		&result.CreatedAt, &result.UpdatedAt, &result.UpdatedBy,
		&result.UpdatedName, &result.CreatedBy,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input whiteListDeviceDAO) GetValidateWhiteListDevice(db *sql.DB, userParam repository.WhiteListDeviceModel) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName = "GetValidateWhiteListDevice"
		query    string
		params   []interface{}
		results  *sql.Row
	)

	query = fmt.Sprintf(`SELECT CASE WHEN wd.id > 0 
		THEN true ELSE false 
		END is_exist
		FROM %s wd
		WHERE 
		wd.name = $1 AND wd.deleted = FALSE `,
		input.TableName)

	params = []interface{}{userParam.Device.String}
	results = db.QueryRow(query, params...)

	dbError := results.Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}