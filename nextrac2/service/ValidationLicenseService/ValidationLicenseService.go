package ValidationLicenseService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type validationLicenseService struct {
	service.AbstractService
}

var ValidationLicenseService = validationLicenseService{}.New()

func (input validationLicenseService) New() (output validationLicenseService) {
	output.FileName = "ValidationLicenseService.go"
	return
}

func (input validationLicenseService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ValidationLicenseRequest) errorModel.ErrorModel) (inputStruct in.ValidationLicenseRequest, err errorModel.ErrorModel) {
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