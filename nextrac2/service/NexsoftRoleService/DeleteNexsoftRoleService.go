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
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type deleteNexsoftRoleService struct {
	service.AbstractService
	service.MultiDeleteData
}

var DeleteNexsoftRoleService = deleteNexsoftRoleService{}.New()

func (input deleteNexsoftRoleService) New() (output deleteNexsoftRoleService) {
	output.FileName = "DeleteNexsoftRoleService.go"
	return
}

func (input deleteNexsoftRoleService) DeleteNexsoftRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName        = "DeleteNexsoftRole"
		nexsoftRoleBody in.RoleRequest
	)

	nexsoftRoleBody, err = input.readBodyAndPathParam(request, contextModel)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, nexsoftRoleBody, contextModel, input.doDeleteNexsoftRole, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Function Additional
	})
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

func (input deleteNexsoftRoleService) doDeleteNexsoftRole(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName        = "doDeleteNexsoftRole"
		nexsoftRoleBody = inputStructInterface.(in.RoleRequest)
		nexsoftRoleOnDB repository.RoleModel
	)

	err = nexsoftRoleBody.ValidateDeleteRole()
	if err.Error != nil {
		return
	}

	nexsoftRoleModel := repository.RoleModel{
		ID:            sql.NullInt64{Int64: nexsoftRoleBody.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: 0},
	}

	nexsoftRoleOnDB, err = dao.NexsoftRoleDAO.GetNexsoftRoleForDelete(tx, nexsoftRoleModel)
	if err.Error != nil {
		return
	}

	if nexsoftRoleOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Role)
		return
	}

	err = input.checkNexsoftDeleteRoleLimitedByLimitedCreatedBy(contextModel, nexsoftRoleOnDB)
	if err.Error != nil {
		return
	}

	nexsoftRoleModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	if nexsoftRoleOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Role)
		return
	}

	if nexsoftRoleOnDB.UpdatedAt.Time != nexsoftRoleBody.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Role)
		return
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(10)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	nexsoftRoleModel.RoleID.String = nexsoftRoleOnDB.RoleID.String + encodedStr
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.NexsoftRoleDAO.TableName, nexsoftRoleOnDB.ID.Int64, 0)...)
	err = dao.NexsoftRoleDAO.DeleteNexsoftRole(tx, nexsoftRoleModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input deleteNexsoftRoleService) readBodyAndPathParam(request *http.Request, contextModel *applicationModel.ContextModel) (nexsoftRoleBody in.RoleRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndPathParam"
		stringBody string
		id         int64
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	readError := json.Unmarshal([]byte(stringBody), &nexsoftRoleBody)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, readError)
		return
	}

	id, err = readPathParam(request)
	if err.Error != nil {
		return
	}

	nexsoftRoleBody.ID = id
	return
}

func (input deleteNexsoftRoleService) checkNexsoftDeleteRoleLimitedByLimitedCreatedBy(contextModel *applicationModel.ContextModel, resultGetOnDB repository.RoleModel) (err errorModel.ErrorModel) {
	var (
		fileName = "NexsoftRoleService.go"
		funcName = "checkNexsoftDeleteRoleLimitedByLimitedCreatedBy"
	)

	// ---------- Check Created By Limited ----------
	if contextModel.LimitedByCreatedBy > 0 && (resultGetOnDB.CreatedBy.Int64 != contextModel.LimitedByCreatedBy) {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	return errorModel.GenerateNonErrorModel()
}
