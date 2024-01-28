package TaskSchedulerService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input taskSchedulerService) SchedulerRetryValidateExpirationLicense()  {
	service.LogMessage("Scheduler Retry Check Expiration Product License", 200)
	var jobProcessError, jobProcessModel repository.JobProcessModel
	var err errorModel.ErrorModel
	var totalCount int
	var errorDataModel []repository.ContentDataOutDetail
	contextModel := input.getContextModel()
	funcName := "SchedulerRetryValidateExpirationLicense"
	task := input.getValidateExpiredProductLicenseTask()
	db := serverconfig.ServerAttribute.DBConnection

	defer func() {
		timeNow := time.Now()
		jobProcessModel.UpdatedAt.Time = timeNow
		if err.Error != nil {
			jobProcessModel.Status.String = constanta.JobProcessErrorStatus
			if jobProcessModel.ContentDataOut.String == "" {
				jobProcessModel.ContentDataOut.String = service.GetErrorMessage(err, contextModel)
			}
			err = dao.JobProcessDAO.UpdateErrorJobProcess(db, jobProcessModel)
			if err.Error != nil {
				input.LogError(err, contextModel)
			}
		}else {
			err = dao.JobProcessDAO.UpdateJobProcessCounter(db, jobProcessModel)
			if err.Error != nil {
				input.LogError(err, contextModel)
			}

			if jobProcessModel.Total.Int32 > 0 {
				successLog := fmt.Sprintf("Success do %s data %d  from %d ", jobProcessModel.Name.String, jobProcessModel.Counter.Int32, jobProcessModel.Total.Int32)
				service.LogMessage(successLog, 200 )
			}
		}
	}()

	jobProcessModel = repository.JobProcessModel{
		Group: sql.NullString{String: task.Group},
		Type: sql.NullString{String: task.Type},
		Name: sql.NullString{String: task.Name},
		Status: sql.NullString{String: constanta.JobProcessErrorStatus},
	}

	jobProcessError, err = dao.JobProcessDAO.GetJobProcessError(db, jobProcessModel)
	if err.Error != nil {
		service.LogMessage(service.GetErrorMessage(err, contextModel), err.Code)
		return
	}

	if jobProcessError.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.JobProcess)
		service.LogMessage("No Data Error Job Process", 200)
		return
	}

	jobProcessModel.JobID = jobProcessError.JobID
	jobProcessModel.Status.String = constanta.JobProcessOnProgressStatus
	jobProcessModel.Total = jobProcessError.Total
	jobProcessModel.Counter = jobProcessError.Counter

	totalCount, errorDataModel = input.retrySchedulerCheckExpiration(jobProcessError, contextModel)
	if len(errorDataModel) > 0 {
		jobProcessModel.ContentDataOut.String = util.StructToJSON(errorDataModel)
		jobProcessModel.Status.String = constanta.JobProcessErrorStatus
	}

	jobProcessModel.Counter.Int32 += int32(totalCount)
}

func (input taskSchedulerService) retrySchedulerCheckExpiration(inputStruct repository.JobProcessModel, contextModel applicationModel.ContextModel) (totalCount int, errorDataModel []repository.ContentDataOutDetail) {
	var dataError []repository.ContentDataOutDetail
	var err errorModel.ErrorModel
	var productLicenses []repository.ProductLicenseModel

	funcName := "retrySchedulerCheckExpiration"

	errorS := json.Unmarshal([]byte(inputStruct.ContentDataOut.String), &dataError)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		service.LogMessage(service.GetErrorMessage(err, contextModel), err.Code)
		return
	}

	for _, detail := range dataError {
		productLicenses = append(productLicenses, repository.ProductLicenseModel{
			ID: sql.NullInt64{Int64: detail.ID},
		})
	}

	totalCount, errorDataModel = input.UpdateExpiredProduct(productLicenses, &contextModel)
	return
}

func (input taskSchedulerService) getContextModel() (output applicationModel.ContextModel) {
	output.AuthAccessTokenModel = model.AuthAccessTokenModel{
		RedisAuthAccessTokenModel:  model.RedisAuthAccessTokenModel{
			ResourceUserID: config.ApplicationConfiguration.GetClientCredentialsAuthUserID(),
			Locale:         constanta.DefaultApplicationsLanguage,
		},
		ClientID:                   config.ApplicationConfiguration.GetClientCredentialsClientID(),
		Locale:                     constanta.DefaultApplicationsLanguage,
	}
	return
}
