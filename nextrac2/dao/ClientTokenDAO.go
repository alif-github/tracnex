package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type clientTokenDAO struct {
	AbstractDAO
}

var ClientTokenDAO = clientTokenDAO{}.New()

func (input clientTokenDAO) New() (output clientTokenDAO) {
	output.FileName = "ClientTokenDAO.go"
	output.TableName = "client_token"
	return
}

func (input clientTokenDAO) InsertClientToken(db *sql.Tx, userParam repository.ClientTokenModel) (error errorModel.ErrorModel) {
	funcName := "InsertClientToken"
	query := "INSERT INTO client_token(client_id, auth_user_id, token, expired_at, created_by, created_client) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"

	stmt, err := db.Prepare(query)
	if err != nil {
		error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	_, err = stmt.Exec(userParam.ClientID.String, userParam.AuthUserID.Int64, userParam.Token.String, userParam.ExpiredAt.Time, userParam.CreatedBy.Int64, userParam.CreatedClient.String)
	if err != nil {
		error = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	error = errorModel.GenerateNonErrorModel()

	return
}

func (input clientTokenDAO) GetListTokenByRoleID(db *sql.Tx, roleID int64) (result []string, err errorModel.ErrorModel) {
	funcName := "GetListTokenByRoleID"
	query :=
		"SELECT " +
			"	token " +
			"FROM " +
			"	client_token " +
			"WHERE " +
			"	client_id IN (SELECT client_id FROM client_role_scope WHERE role_id = $1)"

	rows, errorS := db.Query(query, roleID)
	if errorS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			var temp sql.NullString
			errorS := rows.Scan(&temp)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			result = append(result, temp.String)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTokenDAO) GetListTokenByGroupID(db *sql.Tx, roleID int64) (result []string, err errorModel.ErrorModel) {
	funcName := "GetListTokenByGroupID"
	query :=
		"SELECT " +
			"	token " +
			"FROM " +
			"	client_token " +
			"WHERE " +
			"	client_id IN (SELECT client_id FROM client_role_scope WHERE group_id = $1)"

	rows, errorS := db.Query(query, roleID)
	if errorS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			var temp sql.NullString
			errorS := rows.Scan(&temp)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			result = append(result, temp.String)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTokenDAO) DeleteListTokenByRoleID(db *sql.Tx, roleID int64) (err errorModel.ErrorModel) {
	funcName := "DeleteListTokenByRoleID"
	query :=
		"DELETE FROM " +
			"	client_token " +
			"WHERE " +
			"	client_id IN (SELECT client_id FROM client_role_scope WHERE role_id = $1)"

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}
	_, errs = stmt.Exec(roleID)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
func (input clientTokenDAO) DeleteListTokenByGroupID(db *sql.Tx, groupID int64) (err errorModel.ErrorModel) {
	funcName := "DeleteListTokenByRoleID"
	query :=
		"DELETE FROM " +
			"	client_token " +
			"WHERE " +
			"	client_id IN (SELECT client_id FROM client_role_scope WHERE group_id = $1)"

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}
	_, errs = stmt.Exec(groupID)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTokenDAO) GetListTokenByClientID(db *sql.Tx, clientID string) (result []string, err errorModel.ErrorModel) {
	funcName := "GetListTokenByClientID"
	query :=
		"SELECT " +
			"	token " +
			"FROM " +
			"	client_token " +
			"WHERE " +
			"	client_id = $1"

	rows, errorS := db.Query(query, clientID)
	if errorS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
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
			var temp sql.NullString
			errorS := rows.Scan(&temp)
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
				return
			}
			result = append(result, temp.String)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTokenDAO) DeleteListTokenByClientID(db *sql.Tx, clientID string) (err errorModel.ErrorModel) {
	funcName := "DeleteListTokenByClientID"
	query :=
		"DELETE FROM " +
			"	client_token " +
			"WHERE " +
			"	client_id = $1"

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}
	_, errs = stmt.Exec(clientID)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}