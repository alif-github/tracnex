package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"reflect"
	"strconv"
	"strings"
)

func AuditSystemFieldParamToQuery(index int, userParam map[string]repository.AuditSystemFieldParam) (query string, param []interface{}) {
	if index < 1 {
		index = 1
	}
	for key := range userParam {
		if query != "" {
			query += " AND "
		}
		query += "\"" + key + "\" "
		if userParam[key].IsEqual {
			query += " = $" + strconv.Itoa(index)
			param = append(param, userParam[key].ParamValue)
		} else {
			query += " LIKE $" + strconv.Itoa(index)
			param = append(param, userParam[key].ParamValue)
		}
		index++
	}

	return
}

func ReadQueryToGetAuditField(query string, param []interface{}) (auditField map[string]repository.AuditSystemFieldParam, errorS error) {
	auditField = make(map[string]repository.AuditSystemFieldParam)
	query = strings.Replace(query, "WHERE", "where", -1)
	querySplitWhere := strings.Split(query, "where")
	if len(querySplitWhere) < 1 {
		return
	}

	used := strings.Replace(querySplitWhere[1], "AND", "and", -1)
	splitAND := strings.Split(used, "and")

	for i := 0; i < len(splitAND); i++ {
		split := splitAND[i]
		split = strings.Trim(split, " ")
		word := ""
		last := ""

		var fieldName string
		var isEqual bool
		var paramQuery interface{}
		var index = 1

		for j := 0; j < len(split); j++ {
			if j == len(split) || (last == " " && string(split[j]) != " ") {
				if index == 1 {
					fieldName = word
				} else {
					isEqual = word == "="
				}
				word = ""
				index++
			}
			if string(split[j]) != " " {
				word += string(split[j])
			}
			last = string(split[j])
		}
		if index == 3 {
			word = strings.Replace(word, "$", "", -1)
			var idxParam int
			idxParam, errorS = strconv.Atoi(word)
			if errorS != nil {
				return
			}
			if len(param) < idxParam-1 {
				errorS = errors.New("index param bigger than param length")
				return
			} else {
				paramQuery = param[idxParam-1]
			}
		}
		auditField[fieldName] = repository.AuditSystemFieldParam{
			IsEqual:    isEqual,
			ParamValue: paramQuery,
		}
	}
	return
}

func CheckParamDataType(param []interface{}) (result []interface{}) {
	for i := 0; i < len(param); i++ {
		result = append(result, checkParamDataType(param[i]))
	}
	return
}

func checkParamDataType(param interface{}) interface{} {
	switch param.(type) {
	case sql.NullInt64:
		return param.(sql.NullInt64).Int64
	case sql.NullString:
		return param.(sql.NullString).String
	case sql.NullBool:
		return param.(sql.NullBool).Bool
	case sql.NullTime:
		return param.(sql.NullTime).Time
	case sql.NullInt32:
		return param.(sql.NullInt32).Int32
	case sql.NullFloat64:
		return param.(sql.NullFloat64).Float64
	}
	return param
}

func SearchByParamToQuery(searchByParam []in.SearchByParam, createdBy int64) ([]in.SearchByParam, string, int) {
	var result string
	index := 1
	if len(searchByParam) > 0 || createdBy != 0 {
		result = "WHERE \n"
		if createdBy > 0 {
			if len(searchByParam) == 0 {
				result += "created_by = $" + strconv.Itoa(index)
			} else {
				result += "created_by = $" + strconv.Itoa(index) + " AND "
			}
			index++
		}
		for i := 0; i < len(searchByParam); i++ {
			if searchByParam[i].DataType == "enum" {
				searchByParam[i].SearchKey = "cast( " + searchByParam[i].SearchKey + " AS VARCHAR)"
			}
			if searchByParam[i].SearchOperator == "like" {
				searchByParam[i].SearchKey = "LOWER(" + searchByParam[i].SearchKey + ")"
				searchByParam[i].SearchValue = strings.ToLower(searchByParam[i].SearchValue)
				searchByParam[i].SearchValue = "%" + searchByParam[i].SearchValue + "%"
			}
			operator := searchByParam[i].SearchOperator
			if searchByParam[i].SearchOperator == "eq" {
				operator = "="
			}
			result += " " + searchByParam[i].SearchKey + " " + operator + " $" + strconv.Itoa(index) + " "
			if i < len(searchByParam)-1 {
				result += "AND "
			}
			index++
		}
	}
	if result == "" {
		result += " WHERE deleted = FALSE "
	} else {
		result += " AND deleted = FALSE "
	}
	return searchByParam, result, index
}

func CountOffset(page int, limit int) int {
	return (page - 1) * limit
}

func ArrayStringToArrayInterface(input []string) (output []interface{}) {
	for i := 0; i < len(input); i++ {
		output = append(output, input[i])
	}
	return output
}

func ListDataToInQuery(listScope []interface{}) string {
	result := ""
	for i := 0; i < len(listScope); i++ {
		result += "$" + strconv.Itoa(i+1)
		if i < len(listScope)-1 {
			result += ", "
		}
	}
	return result
}

func ListRangeToInQueryWithStartIndex(rangeData int, startIndex int) (result string, index int) {
	result = ""
	for i := 0; i < rangeData; i++ {
		result += "$" + strconv.Itoa(startIndex)
		if i < rangeData-1 {
			result += ", "
		}

		startIndex++
	}
	return result, startIndex
}

func ListValuesToInsertBulk(inputStruct interface{}, lenValue, index int, inputValue func(inputVal interface{}) []interface{}) (resultQuery string, resultIndex int, resultParam []interface{}) {
	reflectType := reflect.TypeOf(inputStruct)
	reflectValue := reflect.ValueOf(inputStruct)
	var tempQuery string

	switch reflectType.Kind() {
	case reflect.Slice:
		for i := 0; i < reflectValue.Len(); i++ {
			resultQuery += " ( "

			newValue := reflectValue.Index(i).Interface()
			tempParam := inputValue(newValue)
			resultParam = append(resultParam, tempParam...)

			tempQuery, index = ListRangeToInQueryWithStartIndex(lenValue, index)
			resultQuery += tempQuery
			resultQuery += " ) "

			if i < reflectValue.Len()-1 {
				resultQuery += ", "
			}

			resultQuery += fmt.Sprintf("\n")
		}
		break
	}

	resultIndex = index
	return
}

func ArrayInt64ToArrayInterface(input []int64) (output []interface{}) {
	for i := 0; i < len(input); i++ {
		output = append(output, input[i])
	}
	return output
}

func SearchByParamToMapSearch(searchBy []in.SearchByParam) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(searchBy); i++ {
		result[searchBy[i].SearchKey] = searchBy[i].SearchValue
	}
	return result
}

func SearchByParamToQueryWithDefaultField(searchByParam []in.SearchByParam, defaultField DefaultFieldMustCheck, index int) ([]in.SearchByParam, string, int) {
	var result string
	createdBy, _ := defaultField.CreatedBy.Value.(int64)
	if len(searchByParam) > 0 || createdBy != 0 {
		result = " WHERE \n"
		if defaultField.CreatedBy.Value.(int64) > 0 {
			if len(searchByParam) == 0 {
				result += defaultField.CreatedBy.FieldName + " = $" + strconv.Itoa(index)
			} else {
				result += defaultField.CreatedBy.FieldName + " = $" + strconv.Itoa(index) + " AND "
			}
			index++
		}
		for i := 0; i < len(searchByParam); i++ {
			if i == 0 && len(searchByParam) > 1 {
				result += " ( "
			}
			if searchByParam[i].DataType == "enum" {
				searchByParam[i].SearchKey = "cast( " + searchByParam[i].SearchKey + " AS VARCHAR)"
			}
			if searchByParam[i].SearchOperator == "like" {
				searchByParam[i].SearchKey = "LOWER(" + searchByParam[i].SearchKey + ")"
				searchByParam[i].SearchValue = strings.ToLower(searchByParam[i].SearchValue)
				searchByParam[i].SearchValue = "%" + searchByParam[i].SearchValue + "%"
			}
			operator := searchByParam[i].SearchOperator
			if searchByParam[i].SearchOperator == "eq" {
				operator = "="
			}
			result += " " + searchByParam[i].SearchKey + " " + operator + " $" + strconv.Itoa(index) + " "
			if i < len(searchByParam)-1 {
				if searchByParam[i].SearchType == constanta.Search {
					result += "OR "
				} else if searchByParam[i].SearchType == constanta.Filter {
					result += "AND "
				}
			}
			index++
			if i == len(searchByParam)-1 && len(searchByParam) > 1 {
				result += " ) "
			}
		}
	}
	if result == "" {
		result += " WHERE " + defaultField.Deleted.FieldName + " = FALSE "
	} else {
		result += " AND " + defaultField.Deleted.FieldName + " = FALSE "
	}
	return searchByParam, result, index
}

func SearchByParamToJoinQuery(searchByParam []in.SearchByParam, createdBy int64) ([]in.SearchByParam, string, int) {
	return buildSearchByQuery(searchByParam, createdBy)
}

func buildSearchByQuery(searchByParam []in.SearchByParam, createdBy int64) ([]in.SearchByParam, string, int) {
	var result string
	index := 1
	if len(searchByParam) > 0 || createdBy != 0 {
		result = "WHERE \n"
		if createdBy > 0 {
			if len(searchByParam) == 0 {
				result += "created_by = $" + strconv.Itoa(index)
			} else {
				result += "created_by = $" + strconv.Itoa(index) + " AND "
			}
			index++
		}

		for i := 0; i < len(searchByParam); i++ {
			if searchByParam[i].DataType == "enum" {
				searchByParam[i].SearchKey = "cast( " + searchByParam[i].SearchKey + " AS VARCHAR)"
			}
			if searchByParam[i].SearchOperator == "like" {
				searchByParam[i].SearchKey = "LOWER(" + searchByParam[i].SearchKey + ")"
				searchByParam[i].SearchValue = strings.ToLower(searchByParam[i].SearchValue)
				searchByParam[i].SearchValue = "%" + searchByParam[i].SearchValue + "%"
			}
			operator := searchByParam[i].SearchOperator
			if searchByParam[i].SearchOperator == "eq" {
				operator = "="
			}
			result += " " + searchByParam[i].SearchKey + " " + operator + " $" + strconv.Itoa(index) + " "
			if i < len(searchByParam)-1 {
				if searchByParam[i].SearchType == constanta.Search {
					result += "OR "
				} else if searchByParam[i].SearchType == constanta.Filter {
					result += "AND "
				}
			}
			index++
		}
	}

	return searchByParam, result, index
}

func scopeToAddedQuery(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, isView bool, idxStart int, scopeUsed []string) (query string, param []interface{}) {
	for i := 0; i < len(scopeUsed); i++ {
		key := scopeUsed[i]
		listData := scopeLimit[key].([]interface{})
		var temp string
		for i := 0; i < len(listData); i++ {
			data := listData[i].(string)
			if data != "all" {
				if i == 0 {
					value := scopeDB[key].Count
					if isView {
						value = scopeDB[key].View
					}

					// START Handle optional data foreign key
					temp += fmt.Sprintf(" AND %s IN ( $%d, $%d, ", value, idxStart, idxStart+1)
					param = append(param, nil, 0)
					idxStart += 2
					// END Handle optional data foreign key
				}
				id, _ := strconv.Atoi(data)
				if id != 0 {
					param = append(param, int64(id))
					temp += "$" + strconv.Itoa(idxStart) + ", "
					idxStart++
				}
			}
		}
		if temp != "" {
			temp = temp[0 : len(temp)-2]
			temp += ")"
			query += temp + " "
		}
	}

	return
}

func ScopeToAddedQueryView(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int, scopeUsed []string) (query string, param []interface{}) {
	return scopeToAddedQuery(scopeLimit, scopeDB, true, idxStart, scopeUsed)
}

func ScopeToAddedQueryCount(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, idxStart int, scopeUsed []string) (query string, param []interface{}) {
	return scopeToAddedQuery(scopeLimit, scopeDB, false, idxStart, scopeUsed)
}

func PrepareScopeOnDAO(scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB, additionalWhere *string, idxStart int, keyDataScope string, isView bool) {
	var dbParam []interface{}
	var key string

	if isView {
		key = scopeDB[keyDataScope].View
	} else {
		key = scopeDB[keyDataScope].Count
	}

	_, param := ScopeToAddedQueryView(scopeLimit, scopeDB, idxStart, []string{keyDataScope})
	dbParam = append(dbParam, param...)

	if len(dbParam) > 0 {
		*additionalWhere = " " + key + " IN ("
	}

	for idx, valueScope := range dbParam {
		if valueScope == nil {
			if len(dbParam)-(idx+1) == 0 {
				*additionalWhere += "null)"
			} else {
				*additionalWhere += "null,"
			}
		} else {
			idScope, _ := valueScope.(int64)
			if len(dbParam)-(idx+1) == 0 {
				*additionalWhere += strconv.Itoa(int(idScope)) + ")"
			} else {
				*additionalWhere += strconv.Itoa(int(idScope)) + ","
			}
		}
	}
}

func RowsCatchResult(rows *sql.Rows, resultRows func(rws *sql.Rows) (interface{}, errorModel.ErrorModel)) (result []interface{}, err errorModel.ErrorModel) {
	fileName := "DAOUtil.go"
	funcName := "RowsCatchResult"
	var errorS error

	if rows != nil {
		defer func() {
			errorS = rows.Close()
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
				return
			}
		}()
		for rows.Next() {
			var resultTemp interface{}
			resultTemp, err = resultRows(rows)
			if err.Error != nil {
				return
			}

			result = append(result, resultTemp)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func RowCatchResult(row *sql.Row, wrapFunc func(rws *sql.Row) (interface{}, error), fileName string, funcName string) (result interface{}, err errorModel.ErrorModel) {
	var errorS error

	result, errorS = wrapFunc(row)

	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func CreateDollarParamInMultipleRowsDAO(userParamLen int, amountParam int, startParam int, returningFor string) (query string) {
	staticAmountParam := amountParam

	for i := 1; i <= userParamLen; i++ {
		query += "("
		for j := startParam; j <= amountParam; j++ {
			query += " $" + strconv.Itoa(j) + ""
			if amountParam-j != 0 {
				query += ","
			} else {
				query += ")"
			}
		}

		if userParamLen-i != 0 {
			query += ","
		} else {
			if returningFor != "" {
				query += " returning " + returningFor + " "
			}
		}

		startParam += staticAmountParam
		amountParam += staticAmountParam
	}

	return
}

func CheckOwnPermissionAndGetQuery(createdBy int64, query *string, param *[]interface{}, getDefaultData func(int64) DefaultFieldMustCheck, index int) int {
	if createdBy > 0 {
		defaultField := getDefaultData(createdBy)
		queryOwnPermission := " AND " + defaultField.CreatedBy.FieldName + " = $" + strconv.Itoa(index) + " "
		*query += queryOwnPermission
		*param = append(*param, createdBy)
		index += 1
	}
	return index
}

func HandleOptionalParam(userParam []interface{}, param *[]interface{}) {
	for _, item := range userParam {
		currentValue := reflect.ValueOf(item)
		if currentValue.IsZero() {
			*param = append(*param, nil)
		} else {
			*param = append(*param, item)
		}
	}
}
