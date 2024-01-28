package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type reportHistoryDAO struct {
	AbstractDAO
}

var ReportHistoryDAO = reportHistoryDAO{}.New()

func (input reportHistoryDAO) New() (output reportHistoryDAO) {
	output.FileName = "ReportHistoryDAO.go"
	output.TableName = "report_history"
	return
}

func (input reportHistoryDAO) GetListReportHistory(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(`
		SELECT
		r.id, r.success_ticket, r.created_at, 
		r.created_by, u.nt_username, d.name
		FROM %s r
		LEFT JOIN %s AS d ON d.id = r.department_id
		LEFT JOIN "%s" AS u ON r.created_by = u.id `,
		input.TableName, DepartmentDAO.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ReportHistoryModel
			dbError := rows.Scan(
				&temp.ID, &temp.SuccessTicket, &temp.CreatedAt,
				&temp.CreatedBy, &temp.CreatedName, &temp.DepartmentName,
			)
			return temp, dbError
		}, " ", DefaultFieldMustCheck{
			Deleted:   FieldStatus{FieldName: "r.deleted"},
			CreatedBy: FieldStatus{FieldName: "r.created_by", Value: createdBy},
		})
}

func (input reportHistoryDAO) ViewReportHistory(db *sql.DB, userParam repository.ReportHistoryModel) (result repository.ReportHistoryModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewReportHistory"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT
			r.id, r."data", r.success_ticket, 
			r.created_at, u.nt_username, d.name, 
			r.created_by
		FROM %s r
		LEFT JOIN %s AS d ON d.id = r.department_id
		LEFT JOIN "%s" AS u ON r.created_by = u.id
		WHERE 
		r.id = $1 AND r.deleted = FALSE `,
		input.TableName, DepartmentDAO.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.Data,
		&result.SuccessTicket, &result.CreatedAt,
		&result.CreatedName, &result.DepartmentName,
		&result.CreatedBy,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input reportHistoryDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "r." + (*searchByParam)[i].SearchKey
	}

	userParam.OrderBy = "r." + userParam.OrderBy
}
