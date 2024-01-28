package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type productDAO struct {
	AbstractDAO
}

var ProductDAO = productDAO{}.New()

func (input productDAO) New() (output productDAO) {
	output.FileName = "ProductDAO.go"
	output.TableName = "product"
	return
}

func (input productDAO) UpdateProduct(db *sql.Tx, userParam repository.ProductModel) (err errorModel.ErrorModel) {
	funcName := "UpdateProduct"

	query := fmt.Sprintf(`UPDATE %s 
								SET 
									product_id = $1, product_name = $2, product_description = $3, 
									product_group_id = $4, client_type_id = $5, is_license = $6, 
									license_variant_id = $7, license_type_id = $8, deployment_method = $9, 
									no_of_user = $10, is_user_concurrent = $11, module_id_1 = $12, 
									module_id_2 = $13, module_id_3 = $14, module_id_4 = $15, 
									module_id_5 = $16, module_id_6 = $17, module_id_7 = $18, 
									module_id_8 = $19, module_id_9 = $20, module_id_10 = $21, 
									updated_by = $22, updated_at = $23, updated_client = $24, max_offline_days = $25 
								WHERE id = $26 `, input.TableName)

	param := []interface{}{userParam.ProductID.String, userParam.ProductName.String}

	if userParam.ProductDescription.String != "" {
		param = append(param, userParam.ProductDescription.String)
	} else {
		param = append(param, nil)
	}

	param = append(param,
		userParam.ProductGroupID.Int64, userParam.ClientTypeID.Int64, userParam.IsLicense.Bool,
		userParam.LicenseVariantID.Int64, userParam.LicenseTypeID.Int64, userParam.DeploymentMethod.String,
		userParam.NoOfUser.Int64, userParam.IsUserConcurrent.Bool)

	input.addModuleFilled(&param, userParam)
	param = append(param, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.MaxOfflineDays.Int64, userParam.ID.Int64)

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	defer stmt.Close()

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	return
}

func (input productDAO) InsertProduct(db *sql.Tx, userParam repository.ProductModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertProduct"

	query := "INSERT " +
		"INTO " + input.TableName + " " +
		"(product_name, product_description, product_group_id, " +
		"client_type_id, is_license, license_variant_id, " +
		"license_type_id, deployment_method, no_of_user, " +
		"is_user_concurrent, module_id_1, module_id_2, " +
		"module_id_3, module_id_4, module_id_5, " +
		"module_id_6, module_id_7, module_id_8, " +
		"module_id_9, module_id_10, created_by, " +
		"created_client, created_at, updated_by, " +
		"updated_client, updated_at, product_id, " +
		"max_offline_days) " +
		"VALUES " +
		"($1, $2, $3, " +
		"$4, $5, $6, " +
		"$7, $8, $9, " +
		"$10, $11, $12, " +
		"$13, $14, $15, " +
		"$16, $17, $18, " +
		"$19, $20, $21, " +
		"$22, $23, $24, " +
		"$25, $26, $27, " +
		"$28) " +
		"RETURNING id "

	var params []interface{}

	params = append(params,
		userParam.ProductName.String, userParam.ProductDescription.String, userParam.ProductGroupID.Int64,
		userParam.ClientTypeID.Int64, userParam.IsLicense.Bool, userParam.LicenseVariantID.Int64,
		userParam.LicenseTypeID.Int64, userParam.DeploymentMethod.String, userParam.NoOfUser.Int64,
		userParam.IsUserConcurrent.Bool)

	input.addModuleFilled(&params, userParam)

	params = append(params,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.ProductID.String, userParam.MaxOfflineDays.Int64)

	result := db.QueryRow(query, params...)
	errorS := result.Scan(&id)

	if errorS != nil && errorS.Error() != constanta.NoRowsInDB {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input productDAO) GetCountProduct(db *sql.DB, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {

	query := fmt.Sprintf(`SELECT COUNT(%s.id) FROM %s `, input.TableName, input.TableName)

	colAdditionalWhere := input.setScopeData(scopeLimit, scopeDB, false)

	input.setSearchByProduct(&searchBy)
	input.setCreatedByProduct(createdBy, &searchBy)
	getListData := getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: colAdditionalWhere}
	input.setGetListJoinProduct(&getListData)
	return getListData.GetCountJoinData(db, searchBy, 0)
}

func (input productDAO) GetListProduct(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query              string
		colAdditionalWhere []string
		getListData        getListJoinDataDAO
	)

	query = fmt.Sprintf(`
		SELECT product.id as id, product.product_name as product_name, product.product_description as product_description, 
		product_group.product_group_name as product_group_name, client_type.client_type as client_type_name, license_variant.license_variant_name as license_variant_name, 
		license_type.license_type_name as license_type_name, product.product_id as product_id, product.updated_at, 
		client_type.parent_client_type_id, product.client_type_id
		FROM %s `, input.TableName)

	colAdditionalWhere = input.setScopeData(scopeLimit, scopeDB, true)
	input.setSearchByProduct(&searchBy)
	input.setCreatedByProduct(createdBy, &searchBy)

	getListData = getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: colAdditionalWhere}
	input.setGetListJoinProduct(&getListData)

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.ProductModel
		dbError := rows.Scan(
			&resultTemp.ID, &resultTemp.ProductName, &resultTemp.ProductDescription,
			&resultTemp.ProductGroupName, &resultTemp.ClientTypeName, &resultTemp.LicenseVariantName,
			&resultTemp.LicenseTypeName, &resultTemp.ProductID, &resultTemp.UpdatedAt,
			&resultTemp.ParentClientTypeID, &resultTemp.ClientTypeID,
		)

		return resultTemp, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input productDAO) ViewProduct(db *sql.DB, userParam repository.ProductModel) (result repository.ProductModel, err errorModel.ErrorModel) {
	funcName := "ViewProduct"

	subQuery := "SELECT pr.id, pr.product_id, pr.product_name, " +
		"pr.product_description, pr.product_group_id, pg.product_group_name, " +
		"pr.client_type_id, ct.client_type, pr.is_license, " +
		"pr.license_variant_id, lv.license_variant_name, pr.license_type_id, " +
		"lt.license_type_name, pr.deployment_method, pr.no_of_user, " +
		"pr.is_user_concurrent, pr.module_id_1, mo1.module_name, " +
		"pr.module_id_2, mo2.module_name, pr.module_id_3, " +
		"mo3.module_name, pr.module_id_4, mo4.module_name, " +
		"pr.module_id_5, mo5.module_name, pr.module_id_6, " +
		"mo6.module_name, pr.module_id_7, mo7.module_name, " +
		"pr.module_id_8, mo8.module_name, pr.module_id_9, " +
		"mo9.module_name, pr.module_id_10, mo10.module_name, " +
		"pr.created_at, pr.updated_at, pr.updated_by as updated_by, " +
		"pr.max_offline_days, pr.created_by " +
		"FROM " + input.TableName + " pr " +
		"INNER JOIN " + ProductGroupDAO.TableName + " pg ON pg.id = pr.product_group_id " +
		"INNER JOIN " + ClientTypeDAO.TableName + " ct ON ct.id = pr.client_type_id " +
		"INNER JOIN " + LicenseVariantDAO.TableName + " lv ON lv.id = pr.license_variant_id " +
		"INNER JOIN " + LicenseTypeDAO.TableName + " lt ON lt.id = pr.license_type_id " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo1 ON mo1.id = pr.module_id_1 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo2 ON mo2.id = pr.module_id_2 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo3 ON mo3.id = pr.module_id_3 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo4 ON mo4.id = pr.module_id_4 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo5 ON mo5.id = pr.module_id_5 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo6 ON mo6.id = pr.module_id_6 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo7 ON mo7.id = pr.module_id_7 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo8 ON mo8.id = pr.module_id_8 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo9 ON mo9.id = pr.module_id_9 " +
		"LEFT JOIN " + ModuleDAO.TableName + " mo10 ON mo10.id = pr.module_id_10 WHERE " +
		"pr.id = $1 AND pr.deleted = FALSE "

	params := []interface{}{userParam.ID.Int64}

	if userParam.CreatedBy.Int64 > 0 {
		subQuery += " AND pr.created_by = $2 "
		params = append(params, userParam.CreatedBy.Int64)
	}

	query := "SELECT *, (SELECT nt_username as updated_name FROM \"user\" WHERE id = a.updated_by) FROM (" + subQuery + ") a "

	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&result.ID, &result.ProductID, &result.ProductName,
		&result.ProductDescription, &result.ProductGroupID, &result.ProductGroupName,
		&result.ClientTypeID, &result.ClientTypeName, &result.IsLicense,
		&result.LicenseVariantID, &result.LicenseVariantName, &result.LicenseTypeID,
		&result.LicenseTypeName, &result.DeploymentMethod, &result.NoOfUser,
		&result.IsUserConcurrent, &result.Module1, &result.ModuleName1,
		&result.Module2, &result.ModuleName2, &result.Module3,
		&result.ModuleName3, &result.Module4, &result.ModuleName4,
		&result.Module5, &result.ModuleName5, &result.Module6,
		&result.ModuleName6, &result.Module7, &result.ModuleName7,
		&result.Module8, &result.ModuleName8, &result.Module9,
		&result.ModuleName9, &result.Module10, &result.ModuleName10,
		&result.CreatedAt, &result.UpdatedAt, &result.UpdatedBy,
		&result.MaxOfflineDays, &result.CreatedBy, &result.UpdatedName)

	if dbError != nil && dbError.Error() != constanta.NoRowsInDB {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) GetProductForUpdate(db *sql.DB, productModel repository.ProductModel) (result []repository.GetForUpdateProduct, err errorModel.ErrorModel) {
	funcName := "GetProductForUpdate"

	query := fmt.Sprintf(`SELECT product.id, product.updated_at, product.created_by, pc.id, CASE WHEN 
		(SELECT count(id) FROM %s WHERE product_id = product.id AND deleted = FALSE) > 0 
		OR 
		(SELECT count(id) FROM %s WHERE product_id = product.id AND deleted = FALSE) > 0 
		THEN TRUE ELSE FALSE END is_used, product.product_name, product.product_id 
		FROM %s 
		LEFT JOIN %s pc 
		ON product.id = pc.product_id WHERE 
		product.id = $1 AND product.deleted = FALSE `,
		CustomerInstallationDAO.TableName, LicenseConfigDAO.TableName, input.TableName,
		ProductComponentDAO.TableName)

	params := []interface{}{productModel.ID.Int64}
	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	var resultTemp []interface{}
	resultTemp, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.GetForUpdateProduct))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) DeleteProduct(db *sql.Tx, userParam repository.ProductModel) (err errorModel.ErrorModel) {
	funcName := "DeleteProduct"

	query := fmt.Sprintf(`UPDATE %s 
		SET deleted = TRUE, updated_by = $1, updated_client = $2, 
		updated_at = $3, product_name = $4, product_id = $5 
		WHERE id = $6`, input.TableName)

	param := []interface{}{
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.ProductName.String, userParam.ProductID.String, userParam.ID.Int64}

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

func (input productDAO) addModuleFilled(params *[]interface{}, userParam repository.ProductModel) {
	if userParam.Module1.Int64 > 0 {
		*params = append(*params, userParam.Module1.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module2.Int64 > 0 {
		*params = append(*params, userParam.Module2.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module3.Int64 > 0 {
		*params = append(*params, userParam.Module3.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module4.Int64 > 0 {
		*params = append(*params, userParam.Module4.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module5.Int64 > 0 {
		*params = append(*params, userParam.Module5.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module6.Int64 > 0 {
		*params = append(*params, userParam.Module6.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module7.Int64 > 0 {
		*params = append(*params, userParam.Module7.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module8.Int64 > 0 {
		*params = append(*params, userParam.Module8.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module9.Int64 > 0 {
		*params = append(*params, userParam.Module9.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.Module10.Int64 > 0 {
		*params = append(*params, userParam.Module10.Int64)
	} else {
		*params = append(*params, nil)
	}
}

func (input productDAO) setSearchByProduct(searchBy *[]in.SearchByParam) {
	temp := *searchBy
	for index := range temp {
		switch temp[index].SearchKey {
		case "product_id":
			temp[index].SearchKey = "product." + temp[index].SearchKey
		case "product_name":
			temp[index].SearchKey = "product." + temp[index].SearchKey
		case "product_group_id":
			temp[index].SearchKey = "product." + temp[index].SearchKey
		}
	}
}

func (input productDAO) setCreatedByProduct(createdBy int64, searchBy *[]in.SearchByParam) {
	if createdBy > 0 {
		*searchBy = append(*searchBy, in.SearchByParam{
			SearchKey:      "product.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}
}

func (input productDAO) setGetListJoinProduct(getListData *getListJoinDataDAO) {
	getListData.InnerJoin(ProductGroupDAO.TableName, "product_group.id", "product.product_group_id")
	getListData.InnerJoin(ClientTypeDAO.TableName, "client_type.id", "product.client_type_id")
	getListData.InnerJoin(LicenseVariantDAO.TableName, "license_variant.id", "product.license_variant_id")
	getListData.InnerJoin(LicenseTypeDAO.TableName, "license_type.id", "product.license_type_id")
}

func (input productDAO) CheckProductIsExist(db *sql.DB, userParam repository.ProductModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName        = "CheckProductIsExist"
		query           string
		params          []interface{}
		results         *sql.Row
		dbError         error
		additionalWhere []string
	)

	query = fmt.Sprintf(`
		SELECT (CASE WHEN count(pd.id) > 0 THEN TRUE ELSE FALSE END) is_exist 
		FROM %s pd 
		WHERE pd.id = $1 AND pd.deleted = FALSE `, input.TableName)

	if scopeLimit != nil || scopeDB != nil {
		additionalWhere = input.setScopeData(scopeLimit, scopeDB, true)
		if len(additionalWhere) > 0 {
			strWhere := " AND " + strings.Join(additionalWhere, " AND ")
			strWhere = strings.TrimRight(strWhere, " AND ")
			query += strWhere
		}
	}

	params = []interface{}{userParam.ID.Int64}
	results = db.QueryRow(query, params...)
	dbError = results.Scan(&isExist)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) resultRowsInput(rows *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
	funcName := "resultRowsInput"
	var errorS error
	var temp repository.GetForUpdateProduct

	errorS = rows.Scan(&temp.ID, &temp.UpdatedAt, &temp.CreatedBy,
		&temp.ProductComponentID, &temp.IsUsed, &temp.ProductName,
		&temp.ProductID)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp = temp
	return
}

func (input productDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (colAdditionalWhere []string) {
	keyScope := []string{
		constanta.ClientTypeDataScope,
		constanta.ProductGroupDataScope,
	}

	for _, itemKeyScope := range keyScope {
		var additionalWhere string
		PrepareScopeOnDAO(scopeLimit, scopeDB, &additionalWhere, 0, itemKeyScope, isView)
		if additionalWhere != "" {
			colAdditionalWhere = append(colAdditionalWhere, additionalWhere)
		}
	}

	return
}

func (input productDAO) LockProductForUpdate(db *sql.DB, userParam repository.ProductModel) (result repository.ProductModel, err errorModel.ErrorModel) {
	funcName := "LockProductForUpdate"

	query := "SELECT id " +
		"FROM " + input.TableName + " " +
		"WHERE id = $1 AND deleted = FALSE "

	param := []interface{}{userParam.ID.Int64}

	query += fmt.Sprintf(` FOR UPDATE `)

	dbError := db.QueryRow(query, param...).Scan(
		&result.ID,
	)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) CheckProductAndParentClientTypeIsExist(db *sql.DB, userParam repository.ProductModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.ProductModel, err errorModel.ErrorModel) {
	var (
		funcName        = "CheckProductAndParentClientTypeIsExist"
		query           string
		params          []interface{}
		results         *sql.Row
		dbError         error
		additionalWhere []string
	)

	query = fmt.Sprintf(`select pd.id, ct.parent_client_type_id, pd.product_id  
		from %s pd 
		inner join %s ct on pd.client_type_id = ct.id 
		where 
		pd.id = $1 and pd.deleted = false and ct.deleted = false `,
		input.TableName, ClientTypeDAO.TableName)

	params = []interface{}{userParam.ID.Int64}
	if userParam.ParentClientTypeID.Int64 > 0 {
		query += fmt.Sprintf(` and ct.parent_client_type_id = $2 `)
		params = append(params, userParam.ParentClientTypeID.Int64)
	}

	if scopeLimit != nil || scopeDB != nil {
		additionalWhere = input.setScopeData(scopeLimit, scopeDB, true)
		if len(additionalWhere) > 0 {
			strWhere := " AND " + strings.Join(additionalWhere, " AND ")
			strWhere = strings.TrimRight(strWhere, " AND ")
			query += strWhere
		}
	}

	results = db.QueryRow(query, params...)
	dbError = results.Scan(&result.ID, &result.ParentClientTypeID, &result.ProductID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) CheckValidParentProductByIDAndClientType(db *sql.DB, userParam []repository.ProductModel, clientTypeID int64) (isValid bool, err errorModel.ErrorModel) {
	var (
		funcName    = "CheckValidParentProductByIDAndClientType"
		query       string
		params      []interface{}
		results     *sql.Row
		dbError     error
		dollarParam = 1
	)

	query = fmt.Sprintf(`select case when count(pd.id) > 0 
			then true else false 
			end is_valid 
		from %s pd 
		inner join %s ct on pd.client_type_id = ct.id
		where pd.deleted = false and ct.deleted = false `,
		input.TableName, ClientTypeDAO.TableName)

	if len(userParam) > 0 {
		query += fmt.Sprintf(` and pd.id in (`)
		for i := 1; i <= len(userParam); i++ {
			query += fmt.Sprintf(`$%d`, i)
			params = append(params, userParam[i-1].ID.Int64)
			if len(userParam)-i > 0 {
				query += ","
			} else {
				query += ") "
			}
			dollarParam++
		}
	}

	query += fmt.Sprintf(` and ct.id = $%d `, dollarParam)
	params = append(params, clientTypeID)

	results = db.QueryRow(query, params...)
	dbError = results.Scan(&isValid)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) CheckValidParentProductByProductID(db *sql.DB, userParam repository.ProductModel) (isParent bool, result repository.ProductModel, err errorModel.ErrorModel) {
	var (
		funcName = "CheckValidParentProductByProductID"
		query    string
		params   []interface{}
		results  *sql.Row
		dbError  error
	)

	query = fmt.Sprintf(`select ct.id, case when ct.parent_client_type_id is null then true else false end is_parent 
		from %s p 
		inner join %s ct on p.client_type_id = ct.id 
		where p.id = $1 and p.deleted = false and ct.deleted = false `,
		input.TableName, ClientTypeDAO.TableName)

	params = append(params, userParam.ID.Int64)
	results = db.QueryRow(query, params...)
	dbError = results.Scan(&result.ParentClientTypeID, &isParent)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) GetProductChild(db *sql.DB, productModel repository.ProductModel) (result map[int64]bool, err errorModel.ErrorModel) {
	var (
		funcName   = "GetProductChild"
		query      string
		resultTemp []interface{}
	)

	query = fmt.Sprintf(`select id from %s where client_type_id in 
		(select id from %s where parent_client_type_id = $1 and deleted = false) and deleted = false `,
		input.TableName, ClientTypeDAO.TableName)

	params := []interface{}{productModel.ClientTypeID.Int64}
	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			dbError error
			id      int64
		)

		dbError = rws.Scan(&id)
		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}
		resultTemp = id
		return
	})

	if err.Error != nil {
		return
	}

	result = make(map[int64]bool)
	for _, itemResultTemp := range resultTemp {
		result[itemResultTemp.(int64)] = true
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productDAO) GetProductParentByClientTypeID(db *sql.DB, productModel repository.ProductModel) (result map[int64]bool, err errorModel.ErrorModel) {
	var (
		funcName   = "GetProductParentByClientTypeID"
		query      string
		resultTemp []interface{}
	)

	query = fmt.Sprintf(`select id from %s where client_type_id = $1 and deleted = false `, input.TableName)

	params := []interface{}{productModel.ClientTypeID.Int64}
	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			dbError error
			id      int64
		)

		dbError = rws.Scan(&id)
		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}
		resultTemp = id
		return
	})

	if err.Error != nil {
		return
	}

	result = make(map[int64]bool)
	for _, itemResultTemp := range resultTemp {
		result[itemResultTemp.(int64)] = true
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
