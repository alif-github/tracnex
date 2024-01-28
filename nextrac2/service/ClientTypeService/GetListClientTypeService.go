package ClientTypeService

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
	"nexsoft.co.id/nextrac2/util"
)

func (input clientTypeService) GetListClientTypeAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListClientTypeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListClientType(inputStruct, searchByParam, contextModel, true)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_CLIENT_TYPE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeService) GetListClientType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListClientTypeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListClientType(inputStruct, searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_CLIENT_TYPE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeService) doGetListClientType(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin bool) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}
	var scope map[string]interface{}
	var useCreatedBy int64

	if isAdmin {
		scope = make(map[string]interface{})
		scope[constanta.ClientTypeDataScope] = []interface{}{"all"}
		useCreatedBy = 0
	} else {
		scope, err = input.ValidateDataScope(contextModel, constanta.ClientTypeDataScope)
		if err.Error != nil {
			return
		}
		useCreatedBy = contextModel.LimitedByCreatedBy
	}

	dbResult, err = dao.ClientTypeDAO.GetListClientType(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, scope, input.MappingScopeDB, useCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertToListClientType(dbResult)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeService) convertToListClientType(dbResult []interface{}) (result []out.ListClientTypeResponse) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.ClientTypeModel)
		result = append(result, out.ListClientTypeResponse{
			ID:                 repo.ID.Int64,
			ClientType:         repo.ClientType.String,
			Description:        repo.Description.String,
			ParentClientTypeID: repo.ParentClientTypeID.Int64,
			CreatedAt:          repo.CreatedAt.Time,
			UpdatedName:        repo.UpdatedName.String,
			UpdatedAt:          repo.UpdatedAt.Time,
		})
	}

	return result
}
