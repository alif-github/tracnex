package GenerateInternalToken

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type generateInternalTokenService struct {
	service.AbstractService
}

var GenerateInTokenService = generateInternalTokenService{}.New()

func (input generateInternalTokenService) New() (output generateInternalTokenService) {
	output.FileName = "RoleService.go"
	return
}

func (input generateInternalTokenService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.GenerateInternalTokenRequestDTO, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)
	return
}
