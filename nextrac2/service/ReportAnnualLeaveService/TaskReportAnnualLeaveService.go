package ReportAnnualLeaveService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type taskChildAnnualLeave struct {
	Group        string
	Type         string
	Name         string
	GetCountData int64
	DoJob        func(interface{}, *applicationModel.ContextModel, time.Time, *repository.JobProcessModel) errorModel.ErrorModel
}

func (input reportAnnualLeaveService) doServiceJobProcessAnnualLeave(inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time, job *repository.JobProcessModel) {
	var (
		percent = int64(100)
		db      = serverconfig.ServerAttribute.DBConnection
		task    = taskChildAnnualLeave{
			Group:        "Leave-Detail",
			Type:         "File",
			Name:         "Upload File Detail Annual Leave",
			GetCountData: percent,
			DoJob:        input.reportAnnualLeaveService,
		}
	)

	//--- Add Grouping Name
	job.Group.String = task.Group
	job.Type.String = task.Type
	job.Name.String = task.Name

	//--- Job Process
	go input.serviceJobProcessAnnualLeave(inputStructInterface, db, *job, task, *contextModel, timeNow)
}

func (input reportAnnualLeaveService) serviceJobProcessAnnualLeave(inputStructInterface interface{}, db *sql.DB, childJob repository.JobProcessModel, task taskChildAnnualLeave, contextModel applicationModel.ContextModel, timeNow time.Time) {
	var err errorModel.ErrorModel
	defer func() {
		if err.Error != nil {
			childJob.Status.String = constanta.JobProcessErrorStatus
			childJob.UpdatedAt.Time = timeNow
			if childJob.ContentDataOut.String == "" {
				childJob.ContentDataOut.String = service.GetErrorMessage(err, contextModel)
			}

			err = dao.JobProcessDAO.UpdateErrorJobProcess(db, childJob)
			if err.Error != nil {
				input.LogError(err, contextModel) //--- Log Error
			}
		} else {
			childJob.UpdatedAt.Time = timeNow
			err = dao.JobProcessDAO.UpdateJobProcessCounter(db, childJob)
			if err.Error != nil {
				input.LogError(err, contextModel)
			}

			if childJob.Total.Int32 == childJob.Counter.Int32 {
				successLog := fmt.Sprintf("Success do %s data %d from %d", childJob.Name.String, childJob.Counter.Int32, childJob.Total.Int32)
				service.LogMessage(successLog, http.StatusOK)
			}
		}
	}()

	childJob.Level.Int32 = 0
	childJob.Counter.Int32 = 0
	childJob.Total.Int32 = int32(task.GetCountData)

	//--- Insert Job Process
	err = dao.JobProcessDAO.InsertJobProcess(db, childJob)
	if err.Error != nil {
		return
	}

	go input.DoUpdateJobEveryXMinute(constanta.UpdateLastUpdateTimeInMinute, childJob, contextModel)
	if task.DoJob != nil {
		err = task.DoJob(inputStructInterface, &contextModel, timeNow, &childJob)
		if err.Error != nil {
			return
		}
	}
}
