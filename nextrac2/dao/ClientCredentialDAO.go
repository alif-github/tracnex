package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

type clientCredentialDAO struct {
	AbstractDAO
}

var ClientCredentialDAO = clientCredentialDAO{}.New()

func (input clientCredentialDAO) New() (output clientCredentialDAO) {
	output.FileName = "ClientCredentialDAO.go"
	output.TableName = "client_credential"
	return
}

func (input clientCredentialDAO) IsClientCredentialExist(db *sql.DB, userParam repository.ClientCredentialModel) (result bool, err errorModel.ErrorModel) {
	funcName := "IsClientCredentialExist"

	query := "SELECT " +
		" 	CASE WHEN COUNT(id) > 0 " +
		"		THEN TRUE " +
		"		ELSE FALSE " +
		"	END is_exist " +
		" FROM " + input.TableName + " " +
		" WHERE " +
		"	deleted = FALSE AND " +
		"	client_id = $1 AND " +
		"	client_secret = $2 "

	param := []interface{}{userParam.ClientID.String, userParam.ClientSecret.String}
	results := db.QueryRow(query, param...)
	dbError := results.Scan(&result)

	if dbError != nil && dbError.Error() != constanta.NoRowsInDB {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientCredentialDAO) InsertClientCredential(tx *sql.Tx, userParam *repository.ClientCredentialModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertClientCredential"
	query := "INSERT INTO " + input.TableName + " " +
		"(client_id, client_secret, signature_key, " +
		"created_by, created_client, created_at, " +
		"updated_by, updated_client, updated_at) " +
		"VALUES ($1, $2, $3, " +
		"$4, $5, $6, " +
		"$7, $8, $9) " +
		"returning id"
	errorS := tx.QueryRow(query, userParam.ClientID.String, userParam.ClientSecret.String, userParam.SignatureKey.String,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time).Scan(&id)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientCredentialDAO) GetClientCredential(db *sql.DB, userParam []repository.ClientCredentialModel) (result []repository.ClientCredentialModel, err errorModel.ErrorModel) {
	funcName := "GetClientCredential"

	query := "SELECT client_id, client_secret, signature_key " +
		"from " + input.TableName + " " +
		"where " +
		"client_id IN ("
	for i := 1; i <= len(userParam); i++ {
		query += "$" + strconv.Itoa(i)

		if len(userParam)-i > 0 {
			query += ","
		} else {
			query += ")"
		}
	}

	var param []interface{}
	for j := 0; j < len(userParam); j++ {
		param = append(param, userParam[j].ClientID.String)
	}

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
			var temp repository.ClientCredentialModel
			errorS = rows.Scan(&temp.ClientID, &temp.ClientSecret, &temp.SignatureKey)
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

func (input clientCredentialDAO) GetClientCredentialByClientID(db *sql.DB, userParam repository.ClientCredentialModel) (result repository.ClientCredentialModel, err errorModel.ErrorModel) {
	funcName := "GetClientCredentialByClientID"

	query := "SELECT " +
		"id, client_id, client_secret, " +
		"signature_key " +
		"FROM " + ClientCredentialDAO.TableName + " " +
		"WHERE " +
		"client_id = $1 AND deleted = FALSE "

	param := []interface{}{userParam.ClientID.String}
	results := db.QueryRow(query, param...)
	dbError := results.Scan(&result.ID, &result.ClientID, &result.ClientSecret,
		&result.SignatureKey)

	if dbError != nil && dbError.Error() != constanta.NoRowsInDB {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientCredentialDAO) GetClientCredentialForActivationLicense(db *sql.DB, userParam repository.ClientCredentialModel) (result repository.ClientCredentialModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetClientCredentialForActivationLicense"
		query    string
	)

	query = fmt.Sprintf(`
		SELECT DISTINCT cc.client_secret, cc.client_id, cc.signature_key
		FROM %s cc
		INNER JOIN %s cm ON cc.client_id = cc.client_id
		WHERE cc.client_id = $1 AND cc.deleted = FALSE AND 
		cm.client_type_id = $2 AND cc.signature_key = $3 `,
		input.TableName, ClientMappingDAO.TableName)

	param := []interface{}{
		userParam.ClientID.String,
		userParam.ClientTypeID.Int64,
		userParam.SignatureKey.String,
	}
	results := db.QueryRow(query, param...)

	tempResult, err := RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var temp repository.ClientCredentialModel
		dbErrorS := rws.Scan(&temp.ClientSecret, &temp.ClientID, &temp.SignatureKey)
		return temp, dbErrorS
	}, input.FileName, funcName)

	if err.Error != nil {
		return
	}

	result = tempResult.(repository.ClientCredentialModel)
	return
}
