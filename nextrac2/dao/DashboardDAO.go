package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type dashboardDAO struct {
	AbstractDAO
}

var DashboardDAO = dashboardDAO{}.New()

func (input dashboardDAO) New() (output dashboardDAO) {
	output.FileName = "DashboardDAO.go"
	return
}

func (input dashboardDAO) GetReimbursementInMonth(db *sql.DB, timeNow time.Time) (count int64, err errorModel.ErrorModel) {
	var (
		funcName  = "GetReimbursementInMonth"
		firstDate = time.Date(timeNow.Year(), timeNow.Month(), 1, 0, 0, 0, 0, time.Local)
		query     string
	)

	query = fmt.Sprintf(`
		SELECT COUNT(er.id) as count_reimbursement 
		FROM %s er 
		INNER JOIN %s b ON b.id = er.benefit_id 
		WHERE 
		er.created_at >= $1 AND 
		er.created_at <= $2 AND
		er.status = 'Approved' AND 
		er.deleted = FALSE AND 
		b.deleted = FALSE `,
		EmployeeReimbursementDAO.TableName, BenefitDAO.TableName)

	param := []interface{}{firstDate, timeNow}
	dbResult := db.QueryRow(query, param...)
	dbErr := dbResult.Scan(&count)
	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErr)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dashboardDAO) GetLeaveInMonth(db *sql.DB, timeNow time.Time) (count int64, err errorModel.ErrorModel) {
	var (
		funcName  = "GetLeaveInMonth"
		firstDate = time.Date(timeNow.Year(), timeNow.Month(), 1, 0, 0, 0, 0, time.Local)
		query     string
	)

	query = fmt.Sprintf(`
		SELECT COUNT(el.id) as count_leave 
		FROM %s el 
		WHERE 
		el.deleted = FALSE AND 
		el.status = 'Approved' AND 
		el.id IN 
			(SELECT id
				FROM (
				  SELECT 
					id,
					JSONB_ARRAY_ELEMENTS("date"::jsonb)::VARCHAR::TIMESTAMP AS value
			  FROM employee_leave
			) AS subquery
			WHERE value >= $1 and value <= $2 GROUP BY id) AND
		(el.type = 'leave' OR el.type = 'permit' OR el.type = 'sick-leave') `,
		EmployeeLeaveDAO.TableName)

	param := []interface{}{firstDate, timeNow}
	dbResult := db.QueryRow(query, param...)
	dbErr := dbResult.Scan(&count)
	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErr)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dashboardDAO) GetAbsentInLastPeriod(db *sql.DB, param []in.SearchByParam) (avg float64, err errorModel.ErrorModel) {
	var (
		funcName = "GetAbsentInLastPeriod"
		query    string
		addQuery []string
	)

	query = fmt.Sprintf(`
			SELECT 
			    ROUND(AVG(percent_absent), 2) 
			FROM "%s" `,
		"absent")

	for i := 0; i < len(param); i++ {
		switch param[i].SearchKey {
		case "period_start", "period_end":
			addQuery = append(addQuery, fmt.Sprintf(` %s = '%s'::TIMESTAMP `, param[i].SearchKey, param[i].SearchValue))
		default:
		}
	}

	if len(addQuery) != 2 {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errors.New("must period start and end"))
		return
	}

	query += fmt.Sprintf(` WHERE %s AND %s `, addQuery[0], addQuery[1])
	dbResult := db.QueryRow(query)
	dbErr := dbResult.Scan(&avg)
	if dbErr != nil && dbErr.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErr)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
