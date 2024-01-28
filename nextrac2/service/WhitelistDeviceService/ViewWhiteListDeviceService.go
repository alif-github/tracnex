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
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input whitelistDeviceService) ViewWhiteListDevice(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.WhiteListDeviceRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewWhiteListDevice)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewWhiteListDevice(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input whitelistDeviceService) validateViewWhiteListDevice(inputStruct *in.WhiteListDeviceRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}

func (input whitelistDeviceService) doViewWhiteListDevice(inputStruct in.WhiteListDeviceRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName = "doViewWhiteListDevice"
	)

	whiteListDeviceOnDB, err := dao.WhiteListDevice.ViewWhiteListDevice(serverconfig.ServerAttribute.DBConnection, repository.WhiteListDeviceModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})
	if err.Error != nil {
		return
	}

	if whiteListDeviceOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Device)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, whiteListDeviceOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(whiteListDeviceOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input whitelistDeviceService) convertModelToResponseDetail(inputModel repository.WhiteListDeviceModel) out.WhiteListDeviceResponse {
	return out.WhiteListDeviceResponse{
		ID:          inputModel.ID.Int64,
		Device:      inputModel.Device.String,
		Description: inputModel.Description.String,
		CreatedAt:   inputModel.CreatedAt.Time,
		UpdatedAt:   inputModel.UpdatedAt.Time,
		UpdatedName: inputModel.UpdatedName.String,
	}
}