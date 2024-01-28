package EmployeeMasterBenefitService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type employeeBenefitInsertService struct {
	FileName string
	service.AbstractService
}

var EmployeeBenefitInsertService = employeeBenefitInsertService{FileName: "EmployeeBenefitInsertService.go"}

func (input employeeBenefitInsertService) InsertEmployeeBenefit(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	benefitBody, bodySize, err := getEmpBenefitBody(request, input.FileName)
	context.LoggerModel.ByteIn = bodySize
	if err.Error != nil {
		return
	}

	err = benefitBody.ValidateEmployeeMasterBenefit(false)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("InsertEmployeeBenefit", benefitBody, context, input.doInsert, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", context)

	return
}

func (input employeeBenefitInsertService) doInsert(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	authAccessToken := context.AuthAccessTokenModel
	benefitBody := body.(in.EmpBenefitRequest)
	benefit := repository.EmpBenefitModel{
		BenefitName:   sql.NullString{String: benefitBody.BenefitName},
		BenefitType:   sql.NullString{String: benefitBody.BenefitType},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}

	idx, err := dao.EmployeeMasterBenefitDAO.CheckBenefit(tx, benefitBody.BenefitName, "benefit_name")
	if err.Error != nil {
		return
	}

	if idx >= 1 {
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, "doInsert", "benefit")
		return
	}

	ids, err := dao.EmployeeMasterBenefitDAO.CheckBenefit(tx, "medical", "benefit_type")
	if err.Error != nil {
		return
	}

	if ids >= 1 && benefit.BenefitType.String == "medical" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "doInsert", "medical sudah digunakan", "benefit_type", "")
		return
	}

	lastId, err := dao.EmployeeMasterBenefitDAO.InsertEmployeeMasterBenefit(tx, benefit)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeeMasterBenefitDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: lastId},
	})
	return
}

