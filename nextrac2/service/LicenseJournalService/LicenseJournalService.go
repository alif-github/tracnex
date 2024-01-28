package LicenseJournalService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type licenseJournalService struct {
	service.AbstractService
	service.GetListData
}

var LicenseJournalService = licenseJournalService{}.New()

func (input licenseJournalService) New() (output licenseJournalService) {
	output.FileName = "LicenseJournalService.go"
	output.ValidLimit = service.DefaultLimit
	return
}

func (input licenseJournalService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.LicenseJournalRequest) errorModel.ErrorModel) (inputStruct in.LicenseJournalRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if request.Method != "GET" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}
