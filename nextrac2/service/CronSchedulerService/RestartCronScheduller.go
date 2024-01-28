package CronSchedulerService

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/scheduledtask/scheduledconfig"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input cronSchedulerService) RestartSchedulerService(request *http.Request, contextModel *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel) {
	return restartScheduler(request, contextModel, scheduledconfig.RestartScheduler)
}

func (input cronSchedulerService) InternalRestartSchedulerService(request *http.Request, contextModel *applicationModel.ContextModel) (out.Payload, map[string]string, errorModel.ErrorModel) {
	return restartScheduler(request, contextModel, scheduledconfig.RestartOwnScheduler)
}

func restartScheduler(request *http.Request, contextModel *applicationModel.ContextModel, restartFunc func(runType string, db *sql.DB, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel)) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	runType := mux.Vars(request)["ID"]

	err = restartFunc(runType, serverconfig.ServerAttribute.DBConnection, contextModel)
	if err.Error != nil {
		return
	}

	output.Status.Code = util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil)
	output.Status.Message = util.GenerateConstantaI18n("RESTART_SCHEDULER_SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil)

	err = errorModel.GenerateNonErrorModel()
	return
}
