package BacklogService

import (
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
)

func (input backlogService) GetListParentBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListParentBacklogValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListParentBacklog(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) doGetListParentBacklog(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		scope    map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.BacklogDAO.GetListParentBacklog(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input backlogService) convertModelToResponseGetList(dbResult []interface{}) (result []out.ParentBacklogResponse) {
	var standarTimeHour = 8.00

	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.BacklogModel)
		totalMandaysStringValue := fmt.Sprintf(`%.4f`, item.TotalMandays.Float64/standarTimeHour)
		totalMandaysFloatValue, errs := strconv.ParseFloat(totalMandaysStringValue, 64)
		if errs != nil {
			return
		}

		result = append(result, out.ParentBacklogResponse{
			Sprint:       item.Sprint.String,
			TotalMandays: totalMandaysFloatValue,
		})
	}

	return result
}

func (input backlogService) InitiateGetListParentBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListParentBacklogValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateParentBacklog(searchByParam, contextModel)
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
		ValidOperator: applicationModel.GetListParentBacklogValidOperator,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) doInitiateParentBacklog(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var scope map[string]interface{}
	for i, _ := range searchByParam {
		if searchByParam[i].SearchKey == "sprint" {
			searchByParam[i].SearchKey = "sub_query.sprint_time"
		}
	}

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.BacklogDAO.GetCountParentBacklog(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	return
}
