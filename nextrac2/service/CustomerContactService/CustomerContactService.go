package CustomerContactService

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

type customerContactService struct {
	service.AbstractService
	service.GetListData
}

var CustomerContacService = customerContactService{}.New()

func (input customerContactService) New() (output customerContactService) {
	output.FileName = "CustomerContactService.go"
	output.ServiceName = constanta.CustomerContact
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"status",
	}
	output.ValidSearchBy = []string{
		"id",
		"nik",
		"province_id",
		"district_id",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View:  "p.id",
		Count: "p.id",
	}
	output.MappingScopeDB[constanta.DistrictDataScope] = applicationModel.MappingScopeDB{
		View:  "d.id",
		Count: "d.id",
	}

	output.ListScope = input.SetListScope()
	return
}

func (input customerContactService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.CustomerContactRequest) errorModel.ErrorModel) (inputStruct in.CustomerContactRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {return}

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

func (input customerContactService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{
		constanta.ProvinceDataScope,
		constanta.DistrictDataScope,
	})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}