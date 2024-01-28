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
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type insertNexsoftRoleService struct {
	service.AbstractService
}

var InsertNexsoftRoleService = insertNexsoftRoleService{}.New()

func (input insertNexsoftRoleService) New() (output insertNexsoftRoleService) {
	output.FileName = "InsertNexsoftRoleService.go"
	return
}

func (input insertNexsoftRoleService) InsertNexsoftRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertNexsoftRole"
		inputStruct in.RoleRequest
		stringBody  string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = inputStruct.ValidateInsertRole()
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertNexsoftRole, func(_ interface{}, _ applicationModel.ContextModel) {})
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

func (input insertNexsoftRoleService) doInsertNexsoftRole(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName          = "doInsertNexsoftRole"
		inputStruct       = inputStructInterface.(in.RoleRequest)
		roleRepository    repository.RoleModel
		defaultPermission []string
		countPermission   int
		insertedID        int64
	)

	defaultPermission = getDefaultPermission(inputStruct)
	if defaultPermission != nil {
		inputStruct.Permission = append(inputStruct.Permission, defaultPermission...)
	}

	countPermission, err = dao.PermissionDAO.CheckIsPermissionValid(serverconfig.ServerAttribute.DBConnection, inputStruct.Permission)
	if err.Error != nil {
		return
	}

	if countPermission != len(inputStruct.Permission) {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Permission)
		return
	}

	if roleRepository, err = input.convertDtoToRepository(inputStruct, contextModel.AuthAccessTokenModel, timeNow); err.Error != nil {
		return
	}

	insertedID, err = dao.NexsoftRoleDAO.InsertNesoftRole(tx, roleRepository)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.NexsoftRoleDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input insertNexsoftRoleService) convertDtoToRepository(inputStruct in.RoleRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) (result repository.RoleModel, err errorModel.ErrorModel) {
	var (
		funcName   = "convertDtoToRepository"
		permission = service.GenerateHashMapPermissionAndDataScope(inputStruct.Permission, false, false)
	)

	service.ValidateRole(permission, input.FileName, funcName)
	result = repository.RoleModel{
		RoleID:        sql.NullString{String: inputStruct.RoleID},
		Description:   sql.NullString{String: inputStruct.Description},
		Permission:    sql.NullString{String: util.StructToJSON(permission)},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}

	return
}
