package Login

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
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/session"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type groChatLoginService struct {
	service.AbstractService
}

var GroChatLoginService = groChatLoginService{}.New()

func (input groChatLoginService) New() (output groChatLoginService) {
	output.FileName = "GroChatLoginService.go"
	return
}

func (input groChatLoginService) Login(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	body, errModel := input.readBodyAndValidate(request, contextModel, input.validateBody)
	if errModel.Error != nil {
		return
	}

	response, errModel := input.login(body, contextModel)
	if errModel.Error != nil {
		return
	}

	output.Data.Content = response
	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("TOKEN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input groChatLoginService) login(body in.GroChatLoginDTOIn, contextModel *applicationModel.ContextModel) (response out.GroChatLoginResponse, errModel errorModel.ErrorModel) {
	/*
		Request Internal Token
	*/
	internalToken := resource_common_service.GenerateInternalToken("chat", 0, "", config.ApplicationConfiguration.GetServerResourceID(), "id-ID")

	/*
		Request Login
	*/
	groChatResponse, errModel := input.requestLogin(contextModel, internalToken, body)
	if errModel.Error != nil {
		return
	}

	/*
		Request User Detail
	*/
	userDetail, errModel := input.requestUserDetail(contextModel, groChatResponse.UserToken)
	if errModel.Error != nil {
		return
	}

	authUser := userDetail.GroChat.Auth

	/*
		Create new user
	*/
	if errModel = input.createInvitedUser(userDetail); errModel.Error != nil {
		return
	}

	/*
		Get user by client id
	*/
	user, errModel := dao.UserDAO.GetUserByClientID(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
		ClientID: sql.NullString{String: authUser.ClientId},
	})
	if errModel.Error != nil {
		return
	}

	/*
		Validate User
	*/
	if errModel = input.validateUser(user); errModel.Error != nil {
		return
	}

	response = out.GroChatLoginResponse{
		Token:        groChatResponse.UserToken,
		RefreshToken: groChatResponse.RefreshToken,
	}
	return
}

func (input groChatLoginService) validateUser(user repository.UserModel) errorModel.ErrorModel {
	funcName := "validateUser"

	if !user.ID.Valid {
		return errorModel.GenerateUnauthorizedClientError(input.FileName, funcName)
	}

	if user.Status.String == constanta.StatusNonActive {
		return errorModel.GenerateUserStatusNonactive(input.FileName, funcName)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input groChatLoginService) requestLogin(contextModel *applicationModel.ContextModel, internalToken string, body in.GroChatLoginDTOIn) (groChatContent grochat_response.GroChatAuthenticationData, errModel errorModel.ErrorModel) {
	funcName := "requestLogin"

	groChatServer := config.ApplicationConfiguration.GetGrochat()

	path := fmt.Sprintf("%s%s", groChatServer.Host, groChatServer.PathRedirect.Authentication)

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{internalToken}
	headerRequest["Content-Type"] = []string{"application/json"}

	req := util.StructToJSON(grochat_request.GroChatLoginRequest{
		RequestId:         body.RequestId,
		AuthorizationCode: body.AuthorizationCode,
		CodeVerifier:      body.CodeVerifier,
		ResourceId:        config.ApplicationConfiguration.GetServerResourceID(),
	})

	statusCode, _, bodyResult, err := common.HitAPI(path, headerRequest, req, "POST", *contextModel)
	if err != nil {
		errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
		return
	}

	if statusCode != http.StatusOK {
		errModel = input.generateGroChatRequestError(bodyResult)
		return
	}

	groChatResponse, errModel := input.generateGroChatAuthenticationResponse(bodyResult)
	if errModel.Error != nil {
		return
	}

	return groChatResponse.Data, errorModel.GenerateNonErrorModel()
}

func (input groChatLoginService) requestUserDetail(contextModel *applicationModel.ContextModel, PKCEToken string) (result *grochat_response.GroChatUserDetailData, errModel errorModel.ErrorModel) {
	funcName := "requestUserDetail"

	groChatServer := config.ApplicationConfiguration.GetGrochat()

	path := fmt.Sprintf("%s%s", groChatServer.Host, groChatServer.PathRedirect.UserDetail)

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{PKCEToken}

	statusCode, _, bodyResult, err := common.HitAPI(path, headerRequest, "", "GET", *contextModel)
	if err != nil {
		errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
		return
	}

	if statusCode != http.StatusOK {
		errModel = input.generateGroChatRequestError(bodyResult)
		return
	}

	groChatResponse, errModel := input.generateGroChatUserDetailResponse(bodyResult)
	if errModel.Error != nil {
		return
	}

	return groChatResponse.Data, errorModel.GenerateNonErrorModel()
}

func (input groChatLoginService) generateGroChatAuthenticationResponse(bodyResult string) (groChatResponse grochat_response.GroChatAuthenticationResponse, errModel errorModel.ErrorModel) {
	funcName := "generateGroChatAuthenticationResponse"

	var (
		groChatInvalidResource grochat_response.GroChatAuthenticationErrorResponse
	)

	/*
		Invalid Resource
	*/
	errUnmarshalInvalidResource := json.Unmarshal([]byte(bodyResult), &groChatInvalidResource)

	if errUnmarshalInvalidResource == nil {
		if groChatInvalidResource.Status == 300003 {
			errModel = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
			return
		}
	}

	if err := json.Unmarshal([]byte(bodyResult), &groChatResponse); err != nil {
		errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input groChatLoginService) generateGroChatUserDetailResponse(bodyResult string) (groChatResponse grochat_response.GroChatUserDetailResponse, errModel errorModel.ErrorModel) {
	funcName := "generateGroChatUserDetailResponse"

	var (
		groChatErrResponse grochat_response.GroChatErrorResponse
	)

	errUnmarshalNote := json.Unmarshal([]byte(bodyResult), &groChatErrResponse)

	if errUnmarshalNote == nil {
		if groChatErrResponse.Status == 1000027 {
			errModel = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, groChatErrResponse.Code, "GROCHAT", errors.New(groChatErrResponse.Note))
			return
		}

		if groChatErrResponse.Status != 1 {
			errModel = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, groChatErrResponse.Code, "GROCHAT", errors.New(groChatErrResponse.Description))
			return
		}
	}

	if err := json.Unmarshal([]byte(bodyResult), &groChatResponse); err != nil {
		errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input groChatLoginService) generateGroChatRequestError(content string) errorModel.ErrorModel {
	funcName := "generateGroChatRequestError"

	var errorResult grochat_response.GroChatAuthenticationErrorResponse
	_ = json.Unmarshal([]byte(content), &errorResult)

	causedBy := errors.New(errorResult.Description)

	return errorModel.GenerateAuthenticationServerError("GroChatServiceUtil.go", funcName, errorResult.Code, "GROCHAT", causedBy)
}

func (input groChatLoginService) validateBody(body *in.GroChatLoginDTOIn) errorModel.ErrorModel {
	return body.ValidateGroChatLoginDTO()
}

func (input groChatLoginService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.GroChatLoginDTOIn) errorModel.ErrorModel) (body in.GroChatLoginDTOIn, errModel errorModel.ErrorModel) {
	funcName := "readBody"

	content, errModel := input.ReadBody(request, contextModel)
	if errModel.Error != nil {
		return
	}

	if err := json.Unmarshal([]byte(content), &body); err != nil {
		errModel = errorModel.GenerateInvalidRequestError(input.FileName, funcName, err)
		return
	}

	if errModel = validation(&body); errModel.Error != nil {
		return
	}

	return body, errorModel.GenerateNonErrorModel()
}

func (input groChatLoginService) createInvitedUser(groChatUser *grochat_response.GroChatUserDetailData) (errModel errorModel.ErrorModel) {
	funcName := "createInvitedUser"

	now := time.Now()
	authUser := groChatUser.AuthServer.NexSoft.Payload.Data.Content

	tx, err := serverconfig.ServerAttribute.DBConnection.Begin()
	if err != nil {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	defer func() {
		if errModel.Error != nil {
			_ = tx.Rollback()
			return
		}

		if err = tx.Commit(); err != nil {
			errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
			return
		}
	}()

	/*
		Get invitation
	*/
	invitation, errModel := dao.UserInvitationDAO.GetByEmailOrClientIdForUpdate(tx, repository.UserInvitation{
		Email:          sql.NullString{String: authUser.Email},
		ClientId:       sql.NullString{String: authUser.ClientId},
	})
	if errModel.Error != nil {
		return
	}

	if !invitation.Id.Valid {
		return errorModel.GenerateNonErrorModel()
	}

	/*
		Insert User
	*/
	userModel := repository.UserModel{
		ClientID:        sql.NullString{String: authUser.ClientId},
		AuthUserID:      sql.NullInt64{Int64: authUser.UserId},
		Locale:          sql.NullString{String: constanta.DefaultApplicationsLanguage},
		Status:          sql.NullString{String: constanta.StatusActive},
		FirstName:       sql.NullString{String: authUser.FirstName},
		LastName:        sql.NullString{String: authUser.LastName},
		Username:        sql.NullString{String: authUser.Username},
		Email:           sql.NullString{String: authUser.Email},
		Phone:           sql.NullString{String: authUser.Phone},
		CreatedAt:       sql.NullTime{Time: now},
		UpdatedAt:       sql.NullTime{Time: now},
		AliasName:       sql.NullString{String: fmt.Sprintf("%s %s", authUser.FirstName, authUser.LastName)},
		Currency:        sql.NullString{String: constanta.CurrencyIDR},
		PlatformDevice:  sql.NullString{String: constanta.PlatformWebsite},
	}

	_, errModel = dao.UserDAO.InsertUser(tx, userModel)
	if errModel.Error != nil {
		return
	}

	/*
		Insert Client Role Scope
	*/
	clientRoleScopeModel := repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: authUser.ClientId},
		RoleID:        sql.NullInt64{Int64: invitation.RoleId.Int64},
		GroupID:       sql.NullInt64{Int64: invitation.DataGroupId.Int64},
		CreatedAt:     sql.NullTime{Time: now},
		UpdatedAt:     sql.NullTime{Time: now},
	}

	_, errModel = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, clientRoleScopeModel)
	if errModel.Error != nil {
		return
	}

	/*
		Delete Invitation
	*/
	errModel = dao.UserInvitationDAO.DeleteByIdTx(tx, invitation.Id.Int64)
	return
}