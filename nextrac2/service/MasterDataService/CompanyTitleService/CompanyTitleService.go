package CompanyTitleService

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

type companyTitleService struct {
	service.AbstractService
	service.GetListData
}

var CompanyTitleService = companyTitleService{}.New()

func (input companyTitleService) New() (output companyTitleService) {
	output.FileName = "CompanyTitleService.go"
	output.ServiceName = constanta.CompanyTitle
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "title"}
	output.ValidSearchBy = []string{"title"}
	return
}

func (input companyTitleService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *master_data_request.CompanyTitleRequest) errorModel.ErrorModel) (inputStruct master_data_request.CompanyTitleRequest, err errorModel.ErrorModel) {
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