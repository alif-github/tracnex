package JobProcessService

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

func (input jobProcessService) ViewJobProcess(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewJobProcess(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_JOB_PROCESS_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input jobProcessService) doViewJobProcess(inputStruct in.JobProcessRequest, contextModel *applicationModel.ContextModel) (result out.ViewJobProcessResponse, err errorModel.ErrorModel) {
	funcName := "doViewJobProcess"
	var viewJobProcess repository.ViewJobProcessModel

	jobProcessModel := repository.JobProcessModel{
		JobID: sql.NullString{String: inputStruct.JobID},
	}

	jobProcessModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	viewJobProcess, err = dao.JobProcessDAO.ViewJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcessModel)
	if err.Error != nil {
		return
	}

	if viewJobProcess.JobID.String == "" {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.JobID)
		return
	}

	result = reformatRepositoryToDTOOut(viewJobProcess)

	return
}

func reformatRepositoryToDTOOut(jobProcessModel repository.ViewJobProcessModel) out.ViewJobProcessResponse {
	result := out.ViewJobProcessResponse{
		Level:          int(jobProcessModel.Level.Int32),
		JobID:          jobProcessModel.JobID.String,
		Group:          jobProcessModel.Group.String,
		Type:           jobProcessModel.Type.String,
		Name:           jobProcessModel.Name.String,
		UrlIn:          jobProcessModel.UrlIn.String,
		FileNameIn:     jobProcessModel.FileNameIn.String,
		ContentDataOut: jobProcessModel.ContentDataOut.String,
		Counter:        int(jobProcessModel.Counter.Int32),
		Total:          int(jobProcessModel.Total.Int32),
		Status:         jobProcessModel.Status.String,
		CreatedAt:      jobProcessModel.CreatedAt.Time,
		UpdatedAt:      jobProcessModel.UpdatedAt.Time,
		Duration:       jobProcessModel.UpdatedAt.Time.Sub(jobProcessModel.CreatedAt.Time).Seconds(),
		Percentage:     float64(jobProcessModel.Counter.Int32) / float64(jobProcessModel.Total.Int32),
	}

	result.ChildJobProcess = jobProcessModel.ChildJobProcess
	for i := 0; i < len(result.ChildJobProcess); i++ {
		createdAt, _ := in.TimeDBStrToTime(result.ChildJobProcess[i].CreatedAt, "")
		updatedAt, _ := in.TimeDBStrToTime(result.ChildJobProcess[i].UpdatedAt, "")
		result.ChildJobProcess[i].Duration = updatedAt.Sub(createdAt).Seconds()
	}
	return result
}

func (input jobProcessService) validateView(inputStruct *in.JobProcessRequest) errorModel.ErrorModel {
	return inputStruct.ViewDetailJobProcess()
}
