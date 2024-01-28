package ClientTypeService

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

type clientTypeService struct {
	service.AbstractService
	service.GetListData
}

var ClientTypeService = clientTypeService{}.New()

func (input clientTypeService) New() (output clientTypeService) {
	output.FileName = "ClientTypeService.go"
	output.ServiceName = "CLIENT_TYPE"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "client_type", "description", "created_at", "updated_name"}
	output.ValidSearchBy = []string{"client_type_id", "client_type"}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "ct.id",
		Count: "ct.id",
	}
	output.ListScope = input.SetListScope()
	return
}

func (input clientTypeService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ClientTypeRequest) errorModel.ErrorModel) (inputStruct in.ClientTypeRequest, err errorModel.ErrorModel) {
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

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}
