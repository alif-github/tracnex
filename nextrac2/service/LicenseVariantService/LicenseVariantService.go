package LicenseVariantService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"strings"
)

type licenseVariantService struct {
	service.AbstractService
	service.GetListData
}

var LicenseVariantService = licenseVariantService{}.New()

func (input licenseVariantService) New() (output licenseVariantService) {
	output.FileName = "LicenseVariantService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"license_variant_name",
		"created_at",
		"updated_name",
	}
	output.ValidSearchBy = []string{"id","license_variant_name"}
	return
}

func (input licenseVariantService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.LicenseVariantRequest) errorModel.ErrorModel) (inputStruct in.LicenseVariantRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input licenseVariantService) readBodyAndValidateForView(request *http.Request, validation func(input *in.LicenseVariantRequest) errorModel.ErrorModel) (inputStruct in.LicenseVariantRequest, err errorModel.ErrorModel) {

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_license_variant_license_variant_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.LicenseVariantName)
		}
	}

	return err
}

func (input licenseVariantService) convertLicenseVariantInput(inputStruct in.GetListDataDTO) (newSearchByParam []in.SearchByParam) {

	var valueInput string

	if !util.IsStringEmpty(inputStruct.Search) {
		valueSearch := strings.Split(inputStruct.Search, " ")
		for index, valueSearchItem := range valueSearch {
			if index >= 2 && (len(valueSearch)-(index+1)) != 0 {
				valueInput += valueSearchItem + " "
			} else if index >= 2 && (len(valueSearch)-(index+1)) == 0 {
				valueInput += valueSearchItem
			}
		}

		intValueSearch, errorS := strconv.Atoi(valueInput)
		if intValueSearch >= 0 && errorS == nil {
			newSearchByParam = append(newSearchByParam, in.SearchByParam{
				SearchKey:      "id",
				DataType:       "number",
				SearchOperator: "eq",
				SearchValue:    valueInput,
				SearchType:     constanta.Search,
			}, in.SearchByParam{
				SearchKey:      "license_variant_name",
				DataType:       "char",
				SearchOperator: "like",
				SearchValue:    valueInput,
				SearchType:     constanta.Search,
			})
		} else {
			newSearchByParam = append(newSearchByParam, in.SearchByParam{
				SearchKey:      "license_variant_name",
				DataType:       "char",
				SearchOperator: "like",
				SearchValue:    valueInput,
				SearchType:     constanta.Filter,
			})
		}
	}

	return
}
