package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type employeeLevelGradeDAO struct {
	AbstractDAO
}

var EmployeeLevelGradeDAO = employeeLevelGradeDAO{}.New()

func (input employeeLevelGradeDAO) New() (output employeeLevelGradeDAO) {
	output.FileName = "EmployeeLevelGradeDAO.go"
	return
}

func (input employeeLevelGradeDAO) CheckEmployeeLevel(db *sql.DB, userParam repository.EmployeeModel) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName  = "CheckEmployeeLevel"
		tableName = "employee_level"
		query     string
	)

	query = fmt.Sprintf(`
		SELECT CASE WHEN id > 0 THEN true ELSE false END is_exist
		FROM %s  
		WHERE id = $1 AND active = TRUE AND deleted = FALSE `,
		tableName)

	param := []interface{}{userParam.LevelID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLevelGradeDAO) GetNameEmployeeLevelByIDWithoutDeleted(db *sql.DB, userParam repository.EmployeeLevelModel) (result repository.EmployeeLevelModel, err errorModel.ErrorModel) {
	var (
		funcName  = "GetNameEmployeeLevelByIDWithoutDeleted"
		tableName = "employee_level"
		query     string
	)

	query = fmt.Sprintf(`SELECT level FROM %s WHERE id = $1 `, tableName)
	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result.Level)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLevelGradeDAO) CheckEmployeeGrade(db *sql.DB, userParam repository.EmployeeModel) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName  = "CheckEmployeeGrade"
		tableName = "employee_grade"
		query     string
	)

	query = fmt.Sprintf(`
		SELECT CASE WHEN id > 0 THEN true ELSE false END is_exist
		FROM %s  
		WHERE id = $1 AND active = TRUE AND deleted = FALSE `,
		tableName)

	param := []interface{}{userParam.GradeID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLevelGradeDAO) GetNameEmployeeGradeByIDWithoutDeleted(db *sql.DB, userParam repository.EmployeeGradeModel) (result repository.EmployeeGradeModel, err errorModel.ErrorModel) {
	var (
		funcName  = "GetNameEmployeeGradeByIDWithoutDeleted"
		tableName = "employee_grade"
		query     string
	)

	query = fmt.Sprintf(`SELECT grade FROM %s WHERE id = $1 `, tableName)
	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result.Grade)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
