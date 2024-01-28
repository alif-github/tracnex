package ClientTypeService

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

func (input clientTypeService) UpdateClientType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateClientType"
		inputStruct in.ClientTypeRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateClientType)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateClientType, func(interface{}, applicationModel.ContextModel) {
		// additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_CLIENT_TYPE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input clientTypeService) doUpdateClientType(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName         = "doUpdateClientType"
		inputStruct      = inputStructInterface.(in.ClientTypeRequest)
		inputModel       = input.convertDTOToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
		clientTypeOnDB   repository.ClientTypeModel
		isClientTypeUsed bool
	)

	clientTypeOnDB, err = dao.ClientTypeDAO.GetClientTypeByID(serverconfig.ServerAttribute.DBConnection, repository.ClientTypeModel{
		ID: inputModel.ID,
	})
	if err.Error != nil {
		return
	}

	if clientTypeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientType)
		return
	}

	isClientTypeUsed, err = dao.ClientTypeDAO.IsClientTypeUsed(serverconfig.ServerAttribute.DBConnection, repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})

	if err.Error != nil {
		return
	}

	if inputStruct.ClientType != clientTypeOnDB.ClientType.String && isClientTypeUsed {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, "client_type")
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, clientTypeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if clientTypeOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.ClientType)
		return
	}

	err = dao.ClientTypeDAO.UpdateClientType(tx, inputModel)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ClientTypeDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	return
}

func (input clientTypeService) validateUpdateClientType(inputStruct *in.ClientTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
