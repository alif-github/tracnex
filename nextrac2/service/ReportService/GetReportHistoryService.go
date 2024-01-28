package ReportService

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

func (input reportService) GetListReportHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{
			"department_id",
			"created_at",
		}
		validOrderBy = []string{
			"department_id",
			"created_at",
		}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListReportHistoryValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListReportHistory(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input reportService) doGetListReportHistory(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	dbResult, err = dao.ReportHistoryDAO.GetListReportHistory(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input reportService) convertModelToResponseGetList(dbResult []interface{}) (result []out.ReportHistoryResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.ReportHistoryModel)
		result = append(result, out.ReportHistoryResponse{
			ID:                item.ID.Int64,
			Department:        item.DepartmentName.String,
			SuccessTicket:     item.SuccessTicket.String,
			PaymentDate:       item.CreatedAt.Time,
			PersonResponsible: item.CreatedName.String,
		})
	}

	return result
}
