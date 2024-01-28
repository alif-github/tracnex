package ParameterService

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

type parameterService struct {
	service.AbstractService
	service.GetListData
}

var ParameterService = parameterService{}.New()

func (input parameterService) New() (output parameterService) {
	output.FileName = "ParameterService.go"
	output.ValidLimit = []int{10, 20, 50, 100, 200, 500}
	output.ValidOrderBy = []string{"id", "permission", "name"}
	output.ValidSearchBy = []string{"permission", "name"}
	return
}

func (input parameterService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ParameterRequest) errorModel.ErrorModel) (inputStruct in.ParameterRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])

	inputStruct.ID = int64(id)

	err = validation(&inputStruct)

	return
}

func (input parameterService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if service.CheckDBError(err, "uq_parameter_permission_name") {
		return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Permission)
	} else {
		return err
	}
}
