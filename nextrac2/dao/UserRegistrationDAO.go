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
)

type userRegistrationAdminDAO struct {
	AbstractDAO
}

var UserRegistrationAdminDAO = userRegistrationAdminDAO{}.New()

func (input userRegistrationAdminDAO) New() (output userRegistrationAdminDAO) {
	output.FileName = "UserRegistrationDAO.go"
	output.TableName = "user_registration"
	return
}

func (input userRegistrationAdminDAO) PrepareScopeInUserRegistrationAdmin(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int) (additionalWhere []string) {
	for key := range scopeLimit {
		var (
			keyDataScope   string
			additionalTemp string
		)

		switch key {
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

func (input userRegistrationAdminDAO) GetCountUserRegistrationAdmin(db *sql.DB, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT COUNT(%s.id) FROM %s `, input.TableName, input.TableName)

	for index := range searchBy {
		switch searchBy[index].SearchKey {
		case "id":
			searchBy[index].SearchKey = fmt.Sprintf(`%s.%s`, input.TableName, searchBy[index].SearchKey)
		case "company_name":
			searchBy[index].SearchKey = fmt.Sprintf(`%s.%s`, input.TableName, searchBy[index].SearchKey)
		}
	}

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      fmt.Sprintf(`%s.created_by`, input.TableName),
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}

	additionalWhere = input.PrepareScopeInUserRegistrationAdmin(scopeLimit, scopeDB, 1)

	getListData := getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: additionalWhere}
	getListData.InnerJoinAlias(CustomerDAO.TableName, "pc", "pc.id", "user_registration.parent_customer_id")
	getListData.InnerJoinAlias(CustomerDAO.TableName, "cu", "cu.id", "user_registration.customer_id")

	return getListData.GetCountJoinData(db, searchBy, 0)
}

func (input userRegistrationAdminDAO) UpdateUserRegistrationAdmin(tx *sql.Tx, userParam repository.UserRegistrationAdminModel) (result repository.UserRegistrationAdminModel, err errorModel.ErrorModel) {
	funcName := "UpdateUserRegistrationAdmin"

	query := fmt.Sprintf(
		`UPDATE %s SET 
					parent_customer_id = $1,
					site_id = $2,
					unique_id_1 = $3,
					unique_id_2 = $4,
					company_name = $5,
					branch_name = $6,
					user_admin = $7,
					password_admin = $8,
					updated_by = $9,
					updated_client = $10,
					updated_at = $11,
					client_type_id = $12,
					client_mapping_id = $13,
					customer_id = $14
				WHERE id = $15 
				RETURNING id`, input.TableName)

	params := []interface{}{
		userParam.ParentCustomerId.Int64, userParam.SiteId.Int64, userParam.UniqueID1.String,
		userParam.UniqueID2.String, userParam.CompanyName.String, userParam.BranchName.String,
		userParam.UserAdmin.String, userParam.PasswordAdmin.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String, userParam.UpdatedAt.Time, userParam.ClientTypeID.Int64,
		userParam.ClientMappingID.Int64, userParam.CustomerId.Int64, userParam.ID.Int64,
	}

	results := tx.QueryRow(query, params...)

	dbError := results.Scan(&result.ID)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input userRegistrationAdminDAO) GetListUserRegistrationAdmin(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT cu.customer_name, pc.customer_name as parent_customer_name, user_registration.unique_id_1 as company_id, 
		user_registration.unique_id_2 as branch_id, user_registration.company_name, user_registration.branch_name, 
		user_registration.user_admin, user_registration.password_admin, user_registration.id as id 
		FROM %s `, input.TableName)

	input.setSearchByUserRegistrationAdmin(&searchBy)

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      fmt.Sprintf(`%s.created_by`, input.TableName),
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}

	additionalWhere = input.PrepareScopeInUserRegistrationAdmin(scopeLimit, scopeDB, 1)
	getListData := getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: additionalWhere}
	input.setGetListJoinUserRegistrationAdmin(&getListData)

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.UserRegistrationAdminModel

		dbError := rows.Scan(
			&resultTemp.CustomerName, &resultTemp.ParentCustomerName, &resultTemp.UniqueID1,
			&resultTemp.UniqueID2, &resultTemp.CompanyName, &resultTemp.BranchName,
			&resultTemp.UserAdmin, &resultTemp.PasswordAdmin, &resultTemp.ID)

		return resultTemp, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input userRegistrationAdminDAO) setSearchByUserRegistrationAdmin(searchBy *[]in.SearchByParam) {
	var codeSeparate = "."
	temp := *searchBy
	for index := range temp {
		switch temp[index].SearchKey {
		case "id":
			temp[index].SearchKey = input.TableName + codeSeparate + temp[index].SearchKey
		case "company_name":
			temp[index].SearchKey = input.TableName + codeSeparate + temp[index].SearchKey
		}
	}
}

func (input userRegistrationAdminDAO) setGetListJoinUserRegistrationAdmin(getListData *getListJoinDataDAO) {
	getListData.LeftJoinAlias(CustomerDAO.TableName, "pc", "pc.id", "user_registration.parent_customer_id")
	getListData.LeftJoinAlias(CustomerDAO.TableName, "cu", "cu.id", "user_registration.customer_id")
}

func (input userRegistrationAdminDAO) ViewDetailUserRegistrationAdmin(db *sql.DB, userParam repository.UserRegistrationAdminModel, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (resultOnDB repository.UserRegistrationAdminModel, err errorModel.ErrorModel) {
	var (
		funcName        = "ViewDetailUserRegistrationAdmin"
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT ur.id, ur.parent_customer_id, pcu.customer_name, ur.customer_id, ur.site_id, cu.customer_name, 
		ur.company_name, ur.branch_name, ur.unique_id_1, ur.unique_id_2, ur.user_admin, ur.password_admin 
		FROM %s ur 
		JOIN %s cu ON ur.customer_id = cu.id 
		JOIN %s pcu ON ur.parent_customer_id = pcu.id 
		WHERE 
		ur.id = $1 AND ur.deleted = FALSE AND cu.deleted = FALSE AND 
		pcu.deleted = FALSE `,
		input.TableName, CustomerDAO.TableName, CustomerDAO.TableName)

	parameters := []interface{}{userParam.ID.Int64}
	if createdBy > 0 {
		query += fmt.Sprintf(` AND ur.created_by = $2 `)
		parameters = append(parameters, createdBy)
	}

	additionalWhere = input.PrepareScopeInUserRegistrationAdmin(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		for i := 0; i < len(additionalWhere); i++ {
			query += fmt.Sprintf(` AND %s `, additionalWhere[i])
		}
	}

	dbResult := db.QueryRow(query, parameters...)
	dbError := dbResult.Scan(
		&resultOnDB.ID, &resultOnDB.ParentCustomerId, &resultOnDB.ParentCustomerName,
		&resultOnDB.CustomerId, &resultOnDB.SiteId, &resultOnDB.CustomerName,
		&resultOnDB.CompanyName, &resultOnDB.BranchName, &resultOnDB.UniqueID1,
		&resultOnDB.UniqueID2, &resultOnDB.UserAdmin, &resultOnDB.PasswordAdmin,
	)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()

	return
}

func (input userRegistrationAdminDAO) GetUserRegistrationForRegistrationNexMileNexStar(db *sql.DB, userParam repository.UserRegistrationAdminModel) (result repository.UserRegistrationAdminModel, err errorModel.ErrorModel) {
	funcName := "GetUserRegistrationForRegistrationNexMileNexStar"

	query := fmt.Sprintf(`
		SELECT id 
			FROM %s 
		WHERE 
			site_id = $1 AND unique_id_1 = $2 AND unique_id_2 = $3 AND parent_customer_id = $4 AND client_type_id = $5 AND deleted = FALSE`, input.TableName)

	params := []interface{}{
		userParam.SiteId.Int64,
		userParam.UniqueID1.String,
		userParam.UniqueID2.String,
		userParam.ParentCustomerId.Int64,
		userParam.ClientTypeID.Int64,
	}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&result.ID,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationAdminDAO) GetUserRegistrationByUniqueID1AndUniqueID2(db *sql.DB, userParam repository.UserRegistrationAdminModel) (result repository.UserRegistrationAdminModel, err errorModel.ErrorModel) {
	funcName := "GetUserRegistrationForRegistrationNexMileNexStar"

	query := fmt.Sprintf(`SELECT id FROM %s WHERE unique_id_1 = $1 AND unique_id_2 = $2 AND deleted = FALSE`, input.TableName)

	params := []interface{}{
		userParam.UniqueID1.String,
		userParam.UniqueID2.String,
	}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&result.ID,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationAdminDAO) GetFieldForNexmileParameter(db *sql.DB, userParam repository.UserRegistrationAdminModel) (result repository.UserRegistrationAdminModel, err errorModel.ErrorModel) {
	funcName := "GetFieldForNexmileParameter"

	query := fmt.Sprintf(`SELECT 
		id, user_admin, password_admin, company_name, branch_name,
		unique_id_1, unique_id_2
	FROM %s 
	WHERE 
		unique_id_1 = $1 AND unique_id_2 = $2 
		AND client_type_id = $3 AND client_mapping_id = $4 
		AND deleted = FALSE `, input.TableName)

	params := []interface{}{
		userParam.UniqueID1.String,
		userParam.UniqueID2.String,
		userParam.ClientTypeID.Int64,
		userParam.ClientMappingID.Int64,
	}

	dbResult := db.QueryRow(query, params...)
	dbError := dbResult.Scan(
		&result.ID, &result.UserAdmin, &result.PasswordAdmin,
		&result.CompanyName, &result.BranchName, &result.UniqueID1,
		&result.UniqueID2,
	)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	fmt.Println("GetFieldForNexmileParameter", userParam)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationAdminDAO) InsertUserRegistrationAdmin(db *sql.Tx, userParam repository.UserRegistrationAdminModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertUserRegistrationAdmin"

	query := "INSERT INTO " + input.TableName + " " +
		"(unique_id_1, unique_id_2, company_name, " +
		"branch_name, user_admin, password_admin, " +
		"created_by, created_client, created_at, " +
		"updated_by, updated_client, updated_at, " +
		"parent_customer_id, customer_id, site_id, " +
		"client_mapping_id, client_type_id) " +
		"VALUES " +
		"($1, $2, $3, " +
		"$4, $5, $6, " +
		"$7, $8, $9, " +
		"$10, $11, $12, " +
		"$13, $14, $15, " +
		"$16, $17) " +
		"RETURNING id "

	params := []interface{}{
		userParam.UniqueID1.String, userParam.UniqueID2.String, userParam.CompanyName.String,
		userParam.BranchName.String, userParam.UserAdmin.String, userParam.PasswordAdmin.String,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.ParentCustomerId.Int64, userParam.CustomerId.Int64, userParam.SiteId.Int64,
		userParam.ClientMappingID.Int64, userParam.ClientTypeID.Int64,
	}

	results := db.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input userRegistrationAdminDAO) GetUserRegistrationByClientMappingID(db *sql.DB, userParam repository.UserRegistrationAdminModel) (result repository.UserRegistrationAdminModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUserRegistrationByClientMappingID"
		query    string
		params   []interface{}
	)

	query = fmt.Sprintf(`
		SELECT id, company_name, branch_name 
		FROM %s 
		WHERE 
		unique_id_1 = $1 AND
		deleted = FALSE AND 
		client_mapping_id = $2 AND 
		client_type_id = $3 `, input.TableName)

	params = []interface{}{
		userParam.UniqueID1.String, userParam.ClientMappingID.Int64, userParam.ClientTypeID.Int64}

	if userParam.UniqueID2.String != "" {
		query += " AND unique_id_2 = $4 "
		params = append(params, userParam.UniqueID2.String)
	}

	dbResult := db.QueryRow(query, params...)
	fmt.Println("params : ", params)
	fmt.Printf(query, input.TableName)
	fmt.Println("Query GetUserRegistrationByClientMappingID : ", query)
	fmt.Println("DB Result get User Regis : ", dbResult)
	dbError := dbResult.Scan(&result.ID, &result.CompanyName, &result.BranchName)
	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
