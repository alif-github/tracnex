package Task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/HitAddExternalResource"
	util2 "nexsoft.co.id/nextrac2/util"
	"sync"
	"time"
)

type schAddResourceNexcloud struct {
	service.AbstractService
	AbstractScheduledTask
}

var SchAddResourceNexcloud = schAddResourceNexcloud{}.New()

func (schAddResourceNexcloud) New() (output schAddResourceNexcloud) {
	output.RunType = "scheduler.add_resource_nexcloud"
	return
}

func (input schAddResourceNexcloud) Start() {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.StartTask(input.RunType, input.schAddResourceNexcloud)
	}
}

func (input schAddResourceNexcloud) StartMain() {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.schAddResourceNexcloud()
	}
}

func (input schAddResourceNexcloud) logEmptySchAddResourceNexcloud() {
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Message = "Empty client registration log"
	logModel.Status = 200
}

func (input schAddResourceNexcloud) schAddResourceNexcloud() {
	fileName := "SchAddResourceNexcloud.go"
	funcName := "schAddResourceNexcloud"

	serverConfig := serverconfig.ServerAttribute.DBConnection
	var err errorModel.ErrorModel
	var jobProcessAddResource repository.JobProcessModel
	var jobProcessUpdateLog repository.JobProcessModel
	var contextModel applicationModel.ContextModel

	defer func() {
		if err.Error != nil {
			logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
			logModel.Message = "Scheduler: add resource nexcloud failed"
			logModel.Status = 500
			util.LogError(logModel.ToLoggerObject())
			err = dao.JobProcessDAO.UpdateErrorJobProcessWithCounter(serverConfig, jobProcessAddResource)
			if err.Error != nil {
				logModel = applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
				logModel.Message = err.Error.Error()
				if err.CausedBy != nil {
					logModel.Message = err.CausedBy.Error()
				}
				logModel.Status = 500
				util.LogError(logModel.ToLoggerObject())
			}

			err = dao.JobProcessDAO.UpdateErrorJobProcessWithCounter(serverConfig, jobProcessUpdateLog)
			if err.Error != nil {
				logModel = applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
				logModel.Message = err.Error.Error()
				if err.CausedBy != nil {
					logModel.Message = err.CausedBy.Error()
				}
				logModel.Status = 500
				util.LogError(logModel.ToLoggerObject())
			}
		} else {
			jobProcessAddResource.UpdatedAt.Time = time.Now()
			jobProcessUpdateLog.UpdatedAt.Time = time.Now()

			err = dao.JobProcessDAO.UpdateJobProcessCounter(serverConfig, jobProcessAddResource)
			if err.Error != nil {
				return
			}

			err = dao.JobProcessDAO.UpdateJobProcessCounter(serverConfig, jobProcessUpdateLog)
			if err.Error != nil {
				return
			}
		}
	}()

	jobProcessAddResource = repository.GenerateAddResourceNexcloudUserJobProcessModel(1, time.Now(), constanta.SystemID)
	jobProcessUpdateLog = repository.GenerateUpdateLogAfterAddResourceUserJobProcessModel(2, jobProcessAddResource.JobID.String, time.Now(), constanta.SystemID)

	clientLog, err := dao.ClientRegistrationLogDAO.GetClientRegistrationLogForScheduler(serverConfig, repository.ParamClientRegistrationLogModel{ClientTypeID: sql.NullInt64{Int64: 1}})
	if err.Error != nil {
		return
	}

	//---------- Id data who want to be process is empty, then scheduler return exit from function
	if len(clientLog) < 1 {
		input.logEmptySchAddResourceNexcloud()
		return
	}

	jobProcessAddResource.Total.Int32 = int32(len(clientLog))
	jobProcessUpdateLog.Total.Int32 = int32(len(clientLog))

	err = dao.JobProcessDAO.InsertJobProcess(serverConfig, jobProcessAddResource)
	if err.Error != nil {
		return
	}

	err = dao.JobProcessDAO.InsertJobProcess(serverConfig, jobProcessUpdateLog)
	if err.Error != nil {
		return
	}

	//-------------------- Main Service --------------------
	counterSuccessNexcloud, counterSuccessUpdateLog, err := input.addResourceToNexcloud(clientLog, &contextModel)
	jobs := len(clientLog)
	if counterSuccessNexcloud < int32(jobs) || counterSuccessUpdateLog < int32(jobs) || err.Error != nil {

		jobProcessUpdateLog.Counter.Int32 = counterSuccessUpdateLog
		jobProcessAddResource.Counter.Int32 = counterSuccessNexcloud

		errorS := errors.New("data processing has not been completed")
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errorS)
		return
	}

	jobProcessUpdateLog.Counter.Int32 = counterSuccessUpdateLog
	jobProcessAddResource.Counter.Int32 = counterSuccessNexcloud

	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Message = "Scheduler add resource nexcloud success"
	logModel.Status = 200

	util.LogInfo(logModel.ToLoggerObject())
}

func (input schAddResourceNexcloud) addResourceToNexcloud(clientLog []repository.ClientRegistrationLogModel,
	contextModel *applicationModel.ContextModel) (counterSuccessNexcloud int32, counterSuccessUpdateLog int32, err errorModel.ErrorModel) {

	result := make(chan in.ResultWorkerClientLogDTO, len(clientLog))
	jobs := make(chan repository.ClientRegistrationLogModel, len(clientLog))

	var wg sync.WaitGroup

	//---------- Pool workers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go input.workerProcessCore(jobs, result, contextModel, &wg)
	}

	for _, item := range clientLog {
		var clientLogItem repository.ClientRegistrationLogModel
		clientLogItem = item
		jobs <- clientLogItem
	}

	close(jobs)
	wg.Wait()

	//---------- Update registration log
	for i := 0; i < len(clientLog); i++ {
		data := <- result
		if data.Errors.Error == nil {
			counterSuccessNexcloud++
		}

		err = input.updateRegistrationLogWithAudit(data, contextModel)
		if err.Error != nil {
			continue
		} else {
			counterSuccessUpdateLog++
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input schAddResourceNexcloud) workerProcessCore(jobs <- chan repository.ClientRegistrationLogModel, result chan <- in.ResultWorkerClientLogDTO, contextModel *applicationModel.ContextModel, wg *sync.WaitGroup) {
	defer wg.Done()

	for jobsElm := range jobs {
		firstName, err := input.doGetFirstNameUser(jobsElm)
		if err.Error != nil {
			result <- in.ResultWorkerClientLogDTO {
				ID:			jobsElm.ID.Int64,
				ClientID: 	jobsElm.ClientID.String,
				Resource:	jobsElm.Resource.String,
				Errors:		err,
			}
			continue
		}

		err = input.doHitAddResourceToNexcloud(jobsElm, firstName, contextModel)
		if err.Error != nil {
			result <- in.ResultWorkerClientLogDTO {
				ID:			jobsElm.ID.Int64,
				ClientID: 	jobsElm.ClientID.String,
				Resource:	jobsElm.Resource.String,
				Errors:		err,
			}
			continue
		}

		err = errorModel.GenerateNonErrorModel()
		result <- in.ResultWorkerClientLogDTO {
			ID:			jobsElm.ID.Int64,
			ClientID: 	jobsElm.ClientID.String,
			Resource:	jobsElm.Resource.String,
			Errors:		err,
		}
	}
}

func (input schAddResourceNexcloud) doHitAddResourceToNexcloud(clientIDElm repository.ClientRegistrationLogModel, firstName string, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	err = HitAddExternalResource.HitAddExternalResource.AddResourceNexcloud(in.AddResourceNexcloud {
		FirstName: 	firstName,
		LastName: 	constanta.Nexdistribution,
		ClientID: 	clientIDElm.ClientID.String,
	}, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input schAddResourceNexcloud) updateRegistrationLogWithAudit(clientLogData in.ResultWorkerClientLogDTO, contextModel *applicationModel.ContextModel) (errorS errorModel.ErrorModel) {
	var preparedLog in.PreparedRepositoryClientRegisterLog
	var detail string
	result := clientLogData

	type attributeRequest struct {
		ClientTypeID	int64	`json:"client_type_id"`
		ClientID		string	`json:"client_id"`
	}

	attributeRequestStruct := attributeRequest {
		ClientTypeID: 	1,
		ClientID: 		result.ClientID,
	}

	resultByte, _ := json.Marshal(attributeRequestStruct)

	if result.Errors.Error == nil {
		preparedLog = in.PreparedRepositoryClientRegisterLog {
			Status: 	true,
			Resource: 	result.Resource + " " + constanta.NexCloudResourceID,
			Code: 		util2.GenerateConstantaI18n("SUCCESS", constanta.IndonesianLanguage, nil),
			Message:	util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.AddResourceExternalBundle, "SUCCESS_ADD_RESOURCE_NEXCLOUD_MESSAGE", constanta.IndonesianLanguage, nil),
		}
	} else {
		if len(result.Errors.AdditionalInformation) > 0 {
			detail = result.Errors.AdditionalInformation[0]
		}
		preparedLog = in.PreparedRepositoryClientRegisterLog {
			Status: 	false,
			Resource: 	result.Resource,
			Code: 		result.Errors.Error.Error(),
			Message:	util2.GenerateI18NErrorMessage(result.Errors, constanta.IndonesianLanguage),
			Detail: 	detail,
		}
	}

	preparedLog.AttributeRequest = string(resultByte)
	preparedLog.ID = result.ID
	_, errorS = input.ServiceWithDataAuditPreparedByService("updateRegistrationLogWithAudit", preparedLog, contextModel, input.doUpdateToRegistrationLogs, func(_ interface{}, _ applicationModel.ContextModel) {})
	if errorS.Error != nil {
		return
	}
	errorS = errorModel.GenerateNonErrorModel()
	return
}

func (input schAddResourceNexcloud) doUpdateToRegistrationLogs(tx *sql.Tx, inputStruct interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {

	data := inputStruct.(in.PreparedRepositoryClientRegisterLog)

	clientLogModel := repository.ClientRegistrationLogModel {
		ID: 					sql.NullInt64{Int64: data.ID},
		SuccessStatusNexcloud: 	sql.NullBool{Bool: data.Status},
		MessageNexcloud: 		sql.NullString{String: data.Message},
		Details: 				sql.NullString{String: data.Detail},
		Code: 					sql.NullString{String: data.Code},
		RequestTimeStamp: 		sql.NullTime{Time: timeNow},
		UpdatedBy: 				sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: 			sql.NullString{String: constanta.SystemClient},
		UpdatedAt: 				sql.NullTime{Time: timeNow},
		Resource: 				sql.NullString{String: data.Resource},
		AttributeRequest: 		sql.NullString{String: data.AttributeRequest},
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ClientRegistrationLogDAO.TableName, data.ID, 0)...)

	err = dao.ClientRegistrationLogDAO.UpdateRegistrationLogForAddResourceNexcloudScheduler(tx, clientLogModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input schAddResourceNexcloud) doGetFirstNameUser(clientLogData repository.ClientRegistrationLogModel) (firstName string, err errorModel.ErrorModel) {
	fileName := "SchAddResourceNexcloud.go"
	funcName := "doGetFirstNameUser"

	userModel, err := dao.UserDAO.GetIdAndFirstNameUser(serverconfig.ServerAttribute.DBConnection, repository.UserModel {
		ClientID: 		sql.NullString{String: clientLogData.ClientID.String},
	})

	firstName = userModel.FirstName.String

	if err.Error != nil {
		return
	}

	if firstName == "" {
		err = errorModel.GenerateDataNotFound(fileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}