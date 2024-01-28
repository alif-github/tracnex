package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type redmineDAO struct {
	AbstractDAO
}

var RedmineDAO = redmineDAO{}.New()

func (input redmineDAO) New() (output redmineDAO) {
	output.FileName = "RedmineDAO.go"
	return
}

func (input redmineDAO) GetTaskDeveloperOnRedmine(db *sql.DB, userParam in.GetListDataDTO, searchByParam *[]in.SearchByParam, redmineID []int64) (result []repository.RedmineModel, err errorModel.ErrorModel) {
	var (
		funcName  = "GetTaskDeveloperOnRedmine"
		query     string
		queryTemp string
	)

	for idx, itemRedmineID := range redmineID {
		if idx == 0 {
			queryTemp += fmt.Sprintf(` and users.id IN (`)
		}

		if len(redmineID)-(idx+1) == 0 {
			queryTemp += fmt.Sprintf(`%d)`, itemRedmineID)
		} else {
			queryTemp += fmt.Sprintf(`%d,`, itemRedmineID)
		}
	}

	query = fmt.Sprintf(`select
		journals.journalized_id as redmine_ticket, projects.name as project, users.login as signed,
		issues.assigned_to_id as signed_id, issues.subject, custom_values.value as sprint, 
		issue_statuses.name as status, journal_details.value as status_id, journals.created_on, 
		issues.estimated_hours as manhour, cv.value as payment, trac.name as tracker, 
		issues.updated_on
		from journals 
			left join journal_details on journals.id = journal_details.journal_id
			left join issue_statuses on CAST (issue_statuses.id AS INT) = CAST (journal_details.value AS INT)
			left join issues on journals.journalized_id = issues.id
			left join projects on issues.project_id = projects.id
			left join users on issues.assigned_to_id = users.id
			left join custom_values on custom_values.customized_id = issues.id
			left join custom_values AS cv on cv.customized_id = issues.id
			left join custom_fields cf on cf.id = custom_values.custom_field_id 
			left join trackers trac on trac.id = issues.tracker_id
		where 
		journal_details.property = 'attr' 
		and journal_details.prop_key = 'status_id' 
		and journal_details.value IN ('24') 
		and cv.custom_field_id = '7'
		and cv.value = 'UNPAID' 
		and trac.id IN (7)
		and custom_values.value <> ''
		and cf.id = 9 
		and journals.created_on >= $1 
		and journals.created_on < $2 
		and journals.created_on = (
			SELECT MAX(created_on)
			FROM journals j_max 
				LEFT JOIN journal_details ON j_max.id = journal_details.journal_id
				WHERE j_max.journalized_id = journals.journalized_id
				AND journal_details.property = 'attr' 
				AND journal_details.prop_key = 'status_id'
		) %s `, queryTemp)

	return input.processHitToRedmine(funcName, db, userParam, searchByParam, query)
}

func (input redmineDAO) GetTaskQAOnRedmine(db *sql.DB, userParam in.GetListDataDTO, searchByParam *[]in.SearchByParam, redmineID []int64) (result []repository.RedmineModel, err errorModel.ErrorModel) {
	var (
		funcName  = "GetTaskQAOnRedmine"
		query     string
		queryTemp string
	)

	for idx, itemRedmineID := range redmineID {
		if idx == 0 {
			queryTemp += fmt.Sprintf(` and users.id IN (`)
		}

		if len(redmineID)-(idx+1) == 0 {
			queryTemp += fmt.Sprintf(`%d)`, itemRedmineID)
		} else {
			queryTemp += fmt.Sprintf(`%d,`, itemRedmineID)
		}
	}

	query = fmt.Sprintf(`select 
		journals.journalized_id as redmine_ticket, projects.name as project, users.login as signed,
		issues.assigned_to_id as signed_id, issues.subject, custom_values.value as sprint,
		issue_statuses.name as status, journal_details.value as status_id, journals.created_on,
		issues.estimated_hours as manhour, cv.value as payment, trac.name as tracker, 
		issues.updated_on
		from journals 
			left join journal_details on journals.id = journal_details.journal_id
			left join issue_statuses on CAST (issue_statuses.id AS int) = CAST (journal_details.value AS INT)
			left join issues on journals.journalized_id = issues.id
			left join projects on issues.project_id = projects.id
			left join users on issues.assigned_to_id = users.id
			left join custom_values on custom_values.customized_id = issues.id
			left join custom_values as cv on cv.customized_id = issues.id
			left join custom_fields cf on cf.id = custom_values.custom_field_id 
			left join trackers trac on trac.id = issues.tracker_id
		where 
		journal_details.property = 'attr' 
		and journal_details.prop_key = 'status_id' 
		and journal_details.value in ('24')  
		and cv.custom_field_id = '7'
		and cv.value = 'UNPAID'
		and cf.id = 9
		and custom_values.value <> ''
		and trac.id in (13,14) 
		and journals.created_on >= $1 
		and journals.created_on < $2 
		and journals.created_on = (
			SELECT MAX(created_on)
			FROM journals j_max 
				LEFT JOIN journal_details ON j_max.id = journal_details.journal_id
				WHERE j_max.journalized_id = journals.journalized_id
				AND journal_details.property = 'attr' 
				AND journal_details.prop_key = 'status_id'
		) %s `, queryTemp)

	return input.processHitToRedmine(funcName, db, userParam, searchByParam, query)
}

func (input redmineDAO) GetManhourOnRedmine(db *sql.DB, task []repository.RedmineModel) (result []repository.HistoryTimeReportRedmineModel, err errorModel.ErrorModel) {
	var (
		query    string
		values   string
		funcName = "GetManhourOnRedmine"
		start    = `%start%`
		pause    = `%pause%`
		end      = `%end%`
	)

	if len(task) > 0 {
		for i := 0; i < len(task); i++ {
			d := task[i]
			values += fmt.Sprintf(`
			(%d, '%s', '%s', %d, '%s', '%s', %d, '%s'::TIMESTAMP, %f, '%s', '%s', '%s'::TIMESTAMP)`,
				d.RedmineTicket.Int64, d.Project.String, d.Signed.String,
				d.SignedID.Int64, d.Sprint.String, d.Status.String,
				d.StatusID.Int64, d.CreatedAt.Time.Format(constanta.DefaultDBSQLTimeFormat), d.Manhour.Float64,
				d.Payment.String, d.Tracker.String, d.UpdatedOn.Time.Format(constanta.DefaultDBSQLTimeFormat))

			if len(task)-(i+1) != 0 {
				values += ","
			}
		}
	} else {
		values += "(0, '', '', 0, '', '', 0, null, null, '', '', null)"
	}

	query = fmt.Sprintf(`
			with report(redmine_ticket, project, signed, signed_id, sprint, status, status_id, created_on, manhour, payment, tracker, updated_on) 
			as (values %s)
			select 
			sub_query."user",
			sub_query.redmine_ticket,
			sub_query.tracker,
			jsonb_agg(
				jsonb_build_object(
					'subject', sub_query.subject,
					'created_on', sub_query.created_on
				)
			) as time_history 
			from (
				select r.signed_id as "user", r.redmine_ticket, r.tracker,
					case 
						when LOWER(i.subject) like '%s' then 'start'
						when LOWER(i.subject) like '%s' then 'pause'
						when LOWER(i.subject) like '%s' then 'end'
					else '' end subject, 
				i.created_on 
				from issues i  
				inner join report r on i.parent_id = r.redmine_ticket and i.assigned_to_id = r.signed_id
				where 
				(
					LOWER(i.subject) like '%s' 
					OR
					LOWER(i.subject) like '%s'
					OR
					LOWER(i.subject) like '%s'
				) 
				and i.tracker_id = 2 
				order by "user" asc, created_on asc
			) as sub_query 
			group by sub_query."user", sub_query."redmine_ticket", sub_query.tracker `,
		values, start, pause, end, start, pause, end)

	return input.processHitToManhourRedmine(funcName, db, query)
}

func (input redmineDAO) GetTaskInfraOnRedmine(db *sql.DB, userParam in.GetListDataDTO, departmentID int64, searchByParam *[]in.SearchByParam, redmineID []int64) (result []repository.RedmineInfraModel, err errorModel.ErrorModel) {
	var (
		funcName     = "GetTaskInfraOnRedmine"
		query        string
		trackerQuery string
		queryTemp    string
	)

	switch departmentID {
	case constanta.InfraDepartmentID:
		trackerQuery = fmt.Sprintf(` AND trackers.id IN ('4', '5', '6') `)
	case constanta.DevOpsDepartmentID:
		trackerQuery = fmt.Sprintf(` AND trackers.id IN ('10') `)
	default:
	}

	for idx, itemRedmineID := range redmineID {
		if idx == 0 {
			queryTemp += fmt.Sprintf(` AND issues.assigned_to_id IN (`)
		}

		if len(redmineID)-(idx+1) == 0 {
			queryTemp += fmt.Sprintf(`%d)`, itemRedmineID)
		} else {
			queryTemp += fmt.Sprintf(`%d,`, itemRedmineID)
		}
	}

	query = `
		SELECT 
			journals.journalized_id as RedmineTicket,
			projects.name as Project,
			users.login as signed,
			issues.assigned_to_id as signed_id,
			issues.subject,
			issue_statuses.name as status, 
			journal_details.value,
			journals.created_on,
			journals.notes,
			issues.estimated_hours Manhour,
			trackers.name Tracker,
			issues.updated_on,
			GROUP_CONCAT(custom_values.value separator '|') as custom_field_value
		FROM journals 
		LEFT JOIN journal_details ON journals.id = journal_details.journal_id
		LEFT JOIN issue_statuses ON (issue_statuses.id) = (journal_details.value)
		LEFT JOIN issues ON journals.journalized_id = issues.id
		LEFT JOIN projects ON issues.project_id = projects.id
		LEFT JOIN users ON issues.assigned_to_id = users.id
		LEFT JOIN trackers on trackers.id = issues.tracker_id 
		LEFT JOIN custom_fields_trackers cft ON trackers.id = cft.tracker_id
		LEFT JOIN custom_fields cf ON cft.custom_field_id = cf.id
		LEFT JOIN custom_values ON custom_values.customized_id = issues.id
		LEFT JOIN custom_fields ON custom_values.custom_field_id = custom_fields.id
		WHERE 
			journal_details.property = 'attr' AND 
			journal_details.prop_key = 'status_id' AND 
			issue_statuses.name = 'Closed' AND
			users.login != '' AND
			cf.name LIKE '%issue%' AND
			journals.created_on >= ? AND
			journals.created_on < ? AND 
			journals.created_on = (
				SELECT MAX(created_on)
				FROM journals j_max 
					LEFT JOIN journal_details ON j_max.id = journal_details.journal_id
					WHERE j_max.journalized_id = journals.journalized_id
					AND journal_details.property = 'attr' 
					AND journal_details.prop_key = 'status_id'
			) AND 
			custom_fields.name LIKE '%issue%' `

	query += trackerQuery
	query += queryTemp

	return input.processHitToRedmineInfra(funcName, db, userParam, searchByParam, query)
}

func (input redmineDAO) processHitToManhourRedmine(funcName string, db *sql.DB, query string) (result []repository.HistoryTimeReportRedmineModel, err errorModel.ErrorModel) {
	var resultTemp []interface{}
	rows, errorS := db.Query(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			dbError error
			model   repository.HistoryTimeReportRedmineModel
		)

		dbError = rws.Scan(&model.User, &model.RedmineTicket, &model.Tracker, &model.TimeHistory)
		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}

		resultTemp = model
		return
	}); err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.HistoryTimeReportRedmineModel))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input redmineDAO) processHitToRedmine(funcName string, db *sql.DB, userParam in.GetListDataDTO, searchByParam *[]in.SearchByParam, query string) (result []repository.RedmineModel, err errorModel.ErrorModel) {
	var (
		employeeIDBundle []int64
		resultTemp       []interface{}
		redmine          []int64
		sprint           []string
	)

	for i := 0; i < len(*searchByParam); i++ {
		if (*searchByParam)[i].SearchKey == "id" {
			employeeIDBundle, err = ReportDAO.arrNumGenerate((*searchByParam)[i].SearchValue)
			if err.Error != nil {
				return
			}
		}
	}

	if len(employeeIDBundle) > 0 {
		ReportDAO.convertNumberToAdditionalWhere(" AND users.id ", employeeIDBundle, &query)
	}

	start := userParam.UpdatedAtStart.Format(constanta.DefaultInstallationTimeFormat)
	end := userParam.UpdatedAtEnd.Format(constanta.DefaultInstallationTimeFormat)

	param := []interface{}{start, end}
	rows, errorS := db.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			dbError error
			model   repository.RedmineModel
		)

		dbError = rws.Scan(
			&model.RedmineTicket, &model.Project, &model.Signed,
			&model.SignedID, &model.Subject, &model.Sprint,
			&model.Status, &model.StatusID, &model.CreatedAt,
			&model.Manhour, &model.Payment, &model.Tracker,
			&model.UpdatedOn)

		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}

		resultTemp = model
		return
	}); err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.RedmineModel))
		redmine = append(redmine, itemResultTemp.(repository.RedmineModel).RedmineTicket.Int64)
		sprint = append(sprint, itemResultTemp.(repository.RedmineModel).Sprint.String)
	}

	//-- Set to search by param
	if len(redmine) > 0 {
		var redmineBundleByte []byte
		redmineBundleByte, errorS = json.Marshal(redmine)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			return
		}

		*searchByParam = append(*searchByParam, in.SearchByParam{
			SearchKey:      "redmine",
			DataType:       "char",
			SearchOperator: "eq",
			SearchValue:    string(redmineBundleByte),
			SearchType:     constanta.Filter,
		})
	}

	//-- Set to search by param
	if len(sprint) > 0 {
		var sprintBundleByte []byte
		sprintBundleByte, errorS = json.Marshal(sprint)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			return
		}

		*searchByParam = append(*searchByParam, in.SearchByParam{
			SearchKey:      "sprint",
			DataType:       "char",
			SearchOperator: "eq",
			SearchValue:    string(sprintBundleByte),
			SearchType:     constanta.Filter,
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input redmineDAO) processHitToRedmineInfra(funcName string, db *sql.DB, userParam in.GetListDataDTO, searchByParam *[]in.SearchByParam, query string) (result []repository.RedmineInfraModel, err errorModel.ErrorModel) {
	var (
		employeeIDBundle []int64
		resultTemp       []interface{}
		redmine          []int64
	)

	for i := 0; i < len(*searchByParam); i++ {
		if (*searchByParam)[i].SearchKey == "id" {
			employeeIDBundle, err = ReportDAO.arrNumGenerate((*searchByParam)[i].SearchValue)
			if err.Error != nil {
				return
			}
		}
	}

	if len(employeeIDBundle) > 0 {
		ReportDAO.convertNumberToAdditionalWhere(" AND issues.assigned_to_id ", employeeIDBundle, &query)
	}

	query += fmt.Sprintf(` 
		GROUP BY 
			journals.journalized_id, journals.created_on, projects.name,
			users.login, issues.assigned_to_id, issues.subject,
			issue_statuses.name, journal_details.value, journals.notes,
			issues.estimated_hours, trackers.name, issues.updated_on
		ORDER BY journals.journalized_id ASC `)

	start := userParam.UpdatedAtStart.Format(constanta.DefaultInstallationTimeFormat)
	end := userParam.UpdatedAtEnd.Format(constanta.DefaultInstallationTimeFormat)

	param := []interface{}{start, end}
	rows, errorS := db.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			dbError    error
			model      repository.RedmineInfraModel
			createdStr string
			updatedStr string
		)

		dbError = rws.Scan(
			&model.RedmineTicket, &model.Project, &model.Signed,
			&model.SignedID, &model.Subject, &model.Status,
			&model.Value, &createdStr, &model.Notes,
			&model.Manhour, &model.Tracker, &updatedStr,
			&model.Category)

		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}

		model.CreatedAt.Time, dbError = time.Parse(constanta.DefaultDBSQLTimeFormat, createdStr)
		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}

		model.UpdatedOn.Time, dbError = time.Parse(constanta.DefaultDBSQLTimeFormat, updatedStr)
		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
			return
		}

		resultTemp = model
		return
	}); err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.RedmineInfraModel))
		redmine = append(redmine, itemResultTemp.(repository.RedmineInfraModel).RedmineTicket.Int64)
	}

	//-- Set to search by param
	if len(redmine) > 0 {
		var redmineBundleByte []byte
		redmineBundleByte, errorS = json.Marshal(redmine)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			return
		}

		*searchByParam = append(*searchByParam, in.SearchByParam{
			SearchKey:      "redmine",
			DataType:       "char",
			SearchOperator: "eq",
			SearchValue:    string(redmineBundleByte),
			SearchType:     constanta.Filter,
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input redmineDAO) GetCustomFieldsFromRedmine(db *sql.DB, id int) (value string, err errorModel.ErrorModel) {
	var (
		funcName = "GetCustomFieldsFromRedmine"
		query    string
		params   []interface{}
	)

	query = fmt.Sprintf(`SELECT possible_values FROM custom_fields WHERE id = $1`)
	params = []interface{}{id}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(&value)
	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input redmineDAO) UpdateCustomFieldsOnRedmine(db *sql.Tx, id int, value string) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateCustomFieldsOnRedmine"
		query    string
	)

	query = fmt.Sprintf(`UPDATE custom_fields SET possible_values = $1 WHERE id = $2 `)
	param := []interface{}{value, id}

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

func (input redmineDAO) GetEmployeeRedmineByNIK(db *sql.DB, nik string, customFieldID string, isSqlParam bool) (user repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeeRedmineByNIK"
		query    string
		params   []interface{}
		param1   = "$1"
		param2   = "$2"
	)

	if isSqlParam {
		param1 = "?"
		param2 = "?"
	}

	query = fmt.Sprintf(`
		SELECT 
		u.id as redmine_id, u.firstname, u.lastname,
		TRIM(cv.value) as "NIK"
		FROM custom_values cv
		LEFT JOIN custom_fields cf ON cf.id = cv.custom_field_id 
		LEFT JOIN users u ON u.id = cv.customized_id 
		WHERE 
		cf.id = %s AND 
		cv.value != '0' AND 
		cv.value != '' AND 
		TRIM(cv.value) = %s `,
		param1, param2)

	params = []interface{}{customFieldID, nik}
	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&user.RedmineId, &user.FirstName, &user.LastName,
		&user.IDCard)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}
