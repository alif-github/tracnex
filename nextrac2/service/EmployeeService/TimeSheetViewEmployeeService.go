package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input employeeService) ViewEmployeeTimeSheet(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.EmployeeRequest
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewEmployee)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewEmployeeTimeSheet(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input employeeService) doViewEmployeeTimeSheet(inputStruct in.EmployeeRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName     = "doViewEmployeeTimeSheet"
		db           = serverconfig.ServerAttribute.DBConnection
		employeeOnDB repository.EmployeeModel
	)

	employeeOnDB, err = dao.EmployeeDAO.ViewEmployeeTimeSheet(db, repository.EmployeeModel{ID: sql.NullInt64{Int64: inputStruct.ID}})
	if err.Error != nil {
		return
	}

	if employeeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.EmployeeConstanta)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, employeeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetailTimeSheet(employeeOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) convertModelToResponseDetailTimeSheet(inputModel repository.EmployeeModel) out.ViewEmployeeTimeSheetResponse {
	var (
		trackerDev in.TrackerDeveloper
		trackerQA  in.TrackerQA
	)

	if inputModel.DepartmentId.Int64 == 2 {
		_ = json.Unmarshal([]byte(inputModel.MandaysRate.String), &trackerQA)
	} else {
		_ = json.Unmarshal([]byte(inputModel.MandaysRate.String), &trackerDev)
	}

	return out.ViewEmployeeTimeSheetResponse{
		ID:                    inputModel.ID.Int64,
		IDCard:                inputModel.IDCard.String,
		RedmineId:             inputModel.RedmineId.Int64,
		Name:                  inputModel.Name.String,
		DepartmentID:          inputModel.DepartmentId.Int64,
		DepartmentName:        inputModel.DepartmentName.String,
		MandaysRate:           trackerDev.Task,
		MandaysRateAutomation: trackerQA.Automation,
		MandaysRateManual:     trackerQA.Manual,
		CreatedAt:             inputModel.CreatedAt.Time,
		UpdatedAt:             inputModel.UpdatedAt.Time,
		CreatedName:           inputModel.CreatedName.String,
		UpdatedName:           inputModel.UpdatedName.String,
	}
}
