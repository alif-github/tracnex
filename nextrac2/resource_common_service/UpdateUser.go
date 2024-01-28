package resource_common_service

import (
	"encoding/json"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"strconv"
	"strings"
)

func InternalUpdateUser(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel, emailMessage string) (updateResposne authentication_response.UpdateUserAuthenticationResponse, err errorModel.ErrorModel) {
	var (
		statusCode         int
		bodyResult         string
		errs               error
		token              string
		pathUpdateUser     string
		pathViewUser       string
		userOnAuthResponse authentication_response.UserAuthenticationResponse
		userOnAuth         authentication_response.UserContent
		fileName           = "UpdateUser.go"
		funcName           = "InternalUpdateUser"
		cfg                = config.ApplicationConfiguration.GetAuthenticationServer()
		path               = cfg.Host + cfg.PathRedirect.InternalUser.CrudUser
	)

	token = GenerateInternalToken("auth", inputStruct.AuthUserID, contextModel.AuthAccessTokenModel.ClientID, "update user", constanta.DefaultApplicationsLanguage)
	pathUpdateUser = path + "/" + strconv.Itoa(int(inputStruct.AuthUserID))
	pathViewUser = path + "/" + strconv.Itoa(int(inputStruct.AuthUserID))

	userOnAuthResponse, err = internalGetUserByID(token, pathViewUser, contextModel)
	if err.Error != nil {
		return
	}

	userOnAuth = userOnAuthResponse.Nexsoft.Payload.Data.Content
	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{token}

	inputToAuth := authentication_request.UserAuthenticationDTO{
		FirstName:        inputStruct.FirstName,
		LastName:         inputStruct.LastName,
		Email:            inputStruct.Email,
		CountryCode:      inputStruct.CountryCode,
		Phone:            inputStruct.Phone,
		Locale:           inputStruct.Locale,
		UpdatedAt:        userOnAuth.UpdatedAt,
		EmailMessage:     emailMessage,
		EmailLinkMessage: config.ApplicationConfiguration.GetNextracFrontend().Host + config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath,
		PhoneMessage:     constanta.PhoneMessageEmptyDefault,
	}

	if !inputStruct.IsAdmin {
		hostSysUser := config.ApplicationConfiguration.GetNextracFrontend().Host
		pathSysUser := strings.Replace(config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath, "/nexsoft-admin", "", 1)
		inputToAuth.EmailLinkMessage = hostSysUser + pathSysUser
	}

	statusCode, _, bodyResult, errs = common.HitAPI(pathUpdateUser, header, util.StructToJSON(inputToAuth), "PUT", *contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	if statusCode != 200 {
		err = common.ReadAuthServerError(funcName, statusCode, bodyResult, contextModel)
		return
	} else {
		_ = json.Unmarshal([]byte(bodyResult), &updateResposne)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
