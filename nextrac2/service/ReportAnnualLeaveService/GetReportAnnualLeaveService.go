package ReportAnnualLeaveService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input reportAnnualLeaveService) GetListJobReportAnnualLeaveService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListJobReportAnnualLeaveValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListJobAnnualLeaveService(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input reportAnnualLeaveService) doGetListJobAnnualLeaveService(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult  []interface{}
		db        = serverconfig.ServerAttribute.DBConnection
		createdBy = contextModel.AuthAccessTokenModel.ResourceUserID
	)

	dbResult, err = dao.FileUploadDAO.GetListFileUploadAndJobProcess(db, inputStruct, searchByParam, createdBy, dao.JobProcessDAO.TableName)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input reportAnnualLeaveService) convertModelToResponseGetList(dbResult []interface{}) (result []out.GetListFileUploadJobProcess) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.GetListFileUploadJobProcess)
		result = append(result, out.GetListFileUploadJobProcess{
			JobID:       item.JobID.String,
			Status:      item.Status.String,
			Progress:    item.Progress.Float64,
			FileUrl:     item.FileUrl.String,
			CreatedName: item.CreatedName.String,
			CreatedAt:   item.CreatedAt.Time,
		})
	}

	return result
}
