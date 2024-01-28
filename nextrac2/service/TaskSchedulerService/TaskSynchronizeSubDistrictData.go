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
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"sync"
	"time"
)

func (input taskSchedulerService) getSynchronizeSubDistrictDataTask(countryListID []int64) backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constanta.JobProcessSynchronizeGroup,
		Type:  constanta.JobProcessMasterDataType,
		Name:  constanta.JobProcessSynchronizeSubDistrict,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          countryListID,
		},
		GetCountData: input.doGetCountSubDistrict,
		DoJobWithCtx: input.synchronizeSubDistrictFromMasterData,
	}
}

func (input taskSchedulerService) synchronizeSubDistrictFromMasterData(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel, contextModel applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		token            = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, "", constanta.Issue, constanta.DefaultApplicationsLanguage)
		lastSync         time.Time
		subDistrictOnMDB []master_data_response.SubDistrictResponse
	)

	log.Printf("Start Scheduler Sync Regional Data Sub District")
	if childJob.Total.Int32 == 0 {
		service.LogMessage("No data is processed", 200)
		return
	}

	lastSync, err = dao.SubDistrictDAO.GetSubDistrictLastSync(db)
	if err.Error != nil {
		return
	}

	subDistrictOnMDB, err = master_data_dao.GetListAllSubDistrictFromMasterData(master_data_request.SubDistrictRequest{
		UpdatedAtStart: lastSync,
	}, &contextModel, token)

	if err.Error != nil {
		return
	}

	counter, err := input.ServiceWithDataAuditPreparedByService("synchronizeSubDistrictFromMasterData", subDistrictOnMDB, &contextModel, input.doSynchronizeSubDistrictFromMasterData, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	childJob.Counter.Int32 += counter.(int32)
	return
}

func (input taskSchedulerService) doSynchronizeSubDistrictFromMasterData(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct                     = inputStructInterface.([]master_data_response.SubDistrictResponse)
		newSubDistrict, subDistrictOnDB []repository.SubDistrictModel
		counter                         int32
		counterInsert                   int32
	)

	for _, response := range inputStruct {
		newSubDistrict = append(newSubDistrict, repository.SubDistrictModel{
			DistrictID:       sql.NullInt64{Int64: response.DistrictID},
			DistrictName:     sql.NullString{String: response.DistrictName},
			MDBSubDistrictID: sql.NullInt64{Int64: response.ID},
			Code:             sql.NullString{String: response.Code},
			Name:             sql.NullString{String: response.Name},
			Status:           sql.NullString{String: response.Status},
			CreatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			LastSync:         sql.NullTime{Time: timeNow},
		})
	}

	subDistrictOnDB, err = dao.SubDistrictDAO.GetUpdatedMDBSubDistrict(serverconfig.ServerAttribute.DBConnection, newSubDistrict)
	if err.Error != nil {
		return
	}

	subDistrictMustUpdate, subDistrictMustInsert := input.getUpdateAndInsertDataSubDistrict(newSubDistrict, subDistrictOnDB)

	// Update Sub District
	for _, model := range subDistrictMustUpdate {
		err = dao.SubDistrictDAO.UpdateDataSubDistrict(tx, model)
		if err.Error != nil {
			return
		}
		counter += 1
	}

	// Insert Sub District
	//if len(subDistrictMustInsert) > 0 {
	//	limitData := constanta.TotalDataProductLicensePerChannel
	//	totalPage := int(math.Ceil(float64(len(subDistrictMustInsert)) / float64(limitData)))
	//	for i := 0; i < totalPage; i++ {
	//		var insertedID []int64
	//		if i == 0 {
	//			insertedID, err = dao.SubDistrictDAO.InsertBulkSubDistrict(tx, subDistrictMustInsert[:limitData])
	//		} else if i == totalPage-1 {
	//			insertedID, err = dao.SubDistrictDAO.InsertBulkSubDistrict(tx, subDistrictMustInsert[limitData*i:])
	//		} else {
	//			insertedID, err = dao.SubDistrictDAO.InsertBulkSubDistrict(tx, subDistrictMustInsert[limitData*i:limitData*(i+1)])
	//		}
	//
	//		if err.Error != nil {
	//			return
	//		}
	//		counter += int32(len(insertedID))
	//	}
	//}

	// Insert Sub District
	if len(subDistrictMustInsert) > 0 {
		counterInsert, err = input.insertBulkSubDistrict(tx, subDistrictMustInsert)
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

func (input taskSchedulerService) insertBulkSubDistrict(tx *sql.Tx, insertedSubDistrict []repository.SubDistrictModel) (tempCounter int32, err errorModel.ErrorModel) {
	var (
		totalJob    = 1000
		numJobs     = int(math.Ceil(float64(len(insertedSubDistrict)) / float64(totalJob)))
		poolSize    = 5
		job         = make(chan repository.SubDistrictModel, len(insertedSubDistrict))
		results     = make(chan responseDataForThread, len(insertedSubDistrict))
		ctx, cancel = context.WithCancel(context.Background())
		wg          sync.WaitGroup
	)

	if poolSize > numJobs {
		poolSize = numJobs
	}

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go input.doInsertBulkSubDistrict(tx, &wg, ctx, job, results)
	}

	for j := 0; j < len(insertedSubDistrict); j++ {
		job <- insertedSubDistrict[j]
	}

	close(job)
	wg.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		for _ = range insertedSubDistrict {
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
	log.Println(fmt.Sprintf(`[SUB DISTRICT] Data insert raw: %d <===> Data insert berhasil di proses: %d`, len(insertedSubDistrict), tempCounter))
	return
}

func (input taskSchedulerService) doInsertBulkSubDistrict(tx *sql.Tx, wg *sync.WaitGroup, ctx context.Context, job <-chan repository.SubDistrictModel, result chan<- responseDataForThread) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-job:
			var (
				err              errorModel.ErrorModel
				idDistrict       int64
				tempResponseData responseDataForThread
				db               = serverconfig.ServerAttribute.DBConnection
			)

			if !ok {
				return
			}

			idDistrict, err = dao.DistrictDAO.GetOnlyDistrictIDByMdbID(db, repository.DistrictModel{MdbDistrictID: sql.NullInt64{Int64: task.DistrictID.Int64}}, true)
			if err.Error != nil {
				log.Println("Error: =====> ", err.Error.Error())
				tempResponseData = responseDataForThread{
					Err: err,
				}
			}

			if idDistrict > 0 {
				_, err = dao.SubDistrictDAO.InsertSubDistrict(tx, task)
				if err.Error != nil {
					log.Println("Error: ====> ", err.Error.Error())
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

func (input taskSchedulerService) getUpdateAndInsertDataSubDistrict(updatedSubDistrict, updatedSubDistrictOnDB []repository.SubDistrictModel) (updatedSubDistrictData, insertedSubDistrictData []repository.SubDistrictModel) {
	for i := 0; i < len(updatedSubDistrict); i++ {
		var isInsertData = true
		for j := 0; j < len(updatedSubDistrictOnDB); j++ {
			if updatedSubDistrict[i].MDBSubDistrictID.Int64 == updatedSubDistrictOnDB[j].MDBSubDistrictID.Int64 {
				updatedSubDistrict[i].ID = updatedSubDistrictOnDB[j].ID
				isInsertData = false
				updatedSubDistrictData = append(updatedSubDistrictData, updatedSubDistrict[i])
				if j < len(updatedSubDistrictOnDB)-1 {
					updatedSubDistrictOnDB = append(updatedSubDistrictOnDB[:j], updatedSubDistrictOnDB[j+1:]...)
				} else {
					updatedSubDistrictOnDB = append(updatedSubDistrictOnDB[:j])
				}
				break
			}
		}

		if isInsertData {
			insertedSubDistrictData = append(insertedSubDistrictData, updatedSubDistrict[i])
		}
	}

	return
}

func (input taskSchedulerService) doGetCountSubDistrict(db *sql.DB, params []in.SearchByParam, b bool, i int64) (result int, err errorModel.ErrorModel) {
	lastSync, err := dao.SubDistrictDAO.GetSubDistrictLastSync(db)
	if err.Error != nil {
		return
	}

	tempResult, err := master_data_dao.CountAllSubDistrictFromMasterData(master_data_request.SubDistrictRequest{
		UpdatedAtStart: lastSync,
	}, &applicationModel.ContextModel{}, "")
	result = tempResult
	return
}
