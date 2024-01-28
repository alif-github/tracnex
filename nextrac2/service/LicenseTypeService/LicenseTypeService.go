package LicenseTypeService

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

type licenseTypeService struct {
	service.AbstractService
	service.GetListData
}

var LicenseTypeService = licenseTypeService{}.New()

func (input licenseTypeService) New() (output licenseTypeService) {
	output.FileName = "LicenseTypeService.go"
	output.ServiceName = constanta.LicenseTypeConstanta
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"license_type_name",
		"license_type_desc",
		"created_at",
		"updated_at",
		"updated_name",
	}
	output.ValidSearchBy = []string{"license_type_name", "id"}

	return
}

func (input licenseTypeService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.LicenseTypeRequest) errorModel.ErrorModel) (inputStruct in.LicenseTypeRequest, err errorModel.ErrorModel) {
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

func (input licenseTypeService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_license_type_license_type_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.LicenseTypeName)
		}
	}

	return err
}