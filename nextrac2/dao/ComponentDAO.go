package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type componentDAO struct {
	AbstractDAO
}

var ComponentDAO = componentDAO{}.New()

func (input componentDAO) New() (output componentDAO) {
	output.FileName = "ComponentDAO.go"
	output.TableName = "component"
	return
}

func (input componentDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {

	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "comp." + (*searchByParam)[i].SearchKey
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
		userParam.OrderBy = "comp." + userParam.OrderBy
		break
	}

}

func (input componentDAO) DeleteComponent(db *sql.Tx, userParam repository.ComponentModel) (err errorModel.ErrorModel) {
	funcName := "DeleteComponent"

	query := fmt.Sprintf(`UPDATE %s
	SET
		deleted = $1, 
		component_name = $2, 
		updated_by = $3,
		updated_at = $4,
		updated_client = $5 
	WHERE
		id = $6`, input.TableName)
	param := []interface{}{
		true,
		userParam.ComponentName.String,
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

func (input componentDAO) UpdateComponent(db *sql.Tx, userParam repository.ComponentModel) (err errorModel.ErrorModel) {
	funcName := "UpdateComponent"
	query := fmt.Sprintf(`UPDATE %s
	SET 
		component_name = $1,
		updated_by = $2,
		updated_at = $3,
		updated_client = $4
	WHERE
		id = $5
	`, input.TableName)

	params := []interface{}{
		userParam.ComponentName.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input componentDAO) GetComponentForDelete(db *sql.Tx, userParam repository.ComponentModel) (result repository.ComponentModel, err errorModel.ErrorModel) {
	var tempResult interface{}
	funcName := "GetComponentForDelete"
	query := fmt.Sprintf(`SELECT 
		id, updated_at 
	FROM %s
	WHERE
		id = $1 AND deleted = FALSE 
	`, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE "

	rows := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ComponentModel
		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.ComponentModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input componentDAO) GetComponentForUpdate(db *sql.Tx, userParam repository.ComponentModel) (result repository.ComponentModel, err errorModel.ErrorModel) {
	var (
		tempResult interface{}
		funcName   = "GetComponentForUpdate"
	)

	query := fmt.Sprintf(`SELECT cp.id, cp.updated_at, cp.created_by, 
		CASE WHEN 
			(SELECT COUNT(id) FROM %s WHERE component_id = cp.id AND deleted = FALSE) > 0 
				OR 
			(SELECT COUNT(id) FROM %s WHERE component_id = cp.id AND deleted = FALSE) > 0 
		THEN TRUE ELSE FALSE END isUsed, cp.component_name 
		FROM %s cp
		WHERE cp.id = $1 AND cp.deleted = FALSE `,
		LicenseConfigProductComponentDAO.TableName, ProductComponentDAO.TableName, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND cp.created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE "
	rows := db.QueryRow(query, param...)
	if tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ComponentModel
		dbError := rws.Scan(
			&temp.ID, &temp.UpdatedAt, &temp.CreatedBy,
			&temp.IsUsed, &temp.ComponentName,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.ComponentModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input componentDAO) GetCountComponent(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input componentDAO) ViewComponent(db *sql.DB, userParam repository.ComponentModel) (result repository.ComponentModel, err errorModel.ErrorModel) {
	funcName := "ViewComponent"
	var tempResult interface{}

	query := fmt.Sprintf(`SELECT
		comp.id, comp.component_name, 
		comp.created_at, comp.updated_at, 
		comp.updated_by, u.nt_username, 
		comp.created_by 
	FROM %s comp
	LEFT JOIN "%s" AS u ON comp.updated_by = u.id 
	WHERE 
		comp.id = $1 AND comp.deleted = FALSE 
	`, input.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND comp.created_by = $2 "
		params = append(params, userParam.CreatedBy.Int64)
	}

	rows := db.QueryRow(query, params...)

	if tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ComponentModel
		dbError := rws.Scan(
			&temp.ID, &temp.ComponentName,
			&temp.CreatedAt, &temp.UpdatedAt,
			&temp.UpdatedBy, &temp.UpdatedName,
			&temp.CreatedBy,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.ComponentModel)
	}

	return
}

func (input componentDAO) GetLisComponent(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(`SELECT 
		comp.id, comp.component_name, 
		comp.created_at, comp.updated_at, 
		comp.updated_by, u.nt_username 
	FROM %s comp
	LEFT JOIN "%s" AS u ON comp.updated_by = u.id 
	`, input.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ComponentModel
			dbError := rows.Scan(
				&temp.ID, &temp.ComponentName,
				&temp.CreatedAt, &temp.UpdatedAt,
				&temp.UpdatedBy, &temp.UpdatedName,
			)
			return temp, dbError
		}, " ", DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "comp.id"},
			Deleted:   FieldStatus{FieldName: "comp.deleted"},
			CreatedBy: FieldStatus{FieldName: "comp.created_by", Value: createdBy},
		})
}

func (input componentDAO) InsertComponent(db *sql.Tx, userParam repository.ComponentModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertComponent"
		query    string
	)

	query = fmt.Sprintf(`INSERT INTO %s (component_name, created_by, created_at, 
		created_client, updated_by, updated_at, 
		updated_client) VALUES 
		($1, $2, $3, 
		$4, $5, $6, 
		$7) RETURNING id`,
		input.TableName)

	params := []interface{}{userParam.ComponentName.String, userParam.CreatedBy.Int64, userParam.CreatedAt.Time,
		userParam.CreatedClient.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
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

func (input componentDAO) CheckComponentIsExist(db *sql.DB, userParam repository.ComponentModel) (result repository.ComponentModel, err errorModel.ErrorModel) {
	funcName := "CheckComponentIsExist"

	query := fmt.Sprintf(`SELECT id 
	FROM %s
	WHERE id = $1 AND deleted = FALSE
	`, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	dbError := db.QueryRow(query, param...).Scan(&result.ID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
