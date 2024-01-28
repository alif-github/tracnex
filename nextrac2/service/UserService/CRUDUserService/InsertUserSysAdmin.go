package CRUDUserService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input userService) InsertUserSysAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertUserSysAdmin"
		inputStruct in.UserRequest
	)

	inputStruct, err = input.readBodyAndValidateInsert(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, userStruct{inputStruct: inputStruct}, contextModel, input.doMappingUser, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_INSERT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doInsertUserSysAdmin(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		registerUserResponse authentication_response.RegisterUserAuthenticationResponse
		roleScope            repository.ClientRoleScopeModel
		dataGroupOnDB        repository.DataGroupModel
		temp                 = inputStructInterface.(userStruct)
		inputStruct          = temp.inputStruct
		funcName             = "doInsertUserSysAdmin"
	)

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

	//--- Registration user to authentication
	registerUserResponse, err = input.AddUserToAuthenticationServer(inputStruct, contextModel, false)
	if err.Error != nil {
		return
	}

	clientRoleScope := repository.ClientRoleScopeModel{
		RoleID:  roleScope.RoleID,
		GroupID: dataGroupOnDB.ID,
	}

	dataAudit, err = input.saveUserSysAdminToDB(tx, inputStruct, registerUserResponse, clientRoleScope, contextModel, timeNow, "")
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) saveUserSysAdminToDB(tx *sql.Tx, inputStruct in.UserRequest, registerUserResponse authentication_response.RegisterUserAuthenticationResponse,
	ClientRoleScope repository.ClientRoleScopeModel, contextModel *applicationModel.ContextModel, timeNow time.Time, code string) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		content   authentication_response.RegisterUserContent
		userModel repository.UserModel
		id        int64
	)

	content = registerUserResponse.Nexsoft.Payload.Data.Content
	userModel = repository.UserModel{
		FirstName:      sql.NullString{String: inputStruct.FirstName},
		LastName:       sql.NullString{String: inputStruct.LastName},
		Username:       sql.NullString{String: inputStruct.Username},
		Email:          sql.NullString{String: inputStruct.Email},
		Phone:          sql.NullString{String: constanta.IndonesianCodeNumber + "-" + inputStruct.Phone},
		ClientID:       sql.NullString{String: content.ClientID},
		AuthUserID:     sql.NullInt64{Int64: content.UserID},
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
		PlatformDevice: sql.NullString{String: inputStruct.PlatformDevice},
		Currency:       sql.NullString{String: inputStruct.Currency},
	}

	if code == constanta.Remark {
		userModel.Status = sql.NullString{String: constanta.StatusActive}
	}

	id, err = dao.UserDAO.InsertUser(tx, userModel)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: id},
	})

	ClientRoleScope.ClientID.String = content.ClientID
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

func (input userService) validateInsert(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertUser()
}
