package common

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
)

type deleteTokenService struct {
	service.AbstractService
}

var DeleteTokenService = deleteTokenService{}.New()

func (input deleteTokenService) New() (output deleteTokenService) {
	output.FileName = "DeleteTokenService.go"
	return
}

func (input deleteTokenService) StartService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.DeleteTokenDTOIn

	inputStruct, err = input.readBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if inputStruct.Token != "" {
		serverconfig.ServerAttribute.RedisClient.Del(inputStruct.Token)
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateCommonServiceBundleI18NMessage("SUCCESS_DELETE_TOKEN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input deleteTokenService) readBody(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.DeleteTokenDTOIn, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	return
}
