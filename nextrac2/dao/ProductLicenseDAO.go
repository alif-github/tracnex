package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
)

type productLicenseDAO struct {
	AbstractDAO
}

var ProductLicenseDAO = productLicenseDAO{}.New()

func (input productLicenseDAO) New() (output productLicenseDAO) {
	output.FileName = "ProductLicenseDAO.go"
	output.TableName = "product_license"
	return
}

func (input productLicenseDAO) GetCountForCheckExpiration(db *sql.DB, searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (result int, err errorModel.ErrorModel) {
	var (
		funcName = "GetCountForCheckExpiration"
		query    string
	)

	query = fmt.Sprintf(`SELECT COUNT(pl.id)
		FROM %s pl
		INNER JOIN %s lc
		ON pl.license_config_id = lc.id
		WHERE 
		pl.deleted = FALSE AND pl.license_status = 1 AND lc.product_valid_thru < now()::date `,
		input.TableName, LicenseConfigDAO.TableName)

	row := db.QueryRow(query)
	dbError := row.Scan(&result)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseDAO) GetCountProductLicense(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64,
	scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, listScope []string) (int, errorModel.ErrorModel) {
	var (
		query, additionalWhere string
		param                  []interface{}
	)
	query = fmt.Sprintf(` pl 
		JOIN  %s lc ON pl.license_config_id = lc.id 
		JOIN  %s c ON c.id = lc.customer_id 
		JOIN  %s pr ON pr.id = lc.product_id 
		JOIN  %s lv ON lv.id = lc.license_variant_id 
		JOIN  %s lt ON lt.id = lc.license_type_id `,
		LicenseConfigDAO.TableName, CustomerDAO.TableName,
		ProductDAO.TableName, LicenseVariantDAO.TableName, LicenseTypeDAO.TableName)

	additionalWhere, param = ScopeToAddedQueryView(scopeLimit, scopeDB, 1, listScope)

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, param, input.TableName+query,
		searchByParam, additionalWhere,
		input.getCustomerDefaultMustCheck(createdBy))
}

func (input productLicenseDAO) setSearchByProductLicense(searchBy *[]in.SearchByParam) {
	temp := *searchBy

	for index := range temp {
		switch temp[index].SearchKey {
		case "customer_id":
			temp[index].SearchKey = "customer.id"
		case "customer_name":
			temp[index].SearchKey = "customer." + temp[index].SearchKey
		}
	}
}

func (input productLicenseDAO) setCreatedByProductLicense(createdBy int64, searchBy *[]in.SearchByParam) {
	if createdBy > 0 {
		*searchBy = append(*searchBy, in.SearchByParam{
			SearchKey:      "product_license.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}
}

func (input productLicenseDAO) setGetListJoinProductLicense(getListData *getListJoinDataDAO) {
	getListData.InnerJoin(LicenseConfigDAO.TableName, "product_license.license_config_id", "license_configuration.id")
	getListData.InnerJoin(CustomerDAO.TableName, "customer.id", "license_configuration.customer_id")
	getListData.InnerJoin(ProductDAO.TableName, "product.id", "license_configuration.product_id")
	getListData.InnerJoin(LicenseVariantDAO.TableName, "license_variant.id", "license_configuration.license_variant_id")
	getListData.InnerJoin(LicenseTypeDAO.TableName, "license_type.id", "license_configuration.license_type_id")

}

func (input productLicenseDAO) InsertProductLicense(db *sql.Tx, userParam repository.ProductLicenseModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertProductLicense"

	query := fmt.Sprintf(`INSERT INTO %s(
										license_config_id, product_key, product_encrypt, client_id, 
										client_secret, hwid, activation_date, license_status, 
										termination_description, created_by, created_at, 
										created_client, updated_by,updated_at,updated_client, 
										product_signature) 
								VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) 
								RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.LicenseConfigId.Int64,
		userParam.ProductKey.String,
		userParam.ProductEncrypt.String,
		userParam.ClientId.String,
		userParam.ClientSecret.String,
		userParam.HWID.String,
		userParam.ActivationDate.Time,
		userParam.LicenseStatus.Int32,
		userParam.TerminationDescription.String,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.ProductSignature.String,
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return

}

func (input productLicenseDAO) GetIDByLicenseConfigID(db *sql.Tx, userParam repository.ProductLicenseModel) (id int64, err errorModel.ErrorModel) {
	funcName := "GetIDByLicenseConfigID"

	query := fmt.Sprintf(`SELECT id FROM %s WHERE license_config_id = $1 AND deleted=FALSE`, input.TableName)

	params := []interface{}{
		userParam.LicenseConfigId.Int64,
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return

}

func (input productLicenseDAO) InsertBulkProductLicense(db *sql.Tx, userParam []repository.ProductLicenseModel) (output []int64, err errorModel.ErrorModel) {
	funcName := "InsertBulkProductLicense"
	paramLen := 15
	index := 1
	var params []interface{}
	var tempQuery string

	query := fmt.Sprintf(`INSERT INTO %s
	( 
		license_config_id, product_key, product_encrypt, client_id, 
		client_secret, hwid, activation_date, license_status, created_client, 
		created_by, created_at, updated_client, updated_by, updated_at, product_signature
	)
	VALUES `, input.TableName)

	tempQuery, index, params = ListValuesToInsertBulk(userParam, paramLen, index, func(inputVal interface{}) []interface{} {

		tempInputValue := inputVal.(repository.ProductLicenseModel)
		param := []interface{}{
			tempInputValue.LicenseConfigId.Int64,
			tempInputValue.ProductKey.String, tempInputValue.ProductEncrypt.String,
			tempInputValue.ClientId.String, tempInputValue.ClientSecret.String,
			tempInputValue.HWID.String, tempInputValue.ActivationDate.Time,
			tempInputValue.LicenseStatus.Int32, tempInputValue.CreatedClient.String,
			tempInputValue.CreatedBy.Int64, tempInputValue.CreatedAt.Time,
			tempInputValue.UpdatedClient.String, tempInputValue.UpdatedBy.Int64,
			tempInputValue.UpdatedAt.Time, tempInputValue.ProductSignature.String,
		}

		return param
	})

	query += tempQuery

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
			output = append(output, idTemp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseDAO) GetListProductLicense(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, inputStruct repository.ProductLicenseModel,
	scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, listScope []string) (result []interface{}, err errorModel.ErrorModel) {
	var (
		additionalWhere, query  string
		params, additionalParam []interface{}
	)

	query = fmt.Sprintf(
		`SELECT 
			pl.id as id_product_license, lc.id as id, c.customer_name as customer_name, 
			lc.unique_id_1 as unique_id_1, lc.unique_id_2 as unique_id_2, lc.installation_id as installation_id, 
			pr.product_name as product_name, lv.license_variant_name as license_variant_name, lt.license_type_name as license_type_name, 
			lc.product_valid_from as product_valid_from, lc.product_valid_thru as product_valid_thru, pl.license_status as license_status 
		FROM  %s pl 
		JOIN  %s lc ON pl.license_config_id = lc.id 
		JOIN  %s c ON c.id = lc.customer_id 
		JOIN  %s pr ON pr.id = lc.product_id 
		JOIN  %s lv ON lv.id = lc.license_variant_id 
		JOIN  %s lt ON lt.id = lc.license_type_id `,
		input.TableName, LicenseConfigDAO.TableName, CustomerDAO.TableName,
		ProductDAO.TableName, LicenseVariantDAO.TableName, LicenseTypeDAO.TableName)

	for i, param := range searchByParam {
		if searchByParam[i].SearchKey == "customer_id" {
			searchByParam[i].SearchKey = "c.id"
		} else {
			searchByParam[i].SearchKey = "c." + param.SearchKey
		}
	}

	additionalWhere, additionalParam = ScopeToAddedQueryView(scopeLimit, scopeDB, 1, listScope)

	if additionalWhere != "" {
		params = append(params, additionalParam...)
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ProductLicenseModelForView
			dbError := rows.Scan(
				&temp.ID,
				&temp.LicenseConfigId,
				&temp.CustomerName,
				&temp.UniqueId1,
				&temp.UniqueId2,
				&temp.InstallationId,
				&temp.ProductName,
				&temp.LicenseVariantName,
				&temp.LicenseTypeName,
				&temp.ProductValidFrom,
				&temp.ProductValidThru,
				&temp.LicenseStatus,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "pl.id"},
			Deleted:   FieldStatus{FieldName: "pl.deleted"},
			CreatedBy: FieldStatus{FieldName: "pl.created_by", Value: inputStruct.CreatedBy.Int64},
		})
}

func (input productLicenseDAO) ViewDetailProductLicenseById(db *sql.DB, userParam repository.DetailProductLicense,
	scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, listScope []string) (resultOnDB repository.DetailProductLicense, err errorModel.ErrorModel) {
	var (
		funcName                = "ViewDetailProductLicenseById"
		query, additionalWhere  string
		params, additionalParam []interface{}
	)

	query = fmt.Sprintf(
		`SELECT 
			pl.id, pl.product_key, pl.activation_date,
			pl.license_status, pl.termination_description, lc.id,
			lc.installation_id, lc.parent_customer_id, lc.customer_id,
			lc.site_id, lc.client_id, lc.deployment_method,
			lc.no_of_user, lc.is_user_concurrent, lc.unique_id_1,
			lc.unique_id_2, lc.product_valid_from, lc.product_valid_thru,
			pl.created_at, pl.updated_at, lt.license_type_name,
			lv.license_variant_name, ct.client_type, pr.product_name, c.customer_name, 
			cs.customer_name, m1.module_name, m2.module_name, 
			m3.module_name, m4.module_name, m5.module_name, 
			m6.module_name, m7.module_name, m8.module_name, 
			m9.module_name, m10.module_name, u.nt_username 
		FROM  %s pl 
		JOIN %s lc ON pl.license_config_id = lc.id 
		JOIN %s lt ON lt.id = lc.license_type_id 
		JOIN %s lv ON lv.id = lc.license_variant_id 
		JOIN %s ct ON ct.id = lc.client_type_id 
		JOIN %s pr ON pr.id = lc.product_id 
		JOIN %s c ON c.id = lc.customer_id 
		JOIN "%s" u ON u.id = pl.updated_by 
		LEFT JOIN %s cs ON cs.id = lc.parent_customer_id 
		LEFT JOIN %s m1 ON lc.module_id_1 = m1.id 
		LEFT JOIN %s m2 ON lc.module_id_2 = m2.id 
		LEFT JOIN %s m3 ON lc.module_id_3 = m3.id
        LEFT JOIN %s m4 ON lc.module_id_4 = m4.id 
		LEFT JOIN %s m5 ON lc.module_id_5 = m5.id 
		LEFT JOIN %s m6 ON lc.module_id_6 = m6.id 
		LEFT JOIN %s m7 ON lc.module_id_7 = m7.id 
		LEFT JOIN %s m8 ON lc.module_id_8 = m8.id 
		LEFT JOIN %s m9 ON lc.module_id_9 = m9.id 
		LEFT JOIN %s m10 ON lc.module_id_10 = m10.id 
		WHERE pl.id = $1 AND pl.deleted = false `,
		input.TableName, LicenseConfigDAO.TableName, LicenseTypeDAO.TableName,
		LicenseVariantDAO.TableName, ClientTypeDAO.TableName, ProductDAO.TableName,
		CustomerDAO.TableName, UserDAO.TableName, CustomerDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName)

	params = []interface{}{userParam.ID.Int64}

	additionalWhere, additionalParam = ScopeToAddedQueryView(scopeLimit, scopeDB, 2, listScope)

	if additionalWhere != "" {
		query += additionalWhere
		params = append(params, additionalParam...)
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&resultOnDB.ID,
		&resultOnDB.ProductKey,
		&resultOnDB.ActivationDate,
		&resultOnDB.LicenseStatus,
		&resultOnDB.TerminationDescription,
		&resultOnDB.LicenseConfigId,
		&resultOnDB.InstallationId,
		&resultOnDB.ParentCustomerId,
		&resultOnDB.CustomerId,
		&resultOnDB.SiteId,
		&resultOnDB.ClientId,
		&resultOnDB.DeploymentMethod,
		&resultOnDB.NumberOfUser,
		&resultOnDB.ConcurentUser,
		&resultOnDB.UniqueId1,
		&resultOnDB.UniqueId2,
		&resultOnDB.LicenseValidFrom,
		&resultOnDB.LicenseValidThru,
		&resultOnDB.CreatedAt,
		&resultOnDB.UpdatedAt,
		&resultOnDB.LicenseType,
		&resultOnDB.LicenseVariant,
		&resultOnDB.Client,
		&resultOnDB.Product,
		&resultOnDB.Customer,
		&resultOnDB.ParentCustomer,
		&resultOnDB.Module1,
		&resultOnDB.Module2,
		&resultOnDB.Module3,
		&resultOnDB.Module4,
		&resultOnDB.Module5,
		&resultOnDB.Module6,
		&resultOnDB.Module7,
		&resultOnDB.Module8,
		&resultOnDB.Module9,
		&resultOnDB.Module10,
		&resultOnDB.AliasName,
	)

	listComponentDB, err := input.GetComponentByIdLicenseConfig(serverconfig.ServerAttribute.DBConnection, resultOnDB.LicenseConfigId.Int64)
	if err.Error != nil {
		return
	}

	resultOnDB.Components = listComponentDB

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseDAO) CheckExistsProductLicenseById(db *sql.DB, productLicenseId int64) (isExists bool, err errorModel.ErrorModel) {
	funcName := "checkExistsProductLicenseById"

	query := "SELECT EXISTS(SELECT id FROM " + ProductLicenseDAO.TableName + " WHERE id = $1)"

	row := db.QueryRow(query, productLicenseId)

	dbErr := row.Scan(&isExists)
	if dbErr != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbErr)
		return
	}

	return
}

func (input productLicenseDAO) GetComponentByIdLicenseConfig(db *sql.DB, licenseConfigId int64) (result []repository.ProductComponentModel, err errorModel.ErrorModel) {
	funcName := "GetComponentByIdLicenseConfig"

	query := "SELECT c.id, c.component_name, lcpc.component_value " +
		"FROM " + ComponentDAO.TableName + " c " +
		"JOIN " + LicenseConfigProductComponentDAO.TableName + " lcpc ON lcpc.component_id = c.id " +
		"WHERE lcpc.license_config_id = $1"

	params := []interface{}{licenseConfigId}

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
			var temp repository.ProductComponentModel

			errorS = rows.Scan(
				&temp.ComponentID,
				&temp.ComponentName,
				&temp.ComponentValue,
			)

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

func (input productLicenseDAO) GetExpiredProductLicense(db *sql.DB, userParam in.AbstractDTO) (result []repository.ProductLicenseModel, err errorModel.ErrorModel) {
	funcName := "GetExpiredProductLicense"

	query := fmt.Sprintf(`SELECT 
		pl.id, pl.updated_at, pl.created_by
		FROM %s pl 
		INNER JOIN %s lc
		ON pl.license_config_id = lc.id
		WHERE 
			pl.deleted = FALSE AND 
			pl.license_status = $1 AND 
			lc.product_valid_thru < now()::date
		LIMIT $2 OFFSET $3 `,
		input.TableName, LicenseConfigDAO.TableName)

	param := []interface{}{constanta.ProductLicenseStatusActive, userParam.Limit, CountOffset(userParam.Page, userParam.Limit)}

	rows, errorS := db.Query(query, param...)
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
			var temp repository.ProductLicenseModel
			errorS = rows.Scan(
				&temp.ID, &temp.UpdatedAt, &temp.CreatedBy)

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

func (input productLicenseDAO) UpdatedForExpiredProduct(db *sql.Tx, userParam repository.ProductLicenseModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdatedForExpiredProduct"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET 
		license_status = $1, updated_by = $2, updated_at = $3, 
		updated_client = $4
		WHERE 
		id = $5`, input.TableName)

	param := []interface{}{
		userParam.LicenseStatus.Int32, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.ID.Int64,
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

func (input productLicenseDAO) GetProductLicenseForUpdate(db *sql.DB, productLicenseModel repository.ProductLicenseModel,
	scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, listScope []string) (result repository.ProductLicenseModel, err errorModel.ErrorModel) {
	funcName := "GetProductLicenseForUpdate"
	index := 2
	query := fmt.Sprintf(
		`SELECT 
			pl.id, pl.updated_at, pl.created_by,
			lc.product_valid_from, lc.product_valid_thru
		FROM %s pl
		INNER JOIN %s lc ON pl.license_config_id = lc.id
		JOIN  %s c ON c.id = lc.customer_id 
		JOIN  %s pr ON pr.id = lc.product_id 
		JOIN  %s lv ON lv.id = lc.license_variant_id 
		JOIN  %s lt ON lt.id = lc.license_type_id
		WHERE pl.id = $1 AND pl.deleted = FALSE `,
		input.TableName, LicenseConfigDAO.TableName, CustomerDAO.TableName,
		ProductDAO.TableName, LicenseVariantDAO.TableName, LicenseTypeDAO.TableName)

	params := []interface{}{productLicenseModel.ID.Int64}

	if productLicenseModel.CreatedBy.Int64 > 0 {
		query += fmt.Sprintf(" AND created_by = $%d ", index)
		params = append(params, productLicenseModel.CreatedBy.Int64)
		index++
	}

	additionalWhere, additionalParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, listScope)

	if additionalWhere != "" {
		query += additionalWhere
		params = append(params, additionalParam...)
	}

	query += " FOR UPDATE "

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy.Int64,
		&result.ProductValidFrom, &result.ProductValidThru)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseDAO) UpdateProductLicense(db *sql.Tx, userParam repository.ProductLicenseModel) (err errorModel.ErrorModel) {
	funcName := "UpdateProductLicense"

	query := "UPDATE " + input.TableName + " SET " +
		"license_status = $1, " +
		"termination_description = $2, " +
		"updated_by = $3, " +
		"updated_client = $4, " +
		"updated_at = $5 " +
		"WHERE id = $6 "

	param := []interface{}{
		userParam.LicenseStatus.Int32,
		userParam.TerminationDescription.String,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
		userParam.ID.Int64,
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

func (input productLicenseDAO) GetProductLicenseForValidation(db *sql.DB, userParam repository.ProductLicenseModel) (result repository.ProductLicenseModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetProductLicenseForValidation"
		tempResult interface{}
		query      string
	)

	query = fmt.Sprintf(`SELECT 
			pl.id, pl.license_config_id, pl.product_key,
			pl.product_encrypt, pl.product_signature, pl.client_id, 
			pl.client_secret, pl.hwid, lc.product_valid_from, 
			lc.product_valid_thru, p.product_id
		FROM %s pl
		LEFT JOIN %s lc ON pl.license_config_id = lc.id
		LEFT JOIN %s p ON lc.product_id = p.id
		WHERE 
			pl.deleted = FALSE AND pl.product_key = $1 AND
			pl.product_encrypt = $2 AND pl.client_id = $3 AND
			pl.client_secret = $4  AND pl.hwid = $5 `,
		input.TableName, LicenseConfigDAO.TableName, ProductDAO.TableName)

	param := []interface{}{
		userParam.ProductKey.String,
		userParam.ProductEncrypt.String,
		userParam.ClientId.String,
		userParam.ClientSecret.String,
		userParam.HWID.String,
	}

	rows := db.QueryRow(query, param...)
	tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ProductLicenseModel
		dbErrorS := rws.Scan(
			&temp.ID, &temp.LicenseConfigId, &temp.ProductKey,
			&temp.ProductEncrypt, &temp.ProductSignature, &temp.ClientId,
			&temp.ClientSecret, &temp.HWID, &temp.ProductValidFrom,
			&temp.ProductValidThru, &temp.ProductId,
		)
		return temp, dbErrorS
	}, input.FileName, funcName)

	if err.Error != nil {
		return
	}
	result = tempResult.(repository.ProductLicenseModel)
	return
}

func (input productLicenseDAO) UpdateProductLicenseForValidationLicense(db *sql.Tx, userParam repository.ProductLicenseModel) (err errorModel.ErrorModel) {
	funcName := "UpdateProductLicenseForValidationLicense"

	query := fmt.Sprintf(`UPDATE %s 
	SET 
		product_key = $1, 
		product_encrypt = $2, 
		product_signature = $3,
		updated_by = $4, 
		updated_at = $5, 
		updated_client = $6
	WHERE 
		license_config_id = $7 AND deleted = FALSE`, input.TableName)

	param := []interface{}{
		userParam.ProductKey.String, userParam.ProductEncrypt.String,
		userParam.ProductSignature.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.LicenseConfigId.Int64,
	}

	stmt, dbError := db.Prepare(query)
	defer stmt.Close()

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input productLicenseDAO) GetDataForDecryptProductLicense(db *sql.DB, userParam repository.ProductLicenseModel,
	scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, listScope []string) (result repository.ProductLicenseModel, err errorModel.ErrorModel) {
	var (
		fileName = "ProductLicenseDAO.go"
		funcName = "GetDataForDecryptProductLicense"
		query    string
		index    = 2
	)

	query = fmt.Sprintf(`SELECT 
			pl.id, pl.product_encrypt, pl.product_key, 
			pl.product_signature, pl.client_id, pl.client_secret, 
			pl.hwid, cc.signature_key, pr.product_id 
		FROM %s pl 
		JOIN %s cc ON cc.client_id = pl.client_id 
		JOIN %s lc ON lc.id = pl.license_config_id 
		JOIN  %s pr ON pr.id = lc.product_id 
		JOIN  %s c ON c.id = lc.customer_id 
		JOIN  %s lv ON lv.id = lc.license_variant_id 
		JOIN  %s lt ON lt.id = lc.license_type_id
		WHERE pl.id = $1 AND pl.deleted = FALSE`,
		input.TableName, ClientCredentialDAO.TableName, LicenseConfigDAO.TableName, ProductDAO.TableName,
		CustomerDAO.TableName, LicenseVariantDAO.TableName, LicenseTypeDAO.TableName)

	params := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		query += fmt.Sprintf(" AND created_by = $%d ", index)
		params = append(params, userParam.CreatedBy.Int64)
		index++
	}

	additionalWhere, additionalParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index, listScope)

	if additionalWhere != "" {
		query += additionalWhere
		params = append(params, additionalParam...)
	}

	errorS := db.QueryRow(query, params...).Scan(&result.ID, &result.ProductEncrypt, &result.ProductKey, &result.ProductSignature, &result.ClientId, &result.ClientSecret, &result.HWID, &result.SignatureKey, &result.ProductId)
	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseDAO) getCustomerDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "pl.id"},
		Deleted:   FieldStatus{FieldName: "pl.deleted"},
		CreatedBy: FieldStatus{FieldName: "pl.created_by", Value: createdBy},
	}
}

func (input productLicenseDAO) GetActiveLicenseForDashboard(db *sql.DB) (result []repository.ProductLicenseModel, err errorModel.ErrorModel) {
	funcName := "GetActiveLicenseForDashboard"
	query := fmt.Sprintf(
		`SELECT 
			pl.id, lc.product_valid_thru
		FROM %s pl
		LEFT JOIN %s lc ON lc.id = pl.license_config_id
		WHERE 
			date_part('MONTH', lc.product_valid_thru) IN (date_part('MONTH', NOW()) + 1, date_part('MONTH', NOW()) + 2) 
			AND date_part('YEAR', lc.product_valid_thru) IN (date_part('YEAR', NOW())) 
			AND pl.deleted = FALSE AND pl.license_status = $1 `,
		input.TableName, LicenseConfigDAO.TableName)

	rows, errorS := db.Query(query, constanta.ProductLicenseStatusActive)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	var tempResult []interface{}
	tempResult, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultInterface interface{}, errors errorModel.ErrorModel) {
		var temp repository.ProductLicenseModel
		dbError := rows.Scan(&temp.ID, &temp.ProductValidThru)
		if dbError != nil {
			errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}
		resultInterface = temp
		return
	})

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.ProductLicenseModel))
		}
	}

	return
}

func (input productLicenseDAO) UpdateProductLicenseForHWID(db *sql.Tx, userParam repository.ProductLicenseModel) (id int64, err errorModel.ErrorModel) {
	funcName := "UpdateProductLicenseForValidationLicense"

	query := fmt.Sprintf(`UPDATE %s 
	SET 
		product_key = $1, 
		product_encrypt = $2, 
		product_signature = $3,
		updated_by = $4, 
		updated_at = $5, 
		updated_client = $6,
		hwid = $7
	WHERE 
		license_config_id = $8 AND deleted = FALSE
	RETURNING id`, input.TableName)

	param := []interface{}{
		userParam.ProductKey.String, userParam.ProductEncrypt.String,
		userParam.ProductSignature.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.HWID.String, userParam.LicenseConfigId.Int64,
	}

	errorS := db.QueryRow(query, param...).Scan(&id)
	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input productLicenseDAO) GetProductLicenseForUpdateHWID(db *sql.DB, userParam repository.ProductLicenseModel) (result repository.ProductLicenseModel, err errorModel.ErrorModel) {
	var tempResult interface{}
	funcName := "GetProductLicenseForUpdateHWID"
	query := fmt.Sprintf(`SELECT 
			pl.id, pl.license_config_id, pl.product_key,
			pl.product_encrypt, pl.product_signature, pl.client_id, 
			pl.client_secret, pl.hwid, lc.product_valid_from, 
			lc.product_valid_thru, p.product_id
		FROM %s pl
		LEFT JOIN %s lc ON pl.license_config_id = lc.id
		LEFT JOIN %s p ON lc.product_id = p.id
		WHERE 
			pl.deleted = FALSE AND pl.product_key = $1 AND
			pl.product_encrypt = $2 AND pl.client_id = $3 AND
			pl.client_secret = $4 `,
		input.TableName, LicenseConfigDAO.TableName, ProductDAO.TableName)

	param := []interface{}{
		userParam.ProductKey.String,
		userParam.ProductEncrypt.String,
		userParam.ClientId.String,
		userParam.ClientSecret.String,
	}

	rows := db.QueryRow(query, param...)
	tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ProductLicenseModel
		dbErrorS := rws.Scan(
			&temp.ID, &temp.LicenseConfigId, &temp.ProductKey,
			&temp.ProductEncrypt, &temp.ProductSignature, &temp.ClientId,
			&temp.ClientSecret, &temp.HWID, &temp.ProductValidFrom,
			&temp.ProductValidThru, &temp.ProductId,
		)
		return temp, dbErrorS
	}, input.FileName, funcName)

	if err.Error != nil {
		return
	}
	result = tempResult.(repository.ProductLicenseModel)
	return
}

func (input productLicenseDAO) GetCountProductLicenseSalesJournal(db *sql.DB, searchBy []in.SearchByParam) (count int, err errorModel.ErrorModel) {
	var (
		query       string
		getListData getListJoinDataDAO
	)

	query = fmt.Sprintf(`SELECT COUNT(pl.id) FROM %s pl `, input.TableName)
	getListData = getListJoinDataDAO{Table: "pl", Query: query}
	input.setGetListJoinProductLicenseSalesJournal(&getListData)
	return getListData.GetCountJoinData(db, searchBy, 0)
}

func (input productLicenseDAO) GetProductLicenseSalesJournal(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query       string
		getListData getListJoinDataDAO
	)

	query = fmt.Sprintf(`SELECT 
		pl.id as id_license, lc.client_id, pl.license_status, 
		CASE  
			WHEN pl.license_status = 1 THEN 'Activated'
			WHEN pl.license_status = 2 THEN 'Expired'
			WHEN pl.license_status = 3 THEN 'Terminated'
			ELSE '-'
		END status_description,
		lc.unique_id_1, lc.unique_id_2, p.product_name, 
		ct.client_type, lc.allow_activation, lc.no_of_user, 
		lc.product_valid_from, lc.product_valid_thru, lc.is_user_concurrent, 
		ul.total_license, ul.total_activated 
		FROM %s pl `,
		input.TableName)

	getListData = getListJoinDataDAO{Table: "pl", Query: query}
	input.setGetListJoinProductLicenseSalesJournal(&getListData)
	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.LicenseSalesJournal
		dbError := rows.Scan(
			&resultTemp.ID, &resultTemp.ClientID, &resultTemp.LicenseStatusID,
			&resultTemp.LicenseStatus, &resultTemp.UniqueID1, &resultTemp.UniqueID2,
			&resultTemp.ProductName, &resultTemp.ClientType, &resultTemp.AllowActivation,
			&resultTemp.NoOfUser, &resultTemp.ProductValidFrom, &resultTemp.ProductValidThru,
			&resultTemp.IsUserConcurrent, &resultTemp.TotalLicense, &resultTemp.TotalActivated)

		return resultTemp, dbError
	}

	userParam.OrderBy = "pl.id ASC"
	if userParam.Limit > 0 {
		return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
	}

	return getListData.GetListJoinDataWithoutPagination(db, userParam, searchBy, 0, mappingFunc)
}

func (input productLicenseDAO) setGetListJoinProductLicenseSalesJournal(getListData *getListJoinDataDAO) {
	getListData.InnerJoinAlias(LicenseConfigDAO.TableName, "lc", "pl.license_config_id", "lc.id")
	getListData.InnerJoinAlias(ProductDAO.TableName, "p", "lc.product_id", "p.id")
	getListData.InnerJoinAlias(ClientTypeDAO.TableName, "ct", "lc.client_type_id", "ct.id")
	getListData.LeftJoinAliasWithoutDeleted(UserLicenseDAO.TableName, "ul", "ul.product_license_id", "pl.id and ul.deleted = false")
}
