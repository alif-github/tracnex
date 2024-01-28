package Task

import (
	"database/sql"
	"fmt"
	"math"
	"nexsoft.co.id/nexcommon/util"

	"strconv"

	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"time"
)

type historyEmployeeBenefits struct {
	AbstractScheduledTask
}

var HistoryEmployeeBenefits = historyEmployeeBenefits{}.New()

func (input historyEmployeeBenefits) New() (output historyEmployeeBenefits) {
	output.RunType = "scheduler.annual_leave_and_medical"
	return
}

func (input historyEmployeeBenefits) Start() {
	input.StartTask(input.RunType, input.setHistoryBenefit)
}

func (input historyEmployeeBenefits) setHistoryBenefit() {
	input.doSetHistoryEmployeeBenefits()
}

func (input historyEmployeeBenefits) doSetHistoryEmployeeBenefits() {

	var (
		startTime     = time.Now()
		group         = "History Employee Leave Benefits"
		typeJob       = "History Leave Benefits"
		nameSch       = "Scheduler_History_Employee_Benefits"
		jobProcess    repository.JobProcessModel
		err           errorModel.ErrorModel
		jobID         = util.GetUUID()
		counting      = 0
	)

	year , _, _ := startTime.Date()
	countYear, _ := dao.EmployeeHistoryLeaveDAO.GetCountEmployeeByYear(serverconfig.ServerAttribute.DBConnection, strconv.Itoa(year-1))

	fmt.Println("----------> [START] Scheduler process History Employee Benefits")

	benefits, _ := dao.EmployeeDAO.GetAllEmployee(serverconfig.ServerAttribute.DBConnection, " INNER JOIN ")

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

	if countYear == 0 {
		for i:=0; i<len(benefits);i++  {
			benefits[i].Year.String = strconv.Itoa(year-1)
			dao.EmployeeHistoryLeaveDAO.InsertHistoryBenefits(serverconfig.ServerAttribute.DBConnection, benefits[i])
		}
	}

	for i:=0; i<len(benefits);i++  {
		benefitRepo := input.updateRepository(benefits[i])
		_ = dao.EmployeeBenefitsDAO.UpdateLeaveAndMedicalValueDB(serverconfig.ServerAttribute.DBConnection, benefitRepo)
	}

	timeFinish := time.Now()

	jobProcess.Status.String = constanta.JobProcessDoneStatus
	jobProcess.UpdatedAt.Time = timeFinish
	jobProcess.Counter.Int32 = int32(counting)
	jobProcess.Total.Int32 = int32(counting)

	if err = dao.JobProcessDAO.UpdateFullJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcess); err.Error != nil {
		return
	}

	fmt.Println("----------> [FINISH] Scheduler process History Employee Benefits")
}

func (input historyEmployeeBenefits) updateRepository(inputStruct repository.EmployeeBenefitsModel) repository.EmployeeBenefitsModel{
	return repository.EmployeeBenefitsModel{
		EmployeeID :        sql.NullInt64{Int64: inputStruct.ID.Int64},
		CurrentAnnualLeave: sql.NullInt64{Int64: input.getLeaveValue(inputStruct.JoinDate.Time)},
		LastAnnualLeave:    sql.NullInt64{Int64 : inputStruct.CurrentAnnualLeave.Int64},
		CurrentMedicalValue:sql.NullFloat64{Float64:5000000},
		LastMedicalValue:   sql.NullFloat64{Float64:inputStruct.CurrentMedicalValue.Float64},
		UpdatedAt:          sql.NullTime{Time: time.Now()},
	}
}

func (input historyEmployeeBenefits) getLeaveValue(date time.Time)int64{
	currentTime := time.Now()
	yy, mm, dd := date.Date()
	oldTime := time.Date(yy, mm, dd, 0, 0, 0, 0, time.UTC)
	diff := currentTime.Sub(oldTime)
	month := int64(math.Ceil(diff.Hours() / 24 / 30))
	leave := month - 3

	if leave >= 12 {
		leave = 12
	}else if leave < 0 {
		leave = 0
	}

	return leave
}