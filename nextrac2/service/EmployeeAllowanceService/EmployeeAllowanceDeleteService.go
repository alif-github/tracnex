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
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type employeeAllowanceDeleteService struct {
	FileName string
	service.AbstractService
}

var EmployeeAllowanceDeleteService = employeeAllowanceDeleteService{FileName: "EmployeeAllowanceDeleteService.go"}

func (input employeeAllowanceDeleteService) DeleteEmployeeAllowance(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	//authAccessToken := context.AuthAccessTokenModel
	inputStruct, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("DeleteEmployeeAllowance", inputStruct, context, input.deleteAllowance, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", context)

	return
}

func (input employeeAllowanceDeleteService) deleteAllowance(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "deleteAllowance"
	authAccessToken := context.AuthAccessTokenModel

	allowanceBody := body.(in.EmployeeAllowanceRequest)

	allowanceRepository := input.getAllowanceRepository(allowanceBody.ID, authAccessToken, now)
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

	countData, err := dao.EmployeeFacilitiesActiveDAO.GetCountEmployeeMatrixForMaster(serverconfig.ServerAttribute.DBConnection, "allowance_id", allowanceOnDB.ID.Int64)
	if err.Error != nil {
		return
	}

	if countData != 0 {
		err = errorModel.GenerateCannotDeleteData("", funcName)
		return
	}

	err = input.validation(allowanceOnDB, allowanceBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeAllowanceDAO.TableName, allowanceRepository.ID.Int64, userID)...)

	err = dao.EmployeeAllowanceDAO.DeleteEmployeeAllowance(tx, allowanceRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeAllowanceDeleteService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (allowanceBody in.EmployeeAllowanceRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	allowanceBody, bodySize, err := getAllowanceBody(request, input.FileName)
	allowanceBody.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeAllowanceDeleteService) getAllowanceRepository(id int64, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmpAllowanceModel {
	return repository.EmpAllowanceModel{
		ID:            sql.NullInt64{Int64: id},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		Deleted:       sql.NullBool{Bool: true},
	}
}

func (input employeeAllowanceDeleteService) validation(allowanceOnDB repository.EmpAllowanceModel, allowanceBody in.EmployeeAllowanceRequest) (err errorModel.ErrorModel) {
	result, err := TimeStrToTime(allowanceBody.UpdatedAtStr, "updated_at")
	allowanceBody.UpdatedAt = result
	if err.Error != nil {
		return
	}

	err = service.OptimisticLock(allowanceOnDB.UpdatedAt.Time, allowanceBody.UpdatedAt, input.FileName, "allowance")
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