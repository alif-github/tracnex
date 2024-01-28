package ClientMappingService

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
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
)

type uiGetListClientMappingService struct {
	service.AbstractService
	service.GetListData
}

var UIGetListClientMappingService = uiGetListClientMappingService{}.New()

func (input uiGetListClientMappingService) New() (output uiGetListClientMappingService) {
	output.FileName = "UIGetListClientMappingService.go"
	output.ValidSearchBy = []string{
		"client_id",
		"client_alias",
		"success_status_nexcloud",
	}
	output.ValidOrderBy = []string{
		"client_id",
		"socket_id",
		"client_type",
		"company_id",
		"branch_id",
		"client_alias",
	}
	output.ValidLimit = service.DefaultLimit
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "client_mapping.client_type_id",
		Count: "client_mapping.client_type_id",
	}
	return
}

func (input uiGetListClientMappingService) InitiateGetListClientMappings(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		countData     int
		searchByParam []in.SearchByParam
		scope         map[string]interface{}
		key           = "success_status_nexcloud"
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListClientMappingValidOperator)
	if err.Error != nil {
		return
	}

	scope, err = input.validateDataScopeClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	countData, err = dao.ClientMappingDAO.GetCountListClientMapping(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	for i := 0; i < len(input.ValidSearchBy); i++ {
		if input.ValidSearchBy[i] == key {
			input.ValidSearchBy = append(input.ValidSearchBy[:i], input.ValidSearchBy[i+1:]...)
			i = -1
		}
	}

	delete(applicationModel.GetListClientMappingValidOperator, key)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListClientMappingValidOperator,
		EnumData:      nil,
		CountData:     countData,
	}
	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("INITIATE_GET_LIST_CLIENT_MAPPING_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input uiGetListClientMappingService) GetClientMappings(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		clients       []interface{}
		scope         map[string]interface{}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListClientMappingValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	scope, err = input.validateDataScopeClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	clients, err = dao.ClientMappingDAO.GetListClientMapping(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertToDTO(clients)
	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_GET_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input uiGetListClientMappingService) convertToDTO(data []interface{}) (clients []out.ClientMappingForView) {
	for _, item := range data {
		client := item.(repository.ClientMappingForViewModel)
		clients = append(clients, out.ClientMappingForView{
			ID:                    client.ID.Int64,
			ClientID:              client.ClientID.String,
			SocketID:              client.SocketID.String,
			ClientType:            client.ClientType.String,
			CompanyID:             client.CompanyID.String,
			BranchID:              client.BranchID.String,
			Aliases:               client.Aliases.String,
			UpdatedAt:             client.UpdatedAt.Time,
			UpdatedBy:             client.UpdatedBy.Int64,
			CreatedAt:             client.CreatedAt.Time,
			CreatedBy:             client.CreatedBy.Int64,
			SuccessStatusAuth:     client.SuccessStatusAuth.Bool,
			SuccessStatusNexcloud: client.SuccessStatusNexcloud.Bool,
			SuccessStatusNexdrive: client.SuccessStatusNexdrive.Bool,
		})
	}
	return
}

func (input uiGetListClientMappingService) validateDataScopeClientMapping(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "validateDataScopeClientMapping"

	output = service.ValidateScope(contextModel, []string{
		constanta.ClientTypeDataScope,
	})

	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
