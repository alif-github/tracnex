package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type userLicenseDAO struct {
	AbstractDAO
}

var UserLicenseDAO = userLicenseDAO{}.New()

func (input userLicenseDAO) New() (output userLicenseDAO) {
	output.FileName = "UserLicenseDAO.go"
	output.TableName = "user_license"
	return
}

func (input userLicenseDAO) GetCountUserLicense(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (int, errorModel.ErrorModel) {

	query := "SELECT COUNT (ul.id) FROM " + input.TableName + " ul "

	input.setSearchByUserLicense(&searchByParam)
	input.setCreatedByUserLicense(createdBy, &searchByParam)

	additionalWhere := input.setScopeData(scopeLimit, scopeDB, true)
	getListData := getListJoinDataDAO{Table: "ul", Query: query, AdditionalWhere: additionalWhere}

	input.setGetListJoinUserLicense(&getListData)

	return getListData.GetCountJoinData(db, searchByParam, 0)
}

func (input userLicenseDAO) setSearchByUserLicense(searchBy *[]in.SearchByParam) {
	temp := *searchBy

	for index := range temp {
		switch temp[index].SearchKey {
		case "id":
			temp[index].SearchKey = "ul." + temp[index].SearchKey
		case "customer_name":
			temp[index].SearchKey = "cu." + temp[index].SearchKey
		}
	}
}

func (input userLicenseDAO) setCreatedByUserLicense(createdBy int64, searchBy *[]in.SearchByParam) {
	if createdBy > 0 {
		*searchBy = append(*searchBy, in.SearchByParam{
			SearchKey:      "ul.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}
}

func (input userLicenseDAO) setGetListJoinUserLicense(getListData *getListJoinDataDAO) {
	getListData.InnerJoinAlias(ProductLicenseDAO.TableName, "pl", "pl.id", "ul.product_license_id")
	getListData.InnerJoinAlias(LicenseConfigDAO.TableName, "lc", "lc.id", "pl.license_config_id")
	getListData.InnerJoinAlias(CustomerDAO.TableName, "cu", "cu.id", "lc.customer_id")
	getListData.InnerJoinAlias(ProductDAO.TableName, "p", "p.id", "lc.product_id")
	getListData.InnerJoinAlias(ProductGroupDAO.TableName, "pg", "pg.id", "p.product_group_id")
	getListData.InnerJoinAlias(ClientTypeDAO.TableName, "ct", "ct.id", "lc.client_type_id")
}

func (input userLicenseDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (additionalWhereWithScope []string) {
	keyScope := []string{
		constanta.ProvinceDataScope,
		constanta.DistrictDataScope,
		constanta.CustomerGroupDataScope,
		constanta.CustomerCategoryDataScope,
		constanta.SalesmanDataScope,
		constanta.ProductGroupDataScope,
		constanta.ClientTypeDataScope,
	}

	for _, itemKeyScope := range keyScope {
		var additionalWhere string
		PrepareScopeOnDAO(scopeLimit, scopeDB, &additionalWhere, 0, itemKeyScope, isView)
		if additionalWhere != "" {
			additionalWhereWithScope = append(additionalWhereWithScope, additionalWhere)
		}
	}

	return
}

func (input userLicenseDAO) GetListUserLicense(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var params []interface{}

	query := "SELECT ul.id as id, lc.id as license_config_id, cu.customer_name as customer_name, ul.unique_id_1 as unique_id_1, ul.unique_id_2 as unique_id_2, ul.installation_id as installation_id, ul.total_license as total_license, ul.total_activated as total_active " +
		"FROM " + input.TableName + " ul " +
		"JOIN " + ProductLicenseDAO.TableName + " pl ON ul.product_license_id = pl.id " +
		"JOIN " + LicenseConfigDAO.TableName + " lc ON lc.id = pl.license_config_id " +
		"JOIN " + CustomerDAO.TableName + " cu ON cu.id = lc.customer_id " +
		"JOIN " + ProductDAO.TableName + " p ON p.id = lc.product_id " +
		"JOIN " + ProductGroupDAO.TableName + " pg ON pg.id = p.product_group_id " +
		"JOIN " + ClientTypeDAO.TableName + " ct ON ct.id = lc.client_type_id "

	arrStrScope := input.setScopeData(scopeLimit, scopeDB, true)
	if len(arrStrScope) > 0 {
		strWhere := " AND " + strings.Join(arrStrScope, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	for i, param := range searchByParam {
		if searchByParam[i].SearchKey == "id" {
			searchByParam[i].SearchKey = "ul." + param.SearchKey
		} else {
			searchByParam[i].SearchKey = "cu." + param.SearchKey
		}
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.UserLicenseModel
			dbError := rows.Scan(
				&temp.ID,
				&temp.LicenseConfigId,
				&temp.CustomerName,
				&temp.UniqueId1,
				&temp.UniqueId2,
				&temp.InstallationId,
				&temp.TotalLicense,
				&temp.TotalActivated,
			)
			return temp, dbError
		}, "", DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "ul.id"},
			Deleted:   FieldStatus{FieldName: "ul.deleted"},
			CreatedBy: FieldStatus{FieldName: "ul.created_by", Value: createdBy},
		})
}

func (input userLicenseDAO) ViewDetailUserLicense(db *sql.DB, userLicenseParam repository.UserLicenseModel) (userLicenseOnDB repository.UserLicenseModel, err errorModel.ErrorModel) {
	funcName := "ViewDetailUserLicense"

	query := fmt.Sprintf(`SELECT 
			ul.id, lc.id, ul.installation_id, lc.parent_customer_id, 
			pcu.customer_name as parent_customer_name, lc.customer_id, 
			lc.site_id, cu.customer_name, ul.unique_id_1 as company_id, 
			ul.unique_id_2 as branch_id, pr.product_name, 
			ul.total_activated as active_user, ul.total_license,  
			ul.product_valid_from, ul.product_valid_thru, pl.license_status, ul.updated_at,
			pr.client_type_id
		FROM %s ul  
		JOIN %s pl on ul.product_license_id = pl.id  
		JOIN %s lc on lc.id = pl.license_config_id  
		JOIN %s pcu on pcu.id = lc.parent_customer_id  
		JOIN %s cu on cu.id = lc.customer_id  
		JOIN %s pr on lc.product_id = pr.id  
		WHERE ul.id = $1 AND ul.deleted = FALSE `,
		input.TableName, ProductLicenseDAO.TableName, LicenseConfigDAO.TableName,
		CustomerDAO.TableName, CustomerDAO.TableName, ProductDAO.TableName)

	params := []interface{}{userLicenseParam.ID.Int64}

	if userLicenseParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		params = append(params, userLicenseParam.CreatedBy.Int64)
	}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&userLicenseOnDB.ID,
		&userLicenseOnDB.LicenseConfigId,
		&userLicenseOnDB.InstallationId,
		&userLicenseOnDB.ParentCustomerId,
		&userLicenseOnDB.ParentCustomerName,
		&userLicenseOnDB.CustomerId,
		&userLicenseOnDB.SiteId,
		&userLicenseOnDB.CustomerName,
		&userLicenseOnDB.UniqueId1,
		&userLicenseOnDB.UniqueId2,
		&userLicenseOnDB.ProductName,
		&userLicenseOnDB.TotalActivated,
		&userLicenseOnDB.TotalLicense,
		&userLicenseOnDB.ProductValidFrom,
		&userLicenseOnDB.ProductValidThru,
		&userLicenseOnDB.LicenseStatus,
		&userLicenseOnDB.UpdatedAt,
		&userLicenseOnDB.ClientTypeId,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) CheckLicenseNamedUser(db *sql.DB, userLicenseParam repository.UserLicenseModel) (userLicenseOnDb repository.UserLicenseModel, err errorModel.ErrorModel) {
	funcName := "CheckLicenseNamedUser"

	query := "SELECT pl.id, pl.product_key, ul.total_license, " +
		"ul.total_activated, (ul.total_license - ul.total_activated) as quota_license, pl.license_status, " +
		"ul.parent_customer_id, ul.customer_id, ul.site_id, " +
		"ul.installation_id, ul.id, ul.product_valid_from, " +
		"ul.product_valid_thru " +
		"FROM " + ProductLicenseDAO.TableName + " pl " +
		"JOIN " + input.TableName + " ul ON pl.id = ul.product_license_id " +
		"JOIN " + CustomerInstallationDAO.TableName + " ci ON ci.id = ul.installation_id " +
		"JOIN " + ProductDAO.TableName + " pr ON pr.id = ci.product_id " +
		"WHERE pl.license_status = 1 AND " +
		"ul.unique_id_1 = $1 AND " +
		"ul.client_id = $2 AND " +
		"pr.client_type_id = $3 AND "

	params := []interface{}{userLicenseParam.UniqueId1, userLicenseParam.ClientID, userLicenseParam.ClientTypeId}

	if userLicenseParam.UniqueId2.Valid {
		query += "ul.unique_id_2 = $4 AND "
		params = append(params, userLicenseParam.UniqueId2)
	}

	query += "ul.product_valid_from <= CURRENT_DATE AND " +
		"ul.product_valid_thru >= CURRENT_DATE AND " +
		"pl.deleted = FALSE AND " +
		"ul.deleted = FALSE AND " +
		"ci.deleted = FALSE " +
		"ORDER BY ul.id ASC LIMIT 1 "

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&userLicenseOnDb.ProductLicenseID,
		&userLicenseOnDb.ProductKey,
		&userLicenseOnDb.TotalLicense,
		&userLicenseOnDb.TotalActivated,
		&userLicenseOnDb.QuotaLicense,
		&userLicenseOnDb.LicenseStatus,
		&userLicenseOnDb.ParentCustomerId,
		&userLicenseOnDb.CustomerId,
		&userLicenseOnDb.SiteId,
		&userLicenseOnDb.InstallationId,
		&userLicenseOnDb.ID,
		&userLicenseOnDb.ProductValidFrom,
		&userLicenseOnDb.ProductValidThru,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) ReduceTotalLicense(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (err errorModel.ErrorModel) {
	funcName := "ReduceTotalLicense"

	query := fmt.Sprintf(
		`UPDATE %s SET 
					total_activated = total_activated-1, 
					updated_by = $1, 
					updated_client = $2, 
					updated_at = $3 
				WHERE id = $4 `, input.TableName)

	param := []interface{}{
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt,
		userParam.UserLicenseID,
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

func (input userLicenseDAO) GetFieldForValidationUnregister(db *sql.DB, userRegDetailModel repository.UserRegistrationDetailModel) (result repository.UserLicenseModel, err errorModel.ErrorModel) {
	funcName := "GetFieldForValidationUnregister"
	query := "SELECT ul.id, ul.client_id, ul.total_activated FROM user_license ul JOIN user_registration_detail urd ON ul.id = urd.user_license_id WHERE urd.id = $1 AND ul.deleted = FALSE "

	params := []interface{}{
		userRegDetailModel.ID.Int64,
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.ID, &result.ClientID, &result.TotalActivated)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) InsertBulkUserLicenses(db *sql.Tx, userParam []repository.UserLicenseModel) (result []int64, err errorModel.ErrorModel) {
	funcName := "InsertBulkUserLicenses"
	paramLen := 17
	index := 1

	var tempResult []interface{}
	var params []interface{}
	var tempQuery string

	query := fmt.Sprintf(`INSERT INTO %s 
		(
			product_license_id, parent_customer_id, customer_id, 
			site_id, installation_id, client_id, unique_id_1, 
			unique_id_2, product_valid_from, product_valid_thru, 
			total_license, total_activated, created_client, 
			created_by, created_at, updated_client, updated_by, 
			updated_at
		) 
		VALUES `, input.TableName)

	for i := 0; i < len(userParam); i++ {
		query += " ( (SELECT id FROM " + ProductLicenseDAO.TableName + " WHERE " +
			" product_key = $" + strconv.Itoa(index) + " AND " +
			" product_encrypt = $" + strconv.Itoa(index+1) + " AND " +
			" product_signature = $" + strconv.Itoa(index+2) + " AND " +
			" license_config_id = $" + strconv.Itoa(index+3) + "), "

		index += 4

		tempQuery, index = ListRangeToInQueryWithStartIndex(paramLen, index)
		query += tempQuery
		query += " ) "

		if i < len(userParam)-1 {
			query += ", "
		}

		params = append(params, userParam[i].ProductKey.String, userParam[i].ProductEncrypt.String,
			userParam[i].ProductSignature.String, userParam[i].LicenseConfigId.Int64,
			userParam[i].ParentCustomerId.Int64, userParam[i].CustomerId.Int64,
			userParam[i].SiteId.Int64, userParam[i].InstallationId.Int64,
			userParam[i].ClientID.String, userParam[i].UniqueId1.String,
			userParam[i].UniqueId2.String, userParam[i].ProductValidFrom.Time,
			userParam[i].ProductValidThru.Time, userParam[i].TotalLicense.Int64,
			userParam[i].TotalActivated.Int64, userParam[i].CreatedClient.String,
			userParam[i].CreatedBy.Int64, userParam[i].CreatedAt.Time,
			userParam[i].UpdatedClient.String, userParam[i].UpdatedBy.Int64,
			userParam[i].UpdatedAt.Time,
		)
	}

	query += " RETURNING id "

	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	tempResult, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultInterface interface{}, errors errorModel.ErrorModel) {
		var id int64
		dbError := rows.Scan(&id)
		if dbError != nil {
			errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}
		resultInterface = id
		return
	})

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(int64))
		}
	}
	return
}

func (input userLicenseDAO) UpdateTotalActivatedUserLicense(db *sql.Tx, userParam repository.UserLicenseModel) (err errorModel.ErrorModel) {
	funcName := "UpdateTotalActivatedUserLicense"

	query := "UPDATE " + input.TableName + " SET " +
		" total_activated = total_activated + 1, " +
		" updated_by = $1, " +
		" updated_client = $2, " +
		" updated_at = $3 " +
		" WHERE " +
		" id = $4"

	param := []interface{}{
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
		userParam.ID.Int64}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	defer stmt.Close()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) UpdateTotalActivatedMovedUserLicense(db *sql.Tx, userParam repository.UserLicenseModel) (err errorModel.ErrorModel) {
	funcName := "UpdateTotalActivatedMovedUserLicense"

	query := "UPDATE " + input.TableName + " SET " +
		" total_activated = $1, " +
		" updated_by = $2, " +
		" updated_client = $3, " +
		" updated_at = $4 " +
		" WHERE " +
		" id = $5"

	param := []interface{}{
		userParam.TotalActivated.Int64,
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

	defer stmt.Close()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) UpdateTotalActivatedPrevUserLicense(db *sql.Tx, userParam repository.UserLicenseModel) (err errorModel.ErrorModel) {
	funcName := "UpdateTotalActivatedPrevUserLicense"

	query := fmt.Sprintf(`
			UPDATE user_license ul 
				SET total_activated = 0,
				updated_by = $1,
				updated_client = $2,
				updated_at = $3
				FROM (
					SELECT ul.id FROM user_license ul 
						INNER JOIN product_license pl on pl.id = ul.product_license_id 
						INNER JOIN license_configuration lc on lc.id = pl.license_config_id
					WHERE lc.id = $4
				) AS sub_query
			WHERE ul.id = sub_query.id`)

	param := []interface{}{
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
		userParam.LicenseConfigId.Int64,
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	defer stmt.Close()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) CheckLicenseNamedUserForUserLicense(db *sql.DB, userLicenseParam in.CheckLicenseNamedUserRequest) (userLicenseOnDb repository.UserLicenseModel, err errorModel.ErrorModel) {
	funcName := "CheckLicenseNamedUserForUserLicense"

	query := "SELECT ul.id, ul.product_valid_from, ul.product_valid_thru, " +
		"ul.parent_customer_id, ul.customer_id, ul.site_id, " +
		"ul.installation_id " +
		"FROM " + input.TableName + " ul " +
		"INNER JOIN " + ProductLicenseDAO.TableName + " pl ON ul.product_license_id = pl.id " +
		"INNER JOIN " + LicenseConfigDAO.TableName + " lc ON lc.id = pl.license_config_id " +
		"WHERE ul.unique_id_1 = $1 AND " +
		"ul.unique_id_2 = $2 AND " +
		"ul.client_id = $3 AND " +
		"lc.client_type_id = $4 AND " +
		"ul.product_valid_from <= CURRENT_DATE AND " +
		"ul.product_valid_thru >= CURRENT_DATE AND " +
		"ul.deleted = FALSE AND " +
		"pl.deleted = FALSE AND " +
		"pl.license_status = 1 " +
		"ORDER BY ul.id ASC LIMIT 1 "

	params := []interface{}{
		userLicenseParam.UniqueId1,
		userLicenseParam.UniqueId2,
		userLicenseParam.ClientId,
		userLicenseParam.ClientTypeID,
	}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&userLicenseOnDb.ID,
		&userLicenseOnDb.ProductValidFrom,
		&userLicenseOnDb.ProductValidThru,
		&userLicenseOnDb.ParentCustomerId,
		&userLicenseOnDb.CustomerId,
		&userLicenseOnDb.SiteId,
		&userLicenseOnDb.InstallationId,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) NewCheckLicenseNamedUser(db *sql.DB, userLicenseParam in.CheckLicenseNamedUserRequest) (userLicenseOnDb []repository.UserLicenseModel, err errorModel.ErrorModel) {
	var (
		funcName   = "NewCheckLicenseNamedUser"
		resultTemp []interface{}
	)

	query := fmt.Sprintf(`SELECT pl.id, pl.product_key, ul.total_license, 
		ul.total_activated, (ul.total_license - ul.total_activated) as quota_license, pl.license_status, 
		ul.parent_customer_id, ul.customer_id, ul.site_id, 
		ul.installation_id, ul.id, ul.product_valid_from, 
		ul.product_valid_thru 
	FROM %s ul
	LEFT JOIN %s pl ON pl.id = ul.product_license_id
	LEFT JOIN %s lc ON lc.id = pl.license_config_id
	LEFT JOIN %s ci ON ci.id = lc.installation_id
	LEFT JOIN %s cm ON cm.id = ci.client_mapping_id
	LEFT JOIN %s pr ON pr.id = ci.product_id
	WHERE 
		pr.client_type_id = $1 AND ul.client_id = $2 AND 
		ul.unique_id_1 = $3 AND ul.deleted = FALSE AND (ul.total_license - ul.total_activated) > 0 AND
		pl.deleted = FALSE AND pl.license_status = $4 AND ul.product_valid_from <= now()::date AND 
		ul.product_valid_thru >= now()::date
	`, input.TableName, ProductLicenseDAO.TableName, LicenseConfigDAO.TableName,
		CustomerInstallationDAO.TableName, ClientMappingDAO.TableName, ProductDAO.TableName)

	params := []interface{}{
		userLicenseParam.ClientTypeID,
		userLicenseParam.ClientId,
		userLicenseParam.UniqueId1,
		constanta.ProductLicenseStatusActive,
	}

	if userLicenseParam.UniqueId2 != "" {
		query += fmt.Sprintf(` AND ul.unique_id_2 = $5 `)
		params = append(params, userLicenseParam.UniqueId2)
	}

	query += fmt.Sprintf(` ORDER BY ul.id ASC `)
	rows, errorDB := db.Query(query, params...)
	if errorDB != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	resultTemp, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResult := range resultTemp {
		userLicenseOnDb = append(userLicenseOnDb, itemResult.(repository.UserLicenseModel))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) resultRowsInput(rows *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
	funcName := "resultRowsInput"
	var errorS error
	var userLicenseOnDb repository.UserLicenseModel

	errorS = rows.Scan(&userLicenseOnDb.ProductLicenseID,
		&userLicenseOnDb.ProductKey,
		&userLicenseOnDb.TotalLicense,
		&userLicenseOnDb.TotalActivated,
		&userLicenseOnDb.QuotaLicense,
		&userLicenseOnDb.LicenseStatus,
		&userLicenseOnDb.ParentCustomerId,
		&userLicenseOnDb.CustomerId,
		&userLicenseOnDb.SiteId,
		&userLicenseOnDb.InstallationId,
		&userLicenseOnDb.ID,
		&userLicenseOnDb.ProductValidFrom,
		&userLicenseOnDb.ProductValidThru)

	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp = userLicenseOnDb
	return
}

func (input userLicenseDAO) LockUserAndCheckIDUserLicense(db *sql.DB, id int64) (idResult int64, err errorModel.ErrorModel) {
	funcName := "LockUserAndCheckIDUserLicense"

	query := "SELECT id FROM " + input.TableName + " WHERE id = $1 FOR UPDATE "

	params := []interface{}{id}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(&idResult)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) GetCountUserLicenseForActivationNamedUser(db *sql.DB, userParam repository.UserLicenseModel) (result int, err errorModel.ErrorModel) {
	funcName := "GetCountUserLicenseForActivationNamedUser"
	var tempResult interface{}
	query := fmt.Sprintf(`SELECT COUNT(ul.id)
	FROM %s ul
	LEFT JOIN %s pl ON pl.id = ul.product_license_id
	LEFT JOIN %s lc ON lc.id = pl.license_config_id
	LEFT JOIN %s ci ON ci.id = lc.installation_id
	LEFT JOIN %s cm ON cm.id = ci.client_mapping_id
	LEFT JOIN %s pr ON pr.id = ci.product_id
	WHERE 
		pr.client_type_id = $1 AND cm.client_id = $2 AND 
		ul.unique_id_1 = $3 AND ul.unique_id_2 = $4 AND 
		ul.deleted = FALSE AND (ul.total_license - ul.total_activated) > 0 AND
		pl.deleted = FALSE AND pl.license_status = $5 AND
		lc.product_valid_from <= now()::date AND lc.product_valid_thru >= now()::date 
	`, input.TableName, ProductLicenseDAO.TableName, LicenseConfigDAO.TableName,
		CustomerInstallationDAO.TableName, ClientMappingDAO.TableName, ProductDAO.TableName)

	param := []interface{}{
		userParam.ClientTypeId.Int64,
		userParam.ClientID.String,
		userParam.UniqueId1.String,
		userParam.UniqueId2.String,
		constanta.ProductLicenseStatusActive,
	}

	row := db.QueryRow(query, param...)
	tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp int
		errorS := rws.Scan(&temp)
		return temp, errorS
	}, input.FileName, funcName)

	if err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(int)
	}

	return
}

func (input userLicenseDAO) GetAvailableUserLicenseForActivateNamedUser(db *sql.DB, userParam repository.UserLicenseModel) (result repository.UserLicenseModel, err errorModel.ErrorModel) {
	funcName := "GetAvailableUserLicenseForActivateNamedUser"
	var tempResult interface{}
	query := fmt.Sprintf(`SELECT 
		ul.id, ul.total_license, ul.total_activated, 
		pr.client_type_id, ul.customer_id, ul.site_id
	FROM %s ul
	LEFT JOIN %s pl ON pl.id = ul.product_license_id
	LEFT JOIN %s lc ON lc.id = pl.license_config_id
	LEFT JOIN %s ci ON ci.id = lc.installation_id
	LEFT JOIN %s cm ON cm.id = ci.client_mapping_id
	LEFT JOIN %s pr ON pr.id = ci.product_id
	WHERE 
		pr.client_type_id = $1 AND cm.client_id = $2 AND 
		ul.unique_id_1 = $3 AND ul.unique_id_2 = $4 AND 
		ul.deleted = FALSE AND (ul.total_license - ul.total_activated) > 0 AND
		pl.deleted = FALSE AND pl.license_status = $5 AND
		lc.product_valid_from <= now()::date AND lc.product_valid_thru >= now()::date 
	ORDER BY ul.created_at ASC
	LIMIT 1
	`, input.TableName, ProductLicenseDAO.TableName, LicenseConfigDAO.TableName,
		CustomerInstallationDAO.TableName, ClientMappingDAO.TableName, ProductDAO.TableName)

	param := []interface{}{
		userParam.ClientTypeId.Int64,
		userParam.ClientID.String,
		userParam.UniqueId1.String,
		userParam.UniqueId2.String,
		constanta.ProductLicenseStatusActive,
	}

	row := db.QueryRow(query, param...)
	tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserLicenseModel
		errorS := rws.Scan(
			&temp.ID, &temp.TotalLicense, &temp.TotalActivated,
			&temp.ClientTypeId, &temp.CustomerId, &temp.SiteId,
		)
		return temp, errorS
	}, input.FileName, funcName)
	if err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.UserLicenseModel)
	}

	return
}

func (input userLicenseDAO) GetClientTypeById(db *sql.DB, userLicense repository.UserLicenseModel) (clientTypeId int64, err errorModel.ErrorModel) {
	funcName := "GetClientTypeById"

	query := "SELECT cm.client_type_id FROM user_license ul " +
		"JOIN " + ProductLicenseDAO.TableName + " pl ON ul.product_license_id = pl.id " +
		"JOIN " + LicenseConfigDAO.TableName + " lc ON pl.license_config_id = lc.id " +
		"JOIN " + CustomerInstallationDAO.TableName + " ci ON lc.installation_id = ci.id " +
		"JOIN " + ClientMappingDAO.TableName + " cm ON ci.id = cm.installation_id " +
		"WHERE ul.id = $1"

	params := []interface{}{
		userLicense.ID.Int64,
	}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(&clientTypeId)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) GetActiveUserLicenseForUpdate(db *sql.Tx, userParam repository.UserLicenseModel) (result repository.UserLicenseModel, err errorModel.ErrorModel) {
	var tempResult interface{}

	funcName := "GetActiveUserLicenseForUpdate"
	query := fmt.Sprintf(`SELECT 
		ul.id, lc.id, ul.installation_id, lc.parent_customer_id, 
		pcu.customer_name as parent_customer_name, lc.customer_id, 
		lc.site_id, cu.customer_name, ul.unique_id_1 as company_id, 
		ul.unique_id_2 as branch_id, p.product_name, ul.total_activated as active_user, 
		ul.total_license,  ul.product_valid_from, ul.product_valid_thru, 
		ul.updated_at, cm.client_id, p.client_type_id
	FROM %s ul
	JOIN %s pl ON pl.id = ul.product_license_id
	JOIN %s lc ON lc.id = pl.license_config_id
	JOIN %s pcu ON pcu.id = lc.parent_customer_id 
	JOIN %s cu ON cu.id = lc.customer_id 
	JOIN %s ci ON ci.id = lc.installation_id
	JOIN %s cm ON cm.id = ci.client_mapping_id
	JOIN %s p ON p.id = ci.product_id 
	WHERE 
		ul.id = $1 AND ul.deleted = FALSE AND
		pl.deleted = FALSE AND pl.license_status = $2 AND
		ul.product_valid_from <= now()::date AND ul.product_valid_thru >= now()::date
	`, input.TableName, ProductLicenseDAO.TableName, LicenseConfigDAO.TableName,
		CustomerDAO.TableName, CustomerDAO.TableName, CustomerInstallationDAO.TableName,
		ClientMappingDAO.TableName, ProductDAO.TableName)

	param := []interface{}{userParam.ID.Int64, constanta.ProductLicenseStatusActive}

	query += " FOR UPDATE "

	rows := db.QueryRow(query, param...)
	tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserLicenseModel
		dbError := rws.Scan(
			&temp.ID,
			&temp.LicenseConfigId,
			&temp.InstallationId,
			&temp.ParentCustomerId,
			&temp.ParentCustomerName,
			&temp.CustomerId,
			&temp.SiteId,
			&temp.CustomerName,
			&temp.UniqueId1,
			&temp.UniqueId2,
			&temp.ProductName,
			&temp.TotalActivated,
			&temp.TotalLicense,
			&temp.ProductValidThru,
			&temp.ProductValidFrom,
			&temp.UpdatedAt,
			&temp.ClientID,
			&temp.ClientTypeId,
		)

		return temp, dbError
	}, input.FileName, funcName)

	if err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.UserLicenseModel)
	}

	return
}

func (input userLicenseDAO) UpdateTotalLicenseUserLicense(db *sql.Tx, userParam repository.UserLicenseModel) (err errorModel.ErrorModel) {
	funcName := "UpdateTotalActivatedUserLicense"

	query := "UPDATE " + input.TableName + " SET " +
		" total_license = $1, " +
		" updated_by = $2, " +
		" updated_client = $3, " +
		" updated_at = $4 " +
		" WHERE " +
		" id = $5"

	param := []interface{}{
		userParam.TotalLicense.Int64,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
		userParam.ID.Int64}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	defer stmt.Close()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) InsertUserLicense(db *sql.Tx, userParam repository.UserLicenseModel) (result int64, err errorModel.ErrorModel) {
	funcName := "InsertUserLicense"

	var params []interface{}

	query := fmt.Sprintf(`INSERT INTO %s 
		(
			product_license_id, parent_customer_id, customer_id, 
			site_id, installation_id, client_id, unique_id_1, 
			unique_id_2, product_valid_from, product_valid_thru, 
			total_license, total_activated, created_client, 
			created_by, created_at, updated_client, updated_by, 
			updated_at
		) 
		VALUES 
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id `, input.TableName)

	params = append(params, userParam.ProductLicenseID.Int64,
		userParam.ParentCustomerId.Int64, userParam.CustomerId.Int64,
		userParam.SiteId.Int64, userParam.InstallationId.Int64,
		userParam.ClientID.String, userParam.UniqueId1.String,
		userParam.UniqueId2.String, userParam.ProductValidFrom.Time,
		userParam.ProductValidThru.Time, userParam.TotalLicense.Int64,
		userParam.TotalActivated.Int64, userParam.CreatedClient.String,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time,
		userParam.UpdatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
	)

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	return
}

func (input userLicenseDAO) GetCustomerForAccountRegistration(db *sql.DB, userParam repository.UserLicenseModel) (output repository.UserLicenseModel, err errorModel.ErrorModel) {
	funcName := "GetCustomerForAccountRegistration"

	query := fmt.Sprintf(`
			SELECT 
				ul.id, ci.parent_customer_id, ci.customer_id
			FROM %s ul 
				JOIN %s pl ON pl.id = ul.product_license_id  
				JOIN %s lc ON lc.id = pl.license_config_id 
				JOIN %s ci ON ci.id = lc.installation_id  
				JOIN %s cm ON cm.id = ci.client_mapping_id  
			WHERE 
				ul.unique_id_1 = $1 AND lc.client_type_id = $2 AND cm.client_id = $3`,
		input.TableName, ProductLicenseDAO.TableName, LicenseConfigDAO.TableName,
		CustomerInstallationDAO.TableName, ClientMappingDAO.TableName)

	param := []interface{}{userParam.UniqueId1.String, userParam.ClientTypeId.Int64, userParam.ClientID.String}

	if !util.IsStringEmpty(userParam.UniqueId2.String) {
		query += " AND ul.unique_id_2 = $4 "
		param = append(param, userParam.UniqueId2.String)
	}

	rows := db.QueryRow(query, param...)

	var tempResult interface{}
	if tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserLicenseModel
		dbError := rws.Scan(
			&temp.ID, &temp.ParentCustomerId, &temp.CustomerId,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		output = tempResult.(repository.UserLicenseModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseDAO) GetUserLicenseForCheckExpiredLicense(db *sql.DB, productLicenseID int64) (result []repository.UserLicenseModel, err errorModel.ErrorModel) {
	var (
		tempResult []interface{}
		query      string
	)

	query = fmt.Sprintf(`SELECT ul.id, ul.updated_at 
		FROM %s urd 
		LEFT JOIN %s ul ON ul.id = urd.user_license_id 
		LEFT JOIN %s pl ON pl.id = ul.product_license_id 
		WHERE 
		pl.id = $1 AND urd.deleted = FALSE AND (ul.id > 0 OR ul.id IS NOT NULL)
		GROUP BY ul.id, ul.updated_at `,
		UserRegistrationDetailDAO.TableName, UserLicenseDAO.TableName, ProductLicenseDAO.TableName)

	param := []interface{}{productLicenseID}
	tempResult, err = GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.UserLicenseModel
		dbErrors := rows.Scan(&temp.ID, &temp.UpdatedAt)
		return temp, dbErrors
	}, param)
	if err.Error != nil {
		return
	}

	for _, item := range tempResult {
		result = append(result, item.(repository.UserLicenseModel))
	}

	return
}

func (input userLicenseDAO) ResetTotalActivatedUserLicense(db *sql.Tx, userParam repository.UserLicenseModel) (err errorModel.ErrorModel) {
	var (
		funcName = "ResetTotalActivatedUserLicense"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET 
		total_activated = $1, updated_by = $2, updated_at = $3, 
		updated_client = $4 
		WHERE 
		id = $5 `,
		input.TableName)

	param := []interface{}{
		userParam.TotalActivated.Int64, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
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
