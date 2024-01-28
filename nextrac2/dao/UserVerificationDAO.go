package dao

import (
	"database/sql"
	"fmt"

	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type userVerificationDAO struct {
	AbstractDAO
}

var UserVerificationDAO = userVerificationDAO{}.New()

func (input userVerificationDAO) New() (output userVerificationDAO) {
	output.FileName = "UserVerificationDAO.go"
	output.TableName = "user_verification"
	return
}

func (input userVerificationDAO) GetUserVerificationForVerifying(db *sql.DB, userParam repository.UserVerificationModel) (result repository.UserVerificationModel, err errorModel.ErrorModel) {
	funcName := "GetUserVerificationForVerifying"

	query := fmt.Sprintf(
		`SELECT 
		id, email_expires, phone_expires, 
		phone_code, email_code, failed_otp_email, 
		failed_otp_phone
	FROM %s 
	WHERE 
		user_registration_detail_id = $1 `, input.TableName)

	param := []interface{}{userParam.UserRegistrationDetailID.Int64}

	if util.IsStringEmpty(userParam.Email.String) {
		query += " AND phone = $2 "
		param = append(param, userParam.Phone.String)
	} else {
		query += " AND email = $2 "
		param = append(param, userParam.Email.String)
	}

	errorS := db.QueryRow(query, param...).
		Scan(
			&result.ID, &result.EmailExpires, &result.PhoneExpires,
			&result.PhoneCode, &result.EmailCode, &result.FailedOTPEmail,
			&result.FailedOTPPhone,
		)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input userVerificationDAO) GetUserVerificationForUnregister(db *sql.DB, userParam repository.UserVerificationModel) (result repository.UserVerificationModel, err errorModel.ErrorModel) {
	funcName := "GetUserVerificationForVerifying"

	query := fmt.Sprintf(
		`SELECT 
		id 
	FROM %s 
	WHERE 
		user_registration_detail_id = $1 `, input.TableName)

	param := []interface{}{userParam.UserRegistrationDetailID.Int64}

	errorS := db.QueryRow(query, param...).
		Scan(
			&result.ID,
		)
	if errorS != nil && errorS.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	return
}

func (input userVerificationDAO) HardDeleteUserVerification(db *sql.Tx, userParam repository.UserVerificationModel) (err errorModel.ErrorModel) {
	funcName := "HardDeleteUserVerification"
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 `, input.TableName)

	_, errorDB := db.Exec(query, userParam.ID.Int64)
	if errorDB != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	return
}

func (input userVerificationDAO) InsertOTPUser(db *sql.Tx, userParam repository.UserVerificationModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertOTPUser"
	query := fmt.Sprintf(` INSERT INTO %s
      ( user_registration_detail_id, email, email_code, 
      email_expires, phone, phone_code, 
      phone_expires, created_by, created_at, 
      updated_by, updated_at ) 
   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`, input.TableName)

	params := []interface{}{
		userParam.UserRegistrationDetailID.Int64, userParam.Email.String, userParam.EmailCode.String,
		userParam.EmailExpires.Int64,
	}

	fmt.Println("Phone code before : ", params)
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Status = 200
	logModel.Message = "Phone code before : " + userParam.Phone.String
	util.LogInfo(logModel.ToLoggerObject())

	if userParam.Phone.String != "" {
		params = append(params, userParam.Phone.String, userParam.PhoneCode.String, userParam.PhoneExpires.Int64)
		fmt.Println("Phone code ada : ", params)
	} else {
		params = append(params, nil, nil, nil)
		fmt.Println("Phone code after : ", params)
		logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
		logModel.Status = 200
		logModel.Message = "Phone code before : " + userParam.PhoneCode.String
		util.LogInfo(logModel.ToLoggerObject())
	}

	params = append(params,
		userParam.CreatedBy.Int64, userParam.CreatedAt.Time, userParam.UpdatedBy.Int64,
		userParam.UpdatedAt.Time)

	result := db.QueryRow(query, params...)
	var tempResult interface{}
	if tempResult, err = RowCatchResult(result, func(rws *sql.Row) (interface{}, error) {
		var temp int64
		errorS := result.Scan(&temp)
		return temp, errorS
	}, input.FileName, funcName); err.Error != nil {
		return
	}

	if tempResult != nil {
		id = tempResult.(int64)
	}
	return
}

func (input userVerificationDAO) UpdateUserVerification(db *sql.Tx, userParam repository.UserVerificationModel) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateUserVerification"
		query    string
	)

	query = fmt.Sprintf(`UPDATE %s SET
		email_code = $1, email_expires = $2, failed_otp_email = $3,
		phone_code = $4, phone_expires = $5, failed_otp_phone = $6, 
		updated_at = $7   
		WHERE 
		id = $8 `, input.TableName)

	param := []interface{}{userParam.EmailCode.String, userParam.EmailExpires.Int64, userParam.FailedOTPEmail.Int64}

	if !util.IsStringEmpty(userParam.PhoneCode.String) {
		param = append(param, userParam.PhoneCode.String, userParam.PhoneExpires.Int64, userParam.FailedOTPPhone.Int64)
	} else {
		param = append(param, nil, nil, nil)
	}

	param = append(param, userParam.UpdatedAt.Time, userParam.ID.Int64)
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

func (input userVerificationDAO) UpdateFailedOTP(db *sql.DB, userParam repository.UserVerificationModel) (err errorModel.ErrorModel) {
	funcName := "UpdateFailedOTP"
	var params []interface{}

	query := fmt.Sprintf(`UPDATE %s SET `, input.TableName)

	if userParam.EmailExpires.Int64 < 1 {
		query += " failed_otp_phone = failed_otp_phone + 1, phone_expires = $1, "
		params = append(params, userParam.PhoneExpires.Int64)
	} else {
		query += " failed_otp_email = failed_otp_email + 1, email_expires = $1, "
		params = append(params, userParam.EmailExpires.Int64)
	}

	query += `  updated_by = $2, updated_at = $3, updated_client = $4 
  		WHERE id = $5 `

	params = append(params,
		userParam.UpdatedBy.Int64, userParam.UpdatedAt.Time, userParam.UpdatedClient.String,
		userParam.ID.Int64,
	)

	stmt, errorS := db.Prepare(query)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	_, errorS = stmt.Exec(params...)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
