package DataScopeService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type dataScopeService struct {
	service.AbstractService
}

var DataScopeService = dataScopeService{}.New()

func (input dataScopeService) New() (output dataScopeService) {
	output.FileName = "DataScopeService.go"
	output.ServiceName = "DATA_SCOPE"
	return
}

func (input dataScopeService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ScopeRequest) errorModel.ErrorModel) (inputStruct in.ScopeRequest, err errorModel.ErrorModel) {
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
