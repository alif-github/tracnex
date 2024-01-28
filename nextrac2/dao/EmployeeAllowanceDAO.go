package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type employeeAllowanceDAO struct {
	AbstractDAO
}

var EmployeeAllowanceDAO = employeeAllowanceDAO{}.New()

func (input employeeAllowanceDAO) New() (output employeeAllowanceDAO) {
	output.TableName = "allowances"
	output.FileName = "EmployeeAllowanceDAO.go"
	return
}

func (input employeeAllowanceDAO) InsertEmployeeAllowance(db *sql.Tx, inputStruct repository.EmpAllowanceModel) (lastInsertedId int64, err errorModel.ErrorModel) {
	funcName := "InsertEmployeeAllowance"
	query := "INSERT INTO " + input.TableName + " (" +
		"	allowance_name, allowance_type, updated_client, created_client, " +
		"	created_at, created_by, updated_at, updated_by) " +
		"VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 ) " +
		" RETURNING id"
	params := []interface{}{
		inputStruct.AllowanceName.String, inputStruct.Type.String, inputStruct.UpdatedClient.String, inputStruct.CreatedClient,
		inputStruct.CreatedAt.Time, inputStruct.CreatedBy.Int64, inputStruct.UpdatedAt.Time,
		inputStruct.UpdatedBy.Int64,
	}
	result := db.QueryRow(query, params...)
	dbError := result.Scan(&lastInsertedId)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeAllowanceDAO) GetEmployeeAllowance(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) ([]interface{}, errorModel.ErrorModel) {
	query := `SELECT 
			id, allowance_name, allowance_type,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM ` + input.TableName + " "

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var inputStruct repository.EmpAllowanceModel

		dbError := rows.Scan(
			&inputStruct.ID,
			&inputStruct.AllowanceName,
			&inputStruct.Type,
			&inputStruct.CreatedAt,
			&inputStruct.CreatedBy,
			&inputStruct.UpdatedAt,
			&inputStruct.UpdatedBy,)

		return inputStruct, dbError
	}

	return GetListDataDAO.GetListData(db, query, userParam, searchBy, createdBy, mappingFunc, "")
}

func (input employeeAllowanceDAO) UpdateEmployeeAllowance(db *sql.Tx, inputStruct repository.EmpAllowanceModel) errorModel.ErrorModel {
	funcName := "UpdateEmployeeAllowance"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	allowance_name = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4, allowance_type = $5 " +
		"WHERE " +
		"	id = $6 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		inputStruct.AllowanceName.String,
		inputStruct.UpdatedClient.String,
		inputStruct.UpdatedAt.Time,
		inputStruct.UpdatedBy.Int64,
		inputStruct.Type.String,
		inputStruct.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeAllowanceDAO) DeleteEmployeeAllowance(db *sql.Tx, inputStruct repository.EmpAllowanceModel) errorModel.ErrorModel {
	funcName := "DeleteEmployeeAllowance"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	deleted = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4 " +
		"WHERE " +
		"	id = $5 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		inputStruct.Deleted.Bool,
		inputStruct.UpdatedClient.String,
		inputStruct.UpdatedAt.Time,
		inputStruct.UpdatedBy.Int64,
		inputStruct.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeAllowanceDAO) GetDetailEmployeeAllowance(db *sql.Tx, id int64) (inputStruct repository.EmpAllowanceModel, err errorModel.ErrorModel) {
	funcName := "GetDetailEmployeeAllowance"
	query := "SELECT id, updated_at " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND id = $1 "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&inputStruct.ID, &inputStruct.UpdatedAt)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeAllowanceDAO) GetCountEmployeeAllowance(db *sql.DB) (count int64, err errorModel.ErrorModel) {
	funcName := "GetCountEmployeeAllowance"
	query := "SELECT COUNT(*) " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE "

	param := []interface{}{}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&count)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeAllowanceDAO) GetAllowanceForDetail(db *sql.DB) (results []repository.EmpAllowanceModel, err errorModel.ErrorModel) {
	funcName := "GetAllowanceForDetail"

	query := "SELECT id, allowance_name, allowance_type, " +
		" created_at, created_by, updated_at, updated_by " +
		" FROM " + input.TableName +
		" WHERE deleted = false "

	rows, errorRows := db.Query(query)

	if errorRows != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
		return
	}

	if rows != nil {
		defer func() {
			errorRows = rows.Close()
			if errorRows != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
				return
			}
		}()

		for rows.Next() {
			var inputStruct repository.EmpAllowanceModel
			errorRows = rows.Scan(&inputStruct.ID,
				&inputStruct.AllowanceName,
				&inputStruct.Type,
				&inputStruct.CreatedAt,
				&inputStruct.CreatedBy,
				&inputStruct.UpdatedAt,
				&inputStruct.UpdatedBy)

			if errorRows != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
				return
			}
			results = append(results, inputStruct)
		}

	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
		return
	}
	return
}

func (input employeeAllowanceDAO) CheckAllowance(db *sql.Tx, key string, fieldName string) (idCheck int64, err errorModel.ErrorModel) {
	funcName := "CheckAllowance"
	query := "SELECT " +
		"	id FROM " + input.TableName + " " +
		" WHERE LOWER("+ fieldName + ") = LOWER($1) AND deleted = FALSE LIMIT 1 "

	param := []interface{}{key}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&idCheck)
	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}