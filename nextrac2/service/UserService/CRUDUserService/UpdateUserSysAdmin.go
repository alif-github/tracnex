package CRUDUserService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	Login "nexsoft.co.id/nextrac2/service/session/Logout"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

func (input userService) UpdateUserSysAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateUserSysAdmin"
		inputStruct in.UserRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdatedUserSysAdmin, func(_ interface{}, _ applicationModel.ContextModel) {
		//-- func additional
	})
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

func (input userService) doUpdatedUserSysAdmin(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName                  = "UpdateUserSystemAdmin.go"
		funcName                  = "doUpdatedUserSysAdmin"
		inputStruct               = inputStructInterface.(in.UserRequest)
		userModel                 repository.UserModel
		userOnDB                  repository.UserModel
		roleScope                 repository.ClientRoleScopeModel
		dataGroupOnDB             repository.DataGroupModel
		updateUserResponse        authentication_response.UpdateUserAuthenticationResponse
		updateUserContentResponse authentication_response.UpdateUserContent
		clientRoleScopeModel      repository.ClientRoleScopeModel
		listToken                 []string
		userModelStatus           string
	)

	inputStruct.Locale = constanta.IndonesianLanguage
	userModel = input.convertUserDTOToDAO(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
	userModelStatus = userModel.Status.String

	//--- Check user on DB
	userModel.CreatedBy.Int64 = 0
	userOnDB, err = dao.UserDAO.GetUserForUpdate(tx, userModel)
	if err.Error != nil {
		return
	}

	if userOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
		return
	}

	err = input.checkUserLimitedByLimitedCreatedBy(contextModel, userOnDB)
	if err.Error != nil {
		return
	}

	userModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	if userOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.User)
		return
	}

	//--- Check role scope
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

	//--- Check Is Change Level User
	if userModel.IsSystemAdmin.Bool != userOnDB.IsSystemAdmin.Bool {
		clientRoleScopeModel = repository.ClientRoleScopeModel{
			ClientID:      userOnDB.ClientID,
			UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:     sql.NullTime{Time: timeNow},
		}

		err = input.doDeleteRoleScope(tx, userModel.IsSystemAdmin.Bool, clientRoleScopeModel)
		if err.Error != nil {
			return
		}
	}

	detailUserAuth, errors := input.checkUserDetailToAuthServer(inputStruct, contextModel)
	fmt.Println("Detail User Auth : ", detailUserAuth.Nexsoft)
	fmt.Println("Payload Response : ", detailUserAuth.Nexsoft.Payload.PayloadResponse, " with status : ", detailUserAuth.Nexsoft.Payload.PayloadResponse.Status)
	fmt.Println("Payload Status : ", detailUserAuth.Nexsoft.Payload.Status, " with other : ", detailUserAuth.Nexsoft.Payload.PayloadResponse.Other)
	fmt.Println("Detail Content : ", detailUserAuth.Nexsoft.Payload.Data.Content)
	fmt.Println("Error From Auth : ", errors.Error)
	if errors.Error != nil || detailUserAuth.Nexsoft.Payload.Status.Code != "200" {
		//err = errorModel.GenerateUnknownUserAuth(fileName, funcName)
		//return
	} else {
		//--- Update to Authentication Server
		inputStruct.AuthUserID = userOnDB.AuthUserID.Int64
		updateUserResponse, err = input.updateUserToAuthenticationServer(inputStruct, contextModel)
		updateUserContentResponse = updateUserResponse.Nexsoft.Payload.Data.Content
		if err.Error != nil {
			return
		}
	}

	if userOnDB.Status.String == constanta.PendingOnApproval {
		userModel.Status = userOnDB.Status
	}

	if (updateUserContentResponse.EmailStatus.EmailNotifyStatus) && (updateUserContentResponse.EmailStatus.EmailNotify) &&
		(updateUserContentResponse.PhoneStatus.PhoneNotifyStatus) && (updateUserContentResponse.PhoneStatus.PhoneNotify) {
		userModel.Status.String = constanta.PendingOnApproval
	} else if ((updateUserContentResponse.EmailStatus.EmailNotifyStatus) && (updateUserContentResponse.EmailStatus.EmailNotify)) ||
		((updateUserContentResponse.PhoneStatus.PhoneNotifyStatus) && (updateUserContentResponse.PhoneStatus.PhoneNotify)) {
		userModel.Status.String = constanta.StatusActive
	}

	if userModelStatus == constanta.StatusNonActive {
		userModel.Status.String = constanta.StatusNonActive
	}

	//--- Update to DB
	clientRoleScopeModel = input.getClientRoleDAO(roleScope.RoleID.Int64, userOnDB.ClientID.String, contextModel.AuthAccessTokenModel, timeNow)
	clientRoleScopeModel.GroupID = dataGroupOnDB.ID
	dataAudit, err = input.updateUserToDB(tx, userModel, clientRoleScopeModel, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//--- Kick user
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

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) updateUserToDB(tx *sql.Tx, userModel repository.UserModel, roleScopeModel repository.ClientRoleScopeModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		clientRoleOdDB repository.ClientRoleScopeModel
		idInserted     int64
	)

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserDAO.TableName, userModel.ID.Int64, 0)...)
	err = dao.UserDAO.UpdateUserInAdmin(tx, userModel)
	if err.Error != nil {
		return
	}

	if userModel.IsSystemAdmin.Bool {
		clientRoleOdDB, err = dao.NexsoftClientRoleScopeDAO.IsNexsoftClientRoleScopeExist(tx, repository.ClientRoleScopeModel{
			ClientID: roleScopeModel.ClientID,
		})

		if err.Error != nil {
			return
		}

		if clientRoleOdDB.ID.Int64 < 1 {
			idInserted, err = dao.NexsoftClientRoleScopeDAO.InsertNexsoftClientRoleScope(tx, roleScopeModel)
			if err.Error != nil {
				return
			}

			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.NexsoftClientRoleScopeDAO.TableName, idInserted, 0)...)
		} else {
			roleScopeModel.ID = clientRoleOdDB.ID
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.NexsoftClientRoleScopeDAO.TableName, clientRoleOdDB.ID.Int64, 0)...)
			err = dao.NexsoftClientRoleScopeDAO.UpdateNexsoftClientRoleScope(tx, roleScopeModel)
			if err.Error != nil {
				return
			}
		}
	} else {
		clientRoleOdDB, err = dao.ClientRoleScopeDAO.IsClientRoleScopeExist(tx, repository.ClientRoleScopeModel{
			ClientID: roleScopeModel.ClientID,
		})
		if err.Error != nil {
			return
		}

		if clientRoleOdDB.ID.Int64 < 1 {
			idInserted, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, roleScopeModel)
			if err.Error != nil {
				return
			}

			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.ClientRoleScopeDAO.TableName, idInserted, 0)...)
		} else {
			roleScopeModel.ID = clientRoleOdDB.ID
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ClientRoleScopeDAO.TableName, clientRoleOdDB.ID.Int64, 0)...)
			err = dao.ClientRoleScopeDAO.UpdateClientRoleScope(tx, roleScopeModel)
			if err.Error != nil {
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doDeleteRoleScope(tx *sql.Tx, isUserAdmin bool, clientRoleScopeModel repository.ClientRoleScopeModel) (err errorModel.ErrorModel) {
	if isUserAdmin {
		err = dao.ClientRoleScopeDAO.DeleteClientRoleScopeByClientID(tx, clientRoleScopeModel)
		if err.Error != nil {
			return
		}
	} else {
		err = dao.NexsoftClientRoleScopeDAO.DeleteNexsoftClientRoleScopeByClientID(tx, clientRoleScopeModel)
		if err.Error != nil {
			return
		}
	}

	return input.doDeleteToken(tx, clientRoleScopeModel.ClientID.String)
}

func (input userService) doDeleteToken(tx *sql.Tx, clientID string) (err errorModel.ErrorModel) {
	var listToken []string

	listToken, err = dao.ClientTokenDAO.GetListTokenByClientID(tx, clientID)
	if err.Error != nil {
		return
	}

	go service.DeleteTokenFromRedis(listToken)

	err = dao.ClientTokenDAO.DeleteListTokenByClientID(tx, clientID)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) updateUserToAuthenticationServer(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel) (updateResposne authentication_response.UpdateUserAuthenticationResponse, err errorModel.ErrorModel) {
	emailMessages := GetEmailMessage(inputStruct, false, true)
	return resource_common_service.InternalUpdateUser(inputStruct, contextModel, emailMessages)
}

func (input userService) convertUserDTOToDAO(inputStruct in.UserRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.UserModel {
	return repository.UserModel{
		ID:             sql.NullInt64{Int64: inputStruct.ID},
		Locale:         sql.NullString{String: inputStruct.Locale},
		IsSystemAdmin:  sql.NullBool{Bool: inputStruct.IsAdmin},
		FirstName:      sql.NullString{String: inputStruct.FirstName},
		LastName:       sql.NullString{String: inputStruct.LastName},
		Email:          sql.NullString{String: inputStruct.Email},
		Phone:          sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.Phone},
		Status:         sql.NullString{String: inputStruct.Status},
		UpdatedBy:      sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:      sql.NullTime{Time: now},
		UpdatedClient:  sql.NullString{String: authAccessToken.ClientID},
		PlatformDevice: sql.NullString{String: inputStruct.PlatformDevice},
		Currency:       sql.NullString{String: inputStruct.Currency},
	}
}

func (input userService) getClientRoleDAO(roleId int64, clientId string, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.ClientRoleScopeModel {
	return repository.ClientRoleScopeModel{
		ClientID:      sql.NullString{String: clientId},
		RoleID:        sql.NullInt64{Int64: roleId},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
	}
}

func (input userService) validateUpdate(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateUser()
}

func (input userService) UpdateHelpingUserDeleted(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	fileName := "UpdateUserSysAdmin.go"
	funcName := "UpdateHelpingUserDeleted"
	var inputStruct in.UserRequest

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)
	if inputStruct.ID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, "helping")
		return
	}

	err = dao.UserDAO.UpdateHelpingTableUser(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}, true)
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

func (input userService) UpdateHelpingClientND6FirstName(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	fileName := "UpdateUserSysAdmin.go"
	funcName := "UpdateHelpingClientND6FirstName"
	var inputStruct in.UserRequest
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errS != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)
	if inputStruct.ID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, "helping")
		return
	}

	err = dao.UserDAO.UpdateHelpingTableUser(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
		FirstName: sql.NullString{String: inputStruct.FirstName},
	}, false)
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
