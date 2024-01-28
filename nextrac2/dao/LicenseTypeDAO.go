package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type licenseTypeDAO struct {
	AbstractDAO
}

var LicenseTypeDAO = licenseTypeDAO{}.New()

func (input licenseTypeDAO) New() (output licenseTypeDAO) {
	output.FileName = "LicenseTypeDAO.go"
	output.TableName = "license_type"
	return
}

func (input licenseTypeDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "lt.id"},
		Deleted:   FieldStatus{FieldName: "lt.deleted"},
		CreatedBy: FieldStatus{FieldName: "lt.created_by", Value: createdBy},
	}
}

func (input licenseTypeDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "lt." + (*searchByParam)[i].SearchKey
	}

	switch userParam.OrderBy {
	case "updated_name","updated_name ASC","updated_name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "u.nt_username " + strSplit[1]
		}else {
			userParam.OrderBy = "u.nt_username"
		}
		break
	default :
		userParam.OrderBy = "lt." + userParam.OrderBy
		break
	}
}

func (input licenseTypeDAO) GetLicenseTypeForUpdate(db *sql.Tx, userParam repository.LicenseTypeModel) (result repository.LicenseTypeModel, err errorModel.ErrorModel) {
	funcName := "GetLicenseTypeForUpdate"

	query := fmt.Sprintf(`SELECT lt.id, lt.updated_at, lt.created_by, 
									CASE WHEN 
										(SELECT COUNT(p.id) FROM %s p 
											WHERE p.license_type_id = lt.id AND p.deleted = false) > 0 
										OR 
										(SELECT COUNT(lc.id) FROM %s lc 
											WHERE lc.license_type_id = lt.id AND lc.deleted = false) > 0 
									THEN TRUE ELSE FALSE END is_used, 
									lt.license_type_name 
								FROM %s lt 
								WHERE lt.id = $1 AND lt.deleted = FALSE `,
								ProductDAO.TableName,
								LicenseConfigDAO.TableName,
								input.TableName)

	param := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE "

	dbError := db.QueryRow(query, param...).Scan(
		&result.ID, &result.UpdatedAt, &result.CreatedBy, &result.IsUsed, &result.LicenseTypeName,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseTypeDAO) ViewLicenseType(db *sql.DB, userParam repository.LicenseTypeModel) (result repository.LicenseTypeModel, err errorModel.ErrorModel) {
	funcName := "ViewLicenseType"

	query := fmt.Sprintf(
		`SELECT lt.id, lt.license_type_name, lt.license_type_desc,
			lt.created_at, lt.updated_at, lt.updated_by, u.nt_username,
			lt.created_by
		FROM %s lt
		LEFT JOIN "%s" AS u ON lt.updated_by = u.id
		WHERE lt.id = $1 AND lt.deleted = FALSE `,
		input.TableName, UserDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(
		&result.ID, &result.LicenseTypeName, &result.LicenseTypeDesc,
		&result.CreatedAt, &result.UpdatedAt, &result.UpdatedBy,
		&result.UpdatedName, &result.CreatedBy,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input licenseTypeDAO) DeleteLicenseType(db *sql.Tx, userParam repository.LicenseTypeModel) (err errorModel.ErrorModel) {
	funcName := "DeleteLicenseType"

	query := fmt.Sprintf(`UPDATE %s SET 
									deleted = $1, 
									license_type_name = $2, 
									updated_by = $3, 
									updated_at = $4, 
									updated_client = $5 
								WHERE id = $6 `,
								input.TableName)

	param := []interface{}{
		true,
		userParam.LicenseTypeName.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
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

func (input licenseTypeDAO) UpdateLicenseType(db *sql.Tx, userParam repository.LicenseTypeModel) (err errorModel.ErrorModel) {
	funcName := "UpdateLicenseType"

	query := fmt.Sprintf(
		`UPDATE %s
		SET
			license_type_name = $1, license_type_desc = $2,
			updated_by = $3, updated_at = $4,
			updated_client = $5
		WHERE
			id = $6 `, input.TableName)

	params := []interface{}{
		userParam.LicenseTypeName.String, userParam.LicenseTypeDesc.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(params...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return 
}

func (input licenseTypeDAO) GetCountLicenseType(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input licenseTypeDAO) GetListLicenseType(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	query := fmt.Sprintf(
		`SELECT 
			lt.id as id, lt.license_type_name as license_type_name, 
			lt.license_type_desc as license_type_desc,
			lt.created_at as created_at, lt.updated_at as updated_at, 
			lt.updated_by, u.nt_username
		FROM %s lt
		LEFT JOIN "%s" AS u ON lt.updated_by = u.id `,
		input.TableName, UserDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.LicenseTypeModel
			dbError := rows.Scan(
				&temp.ID, &temp.LicenseTypeName, &temp.LicenseTypeDesc,
				&temp.CreatedAt, &temp.UpdatedAt, &temp.UpdatedBy, &temp.UpdatedName,
			)
			return temp, dbError
		}, " ", input.getDefaultMustCheck(createdBy))

}

func (input licenseTypeDAO) InsertLicenseType(db *sql.Tx, userParam repository.LicenseTypeModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertLicenseType"

	query := fmt.Sprintf(
		`INSERT INTO %s
		(
			license_type_name, license_type_desc,
			created_by, created_at, created_client,
			updated_by, updated_at, updated_client
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.LicenseTypeName.String,
		userParam.LicenseTypeDesc.String,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input licenseTypeDAO) CheckLicenseTypeIsExist(db *sql.DB, userParam repository.LicenseTypeModel) (result repository.LicenseTypeModel, err errorModel.ErrorModel) {
	funcName := "CheckLicenseTypeIsExist"

	query := fmt.Sprintf(
		`SELECT 
			id 
		FROM  %s  
		WHERE 
			id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	dbError := db.QueryRow(query, param...).Scan(&result.ID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}