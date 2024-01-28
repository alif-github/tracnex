package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

func (input employeeService) UpdateEmployeeTimeSheet(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateEmployeeTimeSheet"
		inputStruct in.EmployeeRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateEmployeeTimeSheet)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateEmployeeTimeSheet, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) doUpdateEmployeeTimeSheet(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct  = inputStructInterface.(in.EmployeeRequest)
		db           = serverconfig.ServerAttribute.DBConnection
		inputModel   repository.EmployeeModel
		variableOnDB repository.EmployeeVariableModel
	)

	//--- Convert DTO Model
	inputModel, err = input.convertDTOToEmployeeTimeSheetUpdate(inputStruct, *contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//--- Check and Validate
	err = input.checkAndLockEmployeeTimeSheetOnDB(inputStruct, inputModel, contextModel)
	if err.Error != nil {
		return
	}

	//--- Check NIK And Compare ID Redmine
	if err = input.checkNikOnRedmine(inputModel); err.Error != nil {
		return
	}

	//--- Get Variable
	variableOnDB, err = dao.EmployeeVariableDAO.GetEmployeeVariableByEmployeeID(db, repository.EmployeeVariableModel{EmployeeID: inputModel.ID})
	if err.Error != nil {
		return
	}

	//--- Insert Variable
	if variableOnDB.ID.Int64 < 1 {
		var idVariable int64
		if idVariable, err = dao.EmployeeVariableDAO.InsertEmployeeVariable(tx, repository.EmployeeVariableModel{
			EmployeeID:    inputModel.ID,
			RedmineID:     inputModel.RedmineId,
			MandaysRate:   inputModel.MandaysRate,
			CreatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:     inputModel.UpdatedAt,
			UpdatedBy:     inputModel.UpdatedBy,
			UpdatedClient: inputModel.UpdatedClient,
		}); err.Error != nil {
			return
		}

		//--- Audit
		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.EmployeeVariableDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: idVariable},
			Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
		})

	} else {
		//--- Audit
		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.EmployeeVariableDAO.TableName, variableOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)
		if err = dao.EmployeeVariableDAO.UpdateEmployeeVariableByEmployeeID(tx, repository.EmployeeVariableModel{
			ID:            variableOnDB.ID,
			RedmineID:     inputModel.RedmineId,
			MandaysRate:   inputModel.MandaysRate,
			UpdatedAt:     inputModel.UpdatedAt,
			UpdatedBy:     inputModel.UpdatedBy,
			UpdatedClient: inputModel.UpdatedClient,
		}); err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) convertDTOToEmployeeTimeSheetUpdate(inputStruct in.EmployeeRequest, contextModel applicationModel.ContextModel, timeNow time.Time) (repo repository.EmployeeModel, err errorModel.ErrorModel) {
	var (
		errs     error
		byteData []byte
	)

	if inputStruct.DepartmentId == 2 {
		jsonData := in.TrackerQA{
			Automation: inputStruct.MandaysRateAutomation,
			Manual:     inputStruct.MandaysRateManual,
		}
		byteData, errs = json.Marshal(jsonData)
		if errs != nil {
			err = errorModel.GenerateSimpleErrorModel(500, errs.Error())
			return
		}
	} else {
		jsonData := in.TrackerDeveloper{Task: inputStruct.MandaysRate}
		byteData, errs = json.Marshal(jsonData)
		if errs != nil {
			err = errorModel.GenerateSimpleErrorModel(500, errs.Error())
			return
		}
	}

	repo = repository.EmployeeModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		IDCard:        sql.NullString{String: inputStruct.IDCard},
		RedmineId:     sql.NullInt64{Int64: inputStruct.RedmineId},
		DepartmentId:  sql.NullInt64{Int64: inputStruct.DepartmentId},
		MandaysRate:   sql.NullString{String: string(byteData)},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	return
}

func (input employeeService) checkAndLockEmployeeTimeSheetOnDB(inputStruct in.EmployeeRequest, model repository.EmployeeModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName     = "TimeSheetUpdateEmployeeService.go"
		funcName     = "checkAndLockEmployeeTimeSheetOnDB"
		db           = serverconfig.ServerAttribute.DBConnection
		employeeOnDB repository.EmployeeModel
		nikResp      repository.EmployeeModel
	)

	//--- Get Employee Variable
	employeeOnDB, err = dao.EmployeeVariableDAO.GetEmployeeVariableForUpdate(db, model, false)
	if err.Error != nil {
		return
	}

	if employeeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.EmployeeConstanta)
		return
	}

	if employeeOnDB.IsHaveVariable.Bool {
		err = input.CheckUserLimitedByOwnAccess(contextModel, employeeOnDB.CreatedBy.Int64)
		if err.Error != nil {
			return
		}
	}

	if employeeOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.EmployeeConstanta)
		return
	}

	//--- Lock Employee Variable
	_, err = dao.EmployeeVariableDAO.GetEmployeeVariableForUpdate(db, model, true)
	if err.Error != nil {
		return
	}

	//--- Check By ID Card
	nikResp, err = dao.EmployeeDAO.CheckEmployeeByNIK(db, model)
	if err.Error != nil {
		return
	}

	if nikResp.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.IDCard)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) validateUpdateEmployeeTimeSheet(inputStruct *in.EmployeeRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateUpdateEmployeeTimeSheet()
}

func (input employeeService) checkNikOnRedmine(inputModel repository.EmployeeModel) (err errorModel.ErrorModel) {
	var (
		fileName      = "TimeSheetUpdateEmployeeService.go"
		funcName      = "checkNikOnRedmine"
		db            *sql.DB
		customFieldID string
		redmineResult repository.EmployeeModel
		errorS        error
		isSqlParam    bool
	)

	//--- Set Custom Field ID
	switch int(inputModel.DepartmentId.Int64) {
	case constanta.DeveloperDepartmentID, constanta.QAQCDepartmentID, constanta.UIUXDepartmentID:
		db = serverconfig.ServerAttribute.RedmineDBConnection
		customFieldID = "35"
		fmt.Println("Department -> ", inputModel.DepartmentId.Int64)
	case constanta.InfraDepartmentID, constanta.DevOpsDepartmentID:
		db = serverconfig.ServerAttribute.RedmineInfraDBConnection
		customFieldID = "1"
		isSqlParam = true
		fmt.Println("Department -> ", inputModel.DepartmentId.Int64)
	default:
	}

	if errorS = db.Ping(); errorS != nil {
		fmt.Println("Error -> ", errorS.Error())
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	fmt.Println("Custom Fields -> ", customFieldID)
	redmineResult, err = dao.RedmineDAO.GetEmployeeRedmineByNIK(db, inputModel.IDCard.String, customFieldID, isSqlParam)
	if err.Error != nil {
		return
	}

	if redmineResult.IDCard.String == "" {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.IDCard)
		return
	}

	if redmineResult.RedmineId.Int64 != inputModel.RedmineId.Int64 {
		err = errorModel.GenerateDifferentRequestAndDBResult(fileName, funcName, strconv.Itoa(int(redmineResult.RedmineId.Int64)), strconv.Itoa(int(inputModel.RedmineId.Int64)))
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
