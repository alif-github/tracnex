package EmployeeService

import (
	"database/sql"
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

func (input employeeService) GetEmployeeTimeSheetRedmineByNIK(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.EmployeeRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validationGetEmployeeTimeSheetRedmineByNIK)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetEmployeeTimeSheetRedmineByNIK(inputStruct)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input employeeService) doGetEmployeeTimeSheetRedmineByNIK(inputStruct in.EmployeeRequest) (output interface{}, err errorModel.ErrorModel) {
	var (
		fileName      = "TimeSheetCheckEmployeeByNIKService.go"
		funcName      = "doGetEmployeeTimeSheetRedmineByNIK"
		db            *sql.DB
		dbResult      repository.EmployeeModel
		customFieldID string
		isSqlParam    bool
	)

	switch int(inputStruct.DepartmentId) {
	case constanta.DeveloperDepartmentID, constanta.QAQCDepartmentID, constanta.UIUXDepartmentID:
		db = serverconfig.ServerAttribute.RedmineDBConnection
		customFieldID = "35"
	case constanta.InfraDepartmentID, constanta.DevOpsDepartmentID:
		db = serverconfig.ServerAttribute.RedmineInfraDBConnection
		customFieldID = "1"
		isSqlParam = true
	default:
	}

	dbResult, err = dao.RedmineDAO.GetEmployeeRedmineByNIK(db, inputStruct.IDCard, customFieldID, isSqlParam)
	if err.Error != nil {
		return
	}

	if dbResult.IDCard.String == "" {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.IDCard)
		return
	}

	output = input.convertModelToResponseGetListRedmine(dbResult)
	return
}

func (input employeeService) validationGetEmployeeTimeSheetRedmineByNIK(inputStruct *in.EmployeeRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateGetEmployeeTimeSheetRedmineByNIK()
}

func (input employeeService) convertModelToResponseGetListRedmine(dbResult repository.EmployeeModel) (result interface{}) {
	return out.ViewEmployeeByNIKResponse{
		IDCard:    dbResult.IDCard.String,
		RedmineId: dbResult.RedmineId.Int64,
		FirstName: dbResult.FirstName.String,
		LastName:  dbResult.LastName.String,
	}
}
