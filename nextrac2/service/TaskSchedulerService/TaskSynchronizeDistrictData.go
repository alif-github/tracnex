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

func (input taskSchedulerService) SchedulerSyncDistrictData(parentJob *repository.JobProcessModel) {
	var (
		timeNow      = time.Now()
		contextModel applicationModel.ContextModel
		task         backgroundJobModel.ChildTask
		job          repository.JobProcessModel
	)

	task = input.getSynchronizeDistrictDataTask()
	job = service.GetJobProcess(task, contextModel, timeNow)

	input.ServiceWithBackgroundProcess(serverconfig.ServerAttribute.DBConnection, false, *parentJob, job, task, contextModel)

	parentJob.Counter.Int32 += 1
}

func (input taskSchedulerService) getSynchronizeDistrictDataTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constanta.JobProcessSynchronizeGroup,
		Type:  constanta.JobProcessMasterDataType,
		Name:  constanta.JobProcessSynchronizeDistrict,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          nil,
		},
		GetCountData: func(db *sql.DB, params []in.SearchByParam, b bool, i int64) (result int, err errorModel.ErrorModel) {

			lastSync, err := dao.DistrictDAO.GetDistrictLastSync(db)
			if err.Error != nil {
				return
			}

			tempResult, err := master_data_dao.CountAllDistrictFromMasterData(master_data_request.ProvinceRequest{
				UpdatedAtStart: lastSync,
			}, &applicationModel.ContextModel{})

			result = int(tempResult)
			return
		},
		DoJobWithCtx: input.synchronizeDistrictFromMasterData,
	}
}

func (input taskSchedulerService) synchronizeDistrictFromMasterData(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel, ctx applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName      = "synchronizeDistrictFromMasterData"
		districtOnMDB []master_data_response.DistrictResponse
		lastSync      time.Time
	)

	log.Printf("Start Scheduler Sync District Data")
	if childJob.Total.Int32 == 0 {
		service.LogMessage("No data is processed", 200)
		return
	}

	// get last sync
	lastSync, err = dao.DistrictDAO.GetDistrictLastSync(db)
	if err.Error != nil {
		return
	}

	// Get All Province on MDB
	districtOnMDB, err = master_data_dao.GetListAllDistrictFromMasterData(master_data_request.DistrictRequest{
		UpdatedAtStart: lastSync,
	}, &applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	if len(districtOnMDB) != int(childJob.Total.Int32) {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.District)
		return
	}

	counter, err := input.ServiceWithDataAuditPreparedByService("synchronizeDistrictFromMasterData", districtOnMDB, &ctx, input.doSynchronizeDistrict, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	childJob.Counter.Int32 += counter.(int32)
	return
}

func (input taskSchedulerService) doSynchronizeDistrict(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		updatedDistrict, updatedDistrictOnDB, insertedDistrict []repository.DistrictModel
		districtOnMDB                                          = inputStructInterface.([]master_data_response.DistrictResponse)
		counter                                                int32
		counterInsert                                          int32
	)

	// Memisahkan updatedDistrict dan notUpdatedDistrict
	for _, district := range districtOnMDB {
		updatedDistrict = append(updatedDistrict, repository.DistrictModel{
			ProvinceID:    sql.NullInt64{Int64: district.ProvinceID},
			MdbDistrictID: sql.NullInt64{Int64: district.ID},
			Code:          sql.NullString{String: district.Code},
			Name:          sql.NullString{String: district.Name},
			Status:        sql.NullString{String: district.Status},
			LastSync:      sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedAt:     sql.NullTime{Time: timeNow},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
		})
	}

	// get updatedProvince On Local DB
	updatedDistrictOnDB, err = dao.DistrictDAO.GetUpdatedMDBDistrict(serverconfig.ServerAttribute.DBConnection, updatedDistrict)
	if err.Error != nil {
		return
	}

	// Get New Data and Apply New Changes FROM MDB
	updatedDistrict, insertedDistrict = input.getUpdateAndInsertDataDistrict(updatedDistrict, updatedDistrictOnDB)

	// Update Province
	for _, district := range updatedDistrict {
		err = dao.DistrictDAO.UpdateDataDistrict(tx, district)
		if err.Error != nil {
			return
		}
		counter += 1
	}

	// Insert Province
	if len(insertedDistrict) > 0 {
		counterInsert, err = input.insertBulkDistrict(tx, insertedDistrict)
		counter += counterInsert
		if err.Error != nil {
			return
		}
	}

	defer func() {
		output = counter
	}()

	return
}

func (input taskSchedulerService) insertBulkDistrict(tx *sql.Tx, insertedDistrict []repository.DistrictModel) (tempCounter int32, err errorModel.ErrorModel) {
	var (
		totalJob    = 1000
		numJobs     = int(math.Ceil(float64(len(insertedDistrict)) / float64(totalJob)))
		poolSize    = 5
		job         = make(chan repository.DistrictModel, len(insertedDistrict))
		results     = make(chan responseDataForThread, len(insertedDistrict))
		ctx, cancel = context.WithCancel(context.Background())
		wg          sync.WaitGroup
	)

	if poolSize > numJobs {
		poolSize = numJobs
	}

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go input.doInsertBulkDistrict(tx, &wg, ctx, job, results)
	}

	for j := 0; j < len(insertedDistrict); j++ {
		job <- insertedDistrict[j]
	}

	close(job)
	wg.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		for _ = range insertedDistrict {
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
	log.Println(fmt.Sprintf(`[DISTRICT] Data raw insert: %d <===> Data insert berhasil di proses: %d`, len(insertedDistrict), tempCounter))
	return
}

func (input taskSchedulerService) doInsertBulkDistrict(tx *sql.Tx, wg *sync.WaitGroup, ctx context.Context, job <-chan repository.DistrictModel, result chan<- responseDataForThread) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-job:
			var (
				err              errorModel.ErrorModel
				idProvince, id   int64
				tempResponseData responseDataForThread
				db               = serverconfig.ServerAttribute.DBConnection
			)

			if !ok {
				return
			}

			idProvince, err = dao.ProvinceDAO.GetProvinceIDByMdbID(db, repository.ProvinceModel{MDBProvinceID: sql.NullInt64{Int64: task.ProvinceID.Int64}}, true)
			if err.Error != nil {
				log.Println("Error: =====> ", err.Error.Error())
				tempResponseData = responseDataForThread{
					Err: err,
				}
			}

			if idProvince > 0 {
				id, err = dao.DistrictDAO.InsertDistrict(tx, task)
				if err.Error != nil {
					log.Println("Error: ====> ", err.Error.Error())
				} else {
					_, err = input.GenerateDataScope(tx, id, dao.DistrictDAO.TableName, constanta.DistrictDataScope, constanta.SystemID, constanta.SystemClient, time.Now())
					if err.Error != nil {
						log.Println("Error: ====> ", err.Error.Error())
					}
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

func (input taskSchedulerService) getUpdateAndInsertDataDistrict(updatedDistrict, updatedDistrictOnDB []repository.DistrictModel) (updatedDistrictData, insertedDistrictData []repository.DistrictModel) {
	for i := 0; i < len(updatedDistrict); i++ {
		var isInsertData = true
		for j := 0; j < len(updatedDistrictOnDB); j++ {
			if updatedDistrict[i].MdbDistrictID.Int64 == updatedDistrictOnDB[j].MdbDistrictID.Int64 {
				updatedDistrict[i].ID = updatedDistrictOnDB[j].ID
				isInsertData = false
				updatedDistrictData = append(updatedDistrictData, updatedDistrict[i])
				if j < len(updatedDistrictOnDB)-1 {
					updatedDistrictOnDB = append(updatedDistrictOnDB[:j], updatedDistrictOnDB[j+1:]...)
				} else {
					updatedDistrictOnDB = append(updatedDistrictOnDB[:j])
				}
				break
			}
		}

		if isInsertData {
			insertedDistrictData = append(insertedDistrictData, updatedDistrict[i])
		}
	}

	return
}

func (input taskSchedulerService) getDistrictFromMDB(childJob *repository.JobProcessModel, countryList []int64) (result []master_data_response.DistrictResponse, err errorModel.ErrorModel) {
	var (
		wg          = &sync.WaitGroup{}
		totalPage   = int(math.Ceil(float64(childJob.Total.Int32) / float64(constanta.TotalDataProductLicensePerChannel)))
		totalThread = 5
		jobs        = make(chan master_data_request.DistrictRequest, totalPage)
		tempResult  chan master_data_response.DistrictResponse
	)

	if totalPage < totalThread {
		totalThread = totalPage
	}

	// Set Job
	for i := 1; i <= totalPage; i++ {
		jobs <- master_data_request.DistrictRequest{
			AbstractDTO: in.AbstractDTO{
				Page:  i,
				Limit: constanta.TotalDataProductLicensePerChannel,
			},
			CountryIDList: countryList,
		}
	}

	// Do Go Routine
	for i := 0; i < totalThread; i++ {
		wg.Add(1)
		go input.doGetDistrictFromMDB(jobs, tempResult, wg)
	}

	wg.Wait()

	// Collect Result
	for i := 0; i < len(tempResult); i++ {
		result = append(result, <-tempResult)
	}

	close(jobs)
	close(tempResult)
	return
}

func (input taskSchedulerService) doGetDistrictFromMDB(inputRequest <-chan master_data_request.DistrictRequest, response chan<- master_data_response.DistrictResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for request := range inputRequest {
		var (
			districtOnDB []master_data_response.DistrictResponse
			err          errorModel.ErrorModel
		)

		districtOnDB, err = master_data_dao.GetListForSyncDistrictFromMasterData(request, &applicationModel.ContextModel{})
		if err.Error != nil {
			return
		}

		for _, districtResponse := range districtOnDB {
			response <- districtResponse
		}
	}

	return
}
