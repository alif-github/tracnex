package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type companyDAO struct {
	AbstractDAO
}

var CompanyDAO = companyDAO{}.New()

func (input companyDAO) New() (output companyDAO) {
	output.TableName = "internal_company"
	output.FileName = "CompanyDAO.go"
	return
}

func (input companyDAO) InsertCompany(db *sql.Tx, inputStruct repository.CompanyModel) (lastInsertedId int64, err errorModel.ErrorModel) {
	funcName := "InsertCompany"
	query := "INSERT INTO " + input.TableName + " (" +
		" title, company_name, address, address2, " +
		" neighbourhood, hamlet, province_id, district_id, sub_district_id, " +
		" urban_village_id, postal_code_id, longitute, latitute, telephone, alternate_telephone, " +
		" fax, email, alternate_email, npwp, tax_name, tax_address, created_client, created_at, updated_client ) " +
		"VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,$15," +
		" $16, $17, $18, $19, $20, $21, $22, $23, $24 ) " +
		" RETURNING id"
	params := []interface{}{
		inputStruct.CompanyTitle.String, inputStruct.CompanyName.String, inputStruct.Address.String,
		inputStruct.Address2.String, inputStruct.Neighbourhood.String, inputStruct.Hamlet.String,
		inputStruct.ProvinceID.Int64, inputStruct.DistrictID.Int64, inputStruct.SubDistrictID.Int64,
		inputStruct.UrbanVillageID.Int64, inputStruct.PostalCodeID.Int64, inputStruct.Longitude.String,
		inputStruct.Latitude.String, inputStruct.Telephone.String, inputStruct.AlternateTelephone.String,
		inputStruct.Fax.String, inputStruct.CompanyEmail.String, inputStruct.AlternativeCompanyEmail.String,
		inputStruct.Npwp.String, inputStruct.TaxName.String, inputStruct.TaxAddress.String,
		inputStruct.CreatedClient.String, inputStruct.CreatedAt.Time, inputStruct.UpdatedClient.String,
	}
	result := db.QueryRow(query, params...)
	dbError := result.Scan(&lastInsertedId)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyDAO) GetCompanyList(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) ([]interface{}, errorModel.ErrorModel) {
	query := `SELECT 
			id, title, company_name,
			address, telephone,
            created_at, created_by,
			updated_at, updated_by
		FROM ` + input.TableName + " "

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var inputStruct repository.CompanyModel

		dbError := rows.Scan(
			&inputStruct.ID,
			&inputStruct.CompanyTitle,
			&inputStruct.CompanyName,
			&inputStruct.Address,
			&inputStruct.Telephone,
			&inputStruct.CreatedAt,
			&inputStruct.CreatedBy,
			&inputStruct.UpdatedAt,
			&inputStruct.UpdatedBy)

		return inputStruct, dbError
	}

	return GetListDataDAO.GetListData(db, query, userParam, searchBy, createdBy, mappingFunc, "")
}

func (input companyDAO) GetCountCompany(db *sql.DB, searchBy []in.SearchByParam) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT 
				COUNT(*) 
			FROM ` + input.TableName + ``

	params := []interface{}{}

	return GetListDataDAO.GetCountData(db, params, query, searchBy, "", DefaultFieldMustCheck{
		CreatedBy: FieldStatus{
			Value:     int64(0),
		},
		Deleted: FieldStatus{
			IsCheck: true,
			FieldName: "deleted",
		},
	})
}

func (input companyDAO) UpdateCompany(db *sql.Tx, inputStruct repository.CompanyModel) errorModel.ErrorModel {
	funcName := "UpdateCompany"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		" title = $1, company_name = $2, address = $3, address2 = $4, " +
		" neighbourhood = $5, hamlet = $6, province_id = $7, district_id = $8, sub_district_id = $9, " +
		" urban_village_id = $10, postal_code_id = $11, longitute = $12, latitute = $13, telephone = $14, alternate_telephone = $15, " +
		" fax = $16, email = $17, alternate_email = $18, npwp = $19, tax_name = $20, tax_address = $21, " +
		"	updated_client = $22," +
		"	updated_at = $23," +
		"	updated_by = $24 " +
		"WHERE " +
		"	id = $25 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		inputStruct.CompanyTitle.String, inputStruct.CompanyName.String, inputStruct.Address.String,
		inputStruct.Address2.String, inputStruct.Neighbourhood.String, inputStruct.Hamlet.String,
		inputStruct.ProvinceID.Int64, inputStruct.DistrictID.Int64, inputStruct.SubDistrictID.Int64,
		inputStruct.UrbanVillageID.Int64, inputStruct.PostalCodeID.Int64, inputStruct.Longitude.String,
		inputStruct.Latitude.String, inputStruct.Telephone.String, inputStruct.AlternateTelephone.String,
		inputStruct.Fax.String, inputStruct.CompanyEmail.String, inputStruct.AlternativeCompanyEmail.String,
		inputStruct.Npwp.String, inputStruct.TaxName.String, inputStruct.TaxAddress.String,
		inputStruct.UpdatedClient.String, inputStruct.UpdatedAt.Time, inputStruct.UpdatedBy.Int64, inputStruct.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input companyDAO) DeleteCompany(db *sql.Tx, inputStruct repository.CompanyModel) errorModel.ErrorModel {
	funcName := "DeleteCompany"

	query := "UPDATE " + input.TableName + " " +
		"SET" +
		"	deleted = $1," +
		"	updated_client = $2," +
		"	updated_at = $3," +
		"	updated_by = $4 " +
		"WHERE " +
		"	id = $5 AND " +
		"	deleted = false"

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	result, dbError := stmt.Exec(
		inputStruct.Deleted.Bool,
		inputStruct.UpdatedClient.String,
		inputStruct.UpdatedAt.Time,
		inputStruct.UpdatedBy.Int64,
		inputStruct.ID.Int64)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffected < 1 || rowsAffectedError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, rowsAffectedError)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input companyDAO) GetDetailCompanyForUpdate(db *sql.Tx, id int64) (inputStruct repository.CompanyModel, err errorModel.ErrorModel) {
	funcName := "GetDetailCompanyForUpdate"
	query := "SELECT id, updated_at " +
		" FROM " + input.TableName +
		" WHERE deleted = FALSE AND id = $1 "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&inputStruct.ID, &inputStruct.UpdatedAt)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyDAO) GetDetailCompany(db *sql.DB, id int64) (inputStruct repository.CompanyModel, err errorModel.ErrorModel) {
	funcName := "GetDetailCompany"
	query := "SELECT c.id, c.updated_at, c.title, c.company_name, c.address, c.address2, " +
		" c.neighbourhood, c.hamlet, c.province_id, c.district_id, c.sub_district_id, " +
		" c.urban_village_id, c.postal_code_id, c.longitute, c.latitute, c.telephone, c.alternate_telephone," +
		" c.fax, c.email, c.alternate_email, c.npwp, c.tax_name, c.tax_address," +
		" p.name, d.name, sd.name, uv.name, pc.code " +
		" FROM " + input.TableName + " AS c" +
		" LEFT JOIN province AS p ON p.id = c.province_id " +
		" LEFT JOIN district AS d ON d.id = c.district_id " +
		" LEFT JOIN sub_district AS sd ON sd.id = c.sub_district_id " +
		" LEFT JOIN urban_village AS uv ON uv.id = c.urban_village_id " +
		" LEFT JOIN postal_code AS pc ON pc.id = c.postal_code_id " +
		" WHERE c.deleted = FALSE AND c.id = $1 "

	param := []interface{}{id}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&inputStruct.ID, &inputStruct.UpdatedAt, &inputStruct.CompanyTitle,
		&inputStruct.CompanyName, &inputStruct.Address, &inputStruct.Address2, &inputStruct.Neighbourhood,
		&inputStruct.Hamlet, &inputStruct.ProvinceID, &inputStruct.DistrictID, &inputStruct.SubDistrictID,
		&inputStruct.UrbanVillageID, &inputStruct.PostalCodeID, &inputStruct.Longitude, &inputStruct.Latitude,
		&inputStruct.Telephone, &inputStruct.AlternateTelephone, &inputStruct.Fax, &inputStruct.CompanyEmail,
		&inputStruct.AlternativeCompanyEmail, &inputStruct.Npwp, &inputStruct.TaxName, &inputStruct.TaxAddress,
		&inputStruct.ProvinceName, &inputStruct.DistrictName, &inputStruct.SubDistrictName, &inputStruct.UrbanVillageName, &inputStruct.PostalCode)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyDAO) CheckCompany(db *sql.Tx, key string, fieldName string) (id int64, err errorModel.ErrorModel) {
	funcName := "CheckCompany"
	query := "SELECT " +
		"	id FROM " + input.TableName + " " +
		" WHERE LOWER(" + fieldName + ") = LOWER($1) AND deleted = FALSE LIMIT 1 "

	param := []interface{}{key}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&id)
	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyDAO) CheckCompanyByID(db *sql.Tx, id int64) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName = "CheckCompanyByID"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT CASE WHEN id > 0 THEN TRUE ELSE FALSE END is_exist 
		FROM %s 
		WHERE id = $1 AND deleted = FALSE `,
		input.TableName)

	param := []interface{}{id}
	results := db.QueryRow(query, param...)
	dbError := results.Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
