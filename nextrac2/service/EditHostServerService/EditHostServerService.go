package EditHostServerService

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

type editHostServerService struct {
	service.AbstractService
	service.GetListData
}

var EditHostServerService = editHostServerService{}.New()

func (input editHostServerService) New() (output editHostServerService) {
	output.FileName = "EditHostServerService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id"}
	output.ValidSearchBy = []string{"cron_id"}
	return
}

func (input editHostServerService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.CronHostRequest) errorModel.ErrorModel) (inputStruct in.CronHostRequest, err errorModel.ErrorModel) {
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

func (input editHostServerService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_hostserver_hostname") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Name)
		}
	}
	return err
}