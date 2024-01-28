package ModuleService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type moduleService struct {
	ModuleDAO dao.ModuleDAOInterface
	service.AbstractService
	service.GetListData
}

var ModuleService = moduleService{}.New()

func (input moduleService) New() (output moduleService) {
	output.FileName = "ModuleService.go"
	output.ServiceName = constanta.Module
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"module_name",
		"created_at",
		"updated_at",
		"updated_name",
	}
	output.ValidSearchBy = []string{"module_name", "id"}
	output.ModuleDAO = dao.ModuleDAO
	return
}

func (input moduleService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ModuleRequest) errorModel.ErrorModel) (inputStruct in.ModuleRequest, err errorModel.ErrorModel) {
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

func (input moduleService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_module_module_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ModuleName)
		}
	}

	return err
}
