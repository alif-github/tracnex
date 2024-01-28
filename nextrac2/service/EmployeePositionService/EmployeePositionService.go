package EmployeePositionService

import (
	"database/sql"
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

type employeePositionService struct {
	service.AbstractService
	service.GetListData
}

var EmployeePositionService = employeePositionService{}.New()

func (input employeePositionService) New() (output employeePositionService) {
	output.FileName = "EmployeePositionService.go"
	output.ServiceName = "Position"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"name",
		"description",
	}
	output.ValidSearchBy = []string{
		"name",
		"id",
	}
	return
}

func (input employeePositionService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.EmployeePosition) errorModel.ErrorModel) (inputStruct in.EmployeePosition, err errorModel.ErrorModel) {
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

func (input employeePositionService) checkCompany(tx *sql.Tx, companyID int64) (err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "checkCompany"
		isExist  bool
	)

	isExist, err = dao.CompanyDAO.CheckCompanyByID(tx, companyID)
	if err.Error != nil {
		return
	}

	if !isExist {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.CompanyID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
