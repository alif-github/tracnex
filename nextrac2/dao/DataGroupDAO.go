package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type dataGroupDAO struct {
	AbstractDAO
}

var DataGroupDAO = dataGroupDAO{}.New()

func (input dataGroupDAO) New() (output dataGroupDAO) {
	output.FileName = "DataGroupDAO.go"
	output.TableName = "data_group"
	return
}

func (input dataGroupDAO) GetRoleByName(db *sql.Tx, userParam repository.DataGroupModel) (result repository.DataGroupModel, err errorModel.ErrorModel) {
	funcName := "GetRoleByName"
	query :=
		" SELECT " +
			"	id, group_id " +
			" FROM " +
			input.TableName +
			" WHERE " +
			"	group_id = $1 AND deleted = FALSE "

	param := []interface{}{userParam.GroupID.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.GroupID)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupDAO) InsertDataGroup(db *sql.Tx, userParam repository.DataGroupModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertDataScope"
	query := "INSERT INTO data_group(group_id, description, scope, created_by, created_client, created_at, updated_by, updated_client, updated_at) VALUES " +
		"($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	param := []interface{}{userParam.GroupID.String, userParam.Description.String, userParam.Scope.String, userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time}

	results := db.QueryRow(query, param...)

	errs := results.Scan(&id)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupDAO) UpdateDataGroup(db *sql.Tx, userParam repository.DataGroupModel) (err errorModel.ErrorModel) {
	funcName := "UpdateDataGroup"
	query := "UPDATE data_group set " +
		" description = $1, scope = $2, updated_by = $3, updated_client = $4, " +
		" updated_at = $5 WHERE id = $6 "
	param := []interface{}{userParam.Description.String, userParam.Scope.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time, userParam.ID.Int64}

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

func (input dataGroupDAO) GetDataGroupForDelete(db *sql.Tx, userParam repository.DataGroupModel) (output repository.DataGroupModel, err errorModel.ErrorModel) {
	funcName := "GetDataGroupForDelete"
	query :=
		"SELECT " +
			"	dataGroup.id, dataGroup.updated_at, " +
			"	CASE WHEN " +
			"	(SELECT count(id) FROM client_role_scope WHERE group_id = dataGroup.id) > 0 OR " +
			"	(SELECT count(id) FROM nexsoft_client_role_scope WHERE group_id = dataGroup.id) > 0 " +
			"	THEN TRUE ELSE FALSE END is_used, dataGroup.group_id " +
			"FROM " +
			"	data_group dataGroup " +
			"WHERE " +
			"	dataGroup.id = $1 AND dataGroup.deleted = FALSE "
	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND dataGroup.created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, param...)

	errs := results.Scan(&output.ID, &output.UpdatedAt, &output.IsUsed, &output.GroupID)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupDAO) GetDataGroupForUpdate(db *sql.Tx, userParam repository.DataGroupModel) (output repository.DataGroupModel, err errorModel.ErrorModel) {
	funcName := "GetDataGroupForUpdate"
	query :=
		"SELECT " +
			"	id, group_id, scope, updated_at " +
			"FROM " +
			"	data_group " +
			"WHERE " +
			"	id = $1 AND deleted = FALSE"
	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND dataGroup.created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, param...)

	errs := results.Scan(&output.ID, &output.GroupID, &output.Scope, &output.UpdatedAt)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupDAO) DeleteDataGroup(db *sql.Tx, userParam repository.DataGroupModel) (err errorModel.ErrorModel) {
	funcName := "DeleteDataGroup"
	query :=
		"UPDATE data_group SET deleted = TRUE, group_id = $1, " +
			"	updated_at = $2, updated_client = $3, " +
			"	updated_by = $4 " +
			"WHERE " +
			" 	id = $5 AND deleted = FALSE "
	param := []interface{}{userParam.GroupID.String, userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.UpdatedBy.Int64, userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND dataGroup.created_by = $6 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupDAO) GetListDataGroup(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(
		`SELECT
			dg.id, dg.group_id, dg.created_at, dg.created_by, dg.updated_at,
			dg.updated_by, u.nt_username AS created_name, dg.description
		FROM %s dg 
		LEFT JOIN "%s" u ON dg.created_by = u.id `, input.TableName, UserDAO.TableName)

	fmt.Println("Query GetListDataGroup : ", query)
	input.convertUserParamAndSearchBy(&userParam, searchBy)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.DataGroupModel
			errors := rows.Scan(&temp.ID, &temp.GroupID, &temp.CreatedAt, &temp.CreatedBy,
				&temp.UpdatedAt, &temp.UpdatedBy, &temp.CreatedName, &temp.Description)
			return temp, errors
		}, "", input.getDefaultMustCheck(createdBy))
}

func (input dataGroupDAO) GetCountDataGroup(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input dataGroupDAO) ViewDataGroup(db *sql.DB, userParam repository.DataGroupModel) (result repository.DataGroupModel, errors errorModel.ErrorModel) {
	funcName := "ViewDataGroup"

	query :=
		"SELECT " +
			"	dg.id, dg.group_id, dg.description, dg.scope, dg.created_by, dg.updated_at, " +
			"	dg.updated_by, dg.created_at, u.nt_username AS updated_name " +
			"FROM " +
			"	data_group dg " +
			"LEFT JOIN \"" + UserDAO.TableName + "\" u " +
			"	ON dg.updated_by = u.id " +
			"WHERE " +
			"	dg.id = $1 AND dg.deleted = FALSE "
	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 != 0 {
		query += "created_by = $4"
		param = append(param, userParam.CreatedBy.Int64)
	}

	results := db.QueryRow(query, param...)

	err := results.Scan(&result.ID, &result.GroupID, &result.Description,
		&result.Scope, &result.CreatedBy, &result.UpdatedAt,
		&result.UpdatedBy, &result.CreatedAt, &result.UpdatedName)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errors = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "dg.id"},
		Deleted:   FieldStatus{FieldName: "dg.deleted"},
		CreatedBy: FieldStatus{FieldName: "dg.created_by", Value: createdBy},
	}
}

func (input dataGroupDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam []in.SearchByParam) {
	for i := 0; i < len(searchByParam); i++ {
		searchByParam[i].SearchKey = "dg." + searchByParam[i].SearchKey
	}

	switch userParam.OrderBy {
	case "created_name", "created_name ASC", "created_name DESC":
		userParam.OrderBy = "u.nt_username"
		break
	default:
		userParam.OrderBy = "dg." + userParam.OrderBy
		break
	}
}

func (input dataGroupDAO) GetByIdTx(tx *sql.Tx, id int64) (result repository.DataGroupModel, errModel errorModel.ErrorModel) {
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