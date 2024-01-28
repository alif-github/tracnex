package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type roleDAO struct {
	AbstractDAO
}

var RoleDAO = roleDAO{}.New()

func (input roleDAO) New() (output roleDAO) {
	output.FileName = "RoleDAO.go"
	output.TableName = "role"
	return
}

func (input roleDAO) GetRoleByName(db *sql.DB, userParam repository.RoleModel) (result repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "GetRoleByName"
	query := fmt.Sprintf(`SELECT id, role_id FROM %s WHERE role_id = $1 AND created_client != 'SYSTEM' AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.RoleID.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.RoleID)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) InsertRole(db *sql.Tx, userParam repository.RoleModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertRole"

	query :=
		"INSERT INTO role(role_id, description, permission, created_by, created_client, created_at, updated_by, updated_client, updated_at) VALUES " +
			"($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

	listInterface := []interface{}{userParam.RoleID.String, userParam.Description.String, userParam.Permission.String, userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time}
	results := db.QueryRow(query, listInterface...)

	errorS := results.Scan(&id)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input roleDAO) GetRoleForUpdate(db *sql.Tx, userParam repository.RoleModel) (output repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "GetRoleForUpdate"

	query := fmt.Sprintf(`SELECT id, role_id, permission, updated_at, created_by FROM %s WHERE id = $1 AND deleted = FALSE`, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += fmt.Sprintf(` FOR UPDATE`)

	results := db.QueryRow(query, param...)

	errs := results.Scan(&output.ID, &output.RoleID, &output.Permission, &output.UpdatedAt, &output.CreatedBy)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) UpdateRole(db *sql.Tx, userParam repository.RoleModel) (err errorModel.ErrorModel) {
	funcName := "UpdateRole"

	query :=
		"UPDATE role SET role_id = $1, description = $2, permission = $3, updated_by = $4, " +
			" updated_client = $5, updated_at = $6 " +
			"WHERE id = $7 "
	param := []interface{}{userParam.RoleID.String, userParam.Description.String, userParam.Permission.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time, userParam.ID.Int64}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	return
}

func (input roleDAO) ViewRole(db *sql.DB, userParam repository.RoleModel) (result repository.RoleModel, errors errorModel.ErrorModel) {
	funcName := "ViewRole"

	query := fmt.Sprintf(`SELECT id, role_id, description, permission, created_by, created_at, updated_by, updated_at 
									FROM %s WHERE id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 != 0 {
		query += " AND created_by = $2"
		param = append(param, userParam.CreatedBy.Int64)
	}

	results := db.QueryRow(query, param...)

	err := results.Scan(&result.ID, &result.RoleID, &result.Description,
		&result.Permission, &result.CreatedBy, &result.CreatedAt,
		&result.UpdatedBy, &result.UpdatedAt)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errors = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) GetListRole(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(
		`SELECT 
			r.id, r.role_id, r.description, 
			r.created_by, r.updated_at, u.nt_username as created_name,
			r.created_at
		FROM %s r 
		LEFT JOIN "%s" u ON r.created_by = u.id `,
		input.TableName, UserDAO.TableName)

	additionalWhere := fmt.Sprintf(` AND r.role_id NOT IN ( 'user_nd6',  'user_nexmile') `)
	input.convertUserParamAndSearchBy(&userParam, searchBy)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.RoleModel
			errors := rows.Scan(&temp.ID, &temp.RoleID, &temp.Description,
				&temp.CreatedBy, &temp.UpdatedAt, &temp.CreatedName, &temp.CreatedAt)
			return temp, errors
		}, additionalWhere, input.getDefaultMustCheck(createdBy))
}

func (input roleDAO) GetCountRole(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	additionalWhere := " AND role_id NOT IN ( 'user_nd6',  'user_nexmile') "
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, additionalWhere, DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input roleDAO) DeleteRole(db *sql.Tx, userParam repository.RoleModel) (err errorModel.ErrorModel) {
	funcName := "DeleteRole"
	query := "UPDATE role set " +
		" deleted = TRUE, role_id = $1, updated_by = $2, " +
		" updated_client = $3, updated_at = $4 " +
		" WHERE id = $5 "
	param := []interface{}{userParam.RoleID.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time, userParam.ID.Int64}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) GetRoleForDelete(db *sql.Tx, userParam repository.RoleModel) (output repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "GetDataGroupForDelete"

	query := fmt.Sprintf(`SELECT role.id, role.updated_at, 
									(SELECT CASE WHEN count(id) > 0 THEN TRUE ELSE FALSE END 
										FROM %s 
										WHERE role_id = role.id) isUsed, role.created_by, role.role_id 
								FROM %s role 
								WHERE role.id = $1 AND role.deleted = FALSE `, ClientRoleScopeDAO.TableName, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND role.created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += fmt.Sprintf(` FOR UPDATE`)

	results := db.QueryRow(query, param...)

	errs := results.Scan(&output.ID, &output.UpdatedAt, &output.IsUsed, &output.CreatedBy, &output.RoleID)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) CheckIsRoleExist(db *sql.Tx, userParam repository.RoleModel) (result repository.RoleModel, err errorModel.ErrorModel) {
	funcName := "CheckIsRoleExist"

	query := fmt.Sprintf(`SELECT id, permission FROM %s WHERE id = $1 AND deleted = FALSE`, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.Permission)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam []in.SearchByParam) {
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

func (input roleDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "r.id"},
		Deleted:   FieldStatus{FieldName: "r.deleted"},
		CreatedBy: FieldStatus{FieldName: "r.created_by", Value: createdBy},
	}
}

func (input roleDAO) GetRolePermission(db *sql.DB, role string) (id int64, isExist bool, err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "GetRolePermission"
		errs     error
	)

	query := fmt.Sprintf(`
		SELECT id, CASE WHEN id > 0 THEN true ELSE false END is_exist 
		FROM %s WHERE role_id = $1 AND deleted = false `,
		input.TableName)

	dbRow := db.QueryRow(query, role)
	errs = dbRow.Scan(&id, &isExist)
	if errs != nil && errs != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) UpdateRolePermission(db *sql.DB, id int64, permission string) (err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "UpdateRolePermission"
	)

	query := fmt.Sprintf(`
		UPDATE %s SET permission = $1 WHERE id = $2`,
		input.TableName)

	param := []interface{}{permission, id}
	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleDAO) GetByIdTx(tx *sql.Tx, id int64) (result repository.RoleModel, errModel errorModel.ErrorModel) {
	funcName := "GetByIdTx"

	query := `SELECT 
				id
			FROM ` + input.TableName + ` 
			WHERE 
				id = $1 AND 
				deleted = FALSE`

	row := tx.QueryRow(query, id)
	err := row.Scan(&result.ID)
	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}