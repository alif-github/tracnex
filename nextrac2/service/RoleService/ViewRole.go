package RoleService

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

func (input roleService) ViewRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.RoleRequest

	inputStruct, err = input.readParamAndValidate(request, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewRole(inputStruct, contextModel)
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

func (input roleService) doViewRole(inputStruct in.RoleRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	fileName := "ViewRole.go"
	funcName := "doViewRole"

	roleModel := repository.RoleModel{ID: sql.NullInt64{Int64: inputStruct.ID}}
	roleModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	roleModel, err = dao.RoleDAO.ViewRole(serverconfig.ServerAttribute.DBConnection, roleModel)
	if err.Error != nil {
		return
	}

	if roleModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Role)
		return
	}

	result = reformatDAOtoDTO(roleModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func reformatDAOtoDTO(roleModel repository.RoleModel) out.ViewDetailRoleDTOOut {
	var rolesPermission map[string][]string
	_ = json.Unmarshal([]byte(roleModel.Permission.String), &rolesPermission)
	permissionProductLicense1, ok1 := rolesPermission[constanta.MenuUserMasterProdukLisensi]
	if ok1 {
		newValue := remove(permissionProductLicense1, "insert")
		rolesPermission[constanta.MenuUserMasterProdukLisensi] = newValue
	}

	permissionProductLicense2, ok2 := rolesPermission[constanta.MenuUserMasterProdukLisensiRedesign]
	if ok2 {
		newValue := remove(permissionProductLicense2, "insert")
		rolesPermission[constanta.MenuUserMasterProdukLisensiRedesign] = newValue
	}

	permission := service.GenerateInitiateRoleDTOOut(rolesPermission, rolesPermission, true)
	return out.ViewDetailRoleDTOOut{
		ID:          roleModel.ID.Int64,
		RoleID:      roleModel.RoleID.String,
		Description: roleModel.Description.String,
		Permissions: permission.Permissions,
		CreatedBy:   roleModel.CreatedBy.Int64,
		CreatedAt:   roleModel.CreatedAt.Time,
		UpdatedBy:   roleModel.UpdatedBy.Int64,
		UpdatedAt:   roleModel.UpdatedAt.Time,
	}
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (input roleService) validateView(inputStruct *in.RoleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewRole()
}
