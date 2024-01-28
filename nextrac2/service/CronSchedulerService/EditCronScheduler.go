package CronSchedulerService

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
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input cronSchedulerService) UpdateCronScheduler(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("UpdateCronScheduler", inputStruct, contextModel, input.doUpdateCronScheduler, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_CRON_SCHEDULER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input cronSchedulerService) doUpdateCronScheduler(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.CronSchedulerRequest)

	cronScheduler := repository.CRONSchedulerModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Name:          sql.NullString{String: inputStruct.Name},
		RunType:       sql.NullString{String: inputStruct.RunType},
		CRON:          sql.NullString{String: inputStruct.Cron},
		Status:        sql.NullBool{Bool: inputStruct.Status},
		UpdatedAt:     sql.NullTime{Time: inputStruct.UpdatedAt},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	dataAudit, err = input.UpdateCronSchedulerOnDB(tx, cronScheduler, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	return
}

func (input cronSchedulerService) UpdateCronSchedulerOnDB(tx *sql.Tx, cronScheduler repository.CRONSchedulerModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "UpdateCronSchedulerOnDB"

	cronScheduler.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	cronSchedulerOnDB, err := dao.CronSchedulerDAO.GetCronSchedulerForUpdate(tx, cronScheduler)
	if err.Error != nil {
		return
	}

	if cronSchedulerOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	if cronSchedulerOnDB.UpdatedAt.Time != cronScheduler.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.CronScheduler)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CronSchedulerDAO.TableName, cronScheduler.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.CronSchedulerDAO.UpdatedCronScheduler(tx, cronScheduler, timeNow)
	if err.Error != nil {
		return
	}

	return
}

func (input cronSchedulerService) validateUpdate(inputStruct *in.CronSchedulerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateCronScheduler()

}
