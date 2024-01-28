package PostalCodeService

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

type postalCodeService struct {
	service.AbstractService
	service.GetListData
}

var PostalCodeService = postalCodeService{}.New()

func (input postalCodeService) New() (output postalCodeService) {
	output.FileName = "PostalCodeService.go"
	output.ServiceName = constanta.PostalCode
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"code",
		"updated_at",
	}
	output.ValidSearchBy = []string{"code", "urban_village_id"}
	return
}

func (input postalCodeService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.PostalCodeRequest) errorModel.ErrorModel) (inputStruct in.PostalCodeRequest, err errorModel.ErrorModel) {
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