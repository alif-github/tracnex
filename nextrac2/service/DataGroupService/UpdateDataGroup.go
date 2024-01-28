package DataGroupService

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
	"time"
)

func (input dataGroupService) UpdateDataGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateDataGroup"
		inputStruct in.DataGroupRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateDataGroup, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) doUpdateDataGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName    = "doUpdateDataGroup"
		inputStruct = inputStructInterface.(in.DataGroupRequest)
		dataGroupDB repository.DataGroupModel
		countScope  int
		listToken   []string
	)

	countScope, err = dao.DataScopeDAO.CheckIsScopeValid(serverconfig.ServerAttribute.DBConnection, inputStruct.Scope)
	if countScope != len(inputStruct.Scope) {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.DataScope)
		return
	}

	scope := service.GenerateHashMapPermissionAndDataScope(inputStruct.Scope, true, true)
	scopeString := util.StructToJSON(scope)
	dataGroup := repository.DataGroupModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Description:   sql.NullString{String: inputStruct.Description},
		Scope:         sql.NullString{String: scopeString},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	dataGroup.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	dataGroupDB, err = dao.DataGroupDAO.GetDataGroupForDelete(tx, dataGroup)
	if err.Error != nil {
		return
	}

	if dataGroupDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.DataGroup)
		return
	}

	if dataGroupDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.DataGroup)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.DataGroupDAO.TableName, dataGroup.ID.Int64, 0)...)
	err = dao.DataGroupDAO.UpdateDataGroup(tx, dataGroup)
	if err.Error != nil {
		return
	}

	listToken, err = dao.ClientTokenDAO.GetListTokenByGroupID(tx, dataGroupDB.ID.Int64)
	if err.Error != nil {
		return
	}

	go service.DeleteTokenFromRedis(listToken)
	Login.LogoutAuthServerAutomatic(listToken, *contextModel)

	err = dao.ClientTokenDAO.DeleteListTokenByGroupID(tx, dataGroupDB.ID.Int64)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) validateUpdate(inputStruct *in.DataGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateDataGroup()
}
