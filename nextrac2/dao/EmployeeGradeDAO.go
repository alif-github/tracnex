package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

type employeeGradeDAO struct {
	AbstractDAO
}

var EmployeeGradeDAO = employeeGradeDAO{}.New()

func (input employeeGradeDAO) New() (output employeeGradeDAO) {
	output.TableName = "employee_grade"
	output.FileName = "employeeGradeDAO.go"
	return
}

func (input employeeGradeDAO) InsertEmployeeGrade(db *sql.Tx, grade repository.EmployeeGradeModel) (lastInsertedId int64, err errorModel.ErrorModel) {
	funcName := "InsertEmployeeGrade"
	query := "INSERT INTO " + input.TableName + " (" +
		"	grade, updated_client, created_client, " +
		"	created_at, created_by, updated_at, updated_by, description) " +
		"VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 ) " +
		" RETURNING id"
	params := []interface{}{
		grade.Grade.String, grade.UpdatedClient.String, grade.CreatedClient,
		grade.CreatedAt.Time, grade.CreatedBy.Int64, grade.UpdatedAt.Time,
		grade.UpdatedBy.Int64, grade.Description.String,
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

func (input employeeGradeDAO) GetEmployeeGrade(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) ([]interface{}, errorModel.ErrorModel) {
	query := `SELECT 
			id, grade,
			created_at,
			created_by,
			updated_at,
			updated_by, description 
		FROM ` + input.TableName + " "

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var grade repository.EmployeeGradeModel

		dbError := rows.Scan(
			&grade.ID,
			&grade.Grade,
			&grade.CreatedAt,
			&grade.CreatedBy,
			&grade.UpdatedAt,
			&grade.UpdatedBy, &grade.Description)

		return grade, dbError
	}

	return GetListDataDAO.GetListData(db, query, userParam, searchBy, createdBy, mappingFunc, "")
}

func (input employeeGradeDAO) GetEmployeeGradeMatrix(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) ([]interface{}, errorModel.ErrorModel) {
	var (
		query       string
		addWhere    []string
		getListData getListJoinDataDAO
	)

	query = fmt.Sprintf(`
		SELECT eg.id, eg.grade, eg.description, eg.updated_at FROM %s efa `,
		EmployeeFacilitiesActiveDAO.TableName)

	//--- Add Where Active Only
	addWhere = append(addWhere, " efa.active IS TRUE ")

	//--- Search By
	for i := 0; i < len(searchBy); i++ {
		if searchBy[i].SearchKey == "level_id" {
			searchBy[i].SearchKey = "el.id"
			continue
		}
		searchBy[i].SearchKey = "eg." + searchBy[i].SearchKey
	}

	//--- Order By
	userParam.OrderBy = "eg." + userParam.OrderBy

	//--- Created By
	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "eg.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     constanta.Filter,
		})
	}

	//--- Group By
	groupBy := fmt.Sprintf(` GROUP BY eg.id, eg.grade `)

	//--- Get List Join Date
	getListData = getListJoinDataDAO{Table: "efa", Query: query, AdditionalWhere: addWhere, GroupBy: groupBy}
	getListData.InnerJoinAlias(EmployeeLevelDAO.TableName, "el", "efa.employee_level_id", "el.id")
	getListData.InnerJoinAlias(EmployeeGradeDAO.TableName, "eg", "efa.employee_grade_id", "eg.id")

	//--- Mapping
	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var grade repository.EmployeeGradeModel
		dbError := rows.Scan(
			&grade.ID, &grade.Grade, &grade.Description,
			&grade.UpdatedAt)
		return grade, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input employeeGradeDAO) UpdateEmployeeGrade(db *sql.Tx, grade repository.EmployeeGradeModel) errorModel.ErrorModel {
	funcName := "UpdateEmployeeGrade"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	grade = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4, description = $5 " +
		"WHERE " +
		"	id = $6 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		grade.Grade.String,
		grade.UpdatedClient.String,
		grade.UpdatedAt.Time,
		grade.UpdatedBy.Int64,
		grade.Description.String,
		grade.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeGradeDAO) DeleteEmployeeGrade(db *sql.Tx, grade repository.EmployeeGradeModel) errorModel.ErrorModel {
	funcName := "DeleteEmployeeGrade"

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
		grade.Deleted.Bool,
		grade.UpdatedClient.String,
		grade.UpdatedAt.Time,
		grade.UpdatedBy.Int64,
		grade.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeGradeDAO) GetDetailEmployeeGrade(db *sql.Tx, id int64) (grade repository.EmployeeGradeModel, err errorModel.ErrorModel) {
	funcName := "GetDetailEmployeeGrade"
	query := "SELECT id, updated_at " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND id = $1 "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&grade.ID, &grade.UpdatedAt)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeGradeDAO) GetCountEmployeeGrade(db *sql.DB) (gradeCount int64, err errorModel.ErrorModel) {
	funcName := "GetCountEmployeeGrade"
	query := "SELECT COUNT(*) " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE "

	param := []interface{}{}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&gradeCount)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeGradeDAO) GetCountGradeMatrix(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	var (
		query       string
		getListData getListJoinDataDAO
		addWhere    []string
	)

	query = fmt.Sprintf(`SELECT COUNT(a.id) FROM (SELECT eg.id FROM %s efa `, EmployeeFacilitiesActiveDAO.TableName)

	//--- Add Where Active Only
	addWhere = append(addWhere, " efa.active IS TRUE ")

	//--- Search By
	for i := 0; i < len(searchBy); i++ {
		if searchBy[i].SearchKey == "level_id" {
			searchBy[i].SearchKey = "el.id"
			continue
		}
		searchBy[i].SearchKey = "eg." + searchBy[i].SearchKey
	}

	//--- Created By
	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "eg.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     constanta.Filter,
		})
	}

	//--- Group By
	groupBy := fmt.Sprintf(` GROUP BY eg.id, eg.grade) a `)

	//--- Get List Join Date
	getListData = getListJoinDataDAO{Table: "efa", Query: query, AdditionalWhere: addWhere, GroupBy: groupBy}
	getListData.InnerJoinAlias(EmployeeLevelDAO.TableName, "el", "efa.employee_level_id", "el.id")
	getListData.InnerJoinAlias(EmployeeGradeDAO.TableName, "eg", "efa.employee_grade_id", "eg.id")
	return getListData.GetCountJoinData(db, searchBy, 0)
}

func (input employeeGradeDAO) CheckGrade(db *sql.Tx, key string) (id int64, err errorModel.ErrorModel) {
	funcName := "CheckGrade"
	query := "SELECT " +
		"	id FROM " + input.TableName + " " +
		" WHERE LOWER(grade) = LOWER($1) AND deleted = FALSE LIMIT 1 "

	param := []interface{}{key}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}
