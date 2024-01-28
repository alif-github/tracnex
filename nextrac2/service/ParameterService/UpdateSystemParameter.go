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

func (input parameterService) UpdateParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ParameterRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("UpdateParameter", inputStruct, contextModel, input.doUpdateParameter, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterService) doUpdateParameter(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdateParameter"

	inputStruct := inputStructInterface.(in.ParameterRequest)
	var ParameterDB repository.ParameterModel

	ParameterModel := repository.ParameterModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Permission:    sql.NullString{String: inputStruct.Permission},
		Name:          sql.NullString{String: inputStruct.Name},
		Value:         sql.NullString{String: inputStruct.Value},
		Description:   sql.NullString{String: inputStruct.Description},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		ParameterModel.CreatedBy.Int64 = userID
	}

	ParameterDB, err = dao.ParameterDAO.GetParameterForUpdate(tx, ParameterModel)
	if err.Error != nil {
		return
	}

	if ParameterDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Parameter)
		return
	}

	if ParameterDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Parameter)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ParameterDAO.TableName, inputStruct.ID, contextModel.LimitedByCreatedBy)...)

	err = dao.ParameterDAO.UpdateParameter(tx, ParameterModel, timeNow)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterService) validateUpdate(inputStruct *in.ParameterRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateParameter()
}
