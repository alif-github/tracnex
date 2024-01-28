package HostServerService

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

type hostServerService struct {
	service.AbstractService
	service.GetListData
}

type viewDetailHostServerService struct {
	service.AbstractService
	service.GetListData
}

var HostServerService = hostServerService{}.New()
var ViewDetailHostServerService = viewDetailHostServerService{}.NewViewDetail()

func (input hostServerService) New() (output hostServerService) {
	output.FileName = "HostServerService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "host_name"}
	output.ValidSearchBy = []string{"host_name"}
	return
}

func (input viewDetailHostServerService) NewViewDetail() (output viewDetailHostServerService) {
	output.FileName = "HostServerService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "name"}
	output.ValidSearchBy = []string{"name"}
	return
}

func (input hostServerService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.HostServerRequest) errorModel.ErrorModel ) (inputStruct in.HostServerRequest, err errorModel.ErrorModel){
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input hostServerService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_hostserver_hostname") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Name)
		}
	}
	return err
}
