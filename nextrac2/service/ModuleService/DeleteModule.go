package ModuleService

import (
	"database/sql"
	"github.com/google/uuid"
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

func (input moduleService) DeleteModule(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteModule"
		inputStruct in.ModuleRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteModule)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteModule, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input moduleService) doDeleteModule(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.ModuleRequest)
		moduleOnDB  repository.ModuleModel
	)

	// Created Input Model
	inputModel := repository.ModuleModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	// Validate Module Check DB
	moduleOnDB, err = input.validateModuleOnDB(tx, inputStruct, inputModel, contextModel)
	if err.Error != nil {
		return
	}

	// Create Random Module DB
	input.randTokenGenerator(&inputModel, moduleOnDB)

	// Delete Module
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ModuleDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.ModuleDAO.DeleteModule(tx, inputModel)
	return
}

func (input moduleService) validateModuleOnDB(tx *sql.Tx, inputStruct in.ModuleRequest, inputModel repository.ModuleModel, contextModel *applicationModel.ContextModel) (moduleOnDB repository.ModuleModel, err errorModel.ErrorModel) {
	funcName := "validateModuleOnDB"
	moduleOnDB, err = input.ModuleDAO.GetModuleForUpdate(tx, repository.ModuleModel{ID: inputModel.ID})
	if err.Error != nil {
		return
	}

	if moduleOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Module)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, moduleOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if moduleOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Module)
		return
	}

	if moduleOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Module)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input moduleService) randTokenGenerator(inputModel *repository.ModuleModel, moduleOnDB repository.ModuleModel) {
	encodedStr := service.RandTimeToken(constanta.RandTokenForDeleteLength, uuid.New().String())
	inputModel.ModuleName.String = moduleOnDB.ModuleName.String + encodedStr
}

func (input moduleService) validateDeleteModule(inputStruct *in.ModuleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
