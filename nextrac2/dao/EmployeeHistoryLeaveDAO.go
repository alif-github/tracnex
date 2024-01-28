package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type employeeHistoryLeaveDAO struct {
	AbstractDAO
}

var EmployeeHistoryLeaveDAO = employeeHistoryLeaveDAO{}.New()

func (input employeeHistoryLeaveDAO) New() (output employeeHistoryLeaveDAO) {
	output.TableName = "employee_history_leave"
	output.FileName = "EmployeeHistoryLeaveDAO.go"
	return
}

func (input employeeHistoryLeaveDAO) GetHistoryLevelByEmployeeId(db *sql.DB, idEmp int64, year string) (result repository.EmployeeLeaveModel, err errorModel.ErrorModel) {
	funcName := "GetHistoryLevelByEmployeeId"
	query := "SELECT " +
		"	id, current_annual_leave, last_annual_leave, " +
		"   current_medical_value, last_medical_value, year " +
		"   FROM " + input.TableName + " " +
		" WHERE employee_id = $1 AND deleted = FALSE "

	if year != "" {
		query += " AND year = '" + year + "' "
	} else {
		query += " AND year = (EXTRACT(year from NOW()))::varchar "
	}

	query += " LIMIT 1"

	param := []interface{}{idEmp}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&result.ID, &result.CurrentAnnualLeave, &result.LastAnnualLeave,
		&result.CurrentMedicalValue, &result.LastMedicalValue, &result.Year)
	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeHistoryLeaveDAO) GetByUserId(db *sql.DB, userId int64) (result repository.EmployeeBenefitsModel, errModel errorModel.ErrorModel) {
	funcName := "GetByUserId"

	query := `SELECT 
				eb.id, eb.current_annual_leave, eb.last_annual_leave 
			FROM ` + input.TableName + ` AS eb
			INNER JOIN ` + EmployeeDAO.TableName + ` AS e 
				ON eb.employee_id = e.id
			INNER JOIN "` + UserDAO.TableName + `" AS u
				ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
				OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
			WHERE 
				u.id = $1 AND
				eb.deleted = FALSE AND eb.year = (EXTRACT(year from NOW()))::varchar `

	row := db.QueryRow(query, userId)
	err := row.Scan(
		&result.ID, &result.CurrentAnnualLeave, &result.LastAnnualLeave)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeHistoryLeaveDAO) UpdateCurrentAnnualLeaveAndLastAnnualLeaveByEmployeeIdTx(tx *sql.Tx, model repository.EmployeeBenefitsModel) errorModel.ErrorModel {
	funcName := "UpdateCurrentAnnualLeaveAndLastAnnualLeaveByEmployeeIdTx"

	query := `UPDATE ` + input.TableName + ` 
    		SET
				current_annual_leave = $1,
				last_annual_leave = $2, 
				updated_at = $3, 
				updated_by = $4,
				updated_client = $5 
			WHERE 
				employee_id = $6 AND 
				deleted = FALSE AND
                year = (EXTRACT(year from NOW()))::varchar`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(
		&model.CurrentAnnualLeave.Int64,
		&model.LastAnnualLeave.Int64,
		&model.UpdatedAt.Time,
		&model.UpdatedBy.Int64,
		&model.UpdatedClient.String,
		&model.EmployeeID.Int64,
	)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeHistoryLeaveDAO) GetYearForFilter(db *sql.DB) (results []string, err errorModel.ErrorModel) {
	funcName := "GetYearForFilter"

	query := `SELECT DISTINCT(year) 
              FROM employee_history_leave 
              WHERE deleted=FALSE ORDER BY year ASC`

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
			var value repository.EmployeeLeaveModel
			rows.Scan(&value.Year)
			results = append(results, value.Year.String)
		}

	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
		return
	}
	return
}

func (input employeeHistoryLeaveDAO) InsertHistoryBenefits(db *sql.DB, inputStuct repository.EmployeeBenefitsModel) (lastInsertedId int64, err errorModel.ErrorModel) {
	funcName := "InsertHistoryBenefits"
	query := "INSERT INTO " + input.TableName + " (" +
		"	employee_id, current_annual_leave, last_annual_leave, " +
		"	current_medical_value, last_medical_value, year) " +
		"VALUES ( $1, $2, $3, $4, $5, $6) " +
		" RETURNING id"
	params := []interface{}{
		inputStuct.ID.Int64, inputStuct.CurrentAnnualLeave.Int64, inputStuct.LastAnnualLeave.Int64,
		inputStuct.CurrentMedicalValue.Float64, inputStuct.LastMedicalValue.Float64,
		inputStuct.Year.String,
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

func (input employeeHistoryLeaveDAO) GetCountEmployeeByYear(db *sql.DB, year string) (result int64, err errorModel.ErrorModel) {
	funcName := "GetCountEmployeeGrade"
	query := "SELECT COUNT(*) " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND year = $1 "

	param := []interface{}{year}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeHistoryLeaveDAO) UpdateCutOffValueHistory(tx *sql.DB, data repository.EmployeeBenefitsModel) errorModel.ErrorModel {
	funcName := "UpdateCutOffValueHistory"

	query := `UPDATE ` + input.TableName + `
			SET note_cutoff = $1,
                cutoff_value = $2
			WHERE deleted = FALSE
                 AND employee_id = $3 
                 AND year = $4 `

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(&data.NoteCutOff.String,
		               &data.CutOffLeaveValue.Int64,
		               &data.EmployeeID.Int64,
		               &data.Year.String)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}