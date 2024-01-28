package TaskSchedulerService

import (
	"database/sql"
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

func (input taskSchedulerService) SchedulerSyncProvinceData(parentJob *repository.JobProcessModel) {
	timeNow := time.Now()
	task := input.getSynchronizeProvinceDataTask()
	job := service.GetJobProcess(task, applicationModel.ContextModel{}, timeNow)
	input.ServiceWithBackgroundProcess(serverconfig.ServerAttribute.DBConnection, false, *parentJob, job, task, applicationModel.ContextModel{})
}

func (input taskSchedulerService) getSynchronizeProvinceDataTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constanta.JobProcessSynchronizeGroup,
		Type:  constanta.JobProcessMasterDataType,
		Name:  constanta.JobProcessSynchronizeProvince,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          nil,
		},
		GetCountData: input.doGetCountProvince,
		DoJobWithCtx: input.synchronizeProvinceFromMasterData,
	}
}

func (input taskSchedulerService) synchronizeProvinceFromMasterData(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel, contextModel applicationModel.ContextModel) (err errorModel.ErrorModel) {
	// Synchronize Province
	log.Printf("Start Scheduler Sync Regional Data Province")

	var (
		provinceOnMDB []master_data_response.ProvinceResponse
		lastSync      time.Time
	)

	if childJob.Total.Int32 == 0 {
		service.LogMessage("No data is processed", 200)
		return
	}

	// todo get last sync
	lastSync, err = dao.ProvinceDAO.GetProvinceLastSync(db)
	if err.Error != nil {
		return
	}

	// Get All Province on MDB
	provinceOnMDB, err = master_data_dao.GetListAllProvinceFromMasterData(master_data_request.ProvinceRequest{
		UpdatedAtStart: lastSync,
	}, &applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	counter, err := input.ServiceWithDataAuditPreparedByService("synchronizeProvinceFromMasterData", provinceOnMDB, &contextModel, input.doSynchronizeProvince, func(i interface{}, model applicationModel.ContextModel) {

	})
	if err.Error != nil {
		return
	}

	childJob.Counter.Int32 += counter.(int32)
	return
}

func (input taskSchedulerService) doSynchronizeProvince(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		provinceOnMDB                                          = inputStructInterface.([]master_data_response.ProvinceResponse)
		updatedProvince, updatedProvinceOnDB, insertedProvince []repository.ProvinceModel
		insertedID                                             []int64
		counter                                                int32
	)

	// Memisahkan updatedProvince dan notUpdatedProvince
	for _, province := range provinceOnMDB {
		updatedProvince = append(updatedProvince, repository.ProvinceModel{
			CountryID:     sql.NullInt64{Int64: province.CountryID},
			MDBProvinceID: sql.NullInt64{Int64: province.ID},
			Code:          sql.NullString{String: province.Code},
			Name:          sql.NullString{String: province.Name},
			Status:        sql.NullString{String: province.Status},
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
	updatedProvinceOnDB, err = dao.ProvinceDAO.GetUpdatedMDBProvince(serverconfig.ServerAttribute.DBConnection, updatedProvince)
	if err.Error != nil {
		return
	}

	// Get New Data and Apply New Changes FROM MDB
	updatedProvince, insertedProvince = input.getUpdateAndInsertDataProvince(updatedProvince, updatedProvinceOnDB)

	// Update Province
	for _, province := range updatedProvince {
		err = dao.ProvinceDAO.UpdateDataProvince(tx, province)
		if err.Error != nil {
			continue
		}
		counter += 1
	}

	// Insert Province
	if len(insertedProvince) > 0 {
		insertedID, err = dao.ProvinceDAO.InsertBulkProvince(tx, insertedProvince)
		if err.Error != nil {
			return
		}

		counter += int32(len(insertedID))
		for _, id := range insertedID {
			_, err = input.GenerateDataScope(tx, id, dao.ProvinceDAO.TableName, constanta.ProvinceDataScope, constanta.SystemID, constanta.SystemClient, timeNow)
			if err.Error != nil {
				return
			}
		}
	}

	defer func() {
		output = counter
	}()
	return
}

func (input taskSchedulerService) getUpdateAndInsertDataProvince(updatedProvince, updatedProvinceOnDB []repository.ProvinceModel) (updatedProvinceData, insertedProvinceData []repository.ProvinceModel) {
	for i := 0; i < len(updatedProvince); i++ {
		var isInsertData = true
		for j := 0; j < len(updatedProvinceOnDB); j++ {
			if updatedProvince[i].MDBProvinceID.Int64 == updatedProvinceOnDB[j].MDBProvinceID.Int64 {
				updatedProvince[i].ID = updatedProvinceOnDB[j].ID
				isInsertData = false
				updatedProvinceData = append(updatedProvinceData, updatedProvince[i])
				if j < len(updatedProvinceOnDB)-1 {
					updatedProvinceOnDB = append(updatedProvinceOnDB[:j], updatedProvinceOnDB[j+1:]...)
				} else {
					updatedProvinceOnDB = append(updatedProvinceOnDB[:j])
				}
				break
			}
		}

		if isInsertData {
			insertedProvinceData = append(insertedProvinceData, updatedProvince[i])
		}
	}

	return
}

func (input taskSchedulerService) getProvinceFromMDB(childJob *repository.JobProcessModel, countryList []int64, lastSync time.Time) (result []master_data_response.ProvinceResponse, err errorModel.ErrorModel) {
	var (
		waitGroup   sync.WaitGroup
		totalPage   = int(math.Ceil(float64(childJob.Total.Int32) / float64(constanta.TotalDataProductLicensePerChannel)))
		totalThread = 5
		jobs        master_data_request.ProvinceRequest
		tempResult  = make(chan master_data_response.ProvinceResponse, childJob.Total.Int32)
	)

	if totalPage < totalThread {
		totalThread = totalPage
	}

	// Do Go Routine
	for i := 1; i <= totalPage; i++ {
		waitGroup.Add(1)
		jobs = master_data_request.ProvinceRequest{
			AbstractDTO: in.AbstractDTO{
				Page:  i,
				Limit: constanta.TotalDataProductLicensePerChannel,
			},
			CountryIDList:  countryList,
			UpdatedAtStart: lastSync,
		}
		go input.doGetProvinceFromMDB(jobs, tempResult, &waitGroup)
	}

	for i := 0; i < int(childJob.Total.Int32); i++ {
		result = append(result, <-tempResult)
	}

	waitGroup.Wait()
	//fmt.Println("After Wait Group")

	if len(tempResult) >= 1 {
		close(tempResult)
	}
	return
}

func (input taskSchedulerService) doGetProvinceFromMDB(inputRequest master_data_request.ProvinceRequest, response chan<- master_data_response.ProvinceResponse, wg *sync.WaitGroup) {
	var (
		provinceOnDB []master_data_response.ProvinceResponse
		err          errorModel.ErrorModel
	)
	provinceOnDB, err = master_data_dao.GetListForSyncProvinceFromMasterData(inputRequest, &applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	for _, provinceResponse := range provinceOnDB {
		response <- provinceResponse
	}

	defer func() {
		wg.Done()
	}()
	return
}

func (input taskSchedulerService) getCountryFromMDB() (result []master_data_response.CountryResponse, err errorModel.ErrorModel) {
	var waitGroup sync.WaitGroup
	totalData, err := master_data_dao.CountAllCountryCountryFromMasterData(&applicationModel.ContextModel{})
	if err.Error != nil {
		return
	}

	tempResult := make(chan master_data_response.CountryResponse, totalData)

	if totalData < 1 {
		service.LogMessage("No data Country on MDB", 200)
		return
	}

	totalThread := int(math.Ceil(float64(totalData) / float64(constanta.TotalDataProductLicensePerChannel)))

	for i := 1; i <= totalThread; i++ {
		waitGroup.Add(1)
		go input.doGetCountryFromMDB(i, tempResult, &waitGroup)
	}

	for i := 0; i < int(totalData); i++ {
		result = append(result, <-tempResult)
	}

	waitGroup.Wait()
	if len(tempResult) > 0 {
		close(tempResult)
	}
	return
}

func (input taskSchedulerService) doGetCountryFromMDB(page int, tempResult chan<- master_data_response.CountryResponse, wg *sync.WaitGroup) {
	countryONMDB, err := master_data_dao.GetListCountryCountryFromMasterData(master_data_request.CountryRequest{
		AbstractDTO: in.AbstractDTO{
			Page:  page,
			Limit: constanta.TotalDataProductLicensePerChannel,
		},
		UpdatedAtStart: time.Time{},
	}, &applicationModel.ContextModel{})

	if err.Error != nil {
		return
	}

	for _, response := range countryONMDB {
		tempResult <- response
	}

	defer func() {
		wg.Done()
	}()
}

func (input taskSchedulerService) doGetCountProvince(db *sql.DB, params []in.SearchByParam, b bool, i int64) (result int, err errorModel.ErrorModel) {
	lastSync, err := dao.ProvinceDAO.GetProvinceLastSync(db)
	if err.Error != nil {
		return
	}
	tempResult, err := master_data_dao.CountAllProvinceFromMasterData(master_data_request.ProvinceRequest{
		UpdatedAtStart: lastSync,
	}, &applicationModel.ContextModel{})
	result = int(tempResult)
	return
}
