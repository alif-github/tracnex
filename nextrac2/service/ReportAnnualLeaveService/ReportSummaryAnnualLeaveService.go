package ReportAnnualLeaveService

import (
	"database/sql"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/EmployeeService"
	"time"
)

func (input reportAnnualLeaveService) reportSummaryAnnualLeaveService(year string, job *repository.JobProcessModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.GetListDataDTO
		userStruct  = make(map[string]interface{})
	)

	job = &repository.JobProcessModel{
		JobID:         sql.NullString{String: util.GetUUID()},
		Status:        sql.NullString{String: constanta.JobProcessOnProgressStatus},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	userStruct["struct"] = inputStruct
	userStruct["year"] = year

	//--- Do Job
	go input.doServiceJobProcessSummaryAnnualLeave(userStruct, contextModel, timeNow, job)

	//--- Job ID
	if job.JobID.String != "" {
		outputTemp := make(map[string]interface{})
		outputTemp["job_id"] = job.JobID.String
		output.Data.Content = outputTemp
	}

	output.Status = input.GetResponseMessage(constanta.CommonSuccessGetViewMessages, contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) doReportSummaryAnnualLeaveService(inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time, jobTask *repository.JobProcessModel) (err errorModel.ErrorModel) {
	var (
		fileName         = "ReportSummaryAnnualLeaveService.go"
		funcName         = "doReportSummaryAnnualLeaveService"
		inputStruct      = inputStructInterface.(map[string]interface{})
		userStruct       in.GetListDataDTO
		searchByParam    []in.SearchByParam
		fileUploadID     int64
		resultJobProcess repository.ViewJobProcessModel
		key, keyword     string
		year             string
		content          interface{}
		db               = serverconfig.ServerAttribute.DBConnection
	)

	//--- Convert
	userStruct = inputStruct["struct"].(in.GetListDataDTO)
	year = inputStruct["year"].(string)

	//--- Start Year Counter
	jobTask.Counter.Int32 += 40

	content, err = EmployeeService.EmployeeService.DoGetListEmployeeLeaveYearly(userStruct, searchByParam, key, keyword, year)
	if err.Error != nil {
		return
	}

	//--- End Year Counter
	jobTask.Counter.Int32 += 40

	//--- Excel Generate
	fileUploadID, err = input.excelSummaryProduce(content.([]out.EmployeeLeaveYearly), contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//--- Excel Producer Counter
	jobTask.Counter.Int32 += 10

	//--- Get Job Process
	resultJobProcess, err = dao.JobProcessDAO.ViewJobProcess(db, repository.JobProcessModel{JobID: sql.NullString{String: jobTask.JobID.String}})
	if err.Error != nil {
		return
	}

	tx, errs := db.Begin()
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	//--- Update File Upload
	if err = dao.FileUploadDAO.UpdateFileUpload(tx, repository.FileUpload{
		ID:            sql.NullInt64{Int64: fileUploadID},
		Category:      sql.NullString{String: "Report-Summary"},
		Konektor:      sql.NullString{String: dao.JobProcessDAO.TableName},
		ParentID:      sql.NullInt64{Int64: resultJobProcess.ID.Int64},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		if errs = tx.Rollback(); errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}
		return
	} else {
		if errs = tx.Commit(); errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}
	}

	//--- Finish Counter
	jobTask.Counter.Int32 += 10
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) doServiceJobProcessSummaryAnnualLeave(inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time, job *repository.JobProcessModel) {
	var (
		percent = int64(100)
		db      = serverconfig.ServerAttribute.DBConnection
		task    = taskChildAnnualLeave{
			Group:        "Leave-Summary",
			Type:         "File",
			Name:         "Upload File Summary Annual Leave",
			GetCountData: percent,
			DoJob:        input.doReportSummaryAnnualLeaveService,
		}
	)

	//--- Add Grouping Name
	job.Group.String = task.Group
	job.Type.String = task.Type
	job.Name.String = task.Name

	//--- Job Process
	go input.serviceJobProcessAnnualLeave(inputStructInterface, db, *job, task, *contextModel, timeNow)
}
