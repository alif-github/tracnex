package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type employeePositionDAO struct {
	AbstractDAO
}

var EmployeePositionDAO = employeePositionDAO{}.New()

func (input employeePositionDAO) New() (output employeePositionDAO) {
	output.FileName = "EmployeePositionDAO.go"
	output.TableName = "employee_position"
	return
}

func (input employeePositionDAO) CheckEmployeePosition(db *sql.DB, idEmployeePosition int64) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName = "CheckEmployeePosition"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT CASE WHEN id > 0 THEN true ELSE false END is_exist 
		FROM %s WHERE id = $1 AND deleted = FALSE `,
		input.TableName)

	param := []interface{}{idEmployeePosition}
	dbError := db.QueryRow(query, param...).Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionDAO) GetListPosition(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query    string
		addWhere string
	)

	query = fmt.Sprintf(`
		SELECT 
		    p.id, p.name, p.description, 
		    p.updated_at 
		FROM %s p 
		INNER JOIN %s c ON p.internal_company_id = c.id `,
		input.TableName, CompanyDAO.TableName)

	addWhere += fmt.Sprintf(` AND c.deleted = FALSE `)
	input.convertUserParamAndSearchBy(&userParam, &searchBy)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.EmployeePositionModel
		dbError := rows.Scan(
			&temp.ID, &temp.Name, &temp.Description,
			&temp.UpdatedAt,
		)
		return temp, dbError
	}, addWhere, DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "p.id"},
		Deleted:   FieldStatus{FieldName: "p.deleted"},
		CreatedBy: FieldStatus{FieldName: "p.created_by", Value: int64(0)},
	})
}

func (input employeePositionDAO) GetCountPosition(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	var (
		addWhere  string
		tableName string
	)

	for i, item := range searchByParam {
		searchByParam[i].SearchKey = "p." + item.SearchKey
	}

	tableName = fmt.Sprintf(`
		%s p 
		INNER JOIN %s c ON p.internal_company_id = c.id `,
		input.TableName, CompanyDAO.TableName)

	addWhere += fmt.Sprintf(` AND c.deleted = FALSE `)
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, addWhere, DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "p.id"},
		Deleted:   FieldStatus{FieldName: "p.deleted"},
		CreatedBy: FieldStatus{FieldName: "p.created_by", Value: createdBy},
	})
}

func (input employeePositionDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "p." + (*searchByParam)[i].SearchKey
	}

	userParam.OrderBy = "p." + userParam.OrderBy
}

func (input employeePositionDAO) InsertEmployeePosition(db *sql.Tx, userParam repository.EmployeePositionModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertEmployeePosition"
		query    string
	)

	query = fmt.Sprintf(`
		INSERT INTO %s 
		(
		 "name", description, internal_company_id,
		 created_by, created_at, created_client, 
		 updated_by, updated_at, updated_client
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id `, input.TableName)

	params := []interface{}{userParam.Name.String}

	if userParam.Description.String != "" {
		params = append(params, userParam.Description.String)
	} else {
		params = append(params, nil)
	}

	params = append(params, userParam.CompanyID.Int64,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String)

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionDAO) GetEmployeePositionForUpdate(db *sql.DB, userParam repository.EmployeePositionModel) (result repository.EmployeePositionModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeePositionForUpdate"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT id, updated_at, created_by 
		FROM %s 
		WHERE id = $1 AND deleted = FALSE `,
		input.TableName)

	param := []interface{}{userParam.ID.Int64}
	query += " FOR UPDATE "

	dbErr := db.QueryRow(query, param...).Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy)
	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErr)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionDAO) UpdateEmployeePosition(db *sql.Tx, userParam repository.EmployeePositionModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateEmployeePosition"
	)

	query := fmt.Sprintf(`
		UPDATE %s
		SET
			"name" = $1, description = $2, internal_company_id = $3, 
		    updated_by = $4, updated_at = $5, updated_client = $6
		WHERE
		id = $7 `, input.TableName)

	params := []interface{}{
		userParam.Name.String, userParam.Description.String, userParam.CompanyID.Int64,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
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

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionDAO) ViewEmployeePosition(db *sql.DB, userParam repository.EmployeePositionModel) (result repository.EmployeePositionModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewEmployeePosition"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
			p.id, p.name, p.description, 
			c.company_name, p.created_at, p.updated_at, 
			p.created_by, p.updated_by, uc.nt_username as created_name, 
			up.nt_username as updated_name 
		FROM %s p 
			INNER JOIN %s c ON p.internal_company_id = c.id 
			LEFT JOIN "%s" uc ON p.created_by = uc.id 
			LEFT JOIN "%s" up ON p.updated_by = up.id 
		WHERE 
		    p.id = $1 AND p.deleted = FALSE `,
		input.TableName, CompanyDAO.TableName, UserDAO.TableName,
		UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbErr := results.Scan(
		&result.ID, &result.Name, &result.Description,
		&result.CompanyName, &result.CreatedAt, &result.UpdatedAt,
		&result.CreatedBy, &result.UpdatedBy, &result.CreatedName,
		&result.UpdatedName)

	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErr)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionDAO) DeleteEmployeePosition(db *sql.Tx, userParam repository.EmployeePositionModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteEmployeePosition"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET
		deleted = TRUE, updated_by = $1, updated_at = $2, 
		updated_client = $3
		WHERE
		id = $4 `,
		input.TableName)

	param := []interface{}{
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.ID.Int64,
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

func (input employeePositionDAO) GetNamePositionByIDWithoutDeleted(db *sql.DB, userParam repository.EmployeePositionModel) (result repository.EmployeePositionModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetNamePositionByIDWithoutDeleted"
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
