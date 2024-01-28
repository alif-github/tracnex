package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
	"time"
)

type EmployeeLeaveDAOInterface interface {
	InsertTx(tx *sql.Tx, model repository.EmployeeLeaveModel) (id int64, errModel errorModel.ErrorModel)
	GetById(db *sql.DB, id int64) (result repository.EmployeeLeaveModel, errModel errorModel.ErrorModel)
}

type employeeLeaveDAO struct {
	AbstractDAO
}

var EmployeeLeaveDAO = employeeLeaveDAO{}.New()

func (input employeeLeaveDAO) New() (output employeeLeaveDAO) {
	output.FileName = "EmployeeLeaveDAO.go"
	output.TableName = "employee_leave"
	return
}

func (input employeeLeaveDAO) InsertTx(tx *sql.Tx, model repository.EmployeeLeaveModel) (id int64, errModel errorModel.ErrorModel) {
	funcName := "InsertTx"

	query := `INSERT INTO ` + input.TableName + `(
				name, allowance_id, description, date, 
				value, status, file_upload_id, employee_id, 
				created_by, created_at, created_client, updated_by, 
				updated_at, updated_client, type 
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				$9, $10, $11, $12,
				$13, $14, $15
			) RETURNING id`

	fileUploadId := &model.FileUploadId.Int64
	allowanceId := &model.AllowanceId.Int64

	if model.FileUploadId.Int64 == 0 {
		fileUploadId = nil
	}

	if model.AllowanceId.Int64 == 0 {
		allowanceId = nil
	}

	params := []interface{}{
		model.Name.String, allowanceId, model.Description.String, model.Date.String,
		model.Value.Int64, model.Status.String, fileUploadId, model.EmployeeId.Int64,
		model.CreatedBy.Int64, model.CreatedAt.Time, model.CreatedClient.String, model.UpdatedBy.Int64,
		model.UpdatedAt.Time, model.CreatedClient.String, model.Type.String,
	}

	row := tx.QueryRow(query, params...)
	if err := row.Scan(&id); err != nil {
		return 0, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return id, errorModel.GenerateNonErrorModel()
}

func (input employeeLeaveDAO) GetByIdAndStatuses(db *sql.Tx, id int64, statuses []string) (result repository.EmployeeLeaveModel, errModel errorModel.ErrorModel) {
	funcName := "GetByIdAndStatuses"

	query := `SELECT 
				el.id, el.employee_id, el.status,
				el.value, el.allowance_id, el.updated_at,
				al.allowance_type, el.cancellation_reason, el.type,
				el.date, u.client_id, e.email,
				e.first_name, e.last_name, el.created_at 
			FROM ` + input.TableName + ` AS el
			LEFT JOIN allowances AS al
				ON el.allowance_id = al.id
			LEFT JOIN employee as e 
				on el.employee_id  = e.id
			LEFT JOIN "user" as u 
				ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
				OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
			WHERE 
				el.id = $1 AND
				el.deleted = FALSE`

	if statuses != nil {
		query += fmt.Sprintf(" AND el.status IN ('%s') ", strings.Join(statuses, "','"))
	}

	row := db.QueryRow(query, id)
	err := row.Scan(
		&result.ID, &result.EmployeeId, &result.Status,
		&result.Value, &result.AllowanceId, &result.UpdatedAt,
		&result.AllowanceType, &result.CancellationReason, &result.Type,
		&result.Date, &result.ClientID, &result.Email,
		&result.Firstname, &result.Lastname, &result.CreatedAt)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLeaveDAO) GetByIdAndStatusesTx(db *sql.Tx, id int64, statuses []string) (result repository.EmployeeLeaveModel, errModel errorModel.ErrorModel) {
	funcName := "GetByIdAndStatusesTx"

	var (
		addQuery = ""
	)

	if len(statuses) > 0 {
		addQuery = fmt.Sprintf(" AND el.status IN ('%s')", strings.Join(statuses, "','"))
	}

	query := `SELECT 
				el.id, el.employee_id, el.status,
				el.value, el.allowance_id, el.updated_at,
				al.allowance_type, el.type, el.cancellation_reason,
				el.type, el.date, e.first_name, 
				e.last_name 
			FROM ` + input.TableName + ` AS el 
			LEFT JOIN ` + AllowanceDAO.TableName + ` AS al
				ON el.allowance_id = al.id
			LEFT JOIN ` + EmployeeDAO.TableName + ` AS e 
				ON el.employee_id = e.id
			WHERE 
				el.id = $1 AND
				el.deleted = FALSE` + addQuery

	row := db.QueryRow(query, id)
	err := row.Scan(
		&result.ID, &result.EmployeeId, &result.Status,
		&result.Value, &result.AllowanceId, &result.UpdatedAt,
		&result.AllowanceType, &result.Type, &result.CancellationReason,
		&result.Type, &result.Date, &result.Firstname,
		&result.Lastname)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLeaveDAO) UpdateStatusTx(tx *sql.Tx, model repository.EmployeeLeaveModel) (errModel errorModel.ErrorModel) {
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

func (input employeeLeaveDAO) UpdateStatusAndCancellationReasonTx(tx *sql.Tx, model repository.EmployeeLeaveModel) (errModel errorModel.ErrorModel) {
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

func (input employeeLeaveDAO) InitiateGetListEmployeeLeave(db *sql.DB, searchBy []in.SearchByParam, employeeLeave repository.EmployeeLeaveModel) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT 
				COUNT(*)
			FROM (
				SELECT 
					distinct on (id) id, *
				FROM (
					SELECT 
						el.id, e.id_card, e.first_name, 
						e.last_name, d.name as department_name, al.allowance_name,
						el.date, el.value, el.created_at, 
						el.updated_at, el.type, el.created_at AS leave_time,
						el.start_date, el.end_date, el.status,
						jsonb_array_elements(date::jsonb)::varchar::timestamp as d, e.id as employee_id, el.deleted
					FROM employee_leave AS el
					LEFT JOIN employee AS e
						ON el.employee_id = e.id
					LEFT JOIN department AS d 
						ON e.department_id = d.id
					LEFT JOIN allowances AS al
						ON el.allowance_id = al.id
				) el
			) el`

	addQuery, params := input.getListEmployeeLeaveAddWhere(employeeLeave)

	return GetListDataDAO.GetCountData(db, params, query, searchBy, addQuery, DefaultFieldMustCheck{
		CreatedBy: FieldStatus{
			FieldName: "created_by",
			Value:     int64(0),
		},
		Deleted: FieldStatus{
			IsCheck:   true,
			FieldName: "el.deleted",
		},
	})
}

func (input employeeLeaveDAO) GetListEmployeeLeave(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, employeeLeave repository.EmployeeLeaveModel) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT 
				DISTINCT on (id) id, *
			FROM (
				SELECT 
					el.id, e.id_card, e.first_name, 
					e.last_name, d.name as department_name, al.allowance_name,
					el.date, el.value, el.created_at, 
					el.updated_at, el.type, el.created_at AS leave_time,
					el.start_date, el.end_date, el.status,
					jsonb_array_elements(date::jsonb)::varchar::timestamp as d, e.id as employee_id, el.deleted
				FROM employee_leave AS el
				LEFT JOIN employee AS e
					ON el.employee_id = e.id
				LEFT JOIN department AS d 
					ON e.department_id = d.id
				LEFT JOIN allowances AS al
					ON el.allowance_id = al.id
			) el`

	addQuery, params := input.getListEmployeeLeaveAddWhere(employeeLeave)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var (
				model      repository.EmployeeLeaveModel
				id         sql.NullInt64
				date       sql.NullTime
				employeeId sql.NullInt64
				deleted    sql.NullBool
			)

			err := rows.Scan(
				&id, &model.ID, &model.IDCard, &model.Firstname,
				&model.Lastname, &model.Department, &model.AllowanceName,
				&model.StrDateList, &model.Value, &model.CreatedAt,
				&model.UpdatedAt, &model.Type, &model.LeaveTime,
				&model.StartDate, &model.EndDate, &model.Status,
				&date, &employeeId, &deleted,
			)

			return model, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				FieldName: "created_by",
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck:   true,
				FieldName: "el.deleted",
			},
		})
}

func (input employeeLeaveDAO) getListEmployeeLeaveAddQuery(model repository.EmployeeLeaveModel) (result string, params []interface{}) {

	if !model.IsYearly.Bool {
		result += fmt.Sprintf(" AND el.status = $%d", len(params)+1)
		params = append(params, constanta.ApprovedRequestStatus)
	}
	if model.OnLeave.Bool {
		result += fmt.Sprintf(" AND el.status IN ($%d, $%d)", len(params)+1, len(params)+2)
		params = append(params, constanta.ApprovedRequestStatus, constanta.PendingCancellationRequestStatus)

		if model.LeaveDate.String != "" {
			result += fmt.Sprintf(" AND el.date ILIKE $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.LeaveDate.String))
		}
	}

	if model.Name.String != "" {
		result += fmt.Sprintf(" AND CONCAT(e.first_name, ' ', e.last_name) ILIKE $%d", len(params)+1)
		params = append(params, fmt.Sprintf("%%%s%%", model.Name.String))
	}

	if model.SearchBy.String != "" && model.Keyword.String != "" {
		if model.SearchBy.String == "d.name" || model.SearchBy.String == "e.id_card" {
			result += "AND " + model.SearchBy.String + " ILIKE '%" + model.Keyword.String + "%'"
		}
		if model.SearchBy.String == "employee_name" {
			result += fmt.Sprintf(" AND CONCAT(e.first_name, ' ', e.last_name) ILIKE $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.Keyword.String))
		}
	}

	if model.MemberList != nil {
		isAll := false

		for _, memberId := range model.MemberList {
			if memberId == "all" {
				isAll = true
				break
			}
		}

		if !isAll {
			result += fmt.Sprintf(" AND e.id IN (%s)", strings.Join(model.MemberList, ","))
		}
	}

	if model.StrStartDate.String != "" && model.StrEndDate.String != "" {
		result += fmt.Sprintf(" AND DATE(el.created_at) BETWEEN $%d AND $%d", len(params)+1, len(params)+2)
		params = append(params, model.StrStartDate.String, model.StrEndDate.String)
	}

	return
}

func (input employeeLeaveDAO) getListEmployeeLeaveAddWhere(model repository.EmployeeLeaveModel) (result string, params []interface{}) {
	if model.SearchBy.String != "" {
		if model.SearchBy.String == "id_card" {
			result += fmt.Sprintf(" AND el.id_card ILIKE $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.Keyword.String))
		}

		if model.SearchBy.String == "name" {
			result += fmt.Sprintf(" AND CONCAT(el.first_name, ' ', el.last_name) ILIKE $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.Keyword.String))
		}

		if model.SearchBy.String == "department" {
			result += fmt.Sprintf(" AND el.department_name ILIKE $%d", len(params)+1)
			params = append(params, fmt.Sprintf("%%%s%%", model.Keyword.String))
		}
	}

	if model.OnLeave.Bool {
		result += fmt.Sprintf(" AND el.status IN ($%d, $%d)", len(params)+1, len(params)+2)
		params = append(params, constanta.ApprovedRequestStatus, constanta.PendingCancellationRequestStatus)

		if model.LeaveDate.String != "" {
			result += fmt.Sprintf(" AND date(el.d) = $%d", len(params)+1)
			params = append(params, model.LeaveDate.String)
		}
	}

	if model.Type.String != "" {
		result += fmt.Sprintf(" AND el.type = $%d", len(params)+1)
		params = append(params, model.Type.String)
	}

	if model.Status.String != "" {
		result += fmt.Sprintf(" AND el.status = $%d", len(params)+1)
		params = append(params, model.Status.String)
	}

	if model.IDCard.String != "" {
		result += fmt.Sprintf(" AND el.id_card ILIKE $%d", len(params)+1)
		params = append(params, fmt.Sprintf("%%%s%%", model.IDCard.String))
	}

	if model.Name.String != "" {
		result += fmt.Sprintf(" AND CONCAT(el.first_name, ' ', el.last_name) ILIKE $%d", len(params)+1)
		params = append(params, fmt.Sprintf("%%%s%%", model.Name.String))
	}

	if model.Department.String != "" {
		result += fmt.Sprintf(" AND el.department_name ILIKE $%d", len(params)+1)
		params = append(params, fmt.Sprintf("%%%s%%", model.Department.String))
	}

	if model.StrStartDate.String != "" && model.StrEndDate.String != "" {
		result += fmt.Sprintf(" AND DATE(el.created_at) BETWEEN $%d AND $%d", len(params)+1, len(params)+2)
		params = append(params, model.StrStartDate.String, model.StrEndDate.String)
	}

	return
}

func (input employeeLeaveDAO) GetReportEmployeeLeaveYearly(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, employeeLeave repository.EmployeeLeaveModel) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT 
				e.id, e.id_card, e.first_name, 
				e.last_name, d.name, ev.level, 
                eg.grade, eb.current_annual_leave, eb.last_annual_leave

			FROM employee e
			LEFT JOIN department AS d 
				ON e.department_id = d.id
            LEFT JOIN employee_benefits AS eb
				ON eb.employee_id = e.id
            LEFT JOIN employee_level AS ev
				ON eb.employee_level_id = ev.id
            LEFT JOIN employee_grade AS eg
				ON eb.employee_grade_id = eg.id`

	addQuery, params := input.getListEmployeeLeaveAddQuery(employeeLeave)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.EmployeeLeaveModel

			err := rows.Scan(
				&model.ID, &model.IDCard, &model.Firstname,
				&model.Lastname, &model.Department,
				&model.Level, &model.Grade,
				&model.CurrentAnnualLeave, &model.LastAnnualLeave,
			)

			return model, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				FieldName: "created_by",
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck:   true,
				FieldName: "e.deleted",
			},
		})
}

func (input employeeLeaveDAO) GetCountTodayLeave(db *sql.DB, searchByParam []in.SearchByParam, timeNow time.Time) (result int, err errorModel.ErrorModel) {
	var (
		tableName    string
		additional   string
		timeStr      = timeNow.Format(constanta.DefaultInstallationTimeFormat)
		timeStrPlus1 = timeNow.Add(time.Hour * 24).Format(constanta.DefaultTimeFormat)
		status       = constanta.ApprovedRequestStatus
		isPermit     bool
	)
	input.convertSearchByAndOrderByTodayLeave(&searchByParam, &in.GetListDataDTO{}, &isPermit)
	timeSet := fmt.Sprintf(`value = '%s'::TIMESTAMP`, timeStr)
	if isPermit {
		timeSet = fmt.Sprintf(` value >= '%s'::TIMESTAMP AND value <= '%s'::TIMESTAMP `, timeStr, timeStrPlus1)
	}

	additional += fmt.Sprintf(` 
		AND el.status = '%s' 
		AND el.id IN 
			(SELECT id
				FROM (
				  SELECT 
					id, 
					JSONB_ARRAY_ELEMENTS("date"::JSONB)::VARCHAR::TIMESTAMP AS value
				  FROM %s
				) AS subquery
				WHERE %s) `,
		status, input.TableName, timeSet)

	tableName = fmt.Sprintf(` 
		%s el 
		LEFT JOIN %s e ON el.employee_id = e.id
		LEFT JOIN %s d ON e.department_id = d.id
		LEFT JOIN %s al ON el.allowance_id = al.id `,
		input.TableName, EmployeeDAO.TableName, DepartmentDAO.TableName,
		AllowanceDAO.TableName)
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, additional, DefaultFieldMustCheck{
		Deleted:   FieldStatus{IsCheck: true, FieldName: "el.deleted"},
		ID:        FieldStatus{FieldName: "el.id"},
		CreatedBy: FieldStatus{Value: int64(0)},
	})
}

func (input employeeLeaveDAO) GetListTodayLeave(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, timeNow time.Time) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query        string
		getListData  getListJoinDataDAO
		colAddWhere  []string
		timeStr      = timeNow.Format(constanta.DefaultInstallationTimeFormat)
		timeStrPlus1 = timeNow.Add(time.Hour * 24).Format(constanta.DefaultTimeFormat)
		isPermit     bool
	)

	input.convertSearchByAndOrderByTodayLeave(&searchBy, &userParam, &isPermit)
	timeSet := fmt.Sprintf(`value = '%s'::TIMESTAMP`, timeStr)
	if isPermit {
		timeSet = fmt.Sprintf(` value >= '%s'::TIMESTAMP AND value <= '%s'::TIMESTAMP `, timeStr, timeStrPlus1)
	}

	query = fmt.Sprintf(`
		SELECT 
			e.id_card, 
			CONCAT(e.first_name, ' ', e.last_name) AS name, 
			d.name AS department,
			el.date, 
			CASE 
				WHEN el.type ILIKE 'leave' THEN 'LEAVE'
				WHEN el.type ILIKE 'sick-leave' THEN 'SICK'
				WHEN el.type ILIKE 'permit' THEN 'PERMIT'
			ELSE NULL
			END AS type_name
		FROM %s AS el `,
		input.TableName)

	colAddWhere = append(colAddWhere,
		fmt.Sprintf(` el.status = 'Approved'`), //--- Add 1
		fmt.Sprintf(` el.id IN 
		(SELECT id
			FROM (
			  SELECT 
				id, 
				JSONB_ARRAY_ELEMENTS("date"::jsonb)::VARCHAR::TIMESTAMP AS value
			  FROM %s
			) AS subquery
			WHERE %s)`,
			input.TableName, timeSet), //--- Add 2
	)
	getListData = getListJoinDataDAO{Table: "el", Query: query, AdditionalWhere: colAddWhere}
	getListData.LeftJoinAliasWithoutDeleted(EmployeeDAO.TableName, "e", "el.employee_id", "e.id")
	getListData.LeftJoinAliasWithoutDeleted(DepartmentDAO.TableName, "d", "e.department_id", "d.id")
	getListData.LeftJoinAliasWithoutDeleted(AllowanceDAO.TableName, "al", "el.allowance_id", "al.id")

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.EmployeeLeaveModel
		dbError := rows.Scan(
			&resultTemp.IDCard, &resultTemp.Name, &resultTemp.Department,
			&resultTemp.Date, &resultTemp.Type,
		)
		return resultTemp, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input employeeLeaveDAO) GetSummaryLeaveToday(db *sql.DB, timeNow time.Time) (result []interface{}, err errorModel.ErrorModel) {
	var (
		funcName = "GetSummaryLeaveToday"
		query    string
		timeStr  = timeNow.Format(constanta.DefaultInstallationTimeFormat)
		timeStr2 = timeNow.Format(constanta.DefaultTimeFormat)
	)

	query = fmt.Sprintf(`
		SELECT 
			COUNT(e.id_card) AS tot,
			CASE 
				WHEN el.type ILIKE 'leave' THEN 'LEAVE'
				WHEN el.type ILIKE 'sick-leave' THEN 'SICK'
				WHEN el.type ILIKE 'permit' THEN 'PERMIT'
			ELSE NULL
			END AS type_name
		FROM %s AS el 
		LEFT JOIN %s AS e ON el.employee_id = e.id
		LEFT JOIN %s AS d ON e.department_id = d.id
		LEFT JOIN %s AS al ON el.allowance_id = al.id 
		WHERE 
		    el.deleted = FALSE AND 
		    el.status = 'Approved' AND 
		    el.id IN 
				(SELECT id
					FROM (
					  SELECT 
						id, 
						JSONB_ARRAY_ELEMENTS("date"::jsonb)::VARCHAR::TIMESTAMP AS value
					  FROM %s
					) AS subquery
					WHERE value >= '%s'::TIMESTAMP AND value <= '%s'::TIMESTAMP) 
			GROUP BY type_name `,
		input.TableName, EmployeeDAO.TableName, DepartmentDAO.TableName,
		AllowanceDAO.TableName, input.TableName, timeStr,
		timeStr2)

	rows, dbErr := db.Query(query)
	if dbErr != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErr)
		return
	}
	resultTemp, err := RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			errs error
			temp repository.EmployeeLeaveModel
		)
		errs = rws.Scan(&temp.CountType, &temp.Type)
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.TableName, funcName, errs)
			return
		}
		resultTemp = temp
		return
	})
	for _, itemResultTemp := range resultTemp {
		t := itemResultTemp.(repository.EmployeeLeaveModel)
		result = append(result, t)
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLeaveDAO) convertSearchByAndOrderByTodayLeave(searchBy *[]in.SearchByParam, userParam *in.GetListDataDTO, isPermit *bool) {
	temp := *searchBy
	for i := 0; i < len(temp); i++ {
		switch temp[i].SearchKey {
		case "id_card":
			temp[i].SearchKey = "e." + temp[i].SearchKey
		case "name":
			temp[i].SearchKey = "CONCAT(e.first_name, ' ', e.last_name)"
		case "department":
			temp[i].SearchKey = "d.name"
		case "type":
			temp[i].SearchKey = "el.type"
			if temp[i].SearchValue == "sick" {
				temp[i].SearchValue = "sick-leave"
			}
			if temp[i].SearchValue == "permit" {
				*isPermit = true
			}
		default:
		}
	}
	switch userParam.OrderBy {
	case "id_card", "id_card ASC", "id_card DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "e.id_card " + strSplit[1]
		} else {
			userParam.OrderBy = "e.id_card"
		}
	case "name", "name ASC", "name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "CONCAT(e.first_name, ' ', e.last_name) " + strSplit[1]
		} else {
			userParam.OrderBy = "CONCAT(e.first_name, ' ', e.last_name)"
		}
	case "department", "department ASC", "department DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "d.name " + strSplit[1]
		} else {
			userParam.OrderBy = "d.name"
		}
	}
}

func (input employeeLeaveDAO) GetEmployeeLeaveYearly(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, employeeLeave repository.EmployeeLeaveModel) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT 
				e.id, e.id_card, e.first_name, e.last_name, 
                d.name, ev.level, eg.grade, eb.current_annual_leave,
                eb.last_annual_leave
			FROM employee e
			LEFT JOIN department AS d 
				ON e.department_id = d.id
            LEFT JOIN employee_benefits AS eb
				ON eb.employee_id = e.id
            LEFT JOIN employee_level AS ev
				ON eb.employee_level_id = ev.id
            LEFT JOIN employee_grade AS eg
				ON eb.employee_grade_id = eg.id`

	addQuery, params := input.getListEmployeeLeaveAddQuery(employeeLeave)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.EmployeeLeaveModel

			err := rows.Scan(
				&model.ID, &model.IDCard, &model.Firstname,
				&model.Lastname, &model.Department,
				&model.Level, &model.Grade, &model.CurrentAnnualLeave,
				&model.LastAnnualLeave,
			)

			if employeeLeave.Year.String != "" {
				leave, _ := EmployeeHistoryLeaveDAO.GetHistoryLevelByEmployeeId(db, model.ID.Int64, employeeLeave.Year.String)
				model.CurrentAnnualLeave.Int64 = leave.CurrentAnnualLeave.Int64
				model.LastAnnualLeave.Int64 = leave.LastAnnualLeave.Int64
				model.Year.String = leave.Year.String
			}

			return model, err
		}, addQuery, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				FieldName: "created_by",
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck:   true,
				FieldName: "e.deleted",
			},
		})
}

func (input employeeLeaveDAO) InitiateEmployeeLeaveYearly(db *sql.DB, searchBy []in.SearchByParam, employeeLeave repository.EmployeeLeaveModel) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT 
				COUNT(*)
			FROM employee e
			LEFT JOIN department AS d 
				ON e.department_id = d.id
            LEFT JOIN employee_benefits AS eb
				ON eb.employee_id = e.id
            LEFT JOIN employee_level AS ev
				ON eb.employee_level_id = ev.id
            LEFT JOIN employee_grade AS eg
				ON eb.employee_grade_id = eg.id
            LEFT JOIN employee_history_leave ehl
                ON e.id = ehl.employee_id`

	addQuery, params := input.getListEmployeeLeaveAddQuery(employeeLeave)

	return GetListDataDAO.GetCountData(db, params, query, searchBy, addQuery, DefaultFieldMustCheck{
		CreatedBy: FieldStatus{
			FieldName: "created_by",
			Value:     int64(0),
		},
		Deleted: FieldStatus{
			IsCheck:   true,
			FieldName: "e.deleted",
		},
	})
}

func (input employeeLeaveDAO) GetDetailEmployeeLeave(db *sql.DB, id int64) (leave repository.EmployeeLeaveModel, err errorModel.ErrorModel) {
	funcName := "GetDetailEmployeeLeave"
	query := `SELECT 
                    el.id, e.id_card, e.first_name, 
					e.last_name, d.name as department_name, al.allowance_name,
					el.value, el.created_at, 
					el.updated_at, el.type, el.created_at AS leave_time,
					el.start_date, el.end_date, el.status, 
					el.cancellation_reason, el.description,
					eb.current_annual_leave, eb.last_annual_leave, el.date, fu.host, fu.path
				FROM employee_leave AS el
				LEFT JOIN employee AS e
					ON el.employee_id = e.id
				LEFT JOIN department AS d 
					ON e.department_id = d.id
				LEFT JOIN employee_benefits AS eb
					ON eb.employee_id = el.employee_id
				LEFT JOIN allowances AS al
					ON el.allowance_id = al.id
				LEFT JOIN file_upload AS fu 
				    ON fu.id = el.file_upload_id
				WHERE el.deleted=FALSE AND el.id = $1`

	param := []interface{}{id}

	results := db.QueryRow(query, param...)

	dbError := results.Scan(&leave.ID, &leave.IDCard, &leave.Firstname,
		&leave.Lastname, &leave.Department, &leave.AllowanceName, &leave.Value, &leave.CreatedAt,
		&leave.UpdatedAt, &leave.Type, &leave.LeaveTime, &leave.StartDate, &leave.EndDate,
		&leave.Status, &leave.CancellationReason, &leave.Description,
		&leave.CurrentAnnualLeave, &leave.LastAnnualLeave, &leave.StrDateList,
		&leave.Host, &leave.Path)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLeaveDAO) GetReportAnnualLeave(db *sql.DB, date time.Time) (result []repository.EmployeeLeaveReportModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetReportAnnualLeave"
		year       = date.Year()
		month      = int(date.Month())
		query      string
		resultTemp []interface{}
	)

	query = fmt.Sprintf(`
			SELECT 
			ROW_NUMBER() OVER (ORDER BY CONCAT(e.first_name, ' ', e.last_name)) AS row_number,
			e.id_card, 
			CONCAT(e.first_name, ' ', e.last_name) AS name, 
			ep."name" AS position,
			d.name AS department,
			e.date_join::DATE::VARCHAR,
			(e.date_join::DATE + INTERVAL '3 months')::DATE::VARCHAR AS probation_done, 
			JSONB_AGG(
				CASE WHEN el."date"::jsonb IS NOT NULL THEN
					jsonb_build_object(
						'date', el.date::jsonb,
						'type', CASE 
									WHEN el.type ILIKE 'leave' THEN 'LEAVE'
									WHEN el.type ILIKE 'sick-leave' THEN 'SICK'
									WHEN el.type ILIKE 'permit' THEN 'PERMIT'
									ELSE NULL
								END,
						'description', el.description
					) ELSE NULL END
			) AS detail_leave,
			ehl.current_annual_leave AS tot_leave,
			CASE 
				WHEN (ehl.last_annual_leave = 0 OR ehl.last_annual_leave IS NULL) THEN ehl.last_annual_leave
				ELSE ehl.last_annual_leave 
			END AS current_leave
		FROM %s e 
		LEFT JOIN %s el ON el.employee_id = e.id 
			AND el.deleted = FALSE 
			AND el.status = 'Approved' 
			AND el."type" = 'leave' 
			AND el.id IN 
				(SELECT id
					FROM (
					  SELECT 
						id,
						JSONB_ARRAY_ELEMENTS("date"::jsonb)::VARCHAR::TIMESTAMP AS value
				  FROM %s
				) AS subquery
				WHERE EXTRACT(MONTH FROM value) = $1 AND EXTRACT(YEAR FROM value) = $2)
		LEFT JOIN %s d ON d.id = e.department_id
		LEFT JOIN %s al ON al.id = el.allowance_id
		LEFT JOIN %s ep ON ep.id = e.employee_position_id
		LEFT JOIN %s ehl ON ehl.employee_id = e.id AND ehl."year" = $2::VARCHAR AND ehl.deleted = FALSE
		WHERE e.deleted = false  
		GROUP BY 
			e.id_card, e.first_name, e.last_name, 
			ep."name", d."name", e.date_join, 
			ehl.current_annual_leave, ehl.last_annual_leave
		ORDER BY CONCAT(e.first_name, ' ', e.last_name) ASC `,
		EmployeeDAO.TableName, input.TableName, input.TableName,
		DepartmentDAO.TableName, AllowanceDAO.TableName, EmployeePositionDAO.TableName,
		EmployeeHistoryLeaveDAO.TableName)

	params := []interface{}{month, year}
	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.EmployeeLeaveReportModel))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLeaveDAO) resultRowsInput(rows *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
	var (
		funcName = "resultRowsInput"
		errorS   error
		temp     repository.EmployeeLeaveReportModel
	)

	errorS = rows.Scan(&temp.RowNumber, &temp.IDCard, &temp.Name,
		&temp.Position, &temp.Department, &temp.DateJoin,
		&temp.DateProbation, &temp.DetailLeave, &temp.TotalLeave,
		&temp.CurrentLeave)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp = temp
	return
}
