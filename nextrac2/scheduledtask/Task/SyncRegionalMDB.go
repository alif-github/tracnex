package Task

import (
	"fmt"
	"log"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/TaskSchedulerService"
)

type syncRegionalMDB struct {
	AbstractScheduledTask
}

var SyncRegionalMDB = syncRegionalMDB{}.New()

func (input syncRegionalMDB) New() (output syncRegionalMDB) {
	output.RunType = "scheduler.regional_mdb"
	return
}

func (input syncRegionalMDB) Start() {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.StartTask(input.RunType, input.doSyncRegionalMDB)
	}
}

func (input syncRegionalMDB) doSyncRegionalMDB() {
	var (
		job repository.JobProcessModel
		err errorModel.ErrorModel
	)

	defer func() {
		if err.Error != nil {
			service.LogMessage(fmt.Sprintf(`auto scheduler error : %s`, err.Error.Error()), err.Code)
		} else {
			service.LogMessage(fmt.Sprintf(`job id : %s`, job.JobID.String), 200)
		}
	}()

	log.Printf("[AUTO-SCHEDULER] Service Sync Regional Data")
	job, err = TaskSchedulerService.TaskSchedulerService.SchedulerSyncRegionalData(&applicationModel.ContextModel{})
}
