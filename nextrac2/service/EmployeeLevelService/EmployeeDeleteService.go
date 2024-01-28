package EmployeeLevelService

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

type employeeLevelDeleteService struct {
	FileName string
	service.AbstractService
}

var EmployeeLevelDeleteService = employeeLevelDeleteService{FileName: "EmployeeLevelDeleteService.go"}

func (input employeeLevelDeleteService) DeleteEmplooyeeLevel(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	//authAccessToken := context.AuthAccessTokenModel
	levelBody, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("DeleteEmplooyeeLevel", levelBody, context, input.deleteLevel, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", context)

	return
}

func (input employeeLevelDeleteService) deleteLevel(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "deleteLevel"
	authAccessToken := context.AuthAccessTokenModel

	levelBody := body.(in.EmployeeLevelRequest)
	levelRepository := input.getLevelRepository(levelBody.ID, authAccessToken, now)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if isOnlyHaveOwnAccess {
		levelRepository.CreatedBy.Int64 = userID
	}

	levelOnDB, err := dao.EmployeeLevelDAO.GetDetailEmployeeLevel(tx, levelRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if levelOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employeeLevel")
		return
	}

	countData, err := dao.EmployeeFacilitiesActiveDAO.GetCountEmployeeMatrixForMaster(serverconfig.ServerAttribute.DBConnection, "employee_level_id", levelRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if countData != 0 {
		err = errorModel.GenerateCannotDeleteData("", funcName)
		return
	}

	err = input.validation(levelOnDB, levelBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeLevelDAO.TableName, levelRepository.ID.Int64, userID)...)

	err = dao.EmployeeLevelDAO.DeleteEmployeeLevel(tx, levelRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeLevelDeleteService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (levelBody in.EmployeeLevelRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	levelBody, bodySize, err := getLevelBody(request, input.FileName)
	levelBody.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeLevelDeleteService) getLevelRepository(categoryId int64, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmployeeLevelModel {
	return repository.EmployeeLevelModel{
		ID:            sql.NullInt64{Int64: categoryId},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		Deleted:       sql.NullBool{Bool: true},
	}
}

func (input employeeLevelDeleteService) validation(categoryOnDB repository.EmployeeLevelModel, levelBody in.EmployeeLevelRequest) (err errorModel.ErrorModel) {
	result, err := TimeStrToTime(levelBody.UpdatedAtStr, "updated_at")
	levelBody.UpdatedAt = result
	if err.Error != nil {
		return
	}

	err = service.OptimisticLock(categoryOnDB.UpdatedAt.Time, levelBody.UpdatedAt, input.FileName, "employee_level")
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