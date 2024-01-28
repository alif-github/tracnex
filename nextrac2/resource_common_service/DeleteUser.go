package resource_common_service

import (
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"strconv"
)

func InternalDeleteUser(authUserID int64, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var viewUserByIDResponse authentication_response.UserAuthenticationResponse

	token := GenerateInternalToken("auth", 0, contextModel.AuthAccessTokenModel.ClientID, "delete user", constanta.DefaultApplicationsLanguage)
	path := config.ApplicationConfiguration.GetAuthenticationServer().Host + config.ApplicationConfiguration.GetAuthenticationServer().PathRedirect.InternalUser.CrudUser
	path += "/" + strconv.Itoa(int(authUserID))

	viewUserByIDResponse, err = internalGetUserByID(token, path, contextModel)
	if err.Error != nil {
		return
	}

	bodyRequest := authentication_request.DeleteUserRequestDTO{
		ResourceID: config.ApplicationConfiguration.GetServerResourceID(),
		UpdatedAt:  viewUserByIDResponse.Nexsoft.Payload.Data.Content.UpdatedAt,
	}

	err = internalDeleteUser(token, path, util.StructToJSON(bodyRequest), contextModel)
	return
}

func InternalDeleteClientByClientID(clientID string, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		viewClientByClientResponse authentication_response.UserAuthenticationResponse
		bodyRequest                authentication_request.DeleteUserRequestDTO
		token                      string
		path                       string
		issuer                     = "delete user"
		resourceDestination        = "auth"
	)

	token = GenerateInternalToken(resourceDestination, 0, contextModel.AuthAccessTokenModel.ClientID, issuer, constanta.DefaultApplicationsLanguage)
	path = config.ApplicationConfiguration.GetAuthenticationServer().Host + config.ApplicationConfiguration.GetAuthenticationServer().PathRedirect.InternalClient.CrudClient
	path += fmt.Sprintf(`/%s`, clientID)

	viewClientByClientResponse, err = internalGetClientByClientID(token, path, contextModel)
	if err.Error != nil {
		return
	}

	bodyRequest = authentication_request.DeleteUserRequestDTO{
		ResourceID: config.ApplicationConfiguration.GetServerResourceID(),
		UpdatedAt:  viewClientByClientResponse.Nexsoft.Payload.Data.Content.UpdatedAt,
	}

	err = internalDeleteClientByClientID(token, path, util.StructToJSON(bodyRequest), contextModel)
	return
}

func InternalGetUserByID(authUserID int64, contextModel *applicationModel.ContextModel) (viewUserByIDResponse authentication_response.UserAuthenticationResponse, err errorModel.ErrorModel) {
	token := GenerateInternalToken(constanta.AuthDestination, 0, contextModel.AuthAccessTokenModel.ClientID, "get user", constanta.DefaultApplicationsLanguage)
	path := config.ApplicationConfiguration.GetAuthenticationServer().Host + config.ApplicationConfiguration.GetAuthenticationServer().PathRedirect.InternalUser.CrudUser
	path += "/" + strconv.Itoa(int(authUserID))
	return internalGetUserByID(token, path, contextModel)
}

func internalGetUserByID(token string, path string, contextModel *applicationModel.ContextModel) (viewUserByIDResponse authentication_response.UserAuthenticationResponse, err errorModel.ErrorModel) {
	var statusCode int
	var bodyResult string
	var errs error

	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{token}
	statusCode, _, bodyResult, errs = common.HitAPI(path, header, "", "GET", *contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError("DeleteUser.go", "InternalDeleteUser", errs)
		return
	}

	if statusCode != 200 {
		err = common.ReadAuthServerError("InternalDeleteUser", statusCode, bodyResult, contextModel)
		return
	} else {
		_ = json.Unmarshal([]byte(bodyResult), &viewUserByIDResponse)
		return
	}
}

func internalGetClientByClientID(token string, path string, contextModel *applicationModel.ContextModel) (viewUserByIDResponse authentication_response.UserAuthenticationResponse, err errorModel.ErrorModel) {
	var statusCode int
	var bodyResult string
	var errs error

	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{token}
	statusCode, _, bodyResult, errs = common.HitAPI(path, header, "", "GET", *contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError("DeleteUser.go", "internalGetClientByClientID", errs)
		return
	}

	if statusCode != 200 {
		err = common.ReadAuthServerError("internalGetClientByClientID", statusCode, bodyResult, contextModel)
		return
	} else {
		_ = json.Unmarshal([]byte(bodyResult), &viewUserByIDResponse)
		return
	}
}

func internalDeleteUser(token string, path string, body string, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var statusCode int
	var bodyResult string
	var errs error

	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{token}
	statusCode, _, bodyResult, errs = common.HitAPI(path, header, body, "DELETE", *contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError("DeleteUser.go", "InternalDeleteUser", errs)
		return
	}

	if statusCode != 200 {
		err = common.ReadAuthServerError("InternalDeleteUser", statusCode, bodyResult, contextModel)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func internalDeleteClientByClientID(token string, path string, body string, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var statusCode int
	var bodyResult string
	var errs error

	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{token}
	statusCode, _, bodyResult, errs = common.HitAPI(path, header, body, "DELETE", *contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError("DeleteUser.go", "internalDeleteClientByClientID", errs)
		return
	}

	if statusCode != 200 {
		err = common.ReadAuthServerError("internalDeleteClientByClientID", statusCode, bodyResult, contextModel)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
