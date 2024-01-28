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
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input roleService) InsertRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertRole"
		inputStruct in.RoleRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertRole, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleService) doInsertRole(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName              = "doInsertRole"
		inputStruct           = inputStructInterface.(in.RoleRequest)
		defaultPermissionTemp []string
		countPermission       int
		roleID                int64
	)

	defaultPermissionTemp = []string{
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
		"profile.profile-setting:view-own",
		"profile.profile-setting:update-own",
		"audit.audit-monitoring:view",
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
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Permission)
		return
	}

	permission := service.GenerateHashMapPermissionAndDataScope(inputStruct.Permission, false, false)
	service.ValidateRole(permission, input.FileName, funcName)

	roleModel := repository.RoleModel{
		RoleID:        sql.NullString{String: inputStruct.RoleID},
		Description:   sql.NullString{String: inputStruct.Description},
		Permission:    sql.NullString{String: util.StructToJSON(permission)},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	roleID, err = dao.RoleDAO.InsertRole(tx, roleModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.RoleDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: roleID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleService) validateInsert(inputStruct *in.RoleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertRole()
}
