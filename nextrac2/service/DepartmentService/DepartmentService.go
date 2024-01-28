package DepartmentService

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

type departmentService struct {
	service.AbstractService
	service.GetListData
}

var DepartmentService = departmentService{}.New()

func (input departmentService) New() (output departmentService) {
	output.FileName = "DepartmentService.go"
	output.ServiceName = constanta.DepartmentConstanta
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"name",
		"updated_at",
	}
	output.ValidSearchBy = []string{
		"id",
		"name",
	}

	return
}

func (input departmentService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.DepartmentRequest) errorModel.ErrorModel) (inputStruct in.DepartmentRequest, err errorModel.ErrorModel) {
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
