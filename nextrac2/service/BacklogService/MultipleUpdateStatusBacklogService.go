package BacklogService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"time"
)

func (input backlogService) MultipleUpdateStatusBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "MultipleUpdateStatusBacklog"
		inputStruct in.MultipleUpdateStatusRequest
	)

	inputStruct, err = input.readBodyAndValidateForMultipleUpdateStatusBacklog(request, contextModel, input.validateMultipleUpdateStatus)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doMultipleUpdateStatus, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input backlogService) validateMultipleUpdateStatus(inputStruct *in.MultipleUpdateStatusRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateUpdate()
}

func (input backlogService) doMultipleUpdateStatus(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.MultipleUpdateStatusRequest)
		listModel   = input.convertToDTOToInputModelForMultipleUpdate(inputStruct, *contextModel, timeNow)
		backlogOnDB repository.BacklogModel
		funcName    = "doMultipleUpdateStatus"
	)

	// cek updated_at on db
	for _, id := range inputStruct.ID {
		backlogOnDB, err = dao.BacklogDAO.ViewDetailBacklog(serverconfig.ServerAttribute.DBConnection, repository.BacklogModel{
			ID: sql.NullInt64{Int64: id}},
		)
		if err.Error != nil {
			return
		}

		if backlogOnDB.ID.Int64 < 1 {
			dataId := fmt.Sprintf(`ID %d`, id)
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, dataId)
			return
		}

		// matiin updated_at
		//inputItem.UpdatedAt, err = in.TimeStrToTime(inputItem.UpdatedAtStr, constanta.UpdatedAt)
		//if backlogOnDB.UpdatedAt.Time.Unix() != inputItem.UpdatedAt.Unix() {
		//	dataId := fmt.Sprintf(`Diubah pada id %d`, backlogOnDB.ID.Int64)
		//	err = errorModel.GenerateDataLockedError(fileName, funcName, dataId)
		//	return
		//}
	}

	//--- Update On DB
	dataAudit, err = input.updateMultipleStatus(tx, contextModel, timeNow, listModel)
	return
}

func (input backlogService) convertToDTOToInputModelForMultipleUpdate(inputStruct in.MultipleUpdateStatusRequest, contextModel applicationModel.ContextModel, timeNow time.Time) (listModel []repository.BacklogModel) {
	for _, id := range inputStruct.ID {
		temp := repository.BacklogModel{
			ID:            sql.NullInt64{Int64: id},
			Status:        sql.NullString{String: inputStruct.Status},
			UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		}

		listModel = append(listModel, temp)
	}

	return
}

func (input backlogService) updateMultipleStatus(tx *sql.Tx, contextModel *applicationModel.ContextModel, timeNow time.Time, listModel []repository.BacklogModel) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	// dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.BacklogDAO.TableName, model.ID.Int64, contextModel.LimitedByCreatedBy)...)
	for _, backlog := range listModel {
		err = dao.BacklogDAO.UpdateStatusBacklog(tx, backlog)
		if err.Error != nil {
			// err = input.checkDuplicateError(err)
			return
		}
	}

	return
}
