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

func (input taskSchedulerService) SchedulerSyncUrbanVillageData(parentJob *repository.JobProcessModel) {
	var (
		timeNow      = time.Now()
		task         = input.getSynchronizeUrbanVillageDataTask()
		job          = service.GetJobProcess(task, applicationModel.ContextModel{}, timeNow)
		db           = serverconfig.ServerAttribute.DBConnection
		contextModel applicationModel.ContextModel
	)

	input.ServiceWithBackgroundProcess(db, false, *parentJob, job, task, contextModel)
	parentJob.Counter.Int32 += 1
}

func (input taskSchedulerService) getSynchronizeUrbanVillageDataTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group:        constanta.JobProcessSynchronizeGroup,
		Type:         constanta.JobProcessMasterDataType,
		Name:         constanta.JobProcessSynchronizeUrbanVillage,
		Data:         backgroundJobModel.BackgroundServiceModel{},
		GetCountData: input.getCountAllUrbanVillageFromMasterData,
		DoJobWithCtx: input.synchronizeUrbanVillageFromMasterData,
	}
}

func (input taskSchedulerService) synchronizeUrbanVillageFromMasterData(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel, ctx applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		urbanVillageOnMDB []master_data_response.UrbanVillageResponse
		resultOnDB        repository.UrbanVillageModel
		counter           interface{}
	)

	if childJob.Total.Int32 < 1 {
		service.LogMessage("[URBAN-VILLAGE] No data is processed", 200)
		return
	}

	//--- Get Date Last Sync
	resultOnDB, err = dao.UrbanVillageDAO.GetDateLastSyncUrbanVillage(db)
	if err.Error != nil {
		return
	}

	//--- Get Postal Code On MDB
	urbanVillageOnMDB, err = input.getUrbanVillageMDB(resultOnDB)
	if err.Error != nil {
		return
	}

	//--- Update Insert Postal Code
	counter, err = input.ServiceWithDataAuditPreparedByService("synchronizeUrbanVillageFromMasterData", urbanVillageOnMDB, &ctx, input.doSyncUrbanVillage, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	if counter != nil {
		childJob.Counter.Int32 += counter.(int32)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input taskSchedulerService) getCountAllUrbanVillageFromMasterData(db *sql.DB, _ []in.SearchByParam, _ bool, _ int64) (result int, err errorModel.ErrorModel) {
	var (
		tempResult int64
		resultOnDB repository.UrbanVillageModel
	)

	//--- Get Date Last Sync
	resultOnDB, err = dao.UrbanVillageDAO.GetDateLastSyncUrbanVillage(db)
	if err.Error != nil {
		return
	}

	//--- Get Count Postal Code MDB
	tempResult, err = master_data_dao.CountAllUrbanVillageFromMasterData(resultOnDB, &applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	result = int(tempResult)
	return
}

func (input taskSchedulerService) getUrbanVillageMDB(resultOnDB repository.UrbanVillageModel) (urbanVillageOnMDB []master_data_response.UrbanVillageResponse, err errorModel.ErrorModel) {
	request := master_data_request.UrbanVillageRequest{
		UpdatedAtStart: resultOnDB.LastSync.Time,
		AbstractDTO: in.AbstractDTO{
			Page: -99,
		},
	}

	urbanVillageOnMDB, err = master_data_dao.GetListForSyncUrbanVillageFromMDB(request, &applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input taskSchedulerService) doSyncUrbanVillage(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		listUrbanVillageMDB    []repository.UrbanVillageModel
		listUrbanVillageDB     []repository.UrbanVillageModel
		listInsertUrbanVillage []repository.UrbanVillageModel
		listUpdateUrbanVillage []repository.UrbanVillageModel
		isInsertedData         bool
		urbanVillageMDB        = inputStructInterface.([]master_data_response.UrbanVillageResponse)
		counter                int32
	)

	//--- List Urban Village
	for _, urbanVillage := range urbanVillageMDB {
		listUrbanVillageMDB = append(listUrbanVillageMDB, repository.UrbanVillageModel{
			SubDistrictID:     sql.NullInt64{Int64: urbanVillage.SubDistrictID},
			SubDistrictName:   sql.NullString{String: urbanVillage.SubDistrictName},
			MDBUrbanVillageID: sql.NullInt64{Int64: urbanVillage.ID},
			Code:              sql.NullString{String: urbanVillage.Code},
			Name:              sql.NullString{String: urbanVillage.Name},
			Status:            sql.NullString{String: urbanVillage.Status},
			CreatedAt:         sql.NullTime{Time: timeNow},
			CreatedBy:         sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:     sql.NullString{String: constanta.SystemClient},
			UpdatedAt:         sql.NullTime{Time: timeNow},
			UpdatedBy:         sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:     sql.NullString{String: constanta.SystemClient},
			LastSync:          sql.NullTime{Time: timeNow},
		})
	}

	tempData := input.poolGetDataUrbanVillage(listUrbanVillageMDB)

	for _, data := range tempData {
		err = data.Err
		if err.Error != nil {
			return
		}
		listUrbanVillageDB = append(listUrbanVillageDB, data.Data.([]repository.UrbanVillageModel)...)
	}

	for i := 0; i < len(listUrbanVillageMDB); i++ {
		isInsertedData = true
		for j := 0; j < len(listUrbanVillageDB); j++ {
			if listUrbanVillageMDB[i].MDBUrbanVillageID.Int64 == listUrbanVillageDB[j].MDBUrbanVillageID.Int64 {
				isInsertedData = false
				listUrbanVillageMDB[i].ID = listUrbanVillageDB[j].ID
				listUpdateUrbanVillage = append(listUpdateUrbanVillage, listUrbanVillageMDB[i])
				if j < len(listUrbanVillageDB)-1 {
					listUrbanVillageDB = append(listUrbanVillageDB[:j], listUrbanVillageDB[j+1:]...)
				} else {
					listUrbanVillageDB = append(listUrbanVillageDB[:j])
				}
				break
			}
		}

		if isInsertedData {
			listInsertUrbanVillage = append(listInsertUrbanVillage, listUrbanVillageMDB[i])
		}
	}

	//--- Do Update Urban Village
	for _, urbanVillageData := range listUpdateUrbanVillage {
		err = dao.UrbanVillageDAO.UpdateDataUrbanVillage(tx, urbanVillageData)
		if err.Error != nil {
			return
		}

		counter += 1
	}

	//--- Do Insert Urban Village
	//if len(listInsertUrbanVillage) > 0 {
	//	limitData := constanta.TotalDataProductLicensePerChannel
	//	totalPage := int(math.Ceil(float64(len(listInsertUrbanVillage)) / float64(limitData)))
	//	for i := 0; i < totalPage; i++ {
	//		var insertedID []int64
	//		if i == 0 {
	//			insertedID, err = dao.UrbanVillageDAO.InsertBulkUrbanVillage(tx, listInsertUrbanVillage[:limitData])
	//		} else if i == totalPage-1 {
	//			insertedID, err = dao.UrbanVillageDAO.InsertBulkUrbanVillage(tx, listInsertUrbanVillage[limitData*i:])
	//		} else {
	//			insertedID, err = dao.UrbanVillageDAO.InsertBulkUrbanVillage(tx, listInsertUrbanVillage[limitData*i:limitData*(i+1)])
	//		}
	//
	//		if err.Error != nil {
	//			return
	//		}
	//		counter += int32(len(insertedID))
	//	}
	//}

	//--- Do Insert Urban Village
	if len(listInsertUrbanVillage) > 0 {
		var insertedID []int64
		insertedID, err = dao.UrbanVillageDAO.InsertBulkUrbanVillage(tx, listInsertUrbanVillage)
		counter += int32(len(insertedID))
		if err.Error != nil {
			return
		}
	}

	//--- Update Last Sync
	//_ = dao.UrbanVillageDAO.UpdateLastSyncUrbanVillage(tx, timeNow)
	//err = errorModel.GenerateNonErrorModel()

	defer func() {
		output = counter
	}()

	log.Println("End job process Urban Village")
	return
}

func (input taskSchedulerService) insertBulkUrbanVillage(tx *sql.Tx, insertedUrbanVillage []repository.UrbanVillageModel) (tempCounter int32, err errorModel.ErrorModel) {
	var (
		totalJob    = 1000
		numJobs     = int(math.Ceil(float64(len(insertedUrbanVillage)) / float64(totalJob)))
		poolSize    = 5
		job         = make(chan repository.UrbanVillageModel, len(insertedUrbanVillage))
		results     = make(chan responseDataForThread, len(insertedUrbanVillage))
		ctx, cancel = context.WithCancel(context.Background())
		wg          sync.WaitGroup
	)

	if poolSize > numJobs {
		poolSize = numJobs
	}

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go input.doInsertBulkUrbanVillage(tx, &wg, ctx, job, results)
	}

	for j := 0; j < len(insertedUrbanVillage); j++ {
		job <- insertedUrbanVillage[j]
	}

	close(job)
	wg.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		for _ = range insertedUrbanVillage {
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
	log.Println(fmt.Sprintf(`[URBAN VILLAGE] Data insert raw: %d <===> Data insert berhasil di proses: %d`, len(insertedUrbanVillage), tempCounter))
	return
}

func (input taskSchedulerService) doInsertBulkUrbanVillage(tx *sql.Tx, wg *sync.WaitGroup, ctx context.Context, job <-chan repository.UrbanVillageModel, result chan<- responseDataForThread) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-job:
			var (
				err              errorModel.ErrorModel
				subDistrict      repository.SubDistrictModel
				tempResponseData responseDataForThread
				db               = serverconfig.ServerAttribute.DBConnection
			)

			if !ok {
				return
			}

			subDistrict, err = dao.SubDistrictDAO.GetSubDistrictByID(db, task.SubDistrictID.Int64, 0, true)
			if err.Error != nil {
				log.Println("Error: =====> ", err.Error.Error())
				tempResponseData = responseDataForThread{
					Err: err,
				}
			}

			if subDistrict.ID.Int64 > 0 {
				_, err = dao.UrbanVillageDAO.InsertUrbanVillage(tx, task)
				if err.Error != nil {
					log.Println("Error: ====> ", err.CausedBy.Error())
				}
				tempResponseData = responseDataForThread{
					Data: task.Name.String,
					Err:  err,
				}
			} else {
				tempResponseData = responseDataForThread{
					Data: task.Name.String,
					Err:  errorModel.GenerateNonErrorModel(),
				}
			}

			result <- tempResponseData
		}
	}
}

func (input taskSchedulerService) poolGetDataUrbanVillage(modelData []repository.UrbanVillageModel) (result []responseDataForThread) {
	var (
		totalJob    = 1000
		numJobs     = int(math.Ceil(float64(len(modelData)) / float64(totalJob)))
		poolSize    = 5
		jobs        = make(chan []repository.UrbanVillageModel, len(modelData))
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
		go input.doGetDataUrbanVillageForSync(&wg, jobs, results, ctx)
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
				break
			}
		}
	}()

	wg.Wait()
	ctx.Done()
	close(results)
	wg2.Wait()
	return
}

func (input taskSchedulerService) doGetDataUrbanVillageForSync(wg *sync.WaitGroup, jobs <-chan []repository.UrbanVillageModel,
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
			tempData, err := dao.UrbanVillageDAO.GetUpdatedDBUrbanVillage(serverconfig.ServerAttribute.DBConnection, task)
			tempResponseData := responseDataForThread{
				Data: tempData,
				Err:  err,
			}

			result <- tempResponseData
		}
	}
}
