package EmployeeFacilitiesActiveService

import (
	"database/sql"
	"net/http"
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

type employeeMatrixInsertService struct {
	FileName string
	service.AbstractService
}

var EmployeeMatrixInsertService = employeeMatrixInsertService{FileName: "EmployeeMatrixInsertService.go"}

func (input employeeMatrixInsertService) InsertEmployeeMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	matrixBody, bodySize, err := getEmpMatrixBody(request, input.FileName)
	context.LoggerModel.ByteIn = bodySize
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("InsertEmployeeMatrix", matrixBody, context, input.doInsert, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", context)

	return
}

func (input employeeMatrixInsertService) doInsert(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	authAccessToken := context.AuthAccessTokenModel
	matrixBody := body.(in.EmployeeMatrixRequest)

	funcName := "doInsert"

	gradeOnDB, err := dao.EmployeeGradeDAO.GetDetailEmployeeGrade(tx, matrixBody.GradeID)
	if err.Error != nil {
		return
	}

	if gradeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "grade")
		return
	}

	levelOnDB, err := dao.EmployeeLevelDAO.GetDetailEmployeeLevel(tx, matrixBody.LevelID)
	if err.Error != nil {
		return
	}

	if levelOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "level")
		return
	}

	count, err := dao.EmployeeFacilitiesActiveDAO.GetDetailEmployeeMatrixForUpdate(serverconfig.ServerAttribute.DBConnection, matrixBody.LevelID, matrixBody.GradeID)
	if err.Error != nil {
		return
	}

	if count >= 1 {
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, funcName, "matrix")
		return
	}

	for i:=0; i < len(matrixBody.AllowanceList); i++{
		allowanceRepo := input.getEmpMatrixRepository(matrixBody, "allowance", int64(i), authAccessToken, now)
		allowanceOnDB, err2 := dao.EmployeeAllowanceDAO.GetDetailEmployeeAllowance(tx, matrixBody.AllowanceList[i].ID)
		if err2.Error != nil {
			err = err2
			return
		}

		errValue := input.validasiValue(matrixBody.AllowanceList[i].Value)
		if errValue.Error != nil {
			err = errValue
			return
		}

		if allowanceOnDB.ID.Int64 < 1 {
			err2 = errorModel.GenerateUnknownDataError(input.FileName, funcName, "allowance")
			err = err2
			return
		}

		lastId, err2 := dao.EmployeeFacilitiesActiveDAO.InsertEmployeeFacilitiesActive(tx, allowanceRepo)
		if err2.Error != nil {
			err = err2
			return
		}

		auditData = append(auditData, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.EmployeeFacilitiesActiveDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: lastId},
		})
	}

	for i:=0; i < len(matrixBody.BenefitList); i++{
		benefitRepo := input.getEmpMatrixRepository(matrixBody, "benefit", int64(i), authAccessToken, now)
		benefitOnDB, err2 := dao.EmployeeMasterBenefitDAO.GetDetailEmployeeBenefit(tx, matrixBody.BenefitList[i].ID)
		if err2.Error != nil {
			err = err2
			return
		}

		errValue := input.validasiValue(matrixBody.BenefitList[i].Value)
		if errValue.Error != nil {
			err = errValue
			return
		}

		if benefitOnDB.ID.Int64 < 1 {
			err2 = errorModel.GenerateUnknownDataError(input.FileName, funcName, "benefit")
			err = err2
			return
		}

		lastId, err2 := dao.EmployeeFacilitiesActiveDAO.InsertEmployeeFacilitiesActive(tx, benefitRepo)
		if err2.Error != nil {
			err = err2
			return
		}

		auditData = append(auditData, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.EmployeeFacilitiesActiveDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: lastId},
		})
	}

	return
}

func (input employeeMatrixInsertService) getEmpMatrixRepository(matrix in.EmployeeMatrixRequest, typ string, index int64, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmployeeFacilitiesActiveModel {
	value := ""
	active := false
	var benefitID, allowanceID int64
	benefitID = 0
	allowanceID = 0

	if typ == "allowance"{
		allowanceID = matrix.AllowanceList[index].ID
        value = matrix.AllowanceList[index].Value
        active = matrix.AllowanceList[index].Active
	}
	if typ == "benefit"{
		benefitID = matrix.BenefitList[index].ID
		value = matrix.BenefitList[index].Value
		active = matrix.BenefitList[index].Active
	}
	return repository.EmployeeFacilitiesActiveModel{
		LevelID:       sql.NullInt64{Int64: matrix.LevelID},
		GradeID:       sql.NullInt64{Int64: matrix.GradeID},
		AllowanceID  : sql.NullInt64{Int64: allowanceID},
		BenefitID:     sql.NullInt64{Int64: benefitID},
		Active       : sql.NullBool{Bool:   active},
		Value        : sql.NullString{String: value},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}
}

func (input employeeMatrixInsertService) validasiValue(value string) (err errorModel.ErrorModel){
	if value != "" && len(value) > 100{
       err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "validasiValue", "value harus kurang dari 100 karakter", "value", "")
       return err
	}

	return err
}