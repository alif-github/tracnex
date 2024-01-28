package ActivationUserNexmileService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type activationUserNexmileService struct {
	service.AbstractService
}

var ActivationUserNexmileService = activationUserNexmileService{}.New()

func (input activationUserNexmileService) New() (output activationUserNexmileService) {
	output.FileName = "ActivationUserNexmileService.go"
	return
}

func (input activationUserNexmileService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ActivationUserNexmileRequest) errorModel.ErrorModel) (inputStruct in.ActivationUserNexmileRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = validation(&inputStruct)
	return
}

func (input activationUserNexmileService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_pkceclientmapping_authuserid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FileName, constanta.ClientMappingAuthUserID)
		} else if service.CheckDBError(err, "uq_user_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientMappingClientID)
		} else if service.CheckDBError(err, "uq_user_authuserid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientMappingAuthUserID)
		} else if service.CheckDBError(err, "uq_clientrolescope_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientMappingClientID)
		}
	}

	return err
}