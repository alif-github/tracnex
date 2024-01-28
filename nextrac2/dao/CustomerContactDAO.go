package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type customerContactDAO struct {
	AbstractDAO
}

var CustomerContactDAO = customerContactDAO{}.New()

func (input customerContactDAO) New() (output customerContactDAO) {
	output.FileName = "CustomerContactDAO.go"
	output.TableName = "customer_contact"
	return
}

func (input customerContactDAO) GetCustomerContactForUpdate(db *sql.Tx, userParam repository.CustomerContactModel) (result repository.CustomerContactModel, err errorModel.ErrorModel) {
	funcName := "GetCustomerContactForUpdate"
	index := 2
	query := fmt.Sprintf(
		`SELECT 
			cc.id, cc.updated_at, cc.updated_by, 
			cc.nik, cc.mdb_person_title_id, cc.first_name,
			cc.last_name, cc.sex, cc.address,
			cc.hamlet, cc.neighbourhood, cc.province_id, 
			cc.district_id
		FROM %s cc 
		WHERE 
			cc.customer_id = $1 AND	
			cc.deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.CustomerID.Int64}

	if userParam.ID.Int64 > 0 {
		query += fmt.Sprintf(" AND cc.id = $%d ", index)
		param = append(param, userParam.ID.Int64)
		index++
	}

	if !util.IsStringEmpty(userParam.Nik.String) {
		query += fmt.Sprintf(" AND cc.nik = $%d ", index)
		param = append(param, userParam.Nik.String)
		index++
	}

	if userParam.CreatedBy.Int64 > 0 {
		query += fmt.Sprintf(" AND cc.created_by = $%d ", index)
		param = append(param, userParam.CreatedBy.Int64)
		index++
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, param...)

	errs := results.Scan(
		&result.ID, &result.UpdatedAt, &result.UpdatedBy,
		&result.Nik, &result.MdbPersonTitle, &result.FirstName,
		&result.LastName, &result.Sex, &result.Address,
		&result.Hamlet, &result.Neighbourhood, &result.ProvinceID,
		&result.DistrictID,
	)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerContactDAO) UpdateCustomerContact(db *sql.Tx, userParam repository.CustomerContactModel) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomerContact"
	var param []interface{}
	query := " UPDATE " + input.TableName + " " +
		" SET " +
		" 	customer_id = $1, " +
		" 	mdb_person_title_id = $2, " +
		" 	province_id = $3, " +
		"	district_id = $4, " +
		"	mdb_position_id = $5, " +
		" 	position_name = $6, " +
		"	sex = $7, " +
		"	address = $8, " +
		"	hamlet = $9, " +
		"	neighbourhood = $10, " +
		" 	nik = $11, " +
		"	person_title = $12, " +
		"	phone = $13, " +
		"	email = $14, " +
		"	first_name = $15, " +
		" 	last_name = $16, " +
		"	status = $17, " +
		"	updated_by = $18, " +
		"	updated_client = $19, " +
		"	updated_at = $20 " +
		" WHERE " +
		"	id = $21 "

	param = append(param, userParam.CustomerID.Int64)

	if userParam.MdbPersonTitle.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.MdbPersonTitle.Int64)
	}

	if userParam.ProvinceID.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.ProvinceID.Int64)
	}

	if userParam.DistrictID.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.DistrictID.Int64)
	}

	if userParam.MdbPositionID.Int64 < 1 {
		param = append(param, nil)
	} else {
		param = append(param, userParam.MdbPositionID.Int64)
	}

	param = append(param,
		userParam.PositionName.String, userParam.Sex.String,
		userParam.Address.String, userParam.Hamlet.String,
		userParam.Neighbourhood.String, userParam.Nik.String,
		userParam.PersonTitle.String, userParam.Phone.String,
		userParam.Email.String, userParam.FirstName.String,
		userParam.LastName.String, userParam.Status.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.ID.Int64,
	)

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}
	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerContactDAO) DeleteCustomerContact(db *sql.Tx, userParam repository.CustomerContactModel) (err errorModel.ErrorModel) {
	funcName := "DeleteCustomerContact"

	query := fmt.Sprintf(
		`UPDATE  %s 
		SET 
			deleted = $1, updated_by = $2, updated_at = $3, 
			updated_client = $4, nik = $5
		WHERE 
			id = $6 `, input.TableName)
	param := []interface{}{
		true,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.Nik.String,
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

func (input customerContactDAO) GetCustomerContactByCustomerID(db *sql.DB, userParam repository.CustomerContactModel) (result []repository.CustomerContactModel, err errorModel.ErrorModel) {
	funcName := "GetCustomerContactByCustomerID"
	query := fmt.Sprintf(`SELECT 
		cc.id, cc.customer_id, cc.mdb_person_profile_id, 
		cc.nik, cc.mdb_person_title_id, cc.person_title, 
		cc.first_name, cc.last_name, cc.sex, 
		cc.address, cc.address_2, cc.address_3, cc.hamlet, cc.neighbourhood, 
		cc.province_id, cc.district_id, cc.phone, 
		cc.email, cc.mdb_position_id, cc.position_name, 
		cc.status, cc.created_by, cc.created_at, 
		cc.updated_by, cc.updated_at, p.name province_name, 
		d.name district_name, uc.nt_username, ud.nt_username 
	FROM  %s cc 
	LEFT JOIN %s p ON cc.province_id = p.id 
	LEFT JOIN %s d ON cc.district_id = d.id 
	LEFT JOIN "%s" uc ON uc.id = cc.created_by
	LEFT JOIN "%s" ud ON ud.id = cc.updated_by
	WHERE 
		cc.customer_id = $1 AND cc.deleted = FALSE `,
		input.TableName, ProvinceDAO.TableName, DistrictDAO.TableName,
		UserDAO.TableName, UserDAO.TableName)

	params := []interface{}{userParam.CustomerID.Int64}

	rows, errorS := db.Query(query, params...)

	if errorS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			var temp repository.CustomerContactModel
			errorS = rows.Scan(
				&temp.ID, &temp.CustomerID, &temp.MdbPersonProfileID,
				&temp.Nik, &temp.MdbPersonTitle, &temp.PersonTitle,
				&temp.FirstName, &temp.LastName,
				&temp.Sex, &temp.Address, &temp.Address2, &temp.Address3, &temp.Hamlet,
				&temp.Neighbourhood, &temp.ProvinceID, &temp.DistrictID,
				&temp.Phone, &temp.Email, &temp.MdbPositionID,
				&temp.PositionName, &temp.Status, &temp.CreatedBy,
				&temp.CreatedAt, &temp.UpdatedBy, &temp.UpdatedAt,
				&temp.ProvinceName, &temp.DistrictName, &temp.CreatedName, &temp.UpdatedName)

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

func (input customerContactDAO) InsertBulkCustomerContact(db *sql.Tx, userParam []repository.CustomerContactModel) (id []int64, err errorModel.ErrorModel) {
	funcName := "InsertBulkCustomerContact"
	index := 1
	lenParam := 26
	var params []interface{}
	var tempQuery string

	query := fmt.Sprintf(
		`INSERT INTO %s
		( 
			customer_id, mdb_person_profile_id, nik, 
			mdb_person_title_id, person_title, first_name, 
			last_name, sex, address, 
			address_2, address_3, hamlet, 
			neighbourhood, province_id, district_id, 
			phone, email, mdb_position_id, 
			position_name, status, created_by, 
			created_client, created_at, updated_by, 
			updated_client, updated_at 
		) 
		VALUES `, input.TableName)

	for i := 0; i < len(userParam); i++ {
		query += " ( "

		HandleOptionalParam([]interface{}{
			userParam[i].CustomerID.Int64, userParam[i].MdbPersonProfileID.Int64, userParam[i].Nik.String,
			userParam[i].MdbPersonTitle.Int64, userParam[i].PersonTitle.String, userParam[i].FirstName.String,
			userParam[i].LastName.String, userParam[i].Sex.String, userParam[i].Address.String,
			userParam[i].Address2.String, userParam[i].Address3.String, userParam[i].Hamlet.String,
			userParam[i].Neighbourhood.String, userParam[i].ProvinceID.Int64, userParam[i].DistrictID.Int64,
			userParam[i].Phone.String, userParam[i].Email.String, userParam[i].MdbPositionID.Int64,
			userParam[i].PositionName.String, userParam[i].Status.String, userParam[i].CreatedBy.Int64,
			userParam[i].CreatedClient.String, userParam[i].CreatedAt.Time, userParam[i].UpdatedBy.Int64,
			userParam[i].UpdatedClient.String, userParam[i].UpdatedAt.Time,
		}, &params)

		tempQuery, index = ListRangeToInQueryWithStartIndex(lenParam, index)
		query += tempQuery
		query += " ) "

		if i < len(userParam)-1 {
			query += ", "
		}
	}

	query += " RETURNING id "

	rows, errorS := db.Query(query, params...)
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
