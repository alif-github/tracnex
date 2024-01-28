package CRUDUserService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input userService) ViewProfileSettingUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.UserRequest
		stringBody  string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)
	inputStruct.ID = contextModel.AuthAccessTokenModel.ResourceUserID
	output.Data.Content, err = input.doViewUser(request, inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_VIEW_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) ViewUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.UserRequest
	inputStruct, err = input.readBodyAndValidateForViewAndResendOTP(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewUser(request, inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_VIEW_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doViewUser(request *http.Request, inputStruct in.UserRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		fileName     = "ViewUser.go"
		funcName     = "doViewUser"
		userModel    = repository.ViewDetailUserModel{ID: sql.NullInt64{Int64: inputStruct.ID}}
		splitRequest = strings.Split(request.RequestURI, "/")
		isAdmin      bool
		isUrlAdmin   bool
	)

	if splitRequest[3] == constanta.AdminName {
		isUrlAdmin = true
		userModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	}

	//------ Own permission
	if contextModel.LimitedByCreatedBy > 0 {
		userModel.ClientID.String = contextModel.AuthAccessTokenModel.ClientID
	}

	isAdmin, err = input.isIDForAdmin(userModel, isUrlAdmin)
	if err.Error != nil {
		return
	}

	if contextModel.AuthAccessTokenModel.ClientID == constanta.AdminClientID && !isUrlAdmin {
		isAdmin = false
	}

	userModel, err = dao.UserDAO.ViewDetailUser(serverconfig.ServerAttribute.DBConnection, userModel, isAdmin, isUrlAdmin)
	if err.Error != nil {
		return
	}

	if userModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	//--- Validate User In Oauth
	if err = input.validateUserInOauth(&userModel, contextModel); err.Error != nil {
		return
	}

	result, err = reformatDAOtoDTO(userModel, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateUserInOauth(userModel *repository.ViewDetailUserModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		detail     authentication_response.UserAuthenticationResponse
		phoneSplit []string
	)

	phoneSplit = strings.Split(userModel.Phone.String, "-")
	fmt.Println("Phone Split -> ", phoneSplit)

	modelOauth := in.UserRequest{Email: userModel.Email.String}
	if len(phoneSplit) == 2 {
		modelOauth.CountryCode = phoneSplit[0]
		modelOauth.Phone = phoneSplit[1]
	}

	detail, err = input.checkUserDetailToAuthServer(modelOauth, contextModel)
	if err.Error != nil {
		return
	}

	content := detail.Nexsoft.Payload.Data.Content
	if content.Phone != "" {
		userModel.IsVerifyPhone.Bool = true
	}

	if content.Email != "" {
		userModel.IsVerifyEmail.Bool = true
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func reformatDAOtoDTO(userModel repository.ViewDetailUserModel, contextModel *applicationModel.ContextModel) (out.ViewUserDTOOut, errorModel.ErrorModel) {
	var (
		fileName = "ViewUser.go"
		funcName = "reformatDAOtoDTO"
		errS     = errors.New("kesalahan status tidak ditemukan")
		err      errorModel.ErrorModel
		s        out.ViewUserDTOOut
	)

	s = out.ViewUserDTOOut{
		Username:       userModel.Username.String,
		Firstname:      userModel.FirstName.String,
		Lastname:       userModel.LastName.String,
		Email:          userModel.Email.String,
		Phone:          userModel.Phone.String,
		Role:           userModel.Role.String,
		GroupID:        userModel.GroupID.String,
		IsAdmin:        userModel.IsAdmin.Bool,
		Status:         userModel.Status.String,
		CreatedBy:      userModel.CreatedBy.Int64,
		CreatedName:    userModel.Username.String,
		CreatedAt:      userModel.CreatedAt.Time,
		UpdatedBy:      userModel.UpdatedBy.Int64,
		UpdatedName:    userModel.Username.String,
		UpdatedAt:      userModel.UpdatedAt.Time,
		PlatformDevice: userModel.PlatformDevice.String,
		IsVerifyPhone:  userModel.IsVerifyPhone.Bool,
		IsVerifyEmail:  userModel.IsVerifyEmail.Bool,
	}

	s.CreatedName = strings.TrimRight(s.CreatedName, " ")
	if util.IsStringEmpty(s.CreatedName) {
		if s.CreatedBy == constanta.SystemID {
			s.CreatedName = constanta.SystemClient
		} else {
			s.CreatedName = "-"
		}
	}

	s.UpdatedName = strings.TrimRight(s.UpdatedName, " ")
	if util.IsStringEmpty(s.UpdatedName) {
		if s.UpdatedBy == constanta.SystemID {
			s.UpdatedName = constanta.SystemClient
		} else {
			s.UpdatedName = "-"
		}
	}

	bundle := serverconfig.ServerAttribute.ConstantaBundle
	switch userModel.Status.String {
	case constanta.StatusActive:
		s.StatusDefine = util2.GenerateI18NServiceMessage(bundle, constanta.StatusActiveString, contextModel.AuthAccessTokenModel.Locale, nil)
	case constanta.StatusNonActive:
		s.StatusDefine = util2.GenerateI18NServiceMessage(bundle, constanta.StatusNonActiveString, contextModel.AuthAccessTokenModel.Locale, nil)
	case constanta.PendingOnApproval:
		s.StatusDefine = util2.GenerateI18NServiceMessage(bundle, constanta.StatusNotVerifiedYet, contextModel.AuthAccessTokenModel.Locale, nil)
		s.IsDisableStatus = true
	default:
		err = errorModel.GenerateUnknownError(fileName, funcName, errS)
	}

	return s, err
}

func (input userService) isIDForAdmin(userModel repository.ViewDetailUserModel, isUrlAdmin bool) (isAdmin bool, err errorModel.ErrorModel) {
	var (
		fileName = "ViewUser.go"
		funcName = "isIDForAdmin"
	)

	userModel, err = dao.UserDAO.IsAdmin(serverconfig.ServerAttribute.DBConnection, userModel, isUrlAdmin)
	if err.Error != nil {
		return
	}

	if userModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	isAdmin = userModel.IsAdmin.Bool
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateView(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewUserAndResendOTP()
}
