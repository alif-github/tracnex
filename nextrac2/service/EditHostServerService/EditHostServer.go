package EditHostServerService

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
	"strconv"
	"time"
)

func (input editHostServerService) EditHostServer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "EditHostServer"
	var inputStruct in.CronHostRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateEditHostServer)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doEditHostServer, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_HOST_SERVER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input editHostServerService) doEditHostServer(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doEditHostServer"
	inputStruct := inputStructInterface.(in.CronHostRequest)
	var cronHostID int64
	var serverRunID int64

	inputStruct.HostID = inputStruct.ID
	hostServer := repository.HostServerModel{
		ID:        sql.NullInt64{Int64: inputStruct.HostID},
		HostName:  sql.NullString{String: inputStruct.HostName},
		UpdatedBy: sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt: sql.NullTime{Time: inputStruct.UpdatedAt},
	}

	var hostServerOnDB repository.HostServerModel
	hostServerOnDB, err = dao.HostServerDAO.GetHostServer(tx, repository.HostServerModel{ID: sql.NullInt64{Int64: inputStruct.HostID}})
	if err.Error != nil {
		return
	}

	if hostServerOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.HostID)
		return
	}

	if hostServerOnDB.UpdatedAt.Time != hostServer.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.HostServer)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.HostServerDAO.TableName, hostServerOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.HostServerDAO.UpdateHostnameAndIP(tx, hostServer, timeNow)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	for i := 0; i < len(inputStruct.ListCron); i++ {
		var cron repository.CRONSchedulerModel
		cron, err = dao.CronSchedulerDAO.GetCronScheduler(tx, repository.CRONSchedulerModel{ID: sql.NullInt64{Int64: inputStruct.ListCron[i].CronID}})
		if err.Error != nil {
			return
		}

		if cron.ID.Int64 == 0 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CronID)
			err.AdditionalInformation = append(err.AdditionalInformation, strconv.Itoa(int(inputStruct.ListCron[i].CronID)))
			return
		}

		var server repository.ServerRunModel
		server, err = dao.ServerRunDAO.CheckIsServerRunExist(serverconfig.ServerAttribute.DBConnection, repository.ServerRunModel{
			HostID:  sql.NullInt64{Int64: inputStruct.HostID},
			RunType: sql.NullString{String: cron.RunType.String},
		})
		if err.Error != nil {
			return
		}

		serverRun := repository.ServerRunModel{
			ID:            sql.NullInt64{Int64: server.ID.Int64},
			Name:          sql.NullString{String: cron.Name.String},
			RunType:       sql.NullString{String: cron.RunType.String},
			HostID:        sql.NullInt64{Int64: inputStruct.HostID},
			Status:        sql.NullBool{Bool: inputStruct.ListCron[i].RunStatus},
			CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedAt:     sql.NullTime{Time: timeNow},
			CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		}

		if server.ID.Int64 == 0 {
			serverRunID, err = dao.ServerRunDAO.InsertServerRun(tx, serverRun)
			if err.Error != nil {
				return
			}

			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.ServerRunDAO.TableName, serverRunID, contextModel.LimitedByCreatedBy)...)
		} else {
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ServerRunDAO.TableName, server.ID.Int64, contextModel.LimitedByCreatedBy)...)
			err = dao.ServerRunDAO.UpdateServerRun(tx, serverRun, timeNow)
			if err.Error != nil {
				return

			}
		}

		var checkedCron repository.CRONHostModel
		checkedCron, err = dao.CronHostDAO.CheckIsCronHostExist(serverconfig.ServerAttribute.DBConnection, repository.CRONHostModel{
			CronID: sql.NullInt64{Int64: inputStruct.ListCron[i].CronID},
			HostID: sql.NullInt64{Int64: inputStruct.HostID},
		})

		if err.Error != nil {
			return
		}

		cronHost := repository.CRONHostModel{
			ID:            sql.NullInt64{Int64: checkedCron.ID.Int64},
			CronID:        sql.NullInt64{Int64: inputStruct.ListCron[i].CronID},
			HostID:        sql.NullInt64{Int64: inputStruct.HostID},
			Status:        sql.NullBool{Bool: inputStruct.ListCron[i].Active},
			CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedAt:     sql.NullTime{Time: timeNow},
			CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		}

		if checkedCron.ID.Int64 == 0 {
			cronHostID, err = dao.CronHostDAO.InsertCronHost(tx, cronHost)
			if err.Error != nil {
				return
			}
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.ServerRunDAO.TableName, cronHostID, contextModel.LimitedByCreatedBy)...)

		} else {
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CronHostDAO.TableName, checkedCron.ID.Int64, contextModel.LimitedByCreatedBy)...)
			err = dao.CronHostDAO.UpdateCronHost(tx, cronHost, timeNow)
			if err.Error != nil {
				return

			}
		}
	}
	return
}

func (input editHostServerService) ValidateEditHostServer(inputStruct *in.CronHostRequest) errorModel.ErrorModel {
	return inputStruct.ValidateEdit()

}
