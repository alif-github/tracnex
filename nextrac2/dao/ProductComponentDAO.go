package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

type productComponentDAO struct {
	AbstractDAO
}

var ProductComponentDAO = productComponentDAO{}.New()

func (input productComponentDAO) New() (output productComponentDAO) {
	output.FileName = "ProductComponentDAO.go"
	output.TableName = "product_component"
	return
}

func (input productComponentDAO) InsertProductComponent(tx *sql.Tx, userParam repository.ProductModel) (id []int64, err errorModel.ErrorModel) {
	funcName := "InsertProductComponent"
	variableParam := 9
	startIndex := 1

	query := "INSERT INTO "+ input.TableName +" " +
		"(product_id, component_id, component_value, " +
		"created_by, created_client, created_at, " +
		"updated_by, updated_client, updated_at) VALUES "

	query += CreateDollarParamInMultipleRowsDAO(len(userParam.ProductComponentModel), variableParam, startIndex, "id")

	var params []interface{}

	for k := 0; k < len(userParam.ProductComponentModel); k++ {
		params = append(params,
			userParam.ID.Int64, userParam.ProductComponentModel[k].ComponentID.Int64, userParam.ProductComponentModel[k].ComponentValue.String,
			userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
			userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time)
	}

	rows, errorS := tx.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	var result []interface{}
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

func (input productComponentDAO) GetListProductComponentByIDProduct(db *sql.DB, userParam repository.ProductModel) (result []repository.GetListProductComponent, err errorModel.ErrorModel) {
	funcName := "GetListProductComponentByIDProduct"

	query := "SELECT " +
		"pc.id, pc.component_id, co.component_name, " +
		"pc.component_value " +
		"FROM "+ input.TableName +" pc " +
		"INNER JOIN "+ComponentDAO.TableName+" co ON co.id = pc.component_id " +
		"WHERE pc.product_id = $1 AND pc.deleted = FALSE "

	param := []interface{}{userParam.ID.Int64}

	rows, errorS := db.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
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
			var temp repository.GetListProductComponent
			errorS = rows.Scan(&temp.ID, &temp.ComponentID, &temp.ComponentName,
				&temp.ComponentValue)
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

func (input productComponentDAO) DeleteProductComponent(db *sql.Tx, userParam repository.ProductModel) (err errorModel.ErrorModel) {
	funcName := "DeleteProductComponent"

	query := "UPDATE "+ input.TableName +" " +
		" SET " +
		" deleted = TRUE, " +
		" updated_by = $1, " +
		" updated_client = $2, " +
		" updated_at = $3 " +
		" WHERE " +
		" product_id = $4"

	param := []interface{}{
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.ID.Int64}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productComponentDAO) UpdateProductComponent(db *sql.Tx, userParam repository.ProductComponentModel, idProduct int64) (err errorModel.ErrorModel) {
	funcName := "UpdateProductComponent"

	query :=
		"UPDATE "+ input.TableName +" " +
			"SET " +
			"component_id = $1, component_value = $2, updated_by = $3, " +
			"updated_at = $4, updated_client = $5 " +
			"WHERE " +
			"id = $6 AND product_id = $7 "

	param := []interface{}{
		userParam.ComponentID.Int64, userParam.ComponentValue.String, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedClient.String, userParam.ID.Int64,
		idProduct,
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	defer stmt.Close()

	return
}

func (input productComponentDAO) DeleteProductComponentMultiple(db *sql.Tx, userParam []repository.ProductComponentModel, inputParam repository.ProductModel) (err errorModel.ErrorModel) {
	funcName := "DeleteProductComponentMultiple"

	query := "UPDATE "+ input.TableName +" " +
		" SET " +
		" deleted = TRUE, " +
		" updated_by = $1, " +
		" updated_client = $2, " +
		" updated_at = $3 " +
		" WHERE " +
		" id IN ("

	for i := 0; i < len(userParam); i++ {
		query += strconv.Itoa(int(userParam[i].ID.Int64))
		if len(userParam)-(i+1) == 0 {
			query += ") "
		} else {
			query += ","
		}
	}

	var param []interface{}
	param = append(param,
		inputParam.UpdatedBy.Int64, inputParam.UpdatedClient.String, inputParam.UpdatedAt.Time)

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

func (input productComponentDAO) GetProductComponentByIDProduct(db *sql.DB, idProduct int64) (result []repository.ProductComponentModel, err errorModel.ErrorModel) {
	funcName := "GetProductComponentByIDProduct"

	query := "SELECT " +
		"pc.component_id, co.component_name, pc.component_value " +
		"FROM "+ input.TableName +" pc " +
		"INNER JOIN "+ComponentDAO.TableName+" co ON pc.component_id = co.id " +
		"WHERE " +
		"pc.product_id = $1 AND " +
		"pc.deleted = FALSE "

	params := []interface{}{idProduct}

	rows, errorS := db.Query(query, params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
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
			var temp repository.ProductComponentModel

			errorS = rows.Scan(&temp.ComponentID, &temp.ComponentName, &temp.ComponentValue)
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

func (input productComponentDAO) CheckComponentIsUsedInProduct(db *sql.DB, userParam repository.ProductComponentModel) (isUsed bool, err errorModel.ErrorModel) {
	funcName := "CheckComponentIsUsedInProduct"

	query := "SELECT " +
		" (CASE WHEN count(id) > 0 THEN TRUE ELSE FALSE END) is_used " +
		" FROM " + input.TableName + " " +
		" WHERE " +
		" component_id = $1 AND deleted = FALSE "

	param := []interface{}{userParam.ComponentID.Int64}

	dbError := db.QueryRow(query, param...).Scan(&isUsed)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productComponentDAO) resultRowsInput(rows *sql.Rows) (idTemp interface{}, err errorModel.ErrorModel) {
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

func (input productComponentDAO) InsertSingleProductComponent(tx *sql.Tx, userParam repository.ProductModel, productComponent repository.ProductComponentModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertSingleProductComponent"

	query := fmt.Sprintf(
		`INSERT INTO %s
		(
			product_id, component_id, component_value,
			created_by, created_client, created_at,
			updated_by, updated_client, updated_at
		) 
			VALUES
		(
		 $1, $2, $3, $4, $5, $6, $7, $8, $9
		) `, input.TableName)

	var params = []interface{}{
		userParam.ID.Int64, productComponent.ComponentID.Int64, productComponent.ComponentValue.String,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
	}

	results := tx.QueryRow(query, params...)

	dbError := results.Scan(&id)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}