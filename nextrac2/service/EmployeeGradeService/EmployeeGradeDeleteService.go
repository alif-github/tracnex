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
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type employeeGradeDeleteService struct {
	FileName string
	service.AbstractService
}

var EmployeeGradeDeleteService = employeeGradeDeleteService{FileName: "EmployeeGradeDeleteService.go"}

func (input employeeGradeDeleteService) DeleteEmployeeGrade(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	//authAccessToken := context.AuthAccessTokenModel
	gradeBody, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("DeleteEmployeeGrade", gradeBody, context, input.deleteGrade, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", context)

	return
}

func (input employeeGradeDeleteService) deleteGrade(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "deleteGrade"
	authAccessToken := context.AuthAccessTokenModel

	gradeBody := body.(in.EmployeeGradeRequest)

	gradeRepository := input.getGradeRepository(gradeBody.ID, authAccessToken, now)
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
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "grade")
		return
	}

	countData, err := dao.EmployeeFacilitiesActiveDAO.GetCountEmployeeMatrixForMaster(serverconfig.ServerAttribute.DBConnection, "employee_grade_id", gradeRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if countData != 0 {
		err = errorModel.GenerateCannotDeleteData("", funcName)
		return
	}

	err = input.validation(gradeOnDB, gradeBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeLevelDAO.TableName, gradeRepository.ID.Int64, userID)...)

	err = dao.EmployeeGradeDAO.DeleteEmployeeGrade(tx, gradeRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeGradeDeleteService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (gradeBody in.EmployeeGradeRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	gradeBody, bodySize, err := getGradeBody(request, input.FileName)
	gradeBody.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeGradeDeleteService) getGradeRepository(id int64, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmployeeGradeModel {
	return repository.EmployeeGradeModel{
		ID:            sql.NullInt64{Int64: id},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		Deleted:       sql.NullBool{Bool: true},
	}
}

func (input employeeGradeDeleteService) validation(gradeOnDB repository.EmployeeGradeModel, gradeBody in.EmployeeGradeRequest) (err errorModel.ErrorModel) {
	result, err := TimeStrToTime(gradeBody.UpdatedAtStr, "updated_at")
	gradeBody.UpdatedAt = result
	if err.Error != nil {
		return
	}

	err = service.OptimisticLock(gradeOnDB.UpdatedAt.Time, gradeBody.UpdatedAt, input.FileName, "employee_grade")
	return
}

func TimeStrToTime(timeStr string, fieldName string) (output time.Time, error errorModel.ErrorModel) {
	output, err := time.Parse(constanta.DefaultTimeFormat, timeStr)

	if err != nil {
		error = errorModel.GenerateFormatFieldError("AbstractDTO.go", "TimeStrToTime", fieldName)
		return
	}
	return output, errorModel.GenerateNonErrorModel()
}
