package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type employeeHistoryDAO struct {
	AbstractDAO
}

var EmployeeHistoryDAO = employeeHistoryDAO{}.New()

func (input employeeHistoryDAO) New() (output employeeHistoryDAO) {
	output.FileName = "EmployeeHistoryDAO.go"
	return
}

func (input employeeHistoryDAO) InitiateEmployeeRequestHistory(db *sql.DB, searchBy []in.SearchByParam, employeeId int64) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT COUNT(*) FROM (
			(
				SELECT 
					created_at, deleted, employee_id 
				FROM employee_leave
			)
			UNION ALL
			(
				SELECT 
					created_at, deleted, employee_id
				FROM employee_reimbursement
			)
		) AS e`

	addQuery := " AND e.employee_id = $1"
	params := []interface{}{ employeeId }

	return GetListDataDAO.GetCountData(db, params, query, searchBy, addQuery, DefaultFieldMustCheck{
		Deleted:   FieldStatus{
			IsCheck:   true,
			FieldName: "e.deleted",
		},
		CreatedBy: FieldStatus{
			Value:     int64(0),
		},
	})
}

func (input employeeHistoryDAO) GetListEmployeeRequestHistory(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, employeeId int64) (result []interface{}, errModel errorModel.ErrorModel) {
	query := fmt.Sprintf(`SELECT * FROM
			(
				(
					SELECT 
						el.id, el.description, el.date AS leave_date,
						NULL AS date, el.type AS request_type, al.allowance_name AS "type", 
						el.status, NULL AS verified_status, el.cancellation_reason, 
						fu.host, fu.path, el.value AS total_leave, 
						NULL AS value, NULL AS approved_value, NULL AS note, 
						el.created_at, el.updated_at, el.created_by, 
						el.deleted, NULL AS receipt_no, el.employee_id
					FROM employee_leave AS el
					LEFT JOIN allowances AS al
						ON el.allowance_id = al.id
					LEFT JOIN file_upload AS fu
						ON el.file_upload_id = fu.id
				)
				UNION ALL
				(
					SELECT 
						er.id, er.description, NULL, 
						er.date AS date, '%s' AS request_type, b.benefit_name, 
						er.status, er.verified_status, er.cancellation_reason, 
						fu.host, fu.path, NULL AS total_leave, 
						er.value, er.approved_value, er.note, 
						er.created_at, er.updated_at, er.created_by, 
						er.deleted, er.receipt_no, er.employee_id 
					FROM employee_reimbursement AS er
					LEFT JOIN benefits AS b
						ON er.benefit_id = b.id
					LEFT JOIN file_upload AS fu
						ON er.file_upload_id = fu.id
				)
			) AS e `, constanta.ReimbursementType)

	addQuery := " AND e.employee_id = $1"
	params := []interface{}{ employeeId }

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db,params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var employeeRequestHistory repository.EmployeeRequestHistory

			err := rows.Scan(
				&employeeRequestHistory.ID, &employeeRequestHistory.Description, &employeeRequestHistory.LeaveDate,
				&employeeRequestHistory.Date, &employeeRequestHistory.RequestType, &employeeRequestHistory.Type,
				&employeeRequestHistory.Status, &employeeRequestHistory.VerifiedStatus, &employeeRequestHistory.CancellationReason,
				&employeeRequestHistory.Host, &employeeRequestHistory.Path, &employeeRequestHistory.TotalLeave,
				&employeeRequestHistory.Value, &employeeRequestHistory.ApprovedValue, &employeeRequestHistory.Note,
				&employeeRequestHistory.CreatedAt, &employeeRequestHistory.UpdatedAt, &employeeRequestHistory.CreatedBy,
				&employeeRequestHistory.Deleted, &employeeRequestHistory.ReceiptNo, &employeeRequestHistory.EmployeeId,
			)

			return employeeRequestHistory, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				FieldName: "created_by",
				Value:     createdBy,
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "e.deleted",
			},
		})
}

func (input employeeHistoryDAO) InitiateEmployeeApprovalHistoryByEmployeeIdList(db *sql.DB, searchBy []in.SearchByParam, employeeIdList []string, statuses []string) (result int, errModel errorModel.ErrorModel) {
	addQuery := input.getEmployeeApprovalHistoryAddQuery(employeeIdList, statuses)
	query := `SELECT COUNT(*) FROM (
			(
				SELECT 
					employee_id, created_at, deleted
				FROM employee_leave
			)
			UNION ALL
			(
				SELECT 
					employee_id, created_at, deleted
				FROM employee_reimbursement
			)
		) AS e`
	return GetListDataDAO.GetCountData(db, []interface{}{}, query, searchBy, addQuery, DefaultFieldMustCheck{
		Deleted:   FieldStatus{
			IsCheck:   true,
			FieldName: "e.deleted",
		},
		CreatedBy: FieldStatus{
			Value:     int64(0),
		},
	})
}

func (input employeeHistoryDAO) GetListEmployeeApprovalHistoryByEmployeeIdListAndStatus(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, employeeIdList []string, statuses []string) (result []interface{}, errModel errorModel.ErrorModel) {
	addQuery := input.getEmployeeApprovalHistoryAddQuery(employeeIdList, statuses)
	query := fmt.Sprintf(`SELECT * FROM (
			(
				SELECT 
					el.id, e.first_name, e.last_name,
					e.id_card, d.name AS department, el.type,
					al.allowance_name AS "type", el.date AS leave_date, NULL AS date,
					el.status, fu.host, fu.path, 
					eb.current_annual_leave + eb.last_annual_leave AS total_remaining_leave, NULL AS receipt_no, el.description, 
					el.created_at, el.updated_at, el.created_by, 
					el.deleted, e.id AS employee_id, el.value AS total_leave,
					NULL AS value, NULL AS approved_value, NULL AS verified_status, 
					el.cancellation_reason, NULL AS note
				FROM employee_leave AS el
				LEFT JOIN employee AS e
					ON el.employee_id = e.id
				LEFT JOIN department AS d
					ON e.department_id = d.id
				LEFT JOIN employee_benefits AS eb
					ON e.id = eb.employee_id
				LEFT JOIN allowances AS al
					ON el.allowance_id = al.id
				LEFT JOIN file_upload AS fu
					ON el.file_upload_id = fu.id
			)
			UNION ALL
			(
				SELECT 
					er.id, e.first_name, e.last_name,
					e.id_card, d.name AS department, '%s', 
					b.benefit_name AS "type", NULL AS leave_date, er.date, 
					er.status, fu.host, fu.path, 
					NULL, er.receipt_no, er.description, 
					er.created_at, er.updated_at, er.created_by,
					er.deleted, e.id AS employee_id, NULL,
					er.value, er.approved_value, er.verified_status, 
					er.cancellation_reason, er.note
				FROM employee_reimbursement AS er
				LEFT JOIN employee AS e
					ON er.employee_id = e.id
				LEFT JOIN department AS d
					ON e.department_id = d.id
				LEFT JOIN employee_benefits AS eb
					ON e.id = eb.employee_id
				LEFT JOIN benefits AS b
					ON er.benefit_id = b.id
				LEFT JOIN file_upload AS fu
					ON er.file_upload_id = fu.id
			)
		) AS e`, constanta.ReimbursementType)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.EmployeeApprovalHistory

			err := rows.Scan(
				&model.ID, &model.Firstname, &model.Lastname,
				&model.IDCard, &model.Department, &model.RequestType,
				&model.Type, &model.LeaveDate, &model.Date,
				&model.Status, &model.Host, &model.Path,
				&model.TotalRemainingLeave, &model.ReceiptNo, &model.Description,
				&model.CreatedAt, &model.UpdatedAt, &model.CreatedBy,
				&model.Deleted, &model.EmployeeId, &model.TotalLeave,
				&model.Value, &model.ApprovedValue, &model.VerifiedStatus,
				&model.CancellationReason, &model.Note)

			return model, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "e.deleted",
			},
		})
}

func (input employeeHistoryDAO) getEmployeeApprovalHistoryAddQuery(employeeIdList, statuses []string) string {
	addQuery := ""

	/*
		Status
	*/
	if statuses != nil {
		status := fmt.Sprintf("'%s'", strings.Join(statuses, "','"))
		addQuery += fmt.Sprintf(" AND e.status IN (%s)", status)
	}

	/*
		Employee Id
	*/
	employeeIdQuery := " AND e.employee_id IN (0)"

	if employeeIdList != nil {
		isAll := false

		for _, employeeId := range employeeIdList {
			if employeeId == "all" {
				employeeIdQuery = ""
				isAll = true

				break
			}
		}

		if !isAll {
			employeeIdQuery = fmt.Sprintf(" AND e.employee_id IN (%s)", strings.Join(employeeIdList, ","))
		}
	}

	addQuery += employeeIdQuery

	return addQuery
}