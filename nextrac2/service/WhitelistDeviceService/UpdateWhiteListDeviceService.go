package WhitelistDeviceService

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
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input whitelistDeviceService) UpdateWhiteListDevice(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateWhiteListDevice"
		inputStruct in.WhiteListDeviceRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateWhiteListDevice, func(interface{}, applicationModel.ContextModel) {
		// additional Function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input whitelistDeviceService) validateUpdate(inputStruct *in.WhiteListDeviceRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}

func (input whitelistDeviceService) doUpdateWhiteListDevice(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName             = "doUpdateWhiteListDevice"
		inputStruct          = inputStructInterface.(in.WhiteListDeviceRequest)
		whiteListDeviceModel = input.convertStructToModelForUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
	)

	whiteListDeviceOnDB, err := dao.WhiteListDevice.GetWhiteListDeviceForUpdateOrDelete(tx, whiteListDeviceModel, nil, nil)
	if err.Error != nil {
		return
	}

	if whiteListDeviceOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, whiteListDeviceOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if whiteListDeviceOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Device)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.WhiteListDevice.TableName, whiteListDeviceModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.WhiteListDevice.UpdateWhiteListDevice(tx, whiteListDeviceModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	return
}

func (input whitelistDeviceService) convertStructToModelForUpdate(inputStruct in.WhiteListDeviceRequest, authAccessModel model2.AuthAccessTokenModel, timeNow time.Time) repository.WhiteListDeviceModel {
	return repository.WhiteListDeviceModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Device:        sql.NullString{String: inputStruct.Device},
		Description:   sql.NullString{String: inputStruct.Description},
		UpdatedBy:     sql.NullInt64{Int64: authAccessModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessModel.ClientID},
	}
}