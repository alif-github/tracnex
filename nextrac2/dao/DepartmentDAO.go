package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type DepartmentDAOInterface interface {
	CheckIDDepartment(*sql.DB, repository.DepartmentModel) (bool, errorModel.ErrorModel)
}

type departmentDAO struct {
	AbstractDAO
}

var DepartmentDAO = departmentDAO{}.New()

func (input departmentDAO) New() (output departmentDAO) {
	output.FileName = "DepartmentDAO.go"
	output.TableName = "department"
	return
}

func (input departmentDAO) GetListDepartment(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	var (
		dbParam []interface{}
	)

	query := fmt.Sprintf(`SELECT id, name, updated_at, description FROM %s `, input.TableName)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.DepartmentModel
			dbError := rows.Scan(
				&temp.ID, &temp.Name, &temp.UpdatedAt, &temp.Description,
			)
			return temp, dbError
		}, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input departmentDAO) GetCountDepartment(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input departmentDAO) CheckIDDepartment(db *sql.DB, userParam repository.DepartmentModel) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName = "CheckIDDepartment"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT CASE WHEN id > 0 
			THEN true 
			ELSE false 
		END is_exist 
		FROM %s 
		WHERE id = $1 AND deleted = FALSE `,
		input.TableName)

	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input departmentDAO) InsertDepartment(db *sql.Tx, userParam repository.DepartmentModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertDepartment"
	)

	query := fmt.Sprintf(
		`INSERT INTO %s
			(
				name,
				created_by, created_at, created_client,
				updated_by, updated_at, updated_client,
				description
			)
		VALUES
			(
				$1, $2, $3, $4, $5, $6, $7, $8
			)
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.Name.String,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
	}

	if util.IsStringEmpty(userParam.Description.String) {
		params = append(params, nil)
	} else {
		params = append(params, userParam.Description.String)
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input departmentDAO) ViewDepartment(db *sql.DB, userParam repository.DepartmentModel) (result repository.DepartmentModel, err errorModel.ErrorModel) {
	funcName := "ViewDepartment"

	query := fmt.Sprintf(
		`SELECT
			d.id, d.name, d.description, 
			d.created_at, d.updated_at,
			d.updated_by, u.nt_username, 
			d.created_by
		FROM %s d
		LEFT JOIN "%s" AS u ON d.updated_by = u.id
		WHERE
			d.id = $1 AND d.deleted = FALSE `, input.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(
		&result.ID, &result.Name, &result.Description,
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

func (input departmentDAO) GetDepartmentForUpdate(db *sql.Tx, userParam repository.DepartmentModel) (result repository.DepartmentModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetDepartmentForUpdate"
	)

	query := fmt.Sprintf(
		`SELECT
			d.id, d.updated_at, d.created_by, 
			d.name, 
			CASE WHEN
			( SELECT COUNT(e.id) FROM %s e WHERE
				(
					e.department_id = d.id
				) AND e.deleted = false
			) > 0
			THEN TRUE ELSE FALSE END is_used 
		FROM %s d
		WHERE
			d.id = $1 AND d.deleted = FALSE `,
		EmployeeDAO.TableName, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	query += " FOR UPDATE "

	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.UpdatedAt,
		&result.CreatedBy, &result.Name, &result.IsUsed,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input departmentDAO) DeleteDepartment(db *sql.Tx, userParam repository.DepartmentModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteDepartment"
		query    string
	)

	query = fmt.Sprintf(
		`UPDATE %s SET
		deleted = $1, name = $2, updated_by = $3, 
		updated_at = $4, updated_client = $5
		WHERE
		id = $6 `,
		input.TableName)

	param := []interface{}{
		true, userParam.Name.String, userParam.UpdatedBy.Int64,
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

func (input departmentDAO) UpdateDepartmentByID(db *sql.Tx, userParam repository.DepartmentModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateDepartmentByID"
	)

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			name = $1, description = $2, 
			updated_by = $3, updated_at = $4, updated_client = $5 
		WHERE
			id = $6 `, input.TableName)

	params := []interface{}{
		userParam.Name.String,
		userParam.Description.String,
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

func (input departmentDAO) GetNameDepartmentByIDWithoutDeleted(db *sql.DB, userParam repository.DepartmentModel) (result repository.DepartmentModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetNameDepartmentByIDWithoutDeleted"
		query    string
	)

	query = fmt.Sprintf(`SELECT name FROM %s WHERE id = $1 `, input.TableName)
	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result.Name)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
