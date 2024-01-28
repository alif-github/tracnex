package CRUDUserService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
)

type userService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var UserService = userService{}.New()

func (input userService) New() (output userService) {
	output.FileName = "UserService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"full_name",
		"role_id",
		"group_id",
		"email",
		"phone",
		"username",
		"created_at",
		"created_name",
	}
	output.ValidSearchBy = []string{
		"nt_username",
		"full_name",
		"phone",
		"email",
	}

	return
}

type userStruct struct {
	inputStruct in.UserRequest
}

func (input userService) readBodyAndValidateInsert(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UserRequest) errorModel.ErrorModel) (inputStruct in.UserRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, "readBodyAndValidate", errS)
		return
	}

	//---------- Default language is Indonesia
	inputStruct.Locale = constanta.DefaultApplicationsLanguage

	err = validation(&inputStruct)
	return
}

func (input userService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UserRequest) errorModel.ErrorModel) (inputStruct in.UserRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, "readBodyAndValidate", errS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)

	return
}

func (input userService) readBodyWithoutValidation(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.UserRequest, err errorModel.ErrorModel) {
	var stringBody string
	var funcName string = "readBodyWithoutValidation"

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	return
}

func (input userService) readBodyAndValidateForViewAndResendOTP(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UserRequest) errorModel.ErrorModel) (inputStruct in.UserRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)

	return
}

func (input userService) readBodyAndValidateForChangePasswordUser(request *http.Request, contextModel *applicationModel.ContextModel,
	validation func(input *in.ChangePasswordUserRequestDTO) errorModel.ErrorModel) (inputStruct in.ChangePasswordUserRequestDTO, err errorModel.ErrorModel) {

	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)

	return
}

func (input userService) checkRole(inputStruct in.UserRequest) (result repository.ClientRoleScopeModel, err errorModel.ErrorModel) {
	funcName := "checkRole"
	var roleModel repository.RoleModel

	if inputStruct.IsAdmin {
		roleModel, err = dao.NexsoftRoleDAO.GetNexsoftRoleByName(serverconfig.ServerAttribute.DBConnection, repository.RoleModel{RoleID: sql.NullString{String: inputStruct.Role}})
		if err.Error != nil {
			return
		}
	} else {
		roleModel, err = dao.RoleDAO.GetRoleByName(serverconfig.ServerAttribute.DBConnection, repository.RoleModel{RoleID: sql.NullString{String: inputStruct.Role}})
		if err.Error != nil {
			return
		}
	}

	if roleModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Role)
		return
	}

	result = repository.ClientRoleScopeModel{
		RoleID: roleModel.ID,
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) AddUserToAuthenticationServer(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel,
	isNexMileRegistration bool) (registerUserResponse authentication_response.RegisterUserAuthenticationResponse, err errorModel.ErrorModel) {
	var (
		funcName               = "AddUserToAuthenticationServer"
		authenticationServer   = config.ApplicationConfiguration.GetAuthenticationServer()
		registerAuthentication authentication_request.UserAuthenticationDTO
		internalToken          string
		registerUserUrl        string
		causedBy               error
	)

	registerAuthentication = authentication_request.UserAuthenticationDTO{
		Username:              inputStruct.Username,
		Password:              inputStruct.Password,
		FirstName:             inputStruct.FirstName,
		LastName:              inputStruct.LastName,
		Email:                 inputStruct.Email,
		CountryCode:           inputStruct.CountryCode,
		Phone:                 inputStruct.Phone,
		Device:                inputStruct.Device,
		Locale:                inputStruct.Locale,
		EmailMessage:          GetEmailMessage(inputStruct, isNexMileRegistration, false),
		PhoneMessage:          constanta.PhoneMessageEmptyDefault,
		IPWhitelist:           inputStruct.IPWhitelist,
		EmailLinkMessage:      config.ApplicationConfiguration.GetNextracFrontend().Host + config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath,
		ResourceID:            config.ApplicationConfiguration.GetServerResourceID(),
		AdditionalInformation: inputStruct.AdditionalInformation,
	}

	if !isNexMileRegistration {
		registerAuthentication.PhoneMessage = GetPhoneMessage(inputStruct)
	}

	if !inputStruct.IsAdmin {
		hostSysUser := config.ApplicationConfiguration.GetNextracFrontend().Host
		pathSysUser := strings.Replace(config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath, "/nexsoft-admin", "", 1)
		registerAuthentication.EmailLinkMessage = hostSysUser + pathSysUser
	}

	internalToken = resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, "", constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
	registerUserUrl = authenticationServer.Host + authenticationServer.PathRedirect.InternalUser.CrudUser
	statusCode, bodyResult, errorS := common.HitRegisterUserAuthenticationServer(internalToken, registerUserUrl, registerAuthentication, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &registerUserResponse)
	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy = errors.New(registerUserResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, registerUserResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) ChangePasswordToAuthenticationServer(inputStruct authentication_request.ChangePasswordDTOin, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName             = "changePasswordToAuthenticationServer"
		authenticationServer = config.ApplicationConfiguration.GetAuthenticationServer()
		data                 authentication_response.AuthenticationErrorResponse
		internalToken        string
		changePasswordUrl    string
	)

	internalToken = resource_common_service.GenerateInternalToken("auth", 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.DefaultApplicationsLanguage)
	changePasswordUrl = authenticationServer.Host + authenticationServer.PathRedirect.InternalUser.ChangePassword
	statusCode, bodyResult, errorS := common.HitChangePasswordAuthenticationServer(internalToken, changePasswordUrl, inputStruct, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &data)
	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(data.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, data.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_user_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientID)
		} else if service.CheckDBError(err, "uq_user_authuserid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.AuthUserID)
		}
	}
	return err
}

func GetEmailMessage(inputStruct in.UserRequest, _ bool, isUpdate bool) string {
	param := make(map[string]interface{})

	param["FIRST_NAME_USER"] = inputStruct.FirstName
	param["RESOURCE_NAME"] = "AUTHENTICATION"
	param["USERNAME"] = inputStruct.Username
	param["EMAIL"] = inputStruct.Email
	param["PHONE"] = constanta.IndonesianCodeNumber + "-" + inputStruct.Phone
	param["ACTIVATION_LINK"] = "{{.ACTIVATION_LINK}}" + "&username=" + inputStruct.Username
	if isUpdate {
		return util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, "VERIFY_UPDATE_EMAIL_MESSAGE", inputStruct.Locale, param)
	}
	return util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, "VERIFY_EMAIL_MESSAGE1", inputStruct.Locale, param)
}

func GetPhoneMessage(inputStruct in.UserRequest) string {
	var (
		param         = make(map[string]interface{})
		messageBundle = serverconfig.ServerAttribute.CommonServiceBundle
	)

	param["USER_ID"] = "{{.USER_ID}}"
	param["OTP_CODE"] = "{{.OTP_CODE}}"
	return util2.GenerateI18NServiceMessage(messageBundle, "VERIFY_PHONE_MESSAGE", inputStruct.Locale, param)
}

func (input userService) checkUserLimitedByLimitedCreatedBy(contextModel *applicationModel.ContextModel, resultGetOnDB repository.UserModel) (err errorModel.ErrorModel) {
	fileName := "UserService.go"
	funcName := "checkUserLimitedByLimitedCreatedBy"

	// ---------- Check Created By Limited ----------
	if contextModel.LimitedByCreatedBy > 0 && (resultGetOnDB.CreatedBy.Int64 != contextModel.LimitedByCreatedBy) {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}
	// -----------------------------------------------

	return errorModel.GenerateNonErrorModel()
}

func (input userService) checkUserDetailToAuthServer(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (userDetailAuth authentication_response.UserAuthenticationResponse, err errorModel.ErrorModel) {
	var (
		fileName                   = "UserService.go"
		funcName                   = "checkUserDetailToAuthServer"
		codeAuth                   = "E-4-AUT-SRV-003"
		authRequestByEmailAndPhone in.UserRequest
		userDetailAuthByPhone      authentication_response.UserAuthenticationResponse
		userDetailAuthByEmail      authentication_response.UserAuthenticationResponse
		authUserIDByPhone          int64
		authUserIDByEmail          int64
	)

	//--- Get Detail User (Verification 2 Field)
	authRequestByEmailAndPhone = input.setRequestForCheckSignatureUserAuth(inputStructInterface)
	userDetailAuth, err = HitAuthenticateServerForGetDetailUserAuth(authRequestByEmailAndPhone, contextModel)
	if err.Error != nil {
		if err.Error.Error() != codeAuth {
			return
		}
	}

	if userDetailAuth.Nexsoft.Payload.Data.Content.UserID < 1 {
		//--- Check Phone
		if !util.IsStringEmpty(authRequestByEmailAndPhone.Phone) {
			userDetailAuthByPhone, err = HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Phone: authRequestByEmailAndPhone.Phone}, contextModel)
			authUserIDByPhone = userDetailAuthByPhone.Nexsoft.Payload.Data.Content.UserID
			if err.Error != nil {
				if err.Error.Error() != codeAuth {
					return
				}
			}
		}

		//--- Check Email
		if !util.IsStringEmpty(authRequestByEmailAndPhone.Email) {
			userDetailAuthByEmail, err = HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Email: authRequestByEmailAndPhone.Email}, contextModel)
			authUserIDByEmail = userDetailAuthByEmail.Nexsoft.Payload.Data.Content.UserID
			if err.Error != nil {
				if err.Error.Error() != codeAuth {
					return
				}
			}
		}

		//--- Validate User
		if authUserIDByPhone > 0 && authUserIDByEmail > 0 {
			if authUserIDByEmail != authUserIDByPhone {
				err = errorModel.GenerateDifferentAuthUserId(fileName, funcName)
				return
			}
			userDetailAuth = userDetailAuthByEmail
		} else {
			if authUserIDByPhone > 0 {
				userDetailAuth = userDetailAuthByPhone
			} else {
				userDetailAuth = userDetailAuthByEmail
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
