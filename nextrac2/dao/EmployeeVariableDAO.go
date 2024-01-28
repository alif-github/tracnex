package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type employeeVariableDAO struct {
	AbstractDAO
}

var EmployeeVariableDAO = employeeVariableDAO{}.New()

func (input employeeVariableDAO) New() (output employeeVariableDAO) {
	output.FileName = "EmployeeVariableDAO.go"
	output.TableName = "employee_variable"
	return
}

func (input employeeVariableDAO) DeleteEmployeeVariableByEmployeeID(db *sql.Tx, userParam repository.EmployeeVariableModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteEmployeeVariableByEmployeeID"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		deleted = TRUE, updated_by = $1, updated_at = $2, 
		updated_client = $3 
		WHERE employee_id = $4 `,
		input.TableName)

	param := []interface{}{
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.EmployeeID.Int64,
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

func (input employeeVariableDAO) InsertEmployeeVariable(db *sql.Tx, userParam repository.EmployeeVariableModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertEmployeeVariable"
		query    string
	)

	query = fmt.Sprintf(
		`INSERT INTO %s
		(
			redmine_id, mandays_rate, created_by, 
			created_at, created_client, updated_by, 
			updated_at, updated_client, employee_id, 
			lead_mandays
		)
		VALUES
		(
			$1, $2, $3, 
			$4, $5, $6, 
			$7, $8, $9,
			$10
		)
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.RedmineID.Int64, userParam.MandaysRate.String, userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time, userParam.CreatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.EmployeeID.Int64,
	}

	if userParam.LeadMandays.Float64 > 0 {
		params = append(params, userParam.LeadMandays.Float64)
	} else {
		params = append(params, nil)
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeVariableDAO) UpdateEmployeeVariableByEmployeeID(db *sql.Tx, userParam repository.EmployeeVariableModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateEmployeeVariableByEmployeeID"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		redmine_id = $1, mandays_rate = $2, updated_by = $3, 
		updated_at = $4, updated_client = $5, lead_mandays = $6 
		WHERE id = $7 `,
		input.TableName)

	param := []interface{}{
		userParam.RedmineID.Int64, userParam.MandaysRate.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
	}

	if userParam.LeadMandays.Float64 > 0 {
		param = append(param, userParam.LeadMandays.Float64)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.ID.Int64)

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

func (input employeeVariableDAO) GetEmployeeVariableByEmployeeID(db *sql.DB, userParam repository.EmployeeVariableModel) (result repository.EmployeeVariableModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeeVariableByEmployeeID"
		query    string
	)

	query = fmt.Sprintf(`SELECT id FROM %s WHERE employee_id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.EmployeeID.Int64}
	dbRow := db.QueryRow(query, param...)
	dbErrs := dbRow.Scan(&result.ID)
	if dbErrs != nil && dbErrs != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErrs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeVariableDAO) GetEmployeeVariableForUpdate(db *sql.DB, userParam repository.EmployeeModel, isForUpdate bool) (result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeeVariableForUpdate"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
			e.id, e.id_card, 
			CASE WHEN ev.id > 0 OR ev.id IS NOT NULL 
				THEN true 
				ELSE false 
			END is_variable_exist, ev.updated_at, 
			ev.created_by
		FROM %s e `,
		EmployeeDAO.TableName)

	if isForUpdate {
		query += fmt.Sprintf(` 
			INNER JOIN %s ev ON ev.employee_id = e.id AND ev.deleted = FALSE 
			WHERE e.id = $1 AND e.deleted = FALSE FOR UPDATE `,
			input.TableName)
	} else {
		query += fmt.Sprintf(` 
			LEFT JOIN %s ev ON ev.employee_id = e.id AND ev.deleted = FALSE 
			WHERE e.id = $1 AND e.deleted = FALSE `,
			input.TableName)
	}

	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.IDCard, &result.IsHaveVariable,
		&result.UpdatedAt, &result.CreatedBy)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
