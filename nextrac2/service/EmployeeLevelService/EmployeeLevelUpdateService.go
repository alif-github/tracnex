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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type employeeLevelUpdateService struct {
	FileName string
	service.AbstractService
}

var EmployeeLevelUpdateService = employeeLevelUpdateService{FileName: "employeeLevelUpdateService.go"}

func (input employeeLevelUpdateService) UpdateEmployeeLevel(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	levelBody, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("UpdateEmployeeLevel", levelBody, context, input.doUpdate, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", context)

	return
}

func (input employeeLevelUpdateService) doUpdate(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdate"
	authAccessToken := context.AuthAccessTokenModel

	levelBody := body.(in.EmployeeLevelRequest)
	levelRepository := input.getLevelRepository(levelBody, authAccessToken, now)
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

	idx, err := dao.EmployeeLevelDAO.CheckLevel(tx, levelBody.Level)
	if err.Error != nil {
		return
	}
	if idx != 0 && idx != levelOnDB.ID.Int64{
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, "doInsert", "level")
		return
	}

	err = input.validation(levelOnDB, levelBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeLevelDAO.TableName, levelRepository.ID.Int64, userID)...)

	err = dao.EmployeeLevelDAO.UpdateEmployeeLevel(tx, levelRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeLevelUpdateService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (levelBody in.EmployeeLevelRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	levelBody, bodySize, err := getLevelBody(request, input.FileName)
	levelBody.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeLevelUpdateService) getLevelRepository(level in.EmployeeLevelRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmployeeLevelModel {
	return repository.EmployeeLevelModel{
		ID:            sql.NullInt64{Int64: level.ID},
		Level:         sql.NullString{String: level.Level},
		Description:   sql.NullString{String: level.Description},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}
}

func (input employeeLevelUpdateService) validation(categoryOnDB repository.EmployeeLevelModel, levelBody in.EmployeeLevelRequest) (err errorModel.ErrorModel) {
	err = levelBody.ValidateEmployeeLevel(true)
	if err.Error != nil {
		return
	}
	err = service.OptimisticLock(categoryOnDB.UpdatedAt.Time, levelBody.UpdatedAt, input.FileName, "employee_level")
	return
}
