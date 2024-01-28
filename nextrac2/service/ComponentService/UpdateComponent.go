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
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input componentService) UpdateComponent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateComponent"
		inputStruct in.ComponentRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateComponent)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateComponent, func(interface{}, applicationModel.ContextModel) {
		//func Additional
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)

	return
}

func (input componentService) doUpdateComponent(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName      = "doUpdateComponent"
		inputStruct   = inputStructInterface.(in.ComponentRequest)
		inputModel    = input.convertDTOToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
		componentOnDB repository.ComponentModel
	)

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

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ComponentDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.ComponentDAO.UpdateComponent(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input componentService) convertDTOToModelUpdate(inputStruct in.ComponentRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.ComponentModel {
	return repository.ComponentModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		ComponentName: sql.NullString{String: inputStruct.ComponentName},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input componentService) validateUpdateComponent(inputStruct *in.ComponentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
