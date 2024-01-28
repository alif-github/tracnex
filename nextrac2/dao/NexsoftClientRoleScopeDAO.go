package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type nexsoftClientRoleScopeDAO struct {
	AbstractDAO
}

var NexsoftClientRoleScopeDAO = nexsoftClientRoleScopeDAO{}.New()

func (input nexsoftClientRoleScopeDAO) New() (output nexsoftClientRoleScopeDAO) {
	output.FileName = "NexsoftClientRoleScope.go"
	output.TableName = "nexsoft_client_role_scope"
	return
}

func (input nexsoftClientRoleScopeDAO) InsertNexsoftClientRoleScope(tx *sql.Tx, userParam repository.ClientRoleScopeModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertNexsoftClientRoleScope"
	query :=
		"INSERT INTO nexsoft_client_role_scope " +
			" (client_id, role_id, group_id, " +
			" created_by, created_client, created_at, " +
			" updated_by, updated_client, updated_at) " +
			" VALUES " +
			" ($1, $2, $3, " +
			" $4, $5, $6, " +
			" $7, $8, $9) returning id"
	errorS := tx.QueryRow(
		query, userParam.ClientID.String, userParam.RoleID.Int64, userParam.GroupID.Int64,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time).Scan(&id)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexsoftClientRoleScopeDAO) DeleteNexsoftClientRoleScopeByClientID(db *sql.Tx, userParam repository.ClientRoleScopeModel) (err errorModel.ErrorModel) {
	funcName := "DeleteClientRoleScope"
	query := "UPDATE " + NexsoftClientRoleScopeDAO.TableName + " set " +
		" deleted = TRUE, updated_by = $1, updated_client = $2, " +
		" updated_at = $3 WHERE client_id = $4"
	param := []interface{}{userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time, userParam.ClientID.String}

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

func (input nexsoftClientRoleScopeDAO) IsNexsoftClientRoleScopeExist(db *sql.Tx, userParam repository.ClientRoleScopeModel) (result repository.ClientRoleScopeModel, err errorModel.ErrorModel) {
	funcName := "IsExistNexsoftClientRoleScope"
	query := "SELECT " +
		" id, updated_at, created_by, group_id " +
		"FROM " + input.TableName + " " +
		"WHERE client_id = $1 "

	param := []interface{}{userParam.ClientID.String}

	//if userParam.CreatedBy.Int64 > 0 {
	//	query += " AND created_by = $2 "
	//	param = append(param, userParam.CreatedBy.Int64)
	//}

	query += " FOR UPDATE "

	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy, &result.GroupID)
	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexsoftClientRoleScopeDAO) UpdateNexsoftClientRoleScope(tx *sql.Tx, userParam repository.ClientRoleScopeModel) (err errorModel.ErrorModel) {
	funcName := "UpdateNexsoftClientRoleScope"

	query := "UPDATE " + input.TableName + " " +
		"SET " +
		"	role_id = $1, " +
		"	updated_at = $2, " +
		"	updated_client = $3, " +
		"	updated_by = $4, " +
		"	group_id = $5, " +
		"	deleted = false " +
		"WHERE " +
		" 	id = $6 "

	param := []interface{}{
		userParam.RoleID.Int64,
		userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String,
		userParam.UpdatedBy.Int64,
		userParam.GroupID.Int64,
		userParam.ID.Int64,
	}

	stmt, errorS := tx.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(param...)

	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}