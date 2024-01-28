package EmployeeAllowanceService

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

type employeeAllowanceUpdateService struct {
	FileName string
	service.AbstractService
}

var EmployeeAllowanceUpdateService = employeeAllowanceUpdateService{FileName: "EmployeeAllowanceUpdateService.go"}

func (input employeeAllowanceUpdateService) UpdateEmployeeALlowance(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	allowanceBody, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("UpdateEmployeeALlowance", allowanceBody, context, input.doUpdate, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", context)

	return
}

func (input employeeAllowanceUpdateService) doUpdate(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdate"
	authAccessToken := context.AuthAccessTokenModel

	allowanceBody := body.(in.EmployeeAllowanceRequest)
	allowanceRepository := input.getAllowanceRepository(allowanceBody, authAccessToken, now)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if isOnlyHaveOwnAccess {
		allowanceRepository.CreatedBy.Int64 = userID
	}

	allowanceOnDB, err := dao.EmployeeAllowanceDAO.GetDetailEmployeeAllowance(tx, allowanceRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if allowanceOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "allowance")
		return
	}

	idx, err := dao.EmployeeAllowanceDAO.CheckAllowance(tx, allowanceBody.AllowanceName, "allowance_name")
	if err.Error != nil {
		return
	}
	if idx != 0 && idx != allowanceOnDB.ID.Int64 {
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, funcName, "allowance")
		return
	}

	ids, err := dao.EmployeeAllowanceDAO.CheckAllowance(tx, "annual leave", "allowance_type")
	if err.Error != nil {
		return
	}
	if ids != 0 && ids != allowanceOnDB.ID.Int64 && allowanceBody.AllowanceType == "annual leave" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "annual leave sudah digunakan", "allowance_type", "")
		return
	}
	if ids != 0 && ids != allowanceOnDB.ID.Int64 && allowanceBody.AllowanceType == "cuti tahunan" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "cuti tahunan sudah digunakan", "allowance_type", "")
		return
	}

	idCT, err := dao.EmployeeAllowanceDAO.CheckAllowance(tx, "cuti tahunan", "allowance_type")
	if err.Error != nil {
		return
	}
	if idCT != 0 && idCT != allowanceOnDB.ID.Int64 && allowanceBody.AllowanceType == "cuti tahunan" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "cuti tahunan sudah digunakan", "allowance_type", "")
		return
	}
	if idCT != 0 && idCT != allowanceOnDB.ID.Int64 && allowanceBody.AllowanceType == "annual leave" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "annual leave sudah digunakan", "allowance_type", "")
		return
	}

	err = input.validation(allowanceOnDB, allowanceBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeAllowanceDAO.TableName, allowanceOnDB.ID.Int64, userID)...)

	err = dao.EmployeeAllowanceDAO.UpdateEmployeeAllowance(tx, allowanceRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeAllowanceUpdateService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (allowanceBody in.EmployeeAllowanceRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	allowanceBody, bodySize, err := getAllowanceBody(request, input.FileName)
	allowanceBody.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeAllowanceUpdateService) getAllowanceRepository(allowance in.EmployeeAllowanceRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmpAllowanceModel {
	return repository.EmpAllowanceModel{
		ID:            sql.NullInt64{Int64: allowance.ID},
		AllowanceName: sql.NullString{String: allowance.AllowanceName},
		Type       :   sql.NullString{String:allowance.AllowanceType},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}
}

func (input employeeAllowanceUpdateService) validation(allowanceOnDB repository.EmpAllowanceModel, allowanceBody in.EmployeeAllowanceRequest) (err errorModel.ErrorModel) {
	err = allowanceBody.ValidateEmployeeAllowance(true)
	if err.Error != nil {
		return
	}
	err = service.OptimisticLock(allowanceOnDB.UpdatedAt.Time, allowanceBody.UpdatedAt, input.FileName, "employee_allowance")
	return
}