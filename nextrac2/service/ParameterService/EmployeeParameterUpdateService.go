package ParameterService

import (
	"database/sql"
	"fmt"
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
	"strconv"
	"strings"
	"time"
)

type employeeParameterUpdateService struct {
	FileName string
	service.AbstractService
}

var EmployeeParameterUpdateService = employeeParameterUpdateService{FileName: "EmployeeParameterUpdateService.go"}

func (input employeeParameterUpdateService) UpdateEmployeeParameter(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	paramBody, bodySize, err := getParameterBody(request, input.FileName)
	context.LoggerModel.ByteIn = bodySize
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("UpdateEmployeeParameter", paramBody, context, input.doUpdate, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", context)

	return
}

func (input employeeParameterUpdateService) doUpdate(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdate"
	authAccessToken := context.AuthAccessTokenModel

	paramsBody := body.([]in.ParameterRequest)

	for i:=0; i<len(paramsBody);i++  {
		errs := paramsBody[i].ValidateUpdateParameterEmployee()
		if errs.Error != nil {
			err = errs
			return
		}

		paramRepository := input.getParameterRepository(paramsBody[i], authAccessToken, now)

		paramOnDB, errParamDb := dao.ParameterDAO.GetDetailUpdateParameter(tx, paramRepository.ID.Int64)
		if errParamDb.Error != nil {
			err = errParamDb
			return
		}

		if paramOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "parameter")
			return
		}

		errValidate := input.validation(paramOnDB, paramsBody[i])
		if errValidate.Error != nil {
			err = errValidate
			return
		}

		auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.ParameterDAO.TableName, paramOnDB.ID.Int64, authAccessToken.ResourceUserID)...)

		errUpdate := dao.ParameterDAO.UpdateParameterEmployee(tx, paramRepository)
		if errUpdate.Error != nil {
			err = errUpdate
			return
		}

		if paramsBody[i].Name == "cutOffAnualLeave"{
			dateArr := strings.Split(paramsBody[i].Value, "-")
			dd, _ := strconv.ParseInt(dateArr[0], 10, 64)
			mm, _ := strconv.ParseInt(dateArr[1], 10, 64)

			errUpdate := dao.CronSchedulerDAO.UpdatedCronSchedulerByRunType(tx, repository.CRONSchedulerModel{
				RunType: sql.NullString{String:"scheduler.cutOff_annual_leave"},
				CRON:    sql.NullString{String: fmt.Sprintf("0 0 %s %s *", strconv.Itoa(int(dd)),strconv.Itoa(int(mm)))},
			})
			if errUpdate.Error != nil {
				err = errUpdate
				return
			}
		}
	}

	return
}

func (input employeeParameterUpdateService) getParameterRepository(param in.ParameterRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.ParameterModel {
	return repository.ParameterModel{
		ID:            sql.NullInt64{Int64: param.ID},
		Value:         sql.NullString{String: param.Value},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}
}

func (input employeeParameterUpdateService) validation(paramOnDB repository.ParameterModel, paramBody in.ParameterRequest) (err errorModel.ErrorModel) {
	err = paramBody.ValidateUpdateParameterEmployee()
	if err.Error != nil {
		return
	}
	err = service.OptimisticLock(paramOnDB.UpdatedAt.Time, paramBody.UpdatedAt, input.FileName, "paramter")
	return
}
