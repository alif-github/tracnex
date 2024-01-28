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
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input whitelistDeviceService) InsertWhiteListDevice(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertWhiteListDevice"
		inputStruct in.WhiteListDeviceRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertWhiteListDevice, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input whitelistDeviceService) validateInsert(inputStruct *in.WhiteListDeviceRequest) errorModel.ErrorModel {
	return inputStruct.ValidateForInsert()
}

func (input whitelistDeviceService) doInsertWhiteListDevice(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct          = inputStructInterface.(in.WhiteListDeviceRequest)
		whiteListDeviceModel = input.convertStructToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
	)

	insertedID, err := dao.WhiteListDevice.InsertWhiteListDevice(tx, whiteListDeviceModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.WhiteListDevice.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	var dataAuditTemp repository.AuditSystemModel
	dataAuditTemp, err = input.GenerateDataScope(tx, insertedID, dao.WhiteListDevice.TableName, constanta.WhiteListDeviceDataScope, contextModel.AuthAccessTokenModel.ResourceUserID, contextModel.AuthAccessTokenModel.ClientID, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, dataAuditTemp)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input whitelistDeviceService) convertStructToModel(inputStruct in.WhiteListDeviceRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.WhiteListDeviceModel {
	return repository.WhiteListDeviceModel{
		Device:        sql.NullString{String: inputStruct.Device},
		Description:   sql.NullString{String: inputStruct.Description},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input whitelistDeviceService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_device_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Device)
		}
	}

	return err
}