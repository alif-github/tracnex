package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type licenseVariantDAO struct {
	AbstractDAO
}

var LicenseVariantDAO = licenseVariantDAO{}.New()

func (input licenseVariantDAO) New() (output licenseVariantDAO) {
	output.FileName = "LicenseVariantDAO.go"
	output.TableName = "license_variant"
	return
}

func (input licenseVariantDAO) InsertLicenseVariant(db *sql.Tx, userParam repository.LicenseVariantModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertLicenseVariant"
	query := fmt.Sprintf(
		`INSERT 
		INTO %s
		(
			license_variant_name, created_by, created_client, 
			created_at, updated_by, updated_client, 
			updated_at
		) 
		VALUES 
		(	
			$1, $2, $3, $4, $5, $6, $7
		) 
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.LicenseVariantName.String, userParam.CreatedBy.Int64, userParam.CreatedClient.String,
		userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
	}

	result := db.QueryRow(query, params...)
	errorS := result.Scan(&id)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input licenseVariantDAO) GetListLicenseVariant(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {

	subQuery := "SELECT id, license_variant_name, created_at, " +
		"updated_by, updated_at " +
		"FROM " + input.TableName + " "

	query := "SELECT *, (SELECT nt_username as updated_name FROM \"user\" WHERE id = a.updated_by) FROM (" + subQuery + " "

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.LicenseVariantListModel
			errors := rows.Scan(
				&temp.ID, &temp.LicenseVariantName, &temp.CreatedAt,
				&temp.UpdatedBy, &temp.UpdatedAt, &temp.UpdatedName)
			return temp, errors
		}, " ) a ", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input licenseVariantDAO) GetCountLicenseVariant(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input licenseVariantDAO) UpdateLicenseVariant(db *sql.Tx, userParam repository.LicenseVariantModel) (err errorModel.ErrorModel) {
	funcName := "UpdateLicenseVariant"

	query := fmt.Sprintf(
		`UPDATE %s  
		SET 
			license_variant_name = $1, updated_by = $2, updated_client = $3, 
			updated_at = $4 
		WHERE 
			id = $5 `, input.TableName)

	param := []interface{}{
		userParam.LicenseVariantName.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.ID.Int64,
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	return
}

func (input licenseVariantDAO) ViewLicenseVariant(db *sql.DB, userParam repository.LicenseVariantModel) (result repository.LicenseVariantModel, err errorModel.ErrorModel) {
	funcName := "ViewLicenseVariant"

	subQuery := fmt.Sprintf(
		`SELECT 
			id, license_variant_name, created_at,
			updated_at, updated_by, created_by 
		FROM %s 
		WHERE 
			id = $1 AND deleted = FALSE `, input.TableName)

	params := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		subQuery += " AND created_by = $2 "
		params = append(params, userParam.CreatedBy.Int64)
	}

	query := fmt.Sprintf(
		`SELECT 
			*, 
			(SELECT 
				nt_username as updated_name 
			FROM "%s" 
			WHERE 
				id = a.updated_by) 
		FROM ( %s ) a `, UserDAO.TableName, subQuery)

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.LicenseVariantName, &result.CreatedAt,
		&result.UpdatedAt, &result.UpdatedBy, &result.CreatedBy,
		&result.UpdatedName)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantDAO) DeleteLicenseVariant(db *sql.Tx, userParam repository.LicenseVariantModel) (err errorModel.ErrorModel) {
	funcName := "DeleteLicenseVariant"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			deleted = TRUE,
			license_variant_name = $1,
			updated_by = $2,
			updated_client = $3,
			updated_at = $4
		WHERE
			id = $5 `, input.TableName)

	param := []interface{}{
		userParam.LicenseVariantName.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.ID.Int64,
	}

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

func (input licenseVariantDAO) GetLicenseVariantForUpdate(db *sql.DB, licenseVariantModel repository.LicenseVariantModel) (result repository.LicenseVariantModel, err errorModel.ErrorModel) {
	funcName := "GetLicenseVariantForUpdate"

	query := fmt.Sprintf(
		`SELECT 
			lv.id, lv.updated_at, lv.created_by,
			CASE WHEN
				(SELECT COUNT(id) FROM %s WHERE license_variant_id = lv.id AND deleted = FALSE) > 0
					OR
				(SELECT COUNT(id) FROM %s WHERE license_variant_id = lv.id AND deleted = FALSE) > 0
			THEN TRUE ELSE FALSE END isUsed, lv.license_variant_name
		FROM %s as lv
		WHERE
			lv.id = $1 AND
			lv.deleted = FALSE `,
		ProductDAO.TableName, LicenseConfigDAO.TableName, input.TableName)

	params := []interface{}{licenseVariantModel.ID.Int64}

	if licenseVariantModel.CreatedBy.Int64 > 0 {
		query += " AND lv.created_by = $2 "
		params = append(params, licenseVariantModel.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy, &result.IsUsed, &result.LicenseVariantName)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantDAO) CheckIsExistLicenseVariant(db *sql.DB, licenseVariantModel repository.LicenseVariantModel) (result repository.LicenseVariantModel, err errorModel.ErrorModel) {
	funcName := "CheckIsExistLicenseVariant"

	query := fmt.Sprintf(
		`SELECT 
			id
		FROM %s
		WHERE 
			id = $1 AND
			deleted = FALSE `, input.TableName)

	params := []interface{}{licenseVariantModel.ID.Int64}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.ID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantDAO) InsertDataForTesting(db *sql.DB, userParam []repository.LicenseVariantModel) (dbError error) {
	index := 1
	lengthParam := 8
	var params []interface{}

	query := fmt.Sprintf(
		`INSERT 
		INTO %s
		(
			id, license_variant_name, created_by, created_client, 
			created_at, updated_by, updated_client, 
			updated_at
		) 
		VALUES `, input.TableName)

	query += CreateDollarParamInMultipleRowsDAO(len(userParam), lengthParam, index, " id ")
	for _, model := range userParam {
		params = append(params, model.ID.Int64,model.LicenseVariantName.String, model.CreatedBy.Int64, model.CreatedClient.String,
			model.CreatedAt.Time, model.UpdatedBy.Int64, model.UpdatedClient.String,
			model.UpdatedAt.Time)
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		return
	}

	return
}