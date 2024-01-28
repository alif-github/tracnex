package TaskSchedulerService

import (
	"database/sql"
	"fmt"
	"math"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/backgroundJobModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"sync"
	"time"
)

func (input taskSchedulerService) getValidateExpiredProductLicenseTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constanta.JobProcessCheckGroup,
		Type:  constanta.JobProcessExpirationType,
		Name:  constanta.JobProcessCheckExpirationProductLicense,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          nil,
		},
		GetCountData: dao.ProductLicenseDAO.GetCountForCheckExpiration,
		DoJob:        input.checkExpirationProductLicense,
	}
}

func (input taskSchedulerService) SchedulerCheckExpiredProductLicense() {
	var (
		timeNow = time.Now()
		task    = input.getValidateExpiredProductLicenseTask()
		job     = service.GetJobProcess(task, applicationModel.ContextModel{}, timeNow)
		db      = serverconfig.ServerAttribute.DBConnection
	)

	input.ServiceWithBackgroundProcess(db, false, repository.JobProcessModel{}, job, task, applicationModel.ContextModel{})
}

func (input taskSchedulerService) checkExpirationProductLicense(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel) (err errorModel.ErrorModel) {
	var (
		wg                   sync.WaitGroup
		contextModel         applicationModel.ContextModel
		detailContentDataOut []repository.ContentDataOutDetail
		totalThread          = 5
		timeNow              = time.Now()
		timeNowStr           = timeNow.Format(constanta.DefaultTimeFormat)
	)

	service.LogMessage(fmt.Sprintf(`Scheduler Check Expiration Product License [%s]`, timeNowStr), 200)
	contextModel.AuthAccessTokenModel.ClientID = constanta.SystemClient
	contextModel.AuthAccessTokenModel.ResourceUserID = constanta.SystemID

	if childJob.Total.Int32 < 1 {
		service.LogMessage(fmt.Sprintf(`No Data Product License [%s]`, timeNowStr), 200)
		return
	}

	totalPage := int(math.Ceil(float64(childJob.Total.Int32) / float64(constanta.TotalDataProductLicensePerChannel)))
	if totalPage < totalThread {
		totalThread = totalPage
	}

	result := make(chan repository.ProductLicenseResponseForScheduler, childJob.Total.Int32)
	jobs := make(chan in.AbstractDTO, childJob.Total.Int32)

	//--- Do the job with thread pool
	for i := 0; i < totalThread; i++ {
		wg.Add(1)
		go input.doCheckExpiredProductLicense(db, jobs, result, &contextModel, &wg)
	}

	//--- Apply the job
	for i := 0; i < totalPage; i++ {
		var tempAbstractDTO in.AbstractDTO
		tempAbstractDTO.Limit = constanta.TotalDataProductLicensePerChannel
		tempAbstractDTO.Page = i + 1
		jobs <- tempAbstractDTO
	}

	//--- Get the result
	for i := 0; i < totalPage; i++ {
		tempResult := <-result
		childJob.Counter.Int32 += int32(tempResult.TotalData)
		if len(tempResult.ContentDataOutDetail) > 0 {
			detailContentDataOut = append(detailContentDataOut, tempResult.ContentDataOutDetail...)
		}
	}

	close(jobs)
	wg.Wait()

	if len(detailContentDataOut) > 0 {
		childJob.ContentDataOut.String = util.StructToJSON(detailContentDataOut)
		childJob.Status.String = constanta.JobProcessErrorStatus
	}

	return
}

func (input taskSchedulerService) doCheckExpiredProductLicense(db *sql.DB, jobs <-chan in.AbstractDTO, result chan<- repository.ProductLicenseResponseForScheduler, contextModel *applicationModel.ContextModel, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		expProdLicenseOnDB []repository.ProductLicenseModel
		err                errorModel.ErrorModel
	)

	for job := range jobs {
		var response repository.ProductLicenseResponseForScheduler

		//--- Get expired product license
		expProdLicenseOnDB, err = dao.ProductLicenseDAO.GetExpiredProductLicense(db, job)
		if err.Error != nil {
			if err.CausedBy != nil {
				service.LogMessage(err.CausedBy.Error(), err.Code)
			}
			continue
		}

		//--- Update expired product license
		tempTotalCountJob, dataContentDataOutDetail := input.UpdateExpiredProduct(expProdLicenseOnDB, contextModel)

		response.TotalData = int64(tempTotalCountJob)
		response.ContentDataOutDetail = dataContentDataOutDetail
		result <- response
	}
}

func (input taskSchedulerService) UpdateExpiredProduct(expiredProductLicenses []repository.ProductLicenseModel, contextModel *applicationModel.ContextModel) (totalCount int, detailResponse []repository.ContentDataOutDetail) {
	var (
		funcName = "UpdateExpiredProduct"
		err      errorModel.ErrorModel
	)

	totalCount = 0
	for _, product := range expiredProductLicenses {
		_, err = input.ServiceWithDataAuditPreparedByService(funcName, product, contextModel, input.doUpdatedExpiredProduct, func(_ interface{}, _ applicationModel.ContextModel) {})
		if err.Error != nil {
			if err.CausedBy != nil {
				service.LogMessage(err.CausedBy.Error(), err.Code)
			}
			detailResponse = append(detailResponse, repository.ContentDataOutDetail{
				ID:      product.ID.Int64,
				Status:  err.Code,
				Message: service.GetErrorMessage(err, *contextModel),
			})
			continue
		}
		totalCount++
	}
	return
}

func (input taskSchedulerService) doUpdatedExpiredProduct(tx *sql.Tx, inputStruct interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName           = "doUpdatedExpiredProduct"
		expiredProduct     = inputStruct.(repository.ProductLicenseModel)
		db                 = serverconfig.ServerAttribute.DBConnection
		expiredUser        []repository.UserRegistrationDetailModel
		expiredUserLicense []repository.UserLicenseModel
		tempDataAudit      []repository.AuditSystemModel
		userCount          int
	)

	inputModel := repository.ProductLicenseModel{
		ID:            expiredProduct.ID,
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		LicenseStatus: sql.NullInt32{Int32: constanta.ProductLicenseStatusExpired},
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ProductLicenseDAO.TableName, inputModel.ID.Int64, 0)...)
	err = dao.ProductLicenseDAO.UpdatedForExpiredProduct(tx, inputModel)
	if err.Error != nil {
		return
	}

	userCount, err = dao.UserRegistrationDetailDAO.GetCountUserRegisDetailForCheckExpiredLicense(db, expiredProduct.ID.Int64)
	if err.Error != nil {
		return
	}

	if userCount > 0 {
		expiredUser, err = dao.UserRegistrationDetailDAO.GetUserRegistrationDetailForCheckExpiredLicense(db, expiredProduct.ID.Int64)
		if err.Error != nil {
			return
		}

		if len(expiredUser) < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserLicense)
			return
		}

		expiredUserLicense, err = dao.UserLicenseDAO.GetUserLicenseForCheckExpiredLicense(db, expiredProduct.ID.Int64)
		if err.Error != nil {
			return
		}

		if len(expiredUserLicense) < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserLicense)
			return
		}

		_, tempDataAudit, err = input.updateExpiredUserLicense(tx, expiredUser, contextModel, timeNow)
		dataAudit = append(dataAudit, tempDataAudit...)
		if err.Error != nil {
			return
		}

		_, tempDataAudit, err = input.resetTotalActivatedUserLicense(tx, expiredUserLicense, contextModel, timeNow)
		dataAudit = append(dataAudit, tempDataAudit...)
		if err.Error != nil {
			return
		}
	}

	return
}

func (input taskSchedulerService) updateExpiredUserLicense(tx *sql.Tx, inputStruct []repository.UserRegistrationDetailModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	for i := 0; i < len(inputStruct); i++ {
		var tempDataAudit []repository.AuditSystemModel
		_, tempDataAudit, err = input.doUpdateExpiredUserLicense(tx, inputStruct[i], contextModel, timeNow)
		dataAudit = append(dataAudit, tempDataAudit...)
		if err.Error != nil {
			return
		}
	}

	return
}

func (input taskSchedulerService) resetTotalActivatedUserLicense(tx *sql.Tx, inputStruct []repository.UserLicenseModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	for i := 0; i < len(inputStruct); i++ {
		var tempDataAudit []repository.AuditSystemModel
		_, tempDataAudit, err = input.doResetTotalActivatedUserLicense(tx, inputStruct[i], contextModel, timeNow)
		dataAudit = append(dataAudit, tempDataAudit...)
		if err.Error != nil {
			return
		}
	}

	return
}

func (input taskSchedulerService) doUpdateExpiredUserLicense(tx *sql.Tx, inputStruct interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		expiredUser  = inputStruct.(repository.UserRegistrationDetailModel)
		expiredModel repository.UserRegistrationDetailModel
	)

	expiredModel = repository.UserRegistrationDetailModel{
		ID:            expiredUser.ID,
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		Status:        sql.NullString{String: constanta.StatusNonActive},
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserRegistrationDetailDAO.TableName, expiredModel.ID.Int64, 0)...)
	err = dao.UserRegistrationDetailDAO.UpdateStatusUserRegistrationDetail(tx, expiredModel)
	return
}

func (input taskSchedulerService) doResetTotalActivatedUserLicense(tx *sql.Tx, inputStruct interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		expiredUser  = inputStruct.(repository.UserLicenseModel)
		expiredModel repository.UserLicenseModel
	)

	expiredModel = repository.UserLicenseModel{
		ID:             expiredUser.ID,
		UpdatedAt:      sql.NullTime{Time: timeNow},
		UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		TotalActivated: sql.NullInt64{Int64: 0},
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserLicenseDAO.TableName, expiredModel.ID.Int64, 0)...)
	err = dao.UserLicenseDAO.ResetTotalActivatedUserLicense(tx, expiredModel)
	return
}
