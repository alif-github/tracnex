package ReportAnnualLeaveService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

type reportAnnualLeaveService struct {
	service.AbstractService
	service.GetListData
}

var ReportAnnualLeaveService = reportAnnualLeaveService{}.New()

func (input reportAnnualLeaveService) New() (output reportAnnualLeaveService) {
	output.FileName = "ReportAnnualLeaveService.go"
	output.ServiceName = "Report Annual Leave"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{
		"job_id",
		"category",
		"created_at",
	}
	output.ValidOrderBy = []string{
		"created_at",
	}
	return
}

func (input reportAnnualLeaveService) ReportAnnualLeaveService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.ReportAnnualLeave
		timeNow     = time.Now()
		job         *repository.JobProcessModel
	)

	//--- Read Request
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateReport)
	if err.Error != nil {
		return
	}

	//--- If Not Detail Then Return
	if !inputStruct.IsDetail {
		return input.reportSummaryAnnualLeaveService(strconv.Itoa(int(inputStruct.Year[1])), job, contextModel, timeNow)
	}

	job = &repository.JobProcessModel{
		JobID:         sql.NullString{String: util.GetUUID()},
		Status:        sql.NullString{String: constanta.JobProcessOnProgressStatus},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	//--- Do Job
	go input.doServiceJobProcessAnnualLeave(inputStruct, contextModel, timeNow, job)

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

func (input reportAnnualLeaveService) reportAnnualLeaveService(inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time, jobTask *repository.JobProcessModel) (err errorModel.ErrorModel) {
	var (
		fileName         = input.FileName
		funcName         = "reportAnnualLeaveService"
		datasStartYear   []repository.EmployeeLeaveReportModel
		datasEndYear     []repository.EmployeeLeaveReportModel
		fileUploadID     int64
		resultJobProcess repository.ViewJobProcessModel
		inputStruct      = inputStructInterface.(in.ReportAnnualLeave)
		startYear        = int(inputStruct.Year[0])
		endYear          = int(inputStruct.Year[1])
		db               = serverconfig.ServerAttribute.DBConnection
	)

	//--- Main Process (Start Year)
	datasStartYear, err = input.doGetReportAnnualLeaveInYear(startYear)
	if err.Error != nil {
		return
	}

	//--- Start Year Counter
	jobTask.Counter.Int32 += 30

	//--- Main Process (End Year)
	datasEndYear, err = input.doGetReportAnnualLeaveInYear(endYear)
	if err.Error != nil {
		return
	}

	//--- End Year Counter
	jobTask.Counter.Int32 += 30

	//--- Excel Generate
	fileUploadID, err = input.excelProduce(startYear, endYear, datasStartYear, datasEndYear, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//--- Excel Producer Counter
	jobTask.Counter.Int32 += 30

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
		Category:      sql.NullString{String: "Report-Detail"},
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

func (input reportAnnualLeaveService) mainProcessReportAnnualLeave(dateYear time.Time) (datas []repository.EmployeeLeaveReportModel, err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "mainProcessReportAnnualLeave"
		db       = serverconfig.ServerAttribute.DBConnection
		month    = dateYear.Month()
		year     = dateYear.Year()
	)

	//--- Data Collect
	datas, err = dao.EmployeeLeaveDAO.GetReportAnnualLeave(db, dateYear)
	if err.Error != nil {
		return
	}

	//--- If Employee Empty, Return Error
	if len(datas) < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.EmployeeConstanta)
		return
	}

	for i, data := range datas {
		var (
			errs                 error
			detailLeave          []in.LeaveDetail
			detailModelLeaveTemp []repository.DateDetailReportModel
		)

		//--- Fill Month And Year
		datas[i].MonthLeave.Int64 = int64(month)
		datas[i].YearLeave.Int64 = int64(year)

		//--- Null Then Continue
		if data.DetailLeave.String == "[null]" {
			continue
		}

		//--- Unmarshal
		if errs = json.Unmarshal([]byte(data.DetailLeave.String), &detailLeave); errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}

		//--- Date Clean
		detailModelLeaveTemp, err = input.dateProcessCleansing(detailLeave, month, year)
		if err.Error != nil {
			return
		}

		//--- Inject Data
		datas[i].DetailLeaveList = detailModelLeaveTemp
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) dateProcessCleansing(datas []in.LeaveDetail, month time.Month, year int) (result []repository.DateDetailReportModel, err errorModel.ErrorModel) {
	//--- Date Clean
	for _, date := range datas {
		//--- Initiate
		detailLeaveTemp := repository.DateDetailReportModel{
			Type:        sql.NullString{String: date.Type},
			Description: sql.NullString{String: date.Description},
		}

		//--- Do Date Cleansing
		if err = input.doDateProcessCleansing(date.DateStr, &detailLeaveTemp, month, year); err.Error != nil {
			return
		}

		//--- Do Append
		result = append(result, detailLeaveTemp)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) doDateProcessCleansing(datas []string, detailLeaveTemp *repository.DateDetailReportModel, month time.Month, year int) (err errorModel.ErrorModel) {
	var (
		fileName    = input.FileName
		funcName    = "doDateProcessCleansing"
		dateColTemp []sql.NullTime
	)

	//--- Date Clean
	for _, date := range datas {
		//--- Time Parse
		timestamp, errs := time.Parse(constanta.DefaultTimeFormat, date)
		if errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}

		//--- Time Append
		if timestamp.Month() == month && timestamp.Year() == year {
			dateColTemp = append(dateColTemp, sql.NullTime{Time: timestamp})
		}
	}

	detailLeaveTemp.Date = dateColTemp
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ReportAnnualLeave) errorModel.ErrorModel) (inputStruct in.ReportAnnualLeave, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(stringBody) {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}

func (input reportAnnualLeaveService) validateReport(inputStruct *in.ReportAnnualLeave) errorModel.ErrorModel {
	return inputStruct.ValidateReportAnnualLeave()
}
