package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type employeeContractDAO struct {
	AbstractDAO
}

var EmployeeContractDAO = employeeContractDAO{}.New()

func (input employeeContractDAO) New() (output employeeContractDAO) {
	output.FileName = "EmployeeContractDAO.go"
	output.TableName = "employee_contract"
	return
}

func (input employeeContractDAO) InsertEmployeeContract(db *sql.Tx, userParam repository.EmployeeContractModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertEmployeeContract"
		query    string
	)

	query = fmt.Sprintf(`
		INSERT INTO %s 	
		(
			contract_no, information, employee_id, 
			from_date, thru_date, created_by, 
			created_client, created_at, updated_by, 
			updated_client, updated_at
		) VALUES 
		(
			$1, $2, $3, 
			$4, $5, $6, 
			$7, $8, $9, 
			$10, $11
		) 
		RETURNING id `,
		input.TableName)

	param := []interface{}{
		userParam.ContractNo.String, userParam.Information.String, userParam.EmployeeID.Int64,
		userParam.FromDate.Time, userParam.ThruDate.Time, userParam.CreatedBy.Int64,
		userParam.CreatedClient.String, userParam.CreatedAt.Time, userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
	}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractDAO) GetCountEmployeeContract(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		additionalWhere    string
		colAdditionalWhere []string
	)

	colAdditionalWhere = input.setScopeData(scopeLimit, scopeDB, false) //--- Scope Check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += fmt.Sprintf(` AND %s `, itemColAdditionalWhere)
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, additionalWhere, DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input employeeContractDAO) GetListEmployeeContract(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere string
		query           string
	)

	query = fmt.Sprintf(`
		SELECT 
		ec.id, ec.contract_no, ec.information, 
		ec.employee_id, ec.from_date, ec.thru_date, 
		uc.nt_username as created_name, ec.created_at, up.nt_username as updated_name, 
		ec.updated_at
		FROM %s ec 
		LEFT JOIN "%s" uc ON ec.created_by = uc.id AND uc.deleted = FALSE
		LEFT JOIN "%s" up ON ec.updated_by = up.id AND up.deleted = FALSE `,
		input.TableName, UserDAO.TableName, UserDAO.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, true) //--- Scope Check
	for _, itemColAdditionalWhere := range colAdditionalWhere {
		additionalWhere += fmt.Sprintf(` AND %s `, itemColAdditionalWhere)
	}

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.EmployeeContractModel
		dbError := rows.Scan(
			&temp.ID, &temp.ContractNo, &temp.Information,
			&temp.EmployeeID, &temp.FromDate, &temp.ThruDate,
			&temp.CreatedName, &temp.CreatedAt, &temp.UpdatedName,
			&temp.UpdatedAt)
		return temp, dbError
	}, additionalWhere, DefaultFieldMustCheck{
		Deleted:   FieldStatus{FieldName: "ec.deleted"},
		CreatedBy: FieldStatus{FieldName: "ec.created_by", Value: createdBy},
	})
}

func (input employeeContractDAO) ViewEmployeeContract(db *sql.DB, userParam repository.EmployeeContractModel) (result repository.EmployeeContractModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewEmployeeContract"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT 
		ec.id, ec.contract_no, ec.information, 
		ec.from_date, ec.thru_date, ec.employee_id,
		uc.nt_username as created_name, ec.created_at, up.nt_username as updated_name, 
		ec.updated_at, ec.created_by, ec.updated_by
		FROM %s ec
		LEFT JOIN "%s" uc ON ec.created_by = uc.id AND uc.deleted = FALSE
		LEFT JOIN "%s" up ON ec.updated_by = up.id AND up.deleted = FALSE
		WHERE ec.id = $1 AND ec.deleted = FALSE `,
		input.TableName, UserDAO.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.ContractNo, &result.Information,
		&result.FromDate, &result.ThruDate, &result.EmployeeID,
		&result.CreatedName, &result.CreatedAt, &result.UpdatedName,
		&result.UpdatedAt, &result.CreatedBy, &result.UpdatedBy,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractDAO) GetEmployeeContractForUpdate(db *sql.DB, userParam repository.EmployeeContractModel) (result repository.EmployeeContractModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetEmployeeContractForUpdate"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT id, updated_at, created_by 
		FROM %s WHERE id = $1 AND deleted = FALSE 
		FOR UPDATE `,
		input.TableName)

	param := []interface{}{userParam.ID.Int64}
	dbError := db.QueryRow(query, param...).Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy)
	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractDAO) DeleteEmployeeContract(tx *sql.Tx, userParam repository.EmployeeContractModel) errorModel.ErrorModel {
	var (
		funcName = "DeleteEmployeeContract"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		deleted = TRUE, updated_by = $1, updated_client = $2, 
		updated_at = $3 
		WHERE id = $4 `,
		input.TableName)

	param := []interface{}{
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
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

func (input employeeContractDAO) UpdateEmployeeContract(tx *sql.Tx, userParam repository.EmployeeContractModel) errorModel.ErrorModel {
	var (
		funcName = "UpdateEmployeeContract"
		query    string
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		contract_no = $1, information = $2, from_date = $3, 
		thru_date = $4, updated_by = $5, updated_client = $6, 
		updated_at = $7 
		WHERE id = $8 `,
		input.TableName)

	param := []interface{}{
		userParam.ContractNo.String, userParam.Information.String, userParam.FromDate.Time,
		userParam.ThruDate.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.ID.Int64,
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

func (input employeeContractDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (additionalWhere []string) {
	keyScope := []string{constanta.EmployeeDataScope}
	for _, itemKeyScope := range keyScope {
		var additionalWhereTemp string
		PrepareScopeOnDAO(scopeLimit, scopeDB, &additionalWhereTemp, 0, itemKeyScope, isView)
		if additionalWhereTemp != "" {
			additionalWhere = append(additionalWhere, additionalWhereTemp)
		}
	}

	return
}

func (input employeeContractDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		switch (*searchByParam)[i].SearchKey {
		case "contract_no":
			(*searchByParam)[i].SearchKey = "ec.contract_no"
		case "employee_id":
			(*searchByParam)[i].SearchKey = "ec.employee_id"
		default:
			(*searchByParam)[i].SearchKey = "ec." + (*searchByParam)[i].SearchKey
		}
	}

	switch userParam.OrderBy {
	case "contract_no", "contract_no ASC", "contract_no DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "ec.contract_no " + strSplit[1]
		} else {
			userParam.OrderBy = "ec.contract_no"
		}
	case "from_date", "from_date ASC", "from_date DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "ec.from_date " + strSplit[1]
		} else {
			userParam.OrderBy = "ec.from_date"
		}
	case "thru_date", "thru_date ASC", "thru_date DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "ec.thru_date " + strSplit[1]
		} else {
			userParam.OrderBy = "ec.thru_date"
		}
	case "information", "information ASC", "information DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "ec.information " + strSplit[1]
		} else {
			userParam.OrderBy = "ec.information"
		}
	case "created_at", "created_at ASC", "created_at DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "ec.created_at " + strSplit[1]
		} else {
			userParam.OrderBy = "ec.created_at DESC"
		}
	}
}
