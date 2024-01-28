package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type enumDAO struct {
	AbstractDAO
}

var EnumDAO = enumDAO{}.New()

func (input enumDAO) New() (output enumDAO) {
	output.FileName = "EnumDAO.go"
	output.TableName = "pg_enum"

	return
}

func (input enumDAO) GetListEnumLabel(db *sql.DB, inputStruct in.EnumRequest, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	funcName := "GetListEnumLabel"
	query := fmt.Sprintf(`
		SELECT enumlabel
			FROM %s 
		WHERE enumtypid = '%s'::regtype;`, input.TableName, inputStruct.Type,
	)

	rows, errS := db.Query(query)
	if errS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
		return
	}
	defer rows.Close()

	// Loop through the rows and retrieve the data
	for rows.Next() {
		var temp repository.EnumModel
		if errS = rows.Scan(&temp.EnumLabel); errS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
			return
		}

		result = append(result, temp)
	}

	// Handle potential errors
	if errS = rows.Err(); errS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
		return
	}

	return
}
