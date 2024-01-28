package dao

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type getListJoinDataDAO struct {
	AbstractDAO
	Query           string
	Table           string
	Join            []JoinTable
	AdditionalWhere []string
	GroupBy         string
}

type JoinTable struct {
	Type                string
	Table               string
	Key                 string
	Value               string
	Alias               string
	MustNotCheckDeleted bool
}

var GetListJoinDataDAO = getListJoinDataDAO{}.New()
var queryDeleted = ".deleted = FALSE "

func (input getListJoinDataDAO) New() (output getListJoinDataDAO) {
	output.FileName = "GetListJoinDataDAO.go"
	return
}

func (input *getListJoinDataDAO) NewJoinData() *getListJoinDataDAO {
	input.FileName = GetListJoinDataDAO.FileName
	input.Query = ""
	input.Table = ""
	input.Join = []JoinTable{}
	input.AdditionalWhere = []string{}

	return input
}

func (input *getListJoinDataDAO) SetTable(table string) *getListJoinDataDAO {
	input.Table = table
	return input
}

func (input *getListJoinDataDAO) SetQuery(query string) *getListJoinDataDAO {
	input.Query = query
	return input
}

func (input *getListJoinDataDAO) SetWhere(field string, value string) *getListJoinDataDAO {
	input.AdditionalWhere = append(input.AdditionalWhere, field+" = "+"'"+value+"'")
	return input
}

func (input *getListJoinDataDAO) SetDateBetween(field string, start string, end string) *getListJoinDataDAO {
	input.AdditionalWhere = append(input.AdditionalWhere, field+" BETWEEN "+"'"+start+"'"+" AND "+"'"+end+"'")
	return input
}

func (input *getListJoinDataDAO) SetWhereNot(field string, value string) *getListJoinDataDAO {
	input.AdditionalWhere = append(input.AdditionalWhere, field+" != "+"'"+value+"'")
	return input
}

func (input *getListJoinDataDAO) SetWhereAdditional(additionalWhere string) *getListJoinDataDAO {
	input.AdditionalWhere = append(input.AdditionalWhere, additionalWhere)
	return input
}

func (input *getListJoinDataDAO) GetCountJoinData(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	funcName := "GetCountJoinData"

	input.Query += " " + input.buildJoinQuery()
	searchBy, tempQuery, index := SearchByParamToJoinQuery(searchBy, createdBy)

	if index == 1 {
		input.Query += fmt.Sprintf(`WHERE %s`, input.Table+queryDeleted)
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}
				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	} else {
		input.Query += tempQuery
		input.Query += " AND " + input.Table + queryDeleted
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}
				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	}

	if len(input.AdditionalWhere) > 0 {
		strWhere := " AND " + strings.Join(input.AdditionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		input.Query += strWhere
	}

	var queryParam []interface{}
	for i := 0; i < len(searchBy); i++ {
		queryParam = append(queryParam, searchBy[i].SearchValue)
	}

	if createdBy > 0 {
		queryParam = append(queryParam, createdBy)
	}

	if input.GroupBy != "" {
		input.Query += input.GroupBy
	}

	results := db.QueryRow(input.Query, queryParam...)
	errorS := results.Scan(&result)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *getListJoinDataDAO) GetCountJoinDataWithoutDeleted(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (result int, err errorModel.ErrorModel) {
	funcName := "GetCountJoinData"

	input.Query += " " + input.buildJoinQuery()
	searchBy, tempQuery, index := SearchByParamToJoinQuery(searchBy, createdBy)

	if index == 1 {
		input.Query += fmt.Sprintf(`WHERE %s`, input.Table+queryDeleted)
		//if len(input.Join) > 0 {
		//	for _, item := range input.Join {
		//		table := item.Table
		//		if item.Alias != "" {
		//			table = item.Alias
		//		}
		//		if !item.MustNotCheckDeleted {
		//			input.Query += " AND " + table + queryDeleted
		//		}
		//	}
		//}
	} else {
		input.Query += tempQuery
	}

	if len(input.AdditionalWhere) > 0 {
		strWhere := " AND " + strings.Join(input.AdditionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		input.Query += strWhere
	}

	var queryParam []interface{}
	for i := 0; i < len(searchBy); i++ {
		queryParam = append(queryParam, searchBy[i].SearchValue)
	}

	if createdBy > 0 {
		queryParam = append(queryParam, createdBy)
	}

	results := db.QueryRow(input.Query, queryParam...)

	errorS := results.Scan(&result)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *getListJoinDataDAO) GetListJoinData(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, mappingFunc func(rows *sql.Rows) (interface{}, error)) (results []interface{}, err errorModel.ErrorModel) {
	var (
		funcName  = "GetListJoinData"
		params    []interface{}
		tempQuery string
		index     int
	)

	input.Query += " " + input.buildJoinQuery()
	searchBy, tempQuery, index = SearchByParamToJoinQuery(searchBy, createdBy)

	if index == 1 {
		input.Query += fmt.Sprintf(`WHERE %s`, input.Table+queryDeleted)
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}

				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	} else {
		input.Query += tempQuery
		input.Query += " AND " + input.Table + queryDeleted
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}

				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	}

	if len(input.AdditionalWhere) > 0 {
		strWhere := " AND " + strings.Join(input.AdditionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		input.Query += strWhere
	}

	for i := 0; i < len(searchBy); i++ {
		params = append(params, searchBy[i].SearchValue)
	}

	if createdBy > 0 {
		params = append(params, createdBy)
	}

	if input.GroupBy != "" {
		input.Query += input.GroupBy
	}

	if userParam.Limit != -99 && userParam.Page != -99 {
		input.Query += fmt.Sprintf(` ORDER BY %s LIMIT $%s OFFSET $%s `, userParam.OrderBy, strconv.Itoa(index), strconv.Itoa(index+1))
		params = append(params, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
	} else {
		input.Query += fmt.Sprintf(` ORDER BY %s `, userParam.OrderBy)
	}

	rows, errorS := db.Query(input.Query, params...)
	if errorS != nil {
		return results, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			temp, errorS := mappingFunc(rows)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			results = append(results, temp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	return
}

func (input *getListJoinDataDAO) GetListJoinDataWithoutDeleted(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, mappingFunc func(rows *sql.Rows) (interface{}, error)) (results []interface{}, err errorModel.ErrorModel) {
	funcName := "GetListJoinDataWithoutDeleted"
	var params []interface{}
	input.Query += " " + input.buildJoinQuery()
	searchBy, tempQuery, index := SearchByParamToJoinQuery(searchBy, createdBy)

	if index == 1 {
		input.Query += fmt.Sprintf(`WHERE %s`, input.Table+queryDeleted)
		//if len(input.Join) > 0 {
		//	for _, item := range input.Join {
		//		table := item.Table
		//		if item.Alias != "" {
		//			table = item.Alias
		//		}
		//		input.Query += " AND " + table + queryDeleted
		//	}
		//}
	} else {
		input.Query += tempQuery
	}

	if len(input.AdditionalWhere) > 0 {
		strWhere := " AND " + strings.Join(input.AdditionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		input.Query += strWhere
	}

	for i := 0; i < len(searchBy); i++ {
		params = append(params, searchBy[i].SearchValue)
	}

	if createdBy > 0 {
		params = append(params, createdBy)
	}

	input.Query += fmt.Sprintf(
		` ORDER BY  %s
		LIMIT $%s OFFSET $%s `,
		userParam.OrderBy, strconv.Itoa(index), strconv.Itoa(index+1))

	params = append(params, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
	rows, errorS := db.Query(input.Query, params...)
	if errorS != nil {
		return results, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			temp, errorS := mappingFunc(rows)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			results = append(results, temp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	return
}

func (input *getListJoinDataDAO) buildJoinQuery() (results string) {
	if len(input.Join) > 0 {
		for _, item := range input.Join {
			results += item.Type + " JOIN " + item.Table

			if item.Alias != "" {
				results += " AS " + item.Alias
			}

			results += " ON " + item.Key + " = " + item.Value + " "
		}
	}
	return
}

func (input *getListJoinDataDAO) InnerJoin(table string, key string, value string) {
	input.setJoin(table, "INNER", key, value, "", false)
}

func (input *getListJoinDataDAO) InnerJoinAlias(table string, alias string, key string, value string) {
	input.setJoin(table, "INNER", key, value, alias, false)
}

func (input *getListJoinDataDAO) InnerJoinWithoutDeleted(table string, key string, value string) {
	input.setJoin(table, "INNER", key, value, "", true)
}

func (input *getListJoinDataDAO) InnerJoinAliasWithoutDeleted(table string, alias string, key string, value string) {
	input.setJoin(table, "INNER", key, value, alias, true)
}

func (input *getListJoinDataDAO) LeftJoin(table string, key string, value string) {
	input.setJoin(table, "LEFT", key, value, "", false)
}

func (input *getListJoinDataDAO) LeftJoinAlias(table string, alias string, key string, value string) {
	input.setJoin(table, "LEFT", key, value, alias, false)
}

func (input *getListJoinDataDAO) RightJoin(table string, key string, value string) {
	input.setJoin(table, "RIGHT", key, value, "", false)
}

func (input *getListJoinDataDAO) RightJoinAlias(table string, alias string, key string, value string) {
	input.setJoin(table, "RIGHT", key, value, alias, false)
}

func (input *getListJoinDataDAO) LeftJoinWithoutDeleted(table string, key string, value string) {
	input.setJoin(table, "LEFT", key, value, "", true)
}

func (input *getListJoinDataDAO) LeftJoinAliasWithoutDeleted(table string, alias string, key string, value string) {
	input.setJoin(table, "LEFT", key, value, alias, true)
}

func (input *getListJoinDataDAO) setJoin(table string, joinType string, key string, value string, alias string, mustNotCheckDeleted bool) {
	input.Join = append(input.Join, JoinTable{
		Type:                joinType,
		Table:               table,
		Key:                 key,
		Value:               value,
		Alias:               alias,
		MustNotCheckDeleted: mustNotCheckDeleted,
	})
}

//----------------- New Function On Phase 2 -----------------------

func (input *getListJoinDataDAO) GetListJoinDataAndFreeAdditionalWhere(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, mappingFunc func(rows *sql.Rows) (interface{}, error)) (results []interface{}, err errorModel.ErrorModel) {
	var (
		funcName  = "GetListJoinDataAndFreeAdditionalWhere"
		params    []interface{}
		index     int
		tempQuery string
	)

	input.Query += " " + input.buildJoinQuery()
	searchBy, tempQuery, index = SearchByParamToJoinQuery(searchBy, createdBy)

	if index == 1 {
		input.Query += fmt.Sprintf(`WHERE %s`, input.Table+queryDeleted)
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}
				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	} else {
		input.Query += tempQuery
		input.Query += " AND " + input.Table + queryDeleted
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}
				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	}

	if len(input.AdditionalWhere) > 0 {
		for _, valueAdditionalWhere := range input.AdditionalWhere {
			input.Query += valueAdditionalWhere + " "
		}
	}

	for i := 0; i < len(searchBy); i++ {
		params = append(params, searchBy[i].SearchValue)
	}

	if createdBy > 0 {
		params = append(params, createdBy)
	}

	if userParam.Limit != -99 && userParam.Page != -99 {
		input.Query += fmt.Sprintf(` ORDER BY %s LIMIT $%s OFFSET $%s `, userParam.OrderBy, strconv.Itoa(index), strconv.Itoa(index+1))
		params = append(params, userParam.Limit, CountOffset(userParam.Page, userParam.Limit))
	} else {
		input.Query += fmt.Sprintf(` ORDER BY %s `, userParam.OrderBy)
	}

	rows, errorS := db.Query(input.Query, params...)
	if errorS != nil {
		return results, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			temp, errorS := mappingFunc(rows)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			results = append(results, temp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}
	return
}

func (input *getListJoinDataDAO) GetListJoinDataWithoutPagination(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64, mappingFunc func(rows *sql.Rows) (interface{}, error)) (results []interface{}, err errorModel.ErrorModel) {
	var (
		params []interface{}
	)

	params = input.queryBuilder(searchBy, createdBy)
	input.Query += fmt.Sprintf(` ORDER BY  %s `, userParam.OrderBy)
	return input.hitDBQueryGetListData(db, params, mappingFunc)
}

func (input *getListJoinDataDAO) queryBuilder(searchBy []in.SearchByParam, createdBy int64) (params []interface{}) {
	var (
		index     int
		tempQuery string
	)

	input.Query += " " + input.buildJoinQuery()
	searchBy, tempQuery, index = SearchByParamToJoinQuery(searchBy, createdBy)

	if index == 1 {
		input.Query += fmt.Sprintf(`WHERE %s`, input.Table+queryDeleted)
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}
				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	} else {
		input.Query += tempQuery
		input.Query += " AND " + input.Table + queryDeleted
		if len(input.Join) > 0 {
			for _, item := range input.Join {
				table := item.Table
				if item.Alias != "" {
					table = item.Alias
				}
				if !item.MustNotCheckDeleted {
					input.Query += " AND " + table + queryDeleted
				}
			}
		}
	}

	if len(input.AdditionalWhere) > 0 {
		strWhere := " AND " + strings.Join(input.AdditionalWhere, " AND ")
		strWhere = strings.TrimRight(strWhere, " AND ")
		input.Query += strWhere
	}

	for i := 0; i < len(searchBy); i++ {
		params = append(params, searchBy[i].SearchValue)
	}

	if createdBy > 0 {
		params = append(params, createdBy)
	}

	return
}

func (input *getListJoinDataDAO) hitDBQueryGetListData(db *sql.DB, params []interface{}, mappingFunc func(rows *sql.Rows) (interface{}, error)) (results []interface{}, err errorModel.ErrorModel) {
	funcName := "HitDBQueryGetListData"
	rows, errorS := db.Query(input.Query, params...)
	if errorS != nil {
		return results, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			temp, errorS := mappingFunc(rows)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			results = append(results, temp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
