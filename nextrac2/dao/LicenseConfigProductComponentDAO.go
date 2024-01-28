package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type licenseConfigProductComponentDAO struct {
	AbstractDAO
}

var LicenseConfigProductComponentDAO = licenseConfigProductComponentDAO{}.New()

func (input licenseConfigProductComponentDAO) New() (output licenseConfigProductComponentDAO) {
	output.FileName = "LicenseConfigProductComponentDAO.go"
	output.TableName = "license_configuration_productcomponent"
	return
}

func (input licenseConfigProductComponentDAO) InsertLicenseConfigProductComponent(tx *sql.Tx, userParam []repository.LicenseConfigComponent) (id []int64, err errorModel.ErrorModel) {
	funcName := "InsertLicenseConfigProductComponent"
	variableParam := 10
	startIndex := 1

	query := "INSERT INTO "+ input.TableName +" " +
		"(license_config_id, product_id, component_id, " +
		"component_value, created_by, created_client, " +
		"created_at, updated_by, updated_client, " +
		"updated_at) VALUES "

	query += CreateDollarParamInMultipleRowsDAO(len(userParam), variableParam, startIndex, "id")

	var params []interface{}

	for k := 0; k < len(userParam); k++ {
		params = append(params,
			userParam[k].LicenseConfigID.Int64, userParam[k].ProductID.Int64, userParam[k].ComponentID.Int64,
			userParam[k].ComponentValue.String, userParam[k].CreatedBy.Int64, userParam[k].CreatedClient.String,
			userParam[k].CreatedAt.Time, userParam[k].UpdatedBy.Int64, userParam[k].UpdatedClient.String,
			userParam[k].UpdatedAt.Time)
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

func (input licenseConfigProductComponentDAO) resultRowsInput(rows *sql.Rows) (idTemp interface{}, err errorModel.ErrorModel) {
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

func (input licenseConfigProductComponentDAO) GetLicenseConfigProductComponentByIDLicense(db *sql.DB, userParam repository.LicenseConfigComponent) (result []repository.LicenseConfigComponent, err errorModel.ErrorModel) {
	var tempResult []interface{}
	query := fmt.Sprintf(`SELECT
		lcpc.component_value, c.component_name
	LEFT JOIN %s lc ON lc.id = lcpc.license_config_id
	LEFT JOIN %s c ON c.id = lcpc.component_id
	WHERE 
		lcpc.license_config_id = $1 AND lcpc.deleted = FALSE`)

	param := []interface{}{userParam.LicenseConfigID.Int64}

	tempResult, err = GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.LicenseConfigComponent
		dbError := rows.Scan(&temp.ComponentValue, &temp.ComponentName)
		return temp, dbError
	}, param)

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(repository.LicenseConfigComponent))
		}
	}

	return
}