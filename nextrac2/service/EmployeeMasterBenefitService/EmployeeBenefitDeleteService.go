package EmployeeMasterBenefitService

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

type employeeBenefitDeleteService struct {
	FileName string
	service.AbstractService
}

var EmployeeBenefitDeleteService = employeeBenefitDeleteService{FileName: "EmployeeBenefitDeleteService.go"}

func (input employeeBenefitDeleteService) DeleteEmployeeBenefit(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	//authAccessToken := context.AuthAccessTokenModel
	inputStruct, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("DeleteEmployeeBenefit", inputStruct, context, input.deleteBenefit, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", context)

	return
}

func (input employeeBenefitDeleteService) deleteBenefit(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "deleteBenefit"
	authAccessToken := context.AuthAccessTokenModel

	benefitBody := body.(in.EmpBenefitRequest)

	benefitRepository := input.getBenefitRepository(benefitBody.ID, authAccessToken, now)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if isOnlyHaveOwnAccess {
		benefitRepository.CreatedBy.Int64 = userID
	}

	benefitOnDB, err := dao.EmployeeMasterBenefitDAO.GetDetailEmployeeBenefit(tx, benefitRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if benefitOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "benefit")
		return
	}

	countData, err := dao.EmployeeFacilitiesActiveDAO.GetCountEmployeeMatrixForMaster(serverconfig.ServerAttribute.DBConnection, "benefit_id", benefitOnDB.ID.Int64)
	if err.Error != nil {
		return
	}

	if countData != 0 {
		err = errorModel.GenerateCannotDeleteData("", funcName)
		return
	}

	err = input.validation(benefitOnDB, benefitBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeMasterBenefitDAO.TableName, benefitOnDB.ID.Int64, userID)...)

	err = dao.EmployeeMasterBenefitDAO.DeleteEmployeeBenefit(tx, benefitRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeBenefitDeleteService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.EmpBenefitRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	inputStruct, bodySize, err := getEmpBenefitBody(request, input.FileName)
	inputStruct.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeBenefitDeleteService) getBenefitRepository(id int64, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmpBenefitModel {
	return repository.EmpBenefitModel{
		ID:            sql.NullInt64{Int64: id},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		Deleted:       sql.NullBool{Bool: true},
	}
}

func (input employeeBenefitDeleteService) validation(benefitOnDB repository.EmpBenefitModel, benefitBody in.EmpBenefitRequest) (err errorModel.ErrorModel) {
	result, err := TimeStrToTime(benefitBody.UpdatedAtStr, "updated_at")
	benefitBody.UpdatedAt = result
	if err.Error != nil {
		return
	}

	err = service.OptimisticLock(benefitOnDB.UpdatedAt.Time, benefitBody.UpdatedAt, input.FileName, "benefit")
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
