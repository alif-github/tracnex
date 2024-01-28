package endpoint

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/token"
	util2 "nexsoft.co.id/nextrac2/util"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type AbstractEndpoint struct {
	FileName string
	IsAdmin  bool
}

func (input AbstractEndpoint) ServeWhiteListEndpoint(funcName string, isCheckSignature bool, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	serve(input.FileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func (input AbstractEndpoint) ServeWhiteListEndpointWithFile(funcName string, isCheckSignature bool, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (*os.File, map[string]string, errorModel.ErrorModel)) {
	serveWithFile(input.FileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func (input AbstractEndpoint) ServeWhiteListConstantTokenEndpoint(funcName string, isCheckSignature bool, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	if request.Header.Get(constanta.DefaultTokenKeyConstanta) != constanta.DefaultTokenValueConstanta {
		contextModel := request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)
		err := errorModel.GenerateForbiddenAccessClientError("AbstractEndpoint.go", "ServeWhiteListConstantTokenEndpoint")
		writeErrorResponse(responseWriter, err, contextModel, out.Payload{})
		return
	}
	serve(input.FileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func (input AbstractEndpoint) ServeInternalValidationEndpoint(funcName string, isCheckSignature bool, checkClientID bool, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	contextModel := request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)
	defer func() {
		ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, contextModel)
		request = request.WithContext(ctx)
	}()

	contextModel.IsInternal = true
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	addResourceClientUrl := authenticationServer.Host + authenticationServer.PathRedirect.AddResourceClient

	authAccessModel, err := resource_common_service.ValidateJWTInternal(checkClientID, serverconfig.ServerAttribute.RedisClient, addResourceClientUrl, request.Header.Get(constanta.TokenHeaderNameConstanta), config.ApplicationConfiguration.GetJWTToken().Internal, config.ApplicationConfiguration.GetServerResourceID(), resource_common_service.GenerateInternalToken("auth", 0, "", "CheckToken", "en-EN"), RoleMappingInternal, SaveClientToDB, contextModel)
	contextModel.LoggerModel.ClientID = authAccessModel.ClientID
	contextModel.LoggerModel.UserID = strconv.Itoa(int(authAccessModel.AuthenticationServerUserID))
	contextModel.AuthAccessTokenModel = authAccessModel
	contextModel.IsSignatureCheck = isCheckSignature
	contextModel.PermissionHave = authAccessModel.Authentication

	if err.Error != nil {
		readError(err, contextModel, responseWriter)
		return
	}

	serve(input.FileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func (input AbstractEndpoint) ServeJWTTokenWithIPWhitelistValidation(funcName string, isCheckSignature bool, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	input.serveJWTTokenValidationEndpoint(input.FileName, isCheckSignature, true, funcName, scope, permissionMustHave, responseWriter, request, serveFunction)
}
func (input AbstractEndpoint) ServeJWTTokenValidationEndpoint(funcName string, isCheckSignature bool, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	input.serveJWTTokenValidationEndpoint(input.FileName, isCheckSignature, false, funcName, scope, permissionMustHave, responseWriter, request, serveFunction)
}

func (input AbstractEndpoint) ServeJWTTokenWithNexsoftAdminValidationEndpoint(funcName string, isCheckSignature bool, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	fmt.Println(funcName)
	input.serveJWTTokenNexsoftValidationEndpoint(input.FileName, isCheckSignature, false, funcName, scope, permissionMustHave, responseWriter, request, serveFunction)
}

func (input *AbstractEndpoint) serveJWTTokenValidationEndpoint(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	input.IsAdmin = false
	input.serveJWTTokenValidation(fileName, isCheckSignature, checkIPWhitelist, funcName, scope, permissionMustHave, responseWriter, request, serveFunction, RoleMappingUser)
}

func (input *AbstractEndpoint) serveJWTTokenNexsoftValidationEndpoint(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	input.IsAdmin = true
	input.serveJWTTokenValidation(fileName, isCheckSignature, checkIPWhitelist, funcName, scope, permissionMustHave, responseWriter, request, serveFunction, RoleMappingUserNexsoft)
}

func (input AbstractEndpoint) serveJWTTokenValidation(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel), roleMap func(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error)) {
	contextModel := request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)

	contextModel.IsInternal = false
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	checkTokenUrl := authenticationServer.Host + authenticationServer.PathRedirect.CheckToken

	authAccessModel, err := resource_common_service.ValidateJWTToken(serverconfig.ServerAttribute.RedisClient, request.Header.Get(constanta.TokenHeaderNameConstanta), constanta.ExpiredTokenOnRedisConstanta, checkTokenUrl, config.ApplicationConfiguration.GetServerResourceID(), scope, config.ApplicationConfiguration.GetJWTToken().JWT, roleMap, contextModel)
	contextModel.LoggerModel.ClientID = authAccessModel.ClientID
	contextModel.LoggerModel.UserID = strconv.Itoa(int(authAccessModel.AuthenticationServerUserID))
	contextModel.AuthAccessTokenModel = authAccessModel

	if err.Error != nil {
		if err.CausedBy != nil {
			fmt.Println("[Error Validate JWT Token] => ", err.CausedBy.Error())
		}

		readError(err, contextModel, responseWriter)
		return
	}

	var authenticationModel model2.AuthenticationModel
	var permissionHave string
	var errors errorModel.ErrorModel

	if checkIPWhitelist {
		if !common.ValidateIPWhitelist(contextModel.AuthAccessTokenModel.IPWhiteList, contextModel.LoggerModel.IP) {
			errors = errorModel.GenerateUnauthorizedClientError(fileName, funcName)
			writeErrorResponse(responseWriter, errors, contextModel, out.Payload{})
			return
		}
	}

	_ = json.Unmarshal([]byte(authAccessModel.Authentication), &authenticationModel)
	if permissionMustHave != "" {
		permissionHave, errors = ValidatePermissionWithRole(permissionMustHave, authenticationModel.Role)
		if errors.Error != nil {
			fmt.Println("ada error 2 :", errors.Error.Error())
			fmt.Println("permissionHave :", permissionHave)
			fmt.Println("permissionMustHave :", permissionMustHave)
			writeErrorResponse(responseWriter, errors, contextModel, out.Payload{})
			return
		}
	}

	contextModel.AuthAccessTokenModel = authAccessModel
	contextModel.PermissionHave = permissionHave
	contextModel.IsAdmin = input.IsAdmin

	ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, contextModel)
	request = request.WithContext(ctx)

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		contextModel.LimitedByCreatedBy = userID
	}

	serve(fileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func serve(fileName string, funcName string, isCheckSignature bool, responseWriter http.ResponseWriter, request *http.Request, serve func(*http.Request, *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel)) {
	var err errorModel.ErrorModel
	var contextModel *applicationModel.ContextModel
	var output out.Payload
	var header map[string]string

	defer func() {
		if r := recover(); r != nil {
			err = errorModel.GenerateRecoverError()
			contextModel.LoggerModel.Message = string(debug.Stack())
		} else {
			if err.Error != nil {
				contextModel.LoggerModel.Class = "[" + err.FileName + "," + err.FuncName + "]"
				contextModel.LoggerModel.Code = err.Error.Error()
				if err.CausedBy != nil {
					contextModel.LoggerModel.Message = err.CausedBy.Error()
				} else {
					contextModel.LoggerModel.Message = err.Error.Error()
				}
			}
		}

		contextModel.LoggerModel.Status = err.Code
		ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, contextModel)
		request = request.WithContext(ctx)
		finish(request, responseWriter, err, contextModel, output)
	}()

	contextModel = request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)
	contextModel.LoggerModel.Class = "[" + fileName + "," + funcName + "]"
	contextModel.IsSignatureCheck = isCheckSignature
	getDBSchema(contextModel)

	output, header, err = serve(request, contextModel)
	if err.Error != nil {
		return
	}

	setHeader(header, responseWriter)
}

func setHeader(header map[string]string, responseWriter http.ResponseWriter) {
	accessControlExpose := "Access-Control-Expose-Headers"
	accessControlAllow := "Access-Control-Allow-Headers"

	exposeHeader := responseWriter.Header().Get(accessControlExpose)
	allowHeader := responseWriter.Header().Get(accessControlAllow)
	for key := range header {
		lowerKey := strings.ToLower(key)
		responseWriter.Header().Add(key, header[key])
		if exposeHeader == "" {
			exposeHeader = lowerKey
		} else {
			exposeHeader += ", " + lowerKey
		}
		if lowerKey != "authorization" {
			if allowHeader == "" {
				allowHeader = lowerKey
			} else {
				allowHeader += ", " + lowerKey
			}
		}
	}
	if exposeHeader != "" {
		responseWriter.Header().Set(accessControlExpose, exposeHeader)
		responseWriter.Header().Set(accessControlAllow, allowHeader)
	}
}

func finish(request *http.Request, responseWriter http.ResponseWriter, err errorModel.ErrorModel, contextModel *applicationModel.ContextModel, output out.Payload) {
	if err.Error != nil {
		writeErrorResponse(responseWriter, err, contextModel, output)
	} else {
		writeSuccessResponse(request, responseWriter, contextModel, output)
	}
}

func writeErrorResponse(responseWriter http.ResponseWriter, err errorModel.ErrorModel, contextModel *applicationModel.ContextModel, payload out.Payload) {
	if err.Code == 0 {
		responseWriter.WriteHeader(500)
		err.CausedBy = err.Error
		err.Error = errors.New("E-5-TRAC-SRV-001")
	} else {
		responseWriter.WriteHeader(err.Code)
		if err.Code == http.StatusInternalServerError {
			//--- Discord Send Thread
			errS := util2.DiscordSendThread(*contextModel)
			if errS.Error != nil {
				fmt.Println(fmt.Sprintf(`Discord Error -> %s`, errS.CausedBy))
			}
		}
	}

	errCode := err.Error.Error()
	errMessage := util2.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale)
	if errMessage == errCode {
		if err.CausedBy != nil {
			errMessage = err.CausedBy.Error()
		}
	}

	errResponse := out.StatusResponse{
		Success: false,
		Code:    errCode,
		Message: errMessage,
	}

	responseMessage := out.APIResponse{
		Nexsoft: out.NexsoftMessage{
			Header: out.Header{
				RequestID: contextModel.LoggerModel.RequestID,
				Version:   config.ApplicationConfiguration.GetServerVersion(),
				Timestamp: util.GetTimeStamp(),
			},
			Payload: out.Payload{Status: errResponse}},
	}

	if payload.Other != nil {
		responseMessage.Nexsoft.Payload.Other = payload.Other
	}

	if err.AdditionalInformation != nil && len(err.AdditionalInformation) > 0 {
		responseMessage.Nexsoft.Payload.Status.Detail = err.AdditionalInformation
	}

	if err.OtherData != nil {
		responseMessage.Nexsoft.Payload.Other = err.OtherData
	}

	_, errorS := responseWriter.Write([]byte(responseMessage.String()))
	if errorS != nil {
		errModel := errorModel.GenerateUnknownError("AbstractEndpoint.go", "writeErrorResponse", errorS)
		contextModel.LoggerModel.Status = errModel.Code
		contextModel.LoggerModel.Code = errModel.Error.Error()
		contextModel.LoggerModel.Message = errorS.Error()
	}

	contextModel.LoggerModel.ByteOut = len([]byte(responseMessage.String()))
}

func writeSuccessResponse(request *http.Request, responseWriter http.ResponseWriter, contextModel *applicationModel.ContextModel, output out.Payload) {
	if output.Status.Success {
		output.Status.Success = false
	} else {
		output.Status.Success = true
	}

	responseMessage := out.APIResponse{
		Nexsoft: out.NexsoftMessage{
			Header: out.Header{
				RequestID: contextModel.LoggerModel.RequestID,
				Version:   config.ApplicationConfiguration.GetServerVersion(),
				Timestamp: util.GetTimeStamp(),
			},
			Payload: output},
	}
	bodyMessage := responseMessage.String()
	if contextModel.IsSignatureCheck {
		setHeader(GenerateSignature(bodyMessage, contextModel.AuthAccessTokenModel.SignatureKey, request), responseWriter)
	}

	responseWriter.WriteHeader(200)
	_, errorS := responseWriter.Write([]byte(bodyMessage))
	if errorS != nil {
		errModel := errorModel.GenerateUnknownError("AbstractEndpoint.go", "writeSuccessResponse", errorS)
		contextModel.LoggerModel.Status = errModel.Code
		contextModel.LoggerModel.Code = errModel.Error.Error()
		contextModel.LoggerModel.Message = errorS.Error()
	}
	contextModel.LoggerModel.ByteOut = len([]byte(responseMessage.String()))
}

func roleMapping(clientID string, token string, payload token.PayloadJWTToken, roleMap func(db *sql.DB, userParam repository.UserModel) (results repository.RoleMappingPersonProfileModel, err errorModel.ErrorModel)) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error) {
	var (
		authenticationRoleModel model2.AuthenticationRoleModel
		authenticationDataModel model2.AuthenticationDataModel
		authAccessTokenModel    model2.AuthAccessTokenModel
		userModel               repository.UserModel
		db                      = serverconfig.ServerAttribute.DBConnection
		temp                    = make(map[string][]string)
		tempDataScope           = make(map[string]interface{})
	)

	userModel = repository.UserModel{ClientID: sql.NullString{String: clientID}}
	roleModel, err := roleMap(db, userModel)
	if err.Error != nil {
		return authAccessTokenModel, authenticationRoleModel, authenticationDataModel, err.CausedBy
	}

	authAccessTokenModel.ResourceUserID = roleModel.PersonProfileID.Int64
	authAccessTokenModel.IsAdmin = roleModel.IsAdmin.Bool
	authAccessTokenModel.RedisAuthAccessTokenModel = model2.RedisAuthAccessTokenModel{
		ResourceUserID: roleModel.PersonProfileID.Int64,
		IPWhiteList:    roleModel.IPWhitelist.String,
		SignatureKey:   roleModel.SignatureKey.String,
	}

	_ = json.Unmarshal([]byte(roleModel.Permissions.String), &temp)
	authenticationRoleModel.Role = roleModel.RoleName.String
	authenticationRoleModel.Permission = temp

	_ = json.Unmarshal([]byte(roleModel.Scope.String), &tempDataScope)
	authenticationDataModel.Group = roleModel.GroupName.String
	authenticationDataModel.Scope = tempDataScope

	tx, errs := serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		return authAccessTokenModel, authenticationRoleModel, authenticationDataModel, errs
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
		return
	}()

	err = dao.ClientTokenDAO.InsertClientToken(tx, repository.ClientTokenModel{
		ClientID:      sql.NullString{String: clientID},
		AuthUserID:    roleModel.AuthUserID,
		Token:         sql.NullString{String: token},
		ExpiredAt:     sql.NullTime{Time: time.Unix(payload.ExpiresAt, 0)},
		CreatedBy:     roleModel.AuthUserID,
		CreatedClient: sql.NullString{String: clientID},
	})

	if err.Error != nil {
		return authAccessTokenModel, authenticationRoleModel, authenticationDataModel, err.CausedBy
	}

	userModel.ClientID.String = clientID
	err = dao.UserDAO.UpdateLastTokenUser(tx, userModel)
	if err.Error != nil {
		return authAccessTokenModel, authenticationRoleModel, authenticationDataModel, err.CausedBy
	}

	return authAccessTokenModel, authenticationRoleModel, authenticationDataModel, nil
}

func RoleMappingUser(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error) {
	return roleMapping(clientID, token, payload, dao.UserDAO.RoleMappingUser)
}
func RoleMappingUserNexsoft(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error) {
	return roleMapping(clientID, token, payload, dao.UserDAO.RoleMappingUserNexsoft)
}

func RoleMappingInternal(clientID string, userID int64) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error) {
	var authAccessTokenModel model2.AuthAccessTokenModel
	var authenticationRoleModel model2.AuthenticationRoleModel
	var authenticationDataModel model2.AuthenticationDataModel

	userModel := repository.UserModel{
		AuthUserID: sql.NullInt64{Int64: userID},
		ClientID:   sql.NullString{String: clientID}}

	mapping, err := dao.UserDAO.RoleMappingInternalUser(serverconfig.ServerAttribute.DBConnection, userModel)
	if err.Error != nil {
		return authAccessTokenModel, authenticationRoleModel, authenticationDataModel, err.CausedBy
	}

	authAccessTokenModel = model2.AuthAccessTokenModel{
		RedisAuthAccessTokenModel: model2.RedisAuthAccessTokenModel{
			ResourceUserID: mapping.PersonProfileID.Int64,
			IPWhiteList:    mapping.IPWhitelist.String,
			SignatureKey:   mapping.SignatureKey.String,
		},
		ClientID:                   clientID,
		AuthenticationServerUserID: mapping.AuthUserID.Int64,
		Locale:                     mapping.Locale.String,
	}

	_ = json.Unmarshal([]byte(mapping.Permissions.String), &authenticationRoleModel.Permission)
	authenticationRoleModel.Role = mapping.RoleName.String

	return authAccessTokenModel, authenticationRoleModel, authenticationDataModel, nil
}

func SaveClientToDB(addClientResourceResult authentication_response.AddClientAuthenticationResponse, createdBy int64) (int64, error) {
	additionalInfoResult := addClientResourceResult.Nexsoft.Payload.Data.Content
	tx, errors := serverconfig.ServerAttribute.DBConnection.Begin()
	if errors != nil {
		return 0, errors
	}

	defer func() {
		if errors != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
		return
	}()

	_ = repository.UserModel{
		AuthUserID:   sql.NullInt64{Int64: additionalInfoResult.UserID},
		ClientID:     sql.NullString{String: additionalInfoResult.ClientID},
		SignatureKey: sql.NullString{String: additionalInfoResult.SignatureKey},
		Locale:       sql.NullString{String: additionalInfoResult.Locale},
		CreatedBy:    sql.NullInt64{Int64: createdBy},
		UpdatedBy:    sql.NullInt64{Int64: createdBy},
	}

	if additionalInfoResult.UserID > 0 {
		clientRoleModel := repository.ClientRoleScopeModel{
			ClientID:  sql.NullString{String: additionalInfoResult.ClientID},
			RoleID:    sql.NullInt64{Int64: 1},
			GroupID:   sql.NullInt64{Int64: 1},
			CreatedBy: sql.NullInt64{Int64: createdBy},
			UpdatedBy: sql.NullInt64{Int64: createdBy},
		}
		_, err := dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, clientRoleModel)
		if err.Error != nil {
			return 0, err.CausedBy
		}
	}

	return 0, nil
}

func readError(err model2.ResourceCommonErrorModel, contextModel *applicationModel.ContextModel, responseWriter http.ResponseWriter) {
	errs := ReadError(err)
	contextModel.LoggerModel.Status = errs.Code
	contextModel.LoggerModel.Code = errs.Error.Error()
	if errs.CausedBy != nil {
		contextModel.LoggerModel.Message = errs.CausedBy.Error()
	}
	writeErrorResponse(responseWriter, errs, contextModel, out.Payload{})
}

func getDBSchema(model *applicationModel.ContextModel) {
	model.DBSchema = config.ApplicationConfiguration.GetPostgreSQLDefaultSchema()
}

func (input AbstractEndpoint) ServeJWTTokenWithNexsoftAdminValidationEndpointWithFileResult(funcName string, isCheckSignature bool, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (*os.File, map[string]string, errorModel.ErrorModel)) {
	serveJWTTokenNexsoftValidationEndpointWithFileResult(input.FileName, isCheckSignature, false, funcName, scope, permissionMustHave, responseWriter, request, serveFunction)
}

func (input AbstractEndpoint) ServeJWTTokenWithValidationEndpointWithFileCSVResult(funcName string, isCheckSignature bool, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) ([][]string, map[string]string, errorModel.ErrorModel)) {
	serveJWTTokenValidationEndpointWithFileCSVResult(input.FileName, isCheckSignature, false, funcName, scope, permissionMustHave, responseWriter, request, serveFunction)
}

func serveJWTTokenNexsoftValidationEndpointWithFileResult(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (*os.File, map[string]string, errorModel.ErrorModel)) {
	serveJWTTokenValidationWithFileResponse(fileName, isCheckSignature, checkIPWhitelist, funcName, scope, permissionMustHave, responseWriter, request, serveFunction, RoleMappingUserNexsoft)
}

func serveJWTTokenValidationEndpointWithFileCSVResult(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) ([][]string, map[string]string, errorModel.ErrorModel)) {
	serveJWTTokenValidationWithFileCSVResponse(fileName, isCheckSignature, checkIPWhitelist, funcName, scope, permissionMustHave, responseWriter, request, serveFunction, RoleMappingUser)
}

func serveJWTTokenValidationWithFileResponse(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) (*os.File, map[string]string, errorModel.ErrorModel), roleMap func(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error)) {
	contextModel := request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)

	contextModel.IsInternal = false
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	checkTokenUrl := authenticationServer.Host + authenticationServer.PathRedirect.CheckToken

	authAccessModel, err := resource_common_service.ValidateJWTToken(serverconfig.ServerAttribute.RedisClient, request.Header.Get(constanta.TokenHeaderNameConstanta), constanta.ExpiredTokenOnRedisConstanta, checkTokenUrl, config.ApplicationConfiguration.GetServerResourceID(), scope, config.ApplicationConfiguration.GetJWTToken().JWT, roleMap, contextModel)
	contextModel.LoggerModel.ClientID = authAccessModel.ClientID
	contextModel.LoggerModel.UserID = strconv.Itoa(int(authAccessModel.AuthenticationServerUserID))
	contextModel.AuthAccessTokenModel = authAccessModel

	if err.Error != nil {
		readError(err, contextModel, responseWriter)
		return
	}

	var authenticationModel model2.AuthenticationModel
	var permissionHave string
	var errors errorModel.ErrorModel

	if checkIPWhitelist {
		if !common.ValidateIPWhitelist(contextModel.AuthAccessTokenModel.IPWhiteList, contextModel.LoggerModel.IP) {
			errors = errorModel.GenerateUnauthorizedClientError(fileName, funcName)
			writeErrorResponse(responseWriter, errors, contextModel, out.Payload{})
			return
		}
	}

	_ = json.Unmarshal([]byte(authAccessModel.Authentication), &authenticationModel)
	if permissionMustHave != "" {
		permissionHave, errors = ValidatePermissionWithRole(permissionMustHave, authenticationModel.Role)
		if errors.Error != nil {
			writeErrorResponse(responseWriter, errors, contextModel, out.Payload{})
			return
		}
	}

	contextModel.AuthAccessTokenModel = authAccessModel
	contextModel.PermissionHave = permissionHave

	ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, contextModel)
	request = request.WithContext(ctx)

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		contextModel.LimitedByCreatedBy = userID
	}

	serveWithFile(fileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func serveJWTTokenValidationWithFileCSVResponse(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) ([][]string, map[string]string, errorModel.ErrorModel), roleMap func(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error)) {
	contextModel := request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)

	contextModel.IsInternal = false
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	checkTokenUrl := authenticationServer.Host + authenticationServer.PathRedirect.CheckToken

	authAccessModel, err := resource_common_service.ValidateJWTToken(serverconfig.ServerAttribute.RedisClient, request.Header.Get(constanta.TokenHeaderNameConstanta), constanta.ExpiredTokenOnRedisConstanta, checkTokenUrl, config.ApplicationConfiguration.GetServerResourceID(), scope, config.ApplicationConfiguration.GetJWTToken().JWT, roleMap, contextModel)
	contextModel.LoggerModel.ClientID = authAccessModel.ClientID
	contextModel.LoggerModel.UserID = strconv.Itoa(int(authAccessModel.AuthenticationServerUserID))
	contextModel.AuthAccessTokenModel = authAccessModel

	if err.Error != nil {
		readError(err, contextModel, responseWriter)
		return
	}

	var authenticationModel model2.AuthenticationModel
	var permissionHave string
	var errors errorModel.ErrorModel

	if checkIPWhitelist {
		if !common.ValidateIPWhitelist(contextModel.AuthAccessTokenModel.IPWhiteList, contextModel.LoggerModel.IP) {
			errors = errorModel.GenerateUnauthorizedClientError(fileName, funcName)
			writeErrorResponse(responseWriter, errors, contextModel, out.Payload{})
			return
		}
	}

	_ = json.Unmarshal([]byte(authAccessModel.Authentication), &authenticationModel)
	if permissionMustHave != "" {
		permissionHave, errors = ValidatePermissionWithRole(permissionMustHave, authenticationModel.Role)
		if errors.Error != nil {
			writeErrorResponse(responseWriter, errors, contextModel, out.Payload{})
			return
		}
	}

	contextModel.AuthAccessTokenModel = authAccessModel
	contextModel.PermissionHave = permissionHave

	ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, contextModel)
	request = request.WithContext(ctx)

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		contextModel.LimitedByCreatedBy = userID
	}

	serveWithCSV(fileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func serveWithFile(fileName string, funcName string, isCheckSignature bool, responseWriter http.ResponseWriter, request *http.Request, serve func(*http.Request, *applicationModel.ContextModel) (*os.File, map[string]string, errorModel.ErrorModel)) {
	var err errorModel.ErrorModel
	var contextModel *applicationModel.ContextModel
	var output *os.File
	var header map[string]string

	defer func() {
		if r := recover(); r != nil {
			err = errorModel.GenerateRecoverError()
			contextModel.LoggerModel.Message = string(debug.Stack())
		} else {
			if err.Error != nil {
				contextModel.LoggerModel.Class = "[" + err.FileName + "," + err.FuncName + "]"
				contextModel.LoggerModel.Code = err.Error.Error()
				if err.CausedBy != nil {
					contextModel.LoggerModel.Message = err.CausedBy.Error()
				} else {
					contextModel.LoggerModel.Message = err.Error.Error()
				}
			}
		}

		contextModel.LoggerModel.Status = err.Code
		ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, contextModel)
		request = request.WithContext(ctx)
		finishWithFile(request, responseWriter, err, contextModel, output)
	}()

	contextModel = request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)
	contextModel.LoggerModel.Class = "[" + fileName + "," + funcName + "]"
	contextModel.IsSignatureCheck = isCheckSignature
	getDBSchema(contextModel)

	output, header, err = serve(request, contextModel)
	if err.Error != nil {
		return
	}

	setHeader(header, responseWriter)
}

func serveWithCSV(fileName string, funcName string, isCheckSignature bool, responseWriter http.ResponseWriter, request *http.Request, serve func(*http.Request, *applicationModel.ContextModel) ([][]string, map[string]string, errorModel.ErrorModel)) {
	var (
		err          errorModel.ErrorModel
		contextModel *applicationModel.ContextModel
		output       [][]string
		header       map[string]string
	)

	defer func() {
		if r := recover(); r != nil {
			err = errorModel.GenerateRecoverError()
			contextModel.LoggerModel.Message = string(debug.Stack())
		} else {
			if err.Error != nil {
				contextModel.LoggerModel.Class = "[" + err.FileName + "," + err.FuncName + "]"
				contextModel.LoggerModel.Code = err.Error.Error()
				if err.CausedBy != nil {
					contextModel.LoggerModel.Message = err.CausedBy.Error()
				} else {
					contextModel.LoggerModel.Message = err.Error.Error()
				}
			}
		}

		contextModel.LoggerModel.Status = err.Code
		ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, contextModel)
		request = request.WithContext(ctx)
		finishWithCSV(request, responseWriter, err, contextModel, output)
	}()

	contextModel = request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)
	contextModel.LoggerModel.Class = "[" + fileName + "," + funcName + "]"
	contextModel.IsSignatureCheck = isCheckSignature
	getDBSchema(contextModel)

	output, header, err = serve(request, contextModel)
	if err.Error != nil {
		return
	}

	setHeader(header, responseWriter)
}

func finishWithFile(request *http.Request, responseWriter http.ResponseWriter, err errorModel.ErrorModel, contextModel *applicationModel.ContextModel, output *os.File) {
	if err.Error != nil {
		writeErrorResponse(responseWriter, err, contextModel, out.Payload{})
	} else {
		writeSuccessGetFileResponse(request, responseWriter, contextModel, output)
	}
}

func finishWithCSV(request *http.Request, responseWriter http.ResponseWriter, err errorModel.ErrorModel, contextModel *applicationModel.ContextModel, output [][]string) {
	if err.Error != nil {
		writeErrorResponse(responseWriter, err, contextModel, out.Payload{})
	} else {
		writeSuccessGetCSVResponse(request, responseWriter, contextModel, output)
	}
}

func writeSuccessGetFileResponse(_ *http.Request, responseWriter http.ResponseWriter, contextModel *applicationModel.ContextModel, output *os.File) {
	responseWriter.WriteHeader(http.StatusOK)

	defer func() {
		if output != nil {
			_ = output.Close()
		}
	}()

	length, errorS := io.Copy(responseWriter, output)
	if errorS != nil {
		errModel := errorModel.GenerateUnknownError("AbstractEndpoint.go", "writeSuccessResponse", errorS)
		contextModel.LoggerModel.Status = errModel.Code
		contextModel.LoggerModel.Code = errModel.Error.Error()
		contextModel.LoggerModel.Message = errorS.Error()
		return
	}

	contextModel.LoggerModel.ByteOut = int(length)
}

func writeSuccessGetCSVResponse(_ *http.Request, responseWriter http.ResponseWriter, contextModel *applicationModel.ContextModel, output [][]string) {
	responseWriter.Header().Set("Content-Type", "text/csv")
	responseWriter.WriteHeader(http.StatusOK)
	wr := csv.NewWriter(responseWriter)
	wr.Comma = constanta.PipaDelimiter

	errorS := wr.WriteAll(output)
	if errorS != nil {
		errModel := errorModel.GenerateUnknownError("AbstractEndpoint.go", "writeSuccessResponse", errorS)
		contextModel.LoggerModel.Status = errModel.Code
		contextModel.LoggerModel.Code = errModel.Error.Error()
		contextModel.LoggerModel.Message = errorS.Error()
	}
}
