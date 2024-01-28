package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strconv"
)

type trackerDAO struct {
	AbstractDAO
}

var TrackerDAO = trackerDAO{}.New()

func (input trackerDAO) New() (output trackerDAO) {
	output.FileName = "TrackerDAO.go"
	return
}

func (input trackerDAO) GetListTrackerOnRedmine(db *sql.DB, inputStruct in.GetListDataDTO) (result []interface{}, err errorModel.ErrorModel) {
	var (
		funcName  = "GetListTrackerOnRedmine"
		query     string
		param     []interface{}
		index     int
		queryTemp string
	)

	query = fmt.Sprint(`SELECT "name" FROM trackers `)

	index = 1
	queryTemp = fmt.Sprintf(` ORDER BY "name" ASC `)
	if inputStruct.Limit != -99 && inputStruct.Page != -99 {
		queryTemp = fmt.Sprintf(` ORDER BY "name" ASC LIMIT $%s OFFSET $%s `, strconv.Itoa(index), strconv.Itoa(index+1))
		param = append(param, inputStruct.Limit, CountOffset(inputStruct.Page, inputStruct.Limit))
	}

	query += queryTemp
	rows, errorS := db.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return RowsCatchResult(rows, func(rws *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
		var (
			dbError error
			tracker string
		)

		dbError = rws.Scan(&tracker)
		if dbError != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return
		}

		resultTemp = tracker
		return
	})
}
