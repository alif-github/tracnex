package StandarManhourService

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

type standarManhourService struct {
	StandarManhourDAO dao.StandarManhourDAOInterface
	DepartmentDAO     dao.DepartmentDAOInterface
	service.AbstractService
	service.GetListData
}

var StandarManhourService = standarManhourService{}.New()

func (input standarManhourService) New() (output standarManhourService) {
	output.FileName = "StandarManhourService.go"
	output.ServiceName = constanta.StandarManhour
	output.ValidLimit = service.DefaultLimit
	output.StandarManhourDAO = dao.StandarManhourDAO
	output.DepartmentDAO = dao.DepartmentDAO
	output.ValidOrderBy = []string{
		"id",
		"department_id",
		"case",
		"updated_at",
	}
	output.ValidSearchBy = []string{
		"id",
		"case",
		"department_id",
	}
	return
}

func (input standarManhourService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.StandarManhourRequest) errorModel.ErrorModel) (inputStruct in.StandarManhourRequest, err errorModel.ErrorModel) {
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
