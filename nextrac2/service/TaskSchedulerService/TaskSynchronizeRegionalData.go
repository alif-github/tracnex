package TaskSchedulerService

import (
	"database/sql"
	"log"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/backgroundJobModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input taskSchedulerService) SchedulerSyncRegionalData(contextModel *applicationModel.ContextModel) (job repository.JobProcessModel, err errorModel.ErrorModel) {
	var (
		countryList []int64
		timeNow     = time.Now()
		task        = input.getSynchronizeRegionalDataTask()
		db          = serverconfig.ServerAttribute.DBConnection
		countryOnDB []master_data_response.CountryResponse
	)

	job = service.GetJobProcess(task, *contextModel, timeNow)

	//--- Get All Country
	countryOnDB, err = input.getCountryFromMDB()
	if err.Error != nil {
		return
	}

	//--- Set Country List ID
	for _, responseCountryOnDB := range countryOnDB {
		countryList = append(countryList, responseCountryOnDB.ID)
	}

	log.Printf("Start Scheduler Sync Regional Data")
	go input.ServiceWithChildBackgroundProcessWithoutConcurrent(db, false, []backgroundJobModel.ChildTask{
		input.getSynchronizeProvinceDataTask(),
		input.getSynchronizeDistrictDataTask(),
		input.getSynchronizeSubDistrictDataTask(countryList),
		input.getSynchronizeUrbanVillageDataTask(),
		input.getSynchronizePostalCodeDataTask(),
	}, job, *contextModel)
	return
}

func (input taskSchedulerService) getSynchronizeRegionalDataTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constanta.JobProcessSynchronizeGroup,
		Type:  constanta.JobProcessMasterDataType,
		Name:  constanta.JobProcessSynchronizeRegional,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          nil,
		},
		GetCountData: func(db *sql.DB, params []in.SearchByParam, b bool, i int64) (result int, err errorModel.ErrorModel) {
			return 5, errorModel.GenerateNonErrorModel()
		},
	}
}

func (input taskSchedulerService) synchronizeRegionalFromMasterData(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel) (err errorModel.ErrorModel) {

	//--- Scheduler function for Province
	//input.SchedulerSyncProvinceData(childJob)

	//--- Scheduler function for Postal Code
	//input.SchedulerSyncPostalCodeData(childJob)

	//childJob.Counter.Int32 += 4
	return
}
