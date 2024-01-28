package ProvinceService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input provinceService) ResetProvinceService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "ResetProvinceService"
		inputStruct in.ProvinceRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateResetProvince)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doResetProvince, func(interface{}, applicationModel.ContextModel) {
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

func (input provinceService) doResetProvince(tx *sql.Tx, inputStructInterface interface{}, _ *applicationModel.ContextModel, _ time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.ProvinceRequest)
		db          = serverconfig.ServerAttribute.DBConnection
		outputTemp  = make(map[string]interface{})
		lastSync    time.Time
	)

	//--- Get Last Sync
	lastSync, err = dao.ProvinceDAO.GetProvinceLastSync(db)
	if err.Error != nil {
		return
	}

	//--- Reset Last Sync
	err = dao.ProvinceDAO.UpdateLastSyncResetProvince(tx, repository.ProvinceModel{LastSync: sql.NullTime{Time: inputStruct.LastSync}})
	if err.Error != nil {
		return
	}

	timeLastSync := lastSync.Format(constanta.DefaultTimeFormat)
	outputTemp["last_sync"] = timeLastSync
	result = outputTemp
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input provinceService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ProvinceRequest) errorModel.ErrorModel) (inputStruct in.ProvinceRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}

func (input provinceService) validateResetProvince(inputStruct *in.ProvinceRequest) errorModel.ErrorModel {
	return inputStruct.ValidateReset()
}
