package RegistrationNamedUserService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

type registrationNamedUserService struct {
	service.RegistrationPrepared
	service.AbstractService
}

var RegistrationNamedUserService = registrationNamedUserService{}.New()

func (input registrationNamedUserService) New() (output registrationNamedUserService) {
	output.FileName = "RegistrationNamedUserService.go"
	output.IdResourceAllowed = []int64{constanta.ResourceNexmileID, constanta.ResourceNextradeID, constanta.ResourceNexstarID}
	return
}

func (input registrationNamedUserService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.RegisterNamedUserRequest) errorModel.ErrorModel) (inputStruct in.RegisterNamedUserRequest, err errorModel.ErrorModel) {
	var (
		funcName      = "readBodyAndValidate"
		preparedError = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.NewClientType)
		stringBody    string
		isAllowed     bool
		errorS        error
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS = json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = validation(&inputStruct)
	if err.Error != nil {
		return
	}

	if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ParentClientID {
		err = errorModel.GenerateForbiddenClientCredentialAccess(input.FileName, funcName)
		return
	}

	//--- Check to DB, client type exist on table ?
	err = input.CheckIsClientTypeExist(inputStruct.ClientTypeID, preparedError)
	if err.Error != nil {
		return
	}

	//--- Check client type allowing
	for _, idResourceItem := range input.IdResourceAllowed {
		if idResourceItem == inputStruct.ClientTypeID {
			isAllowed = true
			break
		}
	}

	//--- Is not allowed, then forbidden to access
	if !isAllowed {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) readBodyAndValidateBeforeRegisterNamedUser(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.CheckNamedUserBeforeInsertRequest) errorModel.ErrorModel) (inputStruct in.CheckNamedUserBeforeInsertRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidateBeforeRegisterNamedUser"
		stringBody string
		errorS     error
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS = json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = validation(&inputStruct)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) readBodyAndValidateForUnregisterRequest(request *http.Request, validation func(input *in.UnregisterNamedUserRequest) errorModel.ErrorModel) (inputStruct in.UnregisterNamedUserRequest, err errorModel.ErrorModel) {
	id, _ := strconv.Atoi(mux.Vars(request)["id"])
	inputStruct.ID = int64(id)

	err = validation(&inputStruct)
	return
}

func (input registrationNamedUserService) createModelRegistrationNamedUserLicense(inputStruct in.RegisterNamedUserRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (registerNamedModel repository.UserRegistrationDetailModel) {
	pickRegDate := inputStruct.RegDate

	//--- When reg date is zero, then write time now
	if inputStruct.RegDate.IsZero() {
		pickRegDate = timeNow
	}

	registerNamedModel = repository.UserRegistrationDetailModel{
		ParentClientID:   sql.NullString{String: inputStruct.ParentClientID},
		ClientID:         sql.NullString{String: inputStruct.ClientID},
		AuthUserID:       sql.NullInt64{Int64: inputStruct.AuthUserID},
		ClientTypeID:     sql.NullInt64{Int64: inputStruct.ClientTypeID},
		Firstname:        sql.NullString{String: inputStruct.Firstname},
		Lastname:         sql.NullString{String: inputStruct.Lastname},
		Username:         sql.NullString{String: inputStruct.Username},
		UserID:           sql.NullString{String: inputStruct.UserID},
		Password:         sql.NullString{String: inputStruct.Password},
		ClientAliases:    sql.NullString{String: inputStruct.ClientAliases},
		SalesmanID:       sql.NullString{String: inputStruct.SalesmanID},
		AndroidID:        sql.NullString{String: inputStruct.AndroidID},
		RegDate:          sql.NullTime{Time: pickRegDate},
		Email:            sql.NullString{String: inputStruct.Email},
		UniqueID1:        sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:        sql.NullString{String: inputStruct.UniqueID2},
		NoTelp:           sql.NullString{String: fmt.Sprintf(`%s-%s`, inputStruct.CountryCode, inputStruct.NoTelp)},
		SalesmanCategory: sql.NullString{String: inputStruct.SalesmanCategory},
		CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
	}

	return
}

func (input registrationNamedUserService) createModelRegistrationNamedClientMappingUserLicense(inputStruct in.RegisterNamedUserRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (registerNamedModel repository.UserRegistrationDetailModel) {
	registerNamedModel = repository.UserRegistrationDetailModel{
		ParentClientID:   sql.NullString{String: inputStruct.ParentClientID},
		ClientID:         sql.NullString{String: inputStruct.ClientID},
		AuthUserID:       sql.NullInt64{Int64: inputStruct.AuthUserID},
		ClientTypeID:     sql.NullInt64{Int64: inputStruct.ClientTypeID},
		UserID:           sql.NullString{String: inputStruct.UserID},
		Password:         sql.NullString{String: inputStruct.Password},
		SalesmanID:       sql.NullString{String: inputStruct.SalesmanID},
		AndroidID:        sql.NullString{String: inputStruct.AndroidID},
		RegDate:          sql.NullTime{Time: timeNow},
		Email:            sql.NullString{String: inputStruct.Email},
		UniqueID1:        sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:        sql.NullString{String: inputStruct.UniqueID2},
		NoTelp:           sql.NullString{String: inputStruct.NoTelp},
		SalesmanCategory: sql.NullString{String: inputStruct.SalesmanCategory},
		CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
	}

	return
}

func (input registrationNamedUserService) checkLicenseQuota(registerNamedModel *repository.UserRegistrationDetailModel) (userLicenseOnDb repository.UserLicenseModel, err errorModel.ErrorModel) {
	var (
		fileName             = "InsertNamedUserService.go"
		funcName             = "checkLicenseQuota"
		userLicenseOnDbArray []repository.UserLicenseModel
		idxPicked            int
	)

	if userLicenseOnDbArray, err = dao.UserLicenseDAO.NewCheckLicenseNamedUser(serverconfig.ServerAttribute.DBConnection, in.CheckLicenseNamedUserRequest{
		UniqueId1:    registerNamedModel.UniqueID1.String,
		UniqueId2:    registerNamedModel.UniqueID2.String,
		ClientId:     registerNamedModel.ParentClientID.String,
		ClientTypeID: registerNamedModel.ClientTypeID.Int64,
	}); err.Error != nil {
		return
	}

	//--- User License Empty
	if len(userLicenseOnDbArray) == 0 {
		err = errorModel.GenerateUserLicenseNotFound(fileName, funcName)
		return
	}

	//--- CheckAllLicense
	idxPicked, err = input.checkAllLicenseGetInUserLicense(userLicenseOnDbArray)
	if err.Error != nil {
		return
	}

	//--- Re-Check ID User License Picked Up
	if userLicenseOnDbArray[idxPicked].ID.Int64 < 1 {
		err = errorModel.GenerateUserLicenseNotFound(fileName, funcName)
		return
	}

	//--- Lock User License
	_, err = dao.UserLicenseDAO.LockUserAndCheckIDUserLicense(serverconfig.ServerAttribute.DBConnection, userLicenseOnDbArray[idxPicked].ID.Int64)
	if err.Error != nil {
		return
	}

	//--- Fill On Register Named Model
	registerNamedModel.ParentCustomerID.Int64 = userLicenseOnDbArray[idxPicked].ParentCustomerId.Int64
	registerNamedModel.CustomerID.Int64 = userLicenseOnDbArray[idxPicked].CustomerId.Int64
	registerNamedModel.SiteID.Int64 = userLicenseOnDbArray[idxPicked].SiteId.Int64
	registerNamedModel.InstallationID.Int64 = userLicenseOnDbArray[idxPicked].InstallationId.Int64
	registerNamedModel.UserLicenseID.Int64 = userLicenseOnDbArray[idxPicked].ID.Int64
	registerNamedModel.ProductValidFrom.Time = userLicenseOnDbArray[idxPicked].ProductValidFrom.Time
	registerNamedModel.ProductValidThru.Time = userLicenseOnDbArray[idxPicked].ProductValidThru.Time

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) addNewPKCEClientMapping(tx *sql.Tx, registerNamedModel repository.UserRegistrationDetailModel, dataAudit *[]repository.AuditSystemModel, resultCheckerClientMapping repository.PKCEClientMappingModel) (err errorModel.ErrorModel) {
	var (
		isClientDependant               bool
		idResultInsertPKCEClientMapping int64
	)

	if registerNamedModel.ClientTypeID.Int64 == constanta.ResourceNexmileID {
		isClientDependant = true
	}

	idResultInsertPKCEClientMapping, err = dao.PKCEClientMappingDAO.InsertPKCEClientMapping(tx, &resultCheckerClientMapping, isClientDependant)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.PKCEClientMappingDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idResultInsertPKCEClientMapping},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_pkceclientmapping_authuserid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FileName, constanta.ClientMappingAuthUserID)
		} else if service.CheckDBError(err, "uq_user_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientMappingClientID)
		} else if service.CheckDBError(err, "uq_user_authuserid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientMappingAuthUserID)
		} else if service.CheckDBError(err, "uq_clientrolescope_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientMappingClientID)
		}
	}

	return err
}

func (input registrationNamedUserService) checkAllLicenseGetInUserLicense(userLicenseOnDbArray []repository.UserLicenseModel) (idxPicked int, err errorModel.ErrorModel) {
	fileName := "InsertNamedUserService.go"
	funcName := "checkAllLicenseGetInUserLicense"

	//----------------- Check Loop More User License
	for idx, itemUserLicenseOnDbArray := range userLicenseOnDbArray {

		//----------------- Last And No One Quota
		if len(userLicenseOnDbArray)-(idx+1) == 0 && itemUserLicenseOnDbArray.QuotaLicense.Int64 == 0 {
			err = errorModel.GenerateUserLicenseFullFilled(fileName, funcName)
			return
		}

		//----------------- If Any Then Write The Index
		if itemUserLicenseOnDbArray.QuotaLicense.Int64 > 0 {
			idxPicked = idx
			break
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) addNewUser(tx *sql.Tx, userRegistModel repository.UserRegistrationDetailModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	var idUser int64

	//------------ Insert New User
	idUser, err = dao.UserDAO.InsertUser(tx, repository.UserModel{
		ClientID:      sql.NullString{String: userRegistModel.ClientID.String},
		AuthUserID:    sql.NullInt64{Int64: userRegistModel.AuthUserID.Int64},
		Locale:        sql.NullString{String: constanta.IndonesianLanguage},
		Status:        sql.NullString{String: constanta.StatusActive},
		AliasName:     sql.NullString{String: userRegistModel.ClientAliases.String},
		FirstName:     sql.NullString{String: userRegistModel.Firstname.String},
		LastName:      sql.NullString{String: userRegistModel.Lastname.String},
		Username:      sql.NullString{String: userRegistModel.Username.String},
		Email:         sql.NullString{String: userRegistModel.Email.String},
		Phone:         sql.NullString{String: userRegistModel.NoTelp.String},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUser},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	var roleID int64
	roleID = constanta.RoleUserNexMile

	//------------ Insert New Client Credential
	var idClientRoleScope int64
	idClientRoleScope, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: userRegistModel.ClientID.String},
		RoleID:        sql.NullInt64{Int64: roleID},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientRoleScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idClientRoleScope},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) addNewUserStatusPending(tx *sql.Tx, userRegistModel repository.UserRegistrationDetailModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	var idUser int64

	//------------ Insert New User
	idUser, err = dao.UserDAO.InsertUser(tx, repository.UserModel{
		ClientID:      sql.NullString{String: userRegistModel.ClientID.String},
		AuthUserID:    sql.NullInt64{Int64: userRegistModel.AuthUserID.Int64},
		Locale:        sql.NullString{String: constanta.IndonesianLanguage},
		Status:        sql.NullString{String: constanta.PendingOnApproval},
		AliasName:     sql.NullString{String: userRegistModel.ClientAliases.String},
		FirstName:     sql.NullString{String: userRegistModel.Firstname.String},
		LastName:      sql.NullString{String: userRegistModel.Lastname.String},
		Username:      sql.NullString{String: userRegistModel.Username.String},
		Email:         sql.NullString{String: userRegistModel.Email.String},
		Phone:         sql.NullString{String: userRegistModel.NoTelp.String},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUser},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	var roleID int64
	roleID = constanta.RoleUserNexMile

	//------------ Insert New Client Credential
	var idClientRoleScope int64
	idClientRoleScope, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: userRegistModel.ClientID.String},
		RoleID:        sql.NullInt64{Int64: roleID},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientRoleScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idClientRoleScope},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) addNewUserStatusNonAktif(tx *sql.Tx, userRegistModel repository.UserRegistrationDetailModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	var idUser int64

	//------------ Insert New User
	idUser, err = dao.UserDAO.InsertUser(tx, repository.UserModel{
		ClientID:      sql.NullString{String: userRegistModel.ClientID.String},
		AuthUserID:    sql.NullInt64{Int64: userRegistModel.AuthUserID.Int64},
		Locale:        sql.NullString{String: constanta.IndonesianLanguage},
		Status:        sql.NullString{String: constanta.NonactiveUser},
		AliasName:     sql.NullString{String: userRegistModel.ClientAliases.String},
		FirstName:     sql.NullString{String: userRegistModel.Firstname.String},
		LastName:      sql.NullString{String: userRegistModel.Lastname.String},
		Username:      sql.NullString{String: userRegistModel.Username.String},
		Email:         sql.NullString{String: userRegistModel.Email.String},
		Phone:         sql.NullString{String: constanta.IndonesianCodeNumber + "-" + userRegistModel.NoTelp.String},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUser},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	var roleID int64
	roleID = constanta.RoleUserNexMile

	//------------ Insert New Client Credential
	var idClientRoleScope int64
	idClientRoleScope, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: userRegistModel.ClientID.String},
		RoleID:        sql.NullInt64{Int64: roleID},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientRoleScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idClientRoleScope},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

//func (input registrationNamedUserService) sendGroChatMessage(inputStruct in.RegisterNamedUserRequest, otpModel repository.UserVerificationModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
//	var (
//		linkNexstar = config.ApplicationConfiguration.GetNexmile().Host + config.ApplicationConfiguration.GetNexmile().PathRedirect.ActivationUser
//		groRequest  grochat_request.SendMessageGroChatRequest
//	)
//
//	linkNexstar += fmt.Sprintf(`?client_id=%s&unique_id_1=%s&unique_id_2=%s&salesman_id=%s&user_id=%s&otp=%s&phone%s`,
//		inputStruct.ClientID, inputStruct.UniqueID1, inputStruct.UniqueID2, inputStruct.SalesmanID, inputStruct.UserID, otpModel.PhoneCode.String,
//		inputStruct.NoTelp)
//
//	groRequest = grochat_request.SendMessageGroChatRequest{
//		Data: grochat_request.DetailData{
//			PhoneNumber: inputStruct.CountryCode + "-" + inputStruct.NoTelp,
//			MessageContent: grochat_request.MessageContent{
//				Message: common2.GroChatService.GetGroChatMessage(inputStruct, otpModel.PhoneCode.String, linkNexstar),
//			},
//		},
//	}
//
//	//--- Get Default
//	groRequest.GetDefault(&groRequest)
//
//	return common2.GroChatService.SendGroChatMessage(groRequest, contextModel)
//}
