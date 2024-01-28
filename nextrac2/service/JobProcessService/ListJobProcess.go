package JobProcessService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input jobProcessService) GetListJobProcess(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListJobProcessValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListJobProcess(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_LIST_JOB_PROCESS_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input jobProcessService) InitiateGetListJobProcess(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData int

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListJobProcessValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListJobProcess(searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_JOB_PROCESS_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	enumData := make(map[string][]string)
	enumData["status"] = []string{
		"ONPROGRESS",
		"ONPROGRESS-ERROR",
		"OK",
		"ERROR"}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListJobProcessValidOperator,
		EnumData:      enumData,
		CountData:     countData,
	}

	return
}

func (input jobProcessService) doGetListJobProcess(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output []out.ListJobProcessResponse, err errorModel.ErrorModel) {
	var dbResult []interface{}

	dbResult, err = dao.JobProcessDAO.GetListJobProcess(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, false, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input jobProcessService) doInitiateListJobProcess(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output int, err errorModel.ErrorModel) {
	output, err = dao.JobProcessDAO.GetCountJobProcess(serverconfig.ServerAttribute.DBConnection, searchByParam, false, contextModel.LimitedByCreatedBy)
	return
}

func (input jobProcessService) convertToListDTOOut(dbResult []interface{}) (result []out.ListJobProcessResponse) {
	for i := 0; i < len(dbResult); i++ {
		repo := dbResult[i].(repository.ListJobProcessModel)
		result = append(result, out.ListJobProcessResponse{
			Level:     int(repo.Level.Int32),
			JobID:     repo.JobID.String,
			Group:     repo.Group.String,
			Type:      repo.Type.String,
			Name:      repo.Name.String,
			Counter:   int(repo.Counter.Int32),
			Total:     int(repo.Total.Int32),
			Status:    repo.Status.String,
			CreatedAt: repo.CreatedAt.Time,
			UpdatedAt: repo.UpdatedAt.Time,
			Duration:  repo.UpdatedAt.Time.Sub(repo.CreatedAt.Time).Seconds(),
		})
	}
	return result
}