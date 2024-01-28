package endpoint

import (
	"context"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/token"
	"runtime/debug"
	"strconv"
)

func (input AbstractEndpoint) ServeJWTTokenValidationEndpointWithFileResult(funcName string, isCheckSignature bool, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) ([]byte, map[string]string, errorModel.ErrorModel)) {
	serveJWTTokenValidationEndpointWithFileResult(input.FileName, isCheckSignature, false, funcName, scope, permissionMustHave, responseWriter, request, serveFunction)
}

func serveJWTTokenValidationEndpointWithFileResult(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) ([]byte, map[string]string, errorModel.ErrorModel)) {
	serveJWTTokenValidationWithFileBytesResponse(fileName, isCheckSignature, checkIPWhitelist, funcName, scope, permissionMustHave, responseWriter, request, serveFunction, RoleMappingUser)
}

func serveJWTTokenValidationWithFileBytesResponse(fileName string, isCheckSignature bool, checkIPWhitelist bool, funcName string, scope string, permissionMustHave string, responseWriter http.ResponseWriter, request *http.Request, serveFunction func(*http.Request, *applicationModel.ContextModel) ([]byte, map[string]string, errorModel.ErrorModel), roleMap func(clientID string, token string, payload token.PayloadJWTToken) (model2.AuthAccessTokenModel, model2.AuthenticationRoleModel, model2.AuthenticationDataModel, error)) {
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

	serveWithFileBytes(fileName, funcName, isCheckSignature, responseWriter, request, serveFunction)
}

func serveWithFileBytes(fileName string, funcName string, isCheckSignature bool, responseWriter http.ResponseWriter, request *http.Request, serve func(*http.Request, *applicationModel.ContextModel) ([]byte, map[string]string, errorModel.ErrorModel)) {
	var err errorModel.ErrorModel
	var contextModel *applicationModel.ContextModel
	var output []byte
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
		finishWithFileBytes(request, responseWriter, err, contextModel, output)
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


func finishWithFileBytes(request *http.Request, responseWriter http.ResponseWriter, err errorModel.ErrorModel, contextModel *applicationModel.ContextModel, output []byte) {
	if err.Error != nil {
		writeErrorResponse(responseWriter, err, contextModel, out.Payload{})
	} else {
		writeSuccessGetFileBytesResponse(request, responseWriter, contextModel, output)
	}
}

func writeSuccessGetFileBytesResponse(_ *http.Request, responseWriter http.ResponseWriter, contextModel *applicationModel.ContextModel, output []byte) {
	responseWriter.WriteHeader(200)
	_, errorS := responseWriter.Write(output)
	if errorS != nil {
		errModel := errorModel.GenerateUnknownError("AbstractEndpoint.go", "writeSuccessResponse", errorS)
		contextModel.LoggerModel.Status = errModel.Code
		contextModel.LoggerModel.Code = errModel.Error.Error()
		contextModel.LoggerModel.Message = errorS.Error()
	}
	contextModel.LoggerModel.ByteOut = len(output)
}
