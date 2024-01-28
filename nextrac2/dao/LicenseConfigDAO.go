package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type licenseConfigDAO struct {
	AbstractDAO
}

var LicenseConfigDAO = licenseConfigDAO{}.New()

func (input licenseConfigDAO) New() (output licenseConfigDAO) {
	output.FileName = "LicenseConfigDAO.go"
	output.TableName = "license_configuration"
	return
}

func (input licenseConfigDAO) IsLicenseConfigExist(db *sql.DB, userParam repository.LicenseConfigModel) (result bool, err errorModel.ErrorModel) {
	funcName := "IsLicenseConfigExist"
	var tempResult interface{}

	query := fmt.Sprintf(`SELECT CASE WHEN COUNT(id) > 0 THEN TRUE ELSE FALSE END is_exist FROM %s WHERE deleted = FALSE AND id = $1 `, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	results := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var resultTemp bool
		dbError := results.Scan(&resultTemp)
		return resultTemp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(bool)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) CheckPreviousLicenseConfig(db *sql.DB, userParam repository.LicenseConfigModel) (result repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var (
		funcName   = "CheckPreviousLicenseConfig"
		tempResult interface{}
	)

	query := fmt.Sprintf(`
		SELECT
			case when count(urd.id) > 0 then true else false end, lc.old_license_configuration_id
		FROM %s lc 
			left join %s pl ON pl.license_config_id = lc.old_license_configuration_id 
			left join %s ul on ul.product_license_id = pl.id
			left join %s urd on urd.user_license_id = ul.id  
		where lc.id = $1
		group by lc.id
	`, input.TableName, ProductLicenseDAO.TableName, UserLicenseDAO.TableName, UserRegistrationDetailDAO.TableName)

	param := []interface{}{userParam.ID.Int64}

	results := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var resultTemp repository.LicenseConfigModel
		dbError := results.Scan(&resultTemp.IsHasPrevLicenseConfig, &resultTemp.PrevLicenseConfigID)
		return resultTemp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.LicenseConfigModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) InsertLicenseConfig(db *sql.Tx, userParam repository.LicenseConfigModel) (id int64, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertLicenseConfig"
		paramAmount = 34
		startIndex  = 1
	)

	query := fmt.Sprintf(`INSERT INTO %s 
			(installation_id, parent_customer_id, customer_id, 
			site_id, client_id, product_id, 
			client_type_id, license_variant_id, license_type_id, 
			deployment_method, no_of_user, is_user_concurrent, 
			max_offline_days, unique_id_1, unique_id_2, 
			product_valid_from, product_valid_thru, module_id_1, 
			module_id_2, module_id_3, module_id_4, 
			module_id_5, module_id_6, module_id_7, 
			module_id_8, module_id_9, module_id_10, 
			created_by, created_client, created_at, 
			updated_by, updated_client, updated_at, 
			allow_activation) VALUES `, input.TableName)

	query += CreateDollarParamInMultipleRowsDAO(1, paramAmount, startIndex, "id")

	params := []interface{}{
		userParam.InstallationID.Int64, userParam.ParentCustomerID.Int64, userParam.CustomerID.Int64,
		userParam.SiteID.Int64, userParam.ClientID.String, userParam.ProductID.Int64,
		userParam.ClientTypeID.Int64, userParam.LicenseVariantID.Int64, userParam.LicenseTypeID.Int64,
		userParam.DeploymentMethod.String, userParam.NoOfUser.Int64, userParam.IsUserConcurrent.String,
		userParam.MaxOfflineDays.Int64, userParam.UniqueID1.String, userParam.UniqueID2.String,
		userParam.ProductValidFrom.Time, userParam.ProductValidThru.Time,
	}

	input.addModuleFilled(&params, userParam)
	params = append(params,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		"N")

	result := db.QueryRow(query, params...)
	var tempResult interface{}
	if tempResult, err = RowCatchResult(result, func(rws *sql.Row) (interface{}, error) {
		var temp int64
		errorS := result.Scan(&temp)
		return temp, errorS
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		id = tempResult.(int64)
	}

	return
}

func (input licenseConfigDAO) GetLicenseConfigForDelete(db *sql.DB, licenseConfigModel repository.LicenseConfigModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var (
		funcName        = "GetLicenseConfigForDelete"
		query           string
		tempResult      interface{}
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT lc.id, lc.updated_at, lc.created_by, 
			(SELECT CASE WHEN count(id) > 0 THEN TRUE ELSE FALSE END 
			FROM %s WHERE license_config_id = lc.id) isUsed 
		FROM %s as lc 
		INNER JOIN %s cup ON lc.parent_customer_id = cup.id
		INNER JOIN %s cuc ON lc.customer_id = cuc.id 
		INNER JOIN %s pr ON lc.product_id = pr.id 
		WHERE lc.id = $1 AND lc.deleted = FALSE AND cup.deleted = FALSE AND 
		cuc.deleted = FALSE `,
		ProductLicenseDAO.TableName, input.TableName, CustomerDAO.TableName,
		CustomerDAO.TableName, ProductDAO.TableName)

	params := []interface{}{licenseConfigModel.ID.Int64}

	if licenseConfigModel.CreatedBy.Int64 > 0 {
		query += fmt.Sprintf(` AND lc.created_by = $2 `)
		params = append(params, licenseConfigModel.CreatedBy.Int64)
	}

	additionalWhere = input.PrepareScopeInLicenseConfig(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	query += fmt.Sprintf(` FOR UPDATE `)

	results := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var resultTemp repository.LicenseConfigModel
		dbError := results.Scan(&resultTemp.ID, &resultTemp.UpdatedAt, &resultTemp.CreatedBy, &resultTemp.IsUsed)
		return resultTemp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.LicenseConfigModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) DeleteLicenseConfig(db *sql.Tx, userParam repository.LicenseConfigModel) (err errorModel.ErrorModel) {
	funcName := "DeleteLicenseConfig"

	query := fmt.Sprintf(`UPDATE %s SET 
			deleted = TRUE, updated_by = $1, updated_client = $2, updated_at = $3 
			WHERE id = $4 `, input.TableName)

	param := []interface{}{
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
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

	defer func() {
		errs = stmt.Close()
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		}
	}()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) GetLicenseConfigForUpdate(db *sql.DB, licenseConfigModel repository.LicenseConfigModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var (
		funcName        = "GetLicenseConfigForUpdate"
		tempResult      interface{}
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT lc.id, lc.updated_at, lc.created_by, 
		lc.allow_activation 
		FROM %s lc
		INNER JOIN %s cup ON lc.parent_customer_id = cup.id
		INNER JOIN %s cuc ON lc.customer_id = cuc.id 
		INNER JOIN %s pr ON lc.product_id = pr.id
		WHERE lc.id = $1 AND lc.deleted = FALSE AND cup.deleted = FALSE AND 
		cuc.deleted = FALSE `,
		input.TableName, CustomerDAO.TableName, CustomerDAO.TableName,
		ProductDAO.TableName)

	params := []interface{}{licenseConfigModel.ID.Int64}

	if licenseConfigModel.CreatedBy.Int64 > 0 {
		query += fmt.Sprintf(` AND created_by = $2 `)
		params = append(params, licenseConfigModel.CreatedBy.Int64)
	}

	additionalWhere = input.PrepareScopeInLicenseConfig(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	query += fmt.Sprintf(` FOR UPDATE `)

	results := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var resultTemp repository.LicenseConfigModel
		dbError := results.Scan(&resultTemp.ID, &resultTemp.UpdatedAt, &resultTemp.CreatedBy,
			&resultTemp.AllowActivation)
		return resultTemp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.LicenseConfigModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) ActivatingLicenseConfig(db *sql.Tx, userParam repository.LicenseConfigModel) (err errorModel.ErrorModel) {
	funcName := "ActivatingLicenseConfig"

	query := fmt.Sprintf(`UPDATE %s SET 
			allow_activation = $1, updated_by = $2, updated_client = $3, 
			updated_at = $4 WHERE id = $5 `, input.TableName)

	param := []interface{}{
		userParam.AllowActivation.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
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

	defer func() {
		errs = stmt.Close()
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		}
	}()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) GetCountLicenseConfig(db *sql.DB, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		query              string
		colAdditionalWhere []string
	)

	query = fmt.Sprintf(`SELECT COUNT(lc.id) FROM %s lc `, input.TableName)

	colAdditionalWhere = input.setScopeData(scopeLimit, scopeDB, false)
	input.setSearchByLicenseConfig(&searchBy, &in.GetListDataDTO{})
	input.setCreatedByLicenseConfig(createdBy, &searchBy)
	getListData := getListJoinDataDAO{Table: "lc", Query: query, AdditionalWhere: colAdditionalWhere}
	input.setGetListJoinLicenseConfig(&getListData)

	return getListData.GetCountJoinData(db, searchBy, 0)
}

func (input licenseConfigDAO) GetListLicenseConfig(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query              string
		colAdditionalWhere []string
	)

	query = fmt.Sprintf(`SELECT 
			lc.id as id, cuc.customer_name as customer_name, lc.unique_id_1 as unique_id_1, 
			lc.unique_id_2 as unique_id_2, lc.installation_id as installation_id, pr.product_name as product_name, 
			lv.license_variant_name as license_variant_name, lt.license_type_name as license_type_name, lc.product_valid_from as product_valid_from, 
			lc.product_valid_thru as product_valid_thru, lc.allow_activation as allow_activation, lc.updated_at, 
			lc.client_type_id, 
			CASE WHEN lc.allow_activation = 'Y' 
				THEN 'Paid' 
				ELSE 'Unpaid' 
			END payment_status, 
			CASE WHEN lc.allow_activation = 'Y'
				THEN true
				ELSE false
			END is_extend_checklist
			FROM %s lc `, input.TableName)

	colAdditionalWhere = input.setScopeData(scopeLimit, scopeDB, true)
	input.setSearchByLicenseConfig(&searchBy, &userParam)
	input.setCreatedByLicenseConfig(createdBy, &searchBy)

	getListData := getListJoinDataDAO{Table: "lc", Query: query, AdditionalWhere: colAdditionalWhere}
	input.setGetListJoinLicenseConfig(&getListData)

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.LicenseConfigModel

		dbError := rows.Scan(
			&resultTemp.ID, &resultTemp.Customer, &resultTemp.UniqueID1,
			&resultTemp.UniqueID2, &resultTemp.InstallationID, &resultTemp.ProductName,
			&resultTemp.LicenseVariantName, &resultTemp.LicenseTypeName, &resultTemp.ProductValidFrom,
			&resultTemp.ProductValidThru, &resultTemp.AllowActivation, &resultTemp.UpdatedAt,
			&resultTemp.ClientTypeID, &resultTemp.PaymentStatus, &resultTemp.IsExtendChecklist)

		return resultTemp, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input licenseConfigDAO) SelectAllLicenseConfigGetID(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query              string
		colAdditionalWhere []string
		getListData        getListJoinDataDAO
	)

	query = fmt.Sprintf(`SELECT lc.id FROM %s lc `, input.TableName)
	colAdditionalWhere = input.setScopeData(scopeLimit, scopeDB, true)
	input.setSearchByLicenseConfig(&searchBy, &userParam)
	input.setCreatedByLicenseConfig(createdBy, &searchBy)
	searchBy = append(searchBy, in.SearchByParam{
		SearchKey:      "lc.allow_activation",
		DataType:       "char",
		SearchType:     "filter",
		SearchOperator: "eq",
		SearchValue:    "Y",
	})

	getListData = getListJoinDataDAO{Table: "lc", Query: query, AdditionalWhere: colAdditionalWhere}
	input.setGetListJoinLicenseConfig(&getListData)

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var (
			id      int64
			dbError error
		)

		dbError = rows.Scan(&id)
		return id, dbError
	}

	return getListData.GetListJoinDataWithoutPagination(db, userParam, searchBy, 0, mappingFunc)
}

func (input licenseConfigDAO) setSearchByLicenseConfig(searchBy *[]in.SearchByParam, userParam *in.GetListDataDTO) {
	var (
		aliasLicenseConfig  = "lc"
		aliasCustomer       = "cuc"
		aliasParentCustomer = "cup"
		aliasProductLicense = "pl"
		aliasProduct        = "pr"
		aliasLicenseVariant = "lv"
		aliasLicenseType    = "lt"
	)

	temp := *searchBy
	for index := range temp {
		switch temp[index].SearchKey {
		case "parent_customer_id":
			temp[index].SearchKey = fmt.Sprintf(`%s.id`, aliasParentCustomer)
		case "customer_name", "distributor_of", "province_id", "district_id", "customer_group_id", "customer_category_id", "salesman_id":
			temp[index].SearchKey = fmt.Sprintf(`%s.%s`, aliasCustomer, temp[index].SearchKey)
		case "id", "product_id", "product_valid_from", "product_valid_thru", "unique_id_1", "client_type_id", "allow_activation":
			temp[index].SearchKey = fmt.Sprintf(`%s.%s`, aliasLicenseConfig, temp[index].SearchKey)
		case "license_status":
			temp[index].SearchKey = fmt.Sprintf(`%s.%s`, aliasProductLicense, temp[index].SearchKey)
		}
	}

	switch userParam.OrderBy {
	case "product_name", "product_name ASC", "product_name DESC":
		userParam.OrderBy = fmt.Sprintf(`%s.%s`, aliasProduct, userParam.OrderBy)
	case "license_variant_name", "license_variant_name ASC", "license_variant_name DESC":
		userParam.OrderBy = fmt.Sprintf(`%s.%s`, aliasLicenseVariant, userParam.OrderBy)
	case "license_type_name", "license_type_name ASC", "license_type_name DESC":
		userParam.OrderBy = fmt.Sprintf(`%s.%s`, aliasLicenseType, userParam.OrderBy)
	case "customer_name", "customer_name ASC", "customer_name DESC":
		userParam.OrderBy = fmt.Sprintf(`%s.%s`, aliasCustomer, userParam.OrderBy)
	default:
		userParam.OrderBy = fmt.Sprintf(`%s.%s`, aliasLicenseConfig, userParam.OrderBy)
	}
}

func (input licenseConfigDAO) setCreatedByLicenseConfig(createdBy int64, searchBy *[]in.SearchByParam) {
	aliasName := "lc."
	if createdBy > 0 {
		*searchBy = append(*searchBy, in.SearchByParam{
			SearchKey:      aliasName + "created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}
}

func (input licenseConfigDAO) setGetListJoinLicenseConfig(getListData *getListJoinDataDAO) {
	getListData.InnerJoinAlias(CustomerDAO.TableName, "cup", "cup.id", "lc.parent_customer_id")
	getListData.InnerJoinAlias(CustomerDAO.TableName, "cuc", "cuc.id", "lc.customer_id")
	getListData.InnerJoinAlias(ProductDAO.TableName, "pr", "pr.id", "lc.product_id")
	getListData.InnerJoinAlias(LicenseVariantDAO.TableName, "lv", "lv.id", "lc.license_variant_id")
	getListData.InnerJoinAlias(LicenseTypeDAO.TableName, "lt", "lt.id", "lc.license_type_id")
	getListData.LeftJoinAliasWithoutDeleted(ProductLicenseDAO.TableName, "pl", "pl.license_config_id", "lc.id")
}

func (input licenseConfigDAO) addModuleFilled(params *[]interface{}, userParam repository.LicenseConfigModel) {
	if userParam.ModuleID1.Int64 > 0 {
		*params = append(*params, userParam.ModuleID1.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID2.Int64 > 0 {
		*params = append(*params, userParam.ModuleID2.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID3.Int64 > 0 {
		*params = append(*params, userParam.ModuleID3.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID4.Int64 > 0 {
		*params = append(*params, userParam.ModuleID4.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID5.Int64 > 0 {
		*params = append(*params, userParam.ModuleID5.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID6.Int64 > 0 {
		*params = append(*params, userParam.ModuleID6.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID7.Int64 > 0 {
		*params = append(*params, userParam.ModuleID7.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID8.Int64 > 0 {
		*params = append(*params, userParam.ModuleID8.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID9.Int64 > 0 {
		*params = append(*params, userParam.ModuleID9.Int64)
	} else {
		*params = append(*params, nil)
	}

	if userParam.ModuleID10.Int64 > 0 {
		*params = append(*params, userParam.ModuleID10.Int64)
	} else {
		*params = append(*params, nil)
	}
}

func (input licenseConfigDAO) CheckInstallationIsUsed(db *sql.DB, licenseConfigModel repository.LicenseConfigModel) (isUsed bool, err errorModel.ErrorModel) {
	var (
		funcName   = "CheckInstallationIsUsed"
		tempResult interface{}
		query      string
	)

	query = fmt.Sprintf(`
			SELECT 
				CASE WHEN (
					(count(id) > 0) OR 
					(SELECT count(id) 
						FROM %s 
						WHERE 
						id = $1 AND 
						(client_mapping_id > 0 OR client_mapping_id is not null) AND 
						deleted = FALSE)) 
				THEN TRUE ELSE FALSE END 
			FROM %s 
			WHERE 
			installation_id = $2 AND deleted = FALSE `,
		input.TableName, CustomerInstallationDAO.TableName)

	params := []interface{}{licenseConfigModel.InstallationID.Int64, licenseConfigModel.InstallationID.Int64}
	results := db.QueryRow(query, params...)

	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var temp bool
		dbError := results.Scan(&temp)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		isUsed = tempResult.(bool)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) GetInstallationIDByIDLicenseConfig(db *sql.DB, licenseConfigModel repository.LicenseConfigModel) (result repository.LicenseConfigModel, err errorModel.ErrorModel) {
	funcName := "GetInstallationIDByIDLicenseConfig"
	var tempResult interface{}

	query := fmt.Sprintf(`SELECT id, installation_id, created_by, 
			created_at, updated_at, updated_by 
			FROM %s 
			WHERE id = $1 AND deleted = FALSE `, input.TableName)

	params := []interface{}{licenseConfigModel.ID.Int64}

	if licenseConfigModel.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		params = append(params, licenseConfigModel.CreatedBy.Int64)
	}

	results := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var temp repository.LicenseConfigModel
		dbError := results.Scan(&temp.ID, &temp.InstallationID, &temp.CreatedBy,
			&temp.CreatedAt, &temp.UpdatedAt, &temp.UpdatedBy,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.LicenseConfigModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) ViewDetailLicenseConfig(db *sql.DB, licenseConfigModel repository.LicenseConfigModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var (
		funcName   = "ViewDetailLicenseConfig"
		parent     = "_parent"
		child      = "_child"
		query      string
		params     []interface{}
		tempResult interface{}
		indexStart int
	)

	query = fmt.Sprintf(`SELECT 
		lc.parent_customer_id, cup.customer_name as parent_customer, lc.customer_id, 
		lc.site_id, cuc.customer_name, lc.client_id, 
		pr.product_name, ct.client_type, lv.license_variant_name, 
		lt.license_type_name, lc.deployment_method, lc.no_of_user, 
		lc.is_user_concurrent, lc.unique_id_1, lc.unique_id_2, 
		lc.product_valid_from, lc.product_valid_thru, lc.max_offline_days, 
		mo1.module_name, mo2.module_name, mo3.module_name, 
		mo4.module_name, mo5.module_name, mo6.module_name, 
		mo7.module_name, mo8.module_name, mo9.module_name, 
		mo10.module_name, lc.product_id, lc.id, 
		lc.installation_id, lc.created_by, lc.created_at, 
		lc.updated_at, us.nt_username, lc.allow_activation 
		FROM %s lc 
			INNER JOIN %s cup ON lc.parent_customer_id = cup.id 
			INNER JOIN %s cuc ON lc.customer_id = cuc.id 
			INNER JOIN %s pr ON lc.product_id = pr.id 
			INNER JOIN %s ct ON lc.client_type_id = ct.id 
			INNER JOIN %s lv ON lc.license_variant_id = lv.id 
			INNER JOIN %s lt ON lc.license_type_id = lt.id 
			INNER JOIN "%s" us ON lc.updated_by = us.id 
			LEFT JOIN %s mo1 ON mo1.id = lc.module_id_1 
			LEFT JOIN %s mo2 ON mo2.id = lc.module_id_2 
			LEFT JOIN %s mo3 ON mo3.id = lc.module_id_3 
			LEFT JOIN %s mo4 ON mo4.id = lc.module_id_4 
			LEFT JOIN %s mo5 ON mo5.id = lc.module_id_5 
			LEFT JOIN %s mo6 ON mo6.id = lc.module_id_6 
			LEFT JOIN %s mo7 ON mo7.id = lc.module_id_7 
			LEFT JOIN %s mo8 ON mo8.id = lc.module_id_8 
			LEFT JOIN %s mo9 ON mo9.id = lc.module_id_9 
			LEFT JOIN %s mo10 ON mo10.id = lc.module_id_10 
		WHERE lc.id = $1 AND lc.deleted = FALSE `,
		input.TableName, CustomerDAO.TableName, CustomerDAO.TableName,
		ProductDAO.TableName, ClientTypeDAO.TableName, LicenseVariantDAO.TableName,
		LicenseTypeDAO.TableName, UserDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName)

	params = []interface{}{licenseConfigModel.ID.Int64}
	indexStart = 2

	if licenseConfigModel.CreatedBy.Int64 > 0 {
		query += " AND lc.created_by = $2 "
		params = append(params, licenseConfigModel.CreatedBy.Int64)
		indexStart = 3
	}

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, indexStart, []string{
		constanta.ProvinceDataScope + parent,
		constanta.DistrictDataScope + parent,
		constanta.SalesmanDataScope + parent,
		constanta.CustomerGroupDataScope + parent,
		constanta.CustomerCategoryDataScope + parent,
		constanta.ProvinceDataScope + child,
		constanta.DistrictDataScope + child,
		constanta.SalesmanDataScope + child,
		constanta.CustomerGroupDataScope + child,
		constanta.CustomerCategoryDataScope + child,
		constanta.ProductGroupDataScope,
		constanta.ClientTypeDataScope,
	})

	if additionalWhere != "" {
		query += additionalWhere
		for _, itemParam := range param {
			params = append(params, itemParam)
		}
	}

	results := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var temp repository.LicenseConfigModel
		dbError := results.Scan(&temp.ParentCustomerID, &temp.ParentCustomer, &temp.CustomerID,
			&temp.SiteID, &temp.Customer, &temp.ClientID,
			&temp.ProductName, &temp.ClientType, &temp.LicenseVariantName,
			&temp.LicenseTypeName, &temp.DeploymentMethod, &temp.NoOfUser,
			&temp.IsUserConcurrent, &temp.UniqueID1, &temp.UniqueID2,
			&temp.ProductValidFrom, &temp.ProductValidThru, &temp.MaxOfflineDays,
			&temp.ModuleIDName1, &temp.ModuleIDName2, &temp.ModuleIDName3,
			&temp.ModuleIDName4, &temp.ModuleIDName5, &temp.ModuleIDName6,
			&temp.ModuleIDName7, &temp.ModuleIDName8, &temp.ModuleIDName9,
			&temp.ModuleIDName10, &temp.ProductID, &temp.ID,
			&temp.InstallationID, &temp.CreatedBy, &temp.CreatedAt,
			&temp.UpdatedAt, &temp.UpdatedName, &temp.AllowActivation,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.LicenseConfigModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (colAdditionalWhere []string) {
	var (
		parent = "_parent"
		child  = "_child"
	)

	keyScope := []string{
		constanta.ProvinceDataScope + parent,
		constanta.DistrictDataScope + parent,
		constanta.SalesmanDataScope + parent,
		constanta.CustomerGroupDataScope + parent,
		constanta.CustomerCategoryDataScope + parent,
		constanta.ProvinceDataScope + child,
		constanta.DistrictDataScope + child,
		constanta.SalesmanDataScope + child,
		constanta.CustomerGroupDataScope + child,
		constanta.CustomerCategoryDataScope + child,
		constanta.ProductGroupDataScope,
		constanta.ClientTypeDataScope,
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

func (input licenseConfigDAO) GetLicenseForJSONFile(db *sql.DB, userParam repository.LicenseConfigModel) (result []repository.LicenseConfigModel, err errorModel.ErrorModel) {
	query := fmt.Sprintf(`SELECT
		lc.id, lc.installation_id, lc.client_id,
		p.product_id, lv.license_variant_name, lt.license_type_name,
		lc.deployment_method, lc.no_of_user, lc.unique_id_1, 
		lc.unique_id_2, lc.product_valid_from, lc.product_valid_thru, 
		m1.module_name, m2.module_name, m3.module_name, 
		m4.module_name, m5.module_name, m6.module_name, 
		m7.module_name, m8.module_name, m9.module_name, 
		m10.module_name, lc.max_offline_days,
		CASE WHEN COUNT(lcpc.id) > 0 THEN json_agg(json_build_object(
			'name', c.component_name,
			'value', lcpc.component_value
		)) ELSE null END component_value, lc.is_user_concurrent, 
		lc.parent_customer_id, lc.customer_id, lc.site_id,
		lc.client_type_id
	FROM %s lc
	LEFT JOIN %s p ON lc.product_id = p.id
	LEFT JOIN %s lv ON lc.license_variant_id = lv.id
	LEFT JOIN %s lt ON lc.license_type_id = lt.id
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
	LEFT JOIN %s pl ON pl.license_config_id = lc.id
	LEFT JOIN %s lcpc ON lc.id = lcpc.license_config_id
	LEFT JOIN %s c ON lcpc.component_id = c.id
	WHERE 
		lc.client_id = $1 AND lc.unique_id_1 = $2 AND
		lc.allow_activation = 'Y' AND
		lc.deleted = FALSE AND lc.product_valid_from <= now()::date AND 
		lc.product_valid_thru >= now()::date AND pl.id IS NULL `,
		input.TableName, ProductDAO.TableName, LicenseVariantDAO.TableName, LicenseTypeDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ProductLicenseDAO.TableName, LicenseConfigProductComponentDAO.TableName,
		ComponentDAO.TableName)

	param := []interface{}{
		userParam.ClientID.String,
		userParam.UniqueID1.String,
	}

	if userParam.UniqueID2.String != "" {
		query += " AND lc.unique_id_2 = $3 "
		param = append(param, userParam.UniqueID2.String)
	}

	query += ` GROUP BY 
		lc.id, p.product_id, lv.license_variant_name, lt.license_type_name,
		m1.module_name, m2.module_name, m3.module_name, m4.module_name, 
		m5.module_name, m6.module_name, m7.module_name, m8.module_name, 
		m9.module_name, m10.module_name `

	tempResult, err := GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.LicenseConfigModel
		dbErrorS := rows.Scan(
			&temp.ID, &temp.InstallationID, &temp.ClientID,
			&temp.ProductCode, &temp.LicenseVariantName, &temp.LicenseTypeName,
			&temp.DeploymentMethod, &temp.NoOfUser, &temp.UniqueID1,
			&temp.UniqueID2, &temp.ProductValidFrom, &temp.ProductValidThru,
			&temp.ModuleIDName1, &temp.ModuleIDName2, &temp.ModuleIDName3,
			&temp.ModuleIDName4, &temp.ModuleIDName5, &temp.ModuleIDName6,
			&temp.ModuleIDName7, &temp.ModuleIDName8, &temp.ModuleIDName9,
			&temp.ModuleIDName10, &temp.MaxOfflineDays, &temp.ComponentSting,
			&temp.IsUserConcurrent, &temp.ParentCustomerID, &temp.CustomerID,
			&temp.SiteID, &temp.ClientTypeID)
		return temp, dbErrorS
	}, param)

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.LicenseConfigModel))
		}
	}

	return
}

func (input licenseConfigDAO) GetLicenseConfigForValidationLicense(db *sql.DB, userParams []repository.LicenseConfigModel) (result []repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var param []interface{}
	index := 1
	query := fmt.Sprintf(`SELECT
		lc.id, lc.installation_id, lc.client_id,
		p.product_id, lv.license_variant_name, lt.license_type_name,
		lc.deployment_method, lc.no_of_user, lc.unique_id_1, 
		lc.unique_id_2, lc.product_valid_from, lc.product_valid_thru, 
		m1.module_name, m2.module_name, m3.module_name, 
		m4.module_name, m5.module_name, m6.module_name, 
		m7.module_name, m8.module_name, m9.module_name, 
		m10.module_name, lc.max_offline_days,
		CASE WHEN COUNT(lcpc.id) > 0 THEN json_agg(json_build_object(
			'name', c.component_name,
			'value', lcpc.component_value
		)) ELSE null END component_value, lc.is_user_concurrent, 
		lc.parent_customer_id, lc.customer_id, lc.site_id, 
		pl.license_status, lc.max_offline_days
	FROM %s lc
	LEFT JOIN %s p ON lc.product_id = p.id
	LEFT JOIN %s lv ON lc.license_variant_id = lv.id
	LEFT JOIN %s lt ON lc.license_type_id = lt.id
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
	LEFT JOIN %s pl ON pl.license_config_id = lc.id
	LEFT JOIN %s lcpc ON lc.id = lcpc.license_config_id
	LEFT JOIN %s c ON lcpc.component_id = c.id
	WHERE 
		lc.allow_activation = 'Y' AND lc.deleted = FALSE `,
		input.TableName, ProductDAO.TableName, LicenseVariantDAO.TableName, LicenseTypeDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ProductLicenseDAO.TableName, LicenseConfigProductComponentDAO.TableName,
		ComponentDAO.TableName)

	tempQuery, _ := ListRangeToInQueryWithStartIndex(len(userParams), index)

	for i := 0; i < len(userParams); i++ {
		param = append(param, userParams[i].ID.Int64)
	}

	query += fmt.Sprintf(` AND lc.id IN ( %s ) 
	GROUP BY 
		lc.id, p.product_id, lv.license_variant_name, lt.license_type_name,
		m1.module_name, m2.module_name, m3.module_name, m4.module_name, 
		m5.module_name, m6.module_name, m7.module_name, m8.module_name, 
		m9.module_name, m10.module_name, pl.license_status `, tempQuery)

	tempResult, err := GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.LicenseConfigModel
		dbErrorS := rows.Scan(
			&temp.ID, &temp.InstallationID, &temp.ClientID,
			&temp.ProductCode, &temp.LicenseVariantName, &temp.LicenseTypeName,
			&temp.DeploymentMethod, &temp.NoOfUser, &temp.UniqueID1,
			&temp.UniqueID2, &temp.ProductValidFrom, &temp.ProductValidThru,
			&temp.ModuleIDName1, &temp.ModuleIDName2, &temp.ModuleIDName3,
			&temp.ModuleIDName4, &temp.ModuleIDName5, &temp.ModuleIDName6,
			&temp.ModuleIDName7, &temp.ModuleIDName8, &temp.ModuleIDName9,
			&temp.ModuleIDName10, &temp.MaxOfflineDays, &temp.ComponentSting,
			&temp.IsUserConcurrent, &temp.ParentCustomerID, &temp.CustomerID,
			&temp.SiteID, &temp.ProductLicenseStatus, &temp.MaxOfflineDays)
		return temp, dbErrorS
	}, param)

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.LicenseConfigModel))
		}
	}

	return
}

func (input licenseConfigDAO) GetLicenseConfigForTransferUserLicense(db *sql.DB, userParam repository.LicenseConfigModel) (result repository.LicenseConfigModel, err errorModel.ErrorModel) {
	funcName := "GetLicenseConfigForTransferUserLicense"
	var tempResult interface{}

	query := fmt.Sprintf(`SELECT 
		lc.id, lc.product_valid_from, lc.product_valid_thru,
		lc.is_user_concurrent, ci.parent_customer_id
	FROM %s AS lc
	LEFT JOIN %s AS ci ON lc.installation_id = ci.id
	WHERE 
		lc.id = $1 AND lc.deleted = FALSE
	`, input.TableName, CustomerInstallationDAO.TableName)
	param := []interface{}{userParam.ID.Int64}
	rows := db.QueryRow(query, param...)
	if tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.LicenseConfigModel

		errorS := rws.Scan(
			&temp.ID,
			&temp.ProductValidFrom,
			&temp.ProductValidThru,
			&temp.IsUserConcurrent,
			&temp.ParentCustomerID)

		return temp, errorS
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.LicenseConfigModel)
	}

	return
}

func (input licenseConfigDAO) InsertLicenseConfigForTesting(db *sql.Tx, userParam repository.LicenseConfigModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertLicenseConfigForTesting"
	paramAmount := 34
	startIndex := 1

	query := fmt.Sprintf(`INSERT INTO %s 
			(installation_id, parent_customer_id, customer_id, 
			site_id, client_id, product_id, 
			client_type_id, license_variant_id, license_type_id, 
			deployment_method, no_of_user, is_user_concurrent, 
			max_offline_days, unique_id_1, unique_id_2, 
			product_valid_from, product_valid_thru, module_id_1, 
			module_id_2, module_id_3, module_id_4, 
			module_id_5, module_id_6, module_id_7, 
			module_id_8, module_id_9, module_id_10, 
			created_by, created_client, created_at, 
			updated_by, updated_client, updated_at, 
			allow_activation) VALUES `, input.TableName)

	query += CreateDollarParamInMultipleRowsDAO(1, paramAmount, startIndex, "id")

	params := []interface{}{
		userParam.InstallationID.Int64, userParam.ParentCustomerID.Int64, userParam.CustomerID.Int64,
		userParam.SiteID.Int64, userParam.ClientID.String, userParam.ProductID.Int64,
		userParam.ClientTypeID.Int64, userParam.LicenseVariantID.Int64, userParam.LicenseTypeID.Int64,
		userParam.DeploymentMethod.String, userParam.NoOfUser.Int64, userParam.IsUserConcurrent.String,
		userParam.MaxOfflineDays.Int64, userParam.UniqueID1.String, userParam.UniqueID2.String,
		userParam.ProductValidFrom.Time, userParam.ProductValidThru.Time,
	}

	input.addModuleFilled(&params, userParam)
	params = append(params,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.AllowActivation.String)

	result := db.QueryRow(query, params...)
	var tempResult interface{}
	if tempResult, err = RowCatchResult(result, func(rws *sql.Row) (interface{}, error) {
		var temp int64
		errorS := result.Scan(&temp)
		return temp, errorS
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		id = tempResult.(int64)
	}

	return
}

func (input licenseConfigDAO) PrepareScopeInLicenseConfig(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int) (additionalWhere []string) {
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

func (input licenseConfigDAO) GetLicenseConfigForValidation(db *sql.DB, userParams repository.LicenseConfigModel) (result repository.LicenseConfigModel, err errorModel.ErrorModel) {
	funcName := "GetLicenseConfigForValidation"
	//index := 1
	query := fmt.Sprintf(`SELECT
		lc.id, lc.installation_id, lc.client_id,
		p.product_id, lv.license_variant_name, lt.license_type_name,
		lc.deployment_method, lc.no_of_user, lc.unique_id_1, 
		lc.unique_id_2, lc.product_valid_from, lc.product_valid_thru, 
		m1.module_name, m2.module_name, m3.module_name, 
		m4.module_name, m5.module_name, m6.module_name, 
		m7.module_name, m8.module_name, m9.module_name, 
		m10.module_name, lc.max_offline_days,
		CASE WHEN COUNT(lcpc.id) > 0 THEN json_agg(json_build_object(
			'name', c.component_name,
			'value', lcpc.component_value
		)) ELSE null END component_value, lc.is_user_concurrent, 
		lc.parent_customer_id, lc.customer_id, lc.site_id, 
		pl.license_status, lc.max_offline_days, lc.client_type_id, 
		CASE WHEN COUNT(urd.id) > 0 THEN json_agg(json_build_object(
			'id', urd.id,
			'auth_user_id', urd.auth_user_id,
			'user_id', urd.user_id,
			'status', urd.status
		)) ELSE null END salesman_list
	FROM %s lc
	LEFT JOIN %s p ON lc.product_id = p.id
	LEFT JOIN %s lv ON lc.license_variant_id = lv.id
	LEFT JOIN %s lt ON lc.license_type_id = lt.id
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
	LEFT JOIN %s pl ON pl.license_config_id = lc.id
	LEFT JOIN %s lcpc ON lc.id = lcpc.license_config_id
	LEFT JOIN %s c ON lcpc.component_id = c.id
	LEFT JOIN %s ul ON pl.id = ul.product_license_id 
	LEFT JOIN %s urd on ul.id = urd.user_license_id 
	WHERE 
		lc.allow_activation = 'Y' AND lc.deleted = FALSE AND lc.id = $1 
	GROUP BY 
		lc.id, p.product_id, lv.license_variant_name, lt.license_type_name,
		m1.module_name, m2.module_name, m3.module_name, m4.module_name, 
		m5.module_name, m6.module_name, m7.module_name, m8.module_name, 
		m9.module_name, m10.module_name, pl.license_status `,
		input.TableName, ProductDAO.TableName, LicenseVariantDAO.TableName, LicenseTypeDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName, ModuleDAO.TableName,
		ModuleDAO.TableName, ModuleDAO.TableName, ProductLicenseDAO.TableName, LicenseConfigProductComponentDAO.TableName,
		ComponentDAO.TableName, UserLicenseDAO.TableName, UserRegistrationDetailDAO.TableName)

	row := db.QueryRow(query, userParams.ID.Int64)
	dbError := row.Scan(&result.ID, &result.InstallationID, &result.ClientID,
		&result.ProductCode, &result.LicenseVariantName, &result.LicenseTypeName,
		&result.DeploymentMethod, &result.NoOfUser, &result.UniqueID1,
		&result.UniqueID2, &result.ProductValidFrom, &result.ProductValidThru,
		&result.ModuleIDName1, &result.ModuleIDName2, &result.ModuleIDName3,
		&result.ModuleIDName4, &result.ModuleIDName5, &result.ModuleIDName6,
		&result.ModuleIDName7, &result.ModuleIDName8, &result.ModuleIDName9,
		&result.ModuleIDName10, &result.MaxOfflineDays, &result.ComponentSting,
		&result.IsUserConcurrent, &result.ParentCustomerID, &result.CustomerID,
		&result.SiteID, &result.ProductLicenseStatus, &result.MaxOfflineDays,
		&result.ClientTypeID, &result.SalesmanString)

	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigDAO) GetLicenseConfigDataForDuplicate(db *sql.DB, userParam repository.LicenseConfigModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, timeNow time.Time) (result []repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var (
		fileName           = input.FileName
		funcName           = "GetLicenseConfigDataForDuplicate"
		query              string
		tempResult         []interface{}
		param              []interface{}
		colAdditionalWhere []string
	)

	query = fmt.Sprintf(`SELECT 
		lc.id, lc.installation_id, lc.parent_customer_id, lc.customer_id, lc.site_id, lc.client_id, lc.product_id, 
		lc.client_type_id, lc.license_variant_id, lc.license_type_id, lc.deployment_method, lc.no_of_user, lc.is_user_concurrent, 
		lc.max_offline_days, lc.unique_id_1, lc.unique_id_2, lc.product_valid_from, lc.product_valid_thru, lc.module_id_1, 
		lc.module_id_2, lc.module_id_3, lc.module_id_4, lc.module_id_5, lc.module_id_6, lc.module_id_7, 
		lc.module_id_8, lc.module_id_9, lc.module_id_10,
		CASE WHEN COUNT(lcp.id) > 0 THEN 
			json_agg(json_build_object(
				'product_id', lcp.product_id,
				'component_id', lcp.component_id,
				'component_value', lcp.component_value
			)) ELSE null END component_value, 
		lc.allow_activation
		FROM license_configuration lc 
		INNER JOIN customer cup ON cup.id = lc.parent_customer_id
		INNER JOIN customer cuc ON cuc.id = lc.customer_id 
		INNER JOIN product pr ON pr.id = lc.product_id 
		LEFT JOIN license_configuration_productcomponent lcp ON lcp.license_config_id = lc.id
		WHERE lc.allow_activation = 'Y' `)

	colAdditionalWhere = input.setScopeData(scopeLimit, scopeDB, true)
	if len(colAdditionalWhere) > 0 {
		for _, itemAdditionalWhere := range colAdditionalWhere {
			query += fmt.Sprintf(` AND %s `, itemAdditionalWhere)
		}
	}

	if len(userParam.LicenseConfigIDs) > 0 {
		query += fmt.Sprintf(` AND lc.id IN (`)
	}

	for idx, itemLicenseConfig := range userParam.LicenseConfigIDs {
		query += fmt.Sprintf(`%d`, itemLicenseConfig.ID.Int64)
		if len(userParam.LicenseConfigIDs)-(idx+1) > 0 {
			query += fmt.Sprintf(`, `)
			continue
		}
		query += fmt.Sprintf(`)`)
	}

	query += fmt.Sprintf(` GROUP BY lc.id ORDER BY lc.product_valid_thru ASC `)
	tempResult, err = GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.LicenseConfigModel
		dbErrorS := rows.Scan(
			&temp.ID, &temp.InstallationID, &temp.ParentCustomerID, &temp.CustomerID, &temp.SiteID, &temp.ClientID, &temp.ProductID,
			&temp.ClientTypeID, &temp.LicenseVariantID, &temp.LicenseTypeID, &temp.DeploymentMethod, &temp.NoOfUser, &temp.IsUserConcurrent,
			&temp.MaxOfflineDays, &temp.UniqueID1, &temp.UniqueID2, &temp.ProductValidFrom, &temp.ProductValidThru, &temp.ModuleIDName1,
			&temp.ModuleIDName2, &temp.ModuleIDName3, &temp.ModuleIDName4, &temp.ModuleIDName5, &temp.ModuleIDName6, &temp.ModuleIDName7,
			&temp.ModuleIDName8, &temp.ModuleIDName9, &temp.ModuleIDName10,
			&temp.ComponentSting, &temp.AllowActivation)
		return temp, dbErrorS
	}, param)

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			var (
				componentModel []repository.ProductComponentModel
				temp           repository.LicenseConfigModel
			)

			temp = item.(repository.LicenseConfigModel)

			uuidKey, _ := uuid.NewRandom()
			temp.UUIDKey.String = uuidKey.String()

			temp.ProductValidFrom.Time = input.checkDateCompare(temp.ProductValidThru.Time, timeNow)
			temp.ProductValidThru.Time = userParam.ProductValidThru.Time
			temp.CreatedBy.Int64 = userParam.CreatedBy.Int64
			temp.CreatedClient.String = userParam.CreatedClient.String
			temp.CreatedAt.Time = userParam.CreatedAt.Time
			temp.UpdatedBy.Int64 = userParam.UpdatedBy.Int64
			temp.UpdatedClient.String = userParam.UpdatedClient.String
			temp.UpdatedAt.Time = userParam.UpdatedAt.Time
			_ = json.Unmarshal([]byte(temp.ComponentSting.String), &componentModel)

			if temp.ProductValidFrom.Time.After(temp.ProductValidThru.Time) {
				fieldName := util.GenerateConstantaI18n(constanta.ProductValidThru, constanta.DefaultApplicationsLanguage, nil)
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "PRODUCT_VALID_THRU_LESS_THAN_VALID_FROM", fmt.Sprintf(`%s ID %d`, fieldName, temp.ID.Int64), "")
				return
			}

			result = append(result, temp)
		}
	}

	return
}

func (input licenseConfigDAO) checkDateCompare(productValidThruDB time.Time, timeNow time.Time) (productValidFrom time.Time) {
	var (
		dateProductValidThru time.Time
		dateTimeNow          time.Time
		dateDefaultFormat    = constanta.DefaultInstallationTimeFormat
		dur, week            time.Duration
	)

	dateProductValidThru = productValidThruDB
	dateTimeNow, _ = time.Parse(dateDefaultFormat, timeNow.Format(dateDefaultFormat))
	productValidFrom = dateTimeNow

	if dateProductValidThru.After(dateTimeNow) {
		dur = dateProductValidThru.Sub(dateTimeNow)
		week = time.Hour * 24 * 7
		if dur > week || dur == week {
			productValidFrom = dateProductValidThru.Add(-1 * week)
		}
	}

	return
}

func (input licenseConfigDAO) InsertMultipleExtendedLicenseConfig(tx *sql.Tx, licenseConfigModel []repository.LicenseConfigModel) (idReturning []int64, err errorModel.ErrorModel) {
	var (
		funcName = "InsertMultipleExtendedLicenseConfig"
		//parameterField     = 36
		//startParameter     = 1
		query, tempQuery   string
		params, resultTemp []interface{}
		rows               *sql.Rows
		errorS             error
	)

	for i := 0; i < len(licenseConfigModel); i++ {
		tempQuery += fmt.Sprintf(`(
			%d, %d, %d, %d, %d, '%s', 
			%d, %d, %d, %d, '%s'::deploy_method, %d, 
			'%s'::flag_status, %d, '%s', '%s', TO_DATE('%s', 'YYYY-MM-DD'), TO_DATE('%s', 'YYYY-MM-DD'), `,
			licenseConfigModel[i].ID.Int64, licenseConfigModel[i].InstallationID.Int64, licenseConfigModel[i].ParentCustomerID.Int64, licenseConfigModel[i].CustomerID.Int64, licenseConfigModel[i].SiteID.Int64, licenseConfigModel[i].ClientID.String,
			licenseConfigModel[i].ProductID.Int64, licenseConfigModel[i].ClientTypeID.Int64, licenseConfigModel[i].LicenseVariantID.Int64, licenseConfigModel[i].LicenseTypeID.Int64, licenseConfigModel[i].DeploymentMethod.String, licenseConfigModel[i].NoOfUser.Int64,
			licenseConfigModel[i].IsUserConcurrent.String, licenseConfigModel[i].MaxOfflineDays.Int64, licenseConfigModel[i].UniqueID1.String, licenseConfigModel[i].UniqueID2.String, licenseConfigModel[i].ProductValidFrom.Time, licenseConfigModel[i].ProductValidThru.Time)

		if licenseConfigModel[i].ModuleID1.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID1.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID2.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID2.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID3.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID3.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID4.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID4.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID5.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID5.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID6.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID6.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID7.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID7.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID8.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID8.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID9.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID9.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		if licenseConfigModel[i].ModuleID10.Int64 > 0 {
			tempQuery += fmt.Sprintf(`%d, `, licenseConfigModel[i].ModuleID10.Int64)
		} else {
			tempQuery += fmt.Sprintf(`null::int8, `)
		}

		tempQuery += fmt.Sprintf(`'%s'::flag_status, %d, '%s', NOW()::timestamp, %d, '%s', 
			NOW()::timestamp, '%s'::uuid)`,
			licenseConfigModel[i].AllowActivation.String, licenseConfigModel[i].CreatedBy.Int64, licenseConfigModel[i].CreatedClient.String, licenseConfigModel[i].UpdatedBy.Int64, licenseConfigModel[i].UpdatedClient.String,
			licenseConfigModel[i].UUIDKey.String)

		if len(licenseConfigModel)-(i+1) > 0 {
			tempQuery += ", "
		}
	}

	query = fmt.Sprintf(`WITH data(
			id, installation_id, parent_customer_id, 
			customer_id, site_id, client_id, 
			product_id, client_type_id, license_variant_id, 
			license_type_id, deployment_method, no_of_user, 
			is_user_concurrent, max_offline_days, unique_id_1, 
			unique_id_2, product_valid_from, product_valid_thru, 
			module_id_1, module_id_2, module_id_3, 
			module_id_4, module_id_5, module_id_6, 
			module_id_7, module_id_8, module_id_9, 
			module_id_10, allow_activation, created_by, 
			created_client, created_at, updated_by, 
			updated_client, updated_at, uuid_key
		) AS (
		VALUES 
			%s
		)
		, liconfig as (
			insert into %s(
				uuid_key, installation_id, parent_customer_id, 
				customer_id, site_id, client_id, 
				product_id, client_type_id, license_variant_id, 
				license_type_id, deployment_method, no_of_user, 
				is_user_concurrent, max_offline_days, unique_id_1, 
				unique_id_2, product_valid_from, product_valid_thru, 
				allow_activation, module_id_1, module_id_2, 
				module_id_3, module_id_4, module_id_5, 
				module_id_6, module_id_7, module_id_8, 
				module_id_9, module_id_10, created_by, 
				created_client, created_at, updated_by, 
				updated_client, updated_at, old_license_configuration_id
			) select 
			uuid_key, installation_id, parent_customer_id, 
			customer_id, site_id, client_id, 
			product_id, client_type_id, license_variant_id, 
			license_type_id, deployment_method, no_of_user, 
			is_user_concurrent, max_offline_days, unique_id_1, 
			unique_id_2, product_valid_from, product_valid_thru, 
			allow_activation, module_id_1, module_id_2, 
			module_id_3, module_id_4, module_id_5, 
			module_id_6, module_id_7, module_id_8, 
			module_id_9, module_id_10, created_by, 
			created_client, created_at, updated_by, 
			updated_client, updated_at, id 
			from data
			returning installation_id, client_id, product_id, 
			unique_id_1, unique_id_2, product_valid_from, 
			product_valid_thru, uuid_key, id as lcid
		)
		, liconfigpc as (
			insert into %s(
				license_config_id, product_id, component_id, 
				component_value, created_by, created_client, 
				created_at, updated_by, updated_client, 
				updated_at
			) select 
			liconfig.lcid, lcp.product_id, lcp.component_id, 
			lcp.component_value, d.created_by, d.created_client, 
			d.created_at, d.updated_by, d.updated_client, 
			d.updated_at
			from data d 
			join liconfig using(
				installation_id, client_id, product_id, 
				unique_id_1, unique_id_2, product_valid_from, 
				product_valid_thru, uuid_key) 
			join %s lcp on lcp.license_config_id = d.id
		)
		select lcid from liconfig `,
		tempQuery, input.TableName, LicenseConfigProductComponentDAO.TableName,
		LicenseConfigProductComponentDAO.TableName)

	//for _, itemLicenseConfigModel := range licenseConfigModel {
	//	params = append(params,
	//		itemLicenseConfigModel.ID.Int64, itemLicenseConfigModel.InstallationID.Int64, itemLicenseConfigModel.ParentCustomerID.Int64,
	//		itemLicenseConfigModel.CustomerID.Int64, itemLicenseConfigModel.SiteID.Int64, itemLicenseConfigModel.ClientID.String,
	//		itemLicenseConfigModel.ProductID.Int64, itemLicenseConfigModel.ClientTypeID.Int64, itemLicenseConfigModel.LicenseVariantID.Int64,
	//		itemLicenseConfigModel.LicenseTypeID.Int64, itemLicenseConfigModel.DeploymentMethod.String, itemLicenseConfigModel.NoOfUser.Int64,
	//		itemLicenseConfigModel.IsUserConcurrent.String, itemLicenseConfigModel.MaxOfflineDays.Int64, itemLicenseConfigModel.UniqueID1.String,
	//		itemLicenseConfigModel.UniqueID2.String, itemLicenseConfigModel.ProductValidFrom.Time, itemLicenseConfigModel.ProductValidThru.Time,
	//		itemLicenseConfigModel.ModuleID1.Int64, itemLicenseConfigModel.ModuleID2.Int64, itemLicenseConfigModel.ModuleID3.Int64,
	//		itemLicenseConfigModel.ModuleID4.Int64, itemLicenseConfigModel.ModuleID5.Int64, itemLicenseConfigModel.ModuleID6.Int64,
	//		itemLicenseConfigModel.ModuleID7.Int64, itemLicenseConfigModel.ModuleID8.Int64, itemLicenseConfigModel.ModuleID9.Int64,
	//		itemLicenseConfigModel.ModuleID10.Int64, itemLicenseConfigModel.AllowActivation.String, itemLicenseConfigModel.CreatedBy.Int64,
	//		itemLicenseConfigModel.CreatedClient.String, itemLicenseConfigModel.CreatedAt.Time, itemLicenseConfigModel.UpdatedBy.Int64,
	//		itemLicenseConfigModel.UpdatedClient.String, itemLicenseConfigModel.UpdatedAt.Time)
	//}

	rows, errorS = tx.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			errs   error
			idTemp int64
		)

		errs = rws.Scan(&idTemp)
		if errs != nil {
			err = errorModel.GenerateInternalDBServerError(input.TableName, funcName, errs)
			return
		}

		resultTemp = idTemp
		return
	})

	if err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		idTemp := itemResultTemp.(int64)
		idReturning = append(idReturning, idTemp)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
