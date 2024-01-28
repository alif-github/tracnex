package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	"time"
)

func (input employeeService) UpdateEmployee(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateEmployee"
		inputStruct in.EmployeeRequest
		files       []in.MultipartFileDTO
	)

	inputStruct, files, err = input.readBodyWithFileAndValidate(request, contextModel, input.validateUpdateEmployee)
	if err.Error != nil {
		return
	}

	//--- Add Files
	inputStruct.Files = files

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateEmployee, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) doUpdateEmployee(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName     = "ProfileUpdateEmployeeService.go"
		funcName     = "doUpdateEmployee"
		inputStruct  = inputStructInterface.(in.EmployeeRequest)
		db           = serverconfig.ServerAttribute.DBConnection
		inputModel   repository.EmployeeModel
		employeeOnDB repository.EmployeeModel
	)

	//--- Convert To DTO Model
	inputModel = input.convertDTOToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow, true)

	//--- Check Initial
	employeeOnDB, err = input.EmployeeDAO.GetEmployeeForUpdate(db, inputModel)
	if err.Error != nil {
		return
	}

	if employeeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.EmployeeConstanta)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, employeeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if employeeOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.EmployeeConstanta)
		return
	}

	//--- Check Level And Grade
	if err = input.checkLevelAndGrade(inputModel); err.Error != nil {
		return
	}

	//--- Check ID Member
	if err = input.checkIDMember(inputStruct); err.Error != nil {
		return
	}

	//--- Check Employee Position And Department ID
	if err = input.checkPositionAndDepartmentID(inputModel); err.Error != nil {
		return
	}

	//--- Upload Photo Profile
	if err = input.uploadPhotoToLocalCDN(tx, &inputStruct.Files, contextModel, timeNow, &dataAudit); err.Error != nil {
		return
	}

	if len(inputStruct.Files) > 0 {
		if inputStruct.Files[0].FileID > 0 {
			inputModel.FileUploadID.Int64 = inputStruct.Files[0].FileID
		}
	}

	//--- Update Employee to DB
	err = input.doUpdateEmployeeToDB(tx, inputModel, contextModel, timeNow, &dataAudit)
	if err.Error != nil {
		return
	}

	//--- Update Employee Benefits to DB
	err = input.doInsertUpdateEmployeeBenefitsToDB(tx, inputModel, contextModel, timeNow, &dataAudit)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) validateUpdateEmployee(inputStruct *in.EmployeeRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateUpdate()
}

func (input employeeService) doUpdateEmployeeToDB(tx *sql.Tx, inputModel repository.EmployeeModel, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	//--- Get Record Before
	recordBefore := service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.EmployeeDAO.TableName, inputModel.ID.Int64, 0)

	//--- Update Employee
	err = dao.EmployeeDAO.UpdateEmployee(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	//--- Get Record After
	recordAfter := service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.EmployeeDAO.TableName, inputModel.ID.Int64, 0)

	//--- Create Data Audit
	description := input.descriptionEmployeeHistory(recordBefore, recordAfter)
	recordBeforeTemp := recordBefore
	if description != "" {
		recordBeforeTemp[0].Description.String = description
	}

	//--- Store Data Audit
	*dataAudit = append(*dataAudit, recordBeforeTemp...)
	return
}

func (input employeeService) doInsertUpdateEmployeeBenefitsToDB(tx *sql.Tx, inputModel repository.EmployeeModel, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	//--- Get Employee Benefits
	employeeBenefitsDB, err := dao.EmployeeBenefitsDAO.GetByEmployeeIdForUpdate(tx, inputModel.ID.Int64)
	if err.Error != nil {
		return
	}

	if employeeBenefitsDB.ID.Int64 < 1 {
		err = input.doInsertEmployeeBenefitsToDB(tx, inputModel, employeeBenefitsDB, contextModel, timeNow, dataAudit)
		if err.Error != nil {
			return
		}
	} else {
		err = input.doUpdateEmployeeBenefitsToDB(tx, inputModel, employeeBenefitsDB, contextModel, timeNow, dataAudit)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) doUpdateEmployeeBenefitsToDB(tx *sql.Tx, inputModel repository.EmployeeModel, employeeBenefitsDB repository.EmployeeBenefitsModel, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	//--- Get Record Before
	recordBefore := service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.EmployeeBenefitsDAO.TableName, employeeBenefitsDB.ID.Int64, 0)

	//--- Update Employee Benefits
	if err = dao.EmployeeBenefitsDAO.UpdateEmployeeBenefitsByEmployeeID(tx, repository.EmployeeBenefitsModel{
		EmployeeID:      sql.NullInt64{Int64: inputModel.ID.Int64},
		EmployeeLevelID: sql.NullInt64{Int64: inputModel.LevelID.Int64},
		EmployeeGradeID: sql.NullInt64{Int64: inputModel.GradeID.Int64},
		BPJSNo:          sql.NullString{String: inputModel.BPJS.String},
		BPJSTkNo:        sql.NullString{String: inputModel.BPJSTk.String},
		UpdatedBy:       sql.NullInt64{Int64: inputModel.UpdatedBy.Int64},
		UpdatedAt:       sql.NullTime{Time: timeNow},
		UpdatedClient:   sql.NullString{String: inputModel.UpdatedClient.String},
	}); err.Error != nil {
		return
	}

	//--- Get Record After
	recordAfter := service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.EmployeeBenefitsDAO.TableName, employeeBenefitsDB.ID.Int64, 0)

	//--- Create Data Audit
	description := input.descriptionEmployeeHistory(recordBefore, recordAfter)
	recordBeforeTemp := recordBefore
	if description != "" {
		recordBeforeTemp[0].Description.String = description
	}

	//--- Store Data Audit
	*dataAudit = append(*dataAudit, recordBeforeTemp...)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) doInsertEmployeeBenefitsToDB(tx *sql.Tx, inputModel repository.EmployeeModel, employeeBenefitsDB repository.EmployeeBenefitsModel, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	var (
		fileName      = "ProfileUpdateEmployeeService.go"
		funcName      = "doInsertEmployeeBenefitsToDB"
		idBenefits    int64
		errorS        error
		employeeAfter in.EmployeeJSONDB
		history       string
	)

	//--- Insert Employee Benefits
	if idBenefits, err = dao.EmployeeBenefitsDAO.InsertEmployeeBenefits(tx, repository.EmployeeBenefitsModel{
		EmployeeID:      sql.NullInt64{Int64: inputModel.ID.Int64},
		EmployeeLevelID: sql.NullInt64{Int64: inputModel.LevelID.Int64},
		EmployeeGradeID: sql.NullInt64{Int64: inputModel.GradeID.Int64},
		BPJSNo:          sql.NullString{String: inputModel.BPJS.String},
		BPJSTkNo:        sql.NullString{String: inputModel.BPJSTk.String},
		CreatedBy:       sql.NullInt64{Int64: inputModel.UpdatedBy.Int64},
		CreatedAt:       sql.NullTime{Time: timeNow},
		CreatedClient:   sql.NullString{String: inputModel.UpdatedClient.String},
		UpdatedBy:       sql.NullInt64{Int64: inputModel.UpdatedBy.Int64},
		UpdatedAt:       sql.NullTime{Time: timeNow},
		UpdatedClient:   sql.NullString{String: inputModel.UpdatedClient.String},
	}); err.Error != nil {
		return
	}

	//--- Get Record After
	recordAfter := service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.EmployeeBenefitsDAO.TableName, idBenefits, 0)

	//--- Check if record empty then error
	if len(recordAfter) < 1 {
		err = errorModel.GenerateUnknownError(fileName, funcName, errors.New("record empty"))
		return
	}

	//--- Record After Unmarshal
	errorS = json.Unmarshal([]byte(recordAfter[0].Data.String), &employeeAfter)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	historyTemp := input.mappingAndAppendToHistory(in.EmployeeJSONDB{}, employeeAfter)
	if len(historyTemp) > 0 {
		byteTemp, errs := json.Marshal(historyTemp)
		if errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}
		history = string(byteTemp)
	}

	recordAfterTemp := recordAfter
	if history != "" {
		recordAfterTemp[0].Description.String = history
	}

	//--- Store Data Audit
	*dataAudit = append(*dataAudit, recordAfterTemp...)
	err = errorModel.GenerateNonErrorModel()
	return
}
