package EmployeeService

import (
	"database/sql"
	"errors"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

func (input employeeService) InsertEmployee(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.EmployeeRequest
		files       []in.MultipartFileDTO
		funcName    = "InsertEmployee"
	)

	inputStruct, files, err = input.readBodyWithFileAndValidate(request, contextModel, input.ValidateInsertEmployee)
	if err.Error != nil {
		return
	}

	//--- Add Files
	inputStruct.Files = files

	output.Data.Content, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertEmployee, func(i interface{}, contextModel applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input employeeService) doInsertEmployee(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct   = inputStructInterface.(in.EmployeeRequest)
		inputModel    repository.EmployeeModel
		insertedID    int64
		dataAuditTemp repository.AuditSystemModel
		idBenefits    int64
	)

	//--- Convert To DTO Model
	inputModel = input.convertDTOToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow, false)

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

	//--- Insert Employee
	insertedID, err = dao.EmployeeDAO.InsertEmployee(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	//--- Insert Employee Benefits
	inputModel.ID.Int64 = insertedID
	idBenefits, err = dao.EmployeeBenefitsDAO.InsertEmployeeBenefits(tx, repository.EmployeeBenefitsModel{
		EmployeeID:      inputModel.ID,
		EmployeeLevelID: inputModel.LevelID,
		EmployeeGradeID: inputModel.GradeID,
		BPJSNo:          inputModel.BPJS,
		BPJSTkNo:        inputModel.BPJSTk,
		CreatedBy:       inputModel.CreatedBy,
		CreatedClient:   inputModel.CreatedClient,
		CreatedAt:       inputModel.CreatedAt,
		UpdatedBy:       inputModel.UpdatedBy,
		UpdatedClient:   inputModel.UpdatedClient,
		UpdatedAt:       inputModel.UpdatedAt,
	})
	if err.Error != nil {
		return
	}

	//--- Data Audit Store
	dataAudit = append(dataAudit,
		repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.EmployeeDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: insertedID},
		},
		repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.EmployeeBenefitsDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: idBenefits},
		},
	)

	dataAuditTemp, err = input.GenerateDataScope(tx, insertedID, dao.EmployeeDAO.TableName, constanta.EmployeeDataScope, contextModel.AuthAccessTokenModel.ResourceUserID, contextModel.AuthAccessTokenModel.ClientID, timeNow)
	if err.Error != nil {
		return
	}

	//--- Output
	dataAudit = append(dataAudit, dataAuditTemp)
	outputTemp := make(map[string]interface{})
	outputTemp["id"] = insertedID

	output = outputTemp
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) convertDTOToModel(inputStruct in.EmployeeRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time, isUpdate bool) (repo repository.EmployeeModel) {
	repo = repository.EmployeeModel{
		IDCard:             sql.NullString{String: inputStruct.IDCard},
		NPWP:               sql.NullString{String: inputStruct.NPWP},
		FirstName:          sql.NullString{String: inputStruct.FirstName},
		LastName:           sql.NullString{String: inputStruct.LastName},
		Gender:             sql.NullString{String: inputStruct.Gender},
		Email:              sql.NullString{String: inputStruct.Email},
		Phone:              sql.NullString{String: inputStruct.Phone},
		PlaceOfBirth:       sql.NullString{String: inputStruct.PlaceOfBirth},
		DateOfBirth:        sql.NullTime{Time: inputStruct.DateOfBirth},
		Type:               sql.NullString{String: inputStruct.Type},
		MothersMaiden:      sql.NullString{String: inputStruct.MothersMaiden},
		TaxMethod:          sql.NullString{String: inputStruct.TaxMethod},
		AddressResidence:   sql.NullString{String: inputStruct.AddressResidence},
		AddressTax:         sql.NullString{String: inputStruct.AddressTax},
		Religion:           sql.NullString{String: inputStruct.Religion},
		DateJoin:           sql.NullTime{Time: inputStruct.DateJoin},
		DateOut:            sql.NullTime{Time: inputStruct.DateOut},
		ReasonResignation:  sql.NullString{String: inputStruct.ReasonResignation},
		Status:             sql.NullString{String: inputStruct.Status},
		PositionID:         sql.NullInt64{Int64: inputStruct.Position},
		MaritalStatus:      sql.NullString{String: inputStruct.MaritalStatus},
		NumberOfDependents: sql.NullInt64{Int64: inputStruct.NumberOfDependents},
		Nationality:        sql.NullString{String: inputStruct.Nationality},
		DepartmentId:       sql.NullInt64{Int64: inputStruct.DepartmentId},
		BPJS:               sql.NullString{String: inputStruct.NoBpjs},
		BPJSTk:             sql.NullString{String: inputStruct.NoBpjsTk},
		LevelID:            sql.NullInt64{Int64: inputStruct.Level},
		GradeID:            sql.NullInt64{Int64: inputStruct.Grade},
		IsHaveMember:       sql.NullBool{Bool: inputStruct.IsHaveMember},
		Member:             sql.NullString{String: inputStruct.MemberIDStr},
		Education:          sql.NullString{String: inputStruct.Education},
		Active:             sql.NullBool{Bool: inputStruct.Active},
		UpdatedBy:          sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient:      sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:          sql.NullTime{Time: timeNow},
	}

	if !isUpdate {
		repo.CreatedBy.Int64 = authAccessToken.ResourceUserID
		repo.CreatedClient.String = authAccessToken.ClientID
		repo.CreatedAt.Time = timeNow
	} else {
		repo.ID.Int64 = inputStruct.ID
	}

	return
}

func (input employeeService) ValidateInsertEmployee(inputStruct *in.EmployeeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}

func (input employeeService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_nik_employee") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.NIK)
		}

		if service.CheckDBError(err, "uq_employee_redmine_id") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.RedmineId)
		}

		if service.CheckDBError(err, "uq_employee_email_phone") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.EmailPhone)
		}
	}

	return err
}

func (input employeeService) checkIDMember(inputStruct in.EmployeeRequest) (err errorModel.ErrorModel) {
	var (
		fileName = "ProfileInsertEmployeeService.go"
		funcName = "checkIDMember"
		db       = serverconfig.ServerAttribute.DBConnection
	)

	for _, itemMember := range inputStruct.MemberID {
		var (
			isExists  bool
			firstName string
			idMember  int
			errS      error
		)

		if itemMember == "all" {
			break
		}

		idMember, errS = strconv.Atoi(itemMember)
		if errS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errS)
			return
		}

		if idMember < 1 {
			err = errorModel.GenerateUnknownError(fileName, funcName, errors.New("id shouldn't zero (0)"))
			return
		}

		isExists, firstName, err = dao.EmployeeDAO.CheckEmployeeIDByID(db, int64(idMember))
		if err.Error != nil {
			return
		}

		if !isExists {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, firstName)
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) checkLevelAndGrade(inputModel repository.EmployeeModel) (err errorModel.ErrorModel) {
	var (
		fileName      = "ProfileInsertEmployeeService.go"
		funcName      = "checkLevelAndGrade"
		db            = serverconfig.ServerAttribute.DBConnection
		isExistLevel  bool
		isExistGrade  bool
		isExistMatrix bool
	)

	//--- Check By Level ID
	isExistLevel, err = dao.EmployeeLevelGradeDAO.CheckEmployeeLevel(db, inputModel)
	if err.Error != nil {
		return
	}

	if !isExistLevel {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Level)
		return
	}

	//--- Check By Grade ID
	isExistGrade, err = dao.EmployeeLevelGradeDAO.CheckEmployeeGrade(db, inputModel)
	if err.Error != nil {
		return
	}

	if !isExistGrade {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "Grade")
		return
	}

	//--- Check Matrix
	isExistMatrix, err = dao.EmployeeFacilitiesActiveDAO.CheckEmployeeMatrixIsExist(db, repository.EmployeeFacilitiesActiveModel{
		LevelID: inputModel.LevelID,
		GradeID: inputModel.GradeID,
	})
	if err.Error != nil {
		return
	}

	if !isExistMatrix {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "Matrix")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) checkPositionAndDepartmentID(inputModel repository.EmployeeModel) (err errorModel.ErrorModel) {
	var (
		fileName = "ProfileInsertEmployeeService.go"
		funcName = "checkPositionAndDepartmentID"
		db       = serverconfig.ServerAttribute.DBConnection
		isExist  bool
	)

	//--- Check Employee Position
	isExist, err = dao.EmployeePositionDAO.CheckEmployeePosition(db, inputModel.PositionID.Int64)
	if err.Error != nil {
		return
	}

	if !isExist {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.PositionIDStruct)
		return
	}

	//--- Check Department ID
	isExist = false
	isExist, err = dao.DepartmentDAO.CheckIDDepartment(db, repository.DepartmentModel{ID: sql.NullInt64{Int64: inputModel.DepartmentId.Int64}})
	if err.Error != nil {
		return
	}

	if !isExist {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.DepartmentId)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) uploadPhotoToLocalCDN(tx *sql.Tx, files *[]in.MultipartFileDTO, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	var (
		fileName  = "ProfileInsertEmployeeService.go"
		funcName  = "uploadPhotoToLocalCDN"
		container string
	)

	if len(*files) > 1 {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "Foto yang di upload hanya boleh 1", "Foto", "")
		return
	}

	container = constanta.EmployeeAttachmentPrefix + service.GetAzureDateContainer()

	if err = service.UploadFileToLocalCDN(container, files, contextModel.AuthAccessTokenModel.ResourceUserID); err.Error != nil {
		return
	}

	for i := 0; i < len(*files); i++ {
		var (
			photoID    int64
			photoModel repository.FileUpload
		)

		photoModel = repository.FileUpload{
			Category:      sql.NullString{String: "Photo"},
			Konektor:      sql.NullString{String: "employee"},
			FileName:      sql.NullString{String: (*files)[i].Filename},
			Path:          sql.NullString{String: (*files)[i].Path},
			Host:          sql.NullString{String: (*files)[i].Host},
			CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			CreatedAt:     sql.NullTime{Time: timeNow},
			UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:     sql.NullTime{Time: timeNow},
		}

		photoID, err = dao.FileUploadDAO.InsertFileUploadInfoForBacklog(tx, photoModel)
		if err.Error != nil {
			return
		}

		(*files)[i].FileID = photoID
		*dataAudit = append(*dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.FileUploadDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: photoID},
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
