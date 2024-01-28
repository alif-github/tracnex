package PersonTitleService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type personTitleService struct {
	service.AbstractService
	service.GetListData
}

var PersonTitleService = personTitleService{}.New()

func (input personTitleService) New() (output personTitleService) {
	output.FileName = "PersonTitleService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "title"}
	output.ValidSearchBy = []string{"id", "title", "description", "en_description"}
	return
}

func (input personTitleService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.PersonTitleRequest) errorModel.ErrorModel) (inputStruct in.PersonTitleRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)

	return
}