package Task

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
	"time"
)

type cutOffLeaveTask struct {
	AbstractScheduledTask
}

var CutOffLeaveTask = cutOffLeaveTask{}.New()

func (input cutOffLeaveTask) New() (output cutOffLeaveTask) {
	output.RunType = "scheduler.cutOff_annual_leave"
	return
}

func (input cutOffLeaveTask) Start() {
	input.StartTask(input.RunType, input.setCutOffAnnualLeave)
}

func (input cutOffLeaveTask) setCutOffAnnualLeave() {
	input.doSetCutOffAnnualLeave()
}

func (input cutOffLeaveTask) doSetCutOffAnnualLeave() {

	var (
		startTime     = time.Now()
		group         = "Cut Off Annual Leave"
		typeJob       = "CutOffAnnualLeave"
		nameSch       = "scheduler.cutOff_annual_leave"
		jobProcess    repository.JobProcessModel
		err           errorModel.ErrorModel
		jobID         = util.GetUUID()
		counting      = 0
	)

	fmt.Println("----------> [START] Scheduler process Cut Off Annual Leave")

	jobProcess = repository.JobProcessModel{
		Level:     sql.NullInt32{Int32: 0},
		JobID:     sql.NullString{String: jobID},
		Group:     sql.NullString{String: group},
		Type:      sql.NullString{String: typeJob},
		Name:      sql.NullString{String: nameSch},
		Status:    sql.NullString{String: constanta.JobProcessOnProgressStatus},
		CreatedBy: sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt: sql.NullTime{Time: startTime},
		UpdatedAt: sql.NullTime{Time: startTime},
	}

	if err := dao.JobProcessDAO.InsertJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcess); err.Error != nil {
		return
	}

	benefits, _ := dao.EmployeeDAO.GetAllEmployee(serverconfig.ServerAttribute.DBConnection, " INNER JOIN ")

	for i:=0; i<len(benefits);i++ {
		benefitRepo := input.updateRepository(benefits[i])
		dao.EmployeeHistoryLeaveDAO.UpdateCutOffValueHistory(serverconfig.ServerAttribute.DBConnection, benefitRepo)
	}

	dao.EmployeeBenefitsDAO.ResetLastAnnualLeave(serverconfig.ServerAttribute.DBConnection)

	timeFinish := time.Now()

	jobProcess.Status.String = constanta.JobProcessDoneStatus
	jobProcess.UpdatedAt.Time = timeFinish
	jobProcess.Counter.Int32 = int32(counting)
	jobProcess.Total.Int32 = int32(counting)

	if err = dao.JobProcessDAO.UpdateFullJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcess); err.Error != nil {
		return
	}

	fmt.Println("----------> [FINISH] Scheduler process Cut Off Annual Leave")
}

func (input cutOffLeaveTask) updateRepository(inputStruct repository.EmployeeBenefitsModel) repository.EmployeeBenefitsModel{
	time    := time.Now()
	year , _, _ := time.Date()
	return repository.EmployeeBenefitsModel{
		EmployeeID :        sql.NullInt64{Int64: inputStruct.ID.Int64},
		CutOffLeaveValue:   sql.NullInt64{Int64: inputStruct.LastAnnualLeave.Int64},
		NoteCutOff:         sql.NullString{String: "CutOffAnnualLeave done at "+time.Format("02-01-2006")},
		UpdatedAt:          sql.NullTime{Time: time},
		Year:               sql.NullString{String: strconv.Itoa(year-1)},
	}
}