package HostServerService

import (
	"database/sql"
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
	"strconv"
)

func (input viewDetailHostServerService) ViewDetailHostServer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListDataWithID(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListRunningCornValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewDetailHostServer(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_DETAIL_HOST_SERVER", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input viewDetailHostServerService) doViewDetailHostServer(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output out.ViewHostServerResponse, err errorModel.ErrorModel) {
	funcName := "doViewDetailHostServer"

	hostServer := repository.HostServerModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}

	hostServer, err = dao.HostServerDAO.GetHostServerDb(serverconfig.ServerAttribute.DBConnection, hostServer)
	if err.Error != nil {
		return
	}

	if hostServer.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.HostID)
		return
	}
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "name" {
			searchByParam[i].SearchKey = "cs." + searchByParam[i].SearchKey
		}
	}

	list, err := dao.HostServerDAO.ViewDetailHostServer(serverconfig.ServerAttribute.DBConnection, inputStruct, hostServer, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertToDetailHostServerDTOOut(hostServer, list)
	return
}

func (input viewDetailHostServerService) convertToDetailHostServerDTOOut(hostServer repository.HostServerModel, list []interface{}) (result out.ViewHostServerResponse) {
	result.Host = out.Host{
		ID:   hostServer.ID.Int64,
		Name: hostServer.HostName.String,
		Url:  hostServer.HostURL.String,
	}
	for i := 0; i < len(list); i++ {
		scheduler := list[i].(repository.ListScheduler)
		var status, _ = strconv.ParseBool(scheduler.Status.String)
		result.ListScheduler = append(result.ListScheduler, out.ListScheduler{
			Name:      scheduler.Name.String,
			Cron:      scheduler.Cron.String,
			RunStatus: status,
		})
	}
	return result
}

func (input viewDetailHostServerService) InitiateGetListRunningCron(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var countData int
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err := input.ReadAndValidateGetCountDataWithID(request, input.ValidSearchBy, applicationModel.GetListRunningCornValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListRunningCron(inputStruct, searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_DETAIL_HOST_SERVER", contextModel.AuthAccessTokenModel.Locale),
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListRunningCornValidOperator,
		CountData:     countData,
	}
	return
}

func (input viewDetailHostServerService) doInitiateListRunningCron(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output int, err errorModel.ErrorModel) {
	funcName := "doInitiateListRunningCron"
	hostServer := repository.HostServerModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}
	hostServer, err = dao.HostServerDAO.GetHostServerDb(serverconfig.ServerAttribute.DBConnection, hostServer)
	if err.Error != nil {
		return
	}

	if hostServer.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.HostID)
		return
	}
	hostServer.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	output, err = dao.HostServerDAO.CountRunningCron(serverconfig.ServerAttribute.DBConnection, hostServer, searchByParam)
	return
}
