package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

type employeeFacilitiesActiveDAO struct {
	AbstractDAO
}

var EmployeeFacilitiesActiveDAO = employeeFacilitiesActiveDAO{}.New()

func (input employeeFacilitiesActiveDAO) New() (output employeeFacilitiesActiveDAO) {
	output.TableName = "employee_facilities_active"
	output.FileName = "EmployeeFacilitiesActiveDAO.go"
	return
}

func (input employeeFacilitiesActiveDAO) InsertEmployeeFacilitiesActive(db *sql.Tx, inputStuct repository.EmployeeFacilitiesActiveModel) (lastInsertedId int64, err errorModel.ErrorModel) {
	funcName := "InsertEmployeeFacilitiesActive"
	query := "INSERT INTO " + input.TableName + " (" +
		"	value, active, employee_level_id, employee_grade_id," +
		"   updated_client, created_client, " +
		"	created_at, created_by, updated_at, updated_by, allowance_id, benefit_id) " +
		"VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12 ) " +
		" RETURNING id"
	params := []interface{}{
		inputStuct.Value.String, inputStuct.Active.Bool, inputStuct.LevelID.Int64,
		inputStuct.GradeID.Int64,
		inputStuct.UpdatedClient.String, inputStuct.CreatedClient, inputStuct.CreatedAt.Time,
		inputStuct.CreatedBy.Int64, inputStuct.UpdatedAt.Time, inputStuct.UpdatedBy.Int64,
	}

	if inputStuct.AllowanceID.Int64 > 0 {
		params = append(params, inputStuct.AllowanceID.Int64)
	} else {
		params = append(params, nil)
	}

	if inputStuct.BenefitID.Int64 > 0 {
		params = append(params, inputStuct.BenefitID.Int64)
	} else {
		params = append(params, nil)
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

func (input employeeFacilitiesActiveDAO) GetDetailEmployeeMatrixForUpdate(db *sql.DB, idLevel int64, idGrade int64) (count int64, err errorModel.ErrorModel) {
	funcName := "GetDetailEmployeeMatrixForUpdate"
	query := "SELECT COUNT(*) " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND employee_level_id = $1 AND employee_grade_id = $2 "

	param := []interface{}{idLevel, idGrade}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&count)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeFacilitiesActiveDAO) DeleteEmployeeMatrix(db *sql.DB, inputStruct repository.EmployeeFacilitiesActiveModel) errorModel.ErrorModel {
	funcName := "DeleteEmployeeMatrix"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	deleted = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4 " +
		"WHERE " +
		"	employee_level_id = $5 AND employee_grade_id = $6 AND " +
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
		inputStruct.LevelID.Int64,
		inputStruct.GradeID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeFacilitiesActiveDAO) UpdateEmployeeMatrix(db *sql.Tx, inputStruct repository.EmployeeFacilitiesActiveModel) errorModel.ErrorModel {
	funcName := "UpdateEmployeeMatrix"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	value = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4, active = $5 " +
		"WHERE " +
		"	id = $6 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		inputStruct.Value.String,
		inputStruct.UpdatedClient.String,
		inputStruct.UpdatedAt.Time,
		inputStruct.UpdatedBy.Int64,
		inputStruct.Active.Bool,
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

func (input employeeFacilitiesActiveDAO) CekEmployeeMatrixForUpdate(db *sql.Tx, params repository.EmployeeFacilitiesActiveModel, fieldName string, idM int64) (id int64, err errorModel.ErrorModel) {
	funcName := "CekEmployeeMatrixForUpdate"
	query := "SELECT id " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND employee_level_id = $1 AND employee_grade_id = $2 " +
		" AND " + fieldName + " = $3 LIMIT 1"

	param := []interface{}{params.LevelID.Int64, params.GradeID.Int64, idM}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeFacilitiesActiveDAO) CheckEmployeeMatrixIsExist(db *sql.DB, params repository.EmployeeFacilitiesActiveModel) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName = "CheckEmployeeMatrixIsExist"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
			CASE WHEN efa.id > 0 
			THEN TRUE 
			ELSE FALSE 
			END is_exist 
		FROM %s efa 
			INNER JOIN %s el ON efa.employee_level_id = el.id 
			INNER JOIN %s eg ON efa.employee_grade_id = eg.id
		WHERE
		efa.active IS TRUE AND 
		el.deleted = FALSE AND 
		eg.deleted = FALSE AND 
		el.id = $1 AND 
		eg.id = $2
		GROUP BY efa.id, el.id, el."level", eg.id, eg.grade 
		LIMIT 1 `,
		input.TableName, EmployeeLevelDAO.TableName, EmployeeGradeDAO.TableName)

	param := []interface{}{params.LevelID.Int64, params.GradeID.Int64}
	results := db.QueryRow(query, param...)
	dbError := results.Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeFacilitiesActiveDAO) GetDetailEmployeeMatrix(db *sql.DB, idLevel int64, idGrade int64) (matrix repository.EmployeeFacilitiesActiveModel, err errorModel.ErrorModel) {
	funcName := "GetDetailEmployeeMatrix"
	query := "SELECT employee_facilities_active.employee_level_id, " +
		" employee_facilities_active.employee_grade_id, " +
		" employee_grade.grade, employee_level.level " +
		" FROM " + input.TableName +
		" LEFT JOIN employee_grade ON employee_grade.id = employee_facilities_active.employee_grade_id " +
		" LEFT JOIN employee_level ON employee_level.id = employee_facilities_active.employee_level_id " +
		" WHERE employee_facilities_active.deleted = FALSE AND " +
		" employee_facilities_active.employee_level_id = $1 " +
		" AND employee_facilities_active.employee_grade_id = $2 "

	param := []interface{}{idLevel, idGrade}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&matrix.LevelID, &matrix.GradeID, &matrix.Grade, &matrix.Level)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeFacilitiesActiveDAO) GetMatrixForDetail(db *sql.DB, idGrade int64, idLevel int64) (results []repository.EmployeeFacilitiesActiveModel, err errorModel.ErrorModel) {
	funcName := "GetMatrixForDetail"

	query := "SELECT id, value, active, " +
		" allowance_id, benefit_id " +
		" FROM " + input.TableName +
		" WHERE deleted = false " +
		" AND employee_level_id = " + strconv.Itoa(int(idLevel)) +
		" AND employee_grade_id = " + strconv.Itoa(int(idGrade))

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
			var inputStruct repository.EmployeeFacilitiesActiveModel
			errorRows = rows.Scan(&inputStruct.ID,
				&inputStruct.Value,
				&inputStruct.Active,
				&inputStruct.AllowanceID,
				&inputStruct.BenefitID)

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

func (input employeeFacilitiesActiveDAO) GetListMatrixs(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, search_by string, keyword string) (kms []interface{}, err errorModel.ErrorModel) {
	query := `SELECT DISTINCT employee_facilities_active.employee_level_id, 
              employee_facilities_active.employee_grade_id, lvl.level, grade.grade
              FROM employee_facilities_active
              LEFT JOIN employee_level lvl ON lvl.id = employee_facilities_active.employee_level_id
              LEFT JOIN employee_grade grade ON grade.id = employee_facilities_active.employee_grade_id `

	getListJoin := getListJoinDataDAO{Table: input.TableName, Query: query}
	getListJoin.SetWhere("employee_facilities_active.deleted", "FALSE")
	if search_by != "" && keyword != "" {
		if search_by == "lvl.level" || search_by == "grade.grade" {
			getListJoin.SetWhereAdditional(search_by + " ILIKE '%" + keyword + "%'")
		}
	}

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var inputStruct repository.EmployeeFacilitiesActiveModel
		dbError := rows.Scan(
			&inputStruct.LevelID, &inputStruct.GradeID, &inputStruct.Level, &inputStruct.Grade)

		return inputStruct, dbError
	}

	return getListJoin.GetListJoinData(db, userParam, searchBy, createdBy, mappingFunc)
}

func (input employeeFacilitiesActiveDAO) GetCountEmployeeMatrix(db *sql.DB, search_by string, keyword string) (count int64, err errorModel.ErrorModel) {
	funcName := "GetCountEmployeeMatrix"
	query := `SELECT COUNT(DISTINCT (employee_facilities_active.employee_level_id, 
              employee_facilities_active.employee_grade_id, lvl.level, grade.grade))
              FROM employee_facilities_active
              LEFT JOIN employee_level lvl ON lvl.id = employee_facilities_active.employee_level_id
              LEFT JOIN employee_grade grade ON grade.id = employee_facilities_active.employee_grade_id 
              WHERE employee_facilities_active.deleted=FALSE `

	if search_by != "" && keyword != "" {
		if search_by == "lvl.level" || search_by == "grade.grade" {
			query += "AND " + search_by + " ILIKE '%" + keyword + "%'"
		}
	}

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

func (input employeeFacilitiesActiveDAO) GetCountEmployeeMatrixForMaster(db *sql.DB, field string, id int64) (count int64, err errorModel.ErrorModel) {
	funcName := "GetCountEmployeeMatrixForMaster"
	query := "SELECT COUNT(*) " +
		" FROM employee_facilities_active " +
		" WHERE " + field + " = $1 " +
		"  AND employee_facilities_active.deleted=FALSE "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&count)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeFacilitiesActiveDAO) GetByAllowanceIdAndEmployeeLevelIdAndEmployeeGradeIdTx(tx *sql.Tx, model repository.EmployeeFacilitiesActiveModel) (result repository.EmployeeFacilitiesActiveModel, errModel errorModel.ErrorModel) {
	funcName := "GetByAllowanceIdAndEmployeeLevelIdAndEmployeeGradeIdTx"

	query := `SELECT 
				efa.id, efa.value, al.allowance_type
			FROM employee_facilities_active AS efa 
			INNER JOIN allowances AS al
				ON efa.allowance_id = al.id
			WHERE 
				efa.active = TRUE AND 
				al.active = TRUE AND 
				al.id = $1 AND
				efa.employee_level_id = $2 AND 
				efa.employee_grade_id = $3`

	params := []interface{}{
		model.AllowanceID.Int64,
		model.LevelID.Int64,
		model.GradeID.Int64,
	}

	row := tx.QueryRow(query, params...)
	err := row.Scan(&result.ID, &result.Value, &result.AllowanceType)
	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}
