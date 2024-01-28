package ClientCredentialService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type clientCredentialService struct {
	service.AbstractService
}

var ClientCredentialService = clientCredentialService{}.New()

func (input clientCredentialService) New() (output clientCredentialService) {
	output.FileName = "ClientCredentialService.go"
	return
}

func (input clientCredentialService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input in.ClientCredential) errorModel.ErrorModel) (inputStruct in.ClientCredential, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

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

	err = validation(inputStruct)
	return
}
