package EmployeeGradeService

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

type employeeGradeInsertService struct {
	FileName string
	service.AbstractService
}

var EmployeeGradeInsertService = employeeGradeInsertService{FileName: "employeeGradeInsertService.go"}

func (input employeeGradeInsertService) InsertEmployeeGrade(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	gradeBody, bodySize, err := getGradeBody(request, input.FileName)
	context.LoggerModel.ByteIn = bodySize
	if err.Error != nil {
		return
	}

	err = gradeBody.ValidateEmployeeGrade(false)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("InsertEmployeeGrade", gradeBody, context, input.doInsert, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", context)

	return
}

func (input employeeGradeInsertService) doInsert(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	authAccessToken := context.AuthAccessTokenModel
	gradeBody := body.(in.EmployeeGradeRequest)
	grade := repository.EmployeeGradeModel{
		Grade:         sql.NullString{String: gradeBody.Grade},
		Description:   sql.NullString{String:gradeBody.Description},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}

	idx, err := dao.EmployeeGradeDAO.CheckGrade(tx, gradeBody.Grade)
	if err.Error != nil {
		return
	}
	if idx >= 1 {
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, "doInsert", "grade")
		return
	}

	lastId, err := dao.EmployeeGradeDAO.InsertEmployeeGrade(tx, grade)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeeGradeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: lastId},
	})
	return
}