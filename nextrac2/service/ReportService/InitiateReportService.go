package ReportService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input reportService) InitiateReport(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct      in.GetListDataDTO
		searchByParam    []in.SearchByParam
		isMandatoryExist bool
		countData        = 0
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListReportValidOperator)
	if err.Error != nil {
		return
	}

	isMandatoryExist, _, err = input.validateAddParam(request, &inputStruct, &searchByParam)
	if err.Error != nil {
		return
	}

	if isMandatoryExist {
		countData, err = input.doInitiateReport(searchByParam, contextModel)
		if err.Error != nil {
			return
		}
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListReportValidOperator,
		CountData:     countData,
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) doInitiateReport(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (total int, err errorModel.ErrorModel) {
	var (
		scope map[string]interface{}
		db    = serverconfig.ServerAttribute.DBConnection
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	return input.ReportDAO.GetCountReport(db, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
}
