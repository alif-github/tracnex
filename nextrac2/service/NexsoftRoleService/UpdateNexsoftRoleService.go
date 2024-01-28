package NexsoftRoleService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	Login "nexsoft.co.id/nextrac2/service/session/Logout"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type updateNexsoftRoleService struct {
	service.AbstractService
}

var UpdateNexsoftRoleService = updateNexsoftRoleService{}.New()

func (input updateNexsoftRoleService) New() (output updateNexsoftRoleService) {
	output.FileName = "UpdateNexsoftRoleService.go"
	return
}

func (input updateNexsoftRoleService) UpdateNexsoftRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateNexsoftRole"
		inputStruct in.RoleRequest
	)

	inputStruct, err = input.readBodyAndPathParam(request, contextModel)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateNexsoftRole, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input updateNexsoftRoleService) doUpdateNexsoftRole(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName          = "doUpdateNexsoftRole"
		nexsoftRoleBody   = inputStructInterface.(in.RoleRequest)
		defaultPermission = getDefaultPermission(nexsoftRoleBody)
		nexsoftRoleModel  repository.RoleModel
		nexsoftRoleOnDB   repository.RoleModel
		countPermission   int
	)

	if defaultPermission != nil {
		nexsoftRoleBody.Permission = append(nexsoftRoleBody.Permission, defaultPermission...)
	}

	nexsoftRoleModel = input.convertDTOToModel(nexsoftRoleBody, contextModel.AuthAccessTokenModel, timeNow)
	nexsoftRoleModel.CreatedBy.Int64 = 0

	nexsoftRoleOnDB, err = dao.NexsoftRoleDAO.GetNexsoftRoleForUpdate(tx, nexsoftRoleModel)
	if err.Error != nil {
		return
	}

	if nexsoftRoleOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Role)
		return
	}

	err = input.checkNexsoftUpdateRoleLimitedByLimitedCreatedBy(contextModel, nexsoftRoleOnDB)
	if err.Error != nil {
		return
	}

	nexsoftRoleModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	countPermission, err = dao.PermissionDAO.CheckIsPermissionValid(serverconfig.ServerAttribute.DBConnection, nexsoftRoleBody.Permission)
	if countPermission != len(nexsoftRoleBody.Permission) {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Permission)
		return
	}

	permission := service.GenerateHashMapPermissionAndDataScope(nexsoftRoleBody.Permission, false, false)
	service.ValidateRole(permission, input.FileName, funcName)
	nexsoftRoleModel.Permission = sql.NullString{String: util.StructToJSON(permission)}

	if err = input.validaton(nexsoftRoleBody, nexsoftRoleOnDB); err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.NexsoftRoleDAO.TableName, nexsoftRoleOnDB.ID.Int64, 0)...)
	err = dao.NexsoftRoleDAO.UpdateNexoftRole(tx, nexsoftRoleModel)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	if nexsoftRoleOnDB.Permission.String != util.StructToJSON(permission) || nexsoftRoleOnDB.RoleID.String != nexsoftRoleBody.RoleID {
		var listToken []string
		listToken, err = dao.ClientTokenDAO.GetListTokenByRoleID(tx, nexsoftRoleModel.ID.Int64)
		if err.Error != nil {
			return
		}

		go service.DeleteTokenFromRedis(listToken)
		Login.LogoutAuthServerAutomatic(listToken, *contextModel)
		err = dao.ClientTokenDAO.DeleteListTokenByRoleID(tx, nexsoftRoleModel.ID.Int64)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input updateNexsoftRoleService) validaton(nexsoftRoleBody in.RoleRequest, nexsoftRoleOnDB repository.RoleModel) (err errorModel.ErrorModel) {
	funcName := "validaton"
	err = nexsoftRoleBody.ValidateUpdateRole()
	if err.Error != nil {
		return
	}

	if nexsoftRoleBody.UpdatedAt != nexsoftRoleOnDB.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Role)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input updateNexsoftRoleService) convertDTOToModel(requestBody in.RoleRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) (result repository.RoleModel) {
	return repository.RoleModel{
		ID:            sql.NullInt64{Int64: requestBody.ID},
		RoleID:        sql.NullString{String: requestBody.RoleID},
		Description:   sql.NullString{String: requestBody.Description},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input updateNexsoftRoleService) readBodyAndPathParam(request *http.Request, contexModel *applicationModel.ContextModel) (inputStruct in.RoleRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndPathParam"
	stringBody, err := input.ReadBody(request, contexModel)
	if err.Error != nil {
		return
	}

	readError := json.Unmarshal([]byte(stringBody), &inputStruct)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, readError)
		return
	}

	id, err := readPathParam(request)
	if err.Error != nil {
		return
	}

	inputStruct.ID = id
	return
}

func (input updateNexsoftRoleService) checkNexsoftUpdateRoleLimitedByLimitedCreatedBy(contextModel *applicationModel.ContextModel, resultGetOnDB repository.RoleModel) (err errorModel.ErrorModel) {
	fileName := "NexsoftRoleService.go"
	funcName := "checkNexsoftUpdateRoleLimitedByLimitedCreatedBy"

	// ---------- Check Created By Limited ----------
	if contextModel.LimitedByCreatedBy > 0 && (resultGetOnDB.CreatedBy.Int64 != contextModel.LimitedByCreatedBy) {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}
	// -----------------------------------------------

	return errorModel.GenerateNonErrorModel()
}
