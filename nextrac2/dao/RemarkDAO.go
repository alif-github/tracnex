package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type remarkDAO struct {
	AbstractDAO
}

var RemarkDAO = remarkDAO{}.New()

func (input remarkDAO) New() (output remarkDAO) {
	output.FileName = "RemarkDAO.go"
	output.TableName = "remark"
	return
}

func (input remarkDAO) GetParentFromChild(db *sql.DB, userParam repository.RemarkModel) (result []repository.RemarkModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetParentFromChild"
		params   []interface{}
	)

	query := fmt.Sprintf(
		`WITH RECURSIVE whos_your_daddy AS (
		SELECT 
			id, name, value, parent_id, level
		FROM %s 
		WHERE deleted = FALSE AND id = $1
		UNION ALL
		SELECT 
			r.id, r.name, r.value, r.parent_id, r.level
		FROM 
			%s r
		INNER JOIN whos_your_daddy p ON r.id = p.parent_id
	) SELECT * FROM whos_your_daddy`, input.TableName, input.TableName)

	params = append(params, userParam.ID.Int64)

	rows, dbError := db.Query(query, params...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	tempResult, err := RowsCatchResult(rows, func(rws *sql.Rows) (tempOutput interface{}, err errorModel.ErrorModel) {
		var (
			errorS   error
			tempData repository.RemarkModel
		)

		errorS = rows.Scan(&tempData.ID, &tempData.Name, &tempData.Value, &tempData.ParentID, &tempData.Level)
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
		result = append(result, itemResult.(repository.RemarkModel))
	}

	return
}

func (input remarkDAO) GetRemarkFamilyTree(db *sql.DB, userParam repository.RemarkModel) (result []repository.RemarkModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetRemarkFamilyTree"
		params   []interface{}
	)

	query := fmt.Sprintf(
		`WITH RECURSIVE pedigree AS (
		SELECT 
			id, name, value, parent_id
		FROM %s 
		WHERE deleted = FALSE AND id = $1
		UNION ALL
		SELECT 
			r.id, r.name, r.value, r.parent_id
		FROM 
			%s r
		INNER JOIN pedigree p ON p.id = r.parent_id
	) SELECT * FROM pedigree`, input.TableName, input.TableName)

	params = append(params, userParam.ID.Int64)

	rows, dbError := db.Query(query, params...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	tempResult, err := RowsCatchResult(rows, func(rws *sql.Rows) (tempOutput interface{}, err errorModel.ErrorModel) {
		var (
			errorS   error
			tempData repository.RemarkModel
		)

		errorS = rows.Scan(&tempData.ID, &tempData.Name, &tempData.Value, &tempData.ParentID)
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
		result = append(result, itemResult.(repository.RemarkModel))
	}

	return
}

func (input remarkDAO) GetRemarkByID(db *sql.DB, userParam repository.RemarkModel) (result repository.RemarkModel, err errorModel.ErrorModel) {
	funcName := "GetRemarkByID"
	params := []interface{}{userParam.ID.Int64}

	query := fmt.Sprintf(
		`SELECT 
		id, name, value, level, parent_id 
	FROM %s	
	WHERE id = $1 AND deleted = FALSE `, input.TableName)

	rows := db.QueryRow(query, params...)

	tempResult, err := RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp repository.RemarkModel
		dbError := rws.Scan(
			&temp.ID, &temp.Name, &temp.Value,
			&temp.Level, &temp.ParentID,
		)
		return temp, dbError
	}, input.FileName, funcName)

	if err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.RemarkModel)
	}
	return
}
