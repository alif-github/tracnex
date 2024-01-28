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

type employeeLevelDAO struct {
	AbstractDAO
}

var EmployeeLevelDAO = employeeLevelDAO{}.New()

func (input employeeLevelDAO) New() (output employeeLevelDAO) {
	output.TableName = "employee_level"
	output.FileName = "EmployeeLevelDAO.go"
	return
}

func (input employeeLevelDAO) InsertEmployeeLevel(db *sql.Tx, level repository.EmployeeLevelModel) (lastInsertedId int64, err errorModel.ErrorModel) {
	funcName := "InsertEmployeeLevel"
	query := "INSERT INTO " + input.TableName + " (" +
		"	level, description, updated_client, created_client, " +
		"	created_at, created_by, updated_at, updated_by) " +
		"VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 ) " +
		" RETURNING id"
	params := []interface{}{
		level.Level.String, level.Description.String, level.UpdatedClient.String, level.CreatedClient,
		level.CreatedAt.Time, level.CreatedBy.Int64, level.UpdatedAt.Time,
		level.UpdatedBy.Int64,
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

func (input employeeLevelDAO) GetEmployeeLevel(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) ([]interface{}, errorModel.ErrorModel) {
	query := `SELECT 
			id, level,
			created_at,
			created_by,
			updated_at,
			updated_by, description 
		FROM ` + input.TableName + " "

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var level repository.EmployeeLevelModel

		dbError := rows.Scan(
			&level.ID,
			&level.Level,
			&level.CreatedAt,
			&level.CreatedBy,
			&level.UpdatedAt,
			&level.UpdatedBy, &level.Description)

		return level, dbError
	}

	return GetListDataDAO.GetListData(db, query, userParam, searchBy, createdBy, mappingFunc, "")
}

func (input employeeLevelDAO) GetEmployeeLevelMatrix(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) ([]interface{}, errorModel.ErrorModel) {
	var (
		query       string
		addWhere    []string
		getListData getListJoinDataDAO
	)

	query = fmt.Sprintf(`
		SELECT el.id, el."level", el.description, el.updated_at FROM %s efa `,
		EmployeeFacilitiesActiveDAO.TableName)

	//--- Add Where Active Only
	addWhere = append(addWhere, " efa.active IS TRUE ")

	//--- Search By
	for i := 0; i < len(searchBy); i++ {
		searchBy[i].SearchKey = "el." + searchBy[i].SearchKey
	}

	//--- Order By
	userParam.OrderBy = "el." + userParam.OrderBy

	//--- Created By
	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "el.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     constanta.Filter,
		})
	}

	//--- Group By
	groupBy := fmt.Sprintf(` GROUP BY el.id, el."level" `)

	//--- Get List Join Date
	getListData = getListJoinDataDAO{Table: "efa", Query: query, AdditionalWhere: addWhere, GroupBy: groupBy}
	getListData.InnerJoinAlias(EmployeeLevelDAO.TableName, "el", "efa.employee_level_id", "el.id")
	getListData.InnerJoinAlias(EmployeeGradeDAO.TableName, "eg", "efa.employee_grade_id", "eg.id")

	//--- Mapping
	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var level repository.EmployeeLevelModel
		dbError := rows.Scan(
			&level.ID, &level.Level, &level.Description,
			&level.UpdatedAt)
		return level, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input employeeLevelDAO) GetCountLevelMatrix(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	var (
		query       string
		getListData getListJoinDataDAO
		addWhere    []string
	)

	query = fmt.Sprintf(`SELECT COUNT(a.id) FROM (SELECT el.id FROM %s efa `, EmployeeFacilitiesActiveDAO.TableName)

	//--- Add Where Active Only
	addWhere = append(addWhere, " efa.active IS TRUE ")

	//--- Search By
	for i := 0; i < len(searchBy); i++ {
		searchBy[i].SearchKey = "el." + searchBy[i].SearchKey
	}

	//--- Created By
	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "el.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     constanta.Filter,
		})
	}

	//--- Group By
	groupBy := fmt.Sprintf(` GROUP BY el.id, el."level") a `)

	//--- Get List Join Date
	getListData = getListJoinDataDAO{Table: "efa", Query: query, AdditionalWhere: addWhere, GroupBy: groupBy}
	getListData.InnerJoinAlias(EmployeeLevelDAO.TableName, "el", "efa.employee_level_id", "el.id")
	getListData.InnerJoinAlias(EmployeeGradeDAO.TableName, "eg", "efa.employee_grade_id", "eg.id")
	return getListData.GetCountJoinData(db, searchBy, 0)
}

func (input employeeLevelDAO) UpdateEmployeeLevel(db *sql.Tx, level repository.EmployeeLevelModel) errorModel.ErrorModel {
	funcName := "UpdateEmployeeLevel"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	level = $1," +
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
		level.Level.String,
		level.UpdatedClient.String,
		level.UpdatedAt.Time,
		level.UpdatedBy.Int64,
		level.Description.String,
		level.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeLevelDAO) DeleteEmployeeLevel(db *sql.Tx, level repository.EmployeeLevelModel) errorModel.ErrorModel {
	funcName := "DeleteEmployeeLevel"

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
		level.Deleted.Bool,
		level.UpdatedClient.String,
		level.UpdatedAt.Time,
		level.UpdatedBy.Int64,
		level.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeLevelDAO) GetDetailEmployeeLevel(db *sql.Tx, id int64) (level repository.EmployeeLevelModel, err errorModel.ErrorModel) {
	funcName := "GetDetailEmployeeLevel"
	query := "SELECT id, updated_at " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND id = $1 "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&level.ID, &level.UpdatedAt)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLevelDAO) GetCountEmployeeLevel(db *sql.DB) (gradeCount int64, err errorModel.ErrorModel) {
	funcName := "GetCountEmployeeLevel"
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

func (input employeeLevelDAO) CheckLevel(db *sql.Tx, key string) (id int64, err errorModel.ErrorModel) {
	funcName := "CheckLevel"
	query := "SELECT " +
		"	id FROM " + input.TableName + " " +
		" WHERE LOWER(level) = LOWER($1) AND deleted = FALSE LIMIT 1 "

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
