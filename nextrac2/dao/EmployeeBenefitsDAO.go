package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type employeeBenefitsDAO struct {
	AbstractDAO
}

var EmployeeBenefitsDAO = employeeBenefitsDAO{}.New()

func (input employeeBenefitsDAO) New() (output employeeBenefitsDAO) {
	output.FileName = "EmployeeBenefitsDAO.go"
	output.TableName = "employee_benefits"
	return
}

func (input employeeBenefitsDAO) InsertEmployeeBenefits(db *sql.Tx, userParam repository.EmployeeBenefitsModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertEmployeeBenefits"
		query    string
	)

	query = fmt.Sprintf(
		`INSERT INTO %s
		(employee_id, employee_level_id, employee_grade_id,
		created_by, created_client, created_at, 
		updated_by, updated_client, updated_at,
		salary, bpjs_no, bpjs_tk_no, 
		current_annual_leave, last_annual_leave, vehicle_limit)
		VALUES
		($1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12, 
		$13, $14, $15)
		RETURNING id `, input.TableName)

	param := []interface{}{
		userParam.EmployeeID.Int64, userParam.EmployeeLevelID.Int64, userParam.EmployeeGradeID.Int64,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
	}

	if userParam.Salary.Float64 > 0 {
		param = append(param, userParam.Salary.Float64)
	} else {
		param = append(param, nil)
	}

	if userParam.BPJSNo.String != "" {
		param = append(param, userParam.BPJSNo.String)
	} else {
		param = append(param, nil)
	}

	if userParam.BPJSTkNo.String != "" {
		param = append(param, userParam.BPJSTkNo.String)
	} else {
		param = append(param, nil)
	}

	if userParam.CurrentAnnualLeave.Int64 > 0 {
		param = append(param, userParam.CurrentAnnualLeave.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.LastAnnualLeave.Int64 > 0 {
		param = append(param, userParam.LastAnnualLeave.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.VehicleLimit.Float64 > 0 {
		param = append(param, userParam.VehicleLimit.Float64)
	} else {
		param = append(param, nil)
	}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input employeeBenefitsDAO) DeleteEmployeeBenefitsByEmployeeID(db *sql.Tx, userParam repository.EmployeeBenefitsModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteEmployeeBenefitsByEmployeeID"
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

func (input employeeBenefitsDAO) UpdateEmployeeBenefitsByEmployeeID(db *sql.Tx, userParam repository.EmployeeBenefitsModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateEmployeeBenefitsByEmployeeID"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		employee_level_id = $1, employee_grade_id = $2, updated_by = $3, 
		updated_client = $4, updated_at = $5, salary = $6, 
		bpjs_no = $7, bpjs_tk_no = $8, current_annual_leave = $9, 
		last_annual_leave = $10, vehicle_limit = $11
		WHERE employee_id = $12 `,
		input.TableName)

	param := []interface{}{
		userParam.EmployeeLevelID.Int64, userParam.EmployeeGradeID.Int64, userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
	}

	if userParam.Salary.Float64 > 0 {
		param = append(param, userParam.Salary.Float64)
	} else {
		param = append(param, nil)
	}

	if userParam.BPJSNo.String != "" {
		param = append(param, userParam.BPJSNo.String)
	} else {
		param = append(param, nil)
	}

	if userParam.BPJSTkNo.String != "" {
		param = append(param, userParam.BPJSTkNo.String)
	} else {
		param = append(param, nil)
	}

	if userParam.CurrentAnnualLeave.Int64 > 0 {
		param = append(param, userParam.CurrentAnnualLeave.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.LastAnnualLeave.Int64 > 0 {
		param = append(param, userParam.LastAnnualLeave.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.VehicleLimit.Float64 > 0 {
		param = append(param, userParam.VehicleLimit.Float64)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.EmployeeID.Int64)

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

func (input employeeBenefitsDAO) GetByEmployeeIdForUpdate(db *sql.Tx, employeeId int64) (result repository.EmployeeBenefitsModel, errModel errorModel.ErrorModel) {
	var (
		funcName = "GetByEmployeeIdForUpdate"
	)

	query := `SELECT 
				id, current_annual_leave, last_annual_leave, 
				current_medical_value, last_medical_value, employee_level_id,
				employee_grade_id 
			FROM ` + input.TableName + `
			WHERE 
				employee_id = $1 AND 
				deleted = FALSE
			FOR UPDATE`

	row := db.QueryRow(query, employeeId)
	err := row.Scan(
		&result.ID, &result.CurrentAnnualLeave, &result.LastAnnualLeave,
		&result.CurrentMedicalValue, &result.LastMedicalValue, &result.EmployeeLevelID,
		&result.EmployeeGradeID)

	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeBenefitsDAO) GetByUserId(db *sql.DB, userId int64) (result repository.EmployeeBenefitsModel, errModel errorModel.ErrorModel) {
	funcName := "GetByUserId"

	query := `SELECT 
				eb.id, eb.current_annual_leave, eb.last_annual_leave, 
				eb.current_medical_value 
			FROM ` + input.TableName + ` AS eb
			INNER JOIN ` + EmployeeDAO.TableName + ` AS e 
				ON eb.employee_id = e.id
			INNER JOIN "` + UserDAO.TableName + `" AS u
				ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
				OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
			WHERE 
				u.id = $1 AND
				eb.deleted = FALSE`

	row := db.QueryRow(query, userId)
	err := row.Scan(
		&result.ID, &result.CurrentAnnualLeave, &result.LastAnnualLeave,
		&result.CurrentMedicalValue)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeBenefitsDAO) UpdateCurrentAnnualLeaveAndLastAnnualLeaveByEmployeeIdTx(tx *sql.Tx, model repository.EmployeeBenefitsModel) errorModel.ErrorModel {
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
				deleted = FALSE`

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

func (input employeeBenefitsDAO) UpdateCurrentMedicalValueAndLastMedicalValueByEmployeeIdTx(tx *sql.Tx, model repository.EmployeeBenefitsModel) errorModel.ErrorModel {
	funcName := "UpdateCurrentMedicalValueAndLastMedicalValueByEmployeeIdTx"

	query := `UPDATE ` + input.TableName + `
			SET 
				current_medical_value = $1,
				last_medical_value = $2, 
				updated_at = $3, 
				updated_by = $4,
				updated_client = $5 
			WHERE 
				employee_id = $6 AND 
				deleted = FALSE`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(
		&model.CurrentMedicalValue.Float64,
		&model.LastMedicalValue.Float64,
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

func (input employeeBenefitsDAO) GetDetailMedicalValueForVerify(db *sql.DB, id int64) (benefit repository.EmployeeBenefitsModel, err errorModel.ErrorModel) {
	funcName := "GetDetailMedicalValueForVerify"
	query := "SELECT id, updated_at, current_medical_value, last_medical_value, " +
		" current_annual_leave, last_annual_leave " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND employee_id = $1 "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&benefit.ID, &benefit.UpdatedAt,
		&benefit.CurrentMedicalValue, &benefit.LastMedicalValue,
		&benefit.CurrentAnnualLeave, &benefit.LastAnnualLeave)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeBenefitsDAO) UpdateLeaveAndMedicalValueDB(tx *sql.DB, model repository.EmployeeBenefitsModel) errorModel.ErrorModel {
	funcName := "UpdateLeaveAndMedicalValueDB"

	query := `UPDATE ` + input.TableName + `
			SET 
				current_medical_value = $1,
				last_medical_value = $2, 
                current_annual_leave = $3,
				last_annual_leave = $4, 
				updated_at = $5
			WHERE 
				employee_id = $6 AND 
				deleted = FALSE`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(
		&model.CurrentMedicalValue.Float64,
		&model.LastMedicalValue.Float64,
		&model.CurrentAnnualLeave.Int64,
		&model.LastAnnualLeave.Int64,
		&model.UpdatedAt.Time,
		&model.EmployeeID.Int64,
	)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeBenefitsDAO) ResetLastAnnualLeave(tx *sql.DB) errorModel.ErrorModel {
	funcName := "ResetLastAnnualLeave"

	query := `UPDATE ` + input.TableName + `
			SET last_annual_leave = $1
			WHERE deleted = FALSE`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(0)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}