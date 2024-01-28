package StandarManhourService

import (
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input standarManhourService) GetListStandarManhour(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListStandarManhourValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListStandarManhour(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input standarManhourService) doGetListStandarManhour(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	dbResult, err = input.StandarManhourDAO.GetListStandarManhour(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input standarManhourService) convertModelToResponseGetList(dbResult []interface{}) (result []out.StandarManhourResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.StandarManhourModel)
		result = append(result, out.StandarManhourResponse{
			ID:         item.ID.Int64,
			Case:       item.Case.String,
			Department: item.Department.String,
			Manhour:    fmt.Sprintf(`%.4f`, item.Manhour.Float64),
			UpdatedAt:  item.UpdatedAt.Time,
		})
	}

	return result
}
