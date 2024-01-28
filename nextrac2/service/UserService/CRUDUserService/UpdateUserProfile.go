package CRUDUserService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	Login "nexsoft.co.id/nextrac2/service/session/Logout"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input userService) updateProfile(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "updateProfile"

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateProfile, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_UPDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doUpdateProfile(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName                  = "UpdateUserProfile.go"
		funcName                  = "doUpdateProfile"
		inputStruct               = inputStructInterface.(in.UserRequest)
		listToken                 []string
		isOnlyHaveOwnAccess       bool
		userModel                 repository.UserModel
		userOnDB                  repository.UserModel
		updateUserResponse        authentication_response.UpdateUserAuthenticationResponse
		updateUserContentResponse authentication_response.UpdateUserContent
	)

	inputStruct.Locale = constanta.DefaultApplicationsLanguage
	inputStruct.ID = contextModel.AuthAccessTokenModel.ResourceUserID
	userModel = input.getUserProfileModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	_, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		userModel.ClientID.String = contextModel.AuthAccessTokenModel.ClientID
	}

	userOnDB, err = dao.UserDAO.GetUserForUpdateProfile(tx, userModel)
	if err.Error != nil {
		return
	}

	if userOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
		return
	}

	if userOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.User)
		return
	}

	//------ Update to Authentication Server
	inputStruct.AuthUserID = userOnDB.AuthUserID.Int64
	inputStruct.Username = userOnDB.Username.String
	updateUserResponse, err = input.updateUserToAuthenticationServer(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	updateUserContentResponse = updateUserResponse.Nexsoft.Payload.Data.Content
	userModel.Status = userOnDB.Status
	userModel.IsSystemAdmin = userOnDB.IsSystemAdmin
	if updateUserContentResponse.EmailStatus.EmailNotifyStatus && updateUserContentResponse.EmailStatus.EmailNotify {
		userModel.Status.String = constanta.PendingOnApproval
	}

	//------ Update User
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserDAO.TableName, userModel.ID.Int64, 0)...)
	err = dao.UserDAO.UpdateUserInAdmin(tx, userModel)
	if err.Error != nil {
		return
	}

	if userOnDB.Email.String != userModel.Email.String || userOnDB.Phone.String != userModel.Phone.String {
		//------ Kick User
		listToken, err = dao.ClientTokenDAO.GetListTokenByClientID(tx, userOnDB.ClientID.String)
		if err.Error != nil {
			return
		}

		go service.DeleteTokenFromRedis(listToken)
		Login.LogoutAuthServerAutomatic(listToken, *contextModel)
		err = dao.ClientTokenDAO.DeleteListTokenByClientID(tx, userOnDB.ClientID.String)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) getUserProfileModel(inputStruct in.UserRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.UserModel {
	return repository.UserModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Locale:        sql.NullString{String: inputStruct.Locale},
		FirstName:     sql.NullString{String: inputStruct.FirstName},
		LastName:      sql.NullString{String: inputStruct.LastName},
		Email:         sql.NullString{String: inputStruct.Email},
		Phone:         sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.Phone},
		IsSystemAdmin: sql.NullBool{Bool: inputStruct.IsAdmin},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input userService) validateUpdateUserProfile(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateUserProfile()
}
