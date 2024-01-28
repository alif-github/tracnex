package ComponentService

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

type componentService struct {
	service.AbstractService
	service.GetListData
}

var ComponentService = componentService{}.New()

func (input componentService) New() (output componentService) {
	output.FileName = "ComponentService.go"
	output.ServiceName = constanta.Component
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"component_name",
		"created_at",
		"updated_at",
		"updated_name",
	}
	output.ValidSearchBy = []string{"component_name", "id"}
	return
}

func (input componentService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ComponentRequest) errorModel.ErrorModel) (inputStruct in.ComponentRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {return}

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

func (input componentService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_component_component_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ComponentName)
		}
	}

	return err
}