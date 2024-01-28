package ActivationLicenseService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type activationLicenseService struct {
	service.AbstractService
}

var ActivationLicenseService = activationLicenseService{}.New()

func (input activationLicenseService) New() (output activationLicenseService) {
	output.FileName = "ActivationLicenseService.go"
	return
}

func (input activationLicenseService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ActivationLicenseRequest) errorModel.ErrorModel) (inputStruct in.ActivationLicenseRequest, err errorModel.ErrorModel) {
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