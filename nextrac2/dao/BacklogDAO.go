package dao

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type BacklogDAOInterface interface {
	GetListParentBacklog(*sql.DB, in.GetListDataDTO, []in.SearchByParam, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) ([]interface{}, errorModel.ErrorModel)
	GetCountParentBacklog(*sql.DB, []in.SearchByParam, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) (int, errorModel.ErrorModel)
	GetListDetailBacklog(*sql.DB, in.GetListDataDTO, []in.SearchByParam, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) ([]interface{}, errorModel.ErrorModel)
	GetCountDetailBacklog(*sql.DB, []in.SearchByParam, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) (int, errorModel.ErrorModel)
	ViewDetailBacklog(*sql.DB, repository.BacklogModel) (repository.BacklogModel, errorModel.ErrorModel)
	GetDetailBacklogForUpdateOrDelete(*sql.DB, repository.BacklogModel) (repository.BacklogModel, errorModel.ErrorModel)
	UpdateBacklog(*sql.Tx, repository.BacklogModel) errorModel.ErrorModel
	DeleteBacklog(*sql.Tx, repository.BacklogModel) errorModel.ErrorModel
}

type backlogDAO struct {
	AbstractDAO
}

var BacklogDAO = backlogDAO{}.New()

func (input backlogDAO) New() (output backlogDAO) {
	output.FileName = "BacklogDAO.go"
	output.TableName = "backlog"
	return
}

func (input backlogDAO) GetListParentBacklog(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	query = fmt.Sprintf(
		`SELECT 
		b.sprint as sprint, sum(b.mandays) as total_mandays 
		FROM %s b 
		INNER JOIN %s e ON b.employee_id = e.id
		INNER JOIN %s ev ON e.id = ev.employee_id
		INNER JOIN %s d ON e.department_id = d.id `,
		input.TableName, EmployeeDAO.TableName, EmployeeVariableDAO.TableName,
		DepartmentDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	additionalWhere += " AND e.deleted = FALSE AND ev.deleted = FALSE AND d.deleted = FALSE "
	additionalWhere += " GROUP BY b.sprint"
	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.BacklogModel
			dbError := rows.Scan(
				&temp.Sprint, &temp.TotalMandays,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "b.id"},
			Deleted:   FieldStatus{FieldName: "b.deleted"},
			CreatedBy: FieldStatus{FieldName: "b.created_by", Value: createdBy},
		})
}

func (input backlogDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		if (*searchByParam)[i].SearchKey == "pic" {
			(*searchByParam)[i].SearchKey = "CONCAT(e.first_name, ' ', e.last_name)"
		} else {
			(*searchByParam)[i].SearchKey = "b." + (*searchByParam)[i].SearchKey
		}
	}

	if userParam.OrderBy == "id" {
		userParam.OrderBy = "b.id"
	}

}

func (input backlogDAO) GetCountParentBacklog(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		funcName        = "GetCountParentBacklog"
		queryCreated    string
		additionalWhere string
		index           = 1
	)

	if createdBy > 0 {
		queryCreated = fmt.Sprintf(` AND b.created_by = %d `, createdBy)
	}

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, false) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	query := fmt.Sprintf(
		`SELECT COUNT (sub_query.sprint_time) FROM 
		(SELECT b.sprint as sprint_time 
		FROM %s b 
			INNER JOIN %s e ON b.employee_id = e.id
			INNER JOIN %s ev ON e.id = ev.employee_id
			INNER JOIN %s d ON e.department_id = d.id 
		WHERE 
		    b.deleted = FALSE AND e.deleted = FALSE AND ev.deleted = FALSE AND 
		    d.deleted = FALSE %s %s
		GROUP BY b.sprint) AS sub_query `,
		input.TableName, EmployeeDAO.TableName, EmployeeVariableDAO.TableName,
		DepartmentDAO.TableName, queryCreated, additionalWhere)

	var queryParam []interface{}
	for i := 0; i < len(searchByParam); i++ {
		queryParam = append(queryParam, searchByParam[i].SearchValue)
	}

	if len(searchByParam) > 0 {
		query += " WHERE \n"
		for i := 0; i < len(searchByParam); i++ {
			if i == 0 && len(searchByParam) > 1 {
				query += " ( "
			}

			if searchByParam[i].DataType == "enum" {
				searchByParam[i].SearchKey = "cast( " + searchByParam[i].SearchKey + " AS VARCHAR)"
			}

			if searchByParam[i].SearchOperator == "like" {
				searchByParam[i].SearchKey = "LOWER(" + searchByParam[i].SearchKey + ")"
				searchByParam[i].SearchValue = strings.ToLower(searchByParam[i].SearchValue)
				searchByParam[i].SearchValue = "%" + searchByParam[i].SearchValue + "%"
			}

			operator := searchByParam[i].SearchOperator
			if searchByParam[i].SearchOperator == "eq" {
				operator = "="
			}

			query += " " + searchByParam[i].SearchKey + " " + operator + " $" + strconv.Itoa(index) + " "
			if i < len(searchByParam)-1 {
				if searchByParam[i].SearchType == constanta.Search {
					query += "OR "
				} else if searchByParam[i].SearchType == constanta.Filter {
					query += "AND "
				}
			}

			index++

			if i == len(searchByParam)-1 && len(searchByParam) > 1 {
				query += " ) "
			}
		}
	}

	results := db.QueryRow(query, queryParam...)
	errorS := results.Scan(&result)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogDAO) GetListDetailBacklog(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	query = fmt.Sprintf(
		`SELECT
		b.id as id, b.layer_1 as layer_1, b.layer_2 as layer_2,
		b.layer_3 as layer_3, b.redmine_number as redmine_number,
		b.sprint as sprint, CONCAT(e.first_name, ' ', e.last_name) as pic,
		b.status as status, b.mandays,
		b.updated_at as updated_at, d.id
		FROM %s b 
			INNER JOIN %s e ON b.employee_id = e.id
			INNER JOIN %s ev ON e.id = ev.employee_id
			INNER JOIN %s d ON e.department_id = d.id `,
		input.TableName, EmployeeDAO.TableName, EmployeeVariableDAO.TableName,
		DepartmentDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	additionalWhere += " AND e.deleted = FALSE AND ev.deleted = FALSE AND d.deleted = FALSE "
	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.BacklogModel
			dbError := rows.Scan(
				&temp.ID,
				&temp.Layer1, &temp.Layer2,
				&temp.Layer3, &temp.RedmineNumber,
				&temp.Sprint, &temp.EmployeeName,
				&temp.Status, &temp.Mandays,
				&temp.UpdatedAt, &temp.DepartmentId,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "b.id"},
			Deleted:   FieldStatus{FieldName: "b.deleted"},
			CreatedBy: FieldStatus{FieldName: "b.created_by", Value: createdBy},
		})
}

func (input backlogDAO) GetListDetailBacklogBySprint(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, repo repository.BacklogModel) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	query = fmt.Sprintf(
		`SELECT 
			b.id, b.redmine_number as redmine_number 
		FROM %s b 
		LEFT JOIN %s e ON b.employee_id = e.id 
		LEFT JOIN %s d ON e.department_id = d.id `,
		input.TableName, EmployeeDAO.TableName, DepartmentDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	additionalWhere += fmt.Sprintf(`AND b.sprint = '%s' AND b.created_by = %d`, repo.Sprint.String, repo.CreatedBy.Int64)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.BacklogModel
			dbError := rows.Scan(
				&temp.ID,
				&temp.RedmineNumber,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "b.id"},
			Deleted:   FieldStatus{FieldName: "b.deleted"},
			CreatedBy: FieldStatus{FieldName: "b.created_by", Value: createdBy},
		})
}

func (input backlogDAO) GetCountDetailBacklog(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, _ map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		tableName       string
		scopeDBCount    = make(map[string]applicationModel.MappingScopeDB)
	)

	tableName = fmt.Sprintf(`
		%s b
		INNER JOIN %s e ON b.employee_id = e.id
		INNER JOIN %s ev ON e.id = ev.employee_id
		INNER JOIN %s d ON e.department_id = d.id `,
		input.TableName, EmployeeDAO.TableName, EmployeeVariableDAO.TableName,
		DepartmentDAO.TableName)

	scopeDBCount[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{Count: "employee_id"}
	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDBCount, false) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	additionalWhere += " AND e.deleted = FALSE AND ev.deleted = FALSE AND d.deleted = FALSE "
	for i := 0; i < len(searchByParam); i++ {
		if (searchByParam)[i].SearchKey == "pic" {
			(searchByParam)[i].SearchKey = "e.name"
		} else {
			(searchByParam)[i].SearchKey = "b." + (searchByParam)[i].SearchKey
		}
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, additionalWhere, DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "b.id", IsCheck: true},
		Deleted:   FieldStatus{FieldName: "b.deleted", IsCheck: true},
		CreatedBy: FieldStatus{FieldName: "b.created_by", Value: createdBy},
	})
}

func (input backlogDAO) ViewDetailBacklog(db *sql.DB, userParam repository.BacklogModel) (result repository.BacklogModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewDetailBacklog"
		query    string
	)

	query = fmt.Sprintf(`
	SELECT 
		b.id, b.layer_1, b.layer_2, b.layer_3, b.layer_4, b.layer_5,
		b.subject, b.feature, b.redmine_number, b.sprint,
		b.sprint_name, b.reference_ticket,
		e.id, CONCAT(e.first_name, ' ', e.last_name) as pic, b.status, b.estimate_time as mandays_done,
		b.mandays as mandays, b.description, 
		b.flow_changed, b.additional_data, b.note,
		b.url, b.page, d.id, d.name,
		fu.host, fu.path, fu.file_name,
		b.tracker,
		b.updated_at, uc.nt_username as created_name, uu.nt_username as updated_name,
		b.created_by, b.created_at 
		FROM %s b 
			LEFT JOIN %s e ON e.id = b.employee_id 
			LEFT JOIN %s d ON d.id = e.department_id 
			LEFT JOIN "%s" uc ON uc.id = b.created_by
			LEFT JOIN "%s" uu ON uu.id = b.updated_by
			LEFT JOIN %s fu ON b.file_upload_id = fu.id 
		WHERE 
		b.id = $1 AND b.deleted = FALSE AND e.deleted = FALSE AND d.deleted = FALSE `,
		input.TableName, EmployeeDAO.TableName, DepartmentDAO.TableName, UserDAO.TableName,
		UserDAO.TableName, FileUploadDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.Layer1, &result.Layer2, &result.Layer3, &result.Layer4, &result.Layer5,
		&result.Subject, &result.Feature, &result.RedmineNumber, &result.Sprint,
		&result.SprintName, &result.ReferenceTicket,
		&result.EmployeeId, &result.EmployeeName, &result.Status, &result.EstimateTime,
		&result.Mandays, &result.Description,
		&result.FlowChanged, &result.AdditionalData, &result.Note,
		&result.Url, &result.Page, &result.DepartmentId, &result.DepartmentName,
		&result.FileUploadData.Host, &result.FileUploadData.Path, &result.FileUploadData.FileName,
		&result.Tracker,
		&result.UpdatedAt, &result.CreatedName, &result.UpdatedName,
		&result.CreatedBy, &result.CreatedAt)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogDAO) GetDetailBacklogForUpdateOrDelete(db *sql.DB, userParam repository.BacklogModel) (result repository.BacklogModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetDetailBacklogForUpdateOrDelete"
	)

	query := fmt.Sprintf(`
		SELECT 
			b.id, b.updated_at, b.created_by, b.redmine_number 
		FROM %s b 
		WHERE b.id = $1 AND b.deleted = FALSE FOR UPDATE`, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy, &result.RedmineNumber)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogDAO) DeleteBacklog(db *sql.Tx, userParam repository.BacklogModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteBacklog"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET
			deleted = $1, updated_by = $2, 
			updated_at = $3, updated_client = $4,
			redmine_number = $5 
		WHERE
			id = $6 `,
		input.TableName)

	param := []interface{}{
		true,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.RedmineNumber.Int64,
		userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}
	return
}

func (input backlogDAO) DeleteBacklogBySprint(db *sql.Tx, userParam repository.BacklogModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteBacklogBySprint"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET
			deleted = $1, updated_by = $2, 
			updated_at = $3, updated_client = $4 
		WHERE
			sprint = $5 `,
		input.TableName)

	if contextModel.AuthAccessTokenModel.ResourceUserID > 0 {
		query += fmt.Sprintf(` AND created_by = %d`, contextModel.AuthAccessTokenModel.ResourceUserID)
	}

	param := []interface{}{
		true,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.Sprint.String,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}
	return
}

func (input backlogDAO) UpdateBacklog(tx *sql.Tx, userParam repository.BacklogModel) errorModel.ErrorModel {
	var (
		funcName = "UpdateBacklog"
		query    string
		index    = 25
	)

	query = fmt.Sprintf(`
	UPDATE %s
	SET
	layer_1 = $1,
		sprint = $2,
		employee_id = $3,
		estimate_time = $4,
		updated_by = $5,
		updated_at = $6,
		updated_client = $7,
		status = $8,
		redmine_number = $9,
		mandays = $10,
		flow_changed = $11,
		additional_data = $12,
		note = $13,
		page = $14,
		url = $15,
		layer_4 = $16,
		layer_5 = $17,
		layer_2 = $18,
		sprint_name = $19,
		layer_3 = $20,
		feature = $21,
		subject = $22,
		reference_ticket = $23,
		tracker = $24 `, input.TableName)

	// convert estimate time value
	estimateStringValue := fmt.Sprintf(`%.4f`, userParam.EstimateTime.Float64)
	estimateFloatValue, _ := strconv.ParseFloat(estimateStringValue, 64)

	param := []interface{}{
		userParam.Layer1.String,
		userParam.Sprint.String,
		userParam.EmployeeId.Int64,
		estimateFloatValue,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
	}

	// Optional Field Status
	if util.IsStringEmpty(userParam.Status.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Status.String)
	}

	// Optional Field Redmine Number
	if userParam.RedmineNumber.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.RedmineNumber.Int64)
	}

	// Optional Field Mandays
	if userParam.Mandays.Float64 <= 0 {
		param = append(param, nil)
	} else {
		mandaysStringValue := fmt.Sprintf(`%.4f`, userParam.Mandays.Float64)
		mandaysFloatValue, _ := strconv.ParseFloat(mandaysStringValue, 64)
		param = append(param, mandaysFloatValue)
	}

	// Optional Field Flow Changed
	if util.IsStringEmpty(userParam.FlowChanged.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.FlowChanged.String)
	}

	// Optional Field Additional Data
	if util.IsStringEmpty(userParam.AdditionalData.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.AdditionalData.String)
	}

	// Optional Field Note
	if util.IsStringEmpty(userParam.Note.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Note.String)
	}

	// Optional Field Page
	if util.IsStringEmpty(userParam.Page.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Page.String)
	}

	// Optional Field Url
	if util.IsStringEmpty(userParam.Url.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Url.String)
	}

	// Optional Field Layer 4
	if util.IsStringEmpty(userParam.Layer4.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer4.String)
	}

	// Optional Field Layer 5
	if util.IsStringEmpty(userParam.Layer5.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer5.String)
	}

	// Optional Field Layer 2
	if util.IsStringEmpty(userParam.Layer2.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer2.String)
	}

	// Optional Field Sprint Name
	if util.IsStringEmpty(userParam.SprintName.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.SprintName.String)
	}

	// Optional Field Layer 3
	if util.IsStringEmpty(userParam.Layer3.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer3.String)
	}

	// Optional Field Feature
	if userParam.Feature.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Feature.Int64)
	}

	// Optional Field Subject
	if util.IsStringEmpty(userParam.Subject.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Subject.String)
	}

	// Optional Field Reference Ticket
	if userParam.ReferenceTicket.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.ReferenceTicket.Int64)
	}

	// Optional Field Tracker
	if util.IsStringEmpty(userParam.Tracker.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Tracker.String)
	}

	if userParam.FileUploadId.Int64 > 0 {
		param = append(param, userParam.FileUploadId.Int64)
		query += fmt.Sprintf(`, file_upload_id = $%d`, index)
		index++
	}

	query += fmt.Sprintf(` WHERE id = $%d `, index)

	param = append(param, userParam.ID.Int64)

	stmt, dbError := tx.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		if pqError, ok := dbError.(*pgconn.PgError); ok {
			if pqError.Code == "22003" {
				return errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.BacklogRule1, constanta.Mandays, "")
			}
		}

		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input backlogDAO) UpdateStatusBacklog(tx *sql.Tx, userParam repository.BacklogModel) errorModel.ErrorModel {
	var (
		funcName = "UpdateStatusBacklog"
		query    string
	)

	query = fmt.Sprintf(`
	UPDATE % s
	SET
	status = $1,
		updated_at = $2,
		updated_by = $3,
		updated_client = $4
	WHERE
	id = $5
	`, input.TableName)

	param := []interface{}{
		userParam.Status.String,
		userParam.UpdatedAt.Time,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.ID.Int64,
	}

	stmt, dbError := tx.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input backlogDAO) InsertBacklog(tx *sql.Tx, userParam repository.BacklogModel) (idInserted int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertBacklog"
		query    string
	)

	query = fmt.Sprintf(`
	INSERT
	INTO %s(
		layer_1,
		sprint,
		employee_id,
		estimate_time,
		updated_by,
		updated_at,
		updated_client,
		created_by,
		created_at,
		created_client,
		status,
		redmine_number,
		mandays,
		flow_changed,
		additional_data,
		note,
		page,
		url,
		layer_4,
		layer_5,
		layer_2,
		sprint_name,
		layer_3,
		feature,
		subject,
		reference_ticket,
		tracker,
		file_upload_id)
	VALUES(
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28
	) RETURNING
	id
	`, input.TableName)

	// convert estimate time value
	estimateStringValue := fmt.Sprintf(`%.4f`, userParam.EstimateTime.Float64)
	estimateFloatValue, _ := strconv.ParseFloat(estimateStringValue, 64)

	param := []interface{}{
		userParam.Layer1.String,
		userParam.Sprint.String,
		userParam.EmployeeId.Int64,
		estimateFloatValue,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
	}

	// Optional Field Status
	if util.IsStringEmpty(userParam.Status.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Status.String)
	}

	// Optional Field Redmine Number
	if userParam.RedmineNumber.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.RedmineNumber.Int64)
	}

	// Optional Field Mandays
	if userParam.Mandays.Float64 <= 0 {
		param = append(param, nil)
	} else {
		mandaysStringValue := fmt.Sprintf(`%.4f`, userParam.Mandays.Float64)
		mandaysFloatValue, _ := strconv.ParseFloat(mandaysStringValue, 64)
		param = append(param, mandaysFloatValue)
	}

	// Optional Field Flow Changed
	if util.IsStringEmpty(userParam.FlowChanged.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.FlowChanged.String)
	}

	// Optional Field Additional Data
	if util.IsStringEmpty(userParam.AdditionalData.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.AdditionalData.String)
	}

	// Optional Field Note
	if util.IsStringEmpty(userParam.Note.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Note.String)
	}

	// Optional Field Page
	if util.IsStringEmpty(userParam.Page.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Page.String)
	}

	// Optional Field Url
	if util.IsStringEmpty(userParam.Url.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Url.String)
	}

	// Optional Field Layer 4
	if util.IsStringEmpty(userParam.Layer4.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer4.String)
	}

	// Optional Field Layer 5
	if util.IsStringEmpty(userParam.Layer5.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer5.String)
	}

	// Optional Field Layer 2
	if util.IsStringEmpty(userParam.Layer2.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer2.String)
	}

	// Optional Field Sprint Name
	if util.IsStringEmpty(userParam.SprintName.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.SprintName.String)
	}

	// Optional Field Layer 3
	if util.IsStringEmpty(userParam.Layer3.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Layer3.String)
	}

	// Optional Field Feature
	if userParam.Feature.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Feature.Int64)
	}

	// Optional Field Subject
	if util.IsStringEmpty(userParam.Subject.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Subject.String)
	}

	// Optional Field Reference Ticket
	if userParam.ReferenceTicket.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.ReferenceTicket.Int64)
	}

	// Optional Field Reference Tracker
	if util.IsStringEmpty(userParam.Tracker.String) {
		param = append(param, nil)
	} else {
		param = append(param, userParam.Tracker.String)
	}

	// Optional Field File Upload ID
	if userParam.FileUploadId.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.FileUploadId.Int64)
	}

	results := tx.QueryRow(query, param...)
	dbError := results.Scan(&idInserted)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		if pqError, ok := dbError.(*pgconn.PgError); ok {
			if pqError.Code == "22003" {
				err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, constanta.BacklogRule1, constanta.Mandays, "")
				return
			}
		}
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input backlogDAO) UpdatePaymentOnBacklog(db *sql.Tx, ticket []int64, userParam repository.BacklogModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdatePaymentOnBacklog"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET 
		payment_status = $1, updated_by = $2, 
		updated_at = $3, updated_client = $4 `,
		input.TableName)

	for i := 0; i < len(ticket); i++ {
		if i == 0 {
			query += fmt.Sprintf(` WHERE redmine_number IN( `)
		}

		if len(ticket)-(i+1) == 0 {
			query += fmt.Sprintf(` %d) `, ticket[i])
		} else {
			query += fmt.Sprintf(` %d, `, ticket[i])
		}
	}

	param := []interface{}{
		true, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input backlogDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (colAdditionalWhere []string) {
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
