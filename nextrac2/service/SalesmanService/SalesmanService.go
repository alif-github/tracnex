package SalesmanService

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

type salesmanService struct {
	service.AbstractService
	service.GetListData
}

var SalesmanService = salesmanService{}.New()

func (input salesmanService) New() (output salesmanService) {
	output.FileName = "SalesmanService.go"
	output.ServiceName = "SALESMAN"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"first_name",
		"status",
		"address",
		"district",
		"province",
		"phone",
		"email"}
	output.ValidSearchBy = []string{"id", "first_name"}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.SalesmanDataScope] = applicationModel.MappingScopeDB{
		View:  "salesman.id",
		Count: "salesman.id",
	}
	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View:  "salesman.province_id",
		Count: "salesman.province_id",
	}
	output.MappingScopeDB[constanta.DistrictDataScope] = applicationModel.MappingScopeDB{
		View:  "salesman.district_id",
		Count: "salesman.district_id",
	}
	output.ListScope = input.SetListScope()
	return
}

func (input salesmanService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.SalesmanRequest) errorModel.ErrorModel) (inputStruct in.SalesmanRequest, err errorModel.ErrorModel) {
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

func (input salesmanService) readBodyAndValidateForView(request *http.Request, validation func(input *in.SalesmanRequest) errorModel.ErrorModel) (inputStruct in.SalesmanRequest, err errorModel.ErrorModel) {

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_salesman_nik") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.NIK)
		}
	}

	return err
}

func (input salesmanService) validateDataScopeSalesman(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "validateDataScopeSalesman"

	output = service.ValidateScope(contextModel, []string{
		constanta.SalesmanDataScope,
		constanta.ProvinceDataScope,
		constanta.DistrictDataScope,
	})

	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
