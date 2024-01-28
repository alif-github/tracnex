package ParameterService

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
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input parameterService) DeleteParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ParameterRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("DeleteParameter", inputStruct, contextModel, input.doDeleteParameter, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_DELETE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterService) deleteParameterOnDB(tx *sql.Tx, Parameter repository.ParameterModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "deleteParameterOnDB"
	var ParameterOnDB repository.ParameterModel

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		Parameter.CreatedBy.Int64 = userID
	}

	ParameterOnDB, err = dao.ParameterDAO.GetParameterForUpdate(tx, Parameter)
	if err.Error != nil {
		return
	}

	if ParameterOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	if ParameterOnDB.UpdatedAt.Time != Parameter.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Parameter)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ParameterDAO.TableName, Parameter.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.ParameterDAO.DeleteParameter(tx, Parameter, timeNow)
	if err.Error != nil {
		return
	}

	return
}

func (input parameterService) doDeleteParameter(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {

	inputStruct := inputStructInterface.(in.ParameterRequest)

	Parameter := repository.ParameterModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Permission:    sql.NullString{String: inputStruct.Permission},
		Name:          sql.NullString{String: inputStruct.Name},
		Value:         sql.NullString{String: inputStruct.Value},
		Description:   sql.NullString{String: inputStruct.Description},
		UpdatedAt:     sql.NullTime{Time: inputStruct.UpdatedAt},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	dataAudit, err = input.deleteParameterOnDB(tx, Parameter, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterService) validateDelete(inputStruct *in.ParameterRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDeleteParameter()
}
