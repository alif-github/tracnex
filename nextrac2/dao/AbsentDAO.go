package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type absentDAO struct {
	AbstractDAO
}

var AbsentDAO = absentDAO{}.New()

func (input absentDAO) New() (output absentDAO) {
	output.FileName = "AbsentDAO.go"
	output.TableName = "absent"
	return
}

func (input absentDAO) InsertMultipleAbsent(tx *sql.Tx, userParam []repository.AbsentModel) (id []int64, err errorModel.ErrorModel) {
	var (
		funcName   = "InsertMultipleAbsent"
		param      = 26
		startIndex = 1
		query      string
		params     []interface{}
		result     []interface{}
	)

	query = fmt.Sprintf(`
		INSERT INTO %s 
		    (
		     employee_id, absent_id, 
		     normal_days, actual_days, absent, 
		     overdue, leave_early, overtime, 
		     numbers_of_leave, leaving_duties, numbers_in, 
		     numbers_out, scan, sick_leave, 
		     paid_leave, permission_leave, work_hours, 
		     percent_absent, period_start, period_end, 
		     created_by, created_client, created_at, 
		     updated_by, updated_client, updated_at
		    ) VALUES `,
		input.TableName)

	query += CreateDollarParamInMultipleRowsDAO(len(userParam), param, startIndex, "id")
	for _, i := range userParam {
		params = append(params,
			i.EmployeeID.Int64, i.AbsentID.Int64,
			i.NormalDays.Int64, i.ActualDays.Int64, i.Absent.Int64,
			i.Overdue.Int64, i.LeaveEarly.Int64, i.Overtime.Int64,
			i.NumberOfLeave.Int64, i.LeavingDuties.Int64, i.NumbersIn.Int64,
			i.NumbersOut.Int64, i.Scan.Int64, i.SickLeave.Int64,
			i.PaidLeave.Int64, i.PermissionLeave.Int64, i.WorkHours.Int64,
			i.PercentAbsent.Float64, i.PeriodStart.Time, i.PeriodEnd.Time,
			i.CreatedBy.Int64, i.CreatedClient.String, i.CreatedAt.Time,
			i.UpdatedBy.Int64, i.UpdatedClient.String, i.UpdatedAt.Time)
	}

	rows, errs := tx.Query(query, params...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	result, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemID := range result {
		id = append(id, itemID.(int64))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentDAO) UpdateMultipleAbsent(tx *sql.Tx, model []repository.AbsentModel) (err errorModel.ErrorModel) {
	var (
		funcName   = "UpdateMultipleAbsent"
		query      string
		data       string
		param      = 19
		startIndex = 1
		params     []interface{}
	)

	staticAmountParam := param
	for idx := 1; idx <= len(model); idx++ {
		data += "("
		for j := startIndex; j <= param; j++ {
			if param-j != 0 {
				data += "$" + strconv.Itoa(j) + "::bigint"
				data += ", "
			} else {
				data += "$" + strconv.Itoa(j) + "::numeric(6,3)"
				data += ")"
			}
		}

		if len(model)-idx != 0 {
			data += ","
		}

		startIndex += staticAmountParam
		param += staticAmountParam
	}

	for _, i := range model {
		params = append(params,
			i.ID.Int64, i.EmployeeID.Int64, i.AbsentID.Int64,
			i.NormalDays.Int64, i.ActualDays.Int64, i.Absent.Int64,
			i.Overdue.Int64, i.LeaveEarly.Int64, i.Overtime.Int64,
			i.NumberOfLeave.Int64, i.LeavingDuties.Int64, i.NumbersIn.Int64,
			i.NumbersOut.Int64, i.Scan.Int64, i.SickLeave.Int64,
			i.PaidLeave.Int64, i.PermissionLeave.Int64, i.WorkHours.Int64,
			i.PercentAbsent.Float64)
	}

	if data != "" {
		startIndex += param
	}

	query = fmt.Sprintf(`WITH updated_values(
			id, employee_id, absent_id, 
			normal_days, actual_days, "absent", 
			overdue, leave_early, overtime, 
			numbers_of_leave, leaving_duties, numbers_in, 
			numbers_out, scan, sick_leave, 
			paid_leave, permission_leave, work_hours, 
			percent_absent) AS
			(VALUES %s)
		UPDATE %s SET 
			employee_id = updated_values.employee_id, 
			absent_id = updated_values.absent_id, 
			normal_days = updated_values.normal_days,
			actual_days = updated_values.actual_days, 
			"absent" = updated_values."absent", 
			overdue = updated_values.overdue,
			leave_early = updated_values.leave_early,
			overtime = updated_values.overtime,
			numbers_of_leave = updated_values.numbers_of_leave,
			leaving_duties = updated_values.leaving_duties,
			numbers_in = updated_values.numbers_in,
			numbers_out = updated_values.numbers_out, 
			scan = updated_values.scan,
			sick_leave = updated_values.sick_leave,
			paid_leave = updated_values.paid_leave,
			permission_leave = updated_values.permission_leave,
			work_hours = updated_values.work_hours,
			percent_absent = updated_values.percent_absent,
			updated_by = %d, 
			updated_client = '%s',
			updated_at = '%s'::TIMESTAMP
		FROM updated_values 
		WHERE absent.id = updated_values.id `,
		data, input.TableName,
		model[0].UpdatedBy.Int64,
		model[0].UpdatedClient.String,
		model[0].UpdatedAt.Time.Format(constanta.DefaultTimeFormat))

	stmt, dbError := tx.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentDAO) GetAbsentID(db *sql.DB, userParam repository.AbsentModel) (result repository.AbsentModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetAbsentID"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT id FROM %s 
		WHERE 
		    absent_id = $1 AND 
		    period_start = $2 AND 
		    period_end = $3 AND 
		    deleted = FALSE `,
		input.TableName)

	params := []interface{}{userParam.AbsentID.Int64, userParam.PeriodStart.Time, userParam.PeriodEnd.Time}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input absentDAO) GetCountAbsent(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	var funcName = "GetCountAbsent"
	additionalAbsent := input.setSearchByAndOrderByAbsent(&searchByParam, &in.GetListDataDTO{})
	if len(additionalAbsent) != 2 {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errors.New("must period start and end"))
		return
	}

	absent := fmt.Sprintf(`%s AND %s`, additionalAbsent[0], additionalAbsent[1])
	if createdBy > 0 {
		absent += fmt.Sprintf(` AND a.created_by %d`, createdBy)
	}

	tableName := fmt.Sprintf(`
		%s e 
		LEFT JOIN %s a ON a.employee_id = e.id AND a.deleted = FALSE AND %s `,
		EmployeeDAO.TableName, input.TableName, absent)

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, "", DefaultFieldMustCheck{
		Deleted: FieldStatus{
			IsCheck:   true,
			FieldName: "e.deleted",
		},
		ID:        FieldStatus{FieldName: "e.id"},
		CreatedBy: FieldStatus{Value: int64(0)},
	})
}

func (input absentDAO) GetLastPeriodAbsent(db *sql.DB) (result repository.AbsentModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetLastPeriodAbsent"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT period_start, period_end 
		FROM %s 
		WHERE deleted = false 
		GROUP BY period_start, period_end 
		ORDER BY period_start DESC, period_end DESC 
		LIMIT 1 `,
		input.TableName)

	results := db.QueryRow(query)
	dbError := results.Scan(&result.PeriodStart, &result.PeriodEnd)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentDAO) GetPeriodAbsent(db *sql.DB, isLimit bool, limit int) (result []string, err errorModel.ErrorModel) {
	var (
		funcName   = "GetPeriodAbsent"
		query      string
		tempResult []interface{}
	)

	query = fmt.Sprintf(`
		SELECT period_start, period_end 
		FROM %s 
		WHERE deleted = false 
		GROUP BY period_start, period_end 
		ORDER BY period_start DESC, period_end DESC `,
		input.TableName)

	if isLimit {
		query += fmt.Sprintf(` LIMIT %d `, limit)
	}

	rows, dbError := db.Query(query)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	tempResult, err = RowsCatchResult(rows, func(rws *sql.Rows) (tempOutput interface{}, err errorModel.ErrorModel) {
		var (
			errorS   error
			tempData repository.AbsentModel
		)

		errorS = rows.Scan(&tempData.PeriodStart, &tempData.PeriodEnd)
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}

		tempOutput = tempData
		return
	})

	if err.Error != nil {
		return
	}

	for _, itemResult := range tempResult {
		repo := itemResult.(repository.AbsentModel)
		s := repo.PeriodStart.Time.Format(constanta.DefaultTimeSprintFormat)
		e := repo.PeriodEnd.Time.Format(constanta.DefaultTimeSprintFormat)
		result = append(result, fmt.Sprintf(`%s-%s`, s, e))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentDAO) GetListPeriodAbsent(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam) (result []interface{}, err errorModel.ErrorModel) {
	var (
		funcName    = "GetListPeriodAbsent"
		query       string
		colAddWhere []string
		getListData getListJoinDataDAO
		groupBy     string
	)

	query = fmt.Sprintf(`
		SELECT a.period_start, a.period_end FROM %s a `,
		input.TableName)

	//--- Search By Param
	if len(searchBy) > 0 {
		for i := 0; i < len(searchBy); i++ {
			//--- Search By
			if searchBy[i].SearchKey == "period" {
				p := strings.Split(searchBy[i].SearchValue, "-")
				if len(p) != 2 {
					err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "Must YYYYMMDD-YYYYMMDD", "Period", "")
					return
				}

				//--- Additional Where
				addWhereStart := fmt.Sprintf(` a.period_start = '%s'::TIMESTAMP `, p[0])
				addWhereEnd := fmt.Sprintf(` a.period_end = '%s'::TIMESTAMP `, p[1])
				colAddWhere = append(colAddWhere, addWhereStart, addWhereEnd)
			}

			searchBy = []in.SearchByParam{}
		}
	}

	//--- Order By
	userParam.OrderBy = fmt.Sprintf(` a.period_start DESC, a.period_end DESC `)

	//--- Group By
	groupBy = " GROUP BY a.period_start, a.period_end"

	getListData = getListJoinDataDAO{Table: "a", Query: query, AdditionalWhere: colAddWhere, GroupBy: groupBy}
	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var (
			resultTemp repository.AbsentModel
			resultStr  string
		)

		dbError := rows.Scan(&resultTemp.PeriodStart, &resultTemp.PeriodEnd)

		s := resultTemp.PeriodStart.Time.Format(constanta.DefaultTimeSprintFormat)
		e := resultTemp.PeriodEnd.Time.Format(constanta.DefaultTimeSprintFormat)
		resultStr = fmt.Sprintf(`%s-%s`, s, e)

		return resultStr, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input absentDAO) GetListAbsent(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		funcName    = "GetListAbsent"
		query       string
		colAddWhere []string
		getListData getListJoinDataDAO
	)

	additionalAbsent := input.setSearchByAndOrderByAbsent(&searchBy, &userParam)
	if len(additionalAbsent) != 2 {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errors.New("must period start and end"))
		return
	}

	absent := fmt.Sprintf(`%s AND %s`, additionalAbsent[0], additionalAbsent[1])
	if createdBy > 0 {
		absent += fmt.Sprintf(` AND a.created_by %d`, createdBy)
	}

	query = fmt.Sprintf(`
		SELECT 
			e.id_card, CONCAT(e.first_name, ' ', e.last_name) AS "name", a.normal_days,
			a.actual_days, a."absent", a.overdue, 
			a.leave_early, a.overtime, a.numbers_of_leave, 
			a.leaving_duties, a.numbers_in, a.numbers_out, 
			a.scan, a.sick_leave, a.paid_leave, 
			a.permission_leave, a.work_hours, a.percent_absent
		FROM %s e `,
		EmployeeDAO.TableName)

	colAddWhere = input.setScopeData(scopeLimit, scopeDB, true)

	getListData = getListJoinDataDAO{Table: "e", Query: query, AdditionalWhere: colAddWhere}
	getListData.LeftJoinAliasWithoutDeleted(input.TableName, "a", "a.employee_id", fmt.Sprintf(`e.id AND a.deleted = FALSE AND %s `, absent))

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.AbsentModel
		dbError := rows.Scan(
			&resultTemp.IDCard, &resultTemp.EmployeeName, &resultTemp.NormalDays,
			&resultTemp.ActualDays, &resultTemp.Absent, &resultTemp.Overdue,
			&resultTemp.LeaveEarly, &resultTemp.Overtime, &resultTemp.NumberOfLeave,
			&resultTemp.LeavingDuties, &resultTemp.NumbersIn, &resultTemp.NumbersOut,
			&resultTemp.Scan, &resultTemp.SickLeave, &resultTemp.PaidLeave,
			&resultTemp.PermissionLeave, &resultTemp.WorkHours, &resultTemp.PercentAbsent,
		)

		return resultTemp, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input absentDAO) resultRowsInput(rows *sql.Rows) (idTemp interface{}, err errorModel.ErrorModel) {
	var (
		funcName = "resultRowsInput"
		errorS   error
		id       int64
	)

	errorS = rows.Scan(&id)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	idTemp = id
	return
}

func (input absentDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (colAdditionalWhere []string) {
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

func (input absentDAO) setSearchByAndOrderByAbsent(searchBy *[]in.SearchByParam, userParam *in.GetListDataDTO) (absentQuery []string) {
	temp := *searchBy
	for i := 0; i < len(temp); i++ {
		switch temp[i].SearchKey {
		case "name":
			temp[i].SearchKey = "CONCAT(e.first_name, ' ', e.last_name)"
		case "period_start", "period_end":
			absentQuery = append(absentQuery, fmt.Sprintf(`a.%s = '%s'::TIMESTAMP`, temp[i].SearchKey, temp[i].SearchValue))
		default:
			temp[i].SearchKey = "e." + temp[i].SearchKey
		}
	}

	for j := 0; j < len(*searchBy); j++ {
		switch (*searchBy)[j].SearchKey {
		case "period_start", "period_end":
			*searchBy = append((*searchBy)[:j], (*searchBy)[j+1:]...)
			j = -1
		}
	}

	switch userParam.OrderBy {
	case "name", "name ASC", "name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "CONCAT(e.first_name, ' ', e.last_name)" + strSplit[1]
		} else {
			userParam.OrderBy = "CONCAT(e.first_name, ' ', e.last_name)"
		}
	default:
		userParam.OrderBy = "e." + userParam.OrderBy
	}

	return
}
