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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type employeeBenefitUpdateService struct {
	FileName string
	service.AbstractService
}

var EmployeeBenefitUpdateService = employeeBenefitUpdateService{FileName: "EmployeeBenefitUpdateService.go"}

func (input employeeBenefitUpdateService) UpdateEmployeeBenefit(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	BenefitBody, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("UpdateEmployeeBenefit", BenefitBody, context, input.doUpdate, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", context)

	return
}

func (input employeeBenefitUpdateService) doUpdate(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdate"
	authAccessToken := context.AuthAccessTokenModel

	BenefitBody := body.(in.EmpBenefitRequest)
	BenefitRepository := input.getBenefitRepository(BenefitBody, authAccessToken, now)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if isOnlyHaveOwnAccess {
		BenefitRepository.CreatedBy.Int64 = userID
	}

	BenefitOnDB, err := dao.EmployeeMasterBenefitDAO.GetDetailEmployeeBenefit(tx, BenefitRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if BenefitOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employeeBenefit")
		return
	}

	idx, err := dao.EmployeeMasterBenefitDAO.CheckBenefit(tx, BenefitBody.BenefitName, "benefit_name")
	if err.Error != nil {
		return
	}
	if idx != 0 {
		if idx != BenefitOnDB.ID.Int64 {
			err = errorModel.GenerateAlreadyExistDataError(input.FileName, "doInsert", "benefit")
			return
		}
	}

	ids, err := dao.EmployeeMasterBenefitDAO.CheckBenefit(tx, "medical", "benefit_type")
	if err.Error != nil {
		return
	}
	if ids != 0 {
		if ids != BenefitOnDB.ID.Int64 && BenefitBody.BenefitType == "medical" {
			err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "doInsert", "medical sudah digunakan", "benefit_type", "")
			return
		}
	}

	err = input.validation(BenefitOnDB, BenefitBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeMasterBenefitDAO.TableName, BenefitOnDB.ID.Int64, userID)...)

	err = dao.EmployeeMasterBenefitDAO.UpdateEmployeeMasterBenefit(tx, BenefitRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input employeeBenefitUpdateService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (BenefitBody in.EmpBenefitRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	BenefitBody, bodySize, err := getEmpBenefitBody(request, input.FileName)
	BenefitBody.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeBenefitUpdateService) getBenefitRepository(Benefit in.EmpBenefitRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmpBenefitModel {
	return repository.EmpBenefitModel{
		ID:            sql.NullInt64{Int64: Benefit.ID},
		BenefitName:   sql.NullString{String: Benefit.BenefitName},
		BenefitType:   sql.NullString{String:Benefit.BenefitType},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}
}

func (input employeeBenefitUpdateService) validation(BenefitOnDB repository.EmpBenefitModel, BenefitBody in.EmpBenefitRequest) (err errorModel.ErrorModel) {
	err = BenefitBody.ValidateEmployeeMasterBenefit(true)
	if err.Error != nil {
		return
	}
	err = service.OptimisticLock(BenefitOnDB.UpdatedAt.Time, BenefitBody.UpdatedAt, input.FileName, "employee_Benefit")
	return
}
