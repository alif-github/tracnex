package EmployeeGradeService

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
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type employeeGradeUpdateService struct {
	FileName string
	service.AbstractService
}

var EmployeeGradeUpdateService = employeeGradeUpdateService{FileName: "employeeGradeUpdateService.go"}

func (input employeeGradeUpdateService) UpdateEmployeeGrade(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	gradeBody, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("UpdateEmployeeGrade", gradeBody, context, input.doUpdate, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", context)

	return
}

func (input employeeGradeUpdateService) doUpdate(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdate"
	authAccessToken := context.AuthAccessTokenModel

	gradeBody := body.(in.EmployeeGradeRequest)
	gradeRepository := input.getGradeRepository(gradeBody, authAccessToken, now)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if isOnlyHaveOwnAccess {
		gradeRepository.CreatedBy.Int64 = userID
	}

	gradeOnDB, err := dao.EmployeeGradeDAO.GetDetailEmployeeGrade(tx, gradeRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if gradeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employeeGrade")
		return
	}

	idx, err := dao.EmployeeGradeDAO.CheckGrade(tx, gradeBody.Grade)
	if err.Error != nil {
		return
	}
	if idx != 0 && idx != gradeOnDB.ID.Int64 {
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, "doInsert", "grade")
		return
	}

	err = input.validation(gradeOnDB, gradeBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeLevelDAO.TableName, gradeOnDB.ID.Int64, userID)...)

	err = dao.EmployeeGradeDAO.UpdateEmployeeGrade(tx, gradeRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeGradeUpdateService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (gradeBody in.EmployeeGradeRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	gradeBody, bodySize, err := getGradeBody(request, input.FileName)
	gradeBody.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeGradeUpdateService) getGradeRepository(grade in.EmployeeGradeRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmployeeGradeModel {
	return repository.EmployeeGradeModel{
		ID:            sql.NullInt64{Int64: grade.ID},
		Grade:         sql.NullString{String: grade.Grade},
		Description:   sql.NullString{String:grade.Description},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}
}

func (input employeeGradeUpdateService) validation(gradeOnDB repository.EmployeeGradeModel, gradeBody in.EmployeeGradeRequest) (err errorModel.ErrorModel) {
	err = gradeBody.ValidateEmployeeGrade(true)
	if err.Error != nil {
		return
	}
	err = service.OptimisticLock(gradeOnDB.UpdatedAt.Time, gradeBody.UpdatedAt, input.FileName, "employee_grade")
	return
}

