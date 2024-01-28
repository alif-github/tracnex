package PositionService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type positionService struct {
	service.AbstractService
	service.GetListData
}

var PositionService = positionService{}.New()

func (input positionService) New() (output positionService) {
	output.FileName = "PositionService.go"
	output.ServiceName = constanta.Position
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "name"}
	output.ValidSearchBy = []string{"name"}
	return
}

func (input positionService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *master_data_request.PositionGetListRequest) errorModel.ErrorModel) (inputStruct master_data_request.PositionGetListRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateErrorFormatJSON(input.FileName, funcName, errorS)
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
