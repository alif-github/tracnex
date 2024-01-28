package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strconv"
)

type getListDataDAO struct {
	AbstractDAO
}

var GetListDataDAO = getListDataDAO{}.New()

func (input getListDataDAO) New() (output getListDataDAO) {
	output.FileName = "GetListDataDAO.go"
	return
}

func (input getListDataDAO) GetListDataWithDefaultMustCheck(db *sql.DB, queryParam []interface{}, queryGetList string, userParam in.GetListDataDTO, searchBy []in.SearchByParam, wrap func(rows *sql.Rows) (interface{}, error), additionalWhere string, defaultField DefaultFieldMustCheck) (result []interface{}, err errorModel.ErrorModel) {
	funcName := "GetListDataWithDefaultMustCheck"

	searchBy, tempQuery, index := SearchByParamToQueryWithDefaultField(searchBy, defaultField, len(queryParam)+1)
	query := queryGetList + tempQuery

	if defaultField.CreatedBy.Value.(int64) > 0 {
		queryParam = append(queryParam, defaultField.CreatedBy.Value.(int64))
	}

	if defaultField.Status.IsCheck {
		query += " AND " + defaultField.Status.FieldName + " = 'A' "
	}

	for i := 0; i < len(searchBy); i++ {
		queryParam = append(queryParam, searchBy[i].SearchValue)
	}

	if additionalWhere != "" {
		query += additionalWhere
	}

	if !util.IsStringEmpty(userParam.UpdatedAtStartString) && !util.IsStringEmpty(userParam.UpdatedAtEndString) {
		query += fmt.Sprintf(` AND %s BETWEEN $%s AND $%s ORDER BY %s `, defaultField.UpdatedAt.FieldName, strconv.Itoa(index), strconv.Itoa(index+1), userParam.OrderBy)
		queryParam = append(queryParam, userParam.UpdatedAtStart, userParam.UpdatedAtEnd)
		if userParam.Limit != -99 && userParam.Page != -99 {
			query += fmt.Sprintf(` LIMIT $%s OFFSET $%s `, strconv.Itoa(index+2), strconv.Itoa(index+3))
			queryParam = append(queryParam, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
		}
	} else if !util.IsStringEmpty(userParam.UpdatedAtStartString) {
		query += fmt.Sprintf(` AND %s >= $%s ORDER BY %s `, defaultField.UpdatedAt.FieldName, strconv.Itoa(index), userParam.OrderBy)
		queryParam = append(queryParam, userParam.UpdatedAtStart)
		if userParam.Limit != -99 && userParam.Page != -99 {
			query += fmt.Sprintf(` LIMIT $%s OFFSET $%s `, strconv.Itoa(index+1), strconv.Itoa(index+2))
			queryParam = append(queryParam, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
		}
	} else {
		if userParam.Limit != -99 && userParam.Page != -99 {
			query += fmt.Sprintf(` ORDER BY %s LIMIT $%s OFFSET $%s `, userParam.OrderBy, strconv.Itoa(index), strconv.Itoa(index+1))
			queryParam = append(queryParam, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
		} else {
			if userParam.OrderBy != "" {
				query += fmt.Sprintf(` ORDER BY %s `, userParam.OrderBy)
			}
		}
	}

	rows, errorS := db.Query(query, queryParam...)
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
			temp, errorS := wrap(rows)
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

func (input getListDataDAO) GetCountDataWithDefaultMustCheck(db *sql.DB, queryParam []interface{}, tableName string, searchBy []in.SearchByParam, additionalWhere string, defaultField DefaultFieldMustCheck) (result int, err errorModel.ErrorModel) {
	funcName := "GetCountData"
	query :=
		"SELECT " +
			"	count(" + defaultField.ID.FieldName + ") " +
			"FROM " +
			"	" + tableName + " "
	searchBy, tempQuery, _ := SearchByParamToQueryWithDefaultField(searchBy, defaultField, len(queryParam)+1)
	query += tempQuery

	if defaultField.CreatedBy.Value.(int64) > 0 {
		queryParam = append(queryParam, defaultField.CreatedBy.Value.(int64))
	}

	for i := 0; i < len(searchBy); i++ {
		queryParam = append(queryParam, searchBy[i].SearchValue)
	}

	if defaultField.Status.IsCheck {
		query += " AND " + defaultField.Status.FieldName + " = 'A' "
	}

	if additionalWhere != "" {
		query += additionalWhere
	}

	results := db.QueryRow(query, queryParam...)

	errorS := results.Scan(&result)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

//
//func (input getListDataDAO) GetListDataFromElastic(userParam in.GetListDataDTO, searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64, indexName string) (result string, err errorModel.ErrorModel) {
//	var errs error
//
//	funcName := "GetListDataFromElastic"
//	if isCheckStatus {
//		searchBy = append(searchBy, in.SearchByParam{
//			SearchKey:   "status",
//			SearchValue: "A",
//		})
//	}
//	if createdBy > 0 {
//		searchBy = append(searchBy, in.SearchByParam{
//			SearchKey:   "created_by",
//			SearchValue: strconv.Itoa(int(createdBy)),
//		})
//	}
//
//	result, _, errs = input.doGetListDataFromElastic(userParam, indexName, searchBy)
//	if errs != nil {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
//		return
//	}
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input getListDataDAO) GetListDataFromElasticAndDB(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, query string, indexName string, wrap func(rows *sql.Rows) (interface{}, error), defaultField DefaultFieldMustCheck) (result []interface{}, err errorModel.ErrorModel) {
//	var data string
//
//	data, err = input.GetListDataFromElastic(userParam, searchBy, defaultField.Status.IsCheck, defaultField.CreatedBy.Value.(int64), indexName)
//	if err.Error != nil {
//		return
//	}
//
//	var temp []interface{}
//	_ = json.Unmarshal([]byte(data), &temp)
//	if len(temp) != 0 {
//		listID := input.GetIDFromList(temp)
//		query += " WHERE " + defaultField.ID.FieldName + " in (" + ListDataToInQuery(listID) + ") "
//
//		result, err = GetListDataDAO.GetListDataAfterElastic(db, query, listID, wrap, defaultField)
//		if err.Error != nil {
//			return
//		}
//	}
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input getListDataDAO) GetListDataAfterElastic(db *sql.DB, queryGetList string, listID []interface{}, wrap func(rows *sql.Rows) (interface{}, error), defaultField DefaultFieldMustCheck) (result []interface{}, err errorModel.ErrorModel) {
//	funcName := "GetListDataAfterElastic"
//	var queryParam []interface{}
//
//	query := queryGetList
//	queryParam = listID
//
//	if defaultField.Status.IsCheck {
//		query += " AND " + defaultField.Status.FieldName + " = 'A' "
//	}
//
//	if defaultField.CreatedBy.Value.(int64) > 0 {
//		query += "AND created_by = $" + strconv.Itoa(len(listID)+1)
//		queryParam = append(queryParam, defaultField.CreatedBy.Value.(int64))
//	}
//
//	queryParam = append(queryParam)
//	rows, errorS := db.Query(query, queryParam...)
//	if errorS != nil {
//		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//	}
//	if rows != nil {
//		defer func() {
//			errorS = rows.Close()
//			if errorS != nil {
//				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//				return
//			}
//		}()
//		for rows.Next() {
//			temp, errorS := wrap(rows)
//			if errorS != nil {
//				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//				return
//			}
//			result = append(result, temp)
//		}
//	} else {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}

func (input getListDataDAO) GetIDFromList(list []interface{}) []interface{} {
	var result []interface{}
	for i := 0; i < len(list); i++ {
		result = append(result, list[i].(map[string]interface{})["id"])
	}
	return result
}

func (input getListDataDAO) GetListData(db *sql.DB, queryGetList string, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, wrap func(rows *sql.Rows) (interface{}, error), additionalWhere string) (result []interface{}, err errorModel.ErrorModel) {
	funcName := "GetListData"
	var queryParam []interface{}

	searchBy, tempQuery, index := SearchByParamToQuery(searchBy, createdBy)
	query := queryGetList + tempQuery

	for i := 0; i < len(searchBy); i++ {
		queryParam = append(queryParam, searchBy[i].SearchValue)
	}

	if createdBy > 0 {
		queryParam = append(queryParam, createdBy)
	}

	if additionalWhere != "" {
		query += " AND " + additionalWhere
	}

	query += "ORDER BY " + userParam.OrderBy + " " +
		"LIMIT $" + strconv.Itoa(index) + " OFFSET $" + strconv.Itoa(index+1)

	queryParam = append(queryParam, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
	rows, errorS := db.Query(query, queryParam...)
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
			temp, errorS := wrap(rows)
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

func (input getListDataDAO) GetDataRows(db *sql.DB, query string, wrap func(rows *sql.Rows) (interface{}, error), queryParam []interface{}) (result []interface{}, err errorModel.ErrorModel) {
	funcName := "GetDataRows"
	rows, errorS := db.Query(query, queryParam...)
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
			temp, errorS := wrap(rows)
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

func (input getListDataDAO) GetCountDataWithUpdatedAtAndDefaultMustCheck(db *sql.DB, queryParam []interface{}, tableName string, searchBy []in.SearchByParam, additionalWhere string, defaultField DefaultFieldMustCheck, userParam in.GetListDataDTO) (result int, err errorModel.ErrorModel) {
	funcName := "GetCountDataWithUpdatedAtAndDefaultMustCheck"
	query :=
		"SELECT " +
			"	count(" + defaultField.ID.FieldName + ") " +
			"FROM " +
			"	" + tableName + " "
	searchBy, tempQuery, index := SearchByParamToQueryWithDefaultField(searchBy, defaultField, len(queryParam)+1)
	query += tempQuery

	if defaultField.CreatedBy.Value.(int64) > 0 {
		queryParam = append(queryParam, defaultField.CreatedBy.Value.(int64))
	}

	for i := 0; i < len(searchBy); i++ {
		queryParam = append(queryParam, searchBy[i].SearchValue)
	}

	if defaultField.Status.IsCheck {
		query += " AND " + defaultField.Status.FieldName + " = 'A' "
	}

	if additionalWhere != "" {
		query += additionalWhere
	}

	if util.IsStringEmpty(defaultField.UpdatedAt.FieldName) {
		defaultField.UpdatedAt.FieldName = "updated_at"
	}

	if !util.IsStringEmpty(userParam.UpdatedAtStartString) {
		query += " AND " + defaultField.UpdatedAt.FieldName + " BETWEEN $" + strconv.Itoa(index) + " AND $" + strconv.Itoa(index+1)
		queryParam = append(queryParam, userParam.UpdatedAtStart, userParam.UpdatedAtEnd)
	}

	results := db.QueryRow(query, queryParam...)

	errorS := results.Scan(&result)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input getListDataDAO) GetCountData(db *sql.DB, queryParam []interface{}, query string, searchBy []in.SearchByParam, additionalWhere string, defaultField DefaultFieldMustCheck) (result int, err errorModel.ErrorModel) {
	funcName := "GetCountData"

	searchBy, tempQuery, _ := SearchByParamToQueryWithDefaultField(searchBy, defaultField, len(queryParam)+1)
	query += tempQuery

	if defaultField.CreatedBy.Value.(int64) > 0 {
		queryParam = append(queryParam, defaultField.CreatedBy.Value.(int64))
	}

	for i := 0; i < len(searchBy); i++ {
		queryParam = append(queryParam, searchBy[i].SearchValue)
	}

	if defaultField.Status.IsCheck {
		query += " AND " + defaultField.Status.FieldName + " = 'A' "
	}

	if additionalWhere != "" {
		query += additionalWhere
	}

	results := db.QueryRow(query, queryParam...)

	errorS := results.Scan(&result)
	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
