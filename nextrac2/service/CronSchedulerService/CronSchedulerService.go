package CronSchedulerService

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

type cronSchedulerService struct {
	service.AbstractService
	service.GetListData
}

var CronScheduler = cronSchedulerService{}.New()

func (input cronSchedulerService) New() (output cronSchedulerService) {
	output.FileName = "CronSchedulerService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id"}
	output.ValidSearchBy = []string{"name"}
	return
}

func (input cronSchedulerService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.CronSchedulerRequest) errorModel.ErrorModel ) (inputStruct in.CronSchedulerRequest, err errorModel.ErrorModel){
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
