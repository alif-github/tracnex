package HostServerService

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
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input hostServerService) DeleteHostServer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DeleteHostServer"
	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.ValidateDelete)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteHostServer, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_DELETE_HOST_SERVER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input hostServerService) doDeleteHostServer(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.HostServerRequest)

	hostServer := repository.HostServerModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: inputStruct.UpdatedAt},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	dataAudit, err = input.deleteHostServerOnDB(tx, hostServer, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	return
}

func (input hostServerService) deleteHostServerOnDB(tx *sql.Tx, hostServer repository.HostServerModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "deleteHostServerOnDB"

	err = input.checkHostServerUsed(hostServer)
	if err.Error != nil {
		return
	}

	hostServer.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	hostServerOnDB, err := dao.HostServerDAO.GetHostServerForUpdate(tx, hostServer)
	if err.Error != nil {
		return
	}

	if hostServerOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	if hostServerOnDB.UpdatedAt.Time != hostServer.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.HostServer)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.HostServerDAO.TableName, hostServer.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.HostServerDAO.DeleteHostServer(tx, hostServer, timeNow)
	if err.Error != nil {
		return

	}

	return
}

func (input hostServerService) ValidateDelete(inputStruct *in.HostServerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateHostServer()
}

func (input hostServerService) checkHostServerUsed(hostServerModel repository.HostServerModel) (err errorModel.ErrorModel) {
	funcName := "checkHostServerUsed"
	isUsed, err := dao.HostServerDAO.CheckIsHostServerUsed(serverconfig.ServerAttribute.DBConnection, hostServerModel)
	if err.Error != nil {
		return
	}

	if isUsed {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.HostServer)
		return
	}
	return
}
