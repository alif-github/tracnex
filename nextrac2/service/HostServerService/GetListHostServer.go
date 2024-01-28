package HostServerService

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

func (input hostServerService) GetListHostServer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListHostServerValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListHostServer(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_HOST_SERVER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input hostServerService) doGetListHostServer(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output []out.HostServerResponse, err errorModel.ErrorModel) {
	funcName := "doGetListHostServer"
	var dbResult []interface{}
	var isPrimaryKey bool
	var isTableName bool
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "primary_key" {
			isPrimaryKey = true
		}
		if searchByParam[i].SearchKey == "table_name" {
			isTableName = true
		}
	}
	if isPrimaryKey {
		if !isTableName {
			err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.TableName)
			return
		}
	}
	dbResult, err = dao.HostServerDAO.GetListHostServer(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, false, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input hostServerService) convertToListDTOOut(dbResult []interface{}) (result []out.HostServerResponse) {
	for i := 0; i < len(dbResult); i++ {
		repo := dbResult[i].(repository.HostServerModel)
		result = append(result, out.HostServerResponse{
			ID:        repo.ID.Int64,
			HostName:  repo.HostName.String,
			CreatedBy: repo.CreatedBy.Int64,
			UpdatedAt: repo.UpdatedAt.Time,
		})
	}
	return result
}

func (input hostServerService) InitiateGetListHostServer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData int

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListHostServerValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateHostServer(searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_HOST_SERVER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListHostServerValidOperator,
		CountData:     countData,
	}
	return
}

func (input hostServerService) doInitiateHostServer(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output int, err errorModel.ErrorModel) {
	output, err = dao.HostServerDAO.CountHostServer(serverconfig.ServerAttribute.DBConnection, searchByParam, false, contextModel.LimitedByCreatedBy)
	return
}
