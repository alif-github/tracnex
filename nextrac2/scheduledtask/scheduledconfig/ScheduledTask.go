package scheduledconfig

import (
	"database/sql"
	"github.com/robfig/cron"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/scheduledtask/Task"
	"os"
)

var ScheduledCron map[string]*cron.Cron

func startScheduledTask() {
	for key := range ScheduledCron {
		if ScheduledCron[key] != nil {
			ScheduledCron[key].Start()
		}
	}
}

func GenerateSchedulerCron(db *sql.DB) {
	ScheduledCron = make(map[string]*cron.Cron)
	hostName, _ := os.Hostname()

	cronModel, err := dao.SchedulerDAO.GetListSchedulerByHostName(db, hostName)
	if err.Error != nil {
		//todo log
		return
	}

	for i := 0; i < len(cronModel); i++ {
		checkScheduler(cronModel[i])
	}

	startScheduledTask()
}

func generateCron(cronTime string, cmd func()) *cron.Cron {
	c := cron.New()
	_, _ = c.AddFunc(cronTime, cmd)
	return c
}

func checkScheduler(cronModel repository.CRONSchedulerModel) {
	switch cronModel.RunType.String {
	case Task.SyncElasticBank.RunType:
		ScheduledCron[Task.SyncElasticBank.RunType] = generateCron(cronModel.CRON.String, Task.SyncElasticBank.Start)
		break
	case Task.SchAddResourceNexcloud.RunType:
		ScheduledCron[Task.SchAddResourceNexcloud.RunType] = generateCron(cronModel.CRON.String, Task.SchAddResourceNexcloud.Start)
		break
	case Task.CheckExpirationProductLicense.RunType:
		ScheduledCron[Task.CheckExpirationProductLicense.RunType] = generateCron(cronModel.CRON.String, Task.CheckExpirationProductLicense.Start)
		break
	case Task.HouseKeepingFileUpload.RunType:
		ScheduledCron[Task.HouseKeepingFileUpload.RunType] = generateCron(cronModel.CRON.String, Task.HouseKeepingFileUpload.Start)
		break
	case Task.SyncRegionalMDB.RunType:
		ScheduledCron[Task.SyncRegionalMDB.RunType] = generateCron(cronModel.CRON.String, Task.SyncRegionalMDB.Start)
		break
	case Task.HistoryEmployeeBenefits.RunType:
		ScheduledCron[Task.HistoryEmployeeBenefits.RunType] = generateCron(cronModel.CRON.String, Task.HistoryEmployeeBenefits.Start)
		break
	case Task.CutOffLeaveTask.RunType:
		ScheduledCron[Task.CutOffLeaveTask.RunType] = generateCron(cronModel.CRON.String, Task.CutOffLeaveTask.Start)
		break
	default:
		break
	}
}

func RestartScheduler(runType string, db *sql.DB, logModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	cronHostModel, err := dao.SchedulerDAO.GetDataForRestartScheduler(db, runType)
	if err.Error != nil {
		return
	}

	for i := 0; i < len(cronHostModel); i++ {
		go hitOtherHostToRestart(cronHostModel[i].HostURL.String, runType, *logModel)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func RestartOwnScheduler(runType string, db *sql.DB, _ *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	hostName, _ := os.Hostname()

	cronHostModel, err := dao.SchedulerDAO.GetDataForRestartOwnScheduler(db, runType, hostName)
	if err.Error != nil {
		return
	}

	if cronHostModel.CRON.String != "" {
		if ScheduledCron[runType] != nil {
			ScheduledCron[runType].Stop()
		}

		checkScheduler(repository.CRONSchedulerModel{
			RunType: sql.NullString{String: runType},
			CRON:    sql.NullString{String: cronHostModel.CRON.String},
		})

		if ScheduledCron[runType] != nil {
			ScheduledCron[runType].Start()
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func hitOtherHostToRestart(hostUrl string, task string, contextModel applicationModel.ContextModel) {
	header := make(map[string][]string)
	header[common.DefaultTokenKeyConstanta] = []string{common.DefaultTokenValueConstanta}

	status, _, body, err := common.HitAPI(hostUrl+"/v1/internal/scheduler/"+task, header, "", "POST", contextModel)
	if err != nil {
		contextModel.LoggerModel.Status = 500
		contextModel.LoggerModel.Message = err.Error()
		return
	}

	if status != 200 {
		contextModel.LoggerModel.Status = status
		contextModel.LoggerModel.Message = body
		return
	}
}
