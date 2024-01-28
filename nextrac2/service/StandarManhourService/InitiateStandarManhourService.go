package StandarManhourService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input standarManhourService) InitiateStandarManhour(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
		db            = serverconfig.ServerAttribute.DBConnection
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListStandarManhourValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.StandarManhourDAO.GetCountStandarManhour(db, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	if countData == nil {
		countData = 0
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListStandarManhourValidOperator,
		CountData:     countData.(int),
	}

	return
}
