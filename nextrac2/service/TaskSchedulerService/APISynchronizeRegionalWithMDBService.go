package TaskSchedulerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/util"
)

func (input taskSchedulerService) APISynchronizeRegionalData(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var jobProcessData repository.JobProcessModel
	jobProcessData, err = input.SchedulerSyncRegionalData(contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.ListJobProcessResponse{
		Level:     int(jobProcessData.Level.Int32),
		JobID:     jobProcessData.JobID.String,
		Group:     jobProcessData.Group.String,
		Type:      jobProcessData.Type.String,
		Name:      jobProcessData.Name.String,
		Counter:   int(jobProcessData.Counter.Int32),
		Total:     int(jobProcessData.Total.Int32),
		Status:    jobProcessData.Status.String,
		CreatedAt: jobProcessData.CreatedAt.Time,
		UpdatedAt: jobProcessData.UpdatedAt.Time,
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_SYNC_REGIONAL_DATA_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
