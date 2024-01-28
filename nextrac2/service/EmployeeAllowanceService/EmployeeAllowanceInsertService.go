package EmployeeAllowanceService

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

type employeeAllowanceInsertService struct {
	FileName string
	service.AbstractService
}

var EmployeeAllowanceInsertService = employeeAllowanceInsertService{FileName: "EmployeeAllowanceInsertService.go"}

func (input employeeAllowanceInsertService) InsertEmployeeAllowance(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	allowanceBody, bodySize, err := getAllowanceBody(request, input.FileName)
	context.LoggerModel.ByteIn = bodySize
	if err.Error != nil {
		return
	}

	err = allowanceBody.ValidateEmployeeAllowance(false)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("InsertEmployeeAllowance", allowanceBody, context, input.doInsert, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", context)

	return
}

func (input employeeAllowanceInsertService) doInsert(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	authAccessToken := context.AuthAccessTokenModel
	allowanceBody := body.(in.EmployeeAllowanceRequest)
	allowance := repository.EmpAllowanceModel{
		AllowanceName: sql.NullString{String: allowanceBody.AllowanceName},
		Type:          sql.NullString{String: allowanceBody.AllowanceType},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}

	idx, err := dao.EmployeeAllowanceDAO.CheckAllowance(tx, allowanceBody.AllowanceName, "allowance_name")
	if err.Error != nil {
		return
	}
	if idx >= 1 {
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, "doInsert", "allowance")
		return
	}

	ids, err := dao.EmployeeAllowanceDAO.CheckAllowance(tx, "annual leave", "allowance_type")
	if err.Error != nil {
		return
	}

	if ids >= 1 && allowance.Type.String == "annual leave" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "doInsert", "annual leave sudah digunakan", "allowance_type", "")
		return
	}

	if ids >= 1 && allowance.Type.String == "cuti tahunan" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "doInsert", "annual leave / cuti tahunan sudah digunakan", "allowance_type", "")
		return
	}

	idCT, err := dao.EmployeeAllowanceDAO.CheckAllowance(tx, "cuti tahunan", "allowance_type")
	if err.Error != nil {
		return
	}

	if idCT >= 1 && allowance.Type.String == "cuti tahunan" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "doInsert", "cuti tahunan sudah digunakan", "allowance_type", "")
		return
	}

	if idCT >= 1 && allowance.Type.String == "annual leave" {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "doInsert", "annual leave / cuti tahunan sudah digunakan", "allowance_type", "")
		return
	}

	lastId, err := dao.EmployeeAllowanceDAO.InsertEmployeeAllowance(tx, allowance)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeeAllowanceDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: lastId},
	})
	return
}
