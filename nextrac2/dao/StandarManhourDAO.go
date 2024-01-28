package dao

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type StandarManhourDAOInterface interface {
	InsertStandarManhour(*sql.Tx, repository.StandarManhourModel) (int64, errorModel.ErrorModel)
	GetStandarManhourForUpdate(*sql.DB, repository.StandarManhourModel) (repository.StandarManhourModel, errorModel.ErrorModel)
	UpdateStandarManhour(*sql.Tx, repository.StandarManhourModel) errorModel.ErrorModel
	DeleteStandarManhour(*sql.Tx, repository.StandarManhourModel) errorModel.ErrorModel
	GetListStandarManhour(*sql.DB, in.GetListDataDTO, []in.SearchByParam, int64) ([]interface{}, errorModel.ErrorModel)
	GetCountStandarManhour(*sql.DB, []in.SearchByParam, int64) (int, errorModel.ErrorModel)
	ViewStandarManhour(*sql.DB, repository.StandarManhourModel) (repository.StandarManhourModel, errorModel.ErrorModel)
}

type standarManhourDAO struct {
	AbstractDAO
}

var StandarManhourDAO = standarManhourDAO{}.New()

func (input standarManhourDAO) New() (output standarManhourDAO) {
	output.FileName = "StandarManhourDAO.go"
	output.TableName = "standard_manhour"
	return
}

func (input standarManhourDAO) InsertStandarManhour(db *sql.Tx, userParam repository.StandarManhourModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertStandarManhour"
		query    string
	)

	query = fmt.Sprintf(`INSERT INTO %s 
		(
		"case", department_id, standard_time, 
		created_by, created_at, created_client,
		updated_by, updated_at, updated_client
		)
		VALUES
		($1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9)
		RETURNING id `,
		input.TableName)

	params := []interface{}{
		userParam.Case.String, userParam.DepartmentID.Int64, userParam.Manhour.Float64,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		if pqError, ok := dbError.(*pgconn.PgError); ok {
			if pqError.Code == "22003" || pqError.Message == "numeric field overflow" {
				err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.ManhourRule1, constanta.Manhour, "")
				return
			}
		}

		//--- Error Numeric Field Overflow
		if dbError.Error() == "ERROR: numeric field overflow (SQLSTATE 22003)" {
			err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.ManhourRule1, constanta.Manhour, "")
			return
		}

		//--- Error 500
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input standarManhourDAO) UpdateStandarManhour(db *sql.Tx, userParam repository.StandarManhourModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateStandarManhour"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET
		"case" = $1, department_id = $2, standard_time = $3,
		updated_at = $4, updated_client = $5, updated_by = $6
		WHERE 
		id = $7 `, input.TableName)

	params := []interface{}{
		userParam.Case.String, userParam.DepartmentID.Int64, userParam.Manhour.Float64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.UpdatedBy.Int64,
		userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		if pqError, ok := dbError.(*pgconn.PgError); ok {
			if pqError.Code == "22003" || pqError.Message == "numeric field overflow" {
				err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.ManhourRule1, constanta.Manhour, "")
				return
			}
		}

		//-- Error Numeric Field Overflow
		if dbError.Error() == "ERROR: numeric field overflow (SQLSTATE 22003)" {
			err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.ManhourRule1, constanta.Manhour, "")
			return
		}

		//-- Error 500
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		if pqError, ok := dbError.(*pgconn.PgError); ok {
			if pqError.Code == "22003" || pqError.Message == "numeric field overflow" {
				err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.ManhourRule1, constanta.Manhour, "")
				return
			}
		}

		//-- Error Numeric Field Overflow
		if dbError.Error() == "ERROR: numeric field overflow (SQLSTATE 22003)" {
			err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.ManhourRule1, constanta.Manhour, "")
			return
		}

		//-- Error 500
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourDAO) DeleteStandarManhour(db *sql.Tx, userParam repository.StandarManhourModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteStandarManhour"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET
		deleted = $1, updated_by = $2, updated_at = $3, 
		updated_client = $4
		WHERE
		id = $5 `,
		input.TableName)

	param := []interface{}{
		true, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
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

func (input standarManhourDAO) GetStandarManhourForUpdate(db *sql.DB, userParam repository.StandarManhourModel) (result repository.StandarManhourModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetStandarManhourForUpdate"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT id, updated_at, created_by 
		FROM %s 
		WHERE id = $1 AND deleted = FALSE 
		FOR UPDATE`,
		input.TableName)

	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourDAO) GetListStandarManhour(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(`
		SELECT 
		sm.id, sm.case, d.name, 
		sm.standard_time, sm.updated_at 
		FROM %s sm 
		INNER JOIN %s d ON d.id = sm.department_id `,
		input.TableName, DepartmentDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.StandarManhourModel
			dbError := rows.Scan(
				&temp.ID, &temp.Case, &temp.Department,
				&temp.Manhour, &temp.UpdatedAt,
			)
			return temp, dbError
		}, "", input.getDefaultMustCheck("sm", createdBy))
}

func (input standarManhourDAO) GetCountStandarManhour(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	for i := 0; i < len(searchByParam); i++ {
		(searchByParam)[i].SearchKey = "sm." + (searchByParam)[i].SearchKey
	}

	tableName := fmt.Sprintf(`
		%s sm
		INNER JOIN %s d ON d.id = sm.department_id`,
		input.TableName, DepartmentDAO.TableName)

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, "", input.getDefaultMustCheck("sm", createdBy))
}

func (input standarManhourDAO) ViewStandarManhour(db *sql.DB, userParam repository.StandarManhourModel) (result repository.StandarManhourModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewStandarManhour"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
		sm.id, sm.case, sm.department_id, 
		d.name, sm.standard_time, sm.updated_at, 
		sm.created_by
		FROM %s sm 
		INNER JOIN %s d ON d.id = sm.department_id 
		WHERE sm.id = $1 AND sm.deleted = FALSE `,
		input.TableName, DepartmentDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)

	dbError := results.Scan(
		&result.ID, &result.Case, &result.DepartmentID,
		&result.Department, &result.Manhour, &result.UpdatedAt,
		&result.CreatedBy,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourDAO) GetStandarManhourByCase(db *sql.DB, departmentID int64, categoryArray []string) (result repository.StandarManhourModel, err errorModel.ErrorModel) {
	var (
		funcName  = "GetStandarManhourByCase"
		query     string
		queryTemp string
	)

	for idx, itemCategoryArray := range categoryArray {
		if idx == 0 {
			queryTemp += fmt.Sprintf(` AND "case" IN (`)
		}

		if len(categoryArray)-(idx+1) == 0 {
			queryTemp += fmt.Sprintf(`'%s')`, itemCategoryArray)
		} else {
			queryTemp += fmt.Sprintf(`'%s',`, itemCategoryArray)
		}
	}

	query = fmt.Sprintf(`
		SELECT SUM(standard_time) AS manhour 
		FROM %s WHERE 
		department_id = $1 AND deleted = FALSE `,
		input.TableName)

	if !util.IsStringEmpty(queryTemp) {
		query += queryTemp
	}

	params := []interface{}{departmentID}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.Manhour)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "sm." + (*searchByParam)[i].SearchKey
	}

	switch userParam.OrderBy {
	case "id", "id ASC", "id DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "sm.id " + strSplit[1]
		} else {
			userParam.OrderBy = "sm.id"
		}
	case "department_id", "department_id ASC", "department_id DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "sm.department_id " + strSplit[1]
		} else {
			userParam.OrderBy = "sm.department_id"
		}
	case "case", "case ASC", "case DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "sm.case " + strSplit[1]
		} else {
			userParam.OrderBy = "sm.case"
		}
	case "updated_at", "updated_at ASC", "updated_at DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "sm.updated_at " + strSplit[1]
		} else {
			userParam.OrderBy = "sm.updated_at"
		}
	default:
	}
}

func (input standarManhourDAO) getDefaultMustCheck(tableAlias string, createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: tableAlias + ".id"},
		Deleted:   FieldStatus{FieldName: tableAlias + ".deleted"},
		CreatedBy: FieldStatus{FieldName: tableAlias + ".created_by", Value: createdBy},
	}
}
