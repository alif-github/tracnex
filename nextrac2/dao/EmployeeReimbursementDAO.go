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

type employeeReimbursementDAO struct {
	AbstractDAO
}

var EmployeeReimbursementDAO = employeeReimbursementDAO{}.New()

func (input employeeReimbursementDAO) New() (output employeeReimbursementDAO) {
	output.FileName = "EmployeeReimbursementDAO.go"
	output.TableName = "employee_reimbursement"
	return
}

func (input employeeReimbursementDAO) InsertTx(tx *sql.Tx, model repository.EmployeeReimbursement) (id int64, errModel errorModel.ErrorModel) {
	funcName := "InsertTx"

	query := `INSERT INTO ` + input.TableName + `(
				name, receipt_no, benefit_id, 
				description, date, value, 
				status, file_upload_id, employee_id,
				created_at, created_by, created_client,
				updated_at, updated_by, updated_client,
				verified_status 
			) VALUES (
				$1, $2, $3,
				$4, $5, $6,
				$7, $8, $9,
				$10, $11, $12,
				$13, $14, $15,
				$16
			) RETURNING id`

	params := []interface{}{
		model.Name.String, model.ReceiptNo.String, model.BenefitId.Int64,
		model.Description.String, model.Date.Time, model.Value.Float64,
		model.Status.String, model.FileUploadId.Int64, model.EmployeeId.Int64,
		model.CreatedAt.Time, model.CreatedBy.Int64, model.CreatedClient.String,
		model.UpdatedAt.Time, model.UpdatedBy.Int64, model.UpdatedClient.String,
		model.VerifiedStatus.String,
	}

	row := tx.QueryRow(query, params...)
	if err := row.Scan(&id); err != nil {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return

}

func (input employeeReimbursementDAO) GetByIdForUpdate(db *sql.Tx, id int64, status string) (result repository.EmployeeReimbursement, errModel errorModel.ErrorModel) {
	funcName := "GetByIdForUpdate"

	query := `SELECT 
				er.id, er.employee_id, er.status, 
				er.value, er.benefit_id, er.updated_at,
				er."date", e.first_name, e.last_name,
				e.email, er.created_at, u.client_id 
			FROM ` + input.TableName + ` AS er 
			LEFT JOIN employee AS e 
				ON er.employee_id = e.id
			LEFT JOIN "user" AS u 
				ON ((u.email IS NOT NULL OR u.email != '') AND e.email = u.email)
				OR ((u.email IS NULL OR u.email = '') AND e.phone = u.phone) 
			WHERE 
				er.id = $1 AND
				er.status = $2 AND 
				er.deleted = FALSE`

	row := db.QueryRow(query, id, status)
	err := row.Scan(
		&result.ID, &result.EmployeeId, &result.Status,
		&result.Value, &result.BenefitId, &result.UpdatedAt,
		&result.Date, &result.Firstname, &result.Lastname,
		&result.Email, &result.CreatedAt, &result.ClientId)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeReimbursementDAO) UpdateStatusTx(tx *sql.Tx, model repository.EmployeeReimbursement) (errModel errorModel.ErrorModel) {
	funcName := "UpdateStatusTx"

	query := `UPDATE ` + input.TableName + ` 
			SET 
				status = $1, 
				updated_at = $2, 
				updated_by = $3, 
				updated_client = $4 
			WHERE 
				id = $5 AND 
				deleted = FALSE`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(
		model.Status.String,
		model.UpdatedAt.Time,
		model.UpdatedBy.Int64,
		model.UpdatedClient.String,
		model.ID.Int64,
	)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeReimbursementDAO) InitiateGetListEmployeeReimbursement(db *sql.DB, searchBy []in.SearchByParam, employeeReimbursement repository.EmployeeReimbursement) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT
				COUNT(*)
			FROM employee_reimbursement AS er
			LEFT JOIN employee AS e
				ON er.employee_id = e.id
			LEFT JOIN employee_benefits AS eb
				ON e.id = eb.employee_id
			LEFT JOIN file_upload AS fu
				ON er.file_upload_id = fu.id
			LEFT JOIN department AS d
				ON e.department_id = d.id`

	addQuery, params := input.getEmployeeReimbursementAddQuery(employeeReimbursement)

	return GetListDataDAO.GetCountData(db, params, query, searchBy, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				FieldName: "created_by",
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "er.deleted",
			},
		})
}

func (input employeeReimbursementDAO) GetListEmployeeReimbursement(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, employeeReimbursement repository.EmployeeReimbursement) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT
				er.id, e.id_card, e.first_name,
				e.last_name, d.name AS department, eb.current_medical_value,
				er.receipt_no, er.value, er.approved_value,
				er.status, fu.host, fu.path,
				er.note, er.verified_status, er.created_at,
				er.updated_at, e.id
			FROM employee_reimbursement AS er
			LEFT JOIN employee AS e
				ON er.employee_id = e.id
			LEFT JOIN employee_benefits AS eb
				ON e.id = eb.employee_id
			LEFT JOIN file_upload AS fu
				ON er.file_upload_id = fu.id
			LEFT JOIN department AS d
				ON e.department_id = d.id`

	addQuery, params := input.getEmployeeReimbursementAddQuery(employeeReimbursement)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.EmployeeReimbursement

			err := rows.Scan(
				&model.ID, &model.IDCard, &model.Firstname,
				&model.Lastname, &model.Department, &model.CurrentMedicalValue,
				&model.ReceiptNo, &model.Value, &model.ApprovedValue,
				&model.Status, &model.Host, &model.Path,
				&model.Note, &model.VerifiedStatus, &model.CreatedAt,
				&model.UpdatedAt, &model.EmployeeId,
			)

			return model, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				FieldName: "created_by",
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "er.deleted",
			},
		})
}

func (input employeeReimbursementDAO) getEmployeeReimbursementAddQuery(model repository.EmployeeReimbursement) (result string, params []interface{}) {
	if model.SearchBy.String != "" {
		if model.SearchBy.String == "id_card" {
			result += fmt.Sprintf(" AND e.id_card ilike $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.Keyword.String))
		}

		if model.SearchBy.String == "name" {
			result += fmt.Sprintf(" AND CONCAT(e.first_name, ' ', e.last_name) ilike $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.Keyword.String))
		}

		if model.SearchBy.String == "department" {
			result += fmt.Sprintf(" AND d.name ilike $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.Keyword.String))
		}
	}

	if model.FullName.String != "" {
		result += fmt.Sprintf(" AND CONCAT(e.first_name, ' ', e.last_name) ilike $%d", len(params)+1)
		params = append(params, fmt.Sprintf("%%%s%%", model.FullName.String))
	}

	if model.IDCard.String != "" {
		result += fmt.Sprintf(" AND e.id_card ilike $%d", len(params)+1)
		params = append(params, fmt.Sprintf("%%%s%%", model.IDCard.String))
	}

	if model.Department.String != "" {
		result += fmt.Sprintf(" AND d.name ilike $%d", len(params)+1)
		params = append(params, fmt.Sprintf("%%%s%%", model.Department.String))
	}

	if model.Status.String != "" {
		result += fmt.Sprintf(" AND er.status = $%d", len(params)+1)
		params = append(params, model.Status.String)
	}

	if model.VerifiedStatus.String != "" {
		result += fmt.Sprintf(" AND er.verified_status = $%d", len(params)+1)
		params = append(params, model.VerifiedStatus.String)
	}

	if model.StartDate.String != "" && model.EndDate.String != "" {
		result += fmt.Sprintf(" AND DATE(er.created_at) BETWEEN $%d AND $%d", len(params)+1, len(params)+2)
		params = append(params, model.StartDate.String, model.EndDate.String)
	}

	if model.Year.String != "" && model.Month.String != "" && model.IsFilter.Bool {
		result += fmt.Sprintf(" AND EXTRACT(year from er.created_at) = $%d AND EXTRACT(month from er.created_at) = $%d ", len(params)+1, len(params)+2)
		params = append(params, model.Year.String, model.Month.String)
	}

	if model.Year.String != "" && model.Month.String == "" && model.IsFilter.Bool {
		result += fmt.Sprintf(" AND EXTRACT(year from er.created_at) = $%d ", len(params)+1)
		params = append(params, model.Year.String)
	}

	return
}

func (input employeeReimbursementDAO) UpdateVerifiedStatus(db *sql.Tx, er repository.EmployeeReimbursement) errorModel.ErrorModel {
	funcName := "UpdateVerifiedStatus"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	verified_status = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4, note = $5, approved_value = $6 " +
		"WHERE " +
		"	id = $7 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		er.VerifiedStatus.String,
		er.UpdatedClient.String,
		er.UpdatedAt.Time,
		er.UpdatedBy.Int64,
		er.Note.String,
		er.ApprovedValue.Float64,
		er.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeReimbursementDAO) GetDetailReimbursementForVerify(db *sql.Tx, id int64) (reim repository.EmployeeReimbursement, err errorModel.ErrorModel) {
	funcName := "GetDetailReimbursementForVerify"
	query := `SELECT 
			er.id, er.updated_at, er.employee_id, er.value,
			er.created_at, e.first_name, e.last_name, u.client_id,
			e.email 
		 FROM ` + input.TableName + ` AS er
		 LEFT JOIN employee AS e 
			ON er.employee_id = e.id 
		 LEFT JOIN "user" AS u 
			ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
			OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
		 WHERE 
			er.id = $1 AND
			er.status = $2 AND 
			er.deleted = FALSE`

	param := []interface{}{id, constanta.ApprovedRequestStatus}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(
		&reim.ID, &reim.UpdatedAt, &reim.EmployeeId, &reim.Value,
		&reim.CreatedAt, &reim.Firstname, &reim.Lastname, &reim.ClientId,
		&reim.Email)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeReimbursementDAO) GetListEmployeeReimbursementReport(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, employeeReimbursement repository.EmployeeReimbursement) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT
				e.id, e.id_card, e.first_name, e.last_name, d.name AS department,
                e.date_join, e.date_out, eb.current_medical_value, eb.last_medical_value
			FROM  employee AS e
			LEFT JOIN employee_benefits AS eb
				ON e.id = eb.employee_id
			LEFT JOIN department AS d
				ON e.department_id = d.id`

	addQuery, params := input.getEmployeeReimbursementAddQuery(employeeReimbursement)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.EmployeeReimbursement

			err := rows.Scan(
				&model.ID, &model.IDCard, &model.Firstname,
				&model.Lastname, &model.Department,
				&model.DateJoin, &model.DateOut,
				&model.CurrentMedicalValue, &model.LastMedicalValue,
			)

			monthlyReport, monthlyReportArr, _ := input.GetDetailMonthlyReportByEmployeeId(db, model.ID.Int64, employeeReimbursement.Year.String)
			model.MonthlyReport = monthlyReport
			model.MonthlyReportArr = monthlyReportArr

			if employeeReimbursement.Year.String != ""{
				medic, _ := EmployeeHistoryLeaveDAO.GetHistoryLevelByEmployeeId(db, model.ID.Int64, employeeReimbursement.Year.String)
				model.CurrentMedicalValue.Float64 = medic.CurrentMedicalValue.Float64
				model.LastMedicalValue.Float64 = medic.LastMedicalValue.Float64
			}

			return model, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				FieldName: "created_by",
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "e.deleted",
			},
		})
}

func (input employeeReimbursementDAO) GetDetailMonthlyReportByEmployeeId(db *sql.DB, id int64, year string) (result repository.EmployeeReimbursementMonthlyReport, resultArr [12]float64, err errorModel.ErrorModel) {
	funcName := "GetDetailMonthlyReportByEmployeeId"
	query := "SELECT (SELECT SUM(approved_value) FROM employee_reimbursement"+
		      " WHERE EXTRACT(year FROM created_at)='" + year + "'" +
		" AND verified_status = '"+constanta.VerifiedReimbursementVerification+"' "+
		" AND employee_id=$1 AND deleted=FALSE) AS total "

	for i:=1; i<=12;i++  {
		query += "," +
			    " (SELECT SUM(approved_value) FROM employee_reimbursement"+
		        " WHERE EXTRACT(month FROM created_at)='" + strconv.Itoa(i) +"' AND" +
		        " EXTRACT(year FROM created_at)='" + year + "' " +
			    " AND verified_status = '"+constanta.VerifiedReimbursementVerification+"' " +
			    " AND employee_id=$1 AND deleted=FALSE) "
	}

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&result.Total, &result.January, &result.February, &result.March,
		&result.April,&result.May, &result.June,&result.July, &result.August,
		&result.September,&result.October,&result.November, &result.December)

	resultArr[0] = result.January.Float64
	resultArr[1] = result.February.Float64
	resultArr[2] = result.March.Float64
	resultArr[3] = result.April.Float64
	resultArr[4] = result.May.Float64
	resultArr[5] = result.June.Float64
	resultArr[6] = result.July.Float64
	resultArr[7] = result.August.Float64
	resultArr[8] = result.September.Float64
	resultArr[9] = result.October.Float64
	resultArr[10] = result.November.Float64
	resultArr[11] = result.December.Float64

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeReimbursementDAO) GetDetailReportReimbursement(db *sql.DB, year string, month int64, model repository.EmployeeReimbursement) (results []repository.EmployeeReimbursement, err errorModel.ErrorModel) {
	funcName := "GetDetailReportReimbursement"

	query := "SELECT e.id_card, e.first_name, e.last_name, " +
		" json_agg((er.receipt_no, er.approved_value, er.description)) AS item " +
	    " FROM employee_reimbursement er "+
		" LEFT JOIN employee e ON e.id = er.employee_id" +
		" WHERE er.deleted=FALSE AND " +
		" EXTRACT(month FROM er.created_at)='" + strconv.Itoa(int(month)) +"' AND" +
		" EXTRACT(year FROM er.created_at)='" + year + "' " +
		" AND er.verified_status = '"+constanta.VerifiedReimbursementVerification+"' "

	if model.FullName.String != "" {
		query += " AND CONCAT(e.first_name, ' ', e.last_name) ilike '%"+model.FullName.String+"%' "
	}

	if model.IDCard.String != "" {
		query += " AND e.id_card ilike '%"+model.IDCard.String+"%' "
	}

	if model.Department.String != "" {
		//query += " AND d.name ilike '"+model.FullName.String+"' "
	}

	if model.Status.String != "" {
		query += " AND er.status = "+model.Status.String+" "
	}

	if model.VerifiedStatus.String != "" {
		query += " AND er.verified_status = '"+model.VerifiedStatus.String+"%' "
	}

	query +=  " GROUP  BY e.id_card, e.first_name, e.last_name "

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
			var result repository.EmployeeReimbursement
			errorRows = rows.Scan(&result.IDCard,&result.Firstname,&result.Lastname,&result.Description)

			if errorRows != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
				return
			}
			results = append(results, result)
		}

	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
		return
	}
	return
}

func (input employeeReimbursementDAO) UpdateStatusAndCancellationReasonTx(tx *sql.Tx, model repository.EmployeeReimbursement) (errModel errorModel.ErrorModel) {
	funcName := "UpdateStatusAndCancellationReasonTx"

	query := `UPDATE ` + input.TableName + ` 
			SET 
				status = $1,
				cancellation_reason = $2, 
				updated_at = $3, 
				updated_by = $4, 
				updated_client = $5 
			WHERE 
				id = $6 AND 
				deleted = FALSE`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(
		model.Status.String,
		model.CancellationReason.String,
		model.UpdatedAt.Time,
		model.UpdatedBy.Int64,
		model.UpdatedClient.String,
		model.ID.Int64,
	)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeReimbursementDAO) InitiateGetListEmployeeReimbursementReport(db *sql.DB, searchBy []in.SearchByParam, employeeReimbursement repository.EmployeeReimbursement) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT
				COUNT(*)
			FROM  employee AS e
			LEFT JOIN employee_benefits AS eb
				ON e.id = eb.employee_id
			LEFT JOIN department AS d
				ON e.department_id = d.id`

	addQuery, params := input.getEmployeeReimbursementAddQuery(employeeReimbursement)

	return GetListDataDAO.GetCountData(db, params, query, searchBy, addQuery, DefaultFieldMustCheck{
		CreatedBy: FieldStatus{
			FieldName: "created_by",
			Value:     int64(0),
		},
		Deleted: FieldStatus{
			IsCheck: true,
			FieldName: "e.deleted",
		},
	})
}