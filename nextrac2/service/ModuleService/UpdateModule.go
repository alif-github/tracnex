package ModuleService

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

func (input moduleService) UpdateModule(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UpdateModule"
	var inputStruct in.ModuleRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateModule)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateModule, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input moduleService) doUpdateModule(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdateModule"
	inputStruct := inputStructInterface.(in.ModuleRequest)
	inputModel := input.convertDTOToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	moduleOnDB, err := dao.ModuleDAO.GetModuleForUpdate(tx, repository.ModuleModel{
		ID: inputModel.ID,
	})

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

	if moduleOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Module)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ModuleDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.ModuleDAO.UpdateModule(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	return
}

func (input moduleService) convertDTOToModelUpdate(inputStruct in.ModuleRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.ModuleModel {
	return repository.ModuleModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		ModuleName:    sql.NullString{String: inputStruct.ModuleName},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input moduleService) validateUpdateModule(inputStruct *in.ModuleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
