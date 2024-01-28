package PersonProfileService

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

type personProfileService struct {
	service.AbstractService
	service.GetListData
}

var PersonProfileServie = personProfileService{}.New()

func (input personProfileService) New() (output personProfileService) {
	output.FileName = "PersonProfileService.go"
	output.ServiceName = constanta.PersonProfile
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "first_name", "last_name"}
	output.ValidSearchBy = []string{"first_name", "email", "phone", "nik"}
	return
}

func (input personProfileService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *master_data_request.PersonProfileGetListRequest) errorModel.ErrorModel) (inputStruct master_data_request.PersonProfileGetListRequest, err errorModel.ErrorModel) {
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
	if inputStruct.ID  == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}
