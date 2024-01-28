package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type userInvitationDAO struct {
	AbstractDAO
}

var UserInvitationDAO = userInvitationDAO{}.New()

func (input userInvitationDAO) New() (output userInvitationDAO) {
	output.FileName = "UserInvitationDAO.go"
	output.TableName = "user_invitation"
	return
}

func (input userInvitationDAO) GetByEmailForUpdate(db *sql.Tx, email string) (result repository.UserInvitation, errModel errorModel.ErrorModel) {
	funcName := "GetByEmailForUpdate"

	query := `SELECT 
				id 
			FROM ` + input.TableName + ` 
			WHERE 
				email = $1 AND 
				deleted = FALSE 
			FOR UPDATE`

	row := db.QueryRow(query, email)
	err := row.Scan(&result.Id)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input userInvitationDAO) InsertTx(tx *sql.Tx, model repository.UserInvitation) (id int64, errModel errorModel.ErrorModel) {
	funcName := "InsertTx"

	query := `INSERT INTO ` + input.TableName + `(
		invitation_code, email, role_id, 
		data_group_id, expires_on, created_client, 
		updated_client, created_by, updated_by, 
		created_at, updated_at, client_id 
	) VALUES (
		$1, $2, $3,
		$4, $5, $6,
		$7, $8, $9,
		$10, $11, $12 
	) RETURNING id`

	params := []interface{}{
		model.InvitationCode.String, model.Email.String, model.RoleId.Int64,
		model.DataGroupId.Int64, model.ExpiresOn.Time, model.CreatedClient.String,
		model.UpdatedClient.String, model.CreatedBy.Int64, model.UpdatedBy.Int64,
		model.CreatedAt.Time, model.UpdatedAt.Time, model.ClientId.String,
	}

	row := tx.QueryRow(query, params...)
	if err := row.Scan(&id); err != nil {
		return 0, errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return id, errorModel.GenerateNonErrorModel()
}

func (input userInvitationDAO) UpdateByIdTx(tx *sql.Tx, model repository.UserInvitation) errorModel.ErrorModel {
	funcName := "UpdateByIdTx"

	query := `UPDATE ` + input.TableName + ` 
			SET 
				invitation_code = $1, 
				role_id = $2,
				data_group_id = $3, 
				expires_on = $4, 
				updated_at = $5,
				updated_by = $6,
				updated_client = $7,
				client_id = $8 
			WHERE 
				id = $9 AND 
				deleted = FALSE`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(
		model.InvitationCode.String,
		model.RoleId.Int64,
		model.DataGroupId.Int64,
		model.ExpiresOn.Time,
		model.UpdatedAt.Time,
		model.UpdatedBy.Int64,
		model.UpdatedClient.String,
		model.ClientId.String,
		model.Id.Int64,
	)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input userInvitationDAO) GetByEmailOrClientIdForUpdate(db *sql.Tx, invitation repository.UserInvitation) (result repository.UserInvitation, errModel errorModel.ErrorModel) {
	funcName := "GetByEmailOrClientIdForUpdate"

	query := `SELECT 
				id, email, role_id, 
				data_group_id 
			FROM ` + input.TableName + ` 
			WHERE 
				(
					email = $1 OR
					client_id = $2 
				) AND 
				deleted = FALSE 
			FOR UPDATE`

	row := db.QueryRow(query, invitation.Email.String, invitation.ClientId.String)
	err := row.Scan(
		&result.Id, &result.Email, &result.RoleId,
		&result.DataGroupId,
	)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input userInvitationDAO) DeleteByIdTx(tx *sql.Tx, id int64) errorModel.ErrorModel {
	funcName := "DeleteByIdTx"

	query := `DELETE FROM ` + input.TableName + ` 
			WHERE 
				id = $1`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
	}

	return errorModel.GenerateNonErrorModel()
}
