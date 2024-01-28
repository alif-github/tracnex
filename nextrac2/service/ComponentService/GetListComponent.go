package ComponentService

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

func (input componentService) InitiateComponent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListComponentValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateComponent(searchByParam, contextModel)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListComponentValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input componentService) doInitiateComponent(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var createdBy int64
	output = 0
	createdBy = contextModel.LimitedByCreatedBy

	output, err = dao.ComponentDAO.GetCountComponent(serverconfig.ServerAttribute.DBConnection, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	return
}

func (input componentService) GetListComponent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListComponentValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListComponent(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input componentService) doGetListComponent(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}

	dbResult, err = dao.ComponentDAO.GetLisComponent(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input componentService) convertModelToResponseGetList(dbResult []interface{}) (result []out.ComponentResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.ComponentModel)
		result = append(result, out.ComponentResponse{
			ID:            item.ID.Int64,
			ComponentName: item.ComponentName.String,
			CreatedAt:     item.CreatedAt.Time,
			UpdatedAt:     item.UpdatedAt.Time,
			UpdatedBy:     item.UpdatedBy.Int64,
			UpdatedName:   item.UpdatedName.String,
		})
	}

	return result
}