package ComponentService

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
	"time"
)

func (input componentService) DeleteComponent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteComponent"
		inputStruct in.ComponentRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteComponent)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteComponent, func(interface{}, applicationModel.ContextModel) {
		//func optional
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input componentService) doDeleteComponent(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName      = "doDeleteComponent"
		inputStruct   = inputStructInterface.(in.ComponentRequest)
		componentOnDB repository.ComponentModel
		inputModel    repository.ComponentModel
	)

	inputModel = repository.ComponentModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	componentOnDB, err = dao.ComponentDAO.GetComponentForUpdate(tx, repository.ComponentModel{ID: inputModel.ID})
	if err.Error != nil {
		return
	}

	if componentOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Component)
		return
	}

	if componentOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Component)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, componentOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if componentOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Component)
		return
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(8)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	inputModel.ComponentName.String = componentOnDB.ComponentName.String + encodedStr

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ComponentDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.ComponentDAO.DeleteComponent(tx, inputModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input componentService) validateDeleteComponent(inputStruct *in.ComponentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
