package SubDistrictService

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

type subDistrictService struct {
	service.AbstractService
	service.GetListData
}

var SubDistrictService = subDistrictService{}.New()

func (input subDistrictService) New() (output subDistrictService) {
	output.FileName = "SubDistrictService.go"
	output.ServiceName = constanta.SubDistrict
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"name",
		"code",
		"updated_at",
	}
	output.ValidSearchBy = []string{"name", "code", "district_id"}
	return
}

func (input subDistrictService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.SubDistrictRequest) errorModel.ErrorModel) (inputStruct in.SubDistrictRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string
	var errorS error

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {return}

	if stringBody != "" {
		errorS = json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	if validation != nil {
		err = validation(&inputStruct)
	}

	return
}
