package HostServerService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input hostServerService) InsertHostServer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertHostServer"
	var inputStruct in.HostServerRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertHostServer, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_HOST_SERVER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerService) doInsertHostServer(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.HostServerRequest)
	var hostServerID int64

	hostServer := repository.HostServerModel{
		HostName:  sql.NullString{String: inputStruct.HostName},
		HostURL:   sql.NullString{String: inputStruct.HostUrl},
		CreatedBy: sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt: sql.NullTime{Time: timeNow},
		UpdatedBy: sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt: sql.NullTime{Time: timeNow},
	}

	hostServerID, err = dao.HostServerDAO.InsertHostnameAndIP(tx, hostServer)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.HostServerDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: hostServerID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hostServerService) ValidateInsert(inputStruct *in.HostServerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertHostServer()
}
