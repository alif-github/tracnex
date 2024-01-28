package ProductGroupService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input productGroupService) GetInitiateProductGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListProductGroupValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateProductGroup(searchByParam, contextModel)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListProductGroupValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input productGroupService) doInitiateProductGroup(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var createdBy int64
	var scope map[string]interface{}

	createdBy = contextModel.LimitedByCreatedBy

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return 0, err
	}

	output, err = dao.ProductGroupDAO.GetCountProductGroup(serverconfig.ServerAttribute.DBConnection, searchByParam, createdBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	return
}

func (input productGroupService) doGetListProductGroup(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin bool) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}
	var createdBy int64 = 0
	var scope map[string]interface{}

	if isAdmin {
		scope = make(map[string]interface{})
		scope[constanta.ProductGroupDataScope] =[]interface{}{"all"}
	}else{
		//Add userID when have own permission
		createdBy = contextModel.LimitedByCreatedBy

		//Get scope
		scope, err =input.validateDataScope(contextModel)
		if err.Error != nil {
			return
		}
	}

	dbResult, err = dao.ProductGroupDAO.GetListProductGroup(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertToDTOOut(dbResult)
	return
}

func (input productGroupService) convertToDTOOut(dbResult []interface{}) (result []out.ProductGroupResponse) {
	for _, item := range dbResult {
		repo := item.(repository.ProductGroupModel)
		result = append(result, out.ProductGroupResponse{
			ID:               repo.ID.Int64,
			ProductGroupName: repo.ProductGroupName.String,
			CreatedAt:        repo.CreatedAt.Time,
			UpdatedAt:        repo.UpdatedAt.Time,
			UpdatedBy:        repo.UpdatedBy.Int64,
			UpdatedName:      repo.UpdatedName.String,
		})
	}
	return result
}