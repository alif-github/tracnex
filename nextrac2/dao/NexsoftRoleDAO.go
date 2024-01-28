package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type nexsoftRoleDAO struct {
	AbstractDAO
}

var NexsoftRoleDAO = nexsoftRoleDAO{}.New()

func (input nexsoftRoleDAO) New() (output nexsoftRoleDAO) {
	output.FileName = "NexsoftRoleDAO.go"
	output.TableName = "nexsoft_role"
	return
}

func (input nexsoftRoleDAO) GetNexsoftRoleByName(db *sql.DB, userParam repository.RoleModel) (result repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "GetRoleByName"
	query :=
		" SELECT " +
			"	id, role_id " +
			" FROM " + input.TableName + " " +
			"WHERE " +
			"	role_id = $1 AND created_client != 'SYSTEM' AND deleted = FALSE "

	param := []interface{}{userParam.RoleID.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.RoleID)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexsoftRoleDAO) InsertNesoftRole(db *sql.Tx, userParam repository.RoleModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertNesoftRole"
	query :=
		"INSERT INTO " + input.TableName + " " +
			" (role_id, " +
			" description, " +
			" permission, " +
			" created_by, " +
			" created_client, " +
			" created_at, " +
			" updated_by, " +
			" updated_client, " +
			" updated_at) " +
			" VALUES " +
			"($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

	params := []interface{}{
		userParam.RoleID.String,
		userParam.Description.String,
		userParam.Permission.String,
		userParam.CreatedBy.Int64,
		userParam.CreatedClient.String,
		userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time}
	results := db.QueryRow(query, params...)

	errorS := results.Scan(&id)

	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input nexsoftRoleDAO) GetListNexsoftRole(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(`
		SELECT 
			r.id, r.role_id, r.description, 
			r.created_by, r.created_at, r.updated_at, 
			u.nt_username as created_name
		FROM %s r 
		LEFT JOIN "%s" u ON r.created_by = u.id `,
		input.TableName, UserDAO.TableName)

	additionalWhere := fmt.Sprintf(` AND r.created_client != 'SYSTEM' `)
	input.convertUserParamAndSearchBy(&userParam, searchBy)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.RoleModel
			errors := rows.Scan(
				&temp.ID, &temp.RoleID, &temp.Description,
				&temp.CreatedBy, &temp.CreatedAt, &temp.UpdatedAt,
				&temp.CreatedName)
			return temp, errors
		}, additionalWhere, DefaultFieldMustCheck{
			Deleted: FieldStatus{FieldName: "r.deleted"},
			CreatedBy: FieldStatus{
				FieldName: "r.created_by",
				Value:     createdBy,
			},
		})
}

func (input nexsoftRoleDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam []in.SearchByParam) {
	for i := 0; i < len(searchByParam); i++ {
		searchByParam[i].SearchKey = "r." + searchByParam[i].SearchKey
	}

	switch userParam.OrderBy {
	case "created_name", "created_name ASC", "created_name DESC":
		//--- Continue
	default:
		userParam.OrderBy = "r." + userParam.OrderBy
		break
	}
}

func (input nexsoftRoleDAO) GetCountNexsoftRole(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	additionalWhere := " AND created_client != 'SYSTEM' "
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, additionalWhere, DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input nexsoftRoleDAO) GetNexsoftRoleForUpdate(db *sql.Tx, userParam repository.RoleModel) (output repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "GetNexsoftRoleForUpdate"

	query :=
		"SELECT " +
			"	id, role_id, permission, updated_at, created_by " +
			"FROM " + input.TableName + " " +
			"WHERE " +
			"	id = $1 AND created_client != 'SYSTEM' AND deleted = FALSE"
	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, param...)

	errs := results.Scan(&output.ID, &output.RoleID, &output.Permission, &output.UpdatedAt, &output.CreatedBy)
	if errs != nil && errs.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexsoftRoleDAO) UpdateNexoftRole(db *sql.Tx, userParam repository.RoleModel) (err errorModel.ErrorModel) {
	funcName := "UpdateNexoftRole"
	query := "UPDATE " + input.TableName + " " +
		"SET " +
		"	role_id = $1, " +
		"	description = $2, " +
		"	permission = $3, " +
		"	updated_by = $4, " +
		"	updated_client = $5, " +
		"	updated_at = $6 " +
		" WHERE id = $7 "

	params := []interface{}{
		userParam.RoleID.String,
		userParam.Description.String,
		userParam.Permission.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
		userParam.ID.Int64,
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	return
}

func (input nexsoftRoleDAO) GetDetailNexsoftRole(db *sql.DB, userParam repository.RoleModel) (result repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "GetDetailNexsoftRole"

	query := "SELECT " +
		"	id, " +
		"	role_id, " +
		"	description, " +
		"	permission, " +
		"	created_by, " +
		"	created_at, " +
		"	updated_by, " +
		"	updated_at " +
		" FROM " + input.TableName + " " +
		" WHERE " +
		"	id = $1 AND " +
		"	deleted = FALSE "

	params := []interface{}{
		userParam.ID.Int64,
	}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		params = append(params, userParam.CreatedBy.Int64)
	}

	results := db.QueryRow(query, params...)

	errorDB := results.Scan(
		&result.ID,
		&result.RoleID,
		&result.Description,
		&result.Permission,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.UpdatedBy,
		&result.UpdatedAt,
	)
	if errorDB != nil && errorDB.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexsoftRoleDAO) GetNexsoftRoleForDelete(db *sql.Tx, userParam repository.RoleModel) (output repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "GetNexsoftRoleForDelete"
	query :=
		"SELECT " +
			"	nexsoft_role.id, nexsoft_role.updated_at, " +
			"	(SELECT " +
			"		CASE WHEN count(id) > 0 THEN TRUE ELSE FALSE END " +
			"	FROM " +
			"		nexsoft_client_role_scope " +
			"	WHERE " +
			"		role_id = nexsoft_role.id) isUsed, nexsoft_role.created_by, nexsoft_role.role_id " +
			"FROM " + input.TableName + " nexsoft_role " +
			"WHERE " +
			"	nexsoft_role.id = $1 AND nexsoft_role.created_client != 'SYSTEM' AND nexsoft_role.deleted = FALSE "
	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND nexsoft_role.created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, param...)

	errs := results.Scan(&output.ID, &output.UpdatedAt, &output.IsUsed, &output.CreatedBy, &output.RoleID)
	if errs != nil && errs.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexsoftRoleDAO) DeleteNexsoftRole(db *sql.Tx, userParam repository.RoleModel) (err errorModel.ErrorModel) {
	funcName := "DeleteNexsoftRole"

	query := "UPDATE " + input.TableName + " " +
		" SET " +
		"	deleted = TRUE, " +
		"	role_id = $1, " +
		"	updated_by = $2, " +
		"	updated_client = $3, " +
		"	updated_at = $4 " +
		" WHERE " +
		"	id = $5 "

	params := []interface{}{
		userParam.RoleID.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
		userParam.ID.Int64}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
