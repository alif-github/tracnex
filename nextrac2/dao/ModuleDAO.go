package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type ModuleDAOInterface interface {
	GetModuleForUpdate(*sql.Tx, repository.ModuleModel) (repository.ModuleModel, errorModel.ErrorModel)
}

type moduleDAO struct {
	AbstractDAO
}

var ModuleDAO = moduleDAO{}.New()

func (input moduleDAO) New() (output moduleDAO) {
	output.FileName = "ModuleDAO.go"
	output.TableName = "module"
	return
}

func (input moduleDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "m.id"},
		Deleted:   FieldStatus{FieldName: "m.deleted"},
		CreatedBy: FieldStatus{FieldName: "m.created_by", Value: createdBy},
	}
}

func (input moduleDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "m." + (*searchByParam)[i].SearchKey
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
		userParam.OrderBy = "m." + userParam.OrderBy
		break
	}
}

func (input moduleDAO) GetModuleForUpdate(db *sql.Tx, userParam repository.ModuleModel) (result repository.ModuleModel, err errorModel.ErrorModel) {
	funcName := "GetModuleForUpdate"
	query := fmt.Sprintf(
		`SELECT
			m.id, m.updated_at, m.created_by,
			CASE WHEN
			(SELECT COUNT(p.id) FROM %s p WHERE
				(
					p.module_id_1 = m.id OR p.module_id_2 = m.id OR p.module_id_3 = m.id OR
					p.module_id_4 = m.id OR p.module_id_5 = m.id OR p.module_id_6 = m.id OR 
					p.module_id_7 = m.id OR p.module_id_8 = m.id OR p.module_id_9 = m.id OR
					p.module_id_10 = m.id 
				) AND p.deleted = false
			) > 0
			OR
			(SELECT COUNT(lc.id) FROM %s lc WHERE
				(
					lc.module_id_1 = m.id OR lc.module_id_2 = m.id OR lc.module_id_3 = m.id OR
					lc.module_id_4 = m.id OR lc.module_id_5 = m.id OR lc.module_id_6 = m.id OR
					lc.module_id_7 = m.id OR lc.module_id_8 = m.id OR lc.module_id_9 = m.id OR
					lc.module_id_10 = m.id
				) AND lc.deleted = false
			) > 0
			THEN TRUE ELSE FALSE END is_used, m.module_name
		FROM %s m
		WHERE
			m.id = $1 AND m.deleted = FALSE `,
		ProductDAO.TableName, LicenseConfigDAO.TableName, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	query += " FOR UPDATE "

	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.UpdatedAt,
		&result.CreatedBy, &result.IsUsed,
		&result.ModuleName,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input moduleDAO) DeleteModule(db *sql.Tx, userParam repository.ModuleModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteModule"
		query    string
	)

	query = fmt.Sprintf(
		`UPDATE %s SET
		deleted = $1, module_name = $2, updated_by = $3, 
		updated_at = $4, updated_client = $5
		WHERE
		id = $6 `,
		input.TableName)

	param := []interface{}{
		true, userParam.ModuleName.String, userParam.UpdatedBy.Int64,
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

func (input moduleDAO) UpdateModule(db *sql.Tx, userParam repository.ModuleModel) (err errorModel.ErrorModel) {
	funcName := "UpdateModule"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			module_name = $1, updated_by = $2,
			updated_at = $3, updated_client = $4
		WHERE
			id = $5 `, input.TableName)

	params := []interface{}{
		userParam.ModuleName.String,
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

func (input moduleDAO) GetCountModule(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	for i, item := range searchByParam {
		searchByParam[i].SearchKey = "module." + item.SearchKey
	}
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input moduleDAO) ViewModule(db *sql.DB, userParam repository.ModuleModel) (result repository.ModuleModel, err errorModel.ErrorModel) {
	funcName := "ViewModule"

	query := fmt.Sprintf(
		`SELECT
			m.id, m.module_name, m.created_at, m.updated_at,
			m.updated_by, u.nt_username, m.created_by
		FROM %s m
		LEFT JOIN "%s" AS u ON m.updated_by = u.id
		WHERE
			m.id = $1 AND m.deleted = FALSE `, input.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(
		&result.ID, &result.ModuleName,
		&result.CreatedAt, &result.UpdatedAt,
		&result.UpdatedBy, &result.UpdatedName,
		&result.CreatedBy,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input moduleDAO) GetListModule(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(
		`SELECT
			m.id, m.module_name,
			m.created_at, m.updated_at,
			m.updated_by, u.nt_username
		FROM %s m
		LEFT JOIN "%s" AS u
			ON m.updated_by = u.id `, input.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ModuleModel
			dbError := rows.Scan(
				&temp.ID, &temp.ModuleName,
				&temp.CreatedAt, &temp.UpdatedAt,
				&temp.UpdatedBy, &temp.UpdatedName,
			)
			return temp, dbError
		}, " ", input.getDefaultMustCheck(createdBy))
}

func (input moduleDAO) InsertModule(db *sql.Tx, userParam repository.ModuleModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertModule"

	query := fmt.Sprintf(
		`INSERT INTO %s
			(
				module_name, created_by, created_at, created_client,
				updated_by, updated_at, updated_client
			)
		VALUES
			(
				$1, $2, $3, $4, $5, $6, $7
			)
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.ModuleName.String,
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

func (input moduleDAO) CheckModuleIsExist(db *sql.DB, userParam repository.ModuleModel) (result repository.ModuleModel, err errorModel.ErrorModel) {
	funcName := "CheckModuleIsExist"

	query := fmt.Sprintf(
		`SELECT 
			id
		FROM %s
		WHERE
			id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	dbError := db.QueryRow(query, param...).Scan(&result.ID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
