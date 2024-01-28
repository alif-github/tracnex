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
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/service"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input userService) DeleteUserSysAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteUserSysAdmin"
		inputStruct in.UserRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteUserSysAdmin, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_DELETE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doDeleteUserSysAdmin(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.UserRequest)
		userModel   repository.UserModel
		userOnDB    repository.UserModel
		listToken   []string
	)

	userModel = repository.UserModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: inputStruct.UpdatedAt},
		CreatedBy:     sql.NullInt64{Int64: 0},
	}

	dataAudit, userOnDB, err = input.prepareDeleteUser(tx, userModel, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	userModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.UserDAO.TableName, userModel.ID.Int64, 0)...)
	err = dao.UserDAO.DeleteUser(tx, userModel, timeNow)
	if err.Error != nil {
		return
	}

	err = resource_common_service.InternalDeleteClientByClientID(userOnDB.ClientID.String, contextModel)
	if err.Error != nil {
		return
	}

	listToken, err = dao.ClientTokenDAO.GetListTokenByClientID(tx, userOnDB.ClientID.String)
	if err.Error != nil {
		return
	}

	go service.DeleteTokenFromRedis(listToken)
	err = dao.ClientTokenDAO.DeleteListTokenByClientID(tx, userOnDB.ClientID.String)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) prepareDeleteUser(tx *sql.Tx, userModel repository.UserModel, contextModel *applicationModel.ContextModel,
	timeNow time.Time) (dataAudit []repository.AuditSystemModel, userOnDB repository.UserModel, err errorModel.ErrorModel) {

	var (
		fileName = "DeleteUserSysAdmin.go"
		funcName = "prepareDeleteUser"
	)

	userOnDB, err = dao.UserDAO.GetUserForUpdate(tx, userModel)
	if err.Error != nil {
		return
	}

	if userOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	err = input.checkUserLimitedByLimitedCreatedBy(contextModel, userOnDB)
	if err.Error != nil {
		return
	}

	if userOnDB.UpdatedAt.Time != userModel.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.User)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.UserDAO.TableName, userOnDB.ID.Int64, userModel.CreatedBy.Int64)...)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateDelete(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDeleteUser()
}
