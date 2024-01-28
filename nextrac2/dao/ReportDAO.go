package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type ReportDAOInterface interface {
	GetListReport(*sql.DB, in.GetListDataDTO, []in.SearchByParam, int64, int64, bool, map[string]interface{}, map[string]applicationModel.MappingScopeDB) ([]interface{}, errorModel.ErrorModel)
	GetListReportHistory(*sql.DB, in.GetListDataDTO, []in.SearchByParam, []repository.HistoryTimeReportRedmineModel, int64, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) ([]interface{}, errorModel.ErrorModel)
	GetCountReport(*sql.DB, []in.SearchByParam, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) (int, errorModel.ErrorModel)
}

type reportDAO struct {
	AbstractDAO
}

var ReportDAO = reportDAO{}.New()

func (input reportDAO) New() (output reportDAO) {
	output.FileName = "ReportDAO.go"
	return
}

func (input reportDAO) GetCountReport(db *sql.DB, searchByParam []in.SearchByParam, _ int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		query              string
		additionalWhereCol []string
		employeeIDBundle   []int64
	)

	subQuery := fmt.Sprintf(`(SELECT e.id_card FROM %s e `, EmployeeDAO.TableName)
	query = fmt.Sprintf(`SELECT COUNT(sub.id_card) FROM %s `, subQuery)
	groupBy := fmt.Sprintf(` GROUP BY e.id_card, CONCAT(e.first_name, ' ', e.last_name)) sub `)

	additionalWhereCol = input.setScopeData(scopeLimit, scopeDB, false) //-- Scope check
	employeeIDBundle, _, _, err = input.convertUserParamAndSearchBy(nil, &searchByParam)
	if err.Error != nil {
		return
	}

	if len(employeeIDBundle) > 0 {
		var additionalWhere1 string
		input.convertNumberToAdditionalWhere(" ev.redmine_id ", employeeIDBundle, &additionalWhere1)
		additionalWhereCol = append(additionalWhereCol, additionalWhere1)
	}

	additionalWhereCol = append(additionalWhereCol, " ev.redmine_id <> 0 ")
	getListData := getListJoinDataDAO{Query: query, Table: "e", AdditionalWhere: additionalWhereCol, GroupBy: groupBy}
	getListData.InnerJoinAlias(DepartmentDAO.TableName, "d", "d.id", "e.department_id")
	getListData.LeftJoinAlias(EmployeeVariableDAO.TableName, "ev", "e.id", "ev.employee_id")
	return getListData.GetCountJoinData(db, searchByParam, 0)
}

func (input reportDAO) GetListReport(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, departmentID int64, createdBy int64, isBacklogManday bool, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query            string
		additionalWhere  string
		employeeIDBundle []int64
		sprintBundle     []string
		redmineBundle    []int64
		params           []interface{}
		departmentInfo   string
		userParamTemp    = userParam
	)

	userParamTemp.Page = -99
	userParamTemp.Limit = -99

	switch departmentID {
	case constanta.DeveloperDepartmentID, constanta.UIUXDepartmentID:
		departmentInfo += fmt.Sprintf(` AND (b.tracker = 'Task' OR b.tracker IS NULL OR b.tracker = '')  `)
	case constanta.QAQCDepartmentID:
		departmentInfo += fmt.Sprintf(` AND (b.tracker = 'Manual' OR b.tracker = 'Automation' OR b.tracker IS NULL OR b.tracker = '')  `)
	default:
	}

	if createdBy > 0 {
		departmentInfo += fmt.Sprintf(` AND (b.created_by = %d OR b.created_by IS NULL) `, createdBy)
	}

	query = fmt.Sprintf(`
		SELECT 
		e.id_card, CONCAT(e.first_name, ' ', e.last_name), 
			CASE WHEN SUM(b.mandays) IS NOT NULL 
			THEN SUM(b.mandays) 
			ELSE 0 END tot_mandays, 
			ev.mandays_rate as rate,
		d."name" as department, b.tracker, JSONB_AGG(b.redmine_number) as redmine_number, 
		d.id as department_id
		FROM %s e 
		INNER JOIN %s d ON d.id = e.department_id
		LEFT JOIN %s b ON e.id = b.employee_id AND (b.deleted = FALSE OR b.deleted IS NULL) AND (b.payment_status = FALSE OR b.payment_status IS NULL) %s 
		LEFT JOIN %s ev ON e.id = ev.employee_id AND (ev.deleted = FALSE OR ev.deleted IS NULL) `,
		EmployeeDAO.TableName, DepartmentDAO.TableName, BacklogDAO.TableName,
		departmentInfo, EmployeeVariableDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	employeeIDBundle, redmineBundle, sprintBundle, err = input.convertUserParamAndSearchBy(&userParamTemp, &searchByParam)
	if err.Error != nil {
		return
	}

	if len(employeeIDBundle) > 0 {
		input.convertNumberToAdditionalWhere(" AND ev.redmine_id ", employeeIDBundle, &additionalWhere)
	}

	if len(redmineBundle) > 0 && !isBacklogManday {
		input.convertNumberToAdditionalWhere(" AND b.redmine_number ", redmineBundle, &query) //-- If Backlog Not Check This Function
		if len(sprintBundle) > 0 {
			input.convertStringToAdditionalWhere(" AND b.sprint ", sprintBundle, &query) //-- If Redmine Exist, Sprint Must Add
		}
	}

	additionalWhere += fmt.Sprintf(`
		AND d.deleted = FALSE AND ev.redmine_id <> 0
		GROUP BY e.id_card, CONCAT(e.first_name, ' ', e.last_name), ev.mandays_rate, d.name, d.id, b.tracker `)

	result, err = GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParamTemp, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var (
				temp   repository.ReportModel
				manday float64
			)

			dbError := rows.Scan(
				&temp.NIK, &temp.Name, &manday,
				&temp.RateStr, &temp.Department, &temp.Tracker,
				&temp.RedmineNumber, &temp.DepartmentID)

			if isBacklogManday {
				temp.BacklogManday.Float64 = manday
			} else {
				temp.ActualManday.Float64 = manday
			}

			return temp, dbError
		},
		additionalWhere, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{Value: int64(0)},
			Deleted: FieldStatus{
				IsCheck:   true,
				FieldName: "e.deleted",
				Value:     false,
			},
		})

	return
}

func (input reportDAO) GetListReportHistory(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, historyTimeRedmine []repository.HistoryTimeReportRedmineModel, departmentID int64, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query            string
		additionalWhere  string
		employeeIDBundle []int64
		params           []interface{}
		departmentInfo   string
		values           string
		reportTable      = "report"
	)

	if len(historyTimeRedmine) > 0 {
		for i := 0; i < len(historyTimeRedmine); i++ {
			d := historyTimeRedmine[i]
			values += fmt.Sprintf(`
			(%d, %d, '%s', '%s')`,
				d.User.Int64, d.RedmineTicket.Int64, d.Tracker.String, d.TimeHistory.String)

			if len(historyTimeRedmine)-(i+1) != 0 {
				values += ","
			}
		}
	} else {
		values += "(-1, -1, '-', null)"
	}

	switch departmentID {
	case constanta.DeveloperDepartmentID, constanta.UIUXDepartmentID:
		departmentInfo += fmt.Sprintf(` AND (b.tracker = 'Task' OR b.tracker IS NULL OR b.tracker = '')  `)
	case constanta.QAQCDepartmentID:
		departmentInfo += fmt.Sprintf(` AND (b.tracker = 'Manual' OR b.tracker = 'Automation' OR b.tracker IS NULL OR b.tracker = '')  `)
	default:
	}

	if createdBy > 0 {
		departmentInfo += fmt.Sprintf(` AND (b.created_by = %d OR b.created_by IS NULL) `, createdBy)
	}

	query = fmt.Sprintf(`
		with %s("user", redmine_ticket, tracker, time_history) as (values %s)
		SELECT e.id_card, CONCAT(e.first_name, ' ', e.last_name) as name, 
		jsonb_agg(
			jsonb_build_object(
				'ticket', r.redmine_ticket,
				'history', r.time_history::json
			)
		) as tot_mandays,
		ev.mandays_rate as rate, d."name" as department, b.tracker, 
		JSONB_AGG(b.redmine_number) as redmine_number, 
		d.id as department_id
		FROM %s e 
			INNER JOIN %s d ON d.id = e.department_id
			INNER JOIN %s b ON e.id = b.employee_id %s
			INNER JOIN %s r ON b.redmine_number = r.redmine_ticket AND b.tracker = r.tracker  
			LEFT JOIN %s ev ON e.id = ev.employee_id AND ev.deleted = FALSE AND ev.redmine_id = r."user"
		`,
		reportTable, values, EmployeeDAO.TableName,
		DepartmentDAO.TableName, BacklogDAO.TableName, departmentInfo,
		reportTable, EmployeeVariableDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	employeeIDBundle, _, _, err = input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	if err.Error != nil {
		return
	}

	if len(employeeIDBundle) > 0 {
		input.convertNumberToAdditionalWhere(" AND ev.redmine_id ", employeeIDBundle, &additionalWhere)
	}

	additionalWhere += fmt.Sprintf(` 
		AND (b.deleted = FALSE OR b.deleted IS NULL) 
		AND (b.payment_status = FALSE OR b.payment_status IS NULL) 
		AND d.deleted = FALSE AND ev.redmine_id <> 0
		GROUP BY e.id_card, CONCAT(e.first_name, ' ', e.last_name), ev.mandays_rate, d.name, d.id, b.tracker `)

	result, err = GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ReportModel
			dbError := rows.Scan(
				&temp.NIK, &temp.Name, &temp.ActualHistoryManday,
				&temp.RateStr, &temp.Department, &temp.Tracker,
				&temp.RedmineNumber, &temp.DepartmentID)

			return temp, dbError
		},
		additionalWhere, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{Value: int64(0)},
			Deleted: FieldStatus{
				IsCheck:   true,
				FieldName: "e.deleted",
				Value:     false,
			},
		})

	return
}

func (input reportDAO) GetListReportForInfraDevOps(db *sql.DB, userParam in.GetListDataDTO, resultDBInfraDevOps []repository.RedmineInfraModel, searchByParam []in.SearchByParam, departmentID int64, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query            string
		additionalWhere  string
		employeeIDBundle []int64
		params           []interface{}
		values           string
	)

	if len(resultDBInfraDevOps) > 0 {
		for i := 0; i < len(resultDBInfraDevOps); i++ {
			d := resultDBInfraDevOps[i]
			values += fmt.Sprintf(`
			(%d, %d, %d, %.3f, '%s')`,
				d.RedmineTicket.Int64, d.SignedID.Int64, departmentID, d.Manhour.Float64, d.Tracker.String)

			if len(resultDBInfraDevOps)-(i+1) != 0 {
				values += ","
			}
		}
	} else {
		values += "(-1, -1, -1, 0, '')"
	}

	query = fmt.Sprintf(`
		WITH backlog_infra(ticket, signed_id, department_id, manhour, tracker) AS (VALUES %s)
		SELECT e.id_card, CONCAT(e.first_name, ' ', e.last_name) as name, 
			CASE WHEN SUM(bi.manhour) IS NOT NULL 
			THEN SUM(bi.manhour) 
			ELSE 0 END tot_mandays, 
			ev.mandays_rate AS rate, 
		d."name" AS department, JSONB_AGG(bi.ticket) AS redmine_number, 
		d.id AS department_id 
		FROM %s e 
		INNER JOIN %s d ON d.id = e.department_id
		INNER JOIN %s ev ON ev.employee_id = e.id
		LEFT JOIN %s bi ON ev.redmine_id = bi.signed_id `,
		values, EmployeeDAO.TableName, DepartmentDAO.TableName,
		EmployeeVariableDAO.TableName, "backlog_infra")

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	employeeIDBundle, _, _, err = input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	if err.Error != nil {
		return
	}

	if len(employeeIDBundle) > 0 {
		input.convertNumberToAdditionalWhere(" AND ev.redmine_id ", employeeIDBundle, &additionalWhere)
	}

	additionalWhere += fmt.Sprintf(`  
		AND d.id = %d
		AND d.deleted = FALSE AND ev.redmine_id <> 0
		GROUP BY e.id_card, CONCAT(e.first_name, ' ', e.last_name), ev.mandays_rate, d.name, d.id `,
		departmentID)

	result, err = GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ReportModel
			dbError := rows.Scan(
				&temp.NIK, &temp.Name, &temp.ActualManday,
				&temp.RateStr, &temp.Department, &temp.RedmineNumber,
				&temp.DepartmentID)

			return temp, dbError
		},
		additionalWhere, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{Value: int64(0)},
			Deleted: FieldStatus{
				IsCheck:   true,
				FieldName: "e.deleted",
				Value:     false,
			},
		})

	return
}

func (input reportDAO) UnionAllResultReport(db *sql.DB, dataInput []out.ReportResponse, userParam in.GetListDataDTO) (result string, err errorModel.ErrorModel) {
	var (
		funcName   = "UnionAllResultReport"
		query      string
		values     string
		orderQuery string
		params     []interface{}
	)

	//--- Values
	for i := 0; i < len(dataInput); i++ {
		d := dataInput[i]
		values += fmt.Sprintf(`
			(%d, '%s', '%s', %.4f, %.4f, '%s', %.2f, %.2f, '%s', '%s')`,
			d.NIK, d.Name, d.Department, d.BacklogManday,
			d.ActualManday, d.Tracker, d.MandayRate, d.Manday,
			d.BacklogTicket, d.ActualTicket)

		if len(dataInput)-(i+1) != 0 {
			values += ","
		}
	}

	//--- Order By
	switch userParam.OrderBy {
	case "nik", "nik ASC", "nik DESC", "e.nik", "e.nik ASC", "e.nik DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "nik " + strSplit[1]
		} else {
			userParam.OrderBy = "nik"
		}
	case "name", "name ASC", "name DESC", "e.name", "e.name ASC", "e.name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "name " + strSplit[1]
		} else {
			userParam.OrderBy = "name"
		}
	default:
	}

	if userParam.OrderBy != "" {
		orderQuery += fmt.Sprintf(` ORDER BY %s `, userParam.OrderBy)
	}

	if userParam.Limit != -99 && userParam.Page != -99 {
		orderQuery += fmt.Sprintf(` LIMIT $1 OFFSET $2 `)
		params = append(params, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
	}

	query = fmt.Sprintf(`
			WITH result_report(nik, "name", department, backlog_manday, actual_manday, tracker, manday_rate, manday, ticket_backlog, ticket_actual) AS (VALUES %s)
			SELECT json_build_object('results', json_agg(json_build_object(
				'nik', nik,
				'name', "name",
				'department', department,
				'total_manday', total_manday,
				'detail', detail
			)))
			FROM (
				SELECT
					nik,
					"name",
					department,
					sum(manday) as total_manday,
					jsonb_agg(
						jsonb_build_object(
							'backlog_manday', actual_manday,
							'backlog_ticket', CASE WHEN ticket_actual <> '' OR ticket_actual IS NOT NULL OR ticket_actual <> '[null]' THEN ticket_actual ELSE '' END,
							'actual_manday', backlog_manday,
							'actual_ticket', CASE WHEN ticket_backlog <> '' OR ticket_backlog IS NOT NULL OR ticket_backlog <> '[null]' THEN ticket_backlog ELSE '' END,
							'tracker', tracker,
							'manday_rate', manday_rate,
							'manday', manday
						)
					) AS detail
				FROM result_report 
				GROUP BY nik, "name", department %s
			) sub_query `, values, orderQuery)

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(&result)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportDAO) InsertToReportHistory(db *sql.Tx, userParam repository.ReportHistoryModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertToReportHistory"
		query    string
	)

	query = fmt.Sprintf(`
		INSERT INTO report_history ("data", success_ticket, department_id, created_by, created_at, created_client)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id `)

	params := []interface{}{
		userParam.Data.String, userParam.SuccessTicket.String, userParam.DepartmentID.Int64,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.CreatedClient.String,
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportDAO) convertNumberToAdditionalWhere(key string, idBundle []int64, additionalWhere *string) {
	*additionalWhere += " " + key + " IN ("
	for idx, value := range idBundle {
		if len(idBundle)-(idx+1) == 0 {
			*additionalWhere += strconv.Itoa(int(value)) + ")"
		} else {
			*additionalWhere += strconv.Itoa(int(value)) + ","
		}
	}
}

func (input reportDAO) convertStringToAdditionalWhere(key string, sprintBundle []string, additionalWhere *string) {
	*additionalWhere += " " + key + " IN ("
	for idx, value := range sprintBundle {
		if len(sprintBundle)-(idx+1) == 0 {
			*additionalWhere += fmt.Sprintf(`'%s'`, value) + ")"
		} else {
			*additionalWhere += fmt.Sprintf(`'%s'`, value) + ","
		}
	}
}

func (input reportDAO) arrNumGenerate(idBundleStr string) (idBundle []int64, err errorModel.ErrorModel) {
	funcName := "arrNumGenerate"
	errorS := json.Unmarshal([]byte(idBundleStr), &idBundle)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}
	return
}

func (input reportDAO) arrStrGenerate(sprintBundleStr string) (sprintBundle []string, err errorModel.ErrorModel) {
	funcName := "arrStrGenerate"
	errorS := json.Unmarshal([]byte(sprintBundleStr), &sprintBundle)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}
	return
}

func (input reportDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) (employeeIDBundle, redmineBundle []int64, sprintBundle []string, err errorModel.ErrorModel) {
	for i := 0; i < len(*searchByParam); i++ {
		switch (*searchByParam)[i].SearchKey {
		case "department":
			(*searchByParam)[i].SearchKey = "d.id"
		case "sprint":
			//--- Sprint Generate
			sprintBundle, err = input.arrStrGenerate((*searchByParam)[i].SearchValue)
			if err.Error != nil {
				return
			}
		case "id":
			//--- Employee ID Generate
			employeeIDBundle, err = input.arrNumGenerate((*searchByParam)[i].SearchValue)
			if err.Error != nil {
				return
			}
		case "redmine":
			//--- Redmine Ticket Generate
			redmineBundle, err = input.arrNumGenerate((*searchByParam)[i].SearchValue)
			if err.Error != nil {
				return
			}
		default:
		}
	}

	//--- Delete Search By Param Employee ID
	for j := 0; j < len(*searchByParam); j++ {
		if ((*searchByParam)[j].SearchKey == "id") || ((*searchByParam)[j].SearchKey == "sprint" || (*searchByParam)[j].SearchKey == "redmine") {
			*searchByParam = append((*searchByParam)[:j], (*searchByParam)[j+1:]...)
			j = -1
		}
	}

	if userParam == nil {
		return
	}

	switch userParam.OrderBy {
	case "nik", "nik ASC", "nik DESC":
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
	}

	return
}

func (input reportDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (colAdditionalWhere []string) {
	keyScope := []string{constanta.EmployeeDataScope}
	for _, itemKeyScope := range keyScope {
		var additionalWhere string
		PrepareScopeOnDAO(scopeLimit, scopeDB, &additionalWhere, 0, itemKeyScope, isView)
		if additionalWhere != "" {
			colAdditionalWhere = append(colAdditionalWhere, additionalWhere)
		}
	}

	return
}
