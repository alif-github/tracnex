package RoleService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	Login "nexsoft.co.id/nextrac2/service/session/Logout"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input roleService) UpdateRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateRole"
		inputStruct in.RoleRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateRole, func(_ interface{}, _ applicationModel.ContextModel) {})
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

func (input roleService) doUpdateRole(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName        = "UpdateRole.go"
		funcName        = "doUpdateRole"
		inputStruct     = inputStructInterface.(in.RoleRequest)
		countPermission int
		roleDB          repository.RoleModel
	)

	defaultPermissionTemp := []string{
		"admin.pengguna.sub-pengguna:update-own",
		"admin.pengguna.sub-pengguna:changepassword-own",
		"admin.pengguna.sub-pengguna:view-own",
		"admin.pengguna:update-own",
		"admin.pengguna:changepassword-own",
		"admin.pengguna:view-own",
		"admin.kredit:view",
		"admin:view-own",
		"admin:update-own",
		"admin:changepassword-own",
		"nexsoft.province:view",
		"nexsoft.district:view",
		"nexsoft.customer-group:view",
		"nexsoft.customer-category:view",
		"nexsoft.product-group:view",
		"nexsoft.client-type:view",
		"nexsoft:view",
		"nexsoft:update",
		"nexsoft:delete",
		"nexsoft:insert",
	}

	for i := 0; i < len(inputStruct.Permission); i++ {
		for j := 0; j < len(defaultPermissionTemp); j++ {
			if inputStruct.Permission[i] == defaultPermissionTemp[j] {
				inputStruct.Permission = input.removeIndexPerm(inputStruct.Permission, i)
				i = -1
				break
			}
		}
	}

	for _, valueDefaultPerm := range defaultPermissionTemp {
		inputStruct.Permission = append(inputStruct.Permission, valueDefaultPerm)
	}

	countPermission, err = dao.PermissionDAO.CheckIsPermissionValid(serverconfig.ServerAttribute.DBConnection, inputStruct.Permission)
	if err.Error != nil {
		return
	}

	if countPermission != len(inputStruct.Permission) {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Permission)
		return
	}

	permission := service.GenerateHashMapPermissionAndDataScope(inputStruct.Permission, false, false)
	service.ValidateRole(permission, input.FileName, funcName)
	roleModel := repository.RoleModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		RoleID:        sql.NullString{String: inputStruct.RoleID},
		Description:   sql.NullString{String: inputStruct.Description},
		Permission:    sql.NullString{String: util.StructToJSON(permission)},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	roleModel.CreatedBy.Int64 = 0
	roleDB, err = dao.RoleDAO.GetRoleForUpdate(tx, roleModel)
	if err.Error != nil {
		return
	}

	if roleDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Role)
		return
	}

	err = input.checkRoleLimitedByLimitedCreatedBy(contextModel, roleDB)
	if err.Error != nil {
		return
	}

	roleModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	if roleDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.Role)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.RoleDAO.TableName, roleDB.ID.Int64, 0)...)
	err = dao.RoleDAO.UpdateRole(tx, roleModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	if roleDB.Permission.String != util.StructToJSON(permission) || roleDB.RoleID.String != inputStruct.RoleID {
		var listToken []string
		listToken, err = dao.ClientTokenDAO.GetListTokenByRoleID(tx, roleModel.ID.Int64)
		if err.Error != nil {
			return
		}

		go service.DeleteTokenFromRedis(listToken)
		Login.LogoutAuthServerAutomatic(listToken, *contextModel)
		err = dao.ClientTokenDAO.DeleteListTokenByRoleID(tx, roleModel.ID.Int64)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleService) validateUpdate(inputStruct *in.RoleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateRole()
}
