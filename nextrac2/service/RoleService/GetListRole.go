package RoleService

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
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input roleService) GetListRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListRoleValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	err = input.SetDefaultOrder(request, constanta.CreatedAtDesc, &inputStruct, input.ValidOrderBy)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListRole(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleService) InitiateGetListRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListRoleValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListRole(searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListRoleValidOperator,
		CountData:     countData.(int),
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input roleService) doGetListRole(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult  []interface{}
		createdBy int64 = 0
	)

	createdBy = contextModel.LimitedByCreatedBy
	dbResult, err = dao.RoleDAO.GetListRole(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input roleService) convertToListDTOOut(dbResult []interface{}) (result []out.ViewListRoleDTOOut) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.RoleModel)
		result = append(result, out.ViewListRoleDTOOut{
			ID:          repo.ID.Int64,
			RoleID:      repo.RoleID.String,
			Description: repo.Description.String,
			CreatedBy:   repo.CreatedBy.Int64,
			CreatedAt:   repo.CreatedAt.Time,
			UpdatedAt:   repo.UpdatedAt.Time,
			CreatedName: repo.CreatedName.String,
		})
	}

	return result
}

func (input roleService) doInitiateListRole(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = dao.RoleDAO.GetCountRole(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy)
	return
}
