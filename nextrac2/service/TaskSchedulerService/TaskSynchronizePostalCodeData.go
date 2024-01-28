package TaskSchedulerService

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/backgroundJobModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"sync"
	"time"
)

func (input taskSchedulerService) SchedulerSyncPostalCodeData(parentJob *repository.JobProcessModel) {
	var (
		timeNow      = time.Now()
		task         = input.getSynchronizePostalCodeDataTask()
		job          = service.GetJobProcess(task, applicationModel.ContextModel{}, timeNow)
		db           = serverconfig.ServerAttribute.DBConnection
		contextModel applicationModel.ContextModel
	)

	input.ServiceWithBackgroundProcess(db, false, *parentJob, job, task, contextModel)
	parentJob.Counter.Int32 += 1
}

func (input taskSchedulerService) getSynchronizePostalCodeDataTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group:        constanta.JobProcessSynchronizeGroup,
		Type:         constanta.JobProcessMasterDataType,
		Name:         constanta.JobProcessSynchronizePostalCode,
		Data:         backgroundJobModel.BackgroundServiceModel{},
		GetCountData: input.getCountAllPostalCodeFromMasterData,
		DoJobWithCtx: input.synchronizePostalCodeFromMasterData,
	}
}

func (input taskSchedulerService) synchronizePostalCodeFromMasterData(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel, ctx applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		postalCodeOnMDB []master_data_response.PostalCodeResponse
		resultOnDB      repository.PostalCodeModel
		counter         interface{}
	)

	if childJob.Total.Int32 < 1 {
		service.LogMessage("[POSTAL-CODE] No data is processed", 200)
		return
	}

	//--- Get Date Last Sync
	resultOnDB, err = dao.PostalCodeDAO.GetDateLastSyncPostalCode(db)
	if err.Error != nil {
		return
	}

	//--- Get Postal Code On MDB
	postalCodeOnMDB, err = input.getPostalCodeMDB(resultOnDB)
	if err.Error != nil {
		return
	}

	//--- Update Insert Postal Code
	counter, err = input.ServiceWithDataAuditPreparedByService("synchronizePostalCodeFromMasterData", postalCodeOnMDB, &ctx, input.doSyncPostalCode, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	if counter != nil {
		childJob.Counter.Int32 += counter.(int32)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input taskSchedulerService) getCountAllPostalCodeFromMasterData(db *sql.DB, _ []in.SearchByParam, _ bool, _ int64) (result int, err errorModel.ErrorModel) {
	var (
		tempResult int64
		resultOnDB repository.PostalCodeModel
	)

	//--- Get Date Last Sync
	resultOnDB, err = dao.PostalCodeDAO.GetDateLastSyncPostalCode(db)
	if err.Error != nil {
		return
	}

	//--- Get Count Postal Code MDB
	tempResult, err = master_data_dao.CountAllPostalCodeFromMasterData(resultOnDB, &applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	result = int(tempResult)
	return
}

func (input taskSchedulerService) getPostalCodeMDB(resultOnDB repository.PostalCodeModel) (postalCodeOnMDB []master_data_response.PostalCodeResponse, err errorModel.ErrorModel) {
	request := master_data_request.PostalCodeRequest{
		UpdatedAtStart: resultOnDB.LastSync.Time,
		AbstractDTO: in.AbstractDTO{
			Page: -99,
		},
	}

	postalCodeOnMDB, err = master_data_dao.GetListForSyncPostalCodeFromMDB(request, &applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input taskSchedulerService) doSyncPostalCode(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		listPostalCodeMDB    []repository.PostalCodeModel
		listPostalCodeDB     []repository.PostalCodeModel
		listInsertPostalCode []repository.PostalCodeModel
		listUpdatePostalCode []repository.PostalCodeModel
		isInsertedData       bool
		postalCodeMDB        = inputStructInterface.([]master_data_response.PostalCodeResponse)
		counter              int32
	)

	//--- List Postal Code
	for _, postalCode := range postalCodeMDB {
		listPostalCodeMDB = append(listPostalCodeMDB, repository.PostalCodeModel{
			UrbanVillageID:   sql.NullInt64{Int64: postalCode.UrbanVillageID},
			UrbanVillageName: sql.NullString{String: postalCode.UrbanVillageName},
			MDBPostalCodeID:  sql.NullInt64{Int64: postalCode.ID},
			Code:             sql.NullString{String: postalCode.Code},
			Status:           sql.NullString{String: postalCode.Status},
			CreatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			LastSync:         sql.NullTime{Time: timeNow},
		})
	}

	tempData := input.poolGetDataPostalCode(listPostalCodeMDB)

	for _, data := range tempData {
		err = data.Err
		if err.Error != nil {
			return
		}
		listPostalCodeDB = append(listPostalCodeDB, data.Data.([]repository.PostalCodeModel)...)
	}

	for i := 0; i < len(listPostalCodeMDB); i++ {
		isInsertedData = true
		for j := 0; j < len(listPostalCodeDB); j++ {
			if listPostalCodeMDB[i].MDBPostalCodeID.Int64 == listPostalCodeDB[j].MDBPostalCodeID.Int64 {
				isInsertedData = false
				listPostalCodeMDB[i].ID = listPostalCodeDB[j].ID
				listUpdatePostalCode = append(listUpdatePostalCode, listPostalCodeMDB[i])
				if j < len(listPostalCodeDB)-1 {
					listPostalCodeDB = append(listPostalCodeDB[:j], listPostalCodeDB[j+1:]...)
				} else {
					listPostalCodeDB = append(listPostalCodeDB[:j])
				}
				break
			}
		}

		if isInsertedData {
			listInsertPostalCode = append(listInsertPostalCode, listPostalCodeMDB[i])
		}
	}

	//--- Do Update Postal Code
	for _, postalCodeData := range listUpdatePostalCode {
		err = dao.PostalCodeDAO.UpdateDataPostalCode(tx, postalCodeData)
		if err.Error != nil {
			return
		}
		counter += 1
	}

	//--- Do Insert Postal Code
	//if len(listInsertPostalCode) > 0 {
	//	fmt.Println(len(listInsertPostalCode))
	//	limitData := 50
	//	totalPage := int(math.Ceil(float64(len(listInsertPostalCode)) / float64(limitData)))
	//	for i := 0; i < totalPage; i++ {
	//		var insertedID []int64
	//		if i == 0 {
	//			insertedID, err = dao.PostalCodeDAO.InsertBulkPostalCode(tx, listInsertPostalCode[:limitData])
	//		} else if i == totalPage-1 {
	//			insertedID, err = dao.PostalCodeDAO.InsertBulkPostalCode(tx, listInsertPostalCode[limitData*i:])
	//		} else {
	//			insertedID, err = dao.PostalCodeDAO.InsertBulkPostalCode(tx, listInsertPostalCode[limitData*i:limitData*(i+1)])
	//		}
	//
	//		if err.Error != nil {
	//			return
	//		}
	//
	//		counter += int32(len(insertedID))
	//	}
	//}

	//--- Insert Postal Code
	if len(listInsertPostalCode) > 0 {
		var insertedID []int64
		insertedID, err = dao.PostalCodeDAO.InsertBulkPostalCode(tx, listInsertPostalCode)
		if err.Error != nil {
			return
		}

		counter += int32(len(insertedID))
	}

	//--- Update Last Sync
	//_ = dao.PostalCodeDAO.UpdateLastSyncPostalCode(tx, timeNow)
	//err = errorModel.GenerateNonErrorModel()

	defer func() {
		output = counter
	}()

	log.Println("End job process Postal Code")
	return
}

func (input taskSchedulerService) insertBulkPostalCode(tx *sql.Tx, insertedPostalCode []repository.PostalCodeModel) (tempCounter int32, err errorModel.ErrorModel) {
	var (
		totalJob    = 1000
		numJobs     = int(math.Ceil(float64(len(insertedPostalCode)) / float64(totalJob)))
		poolSize    = 5
		job         = make(chan repository.PostalCodeModel, len(insertedPostalCode))
		results     = make(chan responseDataForThread, len(insertedPostalCode))
		ctx, cancel = context.WithCancel(context.Background())
		wg          sync.WaitGroup
	)

	if poolSize > numJobs {
		poolSize = numJobs
	}

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go input.doInsertBulkPostalCode(tx, &wg, ctx, job, results)
	}

	for j := 0; j < len(insertedPostalCode); j++ {
		job <- insertedPostalCode[j]
	}

	close(job)
	wg.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		for _ = range insertedPostalCode {
			select {
			case <-ctx.Done():
				return
			case r := <-results:
				tempCounter++
				if r.Err.Error != nil {
					err = r.Err
					cancel()
				}
			}
		}
	}(&wg)

	wg.Wait()
	ctx.Done()
	close(results)
	log.Println(fmt.Sprintf(`[POSTAL CODE] Data insert raw: %d <===> Data insert berhasil di proses: %d`, len(insertedPostalCode), tempCounter))
	return
}

func (input taskSchedulerService) doInsertBulkPostalCode(tx *sql.Tx, wg *sync.WaitGroup, ctx context.Context, job <-chan repository.PostalCodeModel, result chan<- responseDataForThread) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-job:
			var (
				err              errorModel.ErrorModel
				urbanVillage     repository.UrbanVillageModel
				tempResponseData responseDataForThread
				db               = serverconfig.ServerAttribute.DBConnection
			)

			if !ok {
				return
			}

			urbanVillage, err = dao.UrbanVillageDAO.GetUrbanVillageByID(db, task.UrbanVillageID.Int64, true)
			if err.Error != nil {
				log.Println("Error: =====> ", err.Error.Error())
				tempResponseData = responseDataForThread{
					Err: err,
				}
			}

			if urbanVillage.ID.Int64 > 0 {
				_, err = dao.PostalCodeDAO.InsertPostalCode(tx, task)
				if err.Error != nil {
					log.Println("Error: ====> ", err.Error.Error())
				}
				tempResponseData = responseDataForThread{
					Data: task.Code.String,
					Err:  err,
				}
			} else {
				tempResponseData = responseDataForThread{
					Data: task.Code.String,
					Err:  errorModel.GenerateNonErrorModel(),
				}
			}

			result <- tempResponseData
		}
	}
}

func (input taskSchedulerService) poolGetDataPostalCode(modelData []repository.PostalCodeModel) (result []responseDataForThread) {
	var (
		totalJob    = 1000
		numJobs     = int(math.Ceil(float64(len(modelData)) / float64(totalJob)))
		poolSize    = 5
		jobs        = make(chan []repository.PostalCodeModel, len(modelData))
		results     = make(chan responseDataForThread, len(modelData))
		ctx, cancel = context.WithCancel(context.Background())
		wg, wg2     sync.WaitGroup
	)

	if poolSize > numJobs {
		poolSize = numJobs
	}

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		// Consume Jobs
		go input.doGetDataPostalCodeForSync(&wg, jobs, results, ctx)
	}

	// Produce Jobs
	for i := 0; i < numJobs; i++ {
		if i == numJobs-1 {
			jobs <- modelData[totalJob*i:]
		} else if i == 0 {
			jobs <- modelData[:totalJob]
		} else {
			jobs <- modelData[totalJob*i : totalJob*(i+1)]
		}
	}

	close(jobs)

	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for res := range results {
			result = append(result, res)
			if res.Err.Error != nil {
				cancel()
			}
		}
	}()

	wg.Wait()
	ctx.Done()
	close(results)
	wg2.Wait()
	return
}

func (input taskSchedulerService) doGetDataPostalCodeForSync(wg *sync.WaitGroup, jobs <-chan []repository.PostalCodeModel,
	result chan<- responseDataForThread, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-jobs:
			if !ok {
				return
			}
			tempData, err := dao.PostalCodeDAO.GetUpdatedDBPostalCode(serverconfig.ServerAttribute.DBConnection, task)
			tempResponseData := responseDataForThread{
				Data: tempData,
				Err:  err,
			}

			result <- tempResponseData
		}
	}
}
