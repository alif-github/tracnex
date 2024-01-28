package RoleService

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
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input roleService) DeleteRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteRole"
		inputStruct in.RoleRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteRole, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_DELETE_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleService) doDeleteRole(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName    = "DeleteRole.go"
		funcName    = "doDeleteRole"
		inputStruct = inputStructInterface.(in.RoleRequest)
		roleDB      repository.RoleModel
	)

	roleModel := repository.RoleModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	roleModel.CreatedBy.Int64 = 0
	roleDB, err = dao.RoleDAO.GetRoleForDelete(tx, roleModel)
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
	if roleDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(fileName, funcName, constanta.Role)
		return
	}

	if roleDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.Role)
		return
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(10)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	roleModel.RoleID.String = roleDB.RoleID.String + encodedStr
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.RoleDAO.TableName, roleDB.ID.Int64, 0)...)
	err = dao.RoleDAO.DeleteRole(tx, roleModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleService) validateDelete(inputStruct *in.RoleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDeleteRole()
}
