package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

type customerListDAO struct {
	AbstractDAO
}

var CustomerListDAO = customerListDAO{}.New()

func (input customerListDAO) New() (output customerListDAO) {
	output.FileName = "CustomerListDAO.go"
	output.TableName = "customer_list_registration"
	return
}

func (input customerListDAO) GetCustomerForUpdate(db *sql.DB, userParam repository.CustomerListModel) (result repository.CustomerListModel, err errorModel.ErrorModel) {
	query :=
		"select id " +
		"from "+ input.TableName +" " +
			"where " +
			"company_id = $1 and branch_id = $2 and product = $3"

	param := []interface{}{userParam.CompanyID.String, userParam.BranchID.String, userParam.Product.String}

	errorS := db.QueryRow(query, param...).Scan(&result.ID)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) TruncateTableCustomer(db *sql.Tx) (err errorModel.ErrorModel) {
	funcName := "TruncateTableCustomer"

	query := "TRUNCATE "+ input.TableName +" RESTART IDENTITY"

	result, dbError := db.Exec(query)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	_, rowsAffectedError := result.RowsAffected()
	if rowsAffectedError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) SetValSequenceCustomerListTable(db *sql.Tx) (err errorModel.ErrorModel) {
	funcName := "SetValSequenceCustomerListTable"
	fkName := "customer_list_registration_pkey_seq"

	query := "ALTER SEQUENCE "+ fkName +" RESTART WITH 1"

	result, dbError := db.Exec(query)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	_, rowsAffectedError := result.RowsAffected()
	if rowsAffectedError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) CountRowCustomerList(db *sql.Tx) (count int64, err errorModel.ErrorModel) {
	funcName := "CountRowCustomerList"

	query := "SELECT COUNT(id) FROM "+ input.TableName +""

	dbError := db.QueryRow(query).Scan(&count)
	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) UpdateCustomer(db *sql.DB, userParam repository.CustomerListModel) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomer"

	query :=
		"UPDATE "+ input.TableName +" " +
			"SET " +
			"company_id = $1, branch_id = $2, company_name = $3, " +
			"city = $4, implementer = $5, implementation = $6, " +
			"product = $7, version = $8, license_type = $9, " +
			"user_amount = $10, exp_date = $11, updated_by = $12, " +
			"updated_at = $13, updated_client = $14 " +
			"WHERE " +
			"id = $12 AND deleted = FALSE"

	var param []interface{}
	param = append(param, userParam.CompanyID.String, userParam.BranchID.String, userParam.CompanyName.String)

	if userParam.City.String != "" {
		param = append(param, userParam.City.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Implementer.String != "" {
		param = append(param, userParam.Implementer.String)
	} else {
		param = append(param, nil)
	}

	if !userParam.Implementation.Time.IsZero() {
		param = append(param, userParam.Implementation.Time)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.Product.String)

	if userParam.Version.String != "" {
		param = append(param, userParam.Version.String)
	} else {
		param = append(param, nil)
	}

	if userParam.LicenseType.String != "" {
		param = append(param, userParam.LicenseType.String)
	} else {
		param = append(param, nil)
	}

	if userParam.UserAmount.Int64 != 0 {
		param = append(param, userParam.UserAmount.Int64)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.ExpDate.Time, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String)

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	result, dbError := stmt.Exec(param...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) CheckCustomer(tx *sql.Tx, userParam []repository.CustomerListModel) (result []repository.CustomerListModel, err errorModel.ErrorModel) {
	funcName := "CheckCustomer"
	number := 1
	lengthUserParam := len(userParam)

	query := "select company_id, branch_id, company_name, exp_date from "+ input.TableName +" where id IN ("

	for i := 1; i <= lengthUserParam; i++ {

		query += "(select id from "+ input.TableName +" where company_id = $"+ strconv.Itoa(number) +" " +
			"AND branch_id = $"+ strconv.Itoa(number + 1) +" " +
			"AND (product like 'ND6%' " +
			"OR product like 'NF%') " +
			"order by id limit 1)"

		if lengthUserParam - i != 0 {
			query += ","
		} else {
			query += ")"
		}
		number += 2
	}
	var param []interface{}
	for i := 0; i < lengthUserParam; i++ {
		param = append(param, userParam[i].CompanyID.String, userParam[i].BranchID.String)
	}

	rows, errorS := tx.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if rows != nil {
		defer func() {
			errorS = rows.Close()
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
		}()
		for rows.Next() {
			var temp repository.CustomerListModel
			errorS = rows.Scan(&temp.CompanyID, &temp.BranchID, &temp.CompanyName, &temp.ExpDate)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			result = append(result, temp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) InsertMultipleBranchCustomer(tx *sql.Tx, userParam []repository.CustomerListModel) (id []int64, err errorModel.ErrorModel) {
	funcName := "InsertMultipleBranchCustomer"
	parameterCustomer := 17
	jVar := 1

	query := "INSERT INTO "+ input.TableName +" " +
		"(company_id, branch_id, company_name, " +
		"city, implementer, implementation, " +
		"product, version, license_type, " +
		"user_amount, exp_date, created_by, " +
		"created_at, created_client, updated_by, " +
		"updated_at, updated_client) VALUES "

	for i := 1; i <= len(userParam); i++ {
		query += "("

		for j := jVar; j <= parameterCustomer; j++ {
			query += " $"+ strconv.Itoa(j) +""
			if parameterCustomer - j != 0 {
				query += ","
			} else {
				query += ")"
			}
		}

		if len(userParam) - i != 0 {
			query += ","
		} else {
			query += " returning id"
		}

		jVar += 17
		parameterCustomer += 17
	}
	var param []interface{}
	for i := 0; i < len(userParam); i++ {
		param = append(param, userParam[i].CompanyID.String, userParam[i].BranchID.String)

		if userParam[i].CompanyName.String != "" {
			param = append(param, userParam[i].CompanyName.String)
		} else {
			param = append(param, nil)
		}

		if userParam[i].City.String != "" {
			param = append(param, userParam[i].City.String)
		} else {
			param = append(param, nil)
		}

		if userParam[i].Implementer.String != "" {
			param = append(param, userParam[i].Implementer.String)
		} else {
			param = append(param, nil)
		}

		if !userParam[i].Implementation.Time.IsZero() {
			param = append(param, userParam[i].Implementation.Time)
		} else {
			param = append(param, nil)
		}

		if userParam[i].Product.String != "" {
			param = append(param, userParam[i].Product.String)
		} else {
			param = append(param, nil)
		}

		if userParam[i].Version.String != "" {
			param = append(param, userParam[i].Version.String)
		} else {
			param = append(param, nil)
		}

		if userParam[i].LicenseType.String != "" {
			param = append(param, userParam[i].LicenseType.String)
		} else {
			param = append(param, nil)
		}

		if userParam[i].UserAmount.Int64 != 0 {
			param = append(param, userParam[i].UserAmount.Int64)
		} else {
			param = append(param, nil)
		}

		if !userParam[i].ExpDate.Time.IsZero() {
			param = append(param, userParam[i].ExpDate.Time)
		} else {
			param = append(param, nil)
		}

		param = append(param,
			userParam[i].CreatedBy.Int64, userParam[i].CreatedAt.Time, userParam[i].CreatedClient.String,
			userParam[i].UpdatedBy.Int64, userParam[i].UpdatedAt.Time, userParam[i].UpdatedClient.String)
	}

	rows, errorS := tx.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if rows != nil {
		defer func() {
			errorS = rows.Close()
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
		}()
		for rows.Next() {
			var idTemp int64
			errorS = rows.Scan(&idTemp)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			id = append(id, idTemp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) InsertCustomerByImport(tx *sql.Tx, userParam repository.CustomerListModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertCustomerByImport"

	query :=
		"INSERT INTO "+ input.TableName +" " +
			"(company_id, branch_id, company_name, " + //1, 2, 3
			"city, implementer, implementation, " + //4, 5, 6
			"product, version, license_type, " + //7, 8, 9
			"user_amount, exp_date, created_by, " + //10, 11, 12
			"created_client, created_at, updated_by, " + //13, 14, 15
			"updated_client, updated_at) " + //16, 17
			"VALUES " +
			"($1, $2, $3, " +
			"$4, $5, $6, " +
			"$7, $8, $9, " +
			"$10, $11, $12, " +
			"$13, $14, $15, " +
			"$16, $17) Returning id"

	var param []interface{}
	param = append(param, userParam.CompanyID.String, userParam.BranchID.String, userParam.CompanyName.String) //1, 2, 3

	if userParam.City.String != "" {
		param = append(param, userParam.City.String) //4
	} else {
		param = append(param, nil)
	}

	if userParam.Implementer.String != "" {
		param = append(param, userParam.Implementer.String) //5
	} else {
		param = append(param, nil)
	}

	if !userParam.Implementation.Time.IsZero() {
		param = append(param, userParam.Implementation.Time) //6
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.Product.String) //7

	if userParam.Version.String != "" {
		param = append(param, userParam.Version.String) //8
	} else {
		param = append(param, nil)
	}

	if userParam.LicenseType.String != "" {
		param = append(param, userParam.LicenseType.String) //9
	} else {
		param = append(param, nil)
	}

	if userParam.UserAmount.Int64 != 0 {
		param = append(param, userParam.UserAmount.Int64) //10
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.ExpDate.Time, userParam.CreatedBy.Int64, userParam.CreatedClient.String, //11, 12, 13
		userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, //14, 15, 16
		userParam.UpdatedAt.Time) //17

	errorS := tx.QueryRow(query, param...).Scan(&id)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) CheckCustomerByProductName(tx *sql.Tx, userParam repository.CustomerListModel) (result []repository.CustomerListModel, err errorModel.ErrorModel) {
	funcName := "CheckCustomerByProductName"

	query := "select company_id, branch_id, company_name, exp_date " +
		"from "+ input.TableName +" " +
		"where " +
		"company_id = $1 AND " +
		"branch_id = $2 AND " +
		"(product like 'ND6%' OR " +
		"product like 'NF6%')"

	var param []interface{}
	param = append(param, userParam.CompanyID.String, userParam.BranchID.String)

	rows, errorS := tx.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	if rows != nil {
		defer func() {
			errorS = rows.Close()
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
		}()
		for rows.Next() {
			var temp repository.CustomerListModel
			errorS = rows.Scan(&temp.CompanyID, &temp.BranchID, &temp.CompanyName, &temp.ExpDate)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			result = append(result, temp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) GetListCustomerList(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := "SELECT id, company_id, branch_id, " +
		"company_name, product, user_amount, " +
		"exp_date, updated_at " +
		"FROM " +
		""+ input.TableName +" "

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CustomerListModel
			errors := rows.Scan(
				&temp.ID,
				&temp.CompanyID,
				&temp.BranchID,
				&temp.CompanyName,
				&temp.Product,
				&temp.UserAmount,
				&temp.ExpDate,
				&temp.UpdatedAt)
			return temp, errors
		}, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input customerListDAO) ViewCustomer(db *sql.DB, userParam repository.CustomerListModel) (result repository.CustomerListModel, err errorModel.ErrorModel) {
	funcName := "ViewCustomer"

	query := "SELECT " +
		"company_id, branch_id, company_name, " +
		"city, implementer, implementation, " +
		"product, version, license_type, " +
		"user_amount, exp_date, created_by, " +
		"created_at, created_client, updated_by, " +
		"updated_at, updated_client, id " +
		"FROM "+ input.TableName +" " +
		"WHERE id = $1 AND " +
		"deleted = FALSE "

	params := []interface{}{
		userParam.ID.Int64,
	}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2"
		params = append(params, userParam.CreatedBy.Int64)
	}

	results := db.QueryRow(query, params...)
	errorDB := results.Scan(
		&result.CompanyID, &result.BranchID, &result.CompanyName,
		&result.City, &result.Implementer, &result.Implementation,
		&result.Product, &result.Version, &result.LicenseType,
		&result.UserAmount, &result.ExpDate, &result.CreatedBy,
		&result.CreatedAt, &result.CreatedClient, &result.UpdatedBy,
		&result.UpdatedAt, &result.UpdatedClient, &result.ID)

	if errorDB != nil && errorDB.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListDAO) GetCountCustomer(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}