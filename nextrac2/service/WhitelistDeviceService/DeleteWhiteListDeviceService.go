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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input whitelistDeviceService) DeleteWhiteListDevice(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteProductGroup"
		inputStruct in.WhiteListDeviceRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteWhiteListDevice, func(interface{}, applicationModel.ContextModel) {
		// additional Function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input whitelistDeviceService) validateDelete(inputStruct *in.WhiteListDeviceRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}

func (input whitelistDeviceService) doDeleteWhiteListDevice(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName    = "DeleteWhiteListDevice.go"
		funcName    = "doDeleteWhiteListDevice"
		inputStruct = inputStructInterface.(in.WhiteListDeviceRequest)
	)

	whiteListDeviceModel := repository.WhiteListDeviceModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	whiteListDeviceOnDB, err := dao.WhiteListDevice.GetWhiteListDeviceForUpdateOrDelete(tx, whiteListDeviceModel, nil, nil)
	if err.Error != nil {
		return
	}

	if whiteListDeviceOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, whiteListDeviceOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if whiteListDeviceOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.Device)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.WhiteListDevice.TableName, whiteListDeviceModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.WhiteListDevice.HardDeleteWhiteListDevice(tx, whiteListDeviceModel)
	if err.Error != nil {
		return
	}

	return
}