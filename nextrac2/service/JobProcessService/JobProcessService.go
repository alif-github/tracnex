package JobProcessService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type jobProcessService struct {
	service.AbstractService
	service.GetListData
}

var JobProcessService = jobProcessService{}.New()

func (input jobProcessService) New() (output jobProcessService) {
	output.FileName = "JobProcessService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "job_id", "status"}
	output.ValidSearchBy = []string{"job_id", "status"}
	return
}

func (input jobProcessService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.JobProcessRequest) errorModel.ErrorModel) (inputStruct in.JobProcessRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	inputStruct.JobID = mux.Vars(request)["ID"]

	err = validation(&inputStruct)

	return
}
