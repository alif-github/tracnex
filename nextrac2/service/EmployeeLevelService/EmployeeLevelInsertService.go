package EmployeeLevelService

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

type employeeLevelInsertService struct {
	FileName string
	service.AbstractService
}

var EmployeeLevelInsertService = employeeLevelInsertService{FileName: "EmployeeLevelInsertService.go"}

func (input employeeLevelInsertService) InsertEmployeeLevel(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	levelBody, bodySize, err := getLevelBody(request, input.FileName)
	context.LoggerModel.ByteIn = bodySize
	if err.Error != nil {
		return
	}

	err = levelBody.ValidateEmployeeLevel(false)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("InsertEmployeeLevel", levelBody, context, input.doInsert, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", context)

	return
}

func (input employeeLevelInsertService) doInsert(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	authAccessToken := context.AuthAccessTokenModel
	levelBody := body.(in.EmployeeLevelRequest)
	level := repository.EmployeeLevelModel{
		Level:         sql.NullString{String: levelBody.Level},
		Description:   sql.NullString{String: levelBody.Description},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}

	idx, err := dao.EmployeeLevelDAO.CheckLevel(tx, levelBody.Level)
	if err.Error != nil {
		return
	}
	if idx >= 1 {
		err = errorModel.GenerateAlreadyExistDataError(input.FileName, "doInsert", "level")
		return
	}

	lastId, err := dao.EmployeeLevelDAO.InsertEmployeeLevel(tx, level)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeeLevelDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: lastId},
	})
	return
}
