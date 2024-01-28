package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"strings"
)

type customerDAO struct {
	AbstractDAO
}

var CustomerDAO = customerDAO{}.New()

func (input customerDAO) New() (output customerDAO) {
	output.FileName = "CustomerDAO.go"
	output.TableName = "customer"

	return
}

func (input customerDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	for i := 0; i < len(*searchByParam); i++ {
		(*searchByParam)[i].SearchKey = "c." + (*searchByParam)[i].SearchKey
	}

	switch userParam.OrderBy {
	case "updated_name", "updated_name ASC", "updated_name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "u.nt_username " + strSplit[1]
		} else {
			userParam.OrderBy = "u.nt_username"
		}
		break
	case "district_name", "district_name ASC", "district_name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "d.name " + strSplit[1]
		} else {
			userParam.OrderBy = "d.name "
		}
		break
	case "province_name", "province_name ASC", "province_name DESC":
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = "p.name " + strSplit[1]
		} else {
			userParam.OrderBy = "p.name "
		}
		break
	default:
		userParam.OrderBy = "c." + userParam.OrderBy
		break
	}
}

func (input customerDAO) convertUserParamAndSearchByDistributor(_ *in.GetListDataDTO, searchByParam *[]in.SearchByParam) {
	if searchByParam != nil {
		for i := 0; i < len(*searchByParam); i++ {
			if (*searchByParam)[i].SearchKey == "license_variant" {
				(*searchByParam)[i].SearchKey = "lv.license_variant_name"
			} else {
				(*searchByParam)[i].SearchKey = "ct." + (*searchByParam)[i].SearchKey
			}
		}
	}
}

func (input customerDAO) getCustomerDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "c.id"},
		Deleted:   FieldStatus{FieldName: "c.deleted"},
		CreatedBy: FieldStatus{FieldName: "c.created_by", Value: createdBy},
		UpdatedAt: FieldStatus{FieldName: "c.updated_at"},
	}
}

func (input customerDAO) DeleteCustomer(db *sql.Tx, userParam repository.CustomerModel) (err errorModel.ErrorModel) {
	funcName := "DeleteCustomer"

	query := fmt.Sprintf(`UPDATE %s
	SET
		deleted = $1, updated_by = $2, updated_at = $3, updated_client = $4, npwp = $5 
	WHERE 
		id = $6 
	`, input.TableName)

	param := []interface{}{
		true,
		userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.Npwp.String,
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

func (input customerDAO) InternalGetListDistributor(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam) (output []interface{}, err errorModel.ErrorModel) {
	var (
		dbParam         []interface{}
		query           string
		additionalWhere string
	)

	query = fmt.Sprintf(`
		WITH principals AS (
			SELECT id, is_principal 
			FROM customer 
			where is_principal = true
		)
		SELECT 
			c.id as id, c.mdb_company_profile_id, ct.client_type,
			pc.id as principal_id, c.company_title as dist_title, c.customer_name as dist_name,
			c.npwp as dist_npwp, c.address as dist_address, c.hamlet as dist_hamlet,
			c.neighbourhood as dist_neighbourhood, c.country_id as dist_country, c.province_id as dist_province,
			c.district_id as dist_district, c.sub_district_id as dist_subdistrict, c.urban_village_id as dist_urban_village,
			c.postal_code_id as dist_postal_code, c.long as dist_long, c.lat as dist_lat,
			c.phone as dist_phone, c.fax as dist_fax, c.company_email as dist_email,
			ci.installation_date as dist_joindate, lc.product_valid_from as dist_fromdate, lc.product_valid_thru as dist_expirydate,
			lc.unique_id_1 as company_id, lc.unique_id_2 as branch_id, pl.activation_date as activation_date,
			c.updated_at as updated_date, pl.client_id, us.auth_user_id, lv.license_variant_name
		FROM %s c
		INNER JOIN %s ci ON c.id = ci.customer_id
		LEFT JOIN principals pc on pc.id = c.parent_customer_id
		INNER JOIN %s lc ON lc.installation_id = ci.id
		INNER JOIN %s lv ON lv.id = lc.license_variant_id
		INNER JOIN %s p ON p.id = lc.product_id
		INNER JOIN %s ct ON p.client_type_id = ct.id
		INNER JOIN %s pl ON lc.id = pl.license_config_id 
		INNER JOIN "%s" us ON us.client_id = pl.client_id `,
		input.TableName, CustomerInstallationDAO.TableName, LicenseConfigDAO.TableName,
		LicenseVariantDAO.TableName, ProductDAO.TableName, ClientTypeDAO.TableName,
		ProductLicenseDAO.TableName, UserDAO.TableName)

	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "updated_at" {
			additionalWhere = fmt.Sprintf(` 
				AND pl.updated_at BETWEEN '%s'::TIMESTAMP AND '%s'::TIMESTAMP `,
				searchByParam[i].SearchValue, fmt.Sprintf(`%s 23:59:59`, searchByParam[i].SearchValue))
			searchByParam = append(searchByParam[:i], searchByParam[i+1:]...)
			i = -1
		}
	}

	input.convertUserParamAndSearchByDistributor(&userParam, &searchByParam)
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CustomerModel
			dbError := rows.Scan(
				&temp.ID, &temp.MDBCompanyProfileID, &temp.ClientType,
				&temp.PrincipalID, &temp.CompanyTitle, &temp.CustomerName,
				&temp.Npwp, &temp.Address, &temp.Hamlet,
				&temp.Neighbourhood, &temp.CountryID, &temp.ProvinceID,
				&temp.DistrictID, &temp.SubDistrictID, &temp.UrbanVillageID,
				&temp.PostalCodeID, &temp.Longitude, &temp.Latitude,
				&temp.Phone, &temp.Fax, &temp.CompanyEmail,
				&temp.InstallationDate, &temp.ProductValidFrom, &temp.ProductValidThru,
				&temp.UniqueID1, &temp.UniqueID2, &temp.ActivationDate,
				&temp.UpdatedAt, &temp.ClientID, &temp.AuthUserID, &temp.LicenseVariantName,
			)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "c.id"},
			Deleted:   FieldStatus{FieldName: "c.deleted"},
			CreatedBy: FieldStatus{FieldName: "c.created_by", Value: int64(0)},
			UpdatedAt: FieldStatus{FieldName: "pl.updated_at"},
		})
}

func (input customerDAO) GetCustomerForDelete(db *sql.Tx, userParam repository.CustomerModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (output repository.CustomerModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetCustomerForDelete"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
			c.id, c.updated_at, c.updated_by, c.is_parent,
			CASE WHEN 
				(SELECT COUNT(ci.id) FROM %s ci WHERE (ci.parent_customer_id = c.id OR ci.customer_id = c.id) AND ci.deleted = FALSE) > 0 
					OR
				(SELECT COUNT(cs.id) FROM %s cs WHERE (cs.parent_customer_id = c.id OR cs.customer_id = c.id) AND cs.deleted = FALSE) > 0
					OR
				(SELECT COUNT(lc.id) FROM %s lc WHERE (lc.parent_customer_id = c.id OR lc.customer_id = c.id) AND lc.deleted = FALSE) > 0 
			THEN TRUE ELSE FALSE END is_used, 
			(
				SELECT 
					json_agg(
						json_build_object(
							'id' ,cc.id, 
							'customer_id', cc.customer_id, 
							'nik', cc.nik,
							'updated_by', cc.updated_by, 
							'updated_at', cc.updated_at, 
							'created_by', cc.created_by
						)
					)
				FROM %s cc 
				WHERE cc.customer_id = c.id AND cc.deleted = FALSE
			) AS customer_contact_value, c.created_by, c.npwp
		FROM %s c
		LEFT JOIN %s p ON c.province_id = p.id 
		LEFT JOIN %s d ON c.district_id = d.id 
		LEFT JOIN %s cg ON c.customer_group_id = cg.id 
		LEFT JOIN %s s ON c.salesman_id = s.id 
		WHERE c.id = $1 AND c.deleted = FALSE `,
		CustomerInstallationDAO.TableName, CustomerSiteDAO.TableName, LicenseConfigDAO.TableName,
		CustomerContactDAO.TableName, input.TableName, ProvinceDAO.TableName,
		DistrictDAO.TableName, CustomerGroupDAO.TableName, SalesmanDAO.TableName)

	param := []interface{}{userParam.ID.Int64}
	additionalWhere, additionalParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalWhere != "" {
		query += " " + additionalWhere + " "
		param = append(param, additionalParam...)
	}

	query += fmt.Sprintf(` GROUP BY c.id `)

	results := db.QueryRow(query, param...)
	errs := results.Scan(
		&output.ID, &output.UpdatedAt, &output.UpdatedBy,
		&output.IsParent, &output.IsUsed, &output.CustomerContactStr,
		&output.CreatedBy, &output.Npwp)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) UpdateCustomer(db *sql.Tx, userParam repository.CustomerModel) (err errorModel.ErrorModel) {
	funcName := "UpdateCustomer"
	var param []interface{}
	query := fmt.Sprintf(`UPDATE %s 
	SET
		parent_customer_id = $1, mdb_parent_customer_id = $2, mdb_company_profile_id = $3, 
		mdb_company_title_id = $4, is_principal = $5, is_parent = $6,
		npwp = $7, company_title = $8, customer_name = $9,
		address = $10, hamlet = $11, neighbourhood = $12,
		country_id = $13, province_id = $14, district_id = $15,
		sub_district_id = $16, urban_village_id = $17, postal_code_id = $18,
		lat = $19, long = $20, phone = $21, 
		alternative_phone = $22, fax = $23, company_email = $24, 
		alternative_company_email = $25, customer_source = $26, tax_name = $27,
		tax_address = $28, salesman_id = $29, ref_customer_id = $30, 
		distributor_of = $31, customer_group_id = $32, customer_category_id = $33,
		status = $34, updated_by = $35, updated_client = $36, updated_at = $37, 
		address_2 = $38, address_3 = $39
	WHERE id = $40`, input.TableName)

	if userParam.ParentCustomerID.Int64 > 0 {
		param = append(param, userParam.ParentCustomerID.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.MDBParentCustomerID.Int64 > 0 {
		param = append(param, userParam.MDBParentCustomerID.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.MDBCompanyProfileID.Int64 > 0 {
		param = append(param, userParam.MDBCompanyProfileID.Int64)
	} else {
		param = append(param, nil)
	}

	if userParam.MDBCompanyTitleID.Int64 > 0 {
		param = append(param, userParam.MDBCompanyTitleID.Int64)
	} else {
		param = append(param, nil)
	}

	param = append(param,
		userParam.IsPrincipal.Bool, userParam.IsParent.Bool,
		userParam.Npwp.String, userParam.CompanyTitle.String,
		userParam.CustomerName.String, userParam.Address.String,
		userParam.Hamlet.String, userParam.Neighbourhood.String,
		userParam.CountryID.Int64, userParam.ProvinceID.Int64,
		userParam.DistrictID.Int64, userParam.SubDistrictID.Int64,
		userParam.UrbanVillageID.Int64, userParam.PostalCodeID.Int64,
		userParam.Latitude.Float64, userParam.Longitude.Float64,
		userParam.Phone.String, userParam.AlternativePhone.String,
		userParam.Fax.String, userParam.CompanyEmail.String,
		userParam.AlternativeCompanyEmail.String, userParam.CustomerSource.String,
		userParam.TaxName.String, userParam.TaxAddress.String,
		userParam.SalesmanID.Int64, userParam.RefCustomerID.Int64,
		userParam.DistributorOF.String, userParam.CustomerGroupID.Int64,
		userParam.CustomerCategoryID.Int64, userParam.Status.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.Address2.String,
		userParam.Address3.String, userParam.ID.Int64,
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

func (input customerDAO) GetCustomerForUpdate(db *sql.Tx, userParam repository.CustomerModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (output repository.CustomerModel, err errorModel.ErrorModel) {
	funcName := "GetCustomerForUpdate"
	var tempResult interface{}
	index := 1
	query := fmt.Sprintf(`SELECT 
		c.id, c.updated_at, c.updated_by, 
		c.parent_customer_id, c.npwp, c.mdb_company_profile_id, 
		CASE WHEN
			(SELECT COUNT(id) FROM %s WHERE parent_customer_id = c.id OR customer_id = c.id) > 0 
				OR
			(SELECT COUNT(id) FROM %s WHERE parent_customer_id = c.id OR customer_id = c.id) > 0 
		THEN TRUE ELSE FALSE END is_used 
	FROM %s c
	LEFT JOIN %s p ON c.province_id = p.id
	LEFT JOIN %s d ON c.district_id = d.id
	LEFT JOIN %s cg ON c.customer_group_id = cg.id
	LEFT JOIN %s cc ON c.customer_category_id = cc.id
	LEFT JOIN %s s ON c.salesman_id = s.id
	WHERE c.id = $%s AND c.deleted = FALSE`,
		CustomerInstallationDAO.TableName, CustomerSiteDAO.TableName, input.TableName,
		ProvinceDAO.TableName, DistrictDAO.TableName, CustomerGroupDAO.TableName,
		CustomerCategoryDAO.TableName, SalesmanDAO.TableName, strconv.Itoa(index))

	param := []interface{}{userParam.ID.Int64}
	index += 1

	additionalWhere, additionalParam := ScopeToAddedQueryView(scopeLimit, scopeDB, index,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalWhere != "" {
		query += " " + additionalWhere + " "
		param = append(param, additionalParam...)
		index += len(additionalParam)
	}

	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerDefaultMustCheck, index)

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerModel
		dbError := rws.Scan(&temp.ID, &temp.UpdatedAt, &temp.UpdatedBy,
			&temp.ParentCustomerID, &temp.Npwp, &temp.MDBCompanyProfileID,
			&temp.IsUsed)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		output = tempResult.(repository.CustomerModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) GetCustomerByMdbCompanyProfile(db *sql.DB, userParam repository.CustomerModel) (result repository.CustomerModel, err errorModel.ErrorModel) {
	funcName := "GetCustomerByMdbCompanyProfile"
	var tempResult interface{}

	query := fmt.Sprintf(`SELECT 
		id, customer_name, mdb_company_profile_id 
	FROM %s
	WHERE mdb_company_profile_id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.MDBCompanyProfileID.Int64}

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerModel
		dbError := rws.Scan(
			&temp.ID, &temp.CustomerName,
			&temp.MDBCompanyProfileID,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) ViewCustomer(db *sql.DB, userParam repository.CustomerModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerModel, err errorModel.ErrorModel) {
	var (
		funcName               = "ViewCustomer"
		query, additionalQuery string
		params, tempParams     []interface{}
		tempResult             interface{}
	)

	query = fmt.Sprintf(`SELECT
			c.id, c.is_principal, c.is_parent,
			c.parent_customer_id, c.mdb_parent_customer_id, c.mdb_company_profile_id,
			c.npwp, c.mdb_company_title_id, c.company_title,
			c.customer_name, c.address, c.hamlet,
			c.neighbourhood, c.country_id, c.province_id,
			p.name AS province_name, c.district_id, d.name AS district_name,
			c.sub_district_id, c.urban_village_id, c.postal_code_id,
			c.long, c.lat, c.phone,
			c.alternative_phone, c.fax, c.company_email,
			c.alternative_company_email, c.customer_source, c.tax_name,
			c.tax_address, c.salesman_id, c.ref_customer_id,
			c.distributor_of, c.customer_group_id, c.customer_category_id,
			c.status, c.created_by, c.created_at, 
			c.updated_by, c.updated_at, c.address_2, c.address_3, 
			pcs.customer_name, sd.name AS sub_district_name, uv.name AS urban_village_name, 
			pc.code AS postal_code_name,
			CASE WHEN
				(SELECT COUNT(id) FROM %s WHERE parent_customer_id = c.id OR customer_id = c.id) > 0
					OR
				(SELECT COUNT(id) FROM %s WHERE parent_customer_id = c.id OR customer_id = c.id) > 0
			THEN TRUE ELSE FALSE END is_used, uc.nt_username, ud.nt_username,
			cc.customer_category_name, cg.customer_group_name,
			( SELECT customer_name FROM customer where id = c.ref_customer_id ),
			s.first_name, s.last_name
		FROM %s c
		LEFT JOIN %s p ON c.province_id = p.id
		LEFT JOIN %s d ON c.district_id = d.id
		LEFT JOIN %s cg ON c.customer_group_id = cg.id
		LEFT JOIN %s cc ON c.customer_category_id = cc.id
		LEFT JOIN %s s ON c.salesman_id = s.id
		LEFT JOIN %s pcs ON c.parent_customer_id = pcs.id
		LEFT JOIN %s sd ON c.sub_district_id = sd.id
		LEFT JOIN %s uv ON c.urban_village_id = uv.id
		LEFT JOIN %s pc ON c.postal_code_id = pc.id
		LEFT JOIN "%s" uc ON uc.id = c.created_by
		LEFT JOIN "%s" ud ON ud.id = c.updated_by
		WHERE c.id = $1 AND c.deleted = FALSE `,
		CustomerInstallationDAO.TableName, CustomerSiteDAO.TableName, input.TableName,
		ProvinceDAO.TableName, DistrictDAO.TableName, CustomerGroupDAO.TableName,
		CustomerCategoryDAO.TableName, SalesmanDAO.TableName, input.TableName,
		SubDistrictDAO.TableName, UrbanVillageDAO.TableName, PostalCodeDAO.TableName,
		UserDAO.TableName, UserDAO.TableName)

	params = []interface{}{userParam.ID.Int64}
	additionalQuery, tempParams = ScopeToAddedQueryView(scopeLimit, scopeDB, 2,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalQuery != "" {
		query += additionalQuery
		params = append(params, tempParams...)
	}

	row := db.QueryRow(query, params...)
	tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerModel
		dbError := rws.Scan(
			&temp.ID, &temp.IsPrincipal, &temp.IsParent,
			&temp.ParentCustomerID, &temp.MDBParentCustomerID, &temp.MDBCompanyProfileID,
			&temp.Npwp, &temp.MDBCompanyTitleID, &temp.CompanyTitle,
			&temp.CustomerName, &temp.Address, &temp.Hamlet,
			&temp.Neighbourhood, &temp.CountryID, &temp.ProvinceID,
			&temp.ProvinceName, &temp.DistrictID, &temp.DistrictName,
			&temp.SubDistrictID, &temp.UrbanVillageID, &temp.PostalCodeID,
			&temp.Longitude, &temp.Latitude, &temp.Phone,
			&temp.AlternativePhone, &temp.Fax, &temp.CompanyEmail,
			&temp.AlternativeCompanyEmail, &temp.CustomerSource, &temp.TaxName,
			&temp.TaxAddress, &temp.SalesmanID, &temp.RefCustomerID,
			&temp.DistributorOF, &temp.CustomerGroupID, &temp.CustomerCategoryID,
			&temp.Status, &temp.CreatedBy, &temp.CreatedAt,
			&temp.UpdatedBy, &temp.UpdatedAt, &temp.Address2,
			&temp.Address3, &temp.ParentCustomerName, &temp.SubDistrictName,
			&temp.UrbanVillageName, &temp.PostalCode, &temp.IsUsed, &temp.CreatedName, &temp.UpdatedName,
			&temp.CustomerCategoryName, &temp.CustomerGroupName, &temp.RefCustomerName, &temp.SalesmanFirstName,
			&temp.SalesmanLastName)
		return temp, dbError
	}, input.FileName, funcName)
	if err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) GetCountCustomer(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var dbParam []interface{}
	additionalWhere := ""

	query := fmt.Sprintf(` c 
		LEFT JOIN  %s p ON c.province_id = p.id 
		LEFT JOIN  %s d ON c.district_id = d.id 
		LEFT JOIN  %s cg ON c.customer_group_id = cg.id 
		LEFT JOIN  %s cc ON c.customer_category_id = cc.id 
		LEFT JOIN  %s s ON c.salesman_id = s.id `, ProvinceDAO.TableName, DistrictDAO.TableName,
		CustomerGroupDAO.TableName, CustomerCategoryDAO.TableName, SalesmanDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "c." + param.SearchKey
	}

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, dbParam, input.TableName+query,
		searchByParam, additionalWhere,
		input.getCustomerDefaultMustCheck(createdBy))
}

func (input customerDAO) GetListCustomer(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (output []interface{}, err errorModel.ErrorModel) {
	var (
		dbParam         []interface{}
		param           []interface{}
		query           string
		additionalWhere string
		index           = 1
	)

	query = fmt.Sprintf(`SELECT 
			c.id, c.customer_name, c.address, 
			c.province_id, p.name AS province_name, 
			c.district_id, d.name AS district_name, 
			c.phone, c.status, 
			c.created_by, c.created_at, 
			c.updated_by, c.updated_at, c.npwp 
		FROM %s c 
		LEFT JOIN %s p ON c.province_id = p.id 
		LEFT JOIN %s d ON c.district_id = d.id 
		LEFT JOIN %s cg ON c.customer_group_id = cg.id 
		LEFT JOIN %s cc ON c.customer_category_id = cc.id 
		LEFT JOIN %s s ON c.salesman_id = s.id `,
		input.TableName, ProvinceDAO.TableName, DistrictDAO.TableName,
		CustomerGroupDAO.TableName, CustomerCategoryDAO.TableName, SalesmanDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	additionalWhere, param = ScopeToAddedQueryView(scopeLimit, scopeDB, index,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CustomerModel
			dbError := rows.Scan(
				&temp.ID, &temp.CustomerName, &temp.Address,
				&temp.ProvinceID, &temp.ProvinceName, &temp.DistrictID,
				&temp.DistrictName, &temp.Phone, &temp.Status,
				&temp.CreatedBy, &temp.CreatedAt, &temp.UpdatedBy,
				&temp.UpdatedAt, &temp.Npwp,
			)
			return temp, dbError
		}, additionalWhere, input.getCustomerDefaultMustCheck(createdBy))
}

func (input customerDAO) IsExistCustomerForInsert(db *sql.DB, userParam repository.CustomerModel) (result bool, err errorModel.ErrorModel) {
	funcName := "IsExistCustomerForInsert"
	index := 1

	query := fmt.Sprintf(`SELECT 
		CASE WHEN COUNT(c.id) > 0 THEN TRUE ELSE FALSE END
	FROM %s c
	WHERE
		c.id = $1 AND c.deleted = FALSE`, input.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += 1

	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerDefaultMustCheck, index)
	dbError := db.QueryRow(query, param...).Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) GetCustomerParentForValidate(db *sql.DB, userParam repository.CustomerModel,
	scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.CustomerModel, err errorModel.ErrorModel) {
	var (
		tempResult      interface{}
		funcName        = "GetCustomerParentForValidate"
		additionalWhere string
		tempParam       []interface{}
	)
	index := 1

	query := fmt.Sprintf(`SELECT 
		c.id, c.mdb_company_profile_id,c.updated_at, c.created_by
	FROM %s c
	LEFT JOIN %s p ON c.province_id = p.id 
	LEFT JOIN %s d ON c.district_id = d.id 
	LEFT JOIN %s cg ON c.customer_group_id = cg.id 
	LEFT JOIN %s cc ON c.customer_category_id = cc.id 
	LEFT JOIN %s s ON c.salesman_id = s.id 
	WHERE
		c.id = $1 AND c.deleted = FALSE AND 
		c.is_parent = TRUE `,
		input.TableName, ProvinceDAO.TableName, DistrictDAO.TableName,
		CustomerGroupDAO.TableName, CustomerCategoryDAO.TableName, SalesmanDAO.TableName)

	param := []interface{}{userParam.ID.Int64}
	index += 1

	index = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerDefaultMustCheck, index)

	additionalWhere, tempParam = ScopeToAddedQueryView(scopeLimit, scopeDB, index,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalWhere != "" {
		query += additionalWhere
		param = append(param, tempParam...)
	}

	row := db.QueryRow(query, param...)

	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.CustomerModel
		dbError := rws.Scan(
			&temp.ID, &temp.MDBCompanyProfileID, &temp.UpdatedAt, &temp.CreatedBy,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.CustomerModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) GetCustomerParent(db *sql.DB, userParam repository.CustomerModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result bool, err errorModel.ErrorModel) {
	var (
		funcName        = "GetCustomerParent"
		index           = 1
		query           string
		param           []interface{}
		dbError         error
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT CASE WHEN COUNT(c.id) > 0 THEN TRUE ELSE FALSE END FROM %s c
		WHERE c.id = $1 AND c.deleted = FALSE AND c.is_parent = TRUE `, input.TableName)

	if scopeLimit != nil || scopeDB != nil {
		additionalWhere = input.setScopeData(scopeLimit, scopeDB, true)
		if len(additionalWhere) > 0 {
			strWhere := " AND " + strings.Join(additionalWhere, " AND ")
			strWhere = strings.TrimRight(strWhere, " AND ")
			query += strWhere
		}
	}

	param = []interface{}{userParam.ID.Int64}
	index += 1

	_ = CheckOwnPermissionAndGetQuery(userParam.CreatedBy.Int64, &query, &param, input.getCustomerDefaultMustCheck, index)

	dbError = db.QueryRow(query, param...).Scan(&result)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) InsertCustomer(db *sql.Tx, userParam repository.CustomerModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertCustomer"
	var params []interface{}

	query := fmt.Sprintf(`INSERT INTO %s
	(is_principal, is_parent, parent_customer_id, mdb_parent_customer_id, mdb_company_profile_id, 
	mdb_company_title_id, npwp, company_title, customer_name, address, hamlet, neighbourhood, country_id, province_id, 
	district_id, sub_district_id, urban_village_id, postal_code_id, lat, long, phone, alternative_phone, fax, 
	company_email, alternative_company_email, customer_source, tax_name, tax_address, salesman_id, ref_customer_id, 
	distributor_of, customer_group_id, customer_category_id, status, created_by, created_at, created_client, updated_by,
	updated_at, updated_client, address_2, address_3)
		VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
	$25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42)
		RETURNING id`, input.TableName)

	params = append(params,
		userParam.IsPrincipal.Bool, userParam.IsParent.Bool,
	)

	HandleOptionalParam([]interface{}{
		userParam.ParentCustomerID.Int64, userParam.MDBParentCustomerID.Int64,
		userParam.MDBCompanyProfileID.Int64, userParam.MDBCompanyTitleID.Int64,
	}, &params)

	params = append(params,
		userParam.Npwp.String, userParam.CompanyTitle.String,
		userParam.CustomerName.String, userParam.Address.String,
		userParam.Hamlet.String, userParam.Neighbourhood.String,
		userParam.CountryID.Int64, userParam.ProvinceID.Int64,
		userParam.DistrictID.Int64, userParam.SubDistrictID.Int64,
		userParam.UrbanVillageID.Int64, userParam.PostalCodeID.Int64,
		userParam.Latitude.Float64, userParam.Longitude.Float64,
		userParam.Phone.String, userParam.AlternativePhone.String,
		userParam.Fax.String, userParam.CompanyEmail.String,
		userParam.AlternativeCompanyEmail.String, userParam.CustomerSource.String,
		userParam.TaxName.String, userParam.TaxAddress.String,
		userParam.SalesmanID.Int64, userParam.RefCustomerID.Int64,
		userParam.DistributorOF.String, userParam.CustomerGroupID.Int64,
		userParam.CustomerCategoryID.Int64, userParam.Status.String,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time,
		userParam.CreatedClient.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.Address2.String, userParam.Address3.String)

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input customerDAO) CheckCustomerIsExist(db *sql.DB, userParam repository.CustomerModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (isExist bool, err errorModel.ErrorModel) {
	var (
		funcName        = "CheckCustomerIsExist"
		query           string
		params          []interface{}
		results         *sql.Row
		dbError         error
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT (CASE WHEN count(cu.id) > 0 THEN TRUE ELSE FALSE END) is_exist FROM %s cu
		WHERE cu.id = $1 AND cu.deleted = FALSE `, input.TableName)

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

func (input customerDAO) GetNameCustomer(db *sql.DB, userParam repository.CustomerModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (resultDB repository.CustomerModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetNameCustomer"
		query    string
	)

	query = fmt.Sprintf(`SELECT cu.id, cu.customer_name, ci.updated_at 
		FROM %s cu 
		LEFT JOIN %s ci ON cu.id = ci.parent_customer_id 
		WHERE 
		cu.id = $1 AND cu.deleted = FALSE `,
		input.TableName, CustomerInstallationDAO.TableName)

	additionalWhere := input.setScopeData(scopeLimit, scopeDB, true)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	params := []interface{}{userParam.ID.Int64}
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&resultDB.ID, &resultDB.CustomerName, &resultDB.UpdatedAt)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerDAO) GetListCustomerByStatusParent(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isParent bool) (output []interface{}, err errorModel.ErrorModel) {
	var dbParam []interface{}

	query := fmt.Sprintf(`SELECT
		c.id, c.customer_name, c.address,
		c.province_id, p.name AS province_name,
		c.district_id, d.name AS district_name,
		c.phone, c.status,
		c.created_by, c.created_at,
		c.updated_by, c.updated_at,
		c.mdb_company_profile_id, c.npwp
	FROM %s c
	LEFT JOIN %s p ON c.province_id = p.id
	LEFT JOIN %s d ON c.district_id = d.id
	LEFT JOIN %s cg ON c.customer_group_id = cg.id
	LEFT JOIN %s cc ON c.customer_category_id = cc.id
	LEFT JOIN %s s ON c.salesman_id = s.id `,
		input.TableName, ProvinceDAO.TableName, DistrictDAO.TableName,
		CustomerGroupDAO.TableName, CustomerCategoryDAO.TableName, SalesmanDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}

	additionalWhere += " AND is_parent = $" + strconv.Itoa(len(param)+1) + " "
	dbParam = append(dbParam, isParent)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CustomerModel
			dbError := rows.Scan(
				&temp.ID, &temp.CustomerName, &temp.Address,
				&temp.ProvinceID, &temp.ProvinceName, &temp.DistrictID,
				&temp.DistrictName, &temp.Phone, &temp.Status,
				&temp.CreatedBy, &temp.CreatedAt, &temp.UpdatedBy,
				&temp.UpdatedAt, &temp.MDBCompanyProfileID, &temp.Npwp,
			)
			return temp, dbError
		}, additionalWhere, input.getCustomerDefaultMustCheck(createdBy))
}

func (input customerDAO) GetCountCustomerByStatusParent(db *sql.DB, searchByParam []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isParent bool) (result int, err errorModel.ErrorModel) {
	var dbParam []interface{}
	additionalWhere := ""

	query := fmt.Sprintf(` c 
	LEFT JOIN %s p ON c.province_id = p.id
	LEFT JOIN %s d ON c.district_id = d.id
	LEFT JOIN %s cg ON c.customer_group_id = cg.id
	LEFT JOIN %s cc ON c.customer_category_id = cc.id
	LEFT JOIN %s s ON c.salesman_id = s.id `,
		ProvinceDAO.TableName, DistrictDAO.TableName, CustomerGroupDAO.TableName,
		CustomerCategoryDAO.TableName, SalesmanDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "c." + param.SearchKey
	}

	additionalWhere, param := ScopeToAddedQueryView(scopeLimit, scopeDB, 1,
		[]string{
			constanta.CustomerGroupDataScope,
			constanta.CustomerCategoryDataScope,
			constanta.SalesmanDataScope,
			constanta.ProvinceDataScope,
			constanta.DistrictDataScope,
		})

	if additionalWhere != "" {
		dbParam = append(dbParam, param...)
	}

	additionalWhere += " AND is_parent = $" + strconv.Itoa(len(param)+1) + " "
	dbParam = append(dbParam, isParent)

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, dbParam, input.TableName+query,
		searchByParam, additionalWhere, input.getCustomerDefaultMustCheck(createdBy))
}

func (input customerDAO) setScopeData(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool) (colAdditionalWhere []string) {
	var (
		parent = "_parent"
	)

	keyScope := []string{
		constanta.ProvinceDataScope + parent,
		constanta.DistrictDataScope + parent,
		constanta.SalesmanDataScope + parent,
		constanta.CustomerGroupDataScope + parent,
		constanta.CustomerCategoryDataScope + parent,
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

func (input customerDAO) InternalGetListCustomer(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam) (output []interface{}, err errorModel.ErrorModel) {
	var (
		dbParam         []interface{}
		query           string
		additionalWhere string
	)

	query = fmt.Sprintf(`SELECT 
			c.id, c.mdb_company_profile_id, c.npwp,
			c.is_principal, c.is_parent, c.company_title, 
			c.customer_name, c.address, c.phone, c.company_email
		FROM %s c 
		LEFT JOIN %s p ON c.province_id = p.id 
		LEFT JOIN %s d ON c.district_id = d.id 
		LEFT JOIN %s cg ON c.customer_group_id = cg.id 
		LEFT JOIN %s cc ON c.customer_category_id = cc.id 
		LEFT JOIN %s s ON c.salesman_id = s.id `,
		input.TableName, ProvinceDAO.TableName, DistrictDAO.TableName,
		CustomerGroupDAO.TableName, CustomerCategoryDAO.TableName, SalesmanDAO.TableName)

	input.convertUserParamAndSearchBy(&userParam, &searchByParam)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, dbParam, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.CustomerModel
			dbError := rows.Scan(
				&temp.ID, &temp.MDBCompanyProfileID, &temp.Npwp,
				&temp.IsPrincipal, &temp.IsParent, &temp.CompanyTitle,
				&temp.CustomerName, &temp.Address, &temp.Phone, &temp.CompanyEmail,
			)
			return temp, dbError
		}, additionalWhere, input.getCustomerDefaultMustCheck(0))
}

func (input customerDAO) InternalGetCountCustomer(db *sql.DB, searchByParam []in.SearchByParam, userParam in.GetListDataDTO) (result int, err errorModel.ErrorModel) {
	var dbParam []interface{}
	additionalWhere := ""

	query := fmt.Sprintf(` c 
		LEFT JOIN  %s p ON c.province_id = p.id 
		LEFT JOIN  %s d ON c.district_id = d.id 
		LEFT JOIN  %s cg ON c.customer_group_id = cg.id 
		LEFT JOIN  %s cc ON c.customer_category_id = cc.id 
		LEFT JOIN  %s s ON c.salesman_id = s.id `, ProvinceDAO.TableName, DistrictDAO.TableName,
		CustomerGroupDAO.TableName, CustomerCategoryDAO.TableName, SalesmanDAO.TableName)

	for i, param := range searchByParam {
		searchByParam[i].SearchKey = "c." + param.SearchKey
	}

	return GetListDataDAO.GetCountDataWithUpdatedAtAndDefaultMustCheck(db, dbParam, input.TableName+query,
		searchByParam, additionalWhere,
		input.getCustomerDefaultMustCheck(0), userParam)
}

func (input customerDAO) GetCustomerHasChildHierarchy(db *sql.DB, userParam repository.CustomerModel) (count int, output []out.DetailErrorCustomerResponse, err errorModel.ErrorModel) {
	var (
		fileName  = "CustomerDAO.go"
		funcName  = "GetCustomerHasChildHierarchy"
		param     []interface{}
		query     string
		row       *sql.Row
		errS      error
		dataChild string
	)

	query = fmt.Sprintf(`select 
		count(id),
		json_agg(
			json_build_object(
				'id', id,
				'npwp', npwp,
				'customer_name', customer_name
			)
		) as child from %s 
		where parent_customer_id = $1 and deleted = false 
		group by parent_customer_id `,
		input.TableName)

	param = []interface{}{userParam.ID.Int64}
	row = db.QueryRow(query, param...)
	errS = row.Scan(&count, &dataChild)
	if errS != nil && errS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errS)
		return
	}

	if dataChild != "" {
		_ = json.Unmarshal([]byte(dataChild), &output)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
