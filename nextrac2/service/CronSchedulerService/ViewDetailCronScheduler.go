package CronSchedulerService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input cronSchedulerService) ViewCronSchedulerDetail(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewCronSchedulerDetail(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_CRON_SCHEDULER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input cronSchedulerService) doViewCronSchedulerDetail(inputStruct in.CronSchedulerRequest, contextModel *applicationModel.ContextModel) (result out.ViewCronSchedulerResponse, err errorModel.ErrorModel) {
	funcName := "doViewDataTypeDetail"
	cronSchedulerModel := repository.CRONSchedulerModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}

	cronSchedulerModel, err = dao.CronSchedulerDAO.ViewCronScheduler(serverconfig.ServerAttribute.DBConnection, cronSchedulerModel)
	if err.Error != nil {
		return
	}

	if cronSchedulerModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}
	result = reformatRepositoryToDTOOut(cronSchedulerModel)
	return
}

func reformatRepositoryToDTOOut(cronSchedulerModel repository.CRONSchedulerModel) out.ViewCronSchedulerResponse {
	temp := out.ViewCronSchedulerResponse{
		ID:        cronSchedulerModel.ID.Int64,
		Name:      cronSchedulerModel.Name.String,
		RunType:   cronSchedulerModel.RunType.String,
		Cron:      cronSchedulerModel.CRON.String,
		CreatedBy: cronSchedulerModel.CreatedBy.Int64,
		UpdatedAt: cronSchedulerModel.UpdatedAt.Time,
		Status:    cronSchedulerModel.Status.Bool,
	}
	return temp
}

func (input cronSchedulerService) validateView(inputStruct *in.CronSchedulerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewCronScheduler()
}
