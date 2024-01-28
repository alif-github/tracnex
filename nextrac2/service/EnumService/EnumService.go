package EnumService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type enumService struct {
	service.AbstractService
	service.GetListData
}

var EnumService = enumService{}.New()

func (input enumService) New() (output enumService) {
	output.FileName = "EnumService.go"
	output.ServiceName = constanta.EnumConstanta
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{}
	output.ValidSearchBy = []string{}

	return
}

func (input enumService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.EnumRequest) errorModel.ErrorModel) (inputStruct in.EnumRequest, err errorModel.ErrorModel) {
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

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}
