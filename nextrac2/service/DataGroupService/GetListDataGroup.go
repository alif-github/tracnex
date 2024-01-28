package DataGroupService

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

func (input dataGroupService) GetListDataGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListDataGroupValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	err = input.SetDefaultOrder(request, constanta.CreatedAtDesc, &inputStruct, input.ValidOrderBy)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListDataGroup(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) InitiateGetListDataGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListDataGroupValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListDataGroup(searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListDataGroupValidOperator,
		EnumData:      nil,
		CountData:     countData.(int),
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) doGetListDataGroup(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}
	var createdBy int64 = 0

	createdBy = contextModel.LimitedByCreatedBy

	dbResult, err = dao.DataGroupDAO.GetListDataGroup(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input dataGroupService) doInitiateListDataGroup(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = dao.DataGroupDAO.GetCountDataGroup(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy)
	return
}

func (input dataGroupService) convertToListDTOOut(dbResult []interface{}) (result []out.ViewListDataGroupDTOOut) {
	for i := 0; i < len(dbResult); i++ {
		repo := dbResult[i].(repository.DataGroupModel)
		result = append(result, out.ViewListDataGroupDTOOut{
			ID:          repo.ID.Int64,
			GroupID:     repo.GroupID.String,
			CreatedBy:   repo.CreatedBy.Int64,
			CreatedAt:   repo.CreatedAt.Time,
			UpdatedAt:   repo.UpdatedAt.Time,
			UpdatedBy:   repo.UpdatedBy.Int64,
			Description: repo.Description.String,
			CreatedName: repo.CreatedName.String,
		})
	}
	return result
}
