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

func (input clientTypeService) DeleteClientType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteClientType"
		inputStruct in.ClientTypeRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteClientType)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteClientType, func(interface{}, applicationModel.ContextModel) {
		// additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_DELETE_CLIENT_TYPE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input clientTypeService) doDeleteClientType(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName         = "doDeleteClientType"
		inputStruct      = inputStructInterface.(in.ClientTypeRequest)
		isClientTypeUsed bool
	)

	inputModel := repository.ClientTypeModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	clientTypeOnDB, err := dao.ClientTypeDAO.GetClientTypeByID(serverconfig.ServerAttribute.DBConnection, repository.ClientTypeModel{
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

	if isClientTypeUsed {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.ClientType)
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

	err = dao.ClientTypeDAO.DeleteClientType(tx, inputModel)
	if err.Error != nil {
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ClientTypeDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	return
}

func (input clientTypeService) validateDeleteClientType(inputStruct *in.ClientTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
