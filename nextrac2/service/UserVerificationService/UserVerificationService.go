package UserVerificationService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type userVerificationService struct {
	service.AbstractService
}

var UserVerificationService = userVerificationService{}.New()

func (input userVerificationService) New() (output userVerificationService) {
	output.ServiceName = constanta.UserVerification
	output.FileName = "UserVerificationService.go"
	return
}

func (input userVerificationService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UserVerificationRequest) errorModel.ErrorModel) (inputStruct in.UserVerificationRequest, err errorModel.ErrorModel) {
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

	err = validation(&inputStruct)
	return
}
