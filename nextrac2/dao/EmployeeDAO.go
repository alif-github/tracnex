package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type EmployeeDAOInterface interface {
	InsertEmployee(*sql.Tx, repository.EmployeeModel) (int64, errorModel.ErrorModel)
	UpdateEmployee(*sql.Tx, repository.EmployeeModel) errorModel.ErrorModel
	GetEmployeeForUpdate(*sql.DB, repository.EmployeeModel) (repository.EmployeeModel, errorModel.ErrorModel)
	ViewEmployee(*sql.DB, repository.EmployeeModel, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) (repository.EmployeeModel, errorModel.ErrorModel)
	GetCountEmployee(*sql.DB, []in.SearchByParam, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel)
	DeleteEmployee(*sql.Tx, repository.EmployeeModel) (err errorModel.ErrorModel)
	GetListEmployee(*sql.DB, in.GetListDataDTO, []in.SearchByParam, int64, map[string]interface{}, map[string]applicationModel.MappingScopeDB, bool) (result []interface{}, err errorModel.ErrorModel)
}

type employeeDAO struct {
	AbstractDAO
}

var EmployeeDAO = employeeDAO{}.New()

func (input employeeDAO) New() (output employeeDAO) {
	output.FileName = "EmployeeDAO.go"
	output.TableName = "employee"
	return
}

func (input employeeDAO) InsertEmployee(db *sql.Tx, userParam repository.EmployeeModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertEmployee"
		query    string
	)

	query = fmt.Sprintf(
		`INSERT INTO %s
		(
			id_card, npwp, first_name,
			religion, marital_status, nationality,
			created_by, created_at, created_client, 
			updated_by, updated_at, updated_client, 
			is_have_member, member, education,
			last_name, address_residence, address_tax, 
			date_join, date_out, reason_resignation,
			department_id, employee_position_id, number_of_dependents, 
			status, active, gender, 
		 	place_of_birth, date_of_birth, email, 
		 	phone, "type", mothers_maiden, 
		 	tax_method, file_upload_id
		)
		VALUES
		(
			$1, $2, $3, 
			$4, $5, $6, 
			$7, $8, $9, 
			$10, $11, $12, 
			$13, $14, $15,
			$16, $17, $18, 
			$19, $20, $21, 
			$22, $23, $24,
			$25, $26, $27,
		    $28, $29, $30, 
		    $31, $32, $33,
		    $34, $35
		)
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.IDCard.String, userParam.NPWP.String, userParam.FirstName.String,
		userParam.Religion.String, userParam.MaritalStatus.String, userParam.Nationality.String,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.IsHaveMember.Bool,
	}

	if userParam.IsHaveMember.Bool {
		params = append(params, userParam.Member.String)
	} else {
		params = append(params, nil)
	}

	if userParam.Education.String != "" {
		params = append(params, userParam.Education.String)
	} else {
		params = append(params, nil)
	}

	if userParam.LastName.String != "" {
		params = append(params, userParam.LastName.String)
	} else {
		params = append(params, nil)
	}

	if userParam.AddressResidence.String != "" {
		params = append(params, userParam.AddressResidence.String)
	} else {
		params = append(params, nil)
	}

	if userParam.AddressTax.String != "" {
		params = append(params, userParam.AddressTax.String)
	} else {
		params = append(params, nil)
	}

	if !userParam.DateJoin.Time.IsZero() {
		params = append(params, userParam.DateJoin.Time)
	} else {
		params = append(params, nil)
	}

	if !userParam.DateOut.Time.IsZero() {
		params = append(params, userParam.DateOut.Time)
	} else {
		params = append(params, nil)
	}

	if userParam.ReasonResignation.String != "" {
		params = append(params, userParam.ReasonResignation.String)
	} else {
		params = append(params, nil)
	}

	if userParam.DepartmentId.Int64 > 0 {
		params = append(params, userParam.DepartmentId.Int64)
	} else {
		params = append(params, nil)
	}

	if userParam.PositionID.Int64 > 0 {
		params = append(params, userParam.PositionID.Int64)
	} else {
		params = append(params, nil)
	}

	if userParam.NumberOfDependents.Int64 > 0 {
		params = append(params, userParam.NumberOfDependents.Int64)
	} else {
		params = append(params, nil)
	}

	params = append(params,
		userParam.Status.String, userParam.Active.Bool, userParam.Gender.String,
		userParam.PlaceOfBirth.String, userParam.DateOfBirth.Time)

	if userParam.Email.String != "" {
		params = append(params, userParam.Email.String)
	} else {
		params = append(params, nil)
	}

	if userParam.Phone.String != "" {
		params = append(params, userParam.Phone.String)
	} else {
		params = append(params, nil)
	}

	if userParam.Type.String != "" {
		params = append(params, userParam.Type.String)
	} else {
		params = append(params, nil)
	}

	if userParam.MothersMaiden.String != "" {
		params = append(params, userParam.MothersMaiden.String)
	} else {
		params = append(params, nil)
	}

	if userParam.TaxMethod.String != "" {
		params = append(params, userParam.TaxMethod.String)
	} else {
		params = append(params, nil)
	}

	if userParam.FileUploadID.Int64 > 0 {
		params = append(params, userParam.FileUploadID.Int64)
	} else {
		params = append(params, nil)
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

func (input employeeDAO) UpdateEmployee(tx *sql.Tx, userParam repository.EmployeeModel) errorModel.ErrorModel {
	var (
		funcName = "UpdateEmployee"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		id_card = $1, npwp = $2, first_name = $3, 
		religion = $4, marital_status = $5, nationality = $6, 
		updated_by = $7, updated_client = $8, updated_at = $9, 
		is_have_member = $10, member = $11, education = $12,
		last_name = $13, address_residence = $14, address_tax = $15,
		date_join = $16, date_out = $17, reason_resignation = $18, 
		department_id = $19, employee_position_id = $20, number_of_dependents = $21,
		status = $22, active = $23, gender = $24, 
		place_of_birth = $25, date_of_birth = $26, email = $27, 
		phone = $28, "type" = $29, mothers_maiden = $30, 
		tax_method = $31, file_upload_id = $32
		WHERE id = $33 `,
		input.TableName)

	param = []interface{}{
		userParam.IDCard.String, userParam.NPWP.String, userParam.FirstName.String,
		userParam.Religion.String, userParam.MaritalStatus.String, userParam.Nationality.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.IsHaveMember.Bool,
	}

	if userParam.IsHaveMember.Bool {
		param = append(param, userParam.Member.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Education.String != "" {
		param = append(param, userParam.Education.String)
	} else {
		param = append(param, nil)
	}

	if userParam.LastName.String != "" {
		param = append(param, userParam.LastName.String)
	} else {
		param = append(param, nil)
	}

	if userParam.AddressResidence.String != "" {
		param = append(param, userParam.AddressResidence.String)
	} else {
		param = append(param, nil)
	}

	if userParam.AddressTax.String != "" {
		param = append(param, userParam.AddressTax.String)
	} else {
		param = append(param, nil)
	}

	if !userParam.DateJoin.Time.IsZero() {
		param = append(param, userParam.DateJoin.Time)
	} else {
		param = append(param, nil)
	}

	if !userParam.DateOut.Time.IsZero() {
		param = append(param, userParam.DateOut.Time)
	} else {
		param = append(param, nil)
	}

	if userParam.ReasonResignation.String != "" {
		param = append(param, userParam.ReasonResignation.String)
	} else {
		param = append(param, nil)
	}

	if userParam.DepartmentId.Int64 > 0 {
		param = append(param, userParam.DepartmentId.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.PositionID.Int64 > 0 {
		param = append(param, userParam.PositionID.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.NumberOfDependents.Int64 > 0 {
		param = append(param, userParam.NumberOfDependents.Int64)
	} else {
		param = append(param, nil)
	}

	param = append(param,
		userParam.Status.String, userParam.Active.Bool, userParam.Gender.String,
		userParam.PlaceOfBirth.String, userParam.DateOfBirth.Time)

	if userParam.Email.String != "" {
		param = append(param, userParam.Email.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Phone.String != "" {
		param = append(param, userParam.Phone.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Type.String != "" {
		param = append(param, userParam.Type.String)
	} else {
		param = append(param, nil)
	}

	if userParam.MothersMaiden.String != "" {
		param = append(param, userParam.MothersMaiden.String)
	} else {
		param = append(param, nil)
	}

	if userParam.TaxMethod.String != "" {
		param = append(param, userParam.TaxMethod.String)
	} else {
		param = append(param, nil)
	}

	if userParam.FileUploadID.Int64 > 0 {
		param = append(param, userParam.FileUploadID.Int64)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.ID.Int64)

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

func (input employeeDAO) UpdateEmployeeTimeSheet(tx *sql.Tx, userParam repository.EmployeeModel) errorModel.ErrorModel {
	var (
		funcName = "UpdateEmployeeTimeSheet"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET employee_variable_id = $1 WHERE id = $2 `, input.TableName)
	param := []interface{}{userParam.EmployeeVariableID.Int64, userParam.ID.Int64}
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

func (input employeeDAO) GetEmployeeForUpdate(db *sql.DB, userParam repository.EmployeeModel) (result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeeForUpdate"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
			e.id, e.id_card, e.updated_at, e.created_by,
			CASE WHEN
				(SELECT COUNT (b.id) FROM backlog b WHERE b.employee_id = e.id AND b.deleted = FALSE) > 0
			THEN TRUE ELSE FALSE END is_used 
		FROM %s e 
		WHERE e.id = $1 AND e.deleted = FALSE FOR UPDATE`,
		input.TableName)

	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.IDCard, &result.UpdatedAt,
		&result.CreatedBy, &result.IsUsed)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) GetEmployeeIdByFullName(db *sql.DB, userParam repository.EmployeeModel, _ int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeeIdByFullName"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
			e.id, CONCAT(e.first_name, ' ', e.last_name) as name
		FROM %s e 
		INNER JOIN %s ev ON e.id = ev.employee_id
		INNER JOIN %s d ON e.department_id = d.id
		WHERE 
		    LOWER(CONCAT(e.first_name, ' ', e.last_name)) LIKE $1 AND
		    e.deleted = FALSE AND 
		    ev.deleted = FALSE AND 
		    d.deleted = FALSE `,
		input.TableName, EmployeeVariableDAO.TableName, DepartmentDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		query += " AND " + itemColAdditionalWhere
	}
	query += " LIMIT 1 "

	fmt.Println("Hasil -> ", query)

	param := []interface{}{"%" + userParam.Name.String + "%"}
	dbError := db.QueryRow(query, param...).Scan(&result.ID, &result.Name)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) CheckEmployeeIDByID(db *sql.DB, id int64) (isExists bool, firstName string, err errorModel.ErrorModel) {
	var (
		funcName = "CheckEmployeeIDByID"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT CASE WHEN id > 0 THEN true ELSE false END is_exist, first_name
		FROM %s WHERE id = $1 AND deleted = FALSE `,
		input.TableName)

	param := []interface{}{id}
	dbError := db.QueryRow(query, param...).Scan(&isExists, &firstName)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) GetEmployeeByRedmineID(db *sql.DB, userParam repository.EmployeeModel) (result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeeByRedmineID"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT e.id, e.nik, e.redmine_id, 
		e.name, d.id, d.name as departement_name, 
		e.updated_at, uc.nt_username as created_name, uu.nt_username as updated_name, 
		e.mandays_rate, 
		e.created_by, e.created_at 
		FROM %s e 
			INNER JOIN %s d ON d.id = e.department_id 
			LEFT JOIN "%s" uc ON uc.id = e.created_by
			LEFT JOIN "%s" uu ON uu.id = e.updated_by
		WHERE 
		e.id = $1 AND e.deleted = FALSE AND d.deleted = FALSE `,
		input.TableName, DepartmentDAO.TableName, UserDAO.TableName,
		UserDAO.TableName)

	query = fmt.Sprintf(`SELECT * FROM %s WHERE department_id = $1 AND redmine_id = $2 AND deleted = FALSE `, input.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.NIK, &result.RedmineId,
		&result.Name, &result.DepartmentId, &result.DepartmentName,
		&result.UpdatedAt, &result.CreatedName, &result.UpdatedName,
		&result.MandaysRate, &result.CreatedBy, &result.CreatedAt)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) CheckEmployeeByNIK(db *sql.DB, userParam repository.EmployeeModel) (result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "CheckEmployeeByNIK"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
		e.id, e.id_card, CONCAT(e.first_name, ' ', e.last_name) as employee_name, 
		d.id 
		FROM %s e 
		INNER JOIN %s d ON d.id = e.department_id
		WHERE 
		e.id_card = $1 AND 
		e.id = $2 AND 
		e.deleted = FALSE AND 
		d.deleted = FALSE `,
		input.TableName, DepartmentDAO.TableName)

	params := []interface{}{userParam.IDCard.String, userParam.ID.Int64}
	results := db.QueryRow(query, params...)

	dbError := results.Scan(&result.ID, &result.NIK, &result.Name, &result.DepartmentId)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) GetCountEmployeeTimeSheet(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		tableName       string
		additionalWhere string
	)

	tableName = fmt.Sprintf(` 
		%s e 
		INNER JOIN %s ev ON ev.employee_id = e.id AND ev.deleted = FALSE
		LEFT JOIN %s d ON e.department_id = d.id AND d.deleted = FALSE `,
		input.TableName, EmployeeVariableDAO.TableName, DepartmentDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, false) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	for i, item := range searchByParam {
		if searchByParam[i].SearchKey == "department" {
			searchByParam[i].SearchKey = "d.name"
		} else if searchByParam[i].SearchKey == "department_id" {
			searchByParam[i].SearchKey = "d.id"
		} else if (searchByParam)[i].SearchKey == "nik" {
			(searchByParam)[i].SearchKey = "e.id_card"
		} else if (searchByParam)[i].SearchKey == "name" {
			(searchByParam)[i].SearchKey = "CONCAT(e.first_name, ' ', e.last_name)"
		} else if (searchByParam)[i].SearchKey == "redmine_id" {
			(searchByParam)[i].SearchKey = "ev.redmine_id"
		} else {
			searchByParam[i].SearchKey = "e." + item.SearchKey
		}
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, additionalWhere, DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "e.id"},
		Deleted:   FieldStatus{FieldName: "e.deleted"},
		CreatedBy: FieldStatus{FieldName: "e.created_by", Value: createdBy},
	})
}

func (input employeeDAO) GetCountEmployee(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		tableName       string
		additionalWhere string
	)

	tableName = fmt.Sprintf(` 
		%s e 
		INNER JOIN %s d ON e.department_id = d.id
		LEFT JOIN %s ev ON ev.employee_id = e.id
		LEFT JOIN "%s" u ON e.updated_by = u.id `,
		input.TableName, DepartmentDAO.TableName, EmployeeVariableDAO.TableName,
		UserDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, false) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	additionalWhere += fmt.Sprintf(` AND d.deleted = FALSE `)
	for i, item := range searchByParam {
		if searchByParam[i].SearchKey == "department" {
			searchByParam[i].SearchKey = "d.name"
		} else if searchByParam[i].SearchKey == "department_id" {
			searchByParam[i].SearchKey = "d.id"
		} else if (searchByParam)[i].SearchKey == "nik" {
			(searchByParam)[i].SearchKey = "e.id_card"
		} else if (searchByParam)[i].SearchKey == "name" {
			(searchByParam)[i].SearchKey = "CONCAT(e.first_name, ' ', e.last_name)"
		} else {
			searchByParam[i].SearchKey = "e." + item.SearchKey
		}
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, additionalWhere, DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "e.id"},
		Deleted:   FieldStatus{FieldName: "e.deleted"},
		CreatedBy: FieldStatus{FieldName: "e.created_by", Value: createdBy},
	})
}

func (input employeeDAO) GetListEmployeeTimeSheet(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	query = fmt.Sprintf(`
		SELECT 
			e.id as id, e.id_card as nik, ev.redmine_id as redmine_id, 
			CONCAT(e.first_name, ' ', e.last_name) as name, d.id, d.name as department_name,
			ev.mandays_rate, ev.created_at, ev.updated_at, 
			ev.updated_by, u.nt_username
		FROM %s e
		INNER JOIN %s ev ON ev.employee_id = e.id AND ev.deleted = FALSE
		LEFT JOIN "%s" u ON ev.updated_by = u.id AND u.deleted = FALSE
		LEFT JOIN %s d ON e.department_id = d.id AND d.deleted = FALSE `,
		input.TableName, EmployeeVariableDAO.TableName, UserDAO.TableName,
		DepartmentDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.EmployeeModel
			dbError := rows.Scan(
				&temp.ID, &temp.NIK, &temp.RedmineId,
				&temp.Name, &temp.DepartmentId, &temp.DepartmentName,
				&temp.MandaysRate, &temp.CreatedAt, &temp.UpdatedAt,
				&temp.UpdatedBy, &temp.UpdatedName,
			)
			return temp, dbError
		}, additionalWhere, input.getDefaultMustCheck(createdBy))
}

func (input employeeDAO) GetListEmployee(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isTimeSheet bool) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	if isTimeSheet {
		query = fmt.Sprintf(`
			SELECT 
			e.id as id, e.id_card as nik, CONCAT(e.first_name, ' ', e.last_name) as name, 
			d.id, d.name as department_name, ev.created_at, 
			ev.updated_at, ev.updated_by, u.nt_username, 
			ev.redmine_id
			FROM %s e
			INNER JOIN %s d ON e.department_id = d.id			
			LEFT JOIN %s ev ON ev.employee_id = e.id
			LEFT JOIN "%s" u ON ev.updated_by = u.id `,
			input.TableName, DepartmentDAO.TableName, EmployeeVariableDAO.TableName,
			UserDAO.TableName)
	} else {
		query = fmt.Sprintf(`
			SELECT 
			e.id as id, e.id_card as nik, CONCAT(e.first_name, ' ', e.last_name) as name, 
			d.id, d.name as department_name, e.created_at, 
			e.updated_at, e.updated_by, u.nt_username, 
			ev.redmine_id
			FROM %s e
			INNER JOIN %s d ON e.department_id = d.id 
			LEFT JOIN %s ev ON ev.employee_id = e.id
			LEFT JOIN "%s" u ON e.updated_by = u.id `,
			input.TableName, DepartmentDAO.TableName, EmployeeVariableDAO.TableName,
			UserDAO.TableName)
	}

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	additionalWhere += fmt.Sprintf(` AND d.deleted = FALSE `)
	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.EmployeeModel
			dbError := rows.Scan(
				&temp.ID, &temp.IDCard, &temp.Name,
				&temp.DepartmentId, &temp.DepartmentName,
				&temp.CreatedAt, &temp.UpdatedAt,
				&temp.UpdatedBy, &temp.UpdatedName,
				&temp.RedmineId,
			)
			return temp, dbError
		}, additionalWhere, input.getDefaultMustCheck(createdBy))
}

func (input employeeDAO) GetListEmployeeReport(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	query = fmt.Sprintf(`
			SELECT 
			e.id as id, e.id_card as nik, CONCAT(e.first_name, ' ', e.last_name) as name, 
			d.id, d.name as department_name, ev.created_at, 
			ev.updated_at, ev.updated_by, u.nt_username, 
			ev.redmine_id
			FROM %s e
			INNER JOIN %s d ON e.department_id = d.id			
			INNER JOIN %s ev ON ev.employee_id = e.id
			LEFT JOIN "%s" u ON ev.updated_by = u.id `,
		input.TableName, DepartmentDAO.TableName, EmployeeVariableDAO.TableName,
		UserDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += " AND " + itemColAdditionalWhere
	}

	//--- Redmine ID Greater Than 0
	additionalWhere += " AND ev.redmine_id > 0 "
	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.EmployeeModel
			dbError := rows.Scan(
				&temp.ID, &temp.IDCard, &temp.Name,
				&temp.DepartmentId, &temp.DepartmentName,
				&temp.CreatedAt, &temp.UpdatedAt,
				&temp.UpdatedBy, &temp.UpdatedName,
				&temp.RedmineId,
			)
			return temp, dbError
		}, additionalWhere, input.getDefaultMustCheck(createdBy))
}

func (input employeeDAO) GetListMember(db *sql.DB, idEmployee int64, idMember []int64, isAllMember bool) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
		userParam       in.GetListDataDTO
	)

	query = fmt.Sprintf(`
		SELECT id, first_name, last_name 
		FROM %s `,
		input.TableName)

	if !isAllMember {
		for idx, itemMember := range idMember {
			if idx == 0 {
				additionalWhere += fmt.Sprintf(` AND id IN (%d`, itemMember)
			} else {
				additionalWhere += fmt.Sprintf(`%d`, itemMember)
			}

			if len(idMember)-(idx+1) == 0 {
				additionalWhere += ") "
			} else {
				additionalWhere += ", "
			}
		}
	} else {
		additionalWhere += fmt.Sprintf(` AND id <> %d `, idEmployee)
	}

	userParam = in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:  -99,
			Limit: -99,
		},
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, nil,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.MemberList
			dbError := rows.Scan(&temp.ID, &temp.FirstName, &temp.LastName)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{}.GetDefaultField(false, 0))
}

func (input employeeDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		if (*searchByParam)[i].SearchKey == "department" {
			(*searchByParam)[i].SearchKey = "d.name"
		} else if (*searchByParam)[i].SearchKey == "department_id" {
			(*searchByParam)[i].SearchKey = "d.id"
		} else if (*searchByParam)[i].SearchKey == "nik" {
			(*searchByParam)[i].SearchKey = "e.id_card"
		} else if (*searchByParam)[i].SearchKey == "redmine_id" {
			(*searchByParam)[i].SearchKey = "ev.redmine_id"
		} else if (*searchByParam)[i].SearchKey == "name" {
			(*searchByParam)[i].SearchKey = "CONCAT(e.first_name, ' ', e.last_name)"
		} else {
			(*searchByParam)[i].SearchKey = "e." + (*searchByParam)[i].SearchKey
		}
	}

	switch userParam.OrderBy {
	case "updated_name", "updated_name ASC", "updated_name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "u.nt_username " + strSplit[1]
		} else {
			userParam.OrderBy = "u.nt_username"
		}
		break
	case "department", "department ASC", "department DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "d.name " + strSplit[1]
		} else {
			userParam.OrderBy = "d.name"
		}
		break
	case "nik", "nik ASC", "nik DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "e.id_card " + strSplit[1]
		} else {
			userParam.OrderBy = "e.id_card"
		}
		break
	case "redmine_id", "redmine_id ASC", "redmine_id DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "ev.redmine_id " + strSplit[1]
		} else {
			userParam.OrderBy = "ev.redmine_id"
		}
		break
	case "name", "name ASC", "name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "CONCAT(e.first_name, ' ', e.last_name) " + strSplit[1]
		} else {
			userParam.OrderBy = "CONCAT(e.first_name, ' ', e.last_name)"
		}
		break
	case "updated_at", "updated_at ASC", "updated_at DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "e.updated_at " + strSplit[1]
		} else {
			userParam.OrderBy = "e.updated_at"
		}
		break
	}
}

func (input employeeDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "e.id"},
		Deleted:   FieldStatus{FieldName: "e.deleted"},
		CreatedBy: FieldStatus{FieldName: "e.created_by", Value: createdBy},
	}
}

func (input employeeDAO) DeleteEmployee(db *sql.Tx, userParam repository.EmployeeModel) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteEmployee"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		deleted = $1, id_card = $2, updated_by = $3, 
		updated_at = $4, updated_client = $5
		WHERE 
		id = $6 `,
		input.TableName)

	param := []interface{}{
		true, userParam.IDCard.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
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

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) ViewEmployee(db *sql.DB, userParam repository.EmployeeModel, _ int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewEmployee"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
		e.id, e.id_card, e.first_name, 
		e.last_name, e.gender, e.place_of_birth, 
		e.date_of_birth, e.address_residence, e.npwp,
		e.address_tax, e.religion, e.date_join, 
		e.date_out, e.reason_resignation, e.status,
		e.active, d.id, d.name as departement_name,
		ep.name, e.email, e.phone,
		e.type, e.marital_status, e.number_of_dependents,
		e.education, e.nationality, e.mothers_maiden, 
		e.tax_method, CONCAT(fu.host, fu.path) as url_photo, e.created_by,
		uc.nt_username as created_name, e.created_at, e.updated_by,
		uu.nt_username as updated_name, e.updated_at, eb.bpjs_no, 
		eb.bpjs_tk_no, e.is_have_member, e.member, 
		el.level, eg.grade, e.employee_position_id, 
		e.active, el.id as level_id, eg.id as grade_id
		FROM %s e 
			LEFT JOIN %s d ON d.id = e.department_id AND d.deleted = FALSE
			LEFT JOIN %s eb ON eb.employee_id = e.id AND eb.deleted = FALSE
			LEFT JOIN %s el ON eb.employee_level_id = el.id AND el.deleted = FALSE
			LEFT JOIN %s eg ON eb.employee_grade_id = eg.id AND eg.deleted = FALSE
			LEFT JOIN %s ep ON e.employee_position_id = ep.id AND ep.deleted = FALSE
			LEFT JOIN %s fu ON e.file_upload_id = fu.id AND fu.deleted = FALSE
			LEFT JOIN "%s" uc ON uc.id = e.created_by AND uc.deleted = FALSE
			LEFT JOIN "%s" uu ON uu.id = e.updated_by AND uu.deleted = FALSE
		WHERE 
		e.id = $1 AND e.deleted = FALSE `,
		input.TableName, DepartmentDAO.TableName, EmployeeBenefitsDAO.TableName,
		"employee_level", "employee_grade", "employee_position",
		FileUploadDAO.TableName, UserDAO.TableName, UserDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //-- Scope check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		query += " AND " + itemColAdditionalWhere
	}

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.IDCard, &result.FirstName,
		&result.LastName, &result.Gender, &result.PlaceOfBirth,
		&result.DateOfBirth, &result.AddressResidence, &result.NPWP,
		&result.AddressTax, &result.Religion, &result.DateJoin,
		&result.DateOut, &result.ReasonResignation, &result.Status,
		&result.Active, &result.DepartmentId, &result.DepartmentName,
		&result.Position, &result.Email, &result.Phone,
		&result.Type, &result.MaritalStatus, &result.NumberOfDependents,
		&result.Education, &result.Nationality, &result.MothersMaiden,
		&result.TaxMethod, &result.Photo, &result.CreatedBy,
		&result.CreatedName, &result.CreatedAt, &result.UpdatedBy,
		&result.UpdatedName, &result.UpdatedAt, &result.BPJS,
		&result.BPJSTk, &result.IsHaveMember, &result.Member,
		&result.Level, &result.Grade, &result.PositionID,
		&result.Active, &result.LevelID, &result.GradeID,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) ViewEmployeeTimeSheet(db *sql.DB, userParam repository.EmployeeModel) (result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewEmployeeTimeSheet"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
		e.id, e.id_card, CONCAT(e.first_name, '', e.last_name) as name,
		d.id, d.name as departement_name, ev.mandays_rate, 
		ev.redmine_id, ev.created_by, uc.nt_username as created_name,
		ev.created_at, ev.updated_by, uu.nt_username as updated_name,
		ev.updated_at 
		FROM %s e 
			INNER JOIN %s ev ON ev.employee_id = e.id AND ev.deleted = FALSE
			LEFT JOIN %s d ON d.id = e.department_id AND d.deleted = FALSE
			LEFT JOIN "%s" uc ON uc.id = ev.created_by AND uc.deleted = FALSE
			LEFT JOIN "%s" uu ON uu.id = ev.updated_by AND uu.deleted = FALSE
		WHERE 
		e.id = $1 AND e.deleted = FALSE `,
		input.TableName, EmployeeVariableDAO.TableName,
		DepartmentDAO.TableName, UserDAO.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.IDCard, &result.Name,
		&result.DepartmentId, &result.DepartmentName, &result.MandaysRate,
		&result.RedmineId, &result.CreatedBy, &result.CreatedName,
		&result.CreatedAt, &result.UpdatedBy, &result.UpdatedName,
		&result.UpdatedAt,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) GetRedmineIDEmployeeByDepartmentID(db *sql.DB, departmentID int64) (redmineID []int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetRedmineIDEmployeeByDepartmentID"
		query    string
		result   sql.NullString
	)

	query = fmt.Sprintf(`
		SELECT 
			CASE WHEN COUNT(ev.redmine_id) > 0 
				THEN JSONB_AGG(ev.redmine_id) 
				ELSE null
			END col_id
		FROM %s e
		INNER JOIN %s ev ON ev.employee_id = e.id
		WHERE 
		    e.department_id = $1 AND e.deleted = FALSE AND ev.deleted = FALSE `,
		input.TableName, EmployeeVariableDAO.TableName)

	params := []interface{}{departmentID}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	if result.String != "" {
		_ = json.Unmarshal([]byte(result.String), &redmineID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) GetByUserId(db *sql.DB, id int64) (result repository.EmployeeModel, errModel errorModel.ErrorModel) {
	funcName := "GetByUserId"

	query := `SELECT 
				e.id, e.member, eb.employee_level_id,
				eb.employee_grade_id 
			FROM ` + input.TableName + ` AS e
			LEFT JOIN employee_benefits AS eb 
				ON e.id = eb.employee_id
			INNER JOIN "user" AS u 
				ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
				OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
			WHERE
				u.id = $1 AND
				u.deleted = FALSE AND 
				e.deleted = FALSE`

	row := db.QueryRow(query, id)
	err := row.Scan(
		&result.ID, &result.Member, &result.LevelID,
		&result.GradeID,
	)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (colAdditionalWhere []string) {
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

func (input employeeDAO) CheckEmployeeByNIKParamOnly(db *sql.DB, idCard string) (isExist bool, result repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		funcName = "CheckEmployeeByNIK"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
		    CASE WHEN e.id > 0 
		        THEN true 
		        ELSE false 
		    END is_exist, 
		e.id, CONCAT(e.first_name, ' ', e.last_name) as employee_name 
		FROM %s e
		WHERE 
		    e.id_card = $1 AND e.deleted = FALSE`,
		input.TableName)

	params := []interface{}{idCard}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&isExist, &result.ID, &result.Name)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) GetCountEmployeeHistory(db *sql.DB, searchByParam []in.SearchByParam) (result int, err errorModel.ErrorModel) {
	var (
		tableName       string
		additionalWhere string
	)

	tableName = fmt.Sprintf(`
		%s a
		INNER JOIN %s eb ON a.primary_key = eb.employee_id 
		INNER JOIN %s a2 ON eb.uuid_key = a2.uuid_key AND a.created_at = a2.created_at 
		LEFT JOIN "%s" u ON a.created_by = u.id `,
		AuditSystemDAO.TableName, EmployeeBenefitsDAO.TableName, AuditSystemDAO.TableName,
		UserDAO.TableName)

	for i, item := range searchByParam {
		if (searchByParam)[i].SearchKey == "name" {
			(searchByParam)[i].SearchKey = "CONCAT(u.first_name, ' ', u.last_name)"
		} else {
			searchByParam[i].SearchKey = "a." + item.SearchKey
		}
	}

	additionalWhere += fmt.Sprintf(` 
		AND a.table_name = '%s' 
		AND a2.table_name = '%s' 
		AND a."action" = %d 
		AND a2."action" = %d 
		AND a2.deleted = FALSE 
		AND 
			(
				(a.description IS NOT NULL AND a2.description IS NOT NULL) OR
				(a.description IS NULL AND a2.description IS NOT NULL) OR 
				(a.description IS NOT NULL AND a2.description IS NULL)
			) `, "employee", "employee_benefits", 2, 2)

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, tableName, searchByParam, additionalWhere, DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "a.id"},
		Deleted:   FieldStatus{FieldName: "a.deleted"},
		CreatedBy: FieldStatus{FieldName: "a.created_by", Value: int64(0)},
	})
}

func (input employeeDAO) GetListEmployeeHistory(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	query = fmt.Sprintf(`
		SELECT 
			a.id,
			a.description AS desc_employee, 
			a2.description AS desc_employee_benefits,
			CONCAT(u.first_name, ' ', u.last_name) AS editor,
			a.created_at AS edit_time
		FROM %s a
		INNER JOIN %s eb ON a.primary_key = eb.employee_id 
		INNER JOIN %s a2 ON eb.uuid_key = a2.uuid_key AND a.created_at = a2.created_at 
		LEFT JOIN "%s" u ON a.created_by = u.id `,
		AuditSystemDAO.TableName, EmployeeBenefitsDAO.TableName, AuditSystemDAO.TableName,
		UserDAO.TableName)

	for i, item := range searchByParam {
		if (searchByParam)[i].SearchKey == "name" {
			(searchByParam)[i].SearchKey = "CONCAT(u.first_name, ' ', u.last_name)"
		} else {
			searchByParam[i].SearchKey = "a." + item.SearchKey
		}
	}

	switch userParam.OrderBy {
	case "created_at", "created_at ASC", "created_at DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "a.created_at " + strSplit[1]
		} else {
			userParam.OrderBy = "a.created_at"
		}
	default:
	}

	additionalWhere += fmt.Sprintf(` 
		AND a.table_name = '%s' 
		AND a2.table_name = '%s' 
		AND a."action" = %d 
		AND a2."action" = %d 
		AND a2.deleted = FALSE 
		AND 
			(
				(a.description IS NOT NULL AND a2.description IS NOT NULL) OR
				(a.description IS NULL AND a2.description IS NOT NULL) OR 
				(a.description IS NOT NULL AND a2.description IS NULL)
			) `, "employee", "employee_benefits", 2, 2)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.EmployeeHistoryModel
			dbError := rows.Scan(
				&temp.ID, &temp.Description1, &temp.Description2,
				&temp.Editor, &temp.CreatedAt,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:      FieldStatus{FieldName: "a.id"},
			Deleted: FieldStatus{FieldName: "a.deleted"},
			CreatedBy: FieldStatus{
				FieldName: "a.created_by",
				Value:     int64(0),
			},
		})
}

func (input employeeDAO) GetListByMemberId(db *sql.DB, id int64) (result []repository.EmployeeModel, errModel errorModel.ErrorModel) {
	funcName := "GetListByMemberId"

	query := `select 
				*
			from (
				select 
					e.id, e.first_name, e.last_name,
					u.client_id, e.email, CASE
						WHEN is_valid_json(e.member)
						  THEN (e.member::jsonb->>'member_id')::jsonb
						ELSE
						  NULL
					  end as members
				from ` + EmployeeDAO.TableName + ` e 
				inner join "user" as u
					ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
					OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
				where
					e.is_have_member = true and 
					u.deleted = false and 
					e.deleted = false 
			) as ea
			where 
				ea.members ? $1`

	row, err := db.Query(query, strconv.Itoa(int(id)))
	if err != nil {
		return nil, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	for row.Next() {
		var (
			employee repository.EmployeeModel
			members  sql.NullString
		)

		if err = row.Scan(
			&employee.ID, &employee.FirstName, &employee.LastName,
			&employee.ClientId, &employee.Email, &members,
		); err != nil {
			return nil, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		}

		result = append(result, employee)
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeDAO) GetAllEmployee(db *sql.DB, joinStr string) (results []repository.EmployeeBenefitsModel, err errorModel.ErrorModel) {
	funcName := "GetAllEmployee"

	query := "SELECT e.id, e.date_join FROM employee AS e " +
		     joinStr + "  employee_benefits AS eb ON e.id = eb.employee_id " +
		     " WHERE e.deleted=FALSE ORDER BY e.id ASC"

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
			var value repository.EmployeeBenefitsModel
			rows.Scan(&value.ID, &value.JoinDate)
			benefit, _ := EmployeeBenefitsDAO.GetDetailMedicalValueForVerify(db, value.ID.Int64)
			value.CurrentAnnualLeave.Int64 = benefit.CurrentAnnualLeave.Int64
			value.LastAnnualLeave.Int64 = benefit.LastAnnualLeave.Int64
			value.CurrentMedicalValue.Float64 = benefit.CurrentMedicalValue.Float64
			value.LastMedicalValue.Float64 = benefit.LastMedicalValue.Float64
			results = append(results, value)
		}

	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorRows)
		return
	}
	return
}