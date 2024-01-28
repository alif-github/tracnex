package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strings"
)

type sprintDAO struct {
	AbstractDAO
}

var SprintDAO = sprintDAO{}.New()

func (input sprintDAO) New() (output sprintDAO) {
	output.FileName = "SprintDAO.go"
	return
}

func (input sprintDAO) GetListSprintOnRedmine(db *sql.DB) (sprint []string, err errorModel.ErrorModel) {
	var (
		funcName       = "GetListSprintOnRedmine"
		query          string
		param          []interface{}
		idSprint       = constanta.IDSprintOnRedmineCustomFields
		possibleValues string
	)

	query = fmt.Sprintf(` 
		SELECT possible_values FROM custom_fields WHERE id = %d `,
		idSprint)

	rows := db.QueryRow(query, param...)
	dbError := rows.Scan(&possibleValues)
	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	splitData := strings.Split(possibleValues, "\n")
	for _, itemSplitData := range splitData {
		if itemSplitData == "---" || itemSplitData == "" {
			continue
		}

		d := strings.TrimLeft(itemSplitData, "- ")
		sprint = append(sprint, d)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input sprintDAO) ReArrangeDataSprint(db *sql.DB, inputStruct in.GetListDataDTO, dataInput []string) (result []string, err errorModel.ErrorModel) {
	var (
		funcName = "ReArrangeDataSprint"
		query    string
		values   string
	)

	for i := 0; i < len(dataInput); i++ {
		d := dataInput[i]
		values += fmt.Sprintf(`('%s')`, d)
		if len(dataInput)-(i+1) != 0 {
			values += ","
		}
	}

	query = fmt.Sprintf(`
		WITH result_sprint(sprint) AS (VALUES %s)
		SELECT sprint AS sprint_coll 
		FROM result_sprint 
		ORDER BY sprint DESC LIMIT $1 OFFSET $2 `,
		values)

	params := []interface{}{inputStruct.Limit, CountOffset(inputStruct.Page, inputStruct.Limit)}

	rows, dbError := db.Query(query, params...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	var tempResult []interface{}
	tempResult, err = RowsCatchResult(rows, func(rws *sql.Rows) (tempOutput interface{}, err errorModel.ErrorModel) {
		var (
			errorS   error
			tempData string
		)

		errorS = rows.Scan(&tempData)
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}

		tempOutput = tempData
		return
	})

	if err.Error != nil {
		return
	}

	for _, itemResult := range tempResult {
		result = append(result, itemResult.(string))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
