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

type salesmanDAO struct {
	AbstractDAO
}

var SalesmanDAO = salesmanDAO{}.New()

func (input salesmanDAO) New() (output salesmanDAO) {
	output.FileName = "SalesmanDAO.go"
	output.TableName = "salesman"
	return
}

func (input salesmanDAO) GetSalesmanForInsert(db *sql.DB, salesmanModel repository.SalesmanModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result bool, err errorModel.ErrorModel) {
	funcName := "GetSalesmanForInsert"

	query := fmt.Sprintf(`SELECT CASE WHEN COUNT(salesman.id) > 0 THEN TRUE ELSE FALSE END 
				FROM %s salesman 
				INNER JOIN %s p ON salesman.province_id = p.id 
				WHERE salesman.id = $1 AND salesman.deleted = FALSE `, input.TableName, ProvinceDAO.TableName)

	params := []interface{}{salesmanModel.ID.Int64}

	scopeAdditionalWhere, scopeParam := ScopeToAddedQueryView(scopeLimit, scopeDB, 2, input.getAllScopeConstanta(scopeDB))
	if scopeAdditionalWhere != "" {
		query += " " + scopeAdditionalWhere
		params = append(params, scopeParam...)
	}

	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanDAO) InsertSalesman(db *sql.Tx, userParam repository.SalesmanModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertSalesman"
	query := "INSERT INTO " + input.TableName + "" +
		"(person_title_id, person_title, nik, " +
		"first_name, last_name, sex, " +
		"address, phone, email, " +
		"created_by, created_at, created_client, " +
		"updated_by, updated_at, updated_client, " +
		"status, hamlet, neighbourhood, " +
		"province_id, district_id) " +
		"VALUES " +
		"($1, $2, $3, " +
		"$4, $5, $6, " +
		"$7, $8, $9, " +
		"$10, $11, $12, " +
		"$13, $14, $15, " +
		"$16, $17, $18, " +
		"$19, $20) " +
		"RETURNING id "

	params := []interface{}{
		userParam.PersonTitleID.Int64, userParam.PersonTitle.String, userParam.Nik.String,
		userParam.FirstName.String, userParam.LastName.String, userParam.Sex.String,
		userParam.Address.String, userParam.Phone.String, userParam.Email.String,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.CreatedClient.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
	}

	if userParam.Status.String != "" {
		params = append(params, userParam.Status.String)
	} else {
		params = append(params, "A")
	}

	if userParam.Hamlet.String != "" {
		params = append(params, userParam.Hamlet.String)
	} else {
		params = append(params, nil)
	}

	if userParam.Neighbourhood.String != "" {
		params = append(params, userParam.Neighbourhood.String)
	} else {
		params = append(params, nil)
	}

	params = append(params, userParam.ProvinceID.Int64, userParam.DistrictID.Int64)

	result := db.QueryRow(query, params...)
	errorS := result.Scan(&id)

	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input salesmanDAO) GetListSalesman(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, regionalDeleted bool) (result []interface{}, err errorModel.ErrorModel) {
	var (
		query           string
		additionalWhere []string
		getListData     getListJoinDataDAO
	)

	query = fmt.Sprintf(`SELECT 
		salesman.first_name as first_name, salesman.last_name as last_name, salesman.id as id, 
		salesman.status as status, salesman.address as address, district.name as district, 
		province.name as province, salesman.phone as phone, salesman.email as email, 
		salesman.updated_at FROM %s `, input.TableName)

	for index := range searchBy {
		switch searchBy[index].SearchKey {
		case "id":
			searchBy[index].SearchKey = "salesman." + searchBy[index].SearchKey
		case "first_name":
			searchBy[index].SearchKey = "salesman." + searchBy[index].SearchKey
		}
	}

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "salesman.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}

	if strings.Contains(userParam.OrderBy, "full_name") {
		strSplit := strings.Split(userParam.OrderBy, " ")
		if len(strSplit) > 1 {
			userParam.OrderBy = fmt.Sprintf("TRIM(salesman.first_name) %s, TRIM(salesman.last_name) %s", strSplit[1], strSplit[1])
		} else {
			userParam.OrderBy = fmt.Sprintf("TRIM(salesman.first_name) ASC, TRIM(salesman.last_name) ASC")
		}
	}

	additionalWhere = input.PrepareScopeInSalesman(scopeLimit, scopeDB, 1)
	getListData = getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: additionalWhere}
	if regionalDeleted {
		getListData.InnerJoin(ProvinceDAO.TableName, "province.id", "salesman.province_id")
		getListData.InnerJoin(DistrictDAO.TableName, "district.id", "salesman.district_id")
	} else {
		getListData.InnerJoinWithoutDeleted(ProvinceDAO.TableName, "province.id", "salesman.province_id")
		getListData.InnerJoinWithoutDeleted(DistrictDAO.TableName, "district.id", "salesman.district_id")
	}

	mappingFunc := func(rows *sql.Rows) (interface{}, error) {
		var resultTemp repository.ListSalesmanModel

		dbError := rows.Scan(
			&resultTemp.FirstName, &resultTemp.LastName, &resultTemp.ID,
			&resultTemp.Status, &resultTemp.Address, &resultTemp.District,
			&resultTemp.Province, &resultTemp.Phone, &resultTemp.Email,
			&resultTemp.UpdatedAt)

		return resultTemp, dbError
	}

	return getListData.GetListJoinData(db, userParam, searchBy, 0, mappingFunc)
}

func (input salesmanDAO) GetCountSalesman(db *sql.DB, searchBy []in.SearchByParam, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result int, err errorModel.ErrorModel) {
	var (
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT COUNT(%s.id) FROM %s `, input.TableName, input.TableName)

	for index := range searchBy {
		switch searchBy[index].SearchKey {
		case "id":
			searchBy[index].SearchKey = "salesman." + searchBy[index].SearchKey
		case "first_name":
			searchBy[index].SearchKey = "salesman." + searchBy[index].SearchKey
		}
	}

	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:      "salesman.created_by",
			SearchValue:    strconv.Itoa(int(createdBy)),
			SearchOperator: "eq",
			DataType:       "number",
			SearchType:     "FILTER",
		})
	}

	additionalWhere = input.PrepareScopeInSalesman(scopeLimit, scopeDB, 1)

	getListData := getListJoinDataDAO{Table: input.TableName, Query: query, AdditionalWhere: additionalWhere}
	getListData.InnerJoin(ProvinceDAO.TableName, "province.id", "salesman.province_id")
	getListData.InnerJoin(DistrictDAO.TableName, "district.id", "salesman.district_id")

	return getListData.GetCountJoinData(db, searchBy, 0)
}

func (input salesmanDAO) UpdateSalesman(db *sql.Tx, userParam repository.SalesmanModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateSalesman"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`UPDATE %s SET 
			person_title_id = $1, first_name = $2, last_name = $3, 
			sex = $4, address = $5, hamlet = $6, 
			neighbourhood = $7, province_id = $8, district_id = $9, 
			phone = $10, email = $11, updated_by = $12, 
			updated_at = $13, updated_client = $14, status = $15, 
			person_title = $16 
			WHERE id = $17 `, input.TableName)

	param = []interface{}{
		userParam.PersonTitleID.Int64, userParam.FirstName.String, userParam.LastName.String,
		userParam.Sex.String, userParam.Address.String,
	}

	if userParam.Hamlet.String != "" {
		param = append(param, userParam.Hamlet.String)
	} else {
		param = append(param, nil)
	}

	if userParam.Neighbourhood.String != "" {
		param = append(param, userParam.Neighbourhood.String)
	} else {
		param = append(param, nil)
	}

	param = append(param, userParam.ProvinceID.Int64, userParam.DistrictID.Int64, userParam.Phone.String,
		userParam.Email.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.Status.String, userParam.PersonTitle.String,
		userParam.ID.Int64)

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

func (input salesmanDAO) DeleteSalesman(db *sql.Tx, userParam repository.SalesmanModel) (err errorModel.ErrorModel) {
	funcName := "DeleteSalesman"

	query := "UPDATE " + input.TableName + " SET " +
		" deleted = TRUE, " +
		" updated_by = $1, " +
		" updated_client = $2, " +
		" updated_at = $3, " +
		" nik = $4 " +
		" WHERE " +
		" id = $5"

	param := []interface{}{
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.Nik.String, userParam.ID.Int64}

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

func (input salesmanDAO) ViewSalesman(db *sql.DB, userParam repository.SalesmanModel) (result repository.ViewSalesmanModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewSalesman"
		subQuery string
		query    string
		params   []interface{}
		results  *sql.Row
		dbError  error
	)

	subQuery = fmt.Sprintf(`SELECT sl.id, sl.person_title_id, sl.sex, 
		sl.nik, sl.first_name, sl.last_name, 
		sl.address, sl.hamlet, sl.neighbourhood, 
		sl.province_id, sl.district_id, sl.phone, 
		sl.email, sl.status, sl.created_at, 
		sl.updated_at, sl.updated_by as updated_by, pv.name, 
		ds.name, sl.person_title 
		FROM %s sl 
		INNER JOIN province pv ON pv.id = sl.province_id 
		INNER JOIN district ds ON ds.id = sl.district_id 
		WHERE sl.id = $1 AND sl.deleted = FALSE `,
		input.TableName)

	params = []interface{}{userParam.ID.Int64}
	if userParam.CreatedBy.Int64 > 0 {
		subQuery += " AND sl.created_by = $2 "
		params = append(params, userParam.CreatedBy.Int64)
	}

	query = fmt.Sprintf(`
		SELECT *, (SELECT nt_username as updated_name FROM "%s" WHERE id = a.updated_by) FROM (%s) a `,
		UserDAO.TableName, subQuery)

	results = db.QueryRow(query, params...)
	dbError = results.Scan(
		&result.ID, &result.PersonTitleID, &result.Sex,
		&result.Nik, &result.FirstName, &result.LastName,
		&result.Address, &result.Hamlet, &result.Neighbourhood,
		&result.MdbProvinceID, &result.MdbDistrictID, &result.Phone,
		&result.Email, &result.Status, &result.CreatedAt,
		&result.UpdatedAt, &result.UpdatedBy, &result.Province,
		&result.District, &result.PersonTitle, &result.UpdatedName)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanDAO) GetSalesmanForUpdateDelete(db *sql.DB, salesmanModel repository.SalesmanModel, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (result repository.SalesmanModel, err errorModel.ErrorModel) {
	var (
		funcName        = "GetSalesmanForUpdateDelete"
		query           string
		additionalWhere []string
	)

	query = fmt.Sprintf(`SELECT id, updated_at, created_by, 
			(SELECT 
			CASE WHEN count(id) > 0 THEN TRUE ELSE FALSE END 
			FROM %s WHERE salesman_id = salesman.id) isUsed, nik 
		FROM %s 
		WHERE 
		id = $1 AND deleted = FALSE `,
		CustomerDAO.TableName, input.TableName)

	params := []interface{}{salesmanModel.ID.Int64}

	if salesmanModel.CreatedBy.Int64 > 0 {
		query += fmt.Sprintf(` AND created_by = $2 `)
		params = append(params, salesmanModel.CreatedBy.Int64)
	}

	additionalWhere = input.PrepareScopeInSalesman(scopeLimit, scopeDB, 1)
	if len(additionalWhere) > 0 {
		strWhere := " AND " + strings.Join(additionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		query += strWhere
	}

	query += fmt.Sprintf(` FOR UPDATE `)
	results := db.QueryRow(query, params...)
	dbError := results.Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy,
		&result.IsUsed, &result.Nik)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanDAO) PrepareScopeInSalesman(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int) (additionalWhere []string) {
	for key := range scopeLimit {
		var (
			keyDataScope   string
			additionalTemp string
		)

		switch key {
		case constanta.ProvinceDataScope:
			keyDataScope = constanta.ProvinceDataScope
		case constanta.DistrictDataScope:
			keyDataScope = constanta.DistrictDataScope
		case constanta.SalesmanDataScope:
			keyDataScope = constanta.SalesmanDataScope
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

func (input salesmanDAO) getAllScopeConstanta(scopeDB map[string]applicationModel.MappingScopeDB) (output []string) {
	for key, _ := range scopeDB {
		output = append(output, key)
	}

	return
}
