package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type nexmileParameterDAO struct {
	AbstractDAO
}

var NexmileParameterDAO = nexmileParameterDAO{}.New()

func (input nexmileParameterDAO) New() (output nexmileParameterDAO) {
	output.FileName = "NexmileParameterDAO.go"
	output.TableName = "nexmile_parameter"
	return
}

func (input nexmileParameterDAO) InsertMultiNexmileParameter(tx *sql.Tx, userParam repository.NexmileParameterModelMap) (id []int64, err errorModel.ErrorModel) {
	var (
		funcName                  = "InsertNexmileParameter"
		parameterNexmileParameter = 12
		jVar                      = 1
		param, result             []interface{}
	)

	query := fmt.Sprintf(`INSERT INTO %s 
			(parameter_id, unique_id_1, unique_id_2, 
			client_id, parameter_value, client_type_id, 
			created_by, created_client, created_at, 
			updated_by, updated_client, updated_at) 
			VALUES `, input.TableName)

	query += CreateDollarParamInMultipleRowsDAO(len(userParam.ParameterData), parameterNexmileParameter, jVar, "id")
	for _, itemParameterData := range userParam.ParameterData {
		param = append(param, itemParameterData.ParameterID.String, userParam.UniqueID1.String)
		HandleOptionalParam([]interface{}{userParam.UniqueID2.String}, &param)

		param = append(param, userParam.CLientID.String, itemParameterData.ParameterValue.String)

		param = append(param,
			userParam.ClientTypeID.Int64, userParam.CreatedBy.Int64, userParam.CreatedClient.String,
			userParam.CreatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
			userParam.UpdatedAt.Time)
	}

	rows, errorS := tx.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	result, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResult := range result {
		id = append(id, itemResult.(int64))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterDAO) resultRowsInput(rows *sql.Rows) (idTemp interface{}, err errorModel.ErrorModel) {
	funcName := "resultRowsInput"
	var errorS error
	var id int64

	errorS = rows.Scan(&id)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	idTemp = id
	return
}

func (input nexmileParameterDAO) resultRowsInputGetParameter(rows *sql.Rows) (dataTemp interface{}, err errorModel.ErrorModel) {
	funcName := "resultRowsInputGetParameter"
	var errorS error
	var temp repository.ParameterValueModel

	errorS = rows.Scan(&temp.ParameterID.String, &temp.ParameterValue.String)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	dataTemp = temp
	return
}

func (input nexmileParameterDAO) GetFieldForUserValidationUserNexmileParameter(db *sql.DB, userParam repository.NexmileParameterModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	funcName := "GetFieldForUserValidationUserNexmileParameter"
	var tempResult interface{}

	query := fmt.Sprintf(`SELECT urd.id, urd.user_id, urd.password, urd.unique_id_1, urd.unique_id_2, urd.product_valid_thru 
			FROM %s urd 
			JOIN %s pkce ON urd.client_id = pkce.client_id 
			JOIN %s cm ON pkce.parent_client_id = cm.client_id  
			WHERE 
			urd.client_id = $1 AND urd.user_id = $2 AND urd.status = 'A' 
			AND cm.client_type_id = $3 AND urd.deleted = FALSE `, UserRegistrationDetailDAO.TableName, PKCEClientMappingDAO.TableName, ClientMappingDAO.TableName)

	params := []interface{}{
		userParam.ClientID.String,
		userParam.UserID.String,
		userParam.ClientTypeID.Int64,
	}

	dbResult := db.QueryRow(query, params...)

	if tempResult, err = RowCatchResult(dbResult, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserRegistrationDetailModel
		dbError := dbResult.Scan(
			&temp.ID.Int64,
			&temp.UserID.String,
			&temp.Password.String,
			&temp.UniqueID1.String,
			&temp.UniqueID2.String,
			&temp.ProductValidThru.Time,
		)
		return temp, dbError
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.UserRegistrationDetailModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterDAO) GetFieldNexmileParameter(db *sql.DB, userParam repository.NexmileParameterModel) (result []repository.ParameterValueModel, err errorModel.ErrorModel) {
	funcName := "GetFieldNexmileParameter"

	// Client ID dan Client Type itu punya parent
	query := fmt.Sprintf(`SELECT parameter_id, parameter_value 
	FROM %s 
	WHERE 
		unique_id_1 = $1 AND unique_id_2 = $2 AND client_type_id = $3 AND 
		client_id = $4 AND deleted = FALSE `, input.TableName)

	params := []interface{}{
		userParam.UniqueID1.String,
		userParam.UniqueID2.String,
		userParam.ClientTypeID.Int64,
		userParam.ClientID.String,
	}

	rows, errorDB := db.Query(query, params...)
	if errorDB != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	var resultTemp []interface{}
	if resultTemp, err = RowsCatchResult(rows, input.resultRowsInputGetParameter); err.Error != nil {
		return
	}

	for _, itemResult := range resultTemp {
		result = append(result, itemResult.(repository.ParameterValueModel))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterDAO) CheckNexmileParameterExistOrNot(db *sql.DB, userParam repository.NexmileParameterModel) (result []repository.ParameterValueModel, err errorModel.ErrorModel) {
	var (
		funcName   = "CheckNexmileParameterExistOrNot"
		lastIndex  = 3
		query      string
		params     []interface{}
		resultTemp []interface{}
	)

	query = fmt.Sprintf(`SELECT 
		id, parameter_id FROM nexmile_parameter 
		WHERE 
		client_id = $1 AND unique_id_1 = $2 AND client_type_id = $3 AND 
		deleted = false AND parameter_id IN ( `)

	params = append(params, userParam.ClientID.String, userParam.UniqueID1.String, userParam.ClientTypeID.Int64)
	for idx, itemParameterData := range userParam.ParameterData {
		lastIndex++
		if len(userParam.ParameterData)-(idx+1) == 0 {
			query += fmt.Sprintf(` $%d) `, lastIndex)
			params = append(params, itemParameterData.ParameterID.String)
			continue
		}

		query += fmt.Sprintf(` $%d, `, lastIndex)
		params = append(params, itemParameterData.ParameterID.String)
	}

	if userParam.UniqueID2.String != "" {
		query += fmt.Sprintf(` AND unique_id_2 = $%d `, lastIndex+1)
		params = append(params, userParam.UniqueID2.String)
	}

	rows, errorDB := db.Query(query, params...)
	if errorDB != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	resultTemp, err = RowsCatchResult(rows, func(rws *sql.Rows) (interface{}, errorModel.ErrorModel) {
		var (
			errors error
			temp   repository.ParameterValueModel
		)

		errors = rows.Scan(&temp.ID, &temp.ParameterID)
		if errors != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errors)
			return nil, err
		}

		return temp, errorModel.GenerateNonErrorModel()
	})

	if err.Error != nil {
		return
	}

	for _, itemResult := range resultTemp {
		result = append(result, itemResult.(repository.ParameterValueModel))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterDAO) UpdateNexmileParameter(db *sql.Tx, userParam repository.ParameterValueModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateNexmileParameter"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET
		parameter_value = $1, updated_at = $2, updated_client = $3,
		updated_by = $4
		WHERE 
		id = $5 AND deleted = FALSE `, input.TableName)

	param := []interface{}{
		userParam.ParameterValue.String, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.UpdatedBy.Int64, userParam.ID.Int64}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	return errorModel.GenerateNonErrorModel()
}
