package PKCEClientMappingService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input pkceClientMappingService) InitiateGetListPKCEClientMapping(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		inputStruct   in.GetListDataDTO
		count         int
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListPKCEClientMappingValidOperator)
	if err.Error != nil {
		return
	}

	count, err = input.doInitiateGetListPKCEClientMapping(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListPKCEClientMappingValidOperator,
		CountData:     count,
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("INITIATE_GET_LIST_PKCE_CLIENT_MAPPING_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingService) GetListPKCEClientMapping(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListPKCEClientMappingValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListPKCEClientMapping(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_GET_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingService) doInitiateGetListPKCEClientMapping(_ in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (count int, err errorModel.ErrorModel) {
	var scope map[string]interface{}

	scope, err = input.validateDataScopePKCEClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	count, err = dao.PKCEClientMappingDAO.GetCountPKCEListClientMapping(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	return
}

func (input pkceClientMappingService) doGetListPKCEClientMapping(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		resultGetList []interface{}
		scope         map[string]interface{}
	)

	scope, err = input.validateDataScopePKCEClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	resultGetList, err = dao.PKCEClientMappingDAO.GetListPKCEClientMapping(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertToDTO(resultGetList)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingService) convertToDTO(data []interface{}) (clients []out.PKCEClientMappingForList) {
	for _, item := range data {
		client := item.(repository.ViewPKCEClientMappingModel)
		clients = append(clients, out.PKCEClientMappingForList{
			ID:             client.ID.Int64,
			ParentClientID: client.ParentClientID.String,
			FirstName:      client.FirstName.String,
			LastName:       client.LastName.String,
			Username:       client.Username.String,
			ClientType:     client.ClientType.String,
			CompanyID:      client.CompanyID.String,
			BranchID:       client.BranchID.String,
			ClientAlias:    client.ClientAlias.String,
			UpdatedAt:      client.UpdatedAt.Time,
			UpdatedBy:      client.UpdatedBy.Int64,
			CreatedBy:      client.CreatedBy.Int64,
			CreatedAt:      client.CreatedAt.Time,
		})
	}

	return
}
