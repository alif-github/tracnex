package ParameterService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input parameterService) GetListParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListParameterValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListParameter(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterService) InitiateGetListParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData int

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListParameterValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListParameter(searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListParameterValidOperator,
		EnumData:      nil,
		CountData:     countData,
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterService) doGetListParameter(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output []out.ViewParameterDTOOut, err errorModel.ErrorModel) {
	var dbResult []interface{}
	var createdBy int64 = 0

	if strings.Contains(contextModel.PermissionHave, "own") {
		createdBy = contextModel.AuthAccessTokenModel.ResourceUserID
	}

	dbResult, err = dao.ParameterDAO.GetListParameter(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input parameterService) doInitiateListParameter(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output int, err errorModel.ErrorModel) {
	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(contextModel)
	if !isOnlyHaveOwnAccess {
		userID = 0
	}

	output, err = dao.ParameterDAO.GetCountParameter(serverconfig.ServerAttribute.DBConnection, searchByParam, userID)
	return
}

func (input parameterService) convertToListDTOOut(dbResult []interface{}) (result []out.ViewParameterDTOOut) {
	for i := 0; i < len(dbResult); i++ {
		repo := dbResult[i].(repository.ParameterModel)
		result = append(result, out.ViewParameterDTOOut{
			ID:         repo.ID.Int64,
			Permission: repo.Permission.String,
			Name:       repo.Name.String,
			CreatedBy:  repo.CreatedBy.Int64,
			UpdatedAt:  repo.UpdatedAt.Time,
		})
	}
	return result
}
