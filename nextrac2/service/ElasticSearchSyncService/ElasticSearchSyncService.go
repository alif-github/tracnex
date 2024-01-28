package ElasticSearchSyncService

import (
	"context"
	"database/sql"
	"encoding/json"
	"gopkg.in/olivere/elastic.v7"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/backgroundJobModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

type elasticSearchSyncService struct {
	service.AbstractService
}

var ElasticSearchSyncService = elasticSearchSyncService{}.New()

func (input elasticSearchSyncService) New() (output elasticSearchSyncService) {
	output.FileName = "ElasticSearchSyncService.go"
	return
}

func (input elasticSearchSyncService) DoSyncDBAndElastic(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ElasticSyncRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel)
	if err.Error != nil {
		return
	}

	job := input.DoSyncPartialDBToElastic(*contextModel, inputStruct)

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: util.GenerateConstantaI18n("SUCCESS_SYNC_DATA", contextModel.AuthAccessTokenModel.Locale, nil),
	}

	output.Data.Content = job.JobID.String

	return
}

func (input elasticSearchSyncService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.ElasticSyncRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	err = inputStruct.ValidateElasticSyncRequest()

	return
}

func (input elasticSearchSyncService) clearIndex(client *elastic.Client, indexName string) {
	ctx := context.Background()
	_, _ = client.DeleteIndex(indexName).Do(ctx)
}

func (input elasticSearchSyncService) addPrefixIndexName(indexName string) string {
	return config.ApplicationConfiguration.GetServerResourceID() + "." + indexName
}

func (input elasticSearchSyncService) syncElasticData(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel, getListData func(*sql.DB, in.GetListDataDTO, []in.SearchByParam, bool, int64) ([]interface{}, errorModel.ErrorModel), doJob func(interface{}) errorModel.ErrorModel) errorModel.ErrorModel {
	var err errorModel.ErrorModel

	updateDBEvery := int(float32(childJob.Total.Int32) * 1.0 / 100.0)
	if updateDBEvery < 10 {
		updateDBEvery = 10
	}

	loop := (int(childJob.Total.Int32) / updateDBEvery) + 1

	for i := 0; i < loop; i++ {
		var result []interface{}

		result, err = getListData(db, in.GetListDataDTO{
			AbstractDTO: in.AbstractDTO{
				Page:    i + 1,
				Limit:   updateDBEvery,
				OrderBy: "id",
			}}, nil, false, 0)
		if err.Error != nil {
			return err
		}

		for j := 0; j < len(result); j++ {
			err = doJob(result[j])
			if err.Error != nil {
				return err
			}
		}

		childJob.Counter.Int32 += int32(len(result))
		childJob.UpdatedAt.Time = time.Now()
		err = dao.JobProcessDAO.UpdateJobProcessCounter(db, *childJob)
		if err.Error != nil {
			return err
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input elasticSearchSyncService) DoSyncAllDBToElastic(contextModel applicationModel.ContextModel) repository.JobProcessModel {
	var listTask []backgroundJobModel.ChildTask

	//listTask = append(listTask, input.GetSyncElasticBankChildTask())
	job := service.GetJobProcess(backgroundJobModel.ChildTask{
		Group: constanta.JobProcessSynchronizeGroup,
		Type:  constanta.JobProcessElasticType,
		Name:  constanta.JobProcessSynchronizeGroup + constanta.JobProcessElasticType,
	}, contextModel, time.Now())

	job.Level.Int32 = 1
	go input.ServiceWithChildBackgroundProcess(serverconfig.ServerAttribute.DBConnection, true, listTask, job, contextModel)

	return job
}

func (input elasticSearchSyncService) DoSyncPartialDBToElastic(contextModel applicationModel.ContextModel, requestParam in.ElasticSyncRequest) repository.JobProcessModel {
	if requestParam.All {
		return input.DoSyncAllDBToElastic(contextModel)
	} else {
		var listTask []backgroundJobModel.ChildTask
		if requestParam.Bank {
			//listTask = append(listTask, input.GetSyncElasticBankChildTask())
		}

		job := service.GetJobProcess(backgroundJobModel.ChildTask{
			Group: constanta.JobProcessSynchronizeGroup,
			Type:  constanta.JobProcessElasticType,
			Name:  constanta.JobProcessSynchronizeGroup + constanta.JobProcessElasticType,
		}, contextModel, time.Now())
		job.Level.Int32 = 1

		go input.ServiceWithChildBackgroundProcess(serverconfig.ServerAttribute.DBConnection, false, listTask, job, contextModel)
		return job
	}
}
