package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type customerInstallationDAO struct {
	AbstractDAO
}

var CustomerInstallationDAO = customerInstallationDAO{}.New()

func (input customerInstallationDAO) New() (output customerInstallationDAO) {
	output.FileName = "CustomerInstallationDAO.go"
	output.TableName = "customer_installation"
	return
}

func (input customerInstallationDAO) InsertMultiCustomerInstallation(tx *sql.Tx, model repository.CustomerInstallationModel, idxSite int) (idInstallation []int64, err errorModel.ErrorModel) {

	funcName := "InsertMultiCustomerInstallation"
	parameterCustomerInstallation := 17
	jVar := 1

	query := fmt.Sprintf(`INSERT INTO %s
		(
			parent_customer_id, customer_id, site_id,
			installation_number, product_id, unique_id_1,
			unique_id_2, remark, installation_date,
			installation_status, created_by, created_client,
			created_at, updated_by, updated_client,
			updated_at, client_mapping_id
		)
		VALUES `, input.TableName)

	userParam := model.CustomerInstallationData[idxSite].Installation

	query += CreateDollarParamInMultipleRowsDAO(len(userParam), parameterCustomerInstallation, jVar, "id")

	var param []interface{}
	for k := 0; k < len(userParam); k++ {
		param = append(param,
			model.ParentCustomerID.Int64, model.CustomerInstallationData[idxSite].CustomerID.Int64, model.CustomerInstallationData[idxSite].SiteID.Int64,
			k+1, userParam[k].ProductID.Int64, userParam[k].UniqueID1.String,
			userParam[k].UniqueID2.String)

		if !util.IsStringEmpty(userParam[k].Remark.String) {
			param = append(param, userParam[k].Remark.String)
		} else {
			param = append(param, nil)
		}

		param = append(param, userParam[k].InstallationDate.Time, userParam[k].InstallationStatus.String, model.CreatedBy.Int64,
			model.CreatedClient.String, model.CreatedAt.Time, model.UpdatedBy.Int64,
			model.UpdatedClient.String, model.UpdatedAt.Time)
		
		if model.ClientMappingID.Int64 > 0 {
			param = append(param, userParam[k].ClientMappingID.Int64)
		} else {
			param = append(param, nil)
		}
	}

	rows, errorS := tx.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	var result []interface{}
	result, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResult := range result {
		idInstallation = append(idInstallation, itemResult.(int64))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) DeleteCustomerInstallationBySiteID(tx *sql.Tx, userParam repository.CustomerInstallationModel, idxSite int) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteCustomerInstallationBySiteID"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s set deleted = TRUE, 
		updated_by = $1, updated_client = $2, updated_at = $3 
		WHERE site_id = $4 `, input.TableName)

	param := []interface{}{userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
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

func (input customerInstallationDAO) DeleteCustomerInstallationByInstallationID(tx *sql.Tx, header repository.CustomerInstallationModel, indexSite int, indexInstallation int) (err errorModel.ErrorModel) {
	funcName := "DeleteCustomerInstallationByInstallationID"
	userParam := header.CustomerInstallationData[indexSite].Installation[indexInstallation]

	query := fmt.Sprintf(`UPDATE %s
		SET 
			deleted = TRUE, updated_by = $1,
			updated_client = $2, updated_at = $3
		WHERE id = $4 `, input.TableName)

	param := []interface{}{
		header.UpdatedBy.Int64, header.UpdatedClient.String, header.UpdatedAt.Time,
		userParam.InstallationID.Int64}

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

func (input customerInstallationDAO) InsertCustomerInstallation(db *sql.Tx, userParam repository.CustomerInstallationModel, idxSite int, idxInstallation int, queue int64) (id int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertCustomerInstallation"
		query    string
	)

	query = fmt.Sprintf(`INSERT INTO %s
		(
			parent_customer_id, customer_id, site_id,
			installation_number, product_id, unique_id_1,
			unique_id_2, remark, installation_date,
			installation_status, created_by, created_client,
			created_at, updated_by, updated_client,
			updated_at, client_mapping_id
		) VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16, $17
		)
		RETURNING id `, input.TableName)

	param := userParam.CustomerInstallationData[idxSite].Installation[idxInstallation]

	params := []interface{}{
		userParam.ParentCustomerID.Int64, userParam.CustomerInstallationData[idxSite].CustomerID.Int64, userParam.CustomerInstallationData[idxSite].SiteID.Int64,
		queue + 1, param.ProductID.Int64, param.UniqueID1.String,
		param.UniqueID2.String,
	}

	if param.Remark.String != "" {
		params = append(params, param.Remark.String)
	} else {
		params = append(params, nil)
	}

	params = append(params, param.InstallationDate.Time,
		param.InstallationStatus.String, userParam.CreatedBy.Int64, userParam.CreatedClient.String,
		userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time)

	if param.ClientMappingID.Int64 > 0 {
		params = append(params, param.ClientMappingID.Int64)
	} else {
		params = append(params, nil)
	}

	result := db.QueryRow(query, params...)
	errorS := result.Scan(&id)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input customerInstallationDAO) GetCustomerInstallationForUpdate(db *sql.DB, userParam repository.CustomerInstallationModel, indexSite int, indexInstallation int) (result repository.CustomerInstallationDetail, err errorModel.ErrorModel) {
	var (
		funcName = "GetCustomerInstallationForUpdate"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
		ci.id, ci.updated_at, 
		CASE WHEN 
			(
				(SELECT COUNT(id) FROM %s WHERE installation_id = ci.id AND deleted = FALSE) > 0 OR 
				(ci.client_mapping_id > 0 OR ci.client_mapping_id is not null)
			)
		THEN TRUE ELSE FALSE END is_used, 
		ci.unique_id_1, ci.unique_id_2, ci.product_id 
		FROM %s ci 
		WHERE 
		ci.id = $1 AND ci.deleted = FALSE `,
		LicenseConfigDAO.TableName, input.TableName)

	params := []interface{}{userParam.CustomerInstallationData[indexSite].Installation[indexInstallation].InstallationID.Int64}

	query += fmt.Sprintf(` FOR UPDATE `)

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.InstallationID, &result.UpdatedAt, &result.IsUsed,
		&result.UniqueID1, &result.UniqueID2, &result.ProductID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) UpdateCustomerInstallationByInstallationID(tx *sql.Tx, model repository.CustomerInstallationModel, indexSite int, indexInstallation int) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomerInstallationByInstallationID"

	query := fmt.Sprintf(`UPDATE %s
		SET
			product_id = $1, unique_id_1 = $2, unique_id_2 = $3,
			remark = $4, installation_date = $5, installation_status = $6,
			updated_by = $7, updated_client = $8, updated_at = $9
		WHERE id = $10 `, input.TableName)

	param := []interface{}{
		model.CustomerInstallationData[indexSite].Installation[indexInstallation].ProductID.Int64, model.CustomerInstallationData[indexSite].Installation[indexInstallation].UniqueID1.String, model.CustomerInstallationData[indexSite].Installation[indexInstallation].UniqueID2.String,
		model.CustomerInstallationData[indexSite].Installation[indexInstallation].Remark.String, model.CustomerInstallationData[indexSite].Installation[indexInstallation].InstallationDate.Time, model.CustomerInstallationData[indexSite].Installation[indexInstallation].InstallationStatus.String,
		model.UpdatedBy.Int64, model.UpdatedClient.String, model.UpdatedAt.Time,
		model.CustomerInstallationData[indexSite].Installation[indexInstallation].InstallationID.Int64}

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

	defer func() {
		errorS = stmt.Close()
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		}
	}()

	return errorModel.GenerateNonErrorModel()
}

func (input customerInstallationDAO) ViewCustomerInstallation(db *sql.DB, userParam repository.CustomerInstallationData, userParamDetail repository.CustomerInstallationDetail,
	scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (resultDetail []repository.CustomerInstallationDetail, err errorModel.ErrorModel) {

	var (
		funcName        = "ViewCustomerInstallation"
		lastDollarParam int
		query           string
		params          []interface{}
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT
			ci.id, ci.product_id, pd.product_group_id,
			pd.product_id as product_code, pd.product_name, ci.remark,
			ci.unique_id_1, ci.unique_id_2, ci.installation_date,
			ci.installation_status, ci.product_valid_from, ci.product_valid_thru,
			(ci.product_valid_thru - ci.product_valid_from) as day_range, ci.updated_at, pd.product_description,
			pd.client_type_id, ct.parent_client_type_id
		FROM %s ci
		INNER JOIN %s pd ON ci.product_id = pd.id 
		INNER JOIN %s cu ON ci.parent_customer_id = cu.id 
		INNER JOIN %s cuc ON ci.customer_id = cuc.id
		INNER JOIN %s ct ON pd.client_type_id = ct.id
		WHERE 
		ci.parent_customer_id = $1 AND ci.site_id = $2 AND ci.deleted = FALSE AND 
		pd.deleted = FALSE AND cu.deleted = FALSE AND cuc.deleted = FALSE `,
		CustomerInstallationDAO.TableName, ProductDAO.TableName, CustomerDAO.TableName,
		CustomerDAO.TableName, ClientTypeDAO.TableName)

	lastDollarParam = 2
	params = append(params, userParam.CustomerID.Int64, userParam.SiteID.Int64)
	if userParamDetail.ClientTypeID.Int64 > 0 {
		lastDollarParam++
		query += fmt.Sprintf(` AND pd.client_type_id = $%d `, lastDollarParam)
		params = append(params, userParamDetail.ClientTypeID.Int64)
	}

	if userParamDetail.IsLicense.Bool {
		lastDollarParam++
		query += fmt.Sprintf(` AND pd.is_license = $%d `, lastDollarParam)
		params = append(params, userParamDetail.IsLicense.Bool)
	}

	additionalWhere = input.PrepareScopeInCustomerInstallation(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

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
			var ci repository.CustomerInstallationDetail

			dbError := rows.Scan(
				&ci.InstallationID, &ci.ProductID, &ci.ProductGroupID,
				&ci.ProductCode, &ci.ProductName, &ci.Remark,
				&ci.UniqueID1, &ci.UniqueID2, &ci.InstallationDate,
				&ci.InstallationStatus, &ci.ProductValidFrom, &ci.ProductValidThru,
				&ci.DayRange, &ci.UpdatedAt, &ci.ProductDescription,
				&ci.ClientTypeID, &ci.ParentClientTypeID)

			if dbError != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
				return
			}

			resultDetail = append(resultDetail, ci)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) GetInstallationNumberLastInstallation(db *sql.DB, userParam repository.CustomerInstallationModel, indexSite int) (queue int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetInstallationNumberLastInstallation"
		query    string
		results  *sql.Row
		dbError  error
	)

	query = fmt.Sprintf(`SELECT installation_number FROM %s
		WHERE
		parent_customer_id = $1 AND customer_id = $2 AND site_id = $3 AND
		deleted = FALSE ORDER BY id DESC `, input.TableName)

	params := []interface{}{
		userParam.ParentCustomerID.Int64,
		userParam.CustomerInstallationData[indexSite].CustomerID.Int64,
		userParam.CustomerInstallationData[indexSite].SiteID.Int64,
	}

	results = db.QueryRow(query, params...)
	dbError = results.Scan(&queue)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) GetCustomerInstallationByIDJoinClientMappingAndProduct(db *sql.DB, userParam repository.CustomerInstallationForConfig, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (resultOnDB repository.CustomerInstallationForConfig, err errorModel.ErrorModel) {
	var (
		funcName        = "GetCustomerInstallationByIDJoinClientMappingAndProduct"
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT
			ci.parent_customer_id, ci.customer_id, ci.site_id,
			ci.product_id, ci.unique_id_1, ci.unique_id_2,
			cm.client_id, pr.client_type_id, pr.license_variant_id,
			pr.license_type_id, pr.deployment_method, pr.max_offline_days,
			pr.module_id_1, pr.module_id_2, pr.module_id_3,
			pr.module_id_4, pr.module_id_5, pr.module_id_6,
			pr.module_id_7, pr.module_id_8, pr.module_id_9,
			pr.module_id_10, ci.id, pr.is_user_concurrent
		FROM %s ci
		INNER JOIN %s cm ON ci.client_mapping_id = cm.id
		INNER JOIN %s pr ON ci.product_id = pr.id
		INNER JOIN %s cup ON ci.parent_customer_id = cup.id 
		INNER JOIN %s cuc ON ci.customer_id = cuc.id
		WHERE
			ci.id = $1 AND ci.deleted = FALSE AND cm.deleted = FALSE AND
			pr.deleted = FALSE AND cup.deleted = FALSE AND cuc.deleted = FALSE `,
		input.TableName, ClientMappingDAO.TableName, ProductDAO.TableName,
		CustomerDAO.TableName, CustomerDAO.TableName)

	additionalWhere = input.PrepareScopeInCustomerInstallation(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(
		&resultOnDB.ParentCustomerID, &resultOnDB.CustomerID, &resultOnDB.SiteID,
		&resultOnDB.ProductID, &resultOnDB.UniqueID1, &resultOnDB.UniqueID2,
		&resultOnDB.ClientID, &resultOnDB.ClientTypeID, &resultOnDB.LicenseVariantID,
		&resultOnDB.LicenseTypeID, &resultOnDB.DeploymentMethod, &resultOnDB.MaxOfflineDays,
		&resultOnDB.ModuleID1, &resultOnDB.ModuleID2, &resultOnDB.ModuleID3,
		&resultOnDB.ModuleID4, &resultOnDB.ModuleID5, &resultOnDB.ModuleID6,
		&resultOnDB.ModuleID7, &resultOnDB.ModuleID8, &resultOnDB.ModuleID9,
		&resultOnDB.ModuleID10, &resultOnDB.ID, &resultOnDB.IsUserConcurrent)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) ViewDetailInstallationByInstallationID(db *sql.DB, userParam repository.CustomerInstallationDetailConfig, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (resultOnDB repository.CustomerInstallationDetailConfig, err errorModel.ErrorModel) {
	var (
		funcName        = "ViewDetailInstallationByInstallationID"
		query           string
		params          []interface{}
		results         *sql.Row
		dbError         error
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT 
		ci.parent_customer_id, cup.customer_name as parent_customer, ci.customer_id, 
		ci.site_id, cuc.customer_name, cm.client_id, 
		pr.product_name, ct.client_type, lv.license_variant_name, 
		lt.license_type_name, pr.deployment_method, pr.no_of_user, 
		pr.is_user_concurrent, ci.unique_id_1, ci.unique_id_2, 
		ci.product_valid_from, ci.product_valid_thru, pr.max_offline_days, 
		mo1.module_name, mo2.module_name, mo3.module_name, 
		mo4.module_name, mo5.module_name, mo6.module_name, 
		mo7.module_name, mo8.module_name, mo9.module_name, 
		mo10.module_name, ci.product_id, ci.id, 
		ci.client_mapping_id 
		FROM %s ci 
		INNER JOIN %s cup ON ci.parent_customer_id = cup.id 
		INNER JOIN %s cuc ON ci.customer_id = cuc.id 
		LEFT JOIN %s cm ON cm.id = ci.client_mapping_id 
		INNER JOIN %s pr ON ci.product_id = pr.id 
		INNER JOIN %s ct ON ct.id = pr.client_type_id 
		INNER JOIN %s lv ON lv.id = pr.license_variant_id 
		INNER JOIN %s lt ON lt.id = pr.license_type_id 
		LEFT JOIN %s mo1 ON mo1.id = pr.module_id_1 
		LEFT JOIN %s mo2 ON mo2.id = pr.module_id_2 
		LEFT JOIN %s mo3 ON mo3.id = pr.module_id_3 
		LEFT JOIN %s mo4 ON mo4.id = pr.module_id_4 
		LEFT JOIN %s mo5 ON mo5.id = pr.module_id_5 
		LEFT JOIN %s mo6 ON mo6.id = pr.module_id_6 
		LEFT JOIN %s mo7 ON mo7.id = pr.module_id_7 
		LEFT JOIN %s mo8 ON mo8.id = pr.module_id_8 
		LEFT JOIN %s mo9 ON mo9.id = pr.module_id_9 
		LEFT JOIN %s mo10 ON mo10.id = pr.module_id_10 
		WHERE ci.id = $1 AND ci.deleted = FALSE `,
		CustomerInstallationDAO.TableName, CustomerDAO.TableName, CustomerDAO.TableName,
		ClientMappingDAO.TableName, ProductDAO.TableName, ClientTypeDAO.TableName,
		LicenseVariantDAO.TableName, LicenseTypeDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName)

	params = []interface{}{userParam.InstallationID.Int64}

	additionalWhere = input.PrepareScopeInCustomerInstallation(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	results = db.QueryRow(query, params...)
	dbError = results.Scan(
		&resultOnDB.ParentCustomerID, &resultOnDB.ParentCustomer, &resultOnDB.CustomerID,
		&resultOnDB.SiteID, &resultOnDB.Customer, &resultOnDB.ClientID,
		&resultOnDB.ProductName, &resultOnDB.ClientType, &resultOnDB.LicenseVariantName,
		&resultOnDB.LicenseTypeName, &resultOnDB.DeploymentMethod, &resultOnDB.NoOfUser,
		&resultOnDB.IsUserConcurrent, &resultOnDB.UniqueID1, &resultOnDB.UniqueID2,
		&resultOnDB.ProductValidFrom, &resultOnDB.ProductValidThru, &resultOnDB.MaxOfflineDays,
		&resultOnDB.ModuleIDName1, &resultOnDB.ModuleIDName2, &resultOnDB.ModuleIDName3,
		&resultOnDB.ModuleIDName4, &resultOnDB.ModuleIDName5, &resultOnDB.ModuleIDName6,
		&resultOnDB.ModuleIDName7, &resultOnDB.ModuleIDName8, &resultOnDB.ModuleIDName9,
		&resultOnDB.ModuleIDName10, &resultOnDB.ProductID, &resultOnDB.InstallationID,
		&resultOnDB.ClientMappingID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) UpdateCustomerInstallationFromLicenseConfig(tx *sql.Tx, model repository.LicenseConfigModel) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomerInstallationFromLicenseConfig"

	query := fmt.Sprintf(`UPDATE %s
		SET
			product_valid_from = $1, product_valid_thru = $2, updated_by = $3,
			updated_client = $4, updated_at = $5
		WHERE id = $6 `, input.TableName)

	param := []interface{}{
		model.ProductValidFrom.Time, model.ProductValidThru.Time, model.UpdatedBy.Int64,
		model.UpdatedClient.String, model.UpdatedAt.Time, model.InstallationID.Int64,
	}

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

	defer func() {
		errorS = stmt.Close()
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		}
	}()

	return errorModel.GenerateNonErrorModel()
}

func (input customerInstallationDAO) GetCustomerInstallation(db *sql.DB, userParam repository.CustomerInstallationDetail) (isExist bool, result repository.CustomerInstallationDetailConfig, idInstallation string, err errorModel.ErrorModel) {
	var (
		funcName = "GetCustomerInstallation"
		query    string
		params   []interface{}
		results  *sql.Row
		dbError  error
	)

	query = fmt.Sprintf(`SELECT array_agg(ci.id), ci.parent_customer_id, ci.customer_id, 
		ci.site_id, ci.unique_id_1, ci.unique_id_2, 
			(
			SELECT CASE WHEN count(id) > 0 
			THEN TRUE ELSE FALSE END is_exist 
			FROM %s WHERE 
				company_id = ci.unique_id_1 AND 
				branch_id = ci.unique_id_2 AND 
				deleted = false
			) 
		FROM %s ci WHERE 
		ci.unique_id_1 = $1 AND 
		ci.unique_id_2 = $2 AND 
		ci.installation_status = 'A' AND 
		ci.client_mapping_id is null AND 
		ci.deleted = FALSE 
		GROUP BY 
		ci.parent_customer_id, ci.customer_id, ci.site_id, 
		ci.unique_id_1, ci.unique_id_2 `,
		ClientMappingDAO.TableName, CustomerInstallationDAO.TableName)

	params = []interface{}{userParam.UniqueID1.String, userParam.UniqueID2.String}
	results = db.QueryRow(query, params...)
	dbError = results.Scan(
		&idInstallation, &result.ParentCustomerID, &result.CustomerID,
		&result.SiteID, &result.UniqueID1, &result.UniqueID2, &isExist)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) resultRowsInput(rows *sql.Rows) (idTemp interface{}, err errorModel.ErrorModel) {
	funcName := "resultRowsInput"
	var errorS error
	var id int64

	errorS = rows.Scan(&id)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	idTemp = id
	return
}

func (input customerInstallationDAO) GetCustomerInstallationByClientIDAndID(db *sql.DB, userParam repository.CustomerInstallationForConfig) (result repository.CustomerInstallationModel, err errorModel.ErrorModel) {
	var tempResult interface{}

	funcName := "GetCustomerInstallationByClientIDAndID"
	query := fmt.Sprintf(
		`SELECT 
				ci.id, ci.updated_at
		FROM %s ci
		LEFT JOIN %s cm ON ci.client_mapping_id = cm.id
		LEFT JOIN %s p ON p.id = ci.product_id
		WHERE 
			ci.deleted = FALSE AND cm.deleted = FALSE AND
			ci.id = $1 AND p.client_type_id = $2 AND ci.parent_customer_id = $3
		`, input.TableName, ClientMappingDAO.TableName, ProductDAO.TableName)

	param := []interface{}{
		userParam.ID.Int64,
		userParam.ClientTypeID.Int64,
		userParam.ParentCustomerID.Int64,
	}

	rows := db.QueryRow(query, param...)
	if tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerInstallationModel
		dbError := rws.Scan(
			&temp.ID,
			&temp.UpdatedAt,
		)

		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerInstallationModel)
	}

	return
}

func (input customerInstallationDAO) UpdateCustomerInstallationByMultipleInstallationID(tx *sql.Tx, userParam repository.CustomerInstallationModel, clientMappingID int64, idInstallationCol repository.DetailUniqueID) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomerInstallationByMultipleInstallationID"

	query := fmt.Sprintf(`UPDATE %s
		SET
			client_mapping_id = $1, updated_by = $2, updated_client = $3,
			updated_at = $4
		WHERE id IN ( `, input.TableName)

	for i := 0; i < len(idInstallationCol.InstallationIDCol); i++ {
		installationID := strconv.Itoa(int(idInstallationCol.InstallationIDCol[i].InstallationID.Int64))

		if len(idInstallationCol.InstallationIDCol)-(i+1) == 0 {
			query += installationID + ")"
		} else {
			query += installationID + ","
		}
	}

	param := []interface{}{
		clientMappingID, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time}

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

	defer func() {
		errorS = stmt.Close()
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		}
	}()

	return errorModel.GenerateNonErrorModel()
}

func (input customerInstallationDAO) CheckCustomerInstallationInDifferentSite(db *sql.DB, userParam repository.CustomerInstallationModel, idxSite int, idxInstallation int) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName        = "CheckCustomerInstallationInDifferentSite"
		query           string
		params          []interface{}
		results         *sql.Row
		dbError         error
		parameterNumber = 3
	)

	query = fmt.Sprintf(`SELECT (CASE WHEN count(id) > 0 THEN TRUE ELSE FALSE END) is_exist 
		FROM %s 
		WHERE unique_id_1 = $1 AND (parent_customer_id <> $2 OR customer_id <> $3) AND deleted = FALSE `,
		CustomerInstallationDAO.TableName)

	params = []interface{}{userParam.CustomerInstallationData[idxSite].Installation[idxInstallation].UniqueID1.String,
		userParam.ParentCustomerID.Int64,
		userParam.CustomerInstallationData[idxSite].CustomerID.Int64}

	if userParam.CustomerInstallationData[idxSite].SiteID.Int64 > 0 {
		parameterNumber++
		query += " AND site_id <> $" + strconv.Itoa(parameterNumber) + " "
		params = append(params, userParam.CustomerInstallationData[idxSite].SiteID.Int64)
	}

	if userParam.CustomerInstallationData[idxSite].Installation[idxInstallation].UniqueID2.String != "" {
		parameterNumber++
		query += " AND unique_id_2 = $" + strconv.Itoa(parameterNumber) + " "
		params = append(params, userParam.CustomerInstallationData[idxSite].Installation[idxInstallation].UniqueID2.String)
	}

	results = db.QueryRow(query, params...)
	dbError = results.Scan(&isExist)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) InsertCustomerInstallationForTesting(db *sql.Tx, userParam repository.CustomerInstallationModel, idxSite int, idxInstallation int, queue int64, clientMappingId int64) (id int64, err errorModel.ErrorModel) {

	funcName := "InsertCustomerInstallationForTesting"

	query := fmt.Sprintf(`INSERT INTO %s
		( 
			parent_customer_id, customer_id, site_id, installation_number, product_id, unique_id_1,
			unique_id_2, remark, installation_date, installation_status, created_by, created_client,
			created_at, updated_by, updated_client, updated_at, client_mapping_id 
		) VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9,	$10, $11, $12, $13, $14, $15, $16, $17
		)
		RETURNING id `, input.TableName)

	param := userParam.CustomerInstallationData[idxSite].Installation[idxInstallation]

	params := []interface{}{
		userParam.ParentCustomerID.Int64, userParam.CustomerInstallationData[idxSite].CustomerID.Int64, userParam.CustomerInstallationData[idxSite].SiteID.Int64,
		queue + 1, param.ProductID.Int64, param.UniqueID1.String,
		param.UniqueID2.String,
	}

	if param.Remark.String != "" {
		params = append(params, param.Remark.String)
	} else {
		params = append(params, nil)
	}

	params = append(params, param.InstallationDate.Time,
		param.InstallationStatus.String, userParam.CreatedBy.Int64, userParam.CreatedClient.String,
		userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, clientMappingId)

	result := db.QueryRow(query, params...)
	errorS := result.Scan(&id)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input customerInstallationDAO) GetAllParentCustomerInstallation(db *sql.DB, customer repository.CustomerInstallationForConfig, custInstallation repository.CustomerInstallationDetail) (result map[int64]repository.CustomerInstallationDetail, err errorModel.ErrorModel) {
	var (
		funcName   = "GetAllParentCustomerInstallation"
		query      string
		params     []interface{}
		rows       *sql.Rows
		errorS     error
		resultTemp []interface{}
	)

	query = fmt.Sprintf(`select ci.id as installation_id, ci.unique_id_1, ci.unique_id_2, 
		pr.product_id, pr.product_name  
		from %s ci 
		inner join 
			(select p.id as product_id, p.product_name as product_name, ct.id as client_type_id, 
			ct.parent_client_type_id, p.deleted 
			from %s p 
			inner join %s ct on p.client_type_id = ct.id 
			where 
			ct.id = $1 and ct.parent_client_type_id is null and p.deleted = false and 
			ct.deleted = false) as pr on ci.product_id = pr.product_id
		where 
		ci.unique_id_1 = $2 and ci.deleted = false and pr.deleted = false and 
		ci.parent_customer_id = $3 and ci.customer_id = $4 and ci.site_id = $5 `,
		input.TableName, ProductDAO.TableName, ClientTypeDAO.TableName)

	params = []interface{}{
		custInstallation.ParentClientTypeID.Int64, custInstallation.UniqueID1.String, customer.ParentCustomerID.Int64,
		customer.CustomerID.Int64, customer.SiteID.Int64}

	if custInstallation.UniqueID2.String != "" {
		query += fmt.Sprintf(` and ci.unique_id_2 = $6 `)
		params = append(params, custInstallation.UniqueID2.String)
	}

	rows, errorS = db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			errs        error
			resultCatch repository.CustomerInstallationDetail
		)

		errs = rws.Scan(&resultCatch.InstallationID, &resultCatch.UniqueID1, &resultCatch.UniqueID2, &resultCatch.ProductID, &resultCatch.ProductName)
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.TableName, funcName, errs)
			return
		}

		resultTemp = resultCatch
		return
	})

	if err.Error != nil {
		return
	}

	result = make(map[int64]repository.CustomerInstallationDetail)
	for _, itemResultTemp := range resultTemp {
		t := itemResultTemp.(repository.CustomerInstallationDetail)
		result[t.InstallationID.Int64] = t
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) GetInstallationChild(db *sql.DB, custInstallation repository.CustomerInstallationForConfig) (result map[int64]repository.CustomerInstallationDetail, err errorModel.ErrorModel) {
	var (
		funcName   = "GetInstallationChild"
		query      string
		params     []interface{}
		rows       *sql.Rows
		errorS     error
		resultTemp []interface{}
	)

	query = fmt.Sprintf(`select ci.id, ci.product_id, ci.unique_id_1, ci.unique_id_2 
		from %s ci 
		inner join %s p on ci.product_id = p.id 
		where 
		ci.unique_id_1 = $1 and ci.site_id = $2 and
		p.client_type_id in (select id from %s where 
			parent_client_type_id = (select client_type_id from %s where id = $3 and deleted = false) and deleted = false) and 
		ci.deleted = false and p.deleted = false `,
		input.TableName, ProductDAO.TableName, ClientTypeDAO.TableName,
		ProductDAO.TableName)

	params = []interface{}{custInstallation.UniqueID1.String, custInstallation.SiteID.Int64, custInstallation.ProductID.Int64}

	if custInstallation.UniqueID2.String != "" {
		query += fmt.Sprintf(` and ci.unique_id_2 = $4 `)
		params = append(params, custInstallation.UniqueID2.String)
	}

	rows, errorS = db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			errs        error
			resultCatch repository.CustomerInstallationDetail
		)

		errs = rws.Scan(&resultCatch.InstallationID, &resultCatch.ProductID, &resultCatch.UniqueID1, &resultCatch.UniqueID2)
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.TableName, funcName, errs)
			return
		}

		resultTemp = resultCatch
		return
	})

	if err.Error != nil {
		return
	}

	result = make(map[int64]repository.CustomerInstallationDetail)
	for _, itemResultTemp := range resultTemp {
		t := itemResultTemp.(repository.CustomerInstallationDetail)
		result[t.InstallationID.Int64] = t
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) CheckOtherParentInstallation(db *sql.DB, custInstallation repository.CustomerInstallationForConfig, parentInfo repository.ProductModel) (result map[int64]repository.CustomerInstallationDetail, err errorModel.ErrorModel) {
	var (
		funcName   = "CheckOtherParentInstallation"
		query      string
		params     []interface{}
		rows       *sql.Rows
		errorS     error
		resultTemp []interface{}
	)

	query = fmt.Sprintf(`select ci.id, ci.product_id, ci.unique_id_1, 
		ci.unique_id_2 from %s ci 
		inner join %s pr on ci.product_id = pr.id 
		where 
		ci.unique_id_1 = $1 and (ci.product_id = $2 or pr.client_type_id = $3) and ci.site_id = $4 and 
		ci.id <> $5 and ci.deleted = false and pr.deleted = false `, input.TableName, ProductDAO.TableName)

	params = []interface{}{custInstallation.UniqueID1.String, custInstallation.ProductID.Int64, parentInfo.ParentClientTypeID.Int64,
		custInstallation.SiteID.Int64, custInstallation.ID.Int64}

	if custInstallation.UniqueID2.String != "" {
		query += fmt.Sprintf(` and ci.unique_id_2 = $6 `)
		params = append(params, custInstallation.UniqueID2.String)
	}

	rows, errorS = db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			errs error
			temp repository.CustomerInstallationDetail
		)

		errs = rws.Scan(&temp.InstallationID, &temp.ProductID, &temp.UniqueID1, &temp.UniqueID2)
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.TableName, funcName, errs)
			return
		}

		resultTemp = temp
		return
	})

	if err.Error != nil {
		return
	}

	result = make(map[int64]repository.CustomerInstallationDetail)
	for _, itemResultTemp := range resultTemp {
		t := itemResultTemp.(repository.CustomerInstallationDetail)
		result[t.InstallationID.Int64] = t
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) PrepareScopeInCustomerInstallation(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int) (additionalWhere []string) {
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
		case constanta.ProductGroupDataScope:
			keyDataScope = constanta.ProductGroupDataScope
		case constanta.ClientTypeDataScope:
			keyDataScope = constanta.ClientTypeDataScope
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

func (input customerInstallationDAO) GetCustomerInstallationBySingleUniqueID(db *sql.DB, custInstallation repository.CustomerInstallationDetail, isNullCMID bool) (result repository.CustomerInstallationForConfig, err errorModel.ErrorModel) {
	var (
		funcName = "GetCustomerInstallationByUniqueID"
		query    string
		params   []interface{}
		row      *sql.Row
		errorS   error
	)

	query = fmt.Sprintf(`SELECT 
		id, unique_id_1, unique_id_2, 
		'$1' as branch_name
		FROM %s 
		WHERE 
		unique_id_1 = $2 AND unique_id_2 = $3 `,
		input.TableName)

	if isNullCMID {
		query += " AND nullif(client_mapping_id, 0) IS NULL "
	}

	params = []interface{}{custInstallation.UniqueID1.String, custInstallation.UniqueID2.String}
	row = db.QueryRow(query, params...)
	errorS = row.Scan(&result.ID, &result.UniqueID1, &result.UniqueID2,
		&result.BranchName)
	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) GetCustomerInstallationByUniqueID(db *sql.DB, custInstallation []repository.CustomerInstallationDetail, isNullCMID bool) (result []repository.CustomerInstallationForConfig, err errorModel.ErrorModel) {
	var (
		funcName         = "GetCustomerInstallationByUniqueID"
		query, tempQuery string
		index            = 1
		params           []interface{}
		rows             *sql.Rows
		errorS           error
		resultTemp       []interface{}
	)

	query = fmt.Sprintf(`SELECT 
		id, unique_id_1, unique_id_2
		FROM %s 
		WHERE 
		(unique_id_1, unique_id_2) IN ( `,
		input.TableName)

	for i, detail := range custInstallation {
		query += " ( "

		tempQuery, index = ListRangeToInQueryWithStartIndex(2, index)

		query += tempQuery + " ) "

		if i < len(custInstallation)-1 {
			query += fmt.Sprintf(", \n")
		}

		params = append(params, detail.UniqueID1.String, detail.UniqueID2.String)
	}

	query += " ) "

	if isNullCMID {
		query += " AND nullif(client_mapping_id, 0) IS NULL "
	}

	rows, errorS = db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (output interface{}, err errorModel.ErrorModel) {
		var (
			errs error
			temp repository.CustomerInstallationForConfig
		)

		errs = rws.Scan(&temp.ID, &temp.UniqueID1, &temp.UniqueID2)
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.TableName, funcName, errs)
			return
		}

		output = temp
		return
	})

	if err.Error != nil {
		return
	}

	if len(resultTemp) != 0 {
		for _, item := range resultTemp {
			result = append(result, item.(repository.CustomerInstallationForConfig))
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationDAO) UpdateCustomerInstallationClientMappingID(tx *sql.Tx, model repository.CustomerInstallationModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateCustomerInstallationClientMappingID"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET 
		client_mapping_id = $1, updated_by = $2, updated_client = $3, 
		updated_at = $4 
		WHERE id = $5 `,
		input.TableName)

	param := []interface{}{
		model.ClientMappingID.Int64, model.UpdatedBy.Int64,
		model.UpdatedClient.String, model.UpdatedAt.Time, model.ID.Int64,
	}

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

	defer func() {
		errorS = stmt.Close()
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		}
	}()

	return errorModel.GenerateNonErrorModel()
}

func (input customerInstallationDAO) GetClientMappingIDCustomerInstallationByUniqueID(db *sql.DB, custInstallation repository.CustomerInstallationDetail) (isClientExist bool, clientMappingID int64, err errorModel.ErrorModel) {
	var (
		fileName = "CustomerInstallationDAO.go"
		funcName = "GetClientMappingIDCustomerInstallationByUniqueID"
		query    string
		params   []interface{}
	)
	query = fmt.Sprintf(`select distinct(unique_id_1, unique_id_2) as unique_id, 
		case when (client_mapping_id < 1 or client_mapping_id is null)	
			then false 
			else true end is_client_exist, 
		client_mapping_id
		from %s 
		where unique_id_1 = $1 `,
		input.TableName)
	var (
		uniqueRow           string
		clientMappingIDTemp sql.NullInt64
	)
	params = []interface{}{custInstallation.UniqueID1.String}
	if custInstallation.UniqueID2.String != "" {
		query += fmt.Sprintf(` and unique_id_2 = $2 `)
		params = append(params, custInstallation.UniqueID2.String)
	}
	dbQuery := db.QueryRow(query, params...)
	dbErrs := dbQuery.Scan(&uniqueRow, &isClientExist, &clientMappingIDTemp)
	if dbErrs != nil && dbErrs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, dbErrs)
		return
	}
	clientMappingID = clientMappingIDTemp.Int64
	err = errorModel.GenerateNonErrorModel()
	return
}
