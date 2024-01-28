package EmployeeContractService

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

type employeeContractService struct {
	service.AbstractService
	service.GetListData
}

var EmployeeContractService = employeeContractService{}.New()

func (input employeeContractService) New() (output employeeContractService) {
	output.FileName = "EmployeeContractService.go"
	output.ServiceName = constanta.EmployeeContract
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"created_at",
		"contract_no",
		"from_date",
		"thru_date",
		"information",
	}
	output.ValidSearchBy = []string{
		"contract_no",
		"employee_id",
	}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "ec.employee_id",
		Count: "employee_id",
	}
	return
}

func (input employeeContractService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.EmployeeContractRequest) errorModel.ErrorModel) (inputStruct in.EmployeeContractRequest, err errorModel.ErrorModel) {
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

func CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_employee_contract_contractno") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.EmployeeContract)
		}
	}
	return err
}
