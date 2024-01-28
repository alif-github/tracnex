package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
)

type customerSiteDAO struct {
	AbstractDAO
}

var CustomerSiteDAO = customerSiteDAO{}.New()

func (input customerSiteDAO) New() (output customerSiteDAO) {
	output.FileName = "CustomerSiteDAO.go"
	output.TableName = "customer_site"
	return
}

func (input customerSiteDAO) InsertCustomerSite(db *sql.Tx, userParam repository.CustomerInstallationModel, idxSite int) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertCustomerSite"

	query := "INSERT INTO " + input.TableName + " " +
		"(parent_customer_id, customer_id, created_by, " +
		"created_at, created_client, updated_by, " +
		"updated_at, updated_client) " +
		"VALUES " +
		"($1, $2, $3, " +
		"$4, $5, $6, " +
		"$7, $8) " +
		"RETURNING id "

	params := []interface{}{
		userParam.ParentCustomerID.Int64, userParam.CustomerInstallationData[idxSite].CustomerID.Int64, userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time, userParam.CreatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != constanta.NoRowsInDB {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input customerSiteDAO) CheckCustomerSiteIsExist(db *sql.DB, userParam repository.CustomerInstallationModel, idxSite int, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (idSiteReturn int64, isUsed bool, idInstallation string, err errorModel.ErrorModel) {
	var (
		funcName        = "CheckCustomerSiteIsExist"
		query           string
		params          []interface{}
		results         *sql.Row
		dbError         error
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT 
		cs.id, 
			CASE WHEN (SELECT count(id) FROM %s WHERE site_id = cs.id AND deleted = FALSE) > 0 THEN TRUE ELSE FALSE END is_used,
		array_agg(ci.id) 
		FROM %s cs 
		LEFT JOIN %s ci ON ci.site_id = cs.id 
		INNER JOIN %s cu ON cs.parent_customer_id = cu.id  
		INNER JOIN %s cuc ON cs.customer_id = cuc.id 
		WHERE cs.id = $1 AND cs.deleted = FALSE AND ci.deleted = FALSE AND 
		cu.deleted = FALSE AND cuc.deleted = FALSE `,
		LicenseConfigDAO.TableName, input.TableName, CustomerInstallationDAO.TableName,
		CustomerDAO.TableName, CustomerDAO.TableName)

	additionalWhere = input.PrepareScopeInCustomerSite(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	query += fmt.Sprintf(` GROUP BY ci.site_id, cs.id `)
	params = []interface{}{userParam.CustomerInstallationData[idxSite].SiteID.Int64}
	results = db.QueryRow(query, params...)

	dbError = results.Scan(&idSiteReturn, &isUsed, &idInstallation)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerSiteDAO) DeleteCustomerSite(tx *sql.Tx, userParam repository.CustomerInstallationModel, idxSite int) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteCustomerSite"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s set deleted = TRUE, 
		updated_by = $1, updated_client = $2, updated_at = $3 
		WHERE id = $4 `, input.TableName)

	param := []interface{}{
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.CustomerInstallationData[idxSite].SiteID.Int64}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)

	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	defer stmt.Close()
	return errorModel.GenerateNonErrorModel()
}

func (input customerSiteDAO) GetCustomerSiteForUpdate(db *sql.DB, userParam repository.CustomerInstallationModel, indexSite int) (resultOnDB repository.CustomerSiteModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetCustomerSiteForUpdate"
		query    string
	)

	query = fmt.Sprintf(`SELECT id, updated_at FROM %s WHERE id = $1 AND deleted = FALSE `, input.TableName)

	params := []interface{}{userParam.CustomerInstallationData[indexSite].SiteID.Int64}
	query += fmt.Sprintf(` FOR UPDATE `)
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&resultOnDB.ID.Int64, &resultOnDB.UpdatedAt.Time)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerSiteDAO) ViewDetailCustomerSite(db *sql.DB, userParam repository.CustomerInstallationModel, page int, limit int, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (resultDetail []out.CustomerSite, err errorModel.ErrorModel) {
	var (
		funcName        = "ViewDetailCustomerSite"
		query           string
		params          []interface{}
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT json_build_object(
			'site_id', cs.id, 
			'customer_id', cs.customer_id, 
			'customer_site_name', (SELECT customer_name FROM %s WHERE id = cs.customer_id), 
			'address', (SELECT address FROM %s WHERE id = cs.customer_id), 
			'district', (SELECT name FROM %s INNER JOIN %s ON district.id = customer.district_id WHERE customer.id = cs.customer_id), 
			'province', (SELECT name FROM %s INNER JOIN %s ON province.id = customer.province_id WHERE customer.id = cs.customer_id), 
			'phone', (SELECT phone FROM %s WHERE id = cs.customer_id), 
			'updated_at', cs.updated_at) 
		FROM %s cu 
		INNER JOIN %s cs ON cs.parent_customer_id = cu.id 
		INNER JOIN %s cuc ON cs.customer_id = cuc.id 
		WHERE cu.id = $1 AND cu.deleted = FALSE AND cs.deleted = FALSE `,
		CustomerDAO.TableName, CustomerDAO.TableName, CustomerDAO.TableName,
		DistrictDAO.TableName, CustomerDAO.TableName, ProvinceDAO.TableName,
		CustomerDAO.TableName, CustomerDAO.TableName, CustomerSiteDAO.TableName,
		CustomerDAO.TableName)

	additionalWhere = input.PrepareScopeInCustomerSite(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	query += fmt.Sprintf(` LIMIT $2 OFFSET $3 `)
	params = append(params, userParam.ParentCustomerID.Int64, limit, CountOffset(page, limit))
	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		return resultDetail, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			var (
				dataDetail string
				jsonTemp   out.CustomerSite
			)

			dbError := rows.Scan(&dataDetail)
			if dbError != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
				return
			}

			dbError = json.Unmarshal([]byte(dataDetail), &jsonTemp)
			if dbError != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
				return
			}

			resultDetail = append(resultDetail, jsonTemp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerSiteDAO) CheckCustomerSiteIsExistByParentIDAndCustomerID(db *sql.DB, userParam repository.CustomerInstallationData, parentCustomerID int64) (isExist bool, err errorModel.ErrorModel) {
	funcName := "CheckCustomerSiteIsExistByParentIDAndCustomerID"

	query := "SELECT " +
		"(CASE WHEN count(id) > 0 THEN TRUE ELSE FALSE END) is_exist " +
		"FROM " + input.TableName + " " +
		"WHERE " +
		"parent_customer_id = $1 AND " +
		"customer_id = $2 AND " +
		"deleted = FALSE "

	params := []interface{}{parentCustomerID, userParam.CustomerID.Int64}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&isExist)

	if dbError != nil && dbError.Error() != constanta.NoRowsInDB {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerSiteDAO) CheckCustomerSiteOnly(db *sql.DB, userParam repository.CustomerInstallationModel, idxSite int, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName = "CheckCustomerSiteOnly"
		query    string
		params   []interface{}
		results  *sql.Row
		dbError  error
	)

	query = fmt.Sprintf(`SELECT (CASE WHEN count(cs.id) > 0 THEN TRUE ELSE FALSE END) is_exist 
		FROM %s cs 
		INNER JOIN %s cu ON cs.parent_customer_id = cu.id
		INNER JOIN %s cuc ON cs.customer_id = cuc.id
		WHERE cs.id = $1 AND cs.deleted = FALSE AND cu.deleted = FALSE AND 
		cuc.deleted = FALSE `,
		input.TableName, CustomerDAO.TableName, CustomerDAO.TableName)

	additionalWhere := input.PrepareScopeInCustomerSite(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	params = []interface{}{userParam.CustomerInstallationData[idxSite].SiteID.Int64}
	results = db.QueryRow(query, params...)
	dbError = results.Scan(&isExist)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerSiteDAO) PrepareScopeInCustomerSite(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int) (additionalWhere []string) {
	var (
		parent = "_parent"
		child  = "_child"
	)

	for key := range scopeLimit {
		var (
			keyDataScope   string
			additionalTemp string
		)

		switch key {
		case constanta.ProvinceDataScope + parent:
			keyDataScope = constanta.ProvinceDataScope + parent
		case constanta.DistrictDataScope + parent:
			keyDataScope = constanta.DistrictDataScope + parent
		case constanta.SalesmanDataScope + parent:
			keyDataScope = constanta.SalesmanDataScope + parent
		case constanta.CustomerGroupDataScope + parent:
			keyDataScope = constanta.CustomerGroupDataScope + parent
		case constanta.CustomerCategoryDataScope + parent:
			keyDataScope = constanta.CustomerCategoryDataScope + parent
		case constanta.ProvinceDataScope + child:
			keyDataScope = constanta.ProvinceDataScope + child
		case constanta.DistrictDataScope + child:
			keyDataScope = constanta.DistrictDataScope + child
		case constanta.SalesmanDataScope + child:
			keyDataScope = constanta.SalesmanDataScope + child
		case constanta.CustomerGroupDataScope + child:
			keyDataScope = constanta.CustomerGroupDataScope + child
		case constanta.CustomerCategoryDataScope + child:
			keyDataScope = constanta.CustomerCategoryDataScope + child
		default:
			//--- No Key Limit Found
		}

		if keyDataScope != "" {
			PrepareScopeOnDAO(scopeLimit, scopeDB, &additionalTemp, idxStart, keyDataScope, true)
			if additionalTemp != "" {
				additionalWhere = append(additionalWhere, additionalTemp)
			}
		}
	}

	return
}

func (input customerSiteDAO) GetActiveCustomerSite(db *sql.DB) (output int64, err errorModel.ErrorModel) {
	funcName := "GetActiveCustomerSite"
	query := fmt.Sprintf(
		`SELECT COUNT(DISTINCT(cs.id)) 
	FROM %s cs 
	INNER JOIN %s lc ON cs.id = lc.site_id
	INNER JOIN %s pl ON lc.id = pl.license_config_id 
	WHERE pl.license_status = $1 AND cs.deleted = FALSE `,
		input.TableName, LicenseConfigDAO.TableName, ProductLicenseDAO.TableName)

	results := db.QueryRow(query, constanta.ProductLicenseStatusActive)
	dbError := results.Scan(&output)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}
