package NexsoftRoleService

import (
	"database/sql"
	"encoding/json"
	"net/http"
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
)

type viewNexsoftRoleService struct {
	FileName string
	//service.AbstractService
}

var ViewNexsoftRoleService = viewNexsoftRoleService{FileName: "ViewNexsoftRoleService.go"}

func (input viewNexsoftRoleService) ViewNexsoftRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	nexsoftRoleBody, err := input.readPathParamAndValidate(request)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewNexsoftRole(nexsoftRoleBody, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input viewNexsoftRoleService) doViewNexsoftRole(nexsoftRoleBody in.RoleRequest, contextModel *applicationModel.ContextModel) (result out.ViewDetailRoleDTOOut, err errorModel.ErrorModel) {
	funcName := "doViewNexsoftRole"

	nexsoftRoleModel := repository.RoleModel{
		ID:        sql.NullInt64{Int64: nexsoftRoleBody.ID},
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}

	nexsoftRoleOnDB, err := dao.NexsoftRoleDAO.GetDetailNexsoftRole(serverconfig.ServerAttribute.DBConnection, nexsoftRoleModel)
	if err.Error != nil {
		return
	}

	if nexsoftRoleOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Role)
		return
	}
	
	result = input.convertDAOToDTO(nexsoftRoleOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input viewNexsoftRoleService) convertDAOToDTO(nexsoftRoleModel repository.RoleModel) (output out.ViewDetailRoleDTOOut) {
	var permissionMap map[string][]string

	_ = json.Unmarshal([]byte(nexsoftRoleModel.Permission.String), &permissionMap)
	permision := service.GenerateInitiateRoleDTOOut(permissionMap, permissionMap, true)
	return out.ViewDetailRoleDTOOut{
		ID:          nexsoftRoleModel.ID.Int64,
		RoleID:      nexsoftRoleModel.RoleID.String,
		Description: nexsoftRoleModel.Description.String,
		Permissions: permision.Permissions,
		CreatedBy:   nexsoftRoleModel.CreatedBy.Int64,
		CreatedAt: 	 nexsoftRoleModel.CreatedAt.Time,
		UpdatedBy: 	 nexsoftRoleModel.UpdatedBy.Int64,
		UpdatedAt:   nexsoftRoleModel.UpdatedAt.Time,
	}
}

func (input viewNexsoftRoleService) readPathParamAndValidate(request *http.Request) (nexsoftRoleBody in.RoleRequest, err errorModel.ErrorModel) {
	id, err := readPathParam(request)

	if err.Error != nil {
		return
	}

	nexsoftRoleBody.ID = id
	err = nexsoftRoleBody.ValidateViewRole()
	return
}