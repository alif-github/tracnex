package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
	"time"
)

type userRegistrationDetailDAO struct {
	AbstractDAO
}

var UserRegistrationDetailDAO = userRegistrationDetailDAO{}.New()

func (input userRegistrationDetailDAO) New() (output userRegistrationDetailDAO) {
	output.FileName = "UserRegistrationDetailDAO.go"
	output.TableName = "user_registration_detail"
	return
}

func (input userRegistrationDetailDAO) GetCountUserRegistrationDetail(db *sql.DB, searchBy []in.SearchByParam, inputStruct in.ViewUserLicenseRequest, createdBy int64) (result int, err errorModel.ErrorModel) {
	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, []interface{}{}, input.TableName, searchBy, " AND user_license_id = "+strconv.Itoa(int(inputStruct.UserLicenseId))+" ", DefaultFieldMustCheck{}.GetDefaultField(false, createdBy))
}

func (input userRegistrationDetailDAO) GetListUserRegistrationDetail(db *sql.DB, parameterGetList in.GetListDataDTO, createdBy int64, viewUserLicenseRequest in.ViewUserLicenseRequest) (result []interface{}, err errorModel.ErrorModel) {
	additionalWhere := fmt.Sprintf(` AND urd.user_license_id = $1 `)

	query := fmt.Sprintf(`SELECT urd.id, urd.user_id, urd.salesman_id, urd.email, 
			urd.no_telp, urd.salesman_category, urd.reg_date, 
			urd.android_id, urd.status 
			FROM %s ul 
			JOIN %s urd ON ul.id = urd.user_license_id `, UserLicenseDAO.TableName, input.TableName)

	params := []interface{}{viewUserLicenseRequest.UserLicenseId}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, parameterGetList, nil,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.UserRegistrationDetailModel
			dbError := rows.Scan(
				&temp.UserRegDetailID,
				&temp.UserID,
				&temp.SalesmanId,
				&temp.Email,
				&temp.NoTelp,
				&temp.SalesmanCategory,
				&temp.RegDate,
				&temp.AndroidID,
				&temp.Status)
			return temp, dbError
		}, additionalWhere, DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "urd.id"},
			Deleted:   FieldStatus{FieldName: "urd.deleted"},
			CreatedBy: FieldStatus{FieldName: "urd.created_by", Value: createdBy},
		})
}

func (input userRegistrationDetailDAO) GetCountUserRegisDetailForCheckExpiredLicense(db *sql.DB, productLicenseID int64) (result int, err errorModel.ErrorModel) {
	var (
		funcName   = "GetCountUserRegisDetailForCheckExpiredLicense"
		query      string
		tempResult interface{}
	)

	query = fmt.Sprintf(`SELECT COUNT(urd.id) 
		FROM %s urd 
		LEFT JOIN %s ul ON ul.id = urd.user_license_id 
		LEFT JOIN %s pl ON pl.id = ul.product_license_id 
		WHERE 
		pl.id =  $1 AND urd.deleted = FALSE AND urd.status = 'A' `,
		input.TableName, UserLicenseDAO.TableName, ProductLicenseDAO.TableName)

	param := []interface{}{productLicenseID}

	rows := db.QueryRow(query, param...)
	tempResult, err = RowCatchResult(rows, func(rws *sql.Row) (interface{}, error) {
		var temp int
		dbError := rws.Scan(&temp)
		return temp, dbError
	}, input.TableName, funcName)

	if err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(int)
	}

	return
}

func (input userRegistrationDetailDAO) GetUserRegistrationDetailForUnregister(db *sql.DB, userRegDetailModel repository.UserRegistrationDetailModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	funcName := "GetUserRegistrationDetailForUnregister"
	query := fmt.Sprintf(`SELECT 
		id, status, user_license_id, 
		auth_user_id, client_id 
	FROM %s WHERE id = $1 AND deleted = FALSE FOR UPDATE `, input.TableName)

	params := []interface{}{userRegDetailModel.ID.Int64}

	var tempResult interface{}
	results := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(results, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserRegistrationDetailModel
		dbError := results.Scan(&temp.ID, &temp.Status, &temp.UserLicenseID,
			&temp.AuthUserID, &temp.ClientID)
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

func (input userRegistrationDetailDAO) GetUserRegistrationDetailForCheckExpiredLicense(db *sql.DB, productLicenseID int64) (result []repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	var (
		tempResult []interface{}
		query      string
	)

	query = fmt.Sprintf(`SELECT urd.id, urd.created_by, urd.updated_at 
		FROM %s urd 
		LEFT JOIN %s ul ON ul.id = urd.user_license_id 
		LEFT JOIN %s pl ON pl.id = ul.product_license_id 
		WHERE 
		pl.id = $1 AND urd.deleted = FALSE AND urd.status = 'A' `,
		input.TableName, UserLicenseDAO.TableName, ProductLicenseDAO.TableName)

	param := []interface{}{productLicenseID}
	tempResult, err = GetListDataDAO.GetDataRows(db, query, func(rows *sql.Rows) (interface{}, error) {
		var temp repository.UserRegistrationDetailModel
		dbErrors := rows.Scan(&temp.ID, &temp.CreatedBy, &temp.UpdatedAt)
		return temp, dbErrors
	}, param)
	if err.Error != nil {
		return
	}

	for _, item := range tempResult {
		result = append(result, item.(repository.UserRegistrationDetailModel))
	}

	return
}

func (input userRegistrationDetailDAO) UnregisterNamedUser(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (err errorModel.ErrorModel) {
	funcName := "UnregisterNamedUser"

	query := fmt.Sprintf(`UPDATE %s SET 
			status = $1, updated_by = $2, updated_client = $3, 
			updated_at = $4 
			WHERE id = $5 `, input.TableName)

	param := []interface{}{userParam.Status.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.UpdatedAt.Time, userParam.ID.Int64,
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	return
}

func (input userRegistrationDetailDAO) UpdateStatusUserRegistrationDetail(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateStatusUserRegistrationDetail"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET 
		status = $1, updated_by = $2, updated_at = $3, 
		updated_client = $4 
		WHERE 
		id = $5 `,
		input.TableName)

	param := []interface{}{
		userParam.Status.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input userRegistrationDetailDAO) GetUserForCheckRegistrationNamedUserOrRenewNamedUser(db *sql.DB, userParam repository.UserRegistrationDetailModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	var (
		funcName    = "GetUserForCheckRegistrationNamedUserOrRenewNamedUser"
		tempResult  interface{}
		indexParams = 3
	)

	query := fmt.Sprintf(`
				SELECT 
					id, auth_user_id, status, 
					client_id, email, user_license_id, 
					no_telp, client_id 
				FROM %s 
				WHERE unique_id_1 = $1 AND unique_id_2 = $2 AND deleted = FALSE `, input.TableName)

	params := []interface{}{
		userParam.UniqueID1.String, userParam.UniqueID2.String,
	}

	if userParam.AuthUserID.Int64 > 0 {
		query += fmt.Sprintf(`AND auth_user_id = $%d `, indexParams)
		params = append(params, userParam.AuthUserID.Int64)
		indexParams++
	} else {
		if !util.IsStringEmpty(userParam.UserID.String) {
			query += fmt.Sprintf(`AND user_id = $%d `, indexParams)
			params = append(params, userParam.UserID.String)
			indexParams++
		}
	}

	dbResult := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(dbResult, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserRegistrationDetailModel
		dbError := dbResult.Scan(
			&temp.ID.Int64,
			&temp.AuthUserID.Int64,
			&temp.Status.String,
			&temp.ClientID.String,
			&temp.Email.String,
			&temp.UserLicenseID.Int64,
			&temp.NoTelp.String,
			&temp.ClientID.String,
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

func (input userRegistrationDetailDAO) GetUserForCheckByAuthUserID(db *sql.DB, userParam repository.UserRegistrationDetailModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetUserForCheckByAuthUserID"
		tempResult interface{}
	)

	query := fmt.Sprintf(`
				SELECT 
					id, auth_user_id, status, 
					client_id, email, user_license_id, 
					no_telp, client_id 
				FROM %s 
				WHERE auth_user_id = $1 AND deleted = FALSE`, input.TableName)

	params := []interface{}{
		userParam.AuthUserID.Int64,
	}

	dbResult := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(dbResult, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserRegistrationDetailModel
		dbError := dbResult.Scan(
			&temp.ID.Int64,
			&temp.AuthUserID.Int64,
			&temp.Status.String,
			&temp.ClientID.String,
			&temp.Email.String,
			&temp.UserLicenseID.Int64,
			&temp.NoTelp.String,
			&temp.ClientID.String,
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

func (input userRegistrationDetailDAO) InsertUserRegistrationDetail(db *sql.Tx, userParam repository.UserRegistrationDetailModel, isFromCLientMapping bool) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertUserRegistrationDetail"
	var tempResult interface{}

	query := fmt.Sprintf(`INSERT INTO %s 
			(user_license_id, parent_customer_id, customer_id, 
			site_id, installation_id, client_id, 
			unique_id_1, unique_id_2, auth_user_id, 
			user_id, password, salesman_id, 
			android_id, reg_date, status, 
			email, no_telp, salesman_category, 
			product_valid_from, product_valid_thru, created_by, 
			created_client, created_at, updated_by, 
			updated_client, updated_at) VALUES 
			($1, $2, $3, 
			$4, $5, $6, 
			$7, $8, $9, 
			$10, $11, $12, 
			$13, $14, $15, 
			$16, $17, $18, 
			$19, $20, $21, 
			$22, $23, $24, 
			$25, $26) RETURNING id `, input.TableName)

	params := []interface{}{
		userParam.UserLicenseID.Int64, userParam.ParentCustomerID.Int64, userParam.CustomerID.Int64,
		userParam.SiteID.Int64, userParam.InstallationID.Int64, userParam.ClientID.String,
		userParam.UniqueID1.String,
	}

	HandleOptionalParam([]interface{}{userParam.UniqueID2.String}, &params)

	params = append(params, userParam.AuthUserID.Int64)

	HandleOptionalParam([]interface{}{userParam.UserID.String, userParam.Password.String,
		userParam.SalesmanID.String, userParam.AndroidID.String}, &params)

	if !userParam.RegDate.Time.IsZero() {
		params = append(params, userParam.RegDate.Time)
	} else {
		params = append(params, time.Now())
	}

	if isFromCLientMapping {
		params = append(params, constanta.StatusRegistered)
	} else {
		params = append(params, constanta.StatusActive)
	}

	HandleOptionalParam([]interface{}{userParam.Email.String, userParam.NoTelp.String,
		userParam.SalesmanCategory.String}, &params)

	params = append(params,
		userParam.ProductValidFrom.Time, userParam.ProductValidThru.Time, userParam.CreatedBy.Int64,
		userParam.CreatedClient.String, userParam.CreatedAt.Time, userParam.UpdatedBy.Int64,
		userParam.UpdatedClient.String, userParam.UpdatedAt.Time)

	result := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(result, func(rws *sql.Row) (interface{}, error) {
		var idTemp int64
		errorS := result.Scan(&idTemp)
		return idTemp, errorS
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		id = tempResult.(int64)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationDetailDAO) GetUserRegistrationDetailForActivation(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	funcName := "GetUserRegistrationDetailForActivation"
	var tempResult interface{}

	query := fmt.Sprintf(`SELECT 
			id, user_license_id, auth_user_id, 
			client_id, unique_id_1, unique_id_2 
			FROM %s 
			WHERE deleted = FALSE AND id = $1 AND status = $2 `, input.TableName)

	param := []interface{}{userParam.ID.Int64, constanta.StatusRegistered}

	query += "FOR UPDATE"

	row := db.QueryRow(query, param...)
	if tempResult, err = RowCatchResult(row, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserRegistrationDetailModel
		errorS := rws.Scan(&temp.ID, &temp.UserLicenseID, &temp.AuthUserID,
			&temp.ClientID, &temp.UniqueID1, &temp.UniqueID2)
		return temp, errorS
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		result = tempResult.(repository.UserRegistrationDetailModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationDetailDAO) UpdateActivationUserRegistrationDetail(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (err errorModel.ErrorModel) {
	funcName := "UpdateActivationUserRegistrationDetail"

	query := fmt.Sprintf(`UPDATE %s SET 
			status = $1, updated_by = $2, updated_at = $3, 
			updated_client = $4, user_license_id = $5 
			WHERE id = $6`, input.TableName)

	param := []interface{}{userParam.Status.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.UserLicenseID.Int64, userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input userRegistrationDetailDAO) UpdateRenewUserRegistrationDetail(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (err errorModel.ErrorModel) {
	funcName := "UpdateRenewUserRegistrationDetail"

	query := fmt.Sprintf(`UPDATE %s SET 
			status = $1, updated_by = $2, updated_at = $3, 
			updated_client = $4, user_license_id = $5, user_id = $6, 
			salesman_category = $7, salesman_id = $8, email = $9, 
			no_telp = $10
			WHERE id = $11`, input.TableName)

	param := []interface{}{
		userParam.Status.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.UserLicenseID.Int64, userParam.UserID.String,
		userParam.SalesmanCategory.String, userParam.SalesmanID.String, userParam.Email.String,
		userParam.NoTelp.String,
		userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)

	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	return
}

func (input userRegistrationDetailDAO) GetDataForValidateNamedUser(db *sql.DB, userRegistrationDetail repository.UserRegistrationDetailModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	var (
		funcName   = "GetDataForValidateNamedUser"
		tempResult interface{}
	)

	query := fmt.Sprintf(`SELECT id, client_id, product_valid_from, 
			product_valid_thru, unique_id_1, unique_id_2, 
			status FROM %s 
			WHERE status = 'A' AND unique_id_1 = $1 `, input.TableName)

	params := []interface{}{
		userRegistrationDetail.UniqueID1.String,
	}

	//-------------- Check Optional Field
	idx := 1
	if userRegistrationDetail.UniqueID2.Valid {
		idx++
		query += " AND unique_id_2 = $" + strconv.Itoa(idx)
		params = append(params, userRegistrationDetail.UniqueID2.String)
	}

	if userRegistrationDetail.AuthUserID.Valid {
		idx++
		query += " AND auth_user_id = $" + strconv.Itoa(idx)
		params = append(params, userRegistrationDetail.AuthUserID.Int64)
	}

	//if userRegistrationDetail.UserID.Valid {
	//	idx++
	//	query += " AND user_id = $" + strconv.Itoa(idx)
	//	params = append(params, userRegistrationDetail.UserID.String)
	//
	//}

	dbResult := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(dbResult, func(rws *sql.Row) (interface{}, error) {
		var temp repository.UserRegistrationDetailModel
		dbError := dbResult.Scan(
			&temp.ID.Int64,
			&temp.ClientID.String,
			&temp.ProductValidFrom.Time,
			&temp.ProductValidThru.Time,
			&temp.UniqueID1.String,
			&temp.UniqueID2.String,
			&temp.Status.String,
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

func (input userRegistrationDetailDAO) GetUserActiveRegistrationForVerifying(db *sql.DB, userParam repository.UserRegistrationDetailModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	var (
		funcName = "GetUserActiveRegistrationForVerifying"
		param    []interface{}
	)

	query := fmt.Sprintf(
		`SELECT
		urd.id, urd.client_id, urd.unique_id_1, 
		urd.unique_id_2, lc.max_offline_days, urd.product_valid_thru,
		urd.product_valid_from, pl.license_status
	FROM %s urd
	LEFT JOIN %s ul ON ul.id = urd.user_license_id
	LEFT JOIN %s pl ON pl.id = ul.product_license_id
	LEFT JOIN %s lc ON lc.id = pl.license_config_id
	WHERE 
		urd.user_id = $1 AND urd.unique_id_1 = $2 AND urd.unique_id_2 = $3 
		AND urd.deleted = FALSE AND urd.status = $4 AND pl.license_status = $5 `,
		input.TableName, UserLicenseDAO.TableName, ProductLicenseDAO.TableName,
		LicenseConfigDAO.TableName)

	param = append(param, userParam.UserID.String, userParam.UniqueID1.String, userParam.UniqueID2.String,
		constanta.StatusActive, constanta.ProductLicenseStatusActive)

	errorS := db.QueryRow(query, param...).
		Scan(
			&result.ID, &result.ClientID, &result.UniqueID1,
			&result.UniqueID2, &result.MaxOfflineDays, &result.ProductValidThru,
			&result.ProductValidFrom, &result.LicenseStatus,
		)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input userRegistrationDetailDAO) UpdateAndroidID(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (err errorModel.ErrorModel) {
	funcName := "UpdateAndroidID"

	query := fmt.Sprintf(
		`UPDATE %s 
	SET
		android_id = $1, updated_by = $2, updated_at = $3, 
		updated_client = $4 
	WHERE id = $5 `, input.TableName)

	param := []interface{}{
		userParam.AndroidID.String, userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time,
		userParam.UpdatedClient.String, userParam.ID.Int64,
	}

	stmt, dbError := db.Prepare(query)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}

	_, dbError = stmt.Exec(param...)
	if dbError != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
	}
	return
}

func (input userRegistrationDetailDAO) GetMappingUserValidation(db *sql.DB, validationModel repository.UserRegistrationDetailModel) (err errorModel.ErrorModel, result repository.UserRegistrationDetailMapping) {
	var (
		funcName = "GetMappingUserValidation"
		query    string
		param    []interface{}
	)

	query = fmt.Sprintf(`SELECT 
		urd.id as user_registration_detail_id, urd.auth_user_id, urd.user_id, 
		urd."password", urd.unique_id_1, urd.unique_id_2, 
		urd.product_valid_thru, uv.id_verification, uv.updated_at_verification, 
		pcm.pkce_client_mapping_id, pcm.client_id, u.user_nextrac_id, 
		u.user_status, urd.email, urd.no_telp, 
		u.first_name, pcm.client_type_id, urd.salesman_id, 
		u.alias_name, up.alias_name, uv.phone, 
		cm.id
		FROM %s urd 
			LEFT JOIN (
				SELECT id as id_verification, updated_at as updated_at_verification, user_registration_detail_id, 
				phone
				FROM %s 
				WHERE email = $1 or phone = $2) uv 
				ON urd.id = uv.user_registration_detail_id 
			LEFT JOIN (
				SELECT id as pkce_client_mapping_id, client_id, client_type_id, 
				parent_client_id, company_id, branch_id
				FROM %s 
				WHERE client_type_id = $3 and deleted = FALSE) pcm 
				ON urd.client_id = pcm.client_id 
			LEFT JOIN (
				SELECT id as user_nextrac_id, status as user_status, client_id, 
				first_name, alias_name 
				FROM "%s" 
				WHERE deleted = FALSE) u 
				ON urd.client_id = u.client_id 
			LEFT JOIN (
				SELECT alias_name, client_id 
				FROM "%s" 
				WHERE deleted = FALSE) up 
				ON pcm.client_id = up.client_id 
			LEFT JOIN %s cm 
				ON cm.client_id = pcm.parent_client_id
				AND cm.company_id = pcm.company_id 
				AND cm.branch_id = pcm.branch_id
		WHERE
		urd.deleted = FALSE AND 
		urd.user_id = $4 AND 
		urd.status = $5 AND 
		urd.unique_id_1 = $6 AND 
		urd.auth_user_id = $7 `,
		input.TableName, UserVerificationDAO.TableName, PKCEClientMappingDAO.TableName,
		UserDAO.TableName, UserDAO.TableName, ClientMappingDAO.TableName)

	param = append(param,
		validationModel.Email.String, validationModel.NoTelp.String, validationModel.ClientTypeID.Int64,
		validationModel.UserID.String, validationModel.Status.String, validationModel.UniqueID1.String,
		validationModel.AuthUserID.Int64)

	if validationModel.UniqueID2.String != "" {
		query += fmt.Sprintf(` AND urd.unique_id_2 = $8 `)
		param = append(param, validationModel.UniqueID2.String)
	}

	dbResult := db.QueryRow(query, param...)

	fmt.Println("params GetMappingUserValidation : ", param)
	fmt.Printf(query, input.TableName)

	errorS := dbResult.Scan(
		&result.UserRegistrationDetail.ID, &result.UserRegistrationDetail.AuthUserID, &result.UserRegistrationDetail.UserID,
		&result.UserRegistrationDetail.Password, &result.UserRegistrationDetail.UniqueID1, &result.UserRegistrationDetail.UniqueID2,
		&result.UserRegistrationDetail.ProductValidThru, &result.UserVerification.ID, &result.UserVerification.UpdatedAt,
		&result.PKCEClientMapping.ID, &result.PKCEClientMapping.ClientID, &result.User.ID,
		&result.User.Status, &result.UserRegistrationDetail.Email, &result.UserRegistrationDetail.NoTelp,
		&result.User.FirstName, &result.PKCEClientMapping.ClientTypeID, &result.UserRegistrationDetail.SalesmanID,
		&result.PKCEClientMapping.BranchName, &result.PKCEClientMapping.CompanyName, &result.UserVerification.Phone,
		&result.ClientMapping.ID)

	if errorS != nil && errorS != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationDetailDAO) GetUserRegistrationDetailWithUserIDAndPassword(db *sql.DB, userParam repository.UserRegistrationDetailModel) (result repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	funcName := "GetUserRegistrationDetailWithUserIDAndPassword"

	query := fmt.Sprintf(
		`SELECT 
		urd.id, urd.client_id, urd.unique_id_1, 
		urd.unique_id_2, lc.max_offline_days, urd.product_valid_thru, 
		cm.client_id, urd.product_valid_from, pl.license_status, 
		urd.auth_user_id
	FROM %s urd 
	INNER JOIN %s pkce ON urd.client_id = pkce.client_id
	INNER JOIN %s cm ON pkce.parent_client_id = cm.client_id
	LEFT JOIN %s ul ON ul.id = urd.user_license_id
	LEFT JOIN %s pl ON pl.id = ul.product_license_id
	LEFT JOIN %s lc ON lc.id = pl.license_config_id
	WHERE 
		urd.user_id = $1 AND urd.password = $2 AND urd.android_id = $3 
		AND pkce.client_type_id = $4 AND urd.deleted = FALSE AND urd.status = $5 
		AND urd.auth_user_id = $6 `,
		input.TableName, PKCEClientMappingDAO.TableName, ClientMappingDAO.TableName,
		UserLicenseDAO.TableName, ProductLicenseDAO.TableName, LicenseConfigDAO.TableName)

	param := []interface{}{
		userParam.UserID.String, userParam.Password.String, userParam.AndroidID.String,
		userParam.ClientTypeID.Int64, constanta.StatusActive, userParam.AuthUserID.Int64,
	}

	errorS := db.QueryRow(query, param...).
		Scan(
			&result.ID, &result.ClientID, &result.UniqueID1,
			&result.UniqueID2, &result.MaxOfflineDays, &result.ProductValidThru,
			&result.ParentClientID, &result.ProductValidFrom, &result.LicenseStatus,
			&result.AuthUserID,
		)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input userRegistrationDetailDAO) MoveUserRegistrationDetail(db *sql.Tx, userParam repository.UserRegistrationDetailModel) (countData int64, err errorModel.ErrorModel) {
	var (
		funcName   = "MoveUserRegistrationDetail"
		tempResult interface{}
	)

	query := fmt.Sprintf(`
		WITH rows AS (
			UPDATE %s urd 
				SET user_license_id = $1
			FROM (
				SELECT urd.id from %s urd
					INNER JOIN %s ul on ul.id = urd.user_license_id 
					INNER JOIN %s pl on pl.id = ul.product_license_id 
					INNER JOIN %s lc on lc.id = pl.license_config_id
				WHERE lc.id = (SELECT old_license_configuration_id FROM license_configuration lcu WHERE lcu.id = $2)
			) AS sub_query
				WHERE urd.id = sub_query.id
				RETURNING urd.id , urd.status
			)
		SELECT COUNT(rows.id) FROM rows WHERE rows.status = 'A'`,
		input.TableName, input.TableName, UserLicenseDAO.TableName,
		ProductLicenseDAO.TableName, LicenseConfigDAO.TableName)

	params := []interface{}{
		userParam.UserLicenseID.Int64,
		userParam.LicenseConfigID.Int64,
	}

	result := db.QueryRow(query, params...)
	if tempResult, err = RowCatchResult(result, func(rws *sql.Row) (interface{}, error) {
		var idTemp int64
		errorS := result.Scan(&idTemp)
		return idTemp, errorS
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		countData = tempResult.(int64)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
