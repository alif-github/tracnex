package CRUDUserService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"strings"
	"time"
)

func (input userService) doMappingUser(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(userStruct).inputStruct
	if inputStruct.AuthUserID > 0 {
		return input.doAddUserFromAuth(tx, inputStruct, contextModel, timeNow)
	}

	return input.doInsertUserSysAdmin(tx, inputStructInterface, contextModel, timeNow)
}

func (input userService) doAddUserFromAuth(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName       = "doAddUserFromAuth"
		inputStruct    = inputStructInterface.(in.UserRequest)
		fileName       = input.FileName
		userDetailAuth authentication_response.UserAuthenticationResponse
	)

	//--- Check user detail to auth server
	userDetailAuth, err = input.checkUserDetailToAuthServer(inputStructInterface, contextModel)
	if err.Error != nil {
		return
	}

	//--- Check auth user id
	if userDetailAuth.Nexsoft.Payload.Data.Content.UserID != inputStruct.AuthUserID {
		err = errorModel.GenerateDifferentAuthUserId(fileName, funcName)
		return
	}

	//--- Update information user
	input.updateInformationUser(&inputStruct, userDetailAuth)

	//--- Add resource trac to user auth
	err = HitAuthenticationServerForAddResourceNextracUserAuth(userDetailAuth.Nexsoft.Payload.Data.Content, contextModel)
	if err.Error != nil {
		return
	}

	//--- Insert user from auth to trac
	dataAudit, err = input.saveUserAuthToDB(tx, inputStruct, contextModel, timeNow, constanta.Remark)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateUserFromAuth(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateAddResourceToUser()
}

func (input userService) saveUserAuthToDB(tx *sql.Tx, inputStruct in.UserRequest, contextModel *applicationModel.ContextModel, timeNow time.Time, code string) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName             = "saveUserAuthToDB"
		roleScope            repository.ClientRoleScopeModel
		dataGroupOnDB        repository.DataGroupModel
		registerUserResponse authentication_response.RegisterUserAuthenticationResponse
	)

	err = input.validateInsert(&inputStruct)
	if err.Error != nil {
		return
	}

	roleScope, err = input.checkRole(inputStruct)
	if err.Error != nil {
		return
	}

	//--- Check Data Group ID
	dataGroupOnDB, err = dao.DataGroupDAO.GetRoleByName(tx, repository.DataGroupModel{GroupID: sql.NullString{String: inputStruct.DataGroupID}})
	if err.Error != nil {
		return
	}

	if dataGroupOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.DataGroup)
		return
	}

	clientRoleScope := repository.ClientRoleScopeModel{
		RoleID:  roleScope.RoleID,
		GroupID: dataGroupOnDB.ID,
	}

	registerUserResponse = authentication_response.RegisterUserAuthenticationResponse{
		Nexsoft: authentication_response.RegisterUserBodyResponse{
			Header: model.HeaderResponse{},
			Payload: authentication_response.RegisterUserPayload{
				PayloadResponse: model.PayloadResponse{},
				Data: authentication_response.RegisterUserData{
					Content: authentication_response.RegisterUserContent{
						UserID:       inputStruct.AuthUserID,
						ClientID:     inputStruct.ClientID,
						SignatureKey: inputStruct.SignatureKey,
						NotifyStatus: authentication_response.StatusNotifyEmailOrPhone{},
					},
				},
			},
		},
	}

	dataAudit, err = input.saveUserSysAdminToDB(tx, inputStruct, registerUserResponse, clientRoleScope, contextModel, timeNow, code)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) updateInformationUser(inputStruct *in.UserRequest, userFromAuth authentication_response.UserAuthenticationResponse) {
	var auth = userFromAuth.Nexsoft.Payload.Data.Content

	//--- If Email Auth Not Empty
	if !util.IsStringEmpty(auth.Email) {
		inputStruct.Email = auth.Email
	}

	//--- If Phone Auth Not Empty
	if !util.IsStringEmpty(auth.Phone) {
		s := strings.Split(auth.Phone, "-")
		inputStruct.CountryCode = s[0]
		inputStruct.Phone = s[1]
	}

	inputStruct.Username = auth.Username
	inputStruct.ClientID = auth.ClientID
	inputStruct.Locale = auth.Locale
	inputStruct.ResourceID = auth.ResourceID
	inputStruct.IPWhitelist = auth.IPWhitelist
	inputStruct.SignatureKey = auth.SignatureKey

	return
}

func (input userService) saveUserFromAuthToDB(tx *sql.Tx, inputStruct in.UserRequest, content out.UserAuthDetail,
	ClientRoleScope repository.ClientRoleScopeModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {

	userModel := repository.UserModel{
		FirstName:      sql.NullString{String: inputStruct.FirstName},
		LastName:       sql.NullString{String: inputStruct.LastName},
		Username:       sql.NullString{String: inputStruct.Username},
		Email:          sql.NullString{String: inputStruct.Email},
		Phone:          sql.NullString{String: constanta.IndonesianCodeNumber + "-" + inputStruct.Phone},
		ClientID:       sql.NullString{String: content.ClientId},
		AuthUserID:     sql.NullInt64{Int64: content.UserId},
		Locale:         sql.NullString{String: inputStruct.Locale},
		SignatureKey:   sql.NullString{String: content.SignatureKey},
		AdditionalInfo: sql.NullString{String: inputStruct.AdditionalInformationString()},
		IPWhitelist:    sql.NullString{String: contextModel.LoggerModel.IP},
		IsSystemAdmin:  sql.NullBool{Bool: inputStruct.IsAdmin},
		Status:         sql.NullString{String: constanta.PendingOnApproval},
		CreatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt:      sql.NullTime{Time: timeNow},
		CreatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:      sql.NullTime{Time: timeNow},
		UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	var id int64
	id, err = dao.UserDAO.InsertUser(tx, userModel)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: id},
	})

	ClientRoleScope.ClientID.String = content.ClientId
	ClientRoleScope.CreatedBy = userModel.CreatedBy
	ClientRoleScope.CreatedAt = userModel.CreatedAt
	ClientRoleScope.CreatedClient = userModel.CreatedClient
	ClientRoleScope.UpdatedBy = userModel.UpdatedBy
	ClientRoleScope.UpdatedAt = userModel.UpdatedAt
	ClientRoleScope.UpdatedClient = userModel.UpdatedClient

	if inputStruct.IsAdmin {
		id, err = dao.NexsoftClientRoleScopeDAO.InsertNexsoftClientRoleScope(tx, ClientRoleScope)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.NexsoftClientRoleScopeDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: id},
		})
	} else {
		id, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, ClientRoleScope)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.ClientRoleScopeDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: id},
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func HitAuthenticateServerForGetDetailUserAuth(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel) (structResponse authentication_response.UserAuthenticationResponse, err errorModel.ErrorModel) {
	var (
		fileName   = "InsertUserCheckResource.go"
		funcName   = "hitAuthenticateServerForGetDetailUserAuth"
		resourceID = constanta.AuthDestination
		auth       = config.ApplicationConfiguration.GetAuthenticationServer()
		url        = auth.Host + auth.PathRedirect.InternalUser.CheckUser
		locale     = constanta.DefaultApplicationsLanguage
		method     = http.MethodPost
	)

	authenticationRequest := authentication_request.CheckUser{
		Username: inputStruct.Username,
		Phone:    inputStruct.Phone,
		Email:    inputStruct.Email,
		ClientID: inputStruct.ClientID,
	}

	authenticationRequestJSON := util.StructToJSON(authenticationRequest)

	internalToken := resource_common_service.GenerateInternalToken(resourceID, config.ApplicationConfiguration.GetClientCredentialsAuthUserID(), config.ApplicationConfiguration.GetClientCredentialsClientID(), config.ApplicationConfiguration.GetServerResourceID(), locale)
	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errs := common.HitAPI(url, header, authenticationRequestJSON, method, *contextModel)
	if errs != nil {
		err = errorModel.GenerateErrorModel(statusCode, errs.Error(), fileName, funcName, errs)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &structResponse)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(structResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, structResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func (input userService) hitAuthenticateServerForCekSignatureKey(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel) (structResponse out.APIResponseCheckResource, err errorModel.ErrorModel) {
	var (
		funcName   = "hitAuthenticateServerForCekSignatureKey"
		resourceID = constanta.AuthDestination
		auth       = config.ApplicationConfiguration.GetAuthenticationServer()
		url        = auth.Host + auth.PathRedirect.CheckUser
		locale     = constanta.DefaultApplicationsLanguage
		method     = http.MethodPost
	)

	authenticationRequest := in.UserRequest{
		Email: inputStruct.Email,
		Phone: inputStruct.Phone,
	}

	authenticationRequestJSON := util.StructToJSON(authenticationRequest)

	internalToken := resource_common_service.GenerateInternalToken(resourceID, 0, "", "", locale)
	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errs := common.HitAPI(url, header, authenticationRequestJSON, method, *contextModel)

	_ = json.Unmarshal([]byte(bodyResult), &structResponse)

	if errs != nil {
		err = errorModel.GenerateErrorModel(statusCode, errs.Error(), input.FileName, funcName, errs)
		return
	}

	if statusCode != 200 {
		err = errorModel.GenerateErrorModel(statusCode, structResponse.Nexsoft.Payload.Status.Message, input.FileName, funcName, errs)
		return
	}

	return
}

func HitAuthenticationServerForAddResourceNextracUserAuth(user authentication_response.UserContent, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName       = "InsertUserCheckResource.go"
		funcName       = "hitAuthenticationServerForAddResourceNextracUserAuth"
		locale         = constanta.DefaultApplicationsLanguage
		auth           = config.ApplicationConfiguration.GetAuthenticationServer()
		url            = auth.Host + auth.PathRedirect.AddResourceClient
		method         = http.MethodPost
		header         = make(map[string][]string)
		internalToken  string
		structResponse out.APIResponseAddResource
	)

	authenticationRequestJSON := util.StructToJSON(in.UserRequest{
		ClientID:   user.ClientID,
		ResourceID: config.ApplicationConfiguration.GetServerResourceID(),
	})

	internalToken = resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.AuthDestination, locale)
	header[constanta.TokenHeaderNameConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errs := common.HitAPI(url, header, authenticationRequestJSON, method, *contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &structResponse)
	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(structResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, structResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func HitAuthenticationServerForAddResourceToUserAuth(user authentication_response.UserContent, contextModel *applicationModel.ContextModel, resourceID string) (err errorModel.ErrorModel) {
	var (
		fileName       = "InsertUserCheckResource.go"
		funcName       = "hitAuthenticationServerForAddResourceNextracUserAuth"
		locale         = constanta.DefaultApplicationsLanguage
		auth           = config.ApplicationConfiguration.GetAuthenticationServer()
		url            = auth.Host + auth.PathRedirect.AddResourceClient
		method         = http.MethodPost
		structResponse out.APIResponseAddResource
	)

	authenticationRequest := in.UserRequest{
		ClientID:   user.ClientID,
		ResourceID: resourceID,
	}

	authenticationRequestJSON := util.StructToJSON(authenticationRequest)
	fmt.Println("request add resource ---> ", authenticationRequestJSON)
	internalToken := resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, "", "", locale)

	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errs := common.HitAPI(url, header, authenticationRequestJSON, method, *contextModel)
	fmt.Println("status code ---> ", statusCode)
	fmt.Println("body result ---> ", bodyResult)
	_ = json.Unmarshal([]byte(bodyResult), &structResponse)

	if errs != nil {
		err = errorModel.GenerateErrorModel(statusCode, errs.Error(), fileName, funcName, errs)
		return
	}

	if statusCode != 200 {
		err = errorModel.GenerateErrorModel(statusCode, structResponse.Nexsoft.Payload.Status.Message, fileName, funcName, errs)
		fmt.Println("request id to auth ---> ", structResponse.Nexsoft.Header.RequestID)
		return
	}

	return
}
