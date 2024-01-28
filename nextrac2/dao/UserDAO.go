package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strings"
	"time"
)

type userDAO struct {
	AbstractDAO
}

var UserDAO = userDAO{}.New()

func (input userDAO) New() (output userDAO) {
	output.FileName = "UserDAO.go"
	output.TableName = "user"
	return
}

func (input userDAO) InsertUser(tx *sql.Tx, userParam repository.UserModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertUser"
	query := fmt.Sprintf(
		`INSERT INTO "%s" 
		(
			client_id, auth_user_id, signature_key, locale, ip_whitelist, 
			additional_info, created_client, created_by, created_at,  
			updated_by, updated_client, updated_at, 
			is_admin, first_name, last_name,  
-- 			email, phone, nt_username, status, alias_name
			email, phone, nt_username, status, alias_name, platform_device, currency
		)  
		VALUES 
		(
			$1, $2, $3, 
			$4, $5, $6, 
			$7, $8, $9, 
			$10, $11, $12, 
			$13, $14, $15, 
			$16, $17, $18, 
-- 			$19, $20
			$19, $20, $21, $22
		) 
			RETURNING id `, input.TableName)

	var param []interface{}

	param = append(param, userParam.ClientID.String, userParam.AuthUserID.Int64)

	HandleOptionalParam([]interface{}{userParam.SignatureKey.String}, &param)

	param = append(param, userParam.Locale.String)

	HandleOptionalParam([]interface{}{
		userParam.IPWhitelist.String,
		userParam.AdditionalInfo.String,
	}, &param)

	param = append(param, userParam.CreatedClient.String, userParam.CreatedBy.Int64, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.IsSystemAdmin.Bool)

	if util.IsStringEmpty(userParam.PlatformDevice.String) {
		userParam.PlatformDevice.String = "Website"
	}
	if util.IsStringEmpty(userParam.Currency.String) {
		userParam.Currency.String = "IDR"
	}

	HandleOptionalParam([]interface{}{
		userParam.FirstName.String,
		userParam.LastName.String,
		userParam.Email.String,
		userParam.Phone.String,
		userParam.Username.String,
	}, &param)
	param = append(param, userParam.Status.String)
	HandleOptionalParam([]interface{}{
		userParam.AliasName.String,
		userParam.PlatformDevice.String,
		userParam.Currency.String,
	}, &param)

	errorS := tx.QueryRow(query, param...).Scan(&id)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		if errorS != nil {
			err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		}
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) InsertMultipleUser(tx *sql.Tx, userParam []repository.UserModel) (result []int64, err errorModel.ErrorModel) {
	funcName := "InsertMultipleUser"
	parameterUser := 13
	jVar := 1
	var tempQuery string
	var tempResult []interface{}

	query := fmt.Sprintf(
		`INSERT INTO "%s" (client_id, auth_user_id, signature_key,  +
			locale, ip_whitelist, additional_info,  +
			created_client, created_by, created_at,  +
			updated_by, updated_client, updated_at,  +
			is_admin) VALUES `, input.TableName)

	for i := 1; i <= len(userParam); i++ {
		query += "("

		tempQuery, jVar = ListRangeToInQueryWithStartIndex(parameterUser, jVar)
		query += tempQuery
		query += " ) "

		if len(userParam)-i != 0 {
			query += ","
		} else {
			query += " returning id"
		}
	}
	var param []interface{}
	for i := 0; i < len(userParam); i++ {
		param = append(param,
			userParam[i].ClientID.String, userParam[i].AuthUserID.Int64, userParam[i].SignatureKey.String,
			userParam[i].Locale.String, userParam[i].IPWhitelist.String, userParam[i].AdditionalInfo.String,
			userParam[i].CreatedClient.String, userParam[i].CreatedBy.Int64, userParam[i].CreatedAt.Time,
			userParam[i].UpdatedBy.Int64, userParam[i].UpdatedClient.String, userParam[i].UpdatedAt.Time,
			userParam[i].IsSystemAdmin.Bool)
	}

	rows, errorS := tx.Query(query, param...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	tempResult, err = RowsCatchResult(rows, func(rws *sql.Rows) (interface{}, errorModel.ErrorModel) {
		var id int64
		var errors errorModel.ErrorModel
		dbError := rows.Scan(&id)
		if dbError != nil {
			errors = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			return id, errors
		}
		return id, errors
	})

	if err.Error != nil {
		return
	}

	if len(tempResult) > 0 {
		for _, item := range tempResult {
			result = append(result, item.(int64))
		}
	}
	return
}

func (input userDAO) CheckIsAuthUserExist(tx *sql.DB, userParam repository.UserModel) (province repository.UserModel, err errorModel.ErrorModel) {
	funcName := "CheckIsAuthUserExist"

	query := fmt.Sprintf(
		`SELECT
			id
		FROM
			"%s"
		WHERE
			auth_user_id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.AuthUserID.Int64}

	errorS := tx.QueryRow(query, param...).Scan(&province.ID)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) CheckIsUserExists(db *sql.DB, userParam repository.UserModel) (user repository.UserModel, err errorModel.ErrorModel) {
	var (
		funcName = "CheckIsUserExists"
		param    []interface{}
		errorS   error
	)

	query := fmt.Sprintf(`SELECT id, auth_user_id 
		FROM "%s" 
		WHERE deleted = FALSE AND auth_user_id = $1`,
		input.TableName)

	param = []interface{}{userParam.AuthUserID.Int64}
	errorS = db.QueryRow(query, param...).Scan(&user.ID, &user.AuthUserID)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) CheckIsAuthUserExistForUpdate(tx *sql.DB, userParam repository.UserModel) (province repository.UserModel, err errorModel.ErrorModel) {
	fileName := "UserDAO.go"
	funcName := "CheckIsAuthUserExistForUpdate"

	query := fmt.Sprintf(
		`SELECT
			id
		FROM
			"%s"
		WHERE
			auth_user_id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.AuthUserID.Int64}

	errorS := tx.QueryRow(query, param...).Scan(&province.ID)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) RoleMappingInternalUser(db *sql.DB, userParam repository.UserModel) (results repository.RoleMappingPersonProfileModel, err errorModel.ErrorModel) {
	funcName := "RoleMappingInternalUser"

	query := fmt.Sprintf(
		`SELECT
			userTable.id, userTable.auth_user_id, role.permission, role.role_id,
			userTable.signature_key, userTable.ip_whitelist, userTable.locale
		FROM
			"%s" userTable
		LEFT JOIN client_role_scope ON userTable.client_id = client_role_scope.client_id
		LEFT JOIN role ON client_role_scope.role_id = role.id
		WHERE
			userTable.client_id = $1 AND userTable.deleted = FALSE `, input.TableName)

	errorS := db.QueryRow(query, userParam.ClientID.String).Scan(&results.PersonProfileID,
		&results.AuthUserID, &results.Permissions,
		&results.RoleName, &results.SignatureKey,
		&results.IPWhitelist, &results.Locale)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) RoleMappingUser(db *sql.DB, userParam repository.UserModel) (results repository.RoleMappingPersonProfileModel, err errorModel.ErrorModel) {
	var (
		funcName = "RoleMappingUser"
		query    string
		errorS   error
	)

	query = fmt.Sprintf(
		`SELECT
		userTable.id, userTable.auth_user_id, role.permission, 
		role.role_id, userTable.signature_key, userTable.ip_whitelist, 
		userTable.locale, data_group.group_id, data_group.scope, 
		userTable.is_admin
		FROM "%s" userTable
		INNER JOIN client_role_scope ON userTable.client_id = client_role_scope.client_id
		INNER JOIN role ON client_role_scope.role_id = role.id
		LEFT JOIN data_group ON client_role_scope.group_id = data_group.id
		WHERE
		userTable.client_id = $1 AND userTable.deleted = FALSE AND client_role_scope.deleted = FALSE AND 
		role.deleted = FALSE AND userTable.status = 'A' `,
		input.TableName)

	errorS = db.QueryRow(query, userParam.ClientID.String).Scan(
		&results.PersonProfileID, &results.AuthUserID, &results.Permissions,
		&results.RoleName, &results.SignatureKey, &results.IPWhitelist,
		&results.Locale, &results.GroupName, &results.Scope,
		&results.IsAdmin)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) RoleMappingUserNexsoft(db *sql.DB, userParam repository.UserModel) (results repository.RoleMappingPersonProfileModel, err errorModel.ErrorModel) {
	funcName := "RoleMappingUserNexsoft"

	query := fmt.Sprintf(
		`SELECT
			userTable.id, userTable.auth_user_id, nr.permission, nr.role_id,
			userTable.signature_key, userTable.ip_whitelist, userTable.locale,
			data_group.group_id, data_group.scope, userTable.is_admin
		FROM
			"%s" userTable
		INNER JOIN nexsoft_client_role_scope ncr ON userTable.client_id = ncr.client_id
		INNER JOIN nexsoft_role nr ON ncr.role_id = nr.id
		LEFT JOIN data_group ON ncr.group_id = data_group.id
		WHERE
			userTable.client_id = $1
			AND userTable.deleted = FALSE
			AND ncr.deleted = FALSE
			AND nr.deleted = FALSE
			AND userTable.status = 'A' `, input.TableName)

	errorS := db.QueryRow(query, userParam.ClientID.String).Scan(&results.PersonProfileID,
		&results.AuthUserID, &results.Permissions,
		&results.RoleName, &results.SignatureKey,
		&results.IPWhitelist, &results.Locale,
		&results.GroupName, &results.Scope, &results.IsAdmin)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) IsAdmin(db *sql.DB, userParam repository.ViewDetailUserModel, isUrlAdmin bool) (results repository.ViewDetailUserModel, err errorModel.ErrorModel) {
	funcName := "IsAdmin"

	query := fmt.Sprintf(
		`SELECT id, is_admin
		FROM
		"%s"
		WHERE id = $1 AND deleted = FALSE `, input.TableName)

	param := []interface{}{userParam.ID.Int64}

	var queryClientID string
	var queryCreatedBy string

	if isUrlAdmin {
		queryClientID = " AND ( client_id = $2 "
		queryCreatedBy = " OR created_by = $3 ) "
	} else {
		queryClientID = " AND client_id = $2 "
		queryCreatedBy = " OR created_by = $3 "
	}

	if userParam.ClientID.String != "" {
		query += queryClientID
		param = append(param, userParam.ClientID.String)
	}

	if userParam.CreatedBy.Int64 > 0 {
		query += queryCreatedBy
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := db.QueryRow(query, param...).Scan(&results.ID, &results.IsAdmin)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) ViewDetailUserForRegisterNamedUser(db *sql.DB, userParam repository.UserModel) (results repository.UserModel, err errorModel.ErrorModel) {
	var (
		funcName = "ViewDetailUserForRegisterNamedUser"
	)

	query := fmt.Sprintf(
		`SELECT
			userTable.nt_username, userTable.first_name, userTable.last_name 
		FROM "%s" userTable
		WHERE
			userTable.auth_user_id = $1
			AND userTable.deleted = FALSE `,
		input.TableName)

	param := []interface{}{
		userParam.AuthUserID.Int64,
	}

	errorS := db.QueryRow(query, param...).Scan(
		&results.Username, &results.FirstName, &results.LastName,
	)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) ViewDetailUser(db *sql.DB, userParam repository.ViewDetailUserModel, isAdmin bool, isUrlAdmin bool) (results repository.ViewDetailUserModel, err errorModel.ErrorModel) {
	var (
		funcName             = "ViewDetailUser"
		tableClientRoleScope string
		tableRole            string
		queryClientID        string
		queryCreatedBy       string
	)

	if isAdmin {
		tableClientRoleScope = NexsoftClientRoleScopeDAO.TableName
		tableRole = NexsoftRoleDAO.TableName
	} else {
		tableClientRoleScope = ClientRoleScopeDAO.TableName
		tableRole = RoleDAO.TableName
	}

	query := fmt.Sprintf(
		`SELECT
			userTable.nt_username, userTable.first_name, userTable.last_name,
			userTable.email, userTable.phone, role.role_id,
			userTable.is_admin, userTable.created_by, userTable.created_at,
			userTable.updated_by, userTable.updated_at, userTable.id,
			gr.group_id, userTable.status, userTableCr.first_name, 
			userTableCr.last_name, userTableUp.first_name, userTableUp.last_name, 
			userTable.platform_device, userTable.currency
		FROM "%s" userTable
		INNER JOIN "%s" userTableCr ON userTable.created_by = userTableCr.id 
		INNER JOIN "%s" userTableUp ON userTable.updated_by = userTableUp.id
		LEFT JOIN %s ON userTable.client_id = %s.client_id 
		LEFT JOIN %s role ON %s.role_id = role.id 
		LEFT JOIN %s gr ON %s.group_id = gr.id 
		LEFT JOIN %s cmap ON userTable.client_id = cmap.client_id
		LEFT JOIN %s pcmap ON userTable.client_id = pcmap.client_id
		WHERE
			userTable.id = $1
			AND userTable.deleted = FALSE `,
		input.TableName, input.TableName, input.TableName,
		tableClientRoleScope, tableClientRoleScope, tableRole,
		tableClientRoleScope, DataGroupDAO.TableName, tableClientRoleScope,
		ClientMappingDAO.TableName, PKCEClientMappingDAO.TableName)

	param := []interface{}{userParam.ID.Int64}

	if isUrlAdmin {
		queryClientID = " AND ( userTable.client_id = $2 "
		queryCreatedBy = " OR userTable.created_by = $3 ) "
	} else {
		queryClientID = " AND userTable.client_id = $2 "
		queryCreatedBy = " OR userTable.created_by = $3 "
	}

	if userParam.ClientID.String != "" {
		query += queryClientID
		param = append(param, userParam.ClientID.String)
	}

	if userParam.CreatedBy.Int64 > 0 {
		query += queryCreatedBy
		param = append(param, userParam.CreatedBy.Int64)
	}

	query += fmt.Sprintf(` AND userTable.auth_user_id > 0 AND 
		cmap.id IS NULL AND 
		pcmap.id IS NULL AND 
		((userTable.email IS NOT NULL OR userTable.phone IS NOT NULL) OR 
		userTable.nt_username IN ('admin')) `)

	errorS := db.QueryRow(query, param...).Scan(
		&results.Username, &results.FirstName, &results.LastName,
		&results.Email, &results.Phone, &results.Role,
		&results.IsAdmin, &results.CreatedBy, &results.CreatedAt,
		&results.UpdatedBy, &results.UpdatedAt, &results.ID,
		&results.GroupID, &results.Status, &results.CreatedFirstName,
		&results.CreatedLastName, &results.UpdatedFirstName, &results.UpdatedLastName,
		&results.PlatformDevice, &results.Currency)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) GetUserForResendVerificationCode(db *sql.DB, userParam repository.UserModel) (results repository.UserModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUserForResendVerificationCode"
	)

	query := fmt.Sprintf(`
		SELECT 
			id, email, auth_user_id, 
			status, first_name, nt_username, 
			email, phone, is_admin 
		FROM "user" 
		WHERE id = $1 AND deleted = FALSE`,
	)

	params := []interface{}{
		userParam.ID.Int64,
	}

	errorS := db.QueryRow(query, params...).Scan(
		&results.ID, &results.Email, &results.AuthUserID,
		&results.Status, &results.FirstName, &results.Username,
		&results.Email, &results.Phone, &results.IsSystemAdmin,
	)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) UpdateLastTokenUser(db *sql.Tx, userParam repository.UserModel) (err errorModel.ErrorModel) {
	funcName := "UpdateLastTokenUser"

	query :=
		"UPDATE \"user\" SET last_token = CURRENT_TIMESTAMP WHERE client_id = $1"

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(userParam.ClientID.String)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) GetListUser(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
	var (
		value, operator string
		isFullName      bool
		params          []interface{}
	)

	query := fmt.Sprintf(`SELECT
		userTable.id, userTable.client_id, userTable.auth_user_id, userTable.locale, userTable.ip_whitelist,
		userTable.status, userTable.created_by, userTable.updated_at,
		userTable.nt_username as username, userTable.first_name, userTable.last_name,
		CASE WHEN (radmin.role_id IS NOT NULL OR radmin.role_id <> '') AND userTable.is_admin = TRUE 
			THEN radmin.role_id ELSE ruser.role_id END role_id,
		CASE WHEN (dgadmin.group_id IS NOT NULL OR dgadmin.group_id <> '') AND userTable.is_admin = TRUE 
			THEN dgadmin.group_id ELSE dguser.group_id END group_id,
		userTable.email as email, userTable.phone as phone, userTable.created_at, 
		crname.nt_username as created_name
		from "%s" userTable
			LEFT JOIN (select id, nt_username from "%s" where deleted = false) crname ON userTable.created_by = crname.id
			LEFT JOIN (select id, client_id, role_id, group_id from %s where deleted = false) crs ON userTable.client_id = crs.client_id
			LEFT JOIN (select id, client_id, role_id, group_id from %s where deleted = false) ncrs ON userTable.client_id = ncrs.client_id
			LEFT JOIN (select id, role_id from %s where deleted = false) ruser ON crs.role_id = ruser.id
			LEFT JOIN (select id, role_id from %s where deleted = false) radmin ON ncrs.role_id = radmin.id 
			LEFT JOIN (select id, group_id from %s where deleted = false) dguser ON crs.group_id = dguser.id 
			LEFT JOIN (select id, group_id from %s where deleted = false) dgadmin ON ncrs.group_id = dgadmin.id 
			LEFT JOIN (select id, client_id from %s where deleted = false) cmap ON userTable.client_id = cmap.client_id
			LEFT JOIN (select id, client_id from %s where deleted = false) pcmap ON userTable.client_id = pcmap.client_id `,
		input.TableName, input.TableName, ClientRoleScopeDAO.TableName,
		NexsoftClientRoleScopeDAO.TableName, RoleDAO.TableName, NexsoftRoleDAO.TableName,
		DataGroupDAO.TableName, DataGroupDAO.TableName, ClientMappingDAO.TableName,
		PKCEClientMappingDAO.TableName)

	additionalWhere := fmt.Sprintf(` and userTable.auth_user_id > 0 and 
		cmap.id is null and 
		pcmap.id is null and 
		(userTable.email is not null or userTable.phone is not null) `)

	if operator, value, isFullName, searchBy = input.validateSearchFullName(searchBy); isFullName {
		if operator == "LIKE" {
			additionalWhere += " AND LOWER(RTRIM(CONCAT(BTRIM(userTable.first_name) , ' ' , BTRIM(userTable.last_name)))) " + operator + " $1 "
		} else {
			additionalWhere += " AND RTRIM(CONCAT(BTRIM(userTable.first_name) , ' ' , BTRIM(userTable.last_name))) " + operator + " $1 "
		}
		params = []interface{}{value}
	}

	input.convertUserParamAndSearchBy(&userParam, searchBy)
	defaultFieldMustCheck := input.getDefaultMustCheck(createdBy)

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.ListUserModel
			errorS := rows.Scan(
				&temp.ID, &temp.ClientID, &temp.AuthUserID, &temp.Locale, &temp.IPWhiteList,
				&temp.Status, &temp.CreatedBy, &temp.UpdatedAt,
				&temp.Username, &temp.FirstName, &temp.LastName,
				&temp.RoleID, &temp.GroupID, &temp.Email,
				&temp.Phone, &temp.CreatedAt, &temp.CreatedName)
			return temp, errorS
		}, additionalWhere, defaultFieldMustCheck)
}

func (input userDAO) GetCountUser(db *sql.DB, searchBy []in.SearchByParam, createdBy int64) (int, errorModel.ErrorModel) {
	var (
		value, operator string
		isFullName      bool
		params          []interface{}
	)
	query := fmt.Sprintf(`(SELECT
		userTable.id, userTable.client_id, userTable.auth_user_id, userTable.locale, userTable.ip_whitelist,
		userTable.status, userTable.created_by, userTable.updated_at,
		userTable.nt_username, userTable.first_name, userTable.last_name,
		CASE WHEN (radmin.role_id IS NOT NULL OR radmin.role_id <> '') AND userTable.is_admin = TRUE 
			THEN radmin.role_id ELSE ruser.role_id END role_id,
		CASE WHEN (dgadmin.group_id IS NOT NULL OR dgadmin.group_id <> '') AND userTable.is_admin = TRUE 
			THEN dgadmin.group_id ELSE dguser.group_id END group_id,
		userTable.email, userTable.phone
		from "%s" userTable
			LEFT JOIN (select id, client_id, role_id, group_id from %s where deleted = false) crs ON userTable.client_id = crs.client_id
			LEFT JOIN (select id, client_id, role_id, group_id from %s where deleted = false) ncrs ON userTable.client_id = ncrs.client_id
			LEFT JOIN (select id, role_id from %s where deleted = false) ruser ON crs.role_id = ruser.id
			LEFT JOIN (select id, role_id from %s where deleted = false) radmin ON ncrs.role_id = radmin.id 
			LEFT JOIN (select id, group_id from %s where deleted = false) dguser ON crs.group_id = dguser.id 
			LEFT JOIN (select id, group_id from %s where deleted = false) dgadmin ON ncrs.group_id = dgadmin.id 
			LEFT JOIN (select id, client_id from %s where deleted = false) cmap ON userTable.client_id = cmap.client_id
			LEFT JOIN (select id, client_id from %s where deleted = false) pcmap ON userTable.client_id = pcmap.client_id `,
		input.TableName, ClientRoleScopeDAO.TableName, NexsoftClientRoleScopeDAO.TableName,
		RoleDAO.TableName, NexsoftRoleDAO.TableName, DataGroupDAO.TableName,
		DataGroupDAO.TableName, ClientMappingDAO.TableName, PKCEClientMappingDAO.TableName)

	additionalWhere := fmt.Sprintf(` and userTable.auth_user_id > 0 and 
		cmap.id is null and 
		pcmap.id is null and 
		(userTable.email is not null or userTable.phone is not null)`)

	if operator, value, isFullName, searchBy = input.validateSearchFullName(searchBy); isFullName {
		if operator == "LIKE" {
			additionalWhere += " AND LOWER(RTRIM(CONCAT(BTRIM(userTable.first_name) , ' ' , BTRIM(userTable.last_name)))) " + operator + " $1 "
		} else {
			additionalWhere += " AND RTRIM(CONCAT(BTRIM(userTable.first_name) , ' ' , BTRIM(userTable.last_name))) " + operator + " $1 "
		}
		params = []interface{}{value}
	}

	additionalWhere += ") ust"

	for i := 0; i < len(searchBy); i++ {
		if searchBy[i].SearchKey == "id" {
			searchBy[i].SearchKey = "userTable." + searchBy[i].SearchKey
		} else if searchBy[i].SearchKey == "phone" {
			searchBy[i].SearchKey = "userTable." + searchBy[i].SearchKey
		} else if searchBy[i].SearchKey == "email" {
			searchBy[i].SearchKey = "userTable." + searchBy[i].SearchKey
		}
	}

	defaultFieldMustCheck := DefaultFieldMustCheck{}
	defaultFieldMustCheck.ID.FieldName = "ust.id"
	defaultFieldMustCheck.Deleted.FieldName = "userTable.deleted"
	defaultFieldMustCheck.Deleted.Value = "FALSE"
	defaultFieldMustCheck.Deleted.IsCheck = true
	defaultFieldMustCheck.CreatedBy.FieldName = "userTable.created_by"
	defaultFieldMustCheck.CreatedBy.Value = createdBy
	defaultFieldMustCheck.CreatedBy.IsCheck = true

	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, params, query, searchBy, additionalWhere, defaultFieldMustCheck)
}

func (input userDAO) GetUserForUpdate(db *sql.Tx, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUserForUpdate"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(
		`SELECT 
			userTable.id, userTable.client_id, userTable.auth_user_id, 
			userTable.updated_at, userTable.created_by, userTable.nt_username, userTable.status,
			userTable.is_admin
		FROM "%s" AS userTable
		LEFT JOIN %s ON userTable.client_id = client_role_scope.client_id
		LEFT JOIN %s ON userTable.client_id = nexsoft_client_role_scope.client_id
		LEFT JOIN %s ON client_role_scope.role_id = role.id
		LEFT JOIN %s ON nexsoft_client_role_scope.role_id = nexsoft_role.id
		WHERE ((nexsoft_role.created_client != 'SYSTEM' AND nexsoft_role.deleted = false) OR
		(role.created_client != 'SYSTEM' AND role.deleted = false)) AND
		userTable.id = $1 AND
		userTable.deleted = FALSE `,
		input.TableName, ClientRoleScopeDAO.TableName, NexsoftClientRoleScopeDAO.TableName,
		RoleDAO.TableName, NexsoftRoleDAO.TableName)

	param = []interface{}{userParam.ID.Int64}
	if userParam.CreatedBy.Int64 > 0 {
		query += fmt.Sprintf(` AND userTable.created_by = $2 `)
		param = append(param, userParam.CreatedBy.Int64)
	}

	errorS := db.QueryRow(query, param...).Scan(
		&result.ID, &result.ClientID, &result.AuthUserID,
		&result.UpdatedAt, &result.CreatedBy, &result.Username,
		&result.Status, &result.IsSystemAdmin)

	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) DeleteUser(tx *sql.Tx, userParam repository.UserModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var (
		funcName = "DeleteUser"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`UPDATE "%s" SET 
		deleted = TRUE, updated_by = $1, updated_client = $2, 
		updated_at = $3 WHERE id = $4`,
		input.TableName)

	param = []interface{}{userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, timeNow, userParam.ID.Int64}
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

	return errorModel.GenerateNonErrorModel()
}

func (input userDAO) GetIdAndFirstNameUser(db *sql.DB, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	funcName := "RoleMappingUser"

	query := fmt.Sprintf(
		`SELECT 
			id, first_name, updated_at
		FROM "%s"
		WHERE
			client_id = $1 AND deleted = FALSE `, input.TableName)

	errorS := db.QueryRow(query, userParam.ClientID.String).Scan(&result.ID, &result.FirstName, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) UpdateUserAfterActivation(tx *sql.Tx, userParam repository.UserModel) (err errorModel.ErrorModel) {
	funcName := "UpdateUserAfterActivation"

	query := fmt.Sprintf(
		`UPDATE "user"
		SET
			status = $1
		WHERE
			auth_user_id = $2 AND
			deleted = false `)

	param := []interface{}{userParam.Status.String, userParam.AuthUserID.Int64}

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

func (input userDAO) GetDataUserByUsername(db *sql.DB, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	funcName := "RoleMappingUser"

	query := fmt.Sprintf(
		`SELECT 
			id, updated_at
		FROM "user"
		WHERE
			nt_username = $1 AND deleted = FALSE`)

	errorS := db.QueryRow(query, userParam.Username.String).Scan(&result.ID, &result.UpdatedAt)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) UpdateUser(db *sql.Tx, userParam repository.UserModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateUser"
		query    string
	)

	query = fmt.Sprintf(`UPDATE "%s" SET
		first_name = $1, last_name = $2, email = $3,
		phone = $4, updated_client = $5, updated_at = $6, 
		updated_by = $7 
		WHERE 
		id = $8 AND deleted = false `, input.TableName)

	param := []interface{}{
		userParam.FirstName.String, userParam.LastName.String, userParam.Email.String,
		userParam.Phone.String, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.ID.Int64,
	}

	stmt, errorS := db.Prepare(query)
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

func (input userDAO) UpdateUserInAdmin(db *sql.Tx, userParam repository.UserModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateUserInAdmin"
		query    string
	)

	query = fmt.Sprintf(`UPDATE "%s" SET
		first_name = $1, last_name = $2, email = $3,
		phone = $4, is_admin = $5, updated_client = $6,
 		updated_at = $7, updated_by = $8, platform_device = $9, currency = $10`, input.TableName)
	// 		updated_at = $7, updated_by = $8`, input.TableName)

	if util.IsStringEmpty(userParam.PlatformDevice.String) {
		userParam.PlatformDevice.String = "Website"
	}
	if util.IsStringEmpty(userParam.Currency.String) {
		userParam.Currency.String = "IDR"
	}

	param := []interface{}{
		userParam.FirstName.String, userParam.LastName.String, userParam.Email.String,
		userParam.Phone.String, userParam.IsSystemAdmin.Bool, userParam.UpdatedClient.String,
		//userParam.UpdatedAt.Time, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time, userParam.UpdatedBy.Int64, userParam.PlatformDevice.String, userParam.Currency.String,
	}

	last := 10
	if !util.IsStringEmpty(userParam.Status.String) {
		last++
		query += fmt.Sprintf(`, status = $%d `, last)
		param = append(param, userParam.Status.String)
	}

	last++
	query += fmt.Sprintf(` WHERE id = $%d AND deleted = false `, last)
	param = append(param, userParam.ID.Int64)

	stmt, errorS := db.Prepare(query)
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

func (input userDAO) GetUserForUpdateProfile(db *sql.Tx, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUserForUpdateProfile"
		query    string
	)

	query = fmt.Sprintf(`SELECT 
		id, client_id, auth_user_id, 
		updated_at, created_by, nt_username, 
		status, is_admin, email, 
		phone 
		FROM "%s"
		WHERE id = $1 AND is_admin = $2 AND deleted = FALSE`,
		input.TableName)

	param := []interface{}{userParam.ID.Int64, userParam.IsSystemAdmin.Bool}

	if userParam.ClientID.String != "" {
		query += " AND client_id = $3 "
		param = append(param, userParam.ClientID.String)
	}

	query += " FOR UPDATE "

	errorS := db.QueryRow(query, param...).Scan(
		&result.ID, &result.ClientID, &result.AuthUserID,
		&result.UpdatedAt, &result.CreatedBy, &result.Username,
		&result.Status, &result.IsSystemAdmin, &result.Email,
		&result.Phone)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) ReActiveUser(db *sql.Tx, userParam repository.UserModel) (err errorModel.ErrorModel) {
	funcName := "UpdateUser"

	query := fmt.Sprintf(
		`UPDATE "%s"
		SET
			updated_client = $1,
			updated_at = $2, 
			updated_by = $3,
			created_by = $4,
			created_at = $5,
			created_client = $6,
			deleted = false 
		WHERE
		id = $7 `, input.TableName)

	param := []interface{}{
		userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time,
		userParam.UpdatedBy.Int64,
		userParam.CreatedBy.Int64,
		userParam.CreatedAt.Time,
		userParam.CreatedClient.String,
		userParam.ID.Int64,
	}

	stmt, errorS := db.Prepare(query)
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

func (input userDAO) GetUserByClientID(db *sql.DB, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	funcName := "GetUserByClientID"

	query := fmt.Sprintf(`SELECT u.id, u.first_name, u.last_name, 
		u.updated_at, u.is_admin, u.nt_username, u.status, e.id_card,
		ep.name, d.name, e.is_have_member, u.platform_device, u.currency
		FROM "%s" AS u
		LEFT JOIN %s AS e
			ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
			OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
		LEFT JOIN %s AS ep 
			ON e.employee_position_id = ep.id
		LEFT JOIN %s AS d
			ON e.department_id = d.id
		WHERE u.client_id = $1 AND u.deleted = FALSE AND e.deleted = FALSE`, input.TableName,
		EmployeeDAO.TableName, EmployeePositionDAO.TableName, DepartmentDAO.TableName)

	errorS := db.QueryRow(query, userParam.ClientID.String).Scan(&result.ID, &result.FirstName, &result.LastName,
		&result.UpdatedAt, &result.IsSystemAdmin, &result.Username, &result.Status, &result.IdCard,
		&result.Position, &result.Department, &result.IsHaveMember, &result.PlatformDevice, &result.Currency)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) UpdateHelpingTableUser(db *sql.DB, userParam repository.UserModel, isForServiceDeleted bool) (err errorModel.ErrorModel) {
	funcName := "UpdateHelpingTableUser"
	var query string

	deletedFalseUser := "UPDATE \"user\" SET deleted = false WHERE id = $1"
	clientChangeFirstName := "UPDATE \"user\" SET first_name = $1 WHERE id = $2"

	var param []interface{}
	if !isForServiceDeleted {
		query = clientChangeFirstName
		param = append(param, userParam.FirstName.String)
	} else {
		query = deletedFalseUser
	}

	param = append(param, userParam.ID.Int64)

	stmt, errorS := db.Prepare(query)
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

func (input userDAO) GetAuthUserByClientID(db *sql.DB, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	funcName := "GetAuthUserByClientID"

	query := fmt.Sprintf(
		`SELECT 
			id, auth_user_id, status
		FROM "%s"
		WHERE
			client_id = $1 AND deleted = FALSE `, input.TableName)

	errorS := db.QueryRow(query, userParam.ClientID.String).Scan(&result.ID, &result.AuthUserID, &result.Status)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) GetDetailUserForCheckSessionSysuser(db *sql.DB, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	funcName := "GetDetailUserForCheckSessionSysuser"

	query := fmt.Sprintf(
		`SELECT 
			nt_username, phone, email 
		FROM "%s"
		WHERE
			client_id = $1 AND deleted = FALSE `, input.TableName)

	errorS := db.QueryRow(query, userParam.ClientID.String).Scan(&result.Username, &result.Phone, &result.Email)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) GetAuthUserByClientIDForVerifying(db *sql.DB, userParam repository.UserModel) (result repository.UserModel, err errorModel.ErrorModel) {
	funcName := "GetAuthUserByClientIDForVerifying"

	query := fmt.Sprintf(
		`SELECT 
			u.id, u.auth_user_id, u.status, 
			u.email
		FROM "%s" u
		INNER JOIN %s urd ON u.client_id = urd.client_id
		WHERE
			u.client_id = $1 AND u.deleted = FALSE
			AND u.status != $2 `, input.TableName, UserRegistrationDetailDAO.TableName)

	errorS := db.QueryRow(query, userParam.ClientID.String, constanta.StatusActive).Scan(
		&result.ID, &result.AuthUserID, &result.Status,
		&result.Email,
	)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) UpdateUserStatus(tx *sql.Tx, userParam repository.UserModel) (err errorModel.ErrorModel) {
	funcName := "UpdateUserStatus"

	query := fmt.Sprintf(
		`UPDATE "user"
		SET
			status = $1, updated_client = $2,
			updated_at = $3, updated_by = $4 
		WHERE
			auth_user_id = $5 AND
			deleted = false `)

	param := []interface{}{
		userParam.Status.String, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.AuthUserID.Int64,
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

func (input userDAO) UpdateRenewUserStatus(tx *sql.Tx, userParam repository.UserModel) (err errorModel.ErrorModel) {
	funcName := "UpdateUserStatus"

	query := fmt.Sprintf(
		`UPDATE "user"
		SET
			status = $1, updated_client = $2, updated_at = $3, 
			updated_by = $4, first_name = $5, last_name = $6, 
			email = $7, phone = $8 
		WHERE
			auth_user_id = $9 AND
			deleted = false `)

	param := []interface{}{
		userParam.Status.String, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.FirstName.String, userParam.LastName.String,
		userParam.Email.String, userParam.Phone.String,
		userParam.AuthUserID.Int64,
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

func (input userDAO) GetById(db *sql.DB, id int64) (result repository.UserModel, errModel errorModel.ErrorModel) {
	funcName := "GetById"

	query := `SELECT 
				u.id, e.id, eb.employee_level_id,
				eb.employee_grade_id, e.first_name, e.last_name  
			FROM "` + input.TableName + `" AS u 
			LEFT JOIN ` + EmployeeDAO.TableName + ` AS e 
				ON ((u.email IS NOT NULL OR u.email != '') AND u.email = e.email)
				OR ((u.email IS NULL OR u.email = '') AND u.phone = e.phone) 
			LEFT JOIN ` + EmployeeBenefitsDAO.TableName + ` AS eb 
				ON e.id = eb.employee_id
			WHERE 
				u.id = $1 AND 
				e.deleted = FALSE AND 
				u.deleted = FALSE`

	row := db.QueryRow(query, id)
	err := row.Scan(
		&result.ID, &result.EmployeeId, &result.EmployeeLevelId,
		&result.EmployeeGradeId, &result.FirstName, &result.LastName)
	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input userDAO) getDefaultMustCheck(createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "userTable.id"},
		Deleted:   FieldStatus{FieldName: "userTable.deleted"},
		CreatedBy: FieldStatus{FieldName: "userTable.created_by", Value: createdBy},
	}
}

func (input userDAO) convertUserParamAndSearchBy(userParam *in.GetListDataDTO, searchByParam []in.SearchByParam) {
	var strSplit []string
	for i := 0; i < len(searchByParam); i++ {
		searchByParam[i].SearchKey = "userTable." + searchByParam[i].SearchKey
	}

	switch userParam.OrderBy {
	case "full_name", "full_name ASC", "full_name DESC":
		strSplit = strings.Split(userParam.OrderBy, " ")
		if len(strSplit) == 2 {
			userParam.OrderBy = fmt.Sprintf("TRIM(userTable.first_name) %s, TRIM(userTable.last_name) %s", strSplit[1], strSplit[1])
		}
	case "id", "created_at ASC", "created_at DESC":
		userParam.OrderBy = fmt.Sprintf(`userTable.%s`, userParam.OrderBy)
	default:
	}
}

func (input userDAO) validateSearchFullName(searchBy []in.SearchByParam) (operator, value string, result bool, resultSearchBy []in.SearchByParam) {
	for i := 0; i < len(searchBy); i++ {
		if searchBy[i].SearchKey == "full_name" {
			if searchBy[i].SearchOperator == "like" {
				value = "%" + strings.ToLower(searchBy[i].SearchValue) + "%"
				operator = "LIKE"
			} else if searchBy[i].SearchOperator == "eq" {
				value = searchBy[i].SearchValue
				operator = "="
			}

			if i != len(searchBy)-1 {
				searchBy = append(searchBy[:i], searchBy[i+1:]...)
			} else {
				searchBy = searchBy[:i]
			}
			result = true
			break
		}
	}

	resultSearchBy = searchBy
	return
}

//func (input userDAO) GetByClientId(tx *sql.Tx, clientId string) (result repository.UserModel, errModel errorModel.ErrorModel) {
//	funcName := "GetByClientId"
//
//	query := `SELECT id
//			FROM ` + input.TableName + `
//			WHERE
//				client_id = $1 AND
//				deleted = FALSE`
//
//	row := tx.QueryRow(query, clientId)
//	err := row.Scan(&result.ID)
//	if err != nil && err != sql.ErrNoRows {
//		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
//		return
//	}
//
//	errModel = errorModel.GenerateNonErrorModel()
//	return
//}
